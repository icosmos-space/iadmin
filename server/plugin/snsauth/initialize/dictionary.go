package initialize

import (
	"context"

	model "github.com/icosmos-space/iadmin/server/model/system"
	"github.com/icosmos-space/iadmin/server/plugin/plugin-tool/utils"
)

func Dictionary(ctx context.Context) {
	entities := []model.SysDictionary{}
	utils.RegisterDictionaries(entities...)
}
