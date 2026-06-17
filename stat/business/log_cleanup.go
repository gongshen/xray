package business

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const defaultLogCleanupInterval = 24 * time.Hour

var xrayLogFileNames = []string{"access.log", "error.log"}

func CleanupOldDateDirs(root string, now time.Time, retentionMonths int) ([]string, error) {
	if root == "" || retentionMonths <= 0 {
		return nil, nil
	}
	entries, err := os.ReadDir(root)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	cutoff := dateOnly(now.AddDate(0, -retentionMonths, 0))
	removed := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dirDate, err := time.ParseInLocation("2006-01-02", entry.Name(), now.Location())
		if err != nil {
			continue
		}
		if !dateOnly(dirDate).Before(cutoff) {
			continue
		}
		path := filepath.Join(root, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return removed, err
		}
		removed = append(removed, path)
	}
	return removed, nil
}

func StartLogDirectoryCleaner(root string, retentionMonths int, interval time.Duration, stop <-chan struct{}) {
	if interval <= 0 {
		interval = defaultLogCleanupInterval
	}
	runLogDirectoryCleanup(root, retentionMonths)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			runLogDirectoryCleanup(root, retentionMonths)
		}
	}
}

func CleanupXrayLogFiles(logDir string, now time.Time, retentionMonths int) (int, error) {
	return CleanupOldRotatedLogFiles(logDir, xrayLogFileNames, now, retentionMonths)
}

func CleanupOldRotatedLogFiles(logDir string, activeNames []string, now time.Time, retentionMonths int) (int, error) {
	if logDir == "" || retentionMonths <= 0 {
		return 0, nil
	}
	entries, err := os.ReadDir(logDir)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	cutoff := dateOnly(now.AddDate(0, -retentionMonths, 0))
	removed := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		baseName, ok := rotatedLogBaseName(name, activeNames)
		if !ok {
			continue
		}

		remove := false
		if rotatedAt, ok := rotatedLogDate(name, baseName, now.Location()); ok {
			remove = dateOnly(rotatedAt).Before(cutoff)
		} else {
			if !rotatedLogSequence(name, baseName) {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				return removed, err
			}
			remove = dateOnly(info.ModTime()).Before(cutoff)
		}
		if !remove {
			continue
		}
		if err := os.Remove(filepath.Join(logDir, name)); err != nil {
			return removed, err
		}
		removed++
	}
	return removed, nil
}

func rotatedLogBaseName(name string, activeNames []string) (string, bool) {
	for _, activeName := range activeNames {
		if name == activeName {
			return "", false
		}
		if strings.HasPrefix(name, activeName+".") || strings.HasPrefix(name, activeName+"-") {
			return activeName, true
		}
	}
	return "", false
}

func rotatedLogDate(name string, baseName string, location *time.Location) (time.Time, bool) {
	suffix := rotatedLogSuffix(name, baseName)
	for _, layout := range []string{"20060102", "2006-01-02"} {
		rotatedAt, err := time.ParseInLocation(layout, suffix, location)
		if err != nil {
			continue
		}
		return rotatedAt, true
	}
	return time.Time{}, false
}

func rotatedLogSequence(name string, baseName string) bool {
	suffix := rotatedLogSuffix(name, baseName)
	if suffix == "" {
		return false
	}
	for _, char := range suffix {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func rotatedLogSuffix(name string, baseName string) string {
	suffix := strings.TrimPrefix(name, baseName)
	for _, ext := range []string{".gz", ".xz", ".zst", ".bz2", ".zip"} {
		suffix = strings.TrimSuffix(suffix, ext)
	}
	return strings.TrimLeft(suffix, ".-")
}

func StartXrayLogFileCleaner(logDir string, retentionMonths int, interval time.Duration, stop <-chan struct{}) {
	if interval <= 0 {
		interval = defaultLogCleanupInterval
	}
	runXrayLogFileCleanup(logDir, retentionMonths)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			runXrayLogFileCleanup(logDir, retentionMonths)
		}
	}
}

func runLogDirectoryCleanup(root string, retentionMonths int) {
	removed, err := CleanupOldDateDirs(root, time.Now(), retentionMonths)
	if err != nil {
		logrus.WithError(err).Warn("log directory cleanup failed")
		return
	}
	if len(removed) > 0 {
		logrus.WithField("count", len(removed)).Info("old log directories cleaned")
	}
}

func runXrayLogFileCleanup(logDir string, retentionMonths int) {
	removed, err := CleanupXrayLogFiles(logDir, time.Now(), retentionMonths)
	if err != nil {
		logrus.WithError(err).Warn("xray log cleanup failed")
		return
	}
	if removed > 0 {
		logrus.WithField("removed_files", removed).Info("old xray rotated log files cleaned")
	}
}

func dateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
