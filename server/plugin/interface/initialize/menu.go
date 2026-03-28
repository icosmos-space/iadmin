package initialize

import (
	"context"

	model "github.com/icosmos-space/iadmin/server/model/system"
	"github.com/icosmos-space/iadmin/server/plugin/plugin-tool/utils"
)

func Menu(ctx context.Context) {
	entities := []model.SysBaseMenu{
		{
			ParentId:  8,
			Path:      "interface",
			Name:      "interface",
			Hidden:    false,
			Component: "plugin/interface/view/swagger.vue",
			Sort:      0,
			Meta:      model.Meta{Title: "接口管理", Icon: "tickets"},
		},
	}
	utils.RegisterMenus(entities...)
}
