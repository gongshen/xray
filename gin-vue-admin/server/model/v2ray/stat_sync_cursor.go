package v2ray

// StatSyncCursor stores the last local stat event synced from a node.
type StatSyncCursor struct {
	ID          uint   `gorm:"primarykey"`
	ServerID    uint   `json:"server_id" form:"server_id" gorm:"column:server_id;uniqueIndex"`
	ServerIp    string `json:"server_ip" form:"server_ip" gorm:"column:server_ip;type:varchar(64);index"`
	LastEventID uint64 `json:"last_event_id" form:"last_event_id" gorm:"column:last_event_id"`
	UpdatedAt   int64  `json:"updated_at" form:"updated_at" gorm:"column:updated_at;index"`
}

func (StatSyncCursor) TableName() string {
	return "v2ray_stat_sync_cursor"
}
