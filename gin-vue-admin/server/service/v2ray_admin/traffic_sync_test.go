package v2ray_admin

import (
	"testing"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	"github.com/stretchr/testify/require"
)

func TestBuildTrafficSyncPlanAggregatesKnownEvents(t *testing.T) {
	events := []TrafficSyncEvent{
		{ID: 11, Tag: "1", Down: 100, Up: 20, CollectedAt: 1718424000},
		{ID: 12, Tag: "1", Down: 30, Up: 5, CollectedAt: 1718424010},
	}

	plan := BuildTrafficSyncPlan("10.0.0.1", events, tagSet("1"), tagSet("1"))

	key := trafficSyncItemKey("1", createdAtForUnix(1718424000))
	require.Empty(t, plan.Alerts)
	require.Equal(t, uint64(130), plan.Items[key].Down)
	require.Equal(t, uint64(25), plan.Items[key].Up)
	require.Equal(t, uint64(155), plan.UsedQuota)
	require.Equal(t, uint64(12), plan.LastEventID)
	require.Equal(t, createdAtForUnix(1718424000), plan.Items[key].CreatedAt)
	require.Equal(t, "10.0.0.1", plan.Items[key].ServerIp)
}

func TestBuildTrafficSyncPlanAlertsUnknownAndInactiveUsers(t *testing.T) {
	events := []TrafficSyncEvent{
		{ID: 11, Tag: "ghost", Down: 100, CollectedAt: 1718424000},
		{ID: 12, Tag: "2", Up: 50, CollectedAt: 1718424010},
	}

	plan := BuildTrafficSyncPlan("10.0.0.1", events, tagSet("2"), tagSet("1"))

	require.Empty(t, plan.Items)
	require.Equal(t, uint64(12), plan.LastEventID)
	require.Len(t, plan.Alerts, 2)
	require.Equal(t, TrafficAnomalyUnknownUser, plan.Alerts[0].Reason)
	require.Equal(t, TrafficAnomalyInactiveBinding, plan.Alerts[1].Reason)
	require.Equal(t, int64(50), plan.Alerts[1].Value)
}

func TestBuildTrafficSyncPlanSplitsTrafficByDay(t *testing.T) {
	location, err := time.LoadLocation("Asia/Shanghai")
	require.NoError(t, err)
	firstDay := time.Date(2024, 6, 15, 23, 59, 59, 0, location).Unix()
	secondDay := time.Date(2024, 6, 16, 0, 0, 0, 0, location).Unix()
	events := []TrafficSyncEvent{
		{ID: 11, Tag: "1", Down: 100, CollectedAt: firstDay},
		{ID: 12, Tag: "1", Down: 50, CollectedAt: secondDay},
	}

	plan := BuildTrafficSyncPlan("10.0.0.1", events, tagSet("1"), tagSet("1"))

	require.Len(t, plan.Items, 2)
	require.Equal(t, uint64(100), plan.Items[trafficSyncItemKey("1", createdAtForUnix(firstDay))].Down)
	require.Equal(t, uint64(50), plan.Items[trafficSyncItemKey("1", createdAtForUnix(secondDay))].Down)
}

func TestBuildTrafficSyncPlanOutputMatchesStatsCollectorInput(t *testing.T) {
	events := []TrafficSyncEvent{
		{ID: 11, Tag: "1", Down: 100, CollectedAt: 1718424000},
	}

	plan := BuildTrafficSyncPlan("10.0.0.1", events, tagSet("1"), tagSet("1"))

	require.IsType(t, map[string]*v2ray.Stat{}, plan.Items)
}
