//go:build dev

package core

import (
	"fmt"

	"github.com/icosmos-space/iadmin/server/global"
)

func printStartupBanner(address string) {
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
}
