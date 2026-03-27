package initialize

import (
	"context"

	model "github.com/icosmos-space/iadmin/server/model/system"
	"github.com/icosmos-space/iadmin/server/plugin/plugin-tool/utils"
)

func Api(ctx context.Context) {
	entities := []model.SysApi{
		{Path: "/snsAuth/getEnabledProviderList", Description: "获取已启用SNS平台", ApiGroup: "SNS认证", Method: "GET"},
		{Path: "/snsAuth/getProviderList", Description: "获取SNS平台配置", ApiGroup: "SNS认证", Method: "GET"},
		{Path: "/snsAuth/updateProviderConfig", Description: "更新SNS平台配置", ApiGroup: "SNS认证", Method: "PUT"},
		{Path: "/snsAuth/getLoginURL", Description: "获取SNS登录地址", ApiGroup: "SNS认证", Method: "GET"},
		{Path: "/snsAuth/getBindURL", Description: "获取SNS绑定地址", ApiGroup: "SNS认证", Method: "GET"},
		{Path: "/snsAuth/callback/:provider", Description: "SNS授权回调", ApiGroup: "SNS认证", Method: "GET"},
		{Path: "/snsAuth/telegramLogin", Description: "Telegram登录", ApiGroup: "SNS认证", Method: "POST"},
		{Path: "/snsAuth/telegramBind", Description: "Telegram绑定", ApiGroup: "SNS认证", Method: "POST"},
		{Path: "/snsAuth/getMyBindings", Description: "获取我的绑定", ApiGroup: "SNS认证", Method: "GET"},
		{Path: "/snsAuth/unbind", Description: "解绑SNS账号", ApiGroup: "SNS认证", Method: "DELETE"},
	}
	utils.RegisterApis(entities...)
}
