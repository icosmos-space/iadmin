//go:build !dev

package system

import "github.com/gin-gonic/gin"

func registerMcpAutoCodeRoutes(*gin.RouterGroup) {}
