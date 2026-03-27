package router

import (
	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/middleware"
)

type auth struct{}

func (r *auth) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	{
		g := public.Group("snsAuth")
		g.GET("getEnabledProviderList", apiAuth.GetEnabledProviderList)
		g.GET("getLoginURL", apiAuth.GetLoginURL)
		g.GET("callback/:provider", apiAuth.Callback)
		g.POST("telegramLogin", apiAuth.TelegramLogin)
	}
	{
		g := private.Group("snsAuth")
		g.GET("getProviderList", apiAuth.GetProviderList)
		g.GET("getBindURL", apiAuth.GetBindURL)
		g.GET("getMyBindings", apiAuth.GetMyBindings)
		g.DELETE("unbind", apiAuth.Unbind)
		g.POST("telegramBind", apiAuth.TelegramBind)
	}
	{
		g := private.Group("snsAuth").Use(middleware.OperationRecord())
		g.PUT("updateProviderConfig", apiAuth.UpdateProviderConfig)
	}
}
