package initialize

import (
	"sync"

	ginpprof "github.com/gin-contrib/pprof"

	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/global"
	"go.uber.org/zap"
)

var pprofOnce sync.Once

// RegisterPprof pprof 路由入口
func registerPprof(router *gin.Engine) {
	pprofOnce.Do(func() {
		ginpprof.Register(router, "debug/pprof")
		global.IADMIN_LOG.Info("pprof route registered", zap.String("path", "/debug/pprof/"))
	})
}
