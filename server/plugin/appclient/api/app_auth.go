package api

import (
	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/model/common/response"
	appReq "github.com/icosmos-space/iadmin/server/plugin/appclient/model/request"
	"github.com/icosmos-space/iadmin/server/plugin/appclient/service"
)

type appAuth struct{}

// Register 注册
func (a *appAuth) Register(c *gin.Context) {
	var req appReq.AppRegister
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 可按需增加验证码等校验
	res, err := service.AppUserServiceApp.Register(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "注册成功", c)
}

// Login 登录
func (a *appAuth) Login(c *gin.Context) {
	var req appReq.AppLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := service.AppUserServiceApp.Login(req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(res, "登录成功", c)
}
