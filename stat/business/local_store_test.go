package business

import (
	"path/filepath"
	"testing"
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
