package initialize

import "testing"

func TestTrafficCollectCronSpec(t *testing.T) {
	tests := []struct {
		name     string
		interval string
		wantSpec string
		wantOK   bool
	}{
		{name: "empty uses default", interval: "", wantSpec: "@every 1h", wantOK: true},
		{name: "valid duration", interval: "10m", wantSpec: "@every 10m", wantOK: true},
		{name: "trim spaces", interval: " 30m ", wantSpec: "@every 30m", wantOK: true},
		{name: "invalid uses default", interval: "bad", wantSpec: "@every 1h", wantOK: false},
		{name: "zero uses default", interval: "0s", wantSpec: "@every 1h", wantOK: false},
		{name: "negative uses default", interval: "-1m", wantSpec: "@every 1h", wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSpec, gotOK := trafficCollectCronSpec(tt.interval)
			if gotSpec != tt.wantSpec {
				t.Fatalf("trafficCollectCronSpec(%q) spec = %q, want %q", tt.interval, gotSpec, tt.wantSpec)
			}
			if gotOK != tt.wantOK {
				t.Fatalf("trafficCollectCronSpec(%q) ok = %v, want %v", tt.interval, gotOK, tt.wantOK)
			}
		})
	}
}

func TestSysInfoCollectCronSpec(t *testing.T) {
	tests := []struct {
		name     string
		interval string
		wantSpec string
		wantOK   bool
	}{
		{name: "empty uses default", interval: "", wantSpec: "@every 5m", wantOK: true},
		{name: "valid duration", interval: "2m", wantSpec: "@every 2m", wantOK: true},
		{name: "trim spaces", interval: " 10m ", wantSpec: "@every 10m", wantOK: true},
		{name: "invalid uses default", interval: "bad", wantSpec: "@every 5m", wantOK: false},
		{name: "zero uses default", interval: "0s", wantSpec: "@every 5m", wantOK: false},
		{name: "negative uses default", interval: "-1m", wantSpec: "@every 5m", wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSpec, gotOK := sysInfoCollectCronSpec(tt.interval)
			if gotSpec != tt.wantSpec {
				t.Fatalf("sysInfoCollectCronSpec(%q) spec = %q, want %q", tt.interval, gotSpec, tt.wantSpec)
			}
			if gotOK != tt.wantOK {
				t.Fatalf("sysInfoCollectCronSpec(%q) ok = %v, want %v", tt.interval, gotOK, tt.wantOK)
			}
		})
	}
}
