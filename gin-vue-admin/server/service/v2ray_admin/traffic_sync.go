package v2ray_admin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const trafficSyncLimit = 1000

type TrafficSyncEvent struct {
	ID          uint64 `json:"id"`
	Tag         string `json:"tag"`
	Down        uint64 `json:"down"`
	Up          uint64 `json:"up"`
	CollectedAt int64  `json:"collected_at"`
}

type TrafficSyncResponse struct {
	Events []TrafficSyncEvent `json:"events"`
	LastID uint64             `json:"last_id"`
}

type TrafficSyncPlan struct {
	Items       map[string]*v2ray.Stat
	Alerts      []TrafficCollectionAlert
	UsedQuota   uint64
	LastEventID uint64
}

func BuildTrafficSyncPlan(serverIP string, events []TrafficSyncEvent, knownUserTags map[string]struct{}, activeBindingTags map[string]struct{}) TrafficSyncPlan {
	plan := TrafficSyncPlan{
		Items: make(map[string]*v2ray.Stat),
	}
	for _, event := range events {
		if event.ID > plan.LastEventID {
			plan.LastEventID = event.ID
		}
		value := event.Down + event.Up
		if event.Tag == "" || value == 0 {
			continue
		}
		if _, ok := knownUserTags[event.Tag]; !ok {
			plan.Alerts = append(plan.Alerts, TrafficCollectionAlert{
				Reason:   TrafficAnomalyUnknownUser,
				ServerIp: serverIP,
				Tag:      event.Tag,
				Value:    int64(value),
			})
			continue
		}
		if _, ok := activeBindingTags[event.Tag]; !ok {
			plan.Alerts = append(plan.Alerts, TrafficCollectionAlert{
				Reason:   TrafficAnomalyInactiveBinding,
				ServerIp: serverIP,
				Tag:      event.Tag,
				Value:    int64(value),
			})
			continue
		}

		createdAt := createdAtForUnix(event.CollectedAt)
		key := trafficSyncItemKey(event.Tag, createdAt)
		item, ok := plan.Items[key]
		if !ok {
			item = &v2ray.Stat{
				Tag:       event.Tag,
				CreatedAt: createdAt,
				ServerIp:  serverIP,
			}
			plan.Items[key] = item
		}
		item.Down += event.Down
		item.Up += event.Up
		plan.UsedQuota += value
	}
	return plan
}

func trafficSyncItemKey(tag string, createdAt int) string {
	return fmt.Sprintf("%s:%d", tag, createdAt)
}

func createdAtForUnix(ts int64) int {
	if ts <= 0 {
		return createdAtForTime(time.Now())
	}
	return createdAtForTime(time.Unix(ts, 0))
}

func (statService *StatService) SyncServerTrafficWithRetry(srv *v2ray.Server, attempts int) error {
	if attempts <= 0 {
		attempts = 1
	}
	var lastErr error
	for i := 0; i < attempts; i++ {
		if err := statService.SyncServerTraffic(srv); err != nil {
			lastErr = err
			if global.GVA_LOG != nil {
				global.GVA_LOG.Warn("traffic sync attempt failed", zap.Error(err), zap.String("ip", srv.Ip), zap.Int("attempt", i+1))
			}
			if i+1 < attempts {
				time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
			}
			continue
		}
		return nil
	}
	alert := TrafficCollectionAlert{
		Reason:   TrafficAnomalyCollectFail,
		ServerIp: srv.Ip,
		Detail:   lastErr.Error(),
	}
	logTrafficCollectionAlerts([]TrafficCollectionAlert{alert})
	statService.saveTrafficCollectionAlerts([]TrafficCollectionAlert{alert})
	return lastErr
}

