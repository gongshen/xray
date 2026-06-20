package v2ray_admin

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
	"unicode"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
)

const (
	defaultSiliconFlowBaseURL = "https://api.siliconflow.cn"
	defaultSiliconFlowModel   = "deepseek-ai/DeepSeek-V3.2"
	defaultSiliconFlowTimeout = 90 * time.Second
	minSiliconFlowTimeout     = 30 * time.Second
	defaultSiliconFlowTokens  = 2048
	maxClassifyTargetCount    = 500
)

const trafficTargetClassificationSystemPrompt = `对下面的访问对象（域名或IP地址）进行聚合归类，输出分组结果
规则要求：
1. 按实际用途聚合同类访问对象，每组格式统一为：服务名称（该组内全部IP/域名）：简短用途说明
2. 仅输出本次列表真实存在的分组，无对应样本的类别省略不展示
3. 输出简洁直白，禁止JSON、多余注释、总结性文字

待分析访问对象列表：`

type TrafficTargetClassificationRequest struct {
	Targets []string `json:"targets"`
}

type TrafficTargetClassificationResponse struct {
	Targets       []string `json:"targets"`
	Result        string   `json:"result"`
	CachedTargets int      `json:"cached_targets"`
	LLMTargets    int      `json:"llm_targets"`
}

type trafficTargetClassificationEntry struct {
	Target      string
	ServiceName string
	Purpose     string
}

type siliconFlowChatRequest struct {
	Model       string                   `json:"model"`
	Messages    []siliconFlowChatMessage `json:"messages"`
	Temperature float64                  `json:"temperature"`
	MaxTokens   int                      `json:"max_tokens"`
	Stream      bool                     `json:"stream"`
}

type siliconFlowChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type siliconFlowChatResponse struct {
	Choices []struct {
		Message siliconFlowChatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (serverService *ServerService) ClassifyTrafficTargets(req TrafficTargetClassificationRequest) (*TrafficTargetClassificationResponse, error) {
	targets := normalizeClassificationTargets(req.Targets)
	if len(targets) == 0 {
		return nil, fmt.Errorf("没有可分类的域名/IP")
	}
	if len(targets) > maxClassifyTargetCount {
		return nil, fmt.Errorf("域名/IP数量不能超过%d个", maxClassifyTargetCount)
	}

	cachedEntries, cachedByTarget := loadTrafficTargetClassificationCache(targets)
	missingTargets := make([]string, 0, len(targets)-len(cachedByTarget))
	for _, target := range targets {
		if _, ok := cachedByTarget[target]; !ok {
			missingTargets = append(missingTargets, target)
		}
	}

	newEntries := make([]trafficTargetClassificationEntry, 0, len(missingTargets))
	rawLLMResult := ""
	if len(missingTargets) > 0 {
		result, err := requestSiliconFlowTargetClassification(global.GVA_CONFIG.SiliconFlow, missingTargets)
		if err != nil {
			return nil, err
		}
		rawLLMResult = strings.TrimSpace(result)
		newEntries = parseTrafficTargetClassificationResult(rawLLMResult, missingTargets)
		if len(newEntries) > 0 {
			saveTrafficTargetClassificationCache(newEntries)
		}
	}

	result := renderTrafficTargetClassificationResult(append(cachedEntries, newEntries...), targets)
	if result == "" {
		result = rawLLMResult
	} else if len(newEntries) == 0 && rawLLMResult != "" {
		result = strings.TrimSpace(result + "\n" + rawLLMResult)
	}

	return &TrafficTargetClassificationResponse{
		Targets:       targets,
		Result:        result,
		CachedTargets: len(cachedEntries),
		LLMTargets:    len(missingTargets),
	}, nil
}

func loadTrafficTargetClassificationCache(targets []string) ([]trafficTargetClassificationEntry, map[string]trafficTargetClassificationEntry) {
	byTarget := make(map[string]trafficTargetClassificationEntry, len(targets))
	if global.GVA_DB == nil || len(targets) == 0 {
		return nil, byTarget
	}
	var rows []v2ray.TrafficTargetClassificationCache
	if err := global.GVA_DB.Where("target IN ?", targets).Find(&rows).Error; err != nil {
		logTrafficTargetClassificationCacheWarn("load traffic target classification cache failed", err)
		return nil, byTarget
	}
	entries := make([]trafficTargetClassificationEntry, 0, len(rows))
	for _, row := range rows {
		entry := trafficTargetClassificationEntry{
			Target:      strings.TrimSpace(row.Target),
			ServiceName: strings.TrimSpace(row.ServiceName),
			Purpose:     strings.TrimSpace(row.Purpose),
		}
		if entry.Target == "" || entry.ServiceName == "" || entry.Purpose == "" {
			continue
		}
		byTarget[entry.Target] = entry
		entries = append(entries, entry)
	}
	return entries, byTarget
}

func saveTrafficTargetClassificationCache(entries []trafficTargetClassificationEntry) {
	if global.GVA_DB == nil || len(entries) == 0 {
		return
	}
	now := time.Now().Unix()
	rows := make([]v2ray.TrafficTargetClassificationCache, 0, len(entries))
	seen := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		entry.Target = strings.TrimSpace(entry.Target)
		entry.ServiceName = strings.TrimSpace(entry.ServiceName)
		entry.Purpose = strings.TrimSpace(entry.Purpose)
		if entry.Target == "" || entry.ServiceName == "" || entry.Purpose == "" {
			continue
		}
		if _, ok := seen[entry.Target]; ok {
			continue
		}
		seen[entry.Target] = struct{}{}
		rows = append(rows, v2ray.TrafficTargetClassificationCache{
			Target:      entry.Target,
			ServiceName: entry.ServiceName,
			Purpose:     entry.Purpose,
			CreatedAt:   now,
			UpdatedAt:   now,
		})
	}
	if len(rows) == 0 {
		return
	}
	if err := global.GVA_DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "target"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"service_name",
			"purpose",
			"updated_at",
		}),
	}).Create(&rows).Error; err != nil {
		logTrafficTargetClassificationCacheWarn("save traffic target classification cache failed", err)
	}
}

