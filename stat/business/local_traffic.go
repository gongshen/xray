package business

import (
	"regexp"

	statsservice "github.com/xtls/xray-core/app/stats/command"
)

var localTrafficStatNameRegex = regexp.MustCompile(`^user>>>([^>]+)>>>traffic>>>(downlink|uplink)$`)

type LocalTrafficEvent struct {
	ID          uint64 `json:"id"`
	Tag         string `json:"tag"`
	Down        uint64 `json:"down"`
	Up          uint64 `json:"up"`
	CollectedAt int64  `json:"collected_at"`
}

type LocalTrafficPlan struct {
	Events    []LocalTrafficEvent
	Snapshots map[string]uint64
}

func BuildLocalTrafficPlan(stats []*statsservice.Stat, snapshots map[string]uint64, collectedAt int64) LocalTrafficPlan {
	if snapshots == nil {
		snapshots = map[string]uint64{}
	}
	plan := LocalTrafficPlan{
		Snapshots: make(map[string]uint64),
	}
	eventIndex := make(map[string]int)

	for _, stat := range stats {
		if stat == nil || stat.Value <= 0 {
			continue
		}
		matches := localTrafficStatNameRegex.FindStringSubmatch(stat.Name)
		if len(matches) != 3 {
			continue
		}

		currentValue := uint64(stat.Value)
		plan.Snapshots[stat.Name] = currentValue

		previousValue, seen := snapshots[stat.Name]
		if !seen {
			continue
		}

		delta := currentValue - previousValue
		if currentValue < previousValue {
			delta = currentValue
		}
		if delta == 0 {
			continue
		}

		tag := matches[1]
		idx, ok := eventIndex[tag]
		if !ok {
			plan.Events = append(plan.Events, LocalTrafficEvent{Tag: tag, CollectedAt: collectedAt})
			idx = len(plan.Events) - 1
			eventIndex[tag] = idx
		}
		if matches[2] == "downlink" {
			plan.Events[idx].Down += delta
		} else {
			plan.Events[idx].Up += delta
		}
	}

	return plan
}
