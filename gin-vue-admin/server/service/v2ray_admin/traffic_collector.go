package v2ray_admin

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	uuid "github.com/satori/go.uuid"
	"github.com/valyala/fasthttp"
	"github.com/xtls/xray-core/app/stats/command"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	TrafficAnomalyInvalidStatName = "invalid_stat_name"
	TrafficAnomalyUnknownUser     = "unknown_user"
	TrafficAnomalyInactiveBinding = "inactive_binding"
	TrafficAnomalyCounterReset    = "counter_reset"
	TrafficAnomalyInvalidBinding  = "invalid_binding"
	TrafficAnomalyCollectFail     = "collect_failed"
	TrafficAnomalyPreCollectFail  = "pre_update_collect_failed"
)

var trafficStatNameRegex = regexp.MustCompile(`^(inbound|outbound|user)>>>([^>]+)>>>traffic>>>(downlink|uplink)$`)

type TrafficCollectionAlert struct {
	Reason   string
	ServerIp string
	Name     string
	Tag      string
	Value    int64
	Detail   string
}

type TrafficCollectionPlan struct {
	Items     map[string]*v2ray.Stat
	Snapshots map[string]uint64
	Alerts    []TrafficCollectionAlert
	UsedQuota uint64
}

func BuildTrafficCollectionPlan(serverIP string, createdAt int, stats []*command.Stat, snapshots map[string]uint64, knownUserTags map[string]struct{}, activeBindingTags map[string]struct{}) TrafficCollectionPlan {
	plan := TrafficCollectionPlan{
		Items:     make(map[string]*v2ray.Stat),
		Snapshots: make(map[string]uint64),
	}
	if snapshots == nil {
		snapshots = map[string]uint64{}
	}

	for _, stat := range stats {
		if stat == nil || stat.Value <= 0 {
			continue
		}

		currentValue := uint64(stat.Value)
		plan.Snapshots[stat.Name] = currentValue

		matches := trafficStatNameRegex.FindStringSubmatch(stat.Name)
		if len(matches) != 4 {
			plan.Alerts = append(plan.Alerts, TrafficCollectionAlert{
				Reason:   TrafficAnomalyInvalidStatName,
				ServerIp: serverIP,
				Name:     stat.Name,
				Value:    stat.Value,
			})
			continue
		}
		if matches[1] != "user" {
			continue
		}

		tag := matches[2]
		if _, ok := knownUserTags[tag]; !ok {
			plan.Alerts = append(plan.Alerts, TrafficCollectionAlert{
				Reason:   TrafficAnomalyUnknownUser,
				ServerIp: serverIP,
				Name:     stat.Name,
				Tag:      tag,
				Value:    stat.Value,
			})
			continue
		}
		if _, ok := activeBindingTags[tag]; !ok {
			plan.Alerts = append(plan.Alerts, TrafficCollectionAlert{
				Reason:   TrafficAnomalyInactiveBinding,
				ServerIp: serverIP,
				Name:     stat.Name,
				Tag:      tag,
				Value:    stat.Value,
			})
		}

		previousValue, seen := snapshots[stat.Name]
		if !seen {
			continue
		}
		delta := currentValue - previousValue
		if currentValue < previousValue {
			delta = currentValue
			plan.Alerts = append(plan.Alerts, TrafficCollectionAlert{
				Reason:   TrafficAnomalyCounterReset,
				ServerIp: serverIP,
				Name:     stat.Name,
				Tag:      tag,
				Value:    stat.Value,
			})
		}
		if delta == 0 {
			continue
		}

		item, ok := plan.Items[tag]
		if !ok {
			item = &v2ray.Stat{
				Tag:       tag,
				CreatedAt: createdAt,
				ServerIp:  serverIP,
			}
			plan.Items[tag] = item
		}
		if matches[3] == "downlink" {
			item.Down += delta
		} else {
			item.Up += delta
		}
		plan.UsedQuota += delta
	}
	return plan
}

func BuildXRayConfigFromBindings(port int64, bindings []*v2ray.Binding) (*v2ray.XrayConfig, []TrafficCollectionAlert) {
	config := v2ray.NewXRayConfig(port)
	settingsCli := make([]*v2ray.XrayConfigSettingsClient, 0, len(bindings))
	alerts := make([]TrafficCollectionAlert, 0)

	for _, binding := range bindings {
		if binding == nil || binding.UserID <= 0 || binding.User.UUID == uuid.Nil {
			alerts = append(alerts, TrafficCollectionAlert{
				Reason: TrafficAnomalyInvalidBinding,
				Detail: "binding has no valid user id or uuid",
			})
			continue
		}
		settingsCli = append(settingsCli, &v2ray.XrayConfigSettingsClient{
			Id:      binding.User.UUID.String(),
			Level:   binding.Level,
			AlterId: binding.AlterID,
			Email:   strconv.Itoa(binding.UserID),
		})
		ensureUserTrafficPolicy(config, binding.Level)
	}
	config.InboundConfigs[0].Settings.Clients = settingsCli
	return config, alerts
}

