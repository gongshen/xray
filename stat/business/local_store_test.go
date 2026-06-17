package business

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestTrafficStorePersistsEventsAndListsAfterID(t *testing.T) {
	store, err := OpenTrafficStore(filepath.Join(t.TempDir(), "stat.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	plan := LocalTrafficPlan{
		Events: []LocalTrafficEvent{
			{Tag: "1", Down: 100, Up: 20, CollectedAt: 1718424000},
			{Tag: "2", Down: 50, Up: 10, CollectedAt: 1718424001},
		},
		Snapshots: map[string]uint64{
			"user>>>1>>>traffic>>>downlink": 1500,
			"user>>>1>>>traffic>>>uplink":   300,
		},
	}
	if err := store.SavePlan(plan); err != nil {
		t.Fatal(err)
	}

	events, err := store.ListEventsAfter(1, 10)
	if err != nil {
		t.Fatal(err)
	}

	if len(events) != 1 {
		t.Fatalf("expected one event after id 1, got %d", len(events))
	}
	if events[0].ID != 2 || events[0].Tag != "2" || events[0].Down != 50 || events[0].Up != 10 {
		t.Fatalf("unexpected event: %+v", events[0])
	}
}

func TestTrafficStoreLoadsSnapshots(t *testing.T) {
	store, err := OpenTrafficStore(filepath.Join(t.TempDir(), "stat.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	err = store.SavePlan(LocalTrafficPlan{
		Snapshots: map[string]uint64{
			"user>>>1>>>traffic>>>downlink": 1500,
			"user>>>1>>>traffic>>>uplink":   300,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	snapshots, err := store.LoadSnapshots([]string{
		"user>>>1>>>traffic>>>downlink",
		"user>>>1>>>traffic>>>uplink",
	})
	if err != nil {
		t.Fatal(err)
	}

	if snapshots["user>>>1>>>traffic>>>downlink"] != 1500 {
		t.Fatalf("downlink snapshot = %d, want 1500", snapshots["user>>>1>>>traffic>>>downlink"])
	}
	if snapshots["user>>>1>>>traffic>>>uplink"] != 300 {
		t.Fatalf("uplink snapshot = %d, want 300", snapshots["user>>>1>>>traffic>>>uplink"])
	}
}

func TestTrafficStoreCleanupOldTrafficEventsRemovesOnlyEventsOlderThanRetention(t *testing.T) {
	store, err := OpenTrafficStore(filepath.Join(t.TempDir(), "stat.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	now := time.Date(2026, 6, 17, 12, 0, 0, 0, time.UTC)
	cutoff := now.AddDate(0, -12, 0).Unix()
	err = store.SavePlan(LocalTrafficPlan{
		Events: []LocalTrafficEvent{
			{Tag: "old", Down: 1, CollectedAt: cutoff - 1},
			{Tag: "boundary", Down: 1, CollectedAt: cutoff},
			{Tag: "fresh", Down: 1, CollectedAt: cutoff + 1},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	removed, err := store.CleanupOldTrafficEvents(now, 12)
	if err != nil {
		t.Fatalf("CleanupOldTrafficEvents returned error: %v", err)
	}
	if removed != 1 {
		t.Fatalf("removed = %d, want 1", removed)
	}

	events, err := store.ListEventsAfter(0, 10)
	if err != nil {
		t.Fatal(err)
	}
	gotTags := make([]string, 0, len(events))
	for _, event := range events {
		gotTags = append(gotTags, event.Tag)
	}
	wantTags := []string{"boundary", "fresh"}
	if !reflect.DeepEqual(gotTags, wantTags) {
		t.Fatalf("remaining event tags = %#v, want %#v", gotTags, wantTags)
	}
}

func TestStartTrafficEventCleanerRunsImmediately(t *testing.T) {
	store, err := OpenTrafficStore(filepath.Join(t.TempDir(), "stat.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	err = store.SavePlan(LocalTrafficPlan{
		Events: []LocalTrafficEvent{
			{Tag: "old", Down: 1, CollectedAt: time.Now().AddDate(0, -13, 0).Unix()},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	stop := make(chan struct{})
	close(stop)
	StartTrafficEventCleaner(store, 12, time.Hour, stop)

	events, err := store.ListEventsAfter(0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 0 {
		t.Fatalf("events after immediate cleanup = %#v, want empty", events)
	}
}
