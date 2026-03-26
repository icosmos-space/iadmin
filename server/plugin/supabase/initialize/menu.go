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
			Path:      "supabaseUser",
			Name:      "supabaseUser",
			Hidden:    false,
			Component: "plugin/supabase/view/user.vue",
			Meta:      model.Meta{Title: "Supabase用户", Icon: "user"},
		},
	}
	utils.RegisterMenus(entities...)
}
