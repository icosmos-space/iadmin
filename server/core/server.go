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

	fmt.Printf(`
	欢迎使用 iAdmin
	当前版本:%s
	地址：https://github.com/icosmos-space/iadmin
	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
	默认MCP SSE地址:http://127.0.0.1%s%s
	默认MCP Message地址:http://127.0.0.1%s%s
	默认前端文件运行地址:http://127.0.0.1:8080
	--------------------------------------版权声明--------------------------------------
	** icosmos-space开源团队 **
	** 感谢您对iAdmin的支持与关注**
`, global.Version, address, address, global.IADMIN_CONFIG.MCP.SSEPath, address, global.IADMIN_CONFIG.MCP.MessagePath)
	initServer(address, Router, 10*time.Minute, 10*time.Minute)
}
