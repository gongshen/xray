package v2ray_admin

import (
	"testing"

	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/xtls/xray-core/app/stats/command"
)

func TestBuildTrafficCollectionPlanComputesDeltaFromConfirmedSnapshots(t *testing.T) {
	stats := []*command.Stat{
		{Name: "user>>>1>>>traffic>>>downlink", Value: 1500},
		{Name: "user>>>1>>>traffic>>>uplink", Value: 450},
	}
	snapshots := map[string]uint64{
		"user>>>1>>>traffic>>>downlink": 1000,
		"user>>>1>>>traffic>>>uplink":   100,
	}

	plan := BuildTrafficCollectionPlan("10.0.0.1", 20260615, stats, snapshots, tagSet("1"), tagSet("1"))

	require.Empty(t, plan.Alerts)
	require.Len(t, plan.Items, 1)
	require.Equal(t, uint64(500), plan.Items["1"].Down)
	require.Equal(t, uint64(350), plan.Items["1"].Up)
	require.Equal(t, uint64(850), plan.UsedQuota)
	require.Equal(t, uint64(1500), plan.Snapshots["user>>>1>>>traffic>>>downlink"])
	require.Equal(t, uint64(450), plan.Snapshots["user>>>1>>>traffic>>>uplink"])
}

func TestBuildTrafficCollectionPlanSeedsMissingSnapshotsWithoutDelta(t *testing.T) {
	stats := []*command.Stat{
		{Name: "user>>>1>>>traffic>>>downlink", Value: 1500},
		{Name: "user>>>1>>>traffic>>>uplink", Value: 450},
	}

	plan := BuildTrafficCollectionPlan("10.0.0.1", 20260615, stats, nil, tagSet("1"), tagSet("1"))

	require.Empty(t, plan.Alerts)
	require.Empty(t, plan.Items)
	require.Equal(t, uint64(0), plan.UsedQuota)
	require.Equal(t, uint64(1500), plan.Snapshots["user>>>1>>>traffic>>>downlink"])
	require.Equal(t, uint64(450), plan.Snapshots["user>>>1>>>traffic>>>uplink"])
}

func TestBuildTrafficCollectionPlanTreatsLowerCounterAsXrayReset(t *testing.T) {
	stats := []*command.Stat{
		{Name: "user>>>1>>>traffic>>>downlink", Value: 200},
	}
	snapshots := map[string]uint64{
		"user>>>1>>>traffic>>>downlink": 1500,
	}

	plan := BuildTrafficCollectionPlan("10.0.0.1", 20260615, stats, snapshots, tagSet("1"), tagSet("1"))

	require.Len(t, plan.Alerts, 1)
	require.Equal(t, TrafficAnomalyCounterReset, plan.Alerts[0].Reason)
	require.Equal(t, uint64(200), plan.Items["1"].Down)
	require.Equal(t, uint64(200), plan.Snapshots["user>>>1>>>traffic>>>downlink"])
}

func TestBuildTrafficCollectionPlanAlertsUnknownStatsWithoutDroppingKnownUsers(t *testing.T) {
	stats := []*command.Stat{
		{Name: "user>>>ghost>>>traffic>>>downlink", Value: 900},
		{Name: "not-a-valid-stat-name", Value: 700},
		{Name: "user>>>1>>>traffic>>>downlink", Value: 300},
	}
	snapshots := map[string]uint64{
		"user>>>1>>>traffic>>>downlink": 100,
	}

	plan := BuildTrafficCollectionPlan("10.0.0.1", 20260615, stats, snapshots, tagSet("1"), tagSet("1"))

	require.Len(t, plan.Items, 1)
	require.Equal(t, uint64(200), plan.Items["1"].Down)
	require.Len(t, plan.Alerts, 2)
	require.Equal(t, TrafficAnomalyUnknownUser, plan.Alerts[0].Reason)
	require.Equal(t, TrafficAnomalyInvalidStatName, plan.Alerts[1].Reason)
	require.Equal(t, uint64(900), plan.Snapshots["user>>>ghost>>>traffic>>>downlink"])
	require.Equal(t, uint64(700), plan.Snapshots["not-a-valid-stat-name"])
}

func TestBuildTrafficCollectionPlanCountsKnownUserWithInactiveBindingAndAlerts(t *testing.T) {
	stats := []*command.Stat{
		{Name: "user>>>2>>>traffic>>>uplink", Value: 1024},
	}
	snapshots := map[string]uint64{
		"user>>>2>>>traffic>>>uplink": 24,
	}

	plan := BuildTrafficCollectionPlan("10.0.0.1", 20260615, stats, snapshots, tagSet("2"), tagSet("1"))

	require.Len(t, plan.Alerts, 1)
	require.Equal(t, TrafficAnomalyInactiveBinding, plan.Alerts[0].Reason)
	require.Equal(t, uint64(1000), plan.Items["2"].Up)
}

func TestBuildXRayConfigFromBindingsUsesOnlyValidDatabaseBindings(t *testing.T) {
	firstUUID := uuid.NewV4()
	secondUUID := uuid.NewV4()
	bindings := []*v2ray.Binding{
		{
			UserID:  1,
			User:    system.SysUser{UUID: firstUUID},
			AlterID: 0,
			Level:   0,
		},
		{
			UserID:  0,
			User:    system.SysUser{UUID: uuid.Nil},
			AlterID: 0,
			Level:   0,
		},
		{
			UserID:  2,
			User:    system.SysUser{UUID: secondUUID},
			AlterID: 4,
			Level:   1,
		},
	}

	config, alerts := BuildXRayConfigFromBindings(443, bindings)

	require.Len(t, alerts, 1)
	require.Equal(t, TrafficAnomalyInvalidBinding, alerts[0].Reason)
	require.Len(t, config.InboundConfigs[0].Settings.Clients, 2)
	require.Equal(t, firstUUID.String(), config.InboundConfigs[0].Settings.Clients[0].Id)
	require.Equal(t, "1", config.InboundConfigs[0].Settings.Clients[0].Email)
	require.Equal(t, secondUUID.String(), config.InboundConfigs[0].Settings.Clients[1].Id)
	require.Equal(t, "2", config.InboundConfigs[0].Settings.Clients[1].Email)
	require.True(t, config.Policy.Levels["0"].StatsUserDownlink)
	require.True(t, config.Policy.Levels["0"].StatsUserUplink)
	require.True(t, config.Policy.Levels["1"].StatsUserDownlink)
	require.True(t, config.Policy.Levels["1"].StatsUserUplink)
}

func tagSet(tags ...string) map[string]struct{} {
	out := make(map[string]struct{}, len(tags))
	for _, tag := range tags {
		out[tag] = struct{}{}
	}
	return out
}
