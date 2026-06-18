package v2ray_admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
)

func TestNormalizeClassificationTargetsDedupesAndKeepsPublicIPs(t *testing.T) {
	got := normalizeClassificationTargets([]string{
		"Google.com",
		"google.com.",
		"8.8.8.8",
		"149.154.167.99",
		"127.0.0.1",
		"10.0.0.1",
		"100.64.0.1",
		"192.0.2.1",
		"198.51.100.1",
		"203.0.113.1",
		"2001:4860:4860::8888",
		"2001:db8::1",
		"::1",
		"fc00::1",
		"pki.goog",
		"not-domain",
		"https://example.com",
	})
	want := []string{"google.com", "8.8.8.8", "149.154.167.99", "2001:4860:4860::8888", "pki.goog"}
	if len(got) != len(want) {
		t.Fatalf("targets length = %d, want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("targets[%d] = %q, want %q; all = %#v", i, got[i], want[i], got)
		}
	}
}

func TestRequestSiliconFlowTargetClassification(t *testing.T) {
	var gotReq siliconFlowChatRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Fatalf("path = %s, want /v1/chat/completions", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-key" {
			t.Fatalf("Authorization = %q, want Bearer test-key", auth)
		}
		if err := json.NewDecoder(r.Body).Decode(&gotReq); err != nil {
			t.Fatal(err)
		}
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"Google（google.com）：搜索和技术服务\nTelegram（149.154.167.99）：即时通讯服务"}}]}`))
	}))
	defer server.Close()

	result, err := requestSiliconFlowTargetClassification(config.SiliconFlow{
		APIKey:  "test-key",
		BaseURL: server.URL,
		Model:   "deepseek-ai/DeepSeek-V3.2",
		Timeout: "5s",
	}, []string{"google.com", "149.154.167.99"})
	if err != nil {
		t.Fatal(err)
	}
	if result != "Google（google.com）：搜索和技术服务\nTelegram（149.154.167.99）：即时通讯服务" {
		t.Fatalf("result = %q", result)
	}
	if gotReq.Model != "deepseek-ai/DeepSeek-V3.2" {
		t.Fatalf("model = %q", gotReq.Model)
	}
	if gotReq.MaxTokens != defaultSiliconFlowTokens {
		t.Fatalf("max_tokens = %d, want %d", gotReq.MaxTokens, defaultSiliconFlowTokens)
	}
	if len(gotReq.Messages) != 2 {
		t.Fatalf("messages length = %d, want 2", len(gotReq.Messages))
	}
	if !strings.Contains(gotReq.Messages[0].Content, "按实际用途聚合同类访问对象") {
		t.Fatalf("system prompt = %q", gotReq.Messages[0].Content)
	}
	if gotReq.Messages[1].Content != "google.com\n149.154.167.99" {
		t.Fatalf("user content = %q", gotReq.Messages[1].Content)
	}
}

func TestRequestSiliconFlowTargetClassificationRequiresAPIKey(t *testing.T) {
	_, err := requestSiliconFlowTargetClassification(config.SiliconFlow{}, []string{"google.com"})
	if err == nil || !strings.Contains(err.Error(), "api-key") {
		t.Fatalf("error = %v, want api-key error", err)
	}
}
