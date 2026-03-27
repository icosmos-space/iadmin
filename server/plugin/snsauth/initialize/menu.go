package initialize

import (
	"context"

	model "github.com/icosmos-space/iadmin/server/model/system"
	"github.com/icosmos-space/iadmin/server/plugin/plugin-tool/utils"
)

func Menu(ctx context.Context) {
	entities := []model.SysBaseMenu{
		{
			ParentId:  0,
			Path:      "snsAuthConfig",
			Name:      "snsAuthConfig",
			Hidden:    false,
			Component: "plugin/snsauth/view/config.vue",
			Meta:      model.Meta{Title: "SNS登录配置", Icon: "link"},
		},
	}
	utils.RegisterMenus(entities...)
}
