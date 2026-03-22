package initialize

import (
	"context"
	"fmt"

	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/plugin/announcement/model"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	err := global.IADMIN_DB.WithContext(ctx).AutoMigrate(
		new(model.Info),
	)
	if err != nil {
		err = errors.Wrap(err, "注册表失败!")
		zap.L().Error(fmt.Sprintf("%+v", err))
	}
}
