package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/plugin/appclient/router"
)

// Router 注册 C 端路由：{RouterPrefix}/app/v1/...
func Router(engine *gin.Engine) {
	base := global.IADMIN_CONFIG.System.RouterPrefix + "/app/v1"
	router.RouterGroupApp.App.InitStandalone(engine, base)
}