func logTrafficTargetClassificationCacheWarn(message string, err error) {
	if global.GVA_LOG != nil {
		global.GVA_LOG.Warn(message, zap.Error(err))
	}
}

func requestSiliconFlowTargetClassification(cfg config.SiliconFlow, targets []string) (string, error) {
	apiKey := strings.TrimSpace(cfg.APIKey)
	if apiKey == "" {
		return "", fmt.Errorf("silicon-flow.api-key 未配置")
	}
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		baseURL = defaultSiliconFlowBaseURL
	}
	model := strings.TrimSpace(cfg.Model)
	if model == "" {
		model = defaultSiliconFlowModel
	}
	timeout, err := parseSiliconFlowTimeout(cfg.Timeout)
	if err != nil {
		return "", err
	}

	body, err := json.Marshal(siliconFlowChatRequest{
		Model: model,
		Messages: []siliconFlowChatMessage{
			{Role: "system", Content: trafficTargetClassificationSystemPrompt},
			{Role: "user", Content: strings.Join(targets, "\n")},
		},
		Temperature: 0.1,
		MaxTokens:   defaultSiliconFlowTokens,
		Stream:      false,
	})
	if err != nil {
		return "", err
	}

	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.SetRequestURI(baseURL + "/v1/chat/completions")
	req.SetBody(body)

	client := newSiliconFlowHTTPClient(timeout)
	if err := client.DoTimeout(req, resp, timeout); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "timeout") {
			return "", fmt.Errorf("silicon flow请求超时，已等待%s", timeout)
		}
		return "", err
	}
	if status := resp.StatusCode(); status < fasthttp.StatusOK || status >= fasthttp.StatusMultipleChoices {
		return "", fmt.Errorf("silicon flow returned status %d: %s", status, string(resp.Body()))
	}

	var out siliconFlowChatResponse
	if err := json.Unmarshal(resp.Body(), &out); err != nil {
		return "", err
	}
	if out.Error != nil && strings.TrimSpace(out.Error.Message) != "" {
		return "", fmt.Errorf("%s", out.Error.Message)
	}
	if len(out.Choices) == 0 {
		return "", fmt.Errorf("silicon flow returned empty choices")
	}
	content := strings.TrimSpace(out.Choices[0].Message.Content)
	if content == "" {
		return "", fmt.Errorf("silicon flow returned empty content")
	}
	return content, nil
}

