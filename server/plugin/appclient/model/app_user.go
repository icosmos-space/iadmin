package model

import (
	"github.com/google/uuid"
	"github.com/icosmos-space/iadmin/server/global"
)

// AppUser C 端用户，与 system.SysUser 无关
type AppUser struct {
	global.GVA_MODEL
	UUID uuid.UUID `json:"uuid" gorm:"type:char(36);index"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:64;comment:登录名"`
	Password  string         `json:"-" gorm:"size:128;comment:bcrypt"`
	Nickname  string         `json:"nickname" gorm:"size:64;comment:昵称"`
	Phone     string         `json:"phone" gorm:"size:20;index;comment:手机"`
	Enable    int            `json:"enable" gorm:"default:1;comment:1正常 2禁用"`
}

func (AppUser) TableName() string {
	return "app_users"
}
