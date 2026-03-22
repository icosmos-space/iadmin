package core

import (
	"fmt"
	"time"

	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/initialize"
	"github.com/icosmos-space/iadmin/server/service/system"
	"go.uber.org/zap"
)

func RunServer() {
	if global.IADMIN_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
		if global.IADMIN_CONFIG.System.UseMultipoint {
			initialize.RedisList()
		}
	}

	if global.IADMIN_CONFIG.System.UseMongo {
		err := initialize.Mongo.Initialization()
		if err != nil {
			zap.L().Error(fmt.Sprintf("%+v", err))
		}
	}
	// 从db加载jwt数据
	if global.IADMIN_DB != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()

	address := fmt.Sprintf(":%d", global.IADMIN_CONFIG.System.Addr)

	printStartupBanner(address)
	initServer(address, Router, 10*time.Minute, 10*time.Minute)
}
