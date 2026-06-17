package business

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestCleanupOldDateDirsRemovesOnlyDirsOlderThanRetention(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, "2025-12-16"))
	mustMkdir(t, filepath.Join(root, "2025-12-17"))
	mustMkdir(t, filepath.Join(root, "2025-12-18"))
	mustMkdir(t, filepath.Join(root, "not-a-date"))
	mustWriteFile(t, filepath.Join(root, "2025-01-01"), "not a directory")

	removed, err := CleanupOldDateDirs(root, time.Date(2026, 6, 17, 12, 0, 0, 0, time.Local), 6)
	if err != nil {
		t.Fatalf("CleanupOldDateDirs returned error: %v", err)
	}

	sort.Strings(removed)
	expectedRemoved := []string{filepath.Join(root, "2025-12-16")}
	if !reflect.DeepEqual(removed, expectedRemoved) {
		t.Fatalf("removed = %#v, want %#v", removed, expectedRemoved)
	}

	assertNotExists(t, filepath.Join(root, "2025-12-16"))
	assertExists(t, filepath.Join(root, "2025-12-17"))
	assertExists(t, filepath.Join(root, "2025-12-18"))
	assertExists(t, filepath.Join(root, "not-a-date"))
	assertExists(t, filepath.Join(root, "2025-01-01"))
}

func TestCleanupOldDateDirsIgnoresMissingRoot(t *testing.T) {
	removed, err := CleanupOldDateDirs(filepath.Join(t.TempDir(), "missing"), time.Now(), 6)
	if err != nil {
		t.Fatalf("CleanupOldDateDirs returned error for missing root: %v", err)
	}
	if len(removed) != 0 {
		t.Fatalf("removed = %#v, want empty", removed)
	}
}

func TestCleanupXrayLogFilesRemovesOnlyOldRotatedFiles(t *testing.T) {
	logDir := t.TempDir()
	now := time.Date(2026, 6, 17, 12, 0, 0, 0, time.Local)

	mustWriteFile(t, filepath.Join(logDir, "access.log"), "2025/01/01 active access log must stay\n")
	mustWriteFile(t, filepath.Join(logDir, "error.log"), "2025/01/01 active error log must stay\n")
	mustWriteFile(t, filepath.Join(logDir, "access.log-20250616"), "old rotated access\n")
	mustWriteFile(t, filepath.Join(logDir, "access.log-20250617"), "boundary rotated access\n")
	mustWriteFile(t, filepath.Join(logDir, "error.log-2025-06-16.gz"), "old rotated error\n")
	mustWriteFile(t, filepath.Join(logDir, "error.log.1.gz"), "old numeric rotated error\n")
	mustWriteFile(t, filepath.Join(logDir, "other.log-20250616"), "unrelated log\n")
	mustWriteFile(t, filepath.Join(logDir, "access.log.backup"), "manual backup\n")
	mustChtimes(t, filepath.Join(logDir, "error.log.1.gz"), now.AddDate(0, -12, -1))
	mustChtimes(t, filepath.Join(logDir, "access.log.backup"), now.AddDate(0, -12, -1))

	removed, err := CleanupXrayLogFiles(logDir, now, 12)
	if err != nil {
		t.Fatalf("CleanupXrayLogFiles returned error: %v", err)
	}
	if removed != 3 {
		t.Fatalf("removed = %d, want 3", removed)
	}

	assertExists(t, filepath.Join(logDir, "access.log"))
	assertExists(t, filepath.Join(logDir, "error.log"))
	assertNotExists(t, filepath.Join(logDir, "access.log-20250616"))
	assertExists(t, filepath.Join(logDir, "access.log-20250617"))
	assertNotExists(t, filepath.Join(logDir, "error.log-2025-06-16.gz"))
	assertNotExists(t, filepath.Join(logDir, "error.log.1.gz"))
	assertExists(t, filepath.Join(logDir, "other.log-20250616"))
	assertExists(t, filepath.Join(logDir, "access.log.backup"))

	gotBytes, err := os.ReadFile(filepath.Join(logDir, "access.log"))
	if err != nil {
		t.Fatalf("ReadFile(access.log): %v", err)
	}
	if !strings.Contains(string(gotBytes), "active access log must stay") {
		t.Fatalf("active access log was changed:\n%s", string(gotBytes))
	}
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("MkdirAll(%q): %v", path, err)
	}
}

func mustWriteFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("WriteFile(%q): %v", path, err)
	}
}

func mustChtimes(t *testing.T, path string, timestamp time.Time) {
	t.Helper()
	if err := os.Chtimes(path, timestamp, timestamp); err != nil {
		t.Fatalf("Chtimes(%q): %v", path, err)
	}
}

func assertExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected %q to exist: %v", path, err)
	}
}

func assertNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %q to be removed, stat err=%v", path, err)
	}
}
