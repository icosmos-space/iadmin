package router

import (
	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/middleware"
	"github.com/icosmos-space/iadmin/server/plugin/appclient/api"
)

type appRouter struct{}

func (r *appRouter) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	auth := public.Group("auth")
	{
		auth.POST("register", api.ApiGroupApp.AppAuth.Register)
		auth.POST("login", api.ApiGroupApp.AppAuth.Login)
	}
	user := private.Group("user")
	{
		user.GET("profile", api.ApiGroupApp.AppProfile.Profile)
	}
}

// InitStandalone 用于在 plugin 中挂载：public 无鉴权，private 使用 AppJWTAuth（非 Casbin）
func (r *appRouter) InitStandalone(engine *gin.Engine, prefix string) {
	g := engine.Group(prefix)
	public := g.Group("")
	private := g.Group("")
	private.Use(middleware.AppJWTAuth())
	r.Init(public, private)
}
