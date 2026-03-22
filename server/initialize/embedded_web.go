//go:build embedweb

package initialize

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/global"
)

//go:embed all:resource/web/dist
var embeddedWebDist embed.FS

// registerEmbeddedWebUI 将 resource/web/dist 嵌入二进制并托管 SPA（/assets 与 index 回退）。
// 需在 API 路由注册之后调用；若配置了 system.router-prefix，则该前缀下的未命中路径返回 JSON 404，其余 GET/HEAD 回退到 index.html。
func registerEmbeddedWebUI(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		p := c.Request.URL.Path
		prefix := global.IADMIN_CONFIG.System.RouterPrefix
		if prefix != "" {
			if p == prefix || strings.HasPrefix(p, prefix+"/") {
				c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found"})
				c.Abort()
				return
			}
		}
		if strings.HasPrefix(p, "/assets/") {
			c.FileFromFS(strings.TrimPrefix(p, "/"), http.FS(embeddedWebDist))
			return
		}
		if p == "/favicon.ico" {
			if _, err := embeddedWebDist.Open("favicon.ico"); err == nil {
				c.FileFromFS("favicon.ico", http.FS(embeddedWebDist))
				return
			}
		}
		c.Header("Cache-Control", "no-cache")
		c.FileFromFS("index.html", http.FS(embeddedWebDist))
	})
}

// EmbeddedWebDistFS 供测试或其它包使用。
func EmbeddedWebDistFS() fs.FS {
	return embeddedWebDist
}
