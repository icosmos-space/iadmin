package config

type Config struct {
	URL            string `mapstructure:"url" json:"url" yaml:"url"`
	ServiceRoleKey string `mapstructure:"service-role-key" json:"serviceRoleKey" yaml:"service-role-key"`
}
