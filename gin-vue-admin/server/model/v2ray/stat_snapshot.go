package v2ray

// StatSnapshot stores the last successfully persisted cumulative Xray counter.
type StatSnapshot struct {
	ID        uint   `gorm:"primarykey"`
	ServerIp  string `json:"server_ip" form:"server_ip" gorm:"column:server_ip;type:varchar(64);uniqueIndex:udx_stat_snapshot_server_name;index"`
	Name      string `json:"name" form:"name" gorm:"column:name;type:varchar(255);uniqueIndex:udx_stat_snapshot_server_name"`
	Value     uint64 `json:"value" form:"value" gorm:"column:value"`
	UpdatedAt int64  `json:"updated_at" form:"updated_at" gorm:"column:updated_at;index"`
}

func (StatSnapshot) TableName() string {
	return "v2ray_stat_snapshot"
}

// TrafficAnomaly records suspicious traffic stats that need operator attention.
type TrafficAnomaly struct {
	ID        uint   `gorm:"primarykey"`
	ServerIp  string `json:"server_ip" form:"server_ip" gorm:"column:server_ip;type:varchar(64);index"`
	Reason    string `json:"reason" form:"reason" gorm:"column:reason;type:varchar(64);index"`
	Name      string `json:"name" form:"name" gorm:"column:name;type:varchar(255);index"`
	Tag       string `json:"tag" form:"tag" gorm:"column:tag;type:varchar(64);index"`
	Value     int64  `json:"value" form:"value" gorm:"column:value"`
	Detail    string `json:"detail" form:"detail" gorm:"column:detail;type:varchar(255)"`
	CreatedAt int64  `json:"created_at" form:"created_at" gorm:"column:created_at;index"`
}

func (TrafficAnomaly) TableName() string {
	return "v2ray_traffic_anomaly"
}
