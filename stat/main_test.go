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
