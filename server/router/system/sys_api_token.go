package system

import (
	"github.com/icosmos-space/iadmin/server/api/v1"
	"github.com/icosmos-space/iadmin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ApiTokenRouter struct{}

func (s *ApiTokenRouter) InitApiTokenRouter(Router *gin.RouterGroup) {
	apiTokenRouter := Router.Group("sysApiToken").Use(middleware.OperationRecord())
	apiTokenApi := v1.ApiGroupApp.SystemApiGroup.ApiTokenApi
	{
		apiTokenRouter.POST("createApiToken", apiTokenApi.CreateApiToken)   // 签发Token
		apiTokenRouter.POST("getApiTokenList", apiTokenApi.GetApiTokenList) // 获取列表
		apiTokenRouter.POST("deleteApiToken", apiTokenApi.DeleteApiToken)   // 作废Token
	}
}
