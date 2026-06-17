package business

import (
	"testing"

	statsservice "github.com/xtls/xray-core/app/stats/command"
)

func TestBuildLocalTrafficPlanSeedsMissingSnapshotsWithoutDelta(t *testing.T) {
	stats := []*statsservice.Stat{
		{Name: "user>>>1>>>traffic>>>downlink", Value: 1500},
		{Name: "user>>>1>>>traffic>>>uplink", Value: 300},
	}

	plan := BuildLocalTrafficPlan(stats, nil, 1718424000)

	if len(plan.Events) != 0 {
		t.Fatalf("expected no events for first seen counters, got %d", len(plan.Events))
	}
	if got := plan.Snapshots["user>>>1>>>traffic>>>downlink"]; got != 1500 {
		t.Fatalf("downlink snapshot = %d, want 1500", got)
	}
	if got := plan.Snapshots["user>>>1>>>traffic>>>uplink"]; got != 300 {
		t.Fatalf("uplink snapshot = %d, want 300", got)
	}
}

func TestBuildLocalTrafficPlanComputesDeltaFromSnapshots(t *testing.T) {
	stats := []*statsservice.Stat{
		{Name: "user>>>1>>>traffic>>>downlink", Value: 1800},
		{Name: "user>>>1>>>traffic>>>uplink", Value: 450},
		{Name: "inbound>>>api>>>traffic>>>downlink", Value: 999},
	}
	snapshots := map[string]uint64{
		"user>>>1>>>traffic>>>downlink": 1500,
		"user>>>1>>>traffic>>>uplink":   300,
	}

	plan := BuildLocalTrafficPlan(stats, snapshots, 1718424000)

	if len(plan.Events) != 1 {
		t.Fatalf("expected one event, got %d", len(plan.Events))
	}
	if plan.Events[0].Tag != "1" {
		t.Fatalf("tag = %q, want 1", plan.Events[0].Tag)
	}
	if plan.Events[0].Down != 300 {
		t.Fatalf("down = %d, want 300", plan.Events[0].Down)
	}
	if plan.Events[0].Up != 150 {
		t.Fatalf("up = %d, want 150", plan.Events[0].Up)
	}
	if plan.Events[0].CollectedAt != 1718424000 {
		t.Fatalf("collected_at = %d, want 1718424000", plan.Events[0].CollectedAt)
	}
	if got := plan.Snapshots["user>>>1>>>traffic>>>downlink"]; got != 1800 {
		t.Fatalf("downlink snapshot = %d, want 1800", got)
	}
	if got := plan.Snapshots["user>>>1>>>traffic>>>uplink"]; got != 450 {
		t.Fatalf("uplink snapshot = %d, want 450", got)
	}
}

func TestBuildLocalTrafficPlanTreatsLowerCounterAsXrayReset(t *testing.T) {
	stats := []*statsservice.Stat{
		{Name: "user>>>1>>>traffic>>>downlink", Value: 200},
	}
	snapshots := map[string]uint64{
		"user>>>1>>>traffic>>>downlink": 1800,
	}

	plan := BuildLocalTrafficPlan(stats, snapshots, 1718424000)

	if len(plan.Events) != 1 {
		t.Fatalf("expected one event, got %d", len(plan.Events))
	}
	if plan.Events[0].Down != 200 {
		t.Fatalf("down = %d, want 200", plan.Events[0].Down)
	}
	if got := plan.Snapshots["user>>>1>>>traffic>>>downlink"]; got != 200 {
		t.Fatalf("downlink snapshot = %d, want 200", got)
	}
}