func parseSiliconFlowTimeout(value string) (time.Duration, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultSiliconFlowTimeout, nil
	}
	timeout, err := time.ParseDuration(value)
	if err != nil || timeout <= 0 {
		return 0, fmt.Errorf("silicon-flow.timeout 配置无效")
	}
	if timeout < minSiliconFlowTimeout {
		return minSiliconFlowTimeout, nil
	}
	return timeout, nil
}

func newSiliconFlowHTTPClient(timeout time.Duration) *fasthttp.Client {
	dialTimeout := 10 * time.Second
	if timeout > 0 && timeout < dialTimeout {
		dialTimeout = timeout
	}
	return &fasthttp.Client{
		ReadTimeout:                   timeout,
		WriteTimeout:                  timeout,
		MaxIdleConnDuration:           30 * time.Minute,
		NoDefaultUserAgentHeader:      true,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		Dial: func(addr string) (net.Conn, error) {
			return (&fasthttp.TCPDialer{
				Concurrency:      10,
				DNSCacheDuration: time.Hour * 72,
			}).DialTimeout(addr, dialTimeout)
		},
	}
}

func parseTrafficTargetClassificationResult(result string, expectedTargets []string) []trafficTargetClassificationEntry {
	expected := make(map[string]struct{}, len(expectedTargets))
	for _, target := range expectedTargets {
		expected[target] = struct{}{}
	}
	entries := make([]trafficTargetClassificationEntry, 0, len(expectedTargets))
	seen := make(map[string]struct{}, len(expectedTargets))
	for _, line := range strings.Split(result, "\n") {
		serviceName, lineTargets, purpose, ok := parseTrafficTargetClassificationLine(line)
		if !ok {
			continue
		}
		for _, target := range lineTargets {
			normalized := normalizeClassificationTarget(target)
			if normalized == "" {
				continue
			}
			if _, ok := expected[normalized]; !ok {
				continue
			}
			if _, ok := seen[normalized]; ok {
				continue
			}
			seen[normalized] = struct{}{}
			entries = append(entries, trafficTargetClassificationEntry{
				Target:      normalized,
				ServiceName: serviceName,
				Purpose:     purpose,
			})
		}
	}
	return entries
}

func parseTrafficTargetClassificationLine(line string) (string, []string, string, bool) {
	line = stripClassificationListPrefix(strings.TrimSpace(line))
	if line == "" {
		return "", nil, "", false
	}
	open, close, openSize, closeSize := classificationTargetRange(line)
	if open < 0 || close <= open {
		return "", nil, "", false
	}
	serviceName := strings.TrimSpace(line[:open])
	targets := splitClassificationTargetList(line[open+openSize : close])
	purpose := strings.TrimSpace(line[close+closeSize:])
	purpose = strings.TrimLeftFunc(purpose, func(r rune) bool {
		return r == ':' || r == '：' || r == '-' || unicode.IsSpace(r)
	})
	purpose = strings.TrimSpace(purpose)
	if serviceName == "" || len(targets) == 0 || purpose == "" {
		return "", nil, "", false
	}
	return serviceName, targets, purpose, true
}

func stripClassificationListPrefix(line string) string {
	line = strings.TrimSpace(line)
	line = strings.TrimLeftFunc(line, func(r rune) bool {
		return r == '-' || r == '*' || r == '•' || unicode.IsSpace(r)
	})
	digitEnd := 0
	for digitEnd < len(line) && line[digitEnd] >= '0' && line[digitEnd] <= '9' {
		digitEnd++
	}
	if digitEnd > 0 && digitEnd < len(line) {
		switch line[digitEnd] {
		case '.', ')':
			return strings.TrimSpace(line[digitEnd+1:])
		}
		if strings.HasPrefix(line[digitEnd:], "、") || strings.HasPrefix(line[digitEnd:], "）") {
			return strings.TrimSpace(line[digitEnd+len("、"):])
		}
	}
	return line
}

