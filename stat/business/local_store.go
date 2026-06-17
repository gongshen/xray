package business

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type TrafficStore struct {
	db *sql.DB
}

func OpenTrafficStore(path string) (*TrafficStore, error) {
	if path == "" {
		return nil, fmt.Errorf("traffic store path is empty")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	store := &TrafficStore{db: db}
	if err := store.migrate(); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (store *TrafficStore) Close() error {
	if store == nil || store.db == nil {
		return nil
	}
	return store.db.Close()
}

func (store *TrafficStore) migrate() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS traffic_snapshot (
			name TEXT PRIMARY KEY,
			value INTEGER NOT NULL,
			updated_at INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS traffic_event (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tag TEXT NOT NULL,
			down INTEGER NOT NULL DEFAULT 0,
			up INTEGER NOT NULL DEFAULT 0,
			collected_at INTEGER NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_traffic_event_collected_at ON traffic_event(collected_at)`,
	}
	for _, stmt := range stmts {
		if _, err := store.db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

func (store *TrafficStore) SavePlan(plan LocalTrafficPlan) error {
	if store == nil || store.db == nil {
		return fmt.Errorf("traffic store is not initialized")
	}
	if len(plan.Events) == 0 && len(plan.Snapshots) == 0 {
		return nil
	}
	tx, err := store.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, event := range plan.Events {
		if event.Down == 0 && event.Up == 0 {
			continue
		}
		_, err := tx.Exec(
			`INSERT INTO traffic_event(tag, down, up, collected_at) VALUES(?, ?, ?, ?)`,
			event.Tag, int64(event.Down), int64(event.Up), event.CollectedAt,
		)
		if err != nil {
			return err
		}
	}

	updatedAt := time.Now().Unix()
	for name, value := range plan.Snapshots {
		_, err := tx.Exec(
			`INSERT INTO traffic_snapshot(name, value, updated_at) VALUES(?, ?, ?)
			 ON CONFLICT(name) DO UPDATE SET value=excluded.value, updated_at=excluded.updated_at`,
			name, int64(value), updatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (store *TrafficStore) LoadSnapshots(names []string) (map[string]uint64, error) {
	out := make(map[string]uint64)
	if store == nil || store.db == nil {
		return out, fmt.Errorf("traffic store is not initialized")
	}
	if len(names) == 0 {
		return out, nil
	}
	placeholders := strings.TrimRight(strings.Repeat("?,", len(names)), ",")
	args := make([]interface{}, 0, len(names))
	for _, name := range names {
		args = append(args, name)
	}
	rows, err := store.db.Query(`SELECT name, value FROM traffic_snapshot WHERE name IN (`+placeholders+`)`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value int64
		if err := rows.Scan(&name, &value); err != nil {
			return nil, err
		}
		out[name] = uint64(value)
	}
	return out, rows.Err()
}

func (store *TrafficStore) ListEventsAfter(afterID uint64, limit int) ([]LocalTrafficEvent, error) {
	if store == nil || store.db == nil {
		return nil, fmt.Errorf("traffic store is not initialized")
	}
	if limit <= 0 || limit > 5000 {
		limit = 1000
	}
	rows, err := store.db.Query(
		`SELECT id, tag, down, up, collected_at FROM traffic_event WHERE id > ? ORDER BY id ASC LIMIT ?`,
		int64(afterID), limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]LocalTrafficEvent, 0)
	for rows.Next() {
		var event LocalTrafficEvent
		var id int64
		var down int64
		var up int64
		if err := rows.Scan(&id, &event.Tag, &down, &up, &event.CollectedAt); err != nil {
			return nil, err
		}
		event.ID = uint64(id)
		event.Down = uint64(down)
		event.Up = uint64(up)
		events = append(events, event)
	}
	return events, rows.Err()
}
