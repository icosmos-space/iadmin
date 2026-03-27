package initialize

import (
	"fmt"

	"github.com/icosmos-space/iadmin/server/global"
	snsPlugin "github.com/icosmos-space/iadmin/server/plugin/snsauth/plugin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Viper() {
	err := global.IADMIN_VP.UnmarshalKey("sns-auth", &snsPlugin.Config)
	if err != nil {
		err = errors.Wrap(err, "初始化SNS认证配置失败")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
