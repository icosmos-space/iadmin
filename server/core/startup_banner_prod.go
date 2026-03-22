//go:build !dev

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
	默认自动化文档: 生产构建未包含 Swagger UI（本地开发请 go run -tags=dev）
	默认MCP: 生产构建未包含内嵌 MCP（本地开发请 go run -tags=dev）
	默认前端文件运行地址:http://127.0.0.1:8080
	--------------------------------------版权声明--------------------------------------
	** icosmos-space开源团队 **
	** 感谢您对iAdmin的支持与关注**
`, global.Version, address)
}
