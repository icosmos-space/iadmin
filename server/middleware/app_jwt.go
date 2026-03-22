package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/model/common/response"
	appReq "github.com/icosmos-space/iadmin/server/model/appclient/request"
	"github.com/icosmos-space/iadmin/server/utils"
)

// GetAppToken 从请求头读取 C 端 token：优先 x-app-token，其次 Authorization: Bearer
func GetAppToken(c *gin.Context) string {
	t := strings.TrimSpace(c.GetHeader("x-app-token"))
	if t != "" {
		return t
	}
	auth := strings.TrimSpace(c.GetHeader("Authorization"))
	if len(auth) > 7 && strings.EqualFold(auth[:7], "Bearer ") {
		return strings.TrimSpace(auth[7:])
	}
	return ""
}

// AppJWTAuth C 端鉴权，与后台 JWTAuth（x-token + Casbin）独立
func AppJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := GetAppToken(c)
		if token == "" {
			response.NoAuth("未登录或非法访问", c)
			c.Abort()
			return
		}
		j := utils.NewAppJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				response.NoAuth("登录已过期，请重新登录", c)
			} else {
				response.NoAuth(err.Error(), c)
			}
			c.Abort()
			return
		}
		c.Set("appClaims", claims)
		c.Next()
	}
}

// GetAppClaims 从 Context 读取 C 端 claims（需在 AppJWTAuth 之后调用）
func GetAppClaims(c *gin.Context) *appReq.AppClaims {
	if v, ok := c.Get("appClaims"); ok {
		if cl, ok := v.(*appReq.AppClaims); ok {
			return cl
		}
	}
	return nil
}

// GetAppUserID 从 Context 读取 C 端用户 ID
func GetAppUserID(c *gin.Context) uint {
	if cl := GetAppClaims(c); cl != nil {
		return cl.UserID
	}
	return 0
}
