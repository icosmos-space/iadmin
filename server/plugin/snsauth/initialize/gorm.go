package initialize

import (
	"context"
	"fmt"

	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/plugin/snsauth/model"
	"github.com/icosmos-space/iadmin/server/plugin/snsauth/service"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	err := global.IADMIN_DB.WithContext(ctx).AutoMigrate(
		new(model.SnsProviderConfig),
		new(model.SnsUserBind),
	)
	if err != nil {
		err = errors.Wrap(err, "注册SNS表失败")
		zap.L().Error(fmt.Sprintf("%+v", err))
		return
	}
	if err = service.Service.Auth.InitDefaultProviders(); err != nil {
		zap.L().Error("初始化SNS默认平台失败", zap.Error(err))
	}
}
