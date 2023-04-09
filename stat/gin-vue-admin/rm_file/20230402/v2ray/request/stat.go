package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/v2ray"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type StatSearch struct{
    v2ray.Stat
    StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
    EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    request.PageInfo
}
