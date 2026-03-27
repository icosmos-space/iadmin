package config

type ProviderConfig struct {
	Enabled      bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`
	ClientID     string `mapstructure:"client-id" json:"clientID" yaml:"client-id"`
	ClientSecret string `mapstructure:"client-secret" json:"clientSecret" yaml:"client-secret"`
	RedirectURL  string `mapstructure:"redirect-url" json:"redirectURL" yaml:"redirect-url"`
	Scopes       string `mapstructure:"scopes" json:"scopes" yaml:"scopes"`
	AuthURL      string `mapstructure:"auth-url" json:"authURL" yaml:"auth-url"`
	TokenURL     string `mapstructure:"token-url" json:"tokenURL" yaml:"token-url"`
	UserInfoURL  string `mapstructure:"user-info-url" json:"userInfoURL" yaml:"user-info-url"`
}

type Config struct {
	Github   ProviderConfig `mapstructure:"github" json:"github" yaml:"github"`
	Feishu   ProviderConfig `mapstructure:"feishu" json:"feishu" yaml:"feishu"`
	Wechat   ProviderConfig `mapstructure:"wechat" json:"wechat" yaml:"wechat"`
	Telegram ProviderConfig `mapstructure:"telegram" json:"telegram" yaml:"telegram"`
}
