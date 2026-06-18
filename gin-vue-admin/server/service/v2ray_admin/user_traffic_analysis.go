package v2ray_admin

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	"github.com/valyala/fasthttp"
)

const maxUserTrafficAnalysisRangeMinutes = 120

type UserTrafficAnalysisProxyRequest struct {
	ServerID uint   `json:"server_id"`
	UserTag  string `json:"user_tag"`
	Date     string `json:"date"`
	Start    string `json:"start"`
	End      string `json:"end"`
}

type UserTrafficAnalysisTarget struct {
	Target string `json:"target"`
	Count  int    `json:"count"`
}

type UserTrafficAnalysisMinuteRow struct {
	Minute  string                      `json:"minute"`
	Events  int                         `json:"events"`
	Down    uint64                      `json:"down"`
	Up      uint64                      `json:"up"`
	Total   uint64                      `json:"total"`
	Targets []UserTrafficAnalysisTarget `json:"targets"`
}

type UserTrafficAnalysisResponse struct {
	UserTag          string                         `json:"user_tag"`
	StartTime        string                         `json:"start_time"`
	EndTime          string                         `json:"end_time"`
	AccessLogMatched int                            `json:"access_log_matched"`
	Rows             []UserTrafficAnalysisMinuteRow `json:"rows"`
}

func ValidateUserTrafficAnalysisProxyRequest(req UserTrafficAnalysisProxyRequest) (UserTrafficAnalysisProxyRequest, error) {
	req.UserTag = strings.TrimSpace(req.UserTag)
	req.Date = strings.TrimSpace(req.Date)
	req.Start = strings.TrimSpace(req.Start)
	req.End = strings.TrimSpace(req.End)
	if req.ServerID == 0 {
		return req, fmt.Errorf("server_id is required")
	}
	if req.UserTag == "" {
		return req, fmt.Errorf("user_tag is required")
	}
	if !validAnalysisDate(req.Date) {
		return req, fmt.Errorf("invalid date, expected YYYYMMDD")
	}
	startMinutes, err := analysisClockToMinutes(req.Start)
	if err != nil {
		return req, err
	}
	endMinutes, err := analysisClockToMinutes(req.End)
	if err != nil {
		return req, err
	}
	if endMinutes < startMinutes {
		return req, fmt.Errorf("end time cannot be earlier than start time")
	}
	if endMinutes-startMinutes > maxUserTrafficAnalysisRangeMinutes {
		return req, fmt.Errorf("analysis range cannot exceed 2 hours")
	}
	return req, nil
}

func validAnalysisDate(value string) bool {
	if len(value) == 8 && isDigitString(value) {
		return true
	}
	if len(value) == 10 && value[4] == '-' && value[7] == '-' {
		return isDigitString(value[:4]) && isDigitString(value[5:7]) && isDigitString(value[8:10])
	}
	return false
}

func analysisClockToMinutes(value string) (int, error) {
	parts := strings.Split(value, ":")
	if len(parts) != 2 || !isDigitString(parts[0]) || !isDigitString(parts[1]) || len(parts[1]) != 2 {
		return 0, fmt.Errorf("invalid time, expected H:MM")
	}
	hour, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}
	minute, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, err
	}
	if hour < 0 || hour > 23 || minute < 0 || minute > 59 {
		return 0, fmt.Errorf("invalid time, expected H:MM")
	}
	return hour*60 + minute, nil
}

func isDigitString(value string) bool {
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

func BuildUserTrafficAnalysisStatURL(srv *v2ray.Server, req UserTrafficAnalysisProxyRequest) string {
	values := url.Values{}
	values.Set("user_tag", req.UserTag)
	values.Set("date", req.Date)
	values.Set("start", req.Start)
	values.Set("end", req.End)
	return fmt.Sprintf("http://%s:%d/stat/traffic/user-minute?%s", srv.Ip, srv.GetStatPort(), values.Encode())
}

func (serverService *ServerService) AnalyzeUserTraffic(req UserTrafficAnalysisProxyRequest) (*UserTrafficAnalysisResponse, error) {
	validReq, err := ValidateUserTrafficAnalysisProxyRequest(req)
	if err != nil {
		return nil, err
	}
	srv, err := serverService.GetServer(validReq.ServerID)
	if err != nil {
		return nil, err
	}
	httpReq, httpResp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(httpReq)
	defer fasthttp.ReleaseResponse(httpResp)

	httpReq.SetRequestURI(BuildUserTrafficAnalysisStatURL(&srv, validReq))
	if err := global.HTTP_CLI.Do(httpReq, httpResp); err != nil {
		return nil, err
	}
	if status := httpResp.StatusCode(); status < fasthttp.StatusOK || status >= fasthttp.StatusMultipleChoices {
		return nil, fmt.Errorf("traffic analysis returned status %d: %s", status, string(httpResp.Body()))
	}
	if len(httpResp.Body()) == 0 {
		return nil, fmt.Errorf("traffic analysis returned empty body")
	}
	out := new(UserTrafficAnalysisResponse)
	if err := json.Unmarshal(httpResp.Body(), out); err != nil {
		return nil, err
	}
	return out, nil
}
