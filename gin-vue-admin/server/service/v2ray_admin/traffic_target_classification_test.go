package v2ray_admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
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

func TestParseSiliconFlowTimeoutUsesMinimum(t *testing.T) {
	got, err := parseSiliconFlowTimeout("3s")
	if err != nil {
		t.Fatal(err)
	}
	if got != minSiliconFlowTimeout {
		t.Fatalf("timeout = %s, want %s", got, minSiliconFlowTimeout)
	}
}

func TestParseSiliconFlowTimeoutUsesDefault(t *testing.T) {
	got, err := parseSiliconFlowTimeout("")
	if err != nil {
		t.Fatal(err)
	}
	if got != defaultSiliconFlowTimeout {
		t.Fatalf("timeout = %s, want %s", got, defaultSiliconFlowTimeout)
	}
}

func TestParseTrafficTargetClassificationResult(t *testing.T) {
	got := parseTrafficTargetClassificationResult(
		"1. Google（google.com、pki.goog）：搜索和证书服务\n- Telegram (149.154.167.99, 2001:4860:4860::8888): 即时通讯服务\nOther（unknown.com）：未知服务",
		[]string{"google.com", "pki.goog", "149.154.167.99", "2001:4860:4860::8888"},
	)
	if len(got) != 4 {
		t.Fatalf("entries length = %d, want 4: %#v", len(got), got)
	}
	want := map[string]trafficTargetClassificationEntry{
		"google.com":           {Target: "google.com", ServiceName: "Google", Purpose: "搜索和证书服务"},
		"pki.goog":             {Target: "pki.goog", ServiceName: "Google", Purpose: "搜索和证书服务"},
		"149.154.167.99":       {Target: "149.154.167.99", ServiceName: "Telegram", Purpose: "即时通讯服务"},
		"2001:4860:4860::8888": {Target: "2001:4860:4860::8888", ServiceName: "Telegram", Purpose: "即时通讯服务"},
	}
	for _, entry := range got {
		if entry != want[entry.Target] {
			t.Fatalf("entry = %#v, want %#v", entry, want[entry.Target])
		}
	}
}

func TestRenderTrafficTargetClassificationResultKeepsTargetOrder(t *testing.T) {
	got := renderTrafficTargetClassificationResult([]trafficTargetClassificationEntry{
		{Target: "pki.goog", ServiceName: "Google", Purpose: "搜索和证书服务"},
		{Target: "google.com", ServiceName: "Google", Purpose: "搜索和证书服务"},
		{Target: "149.154.167.99", ServiceName: "Telegram", Purpose: "即时通讯服务"},
	}, []string{"google.com", "149.154.167.99", "pki.goog"})
	want := "Google（google.com、pki.goog）：搜索和证书服务\nTelegram（149.154.167.99）：即时通讯服务"
	if got != want {
		t.Fatalf("rendered = %q, want %q", got, want)
	}
}

func TestClassifyTrafficTargetsUsesCacheWithoutLLM(t *testing.T) {
	oldDB := global.GVA_DB
	oldConfig := global.GVA_CONFIG
	defer func() {
		global.GVA_DB = oldDB
		global.GVA_CONFIG = oldConfig
	}()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&v2ray.TrafficTargetClassificationCache{}); err != nil {
		t.Fatal(err)
	}
	now := time.Now().Unix()
	if err := db.Create([]v2ray.TrafficTargetClassificationCache{
		{Target: "google.com", ServiceName: "Google", Purpose: "搜索服务", CreatedAt: now, UpdatedAt: now},
		{Target: "149.154.167.99", ServiceName: "Telegram", Purpose: "即时通讯服务", CreatedAt: now, UpdatedAt: now},
	}).Error; err != nil {
		t.Fatal(err)
	}
	global.GVA_DB = db
	global.GVA_CONFIG.SiliconFlow = config.SiliconFlow{}

	resp, err := (&ServerService{}).ClassifyTrafficTargets(TrafficTargetClassificationRequest{
		Targets: []string{"google.com", "149.154.167.99"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.CachedTargets != 2 || resp.LLMTargets != 0 {
		t.Fatalf("cache stats = cached %d llm %d", resp.CachedTargets, resp.LLMTargets)
	}
	want := "Google（google.com）：搜索服务\nTelegram（149.154.167.99）：即时通讯服务"
	if resp.Result != want {
		t.Fatalf("result = %q, want %q", resp.Result, want)
	}
}
