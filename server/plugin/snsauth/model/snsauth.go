package model

import (
	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/model/common"
)

type SnsProviderConfig struct {
	global.GVA_MODEL
	Provider     string `json:"provider" gorm:"uniqueIndex;size:32;comment:平台标识"`
	Enabled      bool   `json:"enabled" gorm:"comment:是否启用"`
	ClientID     string `json:"clientID" gorm:"size:255;comment:Client ID"`
	ClientSecret string `json:"clientSecret" gorm:"size:255;comment:Client Secret"`
	RedirectURL  string `json:"redirectURL" gorm:"size:512;comment:回调地址"`
	Scopes       string `json:"scopes" gorm:"size:512;comment:授权范围"`
	AuthURL      string `json:"authURL" gorm:"size:512;comment:授权地址"`
	TokenURL     string `json:"tokenURL" gorm:"size:512;comment:Token地址"`
	UserInfoURL  string `json:"userInfoURL" gorm:"size:512;comment:用户信息地址"`
}

type SnsUserBind struct {
	global.GVA_MODEL
	UserID   uint           `json:"userID" gorm:"index;comment:系统用户ID"`
	Provider string         `json:"provider" gorm:"index:idx_provider_openid,priority:1;size:32;comment:平台"`
	OpenID   string         `json:"openID" gorm:"index:idx_provider_openid,priority:2;size:255;comment:第三方用户ID"`
	UnionID  string         `json:"unionID" gorm:"size:255;comment:第三方UnionID"`
	Meta     common.JSONMap `json:"meta" gorm:"type:text;comment:第三方原始信息"`
}

func (SnsProviderConfig) TableName() string {
	return "sns_provider_configs"
}

func (SnsUserBind) TableName() string {
	return "sns_user_binds"
}
