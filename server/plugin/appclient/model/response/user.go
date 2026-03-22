package response

import (
	"github.com/google/uuid"
)

// AppUserPublic 返回给前端的用户信息（不含密码）
type AppUserPublic struct {
	ID       uint      `json:"id"`
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
	Nickname string    `json:"nickname"`
	Phone    string    `json:"phone"`
}

// AppLoginResponse 登录/注册返回
type AppLoginResponse struct {
	Token string         `json:"token"`
	User  AppUserPublic  `json:"user"`
}
