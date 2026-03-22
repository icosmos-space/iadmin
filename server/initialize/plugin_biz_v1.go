package initialize

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/plugin/email"
	"github.com/icosmos-space/iadmin/server/utils/plugin"
)

func PluginInit(group *gin.RouterGroup, Plugin ...plugin.Plugin) {
	for i := range Plugin {
		fmt.Println(Plugin[i].RouterPath(), "注册开始!")
		PluginGroup := group.Group(Plugin[i].RouterPath())
		Plugin[i].Register(PluginGroup)
		fmt.Println(Plugin[i].RouterPath(), "注册成功!")
	}
}

func bizPluginV1(group ...*gin.RouterGroup) {
	private := group[0]
	public := group[1]
	//  添加跟角色挂钩权限的插件 示例 本地示例模式于在线仓库模式注意上方的import 可以自行切换 效果相同
	PluginInit(private, email.CreateEmailPlug(
		global.IADMIN_CONFIG.Email.To,
		global.IADMIN_CONFIG.Email.From,
		global.IADMIN_CONFIG.Email.Host,
		global.IADMIN_CONFIG.Email.Secret,
		global.IADMIN_CONFIG.Email.Nickname,
		global.IADMIN_CONFIG.Email.Port,
		global.IADMIN_CONFIG.Email.IsSSL,
		global.IADMIN_CONFIG.Email.IsLoginAuth,
	))
	holder(public, private)
}
