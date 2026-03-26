package initialize

import (
	"fmt"

	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/plugin/supabase/plugin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Viper() {
	err := global.IADMIN_VP.UnmarshalKey("supabase", &plugin.Config)
	if err != nil {
		err = errors.Wrap(err, "初始化 Supabase 配置失败!")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
