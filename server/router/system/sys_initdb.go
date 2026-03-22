package system

import (
	"github.com/gin-gonic/gin"
)

type InitRouter struct{}

func (s *InitRouter) InitInitRouter(Router *gin.RouterGroup) {
	initRouter := Router.Group("init")
	{
		initRouter.GET("db-types", dbApi.InitSupportedDBTypes) // 可选数据库类型（与后端编译标签一致）
		initRouter.POST("initdb", dbApi.InitDB)                // 初始化数据库
		initRouter.POST("checkdb", dbApi.CheckDB)              // 检测是否需要初始化数据库
	}
}
