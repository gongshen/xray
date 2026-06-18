package business

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"

	"github.com/valyala/fasthttp"
)

const (
	defaultXrayLogDir        = "/var/log/xray"
	analysisLocationName     = "Asia/Shanghai"
	maxAnalysisRangeInMinute = 120
)

type UserTrafficAnalysisRequest struct {
	UserTag string `json:"user_tag"`
	Date    string `json:"date"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

type UserTrafficTarget struct {
	Target string `json:"target"`
	Count  int    `json:"count"`
}

type UserTrafficMinuteRow struct {
	Minute  string              `json:"minute"`
	Events  int                 `json:"events"`
	Down    uint64              `json:"down"`
	Up      uint64              `json:"up"`
	Total   uint64              `json:"total"`
	Targets []UserTrafficTarget `json:"targets"`
}

type UserTrafficAnalysisResponse struct {
	UserTag          string                 `json:"user_tag"`
	StartTime        string                 `json:"start_time"`
	EndTime          string                 `json:"end_time"`
	AccessLogMatched int                    `json:"access_log_matched"`
	Rows             []UserTrafficMinuteRow `json:"rows"`
}

type userTrafficAnalysisRange struct {
	UserTag   string
	Date      string
	StartTime time.Time
	EndTime   time.Time
}

type accessMinuteSummary struct {
	Matched int
	Targets map[string]map[string]int
}

func AnalyzeUserTrafficHandler(ctx *fasthttp.RequestCtx) {
	if LocalStore == nil {
		ctx.Error("traffic store is not initialized", fasthttp.StatusServiceUnavailable)
		return
	}
	req := UserTrafficAnalysisRequest{
		UserTag: string(ctx.QueryArgs().Peek("user_tag")),
		Date:    string(ctx.QueryArgs().Peek("date")),
		Start:   string(ctx.QueryArgs().Peek("start")),
		End:     string(ctx.QueryArgs().Peek("end")),
	}
	resp, err := AnalyzeUserTraffic(LocalStore, defaultXrayLogDir, req)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	data, err := json.Marshal(resp)
	if err != nil {
		ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		return
	}
	ctx.Success("application/json", data)
}

func AnalyzeUserTraffic(store *TrafficStore, logDir string, req UserTrafficAnalysisRequest) (*UserTrafficAnalysisResponse, error) {
	analysisRange, err := ParseUserTrafficAnalysisRange(req)
	if err != nil {
		return nil, err
	}
	if store == nil || store.db == nil {
		return nil, fmt.Errorf("traffic store is not initialized")
	}
	trafficRows, err := store.userTrafficMinuteRows(analysisRange)
	if err != nil {
		return nil, err
	}
	accessSummary, err := summarizeAccessLogsByMinute(logDir, analysisRange)
	if err != nil {
		return nil, err
	}
	rows := mergeTrafficAndAccessRows(trafficRows, accessSummary.Targets)
	return &UserTrafficAnalysisResponse{
		UserTag:          analysisRange.UserTag,
		StartTime:        analysisRange.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:          analysisRange.EndTime.Format("2006-01-02 15:04:05"),
		AccessLogMatched: accessSummary.Matched,
		Rows:             rows,
	}, nil
}

func ParseUserTrafficAnalysisRange(req UserTrafficAnalysisRequest) (userTrafficAnalysisRange, error) {
	userTag := strings.TrimSpace(req.UserTag)
	if userTag == "" {
		return userTrafficAnalysisRange{}, fmt.Errorf("user_tag is required")
	}
	location, err := time.LoadLocation(analysisLocationName)
	if err != nil {
		return userTrafficAnalysisRange{}, err
	}
	date, err := normalizeAnalysisDate(req.Date)
	if err != nil {
		return userTrafficAnalysisRange{}, err
	}
	startClock, err := normalizeAnalysisClock(req.Start, "00")
	if err != nil {
		return userTrafficAnalysisRange{}, err
	}
	endClock, err := normalizeAnalysisClock(req.End, "59")
	if err != nil {
		return userTrafficAnalysisRange{}, err
	}
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", date+" "+startClock, location)
	if err != nil {
		return userTrafficAnalysisRange{}, err
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", date+" "+endClock, location)
	if err != nil {
		return userTrafficAnalysisRange{}, err
	}
	if endTime.Before(startTime) {
		return userTrafficAnalysisRange{}, fmt.Errorf("end time cannot be earlier than start time")
	}
	startMinute := startTime.Hour()*60 + startTime.Minute()
	endMinute := endTime.Hour()*60 + endTime.Minute()
	if endMinute-startMinute > maxAnalysisRangeInMinute {
		return userTrafficAnalysisRange{}, fmt.Errorf("analysis range cannot exceed 2 hours")
	}
	return userTrafficAnalysisRange{
		UserTag:   userTag,
		Date:      date,
		StartTime: startTime,
		EndTime:   endTime,
	}, nil
}

func normalizeAnalysisDate(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) == 8 && isDigits(value) {
		value = value[:4] + "-" + value[4:6] + "-" + value[6:8]
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil || parsed.Format("2006-01-02") != value {
		return "", fmt.Errorf("invalid date, expected YYYYMMDD")
	}
	return value, nil
}

func normalizeAnalysisClock(value string, second string) (string, error) {
	parts := strings.Split(strings.TrimSpace(value), ":")
	if len(parts) != 2 || !isDigits(parts[0]) || !isDigits(parts[1]) {
		return "", fmt.Errorf("invalid time, expected H:MM")
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", err
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 || len(parts[1]) != 2 {
		return "", fmt.Errorf("invalid time, expected H:MM")
	}
	return fmt.Sprintf("%02d:%02d:%s", hour, minute, second), nil
}

func isDigits(value string) bool {
	if value == "" {
		return false
	}
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func (store *TrafficStore) userTrafficMinuteRows(analysisRange userTrafficAnalysisRange) ([]UserTrafficMinuteRow, error) {
	rows, err := store.db.Query(
		`SELECT (collected_at / 60) * 60 AS bucket, COUNT(*), COALESCE(SUM(down), 0), COALESCE(SUM(up), 0), COALESCE(SUM(down + up), 0)
		 FROM traffic_event
		 WHERE tag = ? AND collected_at >= ? AND collected_at <= ?
		 GROUP BY bucket
		 ORDER BY bucket ASC`,
		analysisRange.UserTag, analysisRange.StartTime.Unix(), analysisRange.EndTime.Unix(),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	location := analysisRange.StartTime.Location()
	out := make([]UserTrafficMinuteRow, 0)
	for rows.Next() {
		var bucket int64
		var events int
		var down, up, total int64
		if err := rows.Scan(&bucket, &events, &down, &up, &total); err != nil {
			return nil, err
		}
		out = append(out, UserTrafficMinuteRow{
			Minute:  time.Unix(bucket, 0).In(location).Format("2006-01-02 15:04"),
			Events:  events,
			Down:    uint64(down),
			Up:      uint64(up),
			Total:   uint64(total),
			Targets: []UserTrafficTarget{},
		})
	}
	return out, rows.Err()
}

func summarizeAccessLogsByMinute(logDir string, analysisRange userTrafficAnalysisRange) (accessMinuteSummary, error) {
	out := accessMinuteSummary{Targets: map[string]map[string]int{}}
	for _, path := range accessLogFiles(logDir, analysisRange.Date) {
		if err := scanAccessLogFile(path, analysisRange, &out); err != nil {
			return out, err
		}
	}
	return out, nil
}

func accessLogFiles(logDir string, date string) []string {
	compactDate := strings.ReplaceAll(date, "-", "")
	candidates := []string{
		filepath.Join(logDir, "access.log"),
		filepath.Join(logDir, "access.log-"+compactDate),
		filepath.Join(logDir, "access.log-"+compactDate+".gz"),
		filepath.Join(logDir, "access.log-"+date),
		filepath.Join(logDir, "access.log-"+date+".gz"),
	}
	out := make([]string, 0, len(candidates))
	seen := map[string]struct{}{}
	for _, path := range candidates {
		if _, ok := seen[path]; ok {
			continue
		}
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			out = append(out, path)
			seen[path] = struct{}{}
		}
	}
	sort.Strings(out)
	return out
}

func scanAccessLogFile(path string, analysisRange userTrafficAnalysisRange, out *accessMinuteSummary) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var reader io.Reader = file
	if strings.HasSuffix(path, ".gz") {
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return err
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		logTime, target, ok := parseAccessLogLine(line, analysisRange.StartTime.Location())
		if !ok {
			continue
		}
		if logTime.Before(analysisRange.StartTime) || logTime.After(analysisRange.EndTime) {
			continue
		}
		if !strings.HasSuffix(line, "email: "+analysisRange.UserTag) {
			continue
		}
		minute := logTime.Format("2006-01-02 15:04")
		if out.Targets[minute] == nil {
			out.Targets[minute] = map[string]int{}
		}
		out.Targets[minute][target]++
		out.Matched++
	}
	return scanner.Err()
}

func parseAccessLogLine(line string, location *time.Location) (time.Time, string, bool) {
	if len(line) < len("2006/01/02 15:04:05") {
		return time.Time{}, "", false
	}
	logTime, err := time.ParseInLocation("2006/01/02 15:04:05", line[:19], location)
	if err != nil {
		return time.Time{}, "", false
	}
	fields := strings.Fields(line)
	for _, field := range fields {
		if strings.HasPrefix(field, "tcp:") || strings.HasPrefix(field, "udp:") {
			target := normalizeAccessTarget(extractAccessTargetHost(field))
			if target == "" {
				return time.Time{}, "", false
			}
			return logTime, target, true
		}
	}
	return time.Time{}, "", false
}

func extractAccessTargetHost(field string) string {
	target := strings.TrimPrefix(strings.TrimPrefix(field, "tcp:"), "udp:")
	if target == "" {
		return ""
	}
	if host, _, err := net.SplitHostPort(target); err == nil {
		return strings.Trim(host, "[]")
	}
	if strings.HasPrefix(target, "[") {
		if end := strings.LastIndex(target, "]"); end > 0 {
			return target[1:end]
		}
	}
	if idx := strings.LastIndex(target, ":"); idx > 0 && isDigits(target[idx+1:]) {
		return target[:idx]
	}
	return target
}

func normalizeAccessTarget(target string) string {
	target = strings.TrimSpace(strings.ToLower(target))
	target = strings.Trim(target, "[]")
	target = strings.TrimSuffix(target, ".")
	if target == "" {
		return ""
	}
	if ip := net.ParseIP(target); ip != nil {
		if isInternalAccessTargetIP(ip) {
			return ""
		}
		return target
	}
	if domain, err := publicsuffix.EffectiveTLDPlusOne(target); err == nil {
		return domain
	}
	return fallbackRegistrableDomain(target)
}

func isInternalAccessTargetIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if ip.IsLoopback() || ip.IsUnspecified() || ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsMulticast() {
		return true
	}
	if v4 := ip.To4(); v4 != nil {
		if v4[0] == 0 || v4[0] == 127 || v4[0] >= 224 {
			return true
		}
		if v4[0] == 100 && v4[1] >= 64 && v4[1] <= 127 {
			return true
		}
		return v4[0] == 255 && v4[1] == 255 && v4[2] == 255 && v4[3] == 255
	}
	return false
}

func fallbackRegistrableDomain(target string) string {
	labels := strings.Split(target, ".")
	if len(labels) < 2 {
		return target
	}
	return labels[len(labels)-2] + "." + labels[len(labels)-1]
}

func mergeTrafficAndAccessRows(trafficRows []UserTrafficMinuteRow, targets map[string]map[string]int) []UserTrafficMinuteRow {
	byMinute := make(map[string]*UserTrafficMinuteRow)
	for i := range trafficRows {
		row := trafficRows[i]
		byMinute[row.Minute] = &row
	}
	for minute, targetCounts := range targets {
		row, ok := byMinute[minute]
		if !ok {
			row = &UserTrafficMinuteRow{Minute: minute}
			byMinute[minute] = row
		}
		row.Targets = sortedTargets(targetCounts)
	}
	minutes := make([]string, 0, len(byMinute))
	for minute := range byMinute {
		minutes = append(minutes, minute)
	}
	sort.Strings(minutes)
	out := make([]UserTrafficMinuteRow, 0, len(minutes))
	for _, minute := range minutes {
		row := *byMinute[minute]
		if row.Targets == nil {
			row.Targets = []UserTrafficTarget{}
		}
		out = append(out, row)
	}
	return out
}

func sortedTargets(targetCounts map[string]int) []UserTrafficTarget {
	targets := make([]UserTrafficTarget, 0, len(targetCounts))
	for target, count := range targetCounts {
		targets = append(targets, UserTrafficTarget{Target: target, Count: count})
	}
	sort.Slice(targets, func(i, j int) bool {
		return targets[i].Target < targets[j].Target
	})
	return targets
}
