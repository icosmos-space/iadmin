package router

import (
	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/middleware"
)

type assistant struct{}

func (r *assistant) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	{
		g := public.Group("aiAssistant").Use(middleware.JWTAuth())
		g.POST("chat", apiAssistant.Chat)
	}
	{
		g := public.Group("aiAssistant").Use(middleware.JWTAuth())
		g.GET("getEnabledPrompts", apiAssistant.GetEnabledPrompts)
	}
	{
		g := private.Group("aiAssistant")
		g.GET("getConfig", apiAssistant.GetConfig)
		g.GET("getPromptList", apiAssistant.GetPromptList)
		g.GET("findPrompt", apiAssistant.FindPrompt)
	}
	{
		g := private.Group("aiAssistant").Use(middleware.OperationRecord())
		g.PUT("updateConfig", apiAssistant.UpdateConfig)
		g.POST("createPrompt", apiAssistant.CreatePrompt)
		g.PUT("updatePrompt", apiAssistant.UpdatePrompt)
		g.DELETE("deletePrompt", apiAssistant.DeletePrompt)
		g.DELETE("deletePromptByIds", apiAssistant.DeletePromptByIDs)
	}
}
