package main

import (
	"flag"
	"testing"
)

func TestDefaultCollectIntervalIsTenSeconds(t *testing.T) {
	collectIntervalFlag := flag.Lookup("collect-interval")
	if collectIntervalFlag == nil {
		t.Fatal("collect-interval flag is not registered")
	}
	if collectIntervalFlag.DefValue != "10s" {
		t.Fatalf("collect-interval default = %q, want %q", collectIntervalFlag.DefValue, "10s")
	}
}

func TestDefaultTrafficRetentionMonthsIsOneYear(t *testing.T) {
	trafficRetentionFlag := flag.Lookup("traffic-retention-months")
	if trafficRetentionFlag == nil {
		t.Fatal("traffic-retention-months flag is not registered")
	}
	if trafficRetentionFlag.DefValue != "12" {
		t.Fatalf("traffic-retention-months default = %q, want %q", trafficRetentionFlag.DefValue, "12")
	}
}
