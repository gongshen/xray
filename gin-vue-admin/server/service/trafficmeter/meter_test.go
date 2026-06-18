package trafficmeter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

func TestMeterFlushReportsAggregatedEvent(t *testing.T) {
	var got trafficEventBatch
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stat/traffic/event" {
			t.Fatalf("path = %s, want /stat/traffic/event", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	meter := NewFromConfig(config.TrafficMeter{
		Enable:        true,
		StatURL:       server.URL,
		Tag:           "7",
		FlushInterval: "1h",
	}, &fasthttp.Client{}, nil)
	meter.Add(100, 200)
	meter.Add(3, 4)

	if err := meter.Flush(); err != nil {
		t.Fatal(err)
	}

	if len(got.Events) != 1 {
		t.Fatalf("events length = %d, want 1", len(got.Events))
	}
	event := got.Events[0]
	if event.Tag != "7" || event.Up != 103 || event.Down != 204 || event.CollectedAt <= 0 {
		t.Fatalf("unexpected event: %+v", event)
	}
}

func TestMeterFlushRestoresTrafficAfterFailure(t *testing.T) {
	var calls int32
	var got trafficEventBatch
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&calls, 1) == 1 {
			http.Error(w, "temporary failure", http.StatusInternalServerError)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	meter := NewFromConfig(config.TrafficMeter{
		Enable:        true,
		StatURL:       server.URL,
		Tag:           "1",
		FlushInterval: "1h",
	}, &fasthttp.Client{}, nil)
	meter.Add(10, 20)

	if err := meter.Flush(); err == nil {
		t.Fatal("first flush error = nil, want non-nil")
	}
	if err := meter.Flush(); err != nil {
		t.Fatal(err)
	}

	if len(got.Events) != 1 {
		t.Fatalf("events length = %d, want 1", len(got.Events))
	}
	if got.Events[0].Up != 10 || got.Events[0].Down != 20 {
		t.Fatalf("unexpected event after retry: %+v", got.Events[0])
	}
}

func TestNewFromConfigDefaultsTagAndInterval(t *testing.T) {
	meter := NewFromConfig(config.TrafficMeter{
		Enable:  true,
		StatURL: "http://127.0.0.1:56611/",
	}, &fasthttp.Client{}, nil)

	if meter == nil {
		t.Fatal("meter is nil")
	}
	if meter.tag != DefaultTag {
		t.Fatalf("tag = %q, want %q", meter.tag, DefaultTag)
	}
	if meter.flushInterval != DefaultFlushInterval {
		t.Fatalf("flushInterval = %s, want %s", meter.flushInterval, DefaultFlushInterval)
	}
	if meter.statURL != "http://127.0.0.1:56611" {
		t.Fatalf("statURL = %q, want trimmed URL", meter.statURL)
	}
}

func TestMiddlewareRecordsGinRequestAndResponseBytes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	meter := NewFromConfig(config.TrafficMeter{
		Enable:        true,
		StatURL:       "http://127.0.0.1:56611",
		Tag:           "1",
		FlushInterval: "1h",
	}, &fasthttp.Client{}, nil)
	router := gin.New()
	router.Use(Middleware(meter))
	router.POST("/hello", func(c *gin.Context) {
		c.String(http.StatusCreated, "response-body")
	})

	req := httptest.NewRequest(http.MethodPost, "/hello?x=1", http.NoBody)
	req.ContentLength = 0
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	up, down := meter.take()
	if up == 0 || down == 0 {
		t.Fatalf("meter bytes up=%d down=%d, want positive values", up, down)
	}
}
