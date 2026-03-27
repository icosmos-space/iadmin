package request

type UpdateProviderConfigReq struct {
	Provider     string `json:"provider" binding:"required,oneof=github feishu wechat telegram"`
	Enabled      bool   `json:"enabled"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	RedirectURL  string `json:"redirectURL"`
	Scopes       string `json:"scopes"`
	AuthURL      string `json:"authURL"`
	TokenURL     string `json:"tokenURL"`
	UserInfoURL  string `json:"userInfoURL"`
}

type BuildURLReq struct {
	Provider string `form:"provider" binding:"required,oneof=github feishu wechat telegram"`
}

type TelegramAuthReq struct {
	ID        int64  `json:"id" binding:"required"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoURL  string `json:"photo_url"`
	AuthDate  int64  `json:"auth_date" binding:"required"`
	Hash      string `json:"hash" binding:"required"`
}