func (statService *StatService) SyncServerTraffic(srv *v2ray.Server) error {
	knownUsers, err := statService.knownUserTagSet()
	if err != nil {
		return err
	}
	activeBindings, err := statService.activeBindingTagSet(int(srv.ID))
	if err != nil {
		return err
	}
	afterID, err := statService.trafficSyncCursor(srv.ID)
	if err != nil {
		return err
	}

	for {
		resp, err := statService.fetchTrafficSyncEvents(srv, afterID, trafficSyncLimit)
		if err != nil {
			return err
		}
		plan := BuildTrafficSyncPlan(srv.Ip, resp.Events, knownUsers, activeBindings)
		logTrafficCollectionAlerts(plan.Alerts)
		statService.saveTrafficCollectionAlerts(plan.Alerts)
		if err := statService.saveTrafficSyncPlan(srv.ID, srv.Ip, plan); err != nil {
			return err
		}
		if plan.LastEventID > afterID {
			afterID = plan.LastEventID
		}
		if len(resp.Events) < trafficSyncLimit {
			return nil
		}
	}
}

func (statService *StatService) FlushServerTraffic(srv *v2ray.Server) error {
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI(fmt.Sprintf("http://%s:%d/stat/traffic/collect", srv.Ip, srv.GetStatPort()))
	if err := global.HTTP_CLI.Do(req, resp); err != nil {
		return err
	}
	if status := resp.StatusCode(); status < fasthttp.StatusOK || status >= fasthttp.StatusMultipleChoices {
		return fmt.Errorf("traffic flush returned status %d: %s", status, string(resp.Body()))
	}
	return nil
}

func (statService *StatService) fetchTrafficSyncEvents(srv *v2ray.Server, afterID uint64, limit int) (*TrafficSyncResponse, error) {
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(fmt.Sprintf("http://%s:%d/stat/traffic/sync?after_id=%d&limit=%d", srv.Ip, srv.GetStatPort(), afterID, limit))
	if err := global.HTTP_CLI.Do(req, resp); err != nil {
		return nil, err
	}
	if status := resp.StatusCode(); status < fasthttp.StatusOK || status >= fasthttp.StatusMultipleChoices {
		return nil, fmt.Errorf("traffic sync returned status %d: %s", status, string(resp.Body()))
	}
	if len(resp.Body()) == 0 {
		return nil, fmt.Errorf("traffic sync returned empty body")
	}
	out := new(TrafficSyncResponse)
	if err := json.Unmarshal(resp.Body(), out); err != nil {
		return nil, err
	}
	return out, nil
}

func (statService *StatService) trafficSyncCursor(serverID uint) (uint64, error) {
	cursor := new(v2ray.StatSyncCursor)
	err := global.GVA_DB.Where("server_id = ?", serverID).First(cursor).Error
	if err == nil {
		return cursor.LastEventID, nil
	}
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	return 0, err
}

func (statService *StatService) saveTrafficSyncPlan(serverID uint, serverIP string, plan TrafficSyncPlan) error {
	if len(plan.Items) == 0 && plan.LastEventID == 0 {
		return nil
	}
	now := time.Now().Unix()
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if len(plan.Items) > 0 {
			if err := statsCollector(tx, plan.Items); err != nil {
				return err
			}
		}
		if plan.UsedQuota > 0 {
			if err := tx.Exec("UPDATE v2ray_server set used_quota = used_quota + ? where ID = ?", plan.UsedQuota, serverID).Error; err != nil {
				return err
			}
		}
		if plan.LastEventID > 0 {
			if err := saveTrafficSyncCursor(tx, serverID, serverIP, plan.LastEventID, now); err != nil {
				return err
			}
		}
		return nil
	})
}

func saveTrafficSyncCursor(db *gorm.DB, serverID uint, serverIP string, lastEventID uint64, updatedAt int64) error {
	return db.Exec(
		`INSERT INTO v2ray_stat_sync_cursor (server_id,server_ip,last_event_id,updated_at)
		 VALUES (?,?,?,?)
		 ON DUPLICATE KEY UPDATE server_ip=VALUES(server_ip),last_event_id=VALUES(last_event_id),updated_at=VALUES(updated_at)`,
		serverID, serverIP, lastEventID, updatedAt,
	).Error
}
