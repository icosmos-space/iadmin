//go:build dev

package system

import "github.com/gin-gonic/gin"

func registerMcpAutoCodeRoutes(autoCodeRouter *gin.RouterGroup) {
	autoCodeRouter.POST("mcp", autoCodeTemplateApi.MCP)
	autoCodeRouter.POST("mcpList", autoCodeTemplateApi.MCPList)
	autoCodeRouter.POST("mcpTest", autoCodeTemplateApi.MCPTest)
}
