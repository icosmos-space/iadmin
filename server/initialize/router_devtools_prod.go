//go:build !dev

package initialize

import "github.com/gin-gonic/gin"

func registerDevTools(*gin.Engine) {}
