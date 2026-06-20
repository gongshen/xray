package v2ray

type TrafficTargetClassificationCache struct {
	ID          uint   `gorm:"primarykey"`
	Target      string `json:"target" form:"target" gorm:"column:target;type:varchar(255);uniqueIndex"`
	ServiceName string `json:"service_name" form:"service_name" gorm:"column:service_name;type:varchar(255);index"`
	Purpose     string `json:"purpose" form:"purpose" gorm:"column:purpose;type:text"`
	CreatedAt   int64  `json:"created_at" form:"created_at" gorm:"column:created_at;index"`
	UpdatedAt   int64  `json:"updated_at" form:"updated_at" gorm:"column:updated_at;index"`
}

func (TrafficTargetClassificationCache) TableName() string {
	return "v2ray_traffic_target_classification_cache"
}
