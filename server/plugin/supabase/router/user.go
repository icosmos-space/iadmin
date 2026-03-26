package router

import (
	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/middleware"
)

type user struct{}

func (r *user) Init(public *gin.RouterGroup, private *gin.RouterGroup) {
	{
		group := private.Group("supabaseUser").Use(middleware.OperationRecord())
		group.PUT("updatePassword", apiUser.UpdateSupabaseUserPassword)
	}
	{
		group := private.Group("supabaseUser")
		group.GET("getUserList", apiUser.GetSupabaseUserList)
	}
}
