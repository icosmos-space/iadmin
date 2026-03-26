package api

import (
	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/model/common/response"
	supaReq "github.com/icosmos-space/iadmin/server/plugin/supabase/model/request"
	"go.uber.org/zap"
)

type user struct{}

// GetSupabaseUserList 获取 Supabase 用户列表
func (a *user) GetSupabaseUserList(c *gin.Context) {
	var pageInfo supaReq.SupabaseUserSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceUser.GetUserList(pageInfo)
	if err != nil {
		global.IADMIN_LOG.Error("获取 Supabase 用户列表失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// UpdateSupabaseUserPassword 修改 Supabase 用户密码
func (a *user) UpdateSupabaseUserPassword(c *gin.Context) {
	var data supaReq.UpdateSupabaseUserPassword
	if err := c.ShouldBindJSON(&data); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceUser.UpdateUserPassword(data); err != nil {
		global.IADMIN_LOG.Error("修改 Supabase 用户密码失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("修改成功", c)
}
