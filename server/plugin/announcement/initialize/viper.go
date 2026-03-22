package initialize

import (
	"fmt"

	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/plugin/announcement/plugin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Viper() {
	err := global.IADMIN_VP.UnmarshalKey("announcement", &plugin.Config)
	if err != nil {
		err = errors.Wrap(err, "初始化配置文件失败!")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
