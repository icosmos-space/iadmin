//go:build !dev

package system

import (
	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/model/common/response"
)

func (a *AutoCodeTemplateApi) MCP(c *gin.Context) {
	response.FailWithMessage("MCP 相关接口仅在开发构建（go build -tags=dev）中启用", c)
}

func (a *AutoCodeTemplateApi) MCPList(c *gin.Context) {
	response.FailWithMessage("MCP 相关接口仅在开发构建（go build -tags=dev）中启用", c)
}

func (a *AutoCodeTemplateApi) MCPTest(c *gin.Context) {
	response.FailWithMessage("MCP 相关接口仅在开发构建（go build -tags=dev）中启用", c)
}
