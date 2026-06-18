package business

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gongshen/xray/stat/conn"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	statsservice "github.com/xtls/xray-core/app/stats/command"
)

const DefaultStatAPITrafficTag = "1"
const maxTrafficEventIngestBatchSize = 1000

var (
	LocalStore        *TrafficStore
	StatAPITrafficTag = DefaultStatAPITrafficTag
	collectMu         sync.Mutex
)

func CollectTraffic(reqCtx *fasthttp.RequestCtx) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := conn.StatServiceClient.QueryStats(ctx, NewTrafficQueryStatsRequest())
	if err != nil {
		reqCtx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(resp)
	if err != nil {
		reqCtx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	reqCtx.Success("application/json", data)
}

func NewTrafficQueryStatsRequest() *statsservice.QueryStatsRequest {
	return &statsservice.QueryStatsRequest{
		Reset_: false,
	}
}

func CollectTrafficToLocalStore(store *TrafficStore) error {
	collectMu.Lock()
	defer collectMu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := conn.StatServiceClient.QueryStats(ctx, NewTrafficQueryStatsRequest())
	if err != nil {
		return err
	}

	names := positiveLocalStatNames(resp.Stat)
	snapshots, err := store.LoadSnapshots(names)
	if err != nil {
		return err
	}
	plan := BuildLocalTrafficPlan(resp.Stat, snapshots, time.Now().Unix())
	return store.SavePlan(plan)
}

func CollectTrafficToLocalStoreHandler(ctx *fasthttp.RequestCtx) {
	if LocalStore == nil {
		ctx.Error("traffic store is not initialized", http.StatusServiceUnavailable)
		return
	}
	if err := CollectTrafficToLocalStore(LocalStore); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	ctx.SuccessString("application/json", `{"ok":true}`)
}

func SetStatAPITrafficTag(tag string) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		tag = DefaultStatAPITrafficTag
	}
	StatAPITrafficTag = tag
}

func RecordStatAPIUsage(up uint64, down uint64) {
	if LocalStore == nil || (up == 0 && down == 0) {
		return
	}

	collectMu.Lock()
	defer collectMu.Unlock()

	err := LocalStore.SavePlan(LocalTrafficPlan{
		Events: []LocalTrafficEvent{
			{
				Tag:         StatAPITrafficTag,
				Down:        down,
				Up:          up,
				CollectedAt: time.Now().Unix(),
			},
		},
	})
	if err != nil {
		logrus.WithError(err).Warn("record stat api traffic failed")
	}
}

type TrafficSyncResponse struct {
	Events []LocalTrafficEvent `json:"events"`
	LastID uint64              `json:"last_id"`
}

type TrafficEventIngestRequest struct {
	Events []LocalTrafficEvent `json:"events"`
}

func IngestTrafficEvents(ctx *fasthttp.RequestCtx) {
	if LocalStore == nil {
		ctx.Error("traffic store is not initialized", http.StatusServiceUnavailable)
		return
	}
	if !ctx.IsPost() {
		ctx.Error("method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TrafficEventIngestRequest
	if err := json.Unmarshal(ctx.Request.Body(), &req); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	if len(req.Events) == 0 {
		ctx.Error("events is empty", http.StatusBadRequest)
		return
	}
	if len(req.Events) > maxTrafficEventIngestBatchSize {
		ctx.Error("events is too large", http.StatusBadRequest)
		return
	}

	now := time.Now().Unix()
	plan := LocalTrafficPlan{
		Events: make([]LocalTrafficEvent, 0, len(req.Events)),
	}
	for i, event := range req.Events {
		event.Tag = strings.TrimSpace(event.Tag)
		if event.Tag == "" {
			ctx.Error(fmt.Sprintf("events[%d].tag is empty", i), http.StatusBadRequest)
			return
		}
		if event.Down == 0 && event.Up == 0 {
			ctx.Error(fmt.Sprintf("events[%d] traffic is empty", i), http.StatusBadRequest)
			return
		}
		if event.CollectedAt <= 0 {
			event.CollectedAt = now
		}
		plan.Events = append(plan.Events, LocalTrafficEvent{
			Tag:         event.Tag,
			Down:        event.Down,
			Up:          event.Up,
			CollectedAt: event.CollectedAt,
		})
	}

	collectMu.Lock()
	defer collectMu.Unlock()

	if err := LocalStore.SavePlan(plan); err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	ctx.SuccessString("application/json", `{"ok":true}`)
}

func SyncLocalTraffic(ctx *fasthttp.RequestCtx) {
	if LocalStore == nil {
		ctx.Error("traffic store is not initialized", http.StatusServiceUnavailable)
		return
	}
	afterID, err := parseUintQuery(ctx, "after_id", 0)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	limitValue, err := parseUintQuery(ctx, "limit", 1000)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	events, err := LocalStore.ListEventsAfter(afterID, int(limitValue))
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	resp := TrafficSyncResponse{Events: events, LastID: afterID}
	if len(events) > 0 {
		resp.LastID = events[len(events)-1].ID
	}
	data, err := json.Marshal(resp)
	if err != nil {
		ctx.Error(err.Error(), http.StatusBadRequest)
		return
	}
	ctx.Success("application/json", data)
}

func parseUintQuery(ctx *fasthttp.RequestCtx, key string, defaultValue uint64) (uint64, error) {
	value := string(ctx.QueryArgs().Peek(key))
	if value == "" {
		return defaultValue, nil
	}
	return strconv.ParseUint(value, 10, 64)
}

func positiveLocalStatNames(stats []*statsservice.Stat) []string {
	names := make([]string, 0, len(stats))
	seen := make(map[string]struct{}, len(stats))
	for _, stat := range stats {
		if stat == nil || stat.Value <= 0 {
			continue
		}
		if !localTrafficStatNameRegex.MatchString(stat.Name) {
			continue
		}
		if _, ok := seen[stat.Name]; ok {
			continue
		}
		seen[stat.Name] = struct{}{}
		names = append(names, stat.Name)
	}
	return names
}

func StartLocalTrafficCollector(store *TrafficStore, interval time.Duration, stop <-chan struct{}) {
	if interval <= 0 {
		interval = 10 * time.Second
	}
	_ = CollectTrafficToLocalStore(store)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			_ = CollectTrafficToLocalStore(store)
		}
	}
}
