package business

import (
	"path/filepath"
	"testing"
)

func TestRecordStatAPIUsagePersistsTagOneEvent(t *testing.T) {
	store, err := OpenTrafficStore(filepath.Join(t.TempDir(), "stat.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	previousStore := LocalStore
	LocalStore = store
	defer func() {
		LocalStore = previousStore
	}()

	RecordStatAPIUsage(123, 456)

	events, err := store.ListEventsAfter(0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("events length = %d, want 1", len(events))
	}
	if events[0].Tag != "1" {
		t.Fatalf("tag = %q, want 1", events[0].Tag)
	}
	if events[0].Up != 123 {
		t.Fatalf("up = %d, want 123", events[0].Up)
	}
	if events[0].Down != 456 {
		t.Fatalf("down = %d, want 456", events[0].Down)
	}
	if events[0].CollectedAt <= 0 {
		t.Fatalf("collected_at = %d, want positive unix timestamp", events[0].CollectedAt)
	}
}

func TestRecordStatAPIUsageSkipsEmptyTraffic(t *testing.T) {
	store, err := OpenTrafficStore(filepath.Join(t.TempDir(), "stat.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	previousStore := LocalStore
	LocalStore = store
	defer func() {
		LocalStore = previousStore
	}()

	RecordStatAPIUsage(0, 0)

	events, err := store.ListEventsAfter(0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 0 {
		t.Fatalf("events length = %d, want 0", len(events))
	}
}

func TestRecordStatAPIUsageUsesConfiguredTag(t *testing.T) {
	store, err := OpenTrafficStore(filepath.Join(t.TempDir(), "stat.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	previousStore := LocalStore
	LocalStore = store
	defer func() {
		LocalStore = previousStore
	}()

	previousTag := StatAPITrafficTag
	SetStatAPITrafficTag("88")
	defer func() {
		SetStatAPITrafficTag(previousTag)
	}()

	RecordStatAPIUsage(1, 2)

	events, err := store.ListEventsAfter(0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("events length = %d, want 1", len(events))
	}
	if events[0].Tag != "88" {
		t.Fatalf("tag = %q, want 88", events[0].Tag)
	}
}

func TestSetStatAPITrafficTagFallsBackToDefaultForEmptyValue(t *testing.T) {
	previousTag := StatAPITrafficTag
	defer func() {
		SetStatAPITrafficTag(previousTag)
	}()

	SetStatAPITrafficTag("88")
	SetStatAPITrafficTag(" ")

	if StatAPITrafficTag != "1" {
		t.Fatalf("StatAPITrafficTag = %q, want 1", StatAPITrafficTag)
	}
}
