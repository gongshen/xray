package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/v2ray"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/v2ray_admin"
)

type ApiGroup struct {
	SystemApiGroup     system.ApiGroup
	ExampleApiGroup    example.ApiGroup
	V2rayAdminApiGroup v2ray_admin.ApiGroup
	V2rayApiGroup      v2ray.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
