package initialize

import (
	"context"
	model "github.com/icosmos-space/iadmin/server/model/system"
	"github.com/icosmos-space/iadmin/server/plugin/plugin-tool/utils"
)

func Api(ctx context.Context) {
	entities := []model.SysApi{
		{
			Path:        "/supabaseUser/getUserList",
			Description: "获取 Supabase 用户列表",
			ApiGroup:    "Supabase用户",
			Method:      "GET",
		},
		{
			Path:        "/supabaseUser/updatePassword",
			Description: "修改 Supabase 用户密码",
			ApiGroup:    "Supabase用户",
			Method:      "PUT",
		},
	}
	utils.RegisterApis(entities...)
}
