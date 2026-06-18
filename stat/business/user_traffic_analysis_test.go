package business

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAnalyzeUserTrafficAggregatesTrafficAndTargetsByMinute(t *testing.T) {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("LoadLocation: %v", err)
	}

	store, err := OpenTrafficStore(filepath.Join(t.TempDir(), "stat.db"))
	if err != nil {
		t.Fatalf("OpenTrafficStore: %v", err)
	}
	defer store.Close()

	if err := store.SavePlan(LocalTrafficPlan{Events: []LocalTrafficEvent{
		{Tag: "8", Down: 100, Up: 10, CollectedAt: time.Date(2026, 6, 17, 8, 10, 5, 0, location).Unix()},
		{Tag: "8", Down: 200, Up: 20, CollectedAt: time.Date(2026, 6, 17, 8, 10, 40, 0, location).Unix()},
		{Tag: "8", Down: 50, Up: 5, CollectedAt: time.Date(2026, 6, 17, 8, 11, 0, 0, location).Unix()},
		{Tag: "8", Down: 70, Up: 7, CollectedAt: time.Date(2026, 6, 17, 8, 12, 0, 0, location).Unix()},
		{Tag: "18", Down: 999, Up: 999, CollectedAt: time.Date(2026, 6, 17, 8, 10, 0, 0, location).Unix()},
	}}); err != nil {
		t.Fatalf("SavePlan: %v", err)
	}

	logDir := t.TempDir()
	accessLog := []byte(
		"2026/06/17 08:10:00 1.1.1.1:10001 accepted tcp:example.com:443 email: 8\n" +
			"2026/06/17 08:10:10 1.1.1.1:10002 accepted tcp:rr2---sn-3pm7dne6.googlevideo.com:443 email: 8\n" +
			"2026/06/17 08:10:15 1.1.1.1:10003 accepted tcp:rr3---sn-3pm7dne6.googlevideo.com:443 email: 8\n" +
			"2026/06/17 08:10:20 1.1.1.1:10004 accepted udp:android.clients.google.com:5228 email: 8\n" +
			"2026/06/17 08:11:00 1.1.1.1:10005 accepted udp:mtalk.google.com:5228 email: 8\n" +
			"2026/06/17 08:11:10 1.1.1.1:10006 accepted tcp:8.8.8.8:443 email: 8\n" +
			"2026/06/17 08:10:00 1.1.1.1:10007 accepted tcp:other.example:443 email: 18\n")
	if err := os.WriteFile(filepath.Join(logDir, "access.log"), accessLog, 0644); err != nil {
		t.Fatalf("WriteFile(access.log): %v", err)
	}

	resp, err := AnalyzeUserTraffic(store, logDir, UserTrafficAnalysisRequest{
		UserTag: "8",
		Date:    "20260617",
		Start:   "8:10",
		End:     "8:12",
	})
	if err != nil {
		t.Fatalf("AnalyzeUserTraffic: %v", err)
	}
	if resp.AccessLogMatched != 6 {
		t.Fatalf("AccessLogMatched = %d, want 6", resp.AccessLogMatched)
	}
	if len(resp.Rows) != 3 {
		t.Fatalf("len(Rows) = %d, want 3: %#v", len(resp.Rows), resp.Rows)
	}
	if resp.Rows[0].Minute != "2026-06-17 08:10" {
		t.Fatalf("first minute = %q", resp.Rows[0].Minute)
	}
	if resp.Rows[0].Down != 300 || resp.Rows[0].Up != 30 || resp.Rows[0].Total != 330 || resp.Rows[0].Events != 2 {
		t.Fatalf("first row traffic = %#v", resp.Rows[0])
	}
	if got := targetNames(resp.Rows[0].Targets); len(got) != 3 || got[0] != "example.com" || got[1] != "google.com" || got[2] != "googlevideo.com" {
		t.Fatalf("first row targets = %#v", resp.Rows[0].Targets)
	}
	if got := targetNames(resp.Rows[1].Targets); resp.Rows[1].Minute != "2026-06-17 08:11" || len(got) != 2 || got[0] != "8.8.8.8" || got[1] != "google.com" {
		t.Fatalf("second row = %#v", resp.Rows[1])
	}
	if resp.Rows[2].Targets == nil || len(resp.Rows[2].Targets) != 0 {
		t.Fatalf("traffic-only row targets = %#v, want empty slice", resp.Rows[2].Targets)
	}
}

func targetNames(targets []UserTrafficTarget) []string {
	names := make([]string, 0, len(targets))
	for _, target := range targets {
		names = append(names, target.Target)
	}
	return names
}

func TestNormalizeAccessTargetUsesRegistrableDomain(t *testing.T) {
	cases := map[string]string{
		"www.baidu.com.cn":                  "baidu.com.cn",
		"signaler-pa.clients6.google.com":   "google.com",
		"rr2---sn-3pm7dne6.googlevideo.com": "googlevideo.com",
		"c.pki.goog":                        "pki.goog",
		"8.8.8.8":                           "8.8.8.8",
		"2001:4860:4860::8888":              "2001:4860:4860::8888",
	}
	for target, want := range cases {
		if got := normalizeAccessTarget(target); got != want {
			t.Fatalf("normalizeAccessTarget(%q) = %q, want %q", target, got, want)
		}
	}
}

func TestAnalyzeUserTrafficRejectsRangesLongerThanTwoHours(t *testing.T) {
	_, err := ParseUserTrafficAnalysisRange(UserTrafficAnalysisRequest{
		UserTag: "8",
		Date:    "20260617",
		Start:   "8:10",
		End:     "10:11",
	})
	if err == nil {
		t.Fatal("expected ranges longer than two hours to be rejected")
	}
}
