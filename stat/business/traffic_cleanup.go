package business

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

const defaultTrafficEventCleanupInterval = 24 * time.Hour

func (store *TrafficStore) CleanupOldTrafficEvents(now time.Time, retentionMonths int) (int64, error) {
	if store == nil || store.db == nil {
		return 0, fmt.Errorf("traffic store is not initialized")
	}
	if retentionMonths <= 0 {
		return 0, nil
	}
	cutoff := now.AddDate(0, -retentionMonths, 0).Unix()
	result, err := store.db.Exec(`DELETE FROM traffic_event WHERE collected_at < ?`, cutoff)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func StartTrafficEventCleaner(store *TrafficStore, retentionMonths int, interval time.Duration, stop <-chan struct{}) {
	if interval <= 0 {
		interval = defaultTrafficEventCleanupInterval
	}
	runTrafficEventCleanup(store, retentionMonths)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			runTrafficEventCleanup(store, retentionMonths)
		}
	}
}

func runTrafficEventCleanup(store *TrafficStore, retentionMonths int) {
	removed, err := store.CleanupOldTrafficEvents(time.Now(), retentionMonths)
	if err != nil {
		logrus.WithError(err).Warn("traffic event cleanup failed")
		return
	}
	if removed > 0 {
		logrus.WithField("removed_events", removed).Info("old traffic events cleaned")
	}
}
