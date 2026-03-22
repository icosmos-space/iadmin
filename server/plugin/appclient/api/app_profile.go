package api

import (
	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/middleware"
	"github.com/icosmos-space/iadmin/server/model/common/response"
	"github.com/icosmos-space/iadmin/server/plugin/appclient/service"
)

type appProfile struct{}

// Profile 当前用户信息（需 x-app-token）
func (a *appProfile) Profile(c *gin.Context) {
	uid := middleware.GetAppUserID(c)
	if uid == 0 {
		response.FailWithMessage("未登录", c)
		return
	}
	u, err := service.AppUserServiceApp.GetByID(uid)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(u, "ok", c)
}
