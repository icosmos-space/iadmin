package initialize

import (
	"context"

	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/plugin/appclient/model"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	if err := global.IADMIN_DB.WithContext(ctx).AutoMigrate(new(model.AppUser)); err != nil {
		zap.L().Error("appclient AutoMigrate failed", zap.Error(err))
	}
}
