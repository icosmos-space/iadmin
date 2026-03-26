package request

import "github.com/icosmos-space/iadmin/server/model/common/request"

type SupabaseUserSearch struct {
	request.PageInfo
}

type UpdateSupabaseUserPassword struct {
	UserID      string `json:"userID" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=72"`
}
