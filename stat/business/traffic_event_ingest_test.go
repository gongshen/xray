package business

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/valyala/fasthttp"
)

func TestIngestTrafficEventsPersistsBatch(t *testing.T) {
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

	body, err := json.Marshal(TrafficEventIngestRequest{
		Events: []LocalTrafficEvent{
			{Tag: " 1 ", Down: 100, Up: 20, CollectedAt: 1718424000},
			{Tag: "2", Down: 50, Up: 10},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	var ctx fasthttp.RequestCtx
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBody(body)

	IngestTrafficEvents(&ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusOK {
		t.Fatalf("status = %d, body = %s", ctx.Response.StatusCode(), ctx.Response.Body())
	}

	events, err := store.ListEventsAfter(0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("events length = %d, want 2", len(events))
	}
	if events[0].Tag != "1" || events[0].Down != 100 || events[0].Up != 20 || events[0].CollectedAt != 1718424000 {
		t.Fatalf("unexpected first event: %+v", events[0])
	}
	if events[1].Tag != "2" || events[1].Down != 50 || events[1].Up != 10 || events[1].CollectedAt <= 0 {
		t.Fatalf("unexpected second event: %+v", events[1])
	}
}

func TestIngestTrafficEventsRejectsInvalidEvent(t *testing.T) {
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

	body, err := json.Marshal(TrafficEventIngestRequest{
		Events: []LocalTrafficEvent{
			{Tag: " ", Down: 100},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	var ctx fasthttp.RequestCtx
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBody(body)

	IngestTrafficEvents(&ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusBadRequest {
		t.Fatalf("status = %d, want 400", ctx.Response.StatusCode())
	}

	events, err := store.ListEventsAfter(0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 0 {
		t.Fatalf("events length = %d, want 0", len(events))
	}
}

func TestIngestTrafficEventsRejectsTooLargeBatch(t *testing.T) {
	previousStore := LocalStore
	LocalStore = &TrafficStore{}
	defer func() {
		LocalStore = previousStore
	}()

	req := TrafficEventIngestRequest{
		Events: make([]LocalTrafficEvent, maxTrafficEventIngestBatchSize+1),
	}
	for i := range req.Events {
		req.Events[i] = LocalTrafficEvent{Tag: "1", Down: 1}
	}
	body, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}

	var ctx fasthttp.RequestCtx
	ctx.Request.Header.SetMethod(fasthttp.MethodPost)
	ctx.Request.SetBody(body)

	IngestTrafficEvents(&ctx)

	if ctx.Response.StatusCode() != fasthttp.StatusBadRequest {
		t.Fatalf("status = %d, want 400", ctx.Response.StatusCode())
	}
}