func classificationTargetRange(line string) (int, int, int, int) {
	if open := strings.Index(line, "（"); open >= 0 {
		if closeRel := strings.Index(line[open+len("（"):], "）"); closeRel >= 0 {
			close := open + len("（") + closeRel
			return open, close, len("（"), len("）")
		}
	}
	if open := strings.Index(line, "("); open >= 0 {
		if closeRel := strings.Index(line[open+1:], ")"); closeRel >= 0 {
			close := open + 1 + closeRel
			return open, close, 1, 1
		}
	}
	return -1, -1, 0, 0
}

func splitClassificationTargetList(value string) []string {
	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == '，' || r == '、' || r == ';' || r == '；' || unicode.IsSpace(r)
	})
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func renderTrafficTargetClassificationResult(entries []trafficTargetClassificationEntry, orderedTargets []string) string {
	if len(entries) == 0 {
		return ""
	}
	byTarget := make(map[string]trafficTargetClassificationEntry, len(entries))
	for _, entry := range entries {
		if entry.Target == "" || entry.ServiceName == "" || entry.Purpose == "" {
			continue
		}
		byTarget[entry.Target] = entry
	}
	type group struct {
		serviceName string
		purpose     string
		targets     []string
	}
	groups := make([]group, 0)
	groupIndex := map[string]int{}
	for _, target := range orderedTargets {
		entry, ok := byTarget[target]
		if !ok {
			continue
		}
		key := entry.ServiceName + "\x00" + entry.Purpose
		idx, ok := groupIndex[key]
		if !ok {
			idx = len(groups)
			groupIndex[key] = idx
			groups = append(groups, group{serviceName: entry.ServiceName, purpose: entry.Purpose})
		}
		groups[idx].targets = append(groups[idx].targets, entry.Target)
	}
	lines := make([]string, 0, len(groups))
	for _, group := range groups {
		if len(group.targets) == 0 {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s（%s）：%s", group.serviceName, strings.Join(group.targets, "、"), group.purpose))
	}
	return strings.Join(lines, "\n")
}

func normalizeClassificationTargets(targets []string) []string {
	seen := make(map[string]struct{}, len(targets))
	out := make([]string, 0, len(targets))
	for _, target := range targets {
		normalized := normalizeClassificationTarget(target)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}
	return out
}

func normalizeClassificationTarget(target string) string {
	target = strings.TrimSpace(strings.ToLower(target))
	target = strings.Trim(target, "[]")
	target = strings.TrimSuffix(target, ".")
	if target == "" || strings.Contains(target, "/") {
		return ""
	}
	if ip := net.ParseIP(target); ip != nil {
		if isInternalClassificationIP(ip) {
			return ""
		}
		return ip.String()
	}
	if strings.Contains(target, ":") {
		return ""
	}
	if !strings.Contains(target, ".") {
		return ""
	}
	return target
}

func isInternalClassificationIP(ip net.IP) bool {
	if ip == nil {
		return false
	}
	if !ip.IsGlobalUnicast() || ip.IsLoopback() || ip.IsUnspecified() || ip.IsPrivate() ||
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
		if v4[0] == 192 && v4[1] == 0 && v4[2] == 2 {
			return true
		}
		if v4[0] == 198 && (v4[1] == 18 || v4[1] == 19) {
			return true
		}
		if v4[0] == 198 && v4[1] == 51 && v4[2] == 100 {
			return true
		}
		if v4[0] == 203 && v4[1] == 0 && v4[2] == 113 {
			return true
		}
		return v4[0] == 255 && v4[1] == 255 && v4[2] == 255 && v4[3] == 255
	}
	v6 := ip.To16()
	if len(v6) == net.IPv6len && v6[0] == 0x20 && v6[1] == 0x01 && v6[2] == 0x0d && v6[3] == 0xb8 {
		return true
	}
	return false
}
