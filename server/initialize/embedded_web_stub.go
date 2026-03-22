//go:build !embedweb

package initialize

import "github.com/gin-gonic/gin"

func registerEmbeddedWebUI(*gin.Engine) {}
