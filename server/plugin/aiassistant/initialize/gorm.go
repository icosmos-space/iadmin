package initialize

import (
	"context"
	"fmt"

	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/plugin/aiassistant/model"
	"github.com/icosmos-space/iadmin/server/plugin/aiassistant/service"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func Gorm(ctx context.Context) {
	err := global.IADMIN_DB.WithContext(ctx).AutoMigrate(
		new(model.AIAssistantConfig),
		new(model.AIAssistantPrompt),
	)
	if err != nil {
		err = errors.Wrap(err, "register ai assistant tables failed")
		zap.L().Error(fmt.Sprintf("%+v", err))
		return
	}
	if err = service.Service.Assistant.InitDefaults(); err != nil {
		zap.L().Error("init ai assistant default data failed", zap.Error(err))
	}
}
