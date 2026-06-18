package v2ray_admin

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/valyala/fasthttp"
)

const (
	defaultSiliconFlowBaseURL = "https://api.siliconflow.cn"
	defaultSiliconFlowModel   = "deepseek-ai/DeepSeek-V3.2"
	defaultSiliconFlowTimeout = 30 * time.Second
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
	Targets []string `json:"targets"`
	Result  string   `json:"result"`
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

	result, err := requestSiliconFlowTargetClassification(global.GVA_CONFIG.SiliconFlow, targets)
	if err != nil {
		return nil, err
	}
	return &TrafficTargetClassificationResponse{
		Targets: targets,
		Result:  result,
	}, nil
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
	timeout := defaultSiliconFlowTimeout
	if strings.TrimSpace(cfg.Timeout) != "" {
		parsed, err := time.ParseDuration(strings.TrimSpace(cfg.Timeout))
		if err != nil || parsed <= 0 {
			return "", fmt.Errorf("silicon-flow.timeout 配置无效")
		}
		timeout = parsed
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

	if err := global.HTTP_CLI.DoTimeout(req, resp, timeout); err != nil {
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
