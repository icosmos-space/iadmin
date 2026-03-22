package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/global"
)

func InstallPlugin(PrivateGroup *gin.RouterGroup, PublicRouter *gin.RouterGroup, engine *gin.Engine) {
	if global.IADMIN_DB == nil {
		global.IADMIN_LOG.Info("项目暂未初始化，无法安装插件，初始化后重启项目即可完成插件安装")
		return
	}
	bizPluginV1(PrivateGroup, PublicRouter)
	bizPluginV2(engine)
}
