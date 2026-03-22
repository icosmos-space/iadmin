package initialize

import (
	"github.com/icosmos-space/iadmin/server/global"
	mcpTool "github.com/icosmos-space/iadmin/server/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func McpRun() *server.SSEServer {
	config := global.IADMIN_CONFIG.MCP

	s := server.NewMCPServer(
		config.Name,
		config.Version,
	)

	global.IADMIN_MCP_SERVER = s

	mcpTool.RegisterAllTools(s)

	return server.NewSSEServer(s,
		server.WithSSEEndpoint(config.SSEPath),
		server.WithMessageEndpoint(config.MessagePath),
		server.WithBaseURL(config.UrlPrefix))
}