func ensureUserTrafficPolicy(config *v2ray.XrayConfig, level int64) {
	if config.Policy == nil {
		config.Policy = &v2ray.XrayConfigPolicy{}
	}
	if config.Policy.Levels == nil {
		config.Policy.Levels = make(map[string]*v2ray.XrayConfigPolicyLevel)
	}
	key := strconv.FormatInt(level, 10)
	policyLevel, ok := config.Policy.Levels[key]
	if !ok || policyLevel == nil {
		policyLevel = &v2ray.XrayConfigPolicyLevel{}
		config.Policy.Levels[key] = policyLevel
	}
	policyLevel.StatsUserDownlink = true
	policyLevel.StatsUserUplink = true
}

func (statService *StatService) CollectServerTrafficWithRetry(srv *v2ray.Server, createdAt int, attempts int) error {
	if attempts <= 0 {
		attempts = 1
	}
	var lastErr error
	for i := 0; i < attempts; i++ {
		if err := statService.SyncServerTraffic(srv); err != nil {
			lastErr = err
			if global.GVA_LOG != nil {
				global.GVA_LOG.Warn("traffic sync failed, fallback to legacy collector", zap.Error(err), zap.String("ip", srv.Ip), zap.Int("attempt", i+1))
			}
			if legacyErr := statService.CollectServerTraffic(srv, createdAt); legacyErr != nil {
				lastErr = fmt.Errorf("sync failed: %v; legacy collect failed: %w", err, legacyErr)
				if global.GVA_LOG != nil {
					global.GVA_LOG.Warn("legacy traffic collection attempt failed", zap.Error(legacyErr), zap.String("ip", srv.Ip), zap.Int("attempt", i+1))
				}
				if i+1 < attempts {
					time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
				}
				continue
			}
			return nil
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

func (statService *StatService) PreCollectServerTraffic(srv *v2ray.Server) error {
	if err := statService.FlushServerTraffic(srv); err != nil {
		if global.GVA_LOG != nil {
			global.GVA_LOG.Warn("pre-update traffic flush failed, fallback to legacy collector", zap.Error(err), zap.String("ip", srv.Ip))
		}
		if legacyErr := statService.CollectServerTrafficWithRetry(srv, createdAtForTime(time.Now()), 2); legacyErr != nil {
			alert := TrafficCollectionAlert{
				Reason:   TrafficAnomalyPreCollectFail,
				ServerIp: srv.Ip,
				Detail:   fmt.Sprintf("flush failed: %v; legacy collect failed: %v", err, legacyErr),
			}
			logTrafficCollectionAlerts([]TrafficCollectionAlert{alert})
			statService.saveTrafficCollectionAlerts([]TrafficCollectionAlert{alert})
			return fmt.Errorf("pre-update traffic collection failed: %w", legacyErr)
		}
		return nil
	}
	if err := statService.SyncServerTrafficWithRetry(srv, 2); err != nil {
		if global.GVA_LOG != nil {
			global.GVA_LOG.Warn("pre-update traffic sync failed, fallback to legacy collector", zap.Error(err), zap.String("ip", srv.Ip))
		}
		if legacyErr := statService.CollectServerTrafficWithRetry(srv, createdAtForTime(time.Now()), 2); legacyErr != nil {
			alert := TrafficCollectionAlert{
				Reason:   TrafficAnomalyPreCollectFail,
				ServerIp: srv.Ip,
				Detail:   fmt.Sprintf("sync failed: %v; legacy collect failed: %v", err, legacyErr),
			}
			logTrafficCollectionAlerts([]TrafficCollectionAlert{alert})
			statService.saveTrafficCollectionAlerts([]TrafficCollectionAlert{alert})
			return fmt.Errorf("pre-update traffic collection failed: %w", legacyErr)
		}
	}
	return nil
}

func (statService *StatService) CollectServerTraffic(srv *v2ray.Server, createdAt int) error {
	req, resp := fasthttp.AcquireRequest(), fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(fmt.Sprintf("http://%s:%d/stat/traffic", srv.Ip, srv.GetStatPort()))
	if err := global.HTTP_CLI.Do(req, resp); err != nil {
		return err
	}
	if status := resp.StatusCode(); status < fasthttp.StatusOK || status >= fasthttp.StatusMultipleChoices {
		return fmt.Errorf("traffic collector returned status %d: %s", status, string(resp.Body()))
	}
	if len(resp.Body()) == 0 {
		return fmt.Errorf("traffic collector returned empty body")
	}

	statsResp := new(command.QueryStatsResponse)
	if err := json.Unmarshal(resp.Body(), statsResp); err != nil {
		return err
	}

	knownUsers, err := statService.knownUserTagSet()
	if err != nil {
		return err
	}
	activeBindings, err := statService.activeBindingTagSet(int(srv.ID))
	if err != nil {
		return err
	}
	names := positiveStatNames(statsResp.Stat)
	snapshots, err := statService.trafficSnapshotMap(srv.Ip, names)
	if err != nil {
		return err
	}

	plan := BuildTrafficCollectionPlan(srv.Ip, createdAt, statsResp.Stat, snapshots, knownUsers, activeBindings)
	logTrafficCollectionAlerts(plan.Alerts)
	statService.saveTrafficCollectionAlerts(plan.Alerts)
	return statService.saveTrafficCollectionPlan(srv.ID, srv.Ip, plan)
}

func (statService *StatService) knownUserTagSet() (map[string]struct{}, error) {
	users := make([]system.SysUser, 0)
	if err := global.GVA_DB.Select("id").Find(&users).Error; err != nil {
		return nil, err
	}
	tags := make(map[string]struct{}, len(users))
	for _, user := range users {
		tags[strconv.Itoa(int(user.ID))] = struct{}{}
	}
	return tags, nil
}

func (statService *StatService) activeBindingTagSet(serverID int) (map[string]struct{}, error) {
	bindings := make([]v2ray.Binding, 0)
	if err := global.GVA_DB.Select("user_id").Where("server_id = ? and is_limited = 0", serverID).Find(&bindings).Error; err != nil {
		return nil, err
	}
	tags := make(map[string]struct{}, len(bindings))
	for _, binding := range bindings {
		tags[strconv.Itoa(binding.UserID)] = struct{}{}
	}
	return tags, nil
}

func (statService *StatService) trafficSnapshotMap(serverIP string, names []string) (map[string]uint64, error) {
	out := make(map[string]uint64)
	if len(names) == 0 {
		return out, nil
	}
	snapshots := make([]v2ray.StatSnapshot, 0)
	if err := global.GVA_DB.Where("server_ip = ? and name in ?", serverIP, names).Find(&snapshots).Error; err != nil {
		return nil, err
	}
	for _, snapshot := range snapshots {
		out[snapshot.Name] = snapshot.Value
	}
	return out, nil
}

func (statService *StatService) saveTrafficCollectionPlan(serverID uint, serverIP string, plan TrafficCollectionPlan) error {
	if len(plan.Items) == 0 && len(plan.Snapshots) == 0 {
		return nil
	}
	now := time.Now().Unix()
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if len(plan.Items) > 0 {
			if err := statsCollector(tx, plan.Items); err != nil {
				return err
			}
		}
		if len(plan.Snapshots) > 0 {
			if err := saveTrafficSnapshots(tx, serverIP, plan.Snapshots, now); err != nil {
				return err
			}
		}
		if plan.UsedQuota > 0 {
			if err := tx.Exec("UPDATE v2ray_server set used_quota = used_quota + ? where ID = ?", plan.UsedQuota, serverID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func saveTrafficSnapshots(db *gorm.DB, serverIP string, snapshots map[string]uint64, updatedAt int64) error {
	if len(snapshots) == 0 {
		return nil
	}
	args := make([]interface{}, 0, len(snapshots)*4)
	sql := "INSERT INTO v2ray_stat_snapshot (server_ip,name,value,updated_at) VALUES "
	i := 0
	for name, value := range snapshots {
		if i > 0 {
			sql += ","
		}
		sql += "(?,?,?,?)"
		args = append(args, serverIP, name, value, updatedAt)
		i++
	}
	sql += " ON DUPLICATE KEY UPDATE value=VALUES(value),updated_at=VALUES(updated_at)"
	return db.Exec(sql, args...).Error
}

func (statService *StatService) saveTrafficCollectionAlerts(alerts []TrafficCollectionAlert) {
	if len(alerts) == 0 || global.GVA_DB == nil {
		return
	}
	now := time.Now().Unix()
	rows := make([]v2ray.TrafficAnomaly, 0, len(alerts))
	for _, alert := range alerts {
		rows = append(rows, v2ray.TrafficAnomaly{
			ServerIp:  alert.ServerIp,
			Reason:    alert.Reason,
			Name:      alert.Name,
			Tag:       alert.Tag,
			Value:     alert.Value,
			Detail:    alert.Detail,
			CreatedAt: now,
		})
	}
	if err := global.GVA_DB.Create(&rows).Error; err != nil && global.GVA_LOG != nil {
		global.GVA_LOG.Warn("save traffic anomaly failed", zap.Error(err))
	}
}

func positiveStatNames(stats []*command.Stat) []string {
	names := make([]string, 0, len(stats))
	seen := make(map[string]struct{}, len(stats))
	for _, stat := range stats {
		if stat == nil || stat.Value <= 0 {
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

func logTrafficCollectionAlerts(alerts []TrafficCollectionAlert) {
	if global.GVA_LOG == nil {
		return
	}
	for _, alert := range alerts {
		global.GVA_LOG.Warn("traffic collection anomaly",
			zap.String("reason", alert.Reason),
			zap.String("server_ip", alert.ServerIp),
			zap.String("name", alert.Name),
			zap.String("tag", alert.Tag),
			zap.Int64("value", alert.Value),
			zap.String("detail", alert.Detail),
		)
	}
}

func createdAtForTime(t time.Time) int {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err == nil {
		t = t.In(location)
	}
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}
