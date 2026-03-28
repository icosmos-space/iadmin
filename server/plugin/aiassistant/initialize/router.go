package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/middleware"
	"github.com/icosmos-space/iadmin/server/plugin/aiassistant/router"
)

func Router(engine *gin.Engine) {
	public := engine.Group(global.IADMIN_CONFIG.System.RouterPrefix).Group("")
	private := engine.Group(global.IADMIN_CONFIG.System.RouterPrefix).Group("")
	private.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	router.Router.Assistant.Init(public, private)
}
