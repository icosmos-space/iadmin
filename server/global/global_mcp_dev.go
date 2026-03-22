//go:build dev

package global

import "github.com/mark3labs/mcp-go/server"

// IADMIN_MCP_SERVER 仅在开发构建（-tags=dev）中由 initialize.McpRun 赋值。
var IADMIN_MCP_SERVER *server.MCPServer
