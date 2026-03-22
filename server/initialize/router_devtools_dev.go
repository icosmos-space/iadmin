//go:build dev

package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/docs"
	"github.com/icosmos-space/iadmin/server/global"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// registerDevTools 注册 Swagger UI 与内嵌 MCP（SSE）；仅在 go build -tags=dev 时编译进二进制。
func registerDevTools(Router *gin.Engine) {
	if !global.IADMIN_CONFIG.MCP.Separate {
		sseServer := McpRun()
		Router.GET(global.IADMIN_CONFIG.MCP.SSEPath, func(c *gin.Context) {
			sseServer.SSEHandler().ServeHTTP(c.Writer, c.Request)
		})
		Router.POST(global.IADMIN_CONFIG.MCP.MessagePath, func(c *gin.Context) {
			sseServer.MessageHandler().ServeHTTP(c.Writer, c.Request)
		})
	}
	docs.SwaggerInfo.BasePath = global.IADMIN_CONFIG.System.RouterPrefix
	Router.GET(global.IADMIN_CONFIG.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.IADMIN_LOG.Info("register swagger handler")
}
