package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
	"github.com/flipped-aurora/gin-vue-admin/server/router/v2ray"
)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	V2ray   v2ray.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
