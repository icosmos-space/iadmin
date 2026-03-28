package model

import "github.com/icosmos-space/iadmin/server/global"

type AIAssistantConfig struct {
	global.GVA_MODEL
	Name        string  `json:"name" gorm:"size:64;uniqueIndex;comment:config name"`
	BaseURL     string  `json:"baseURL" gorm:"size:255;comment:ai base url"`
	Token       string  `json:"token" gorm:"size:2048;comment:ai token"`
	Model       string  `json:"model" gorm:"size:128;comment:model name"`
	ChatPath    string  `json:"chatPath" gorm:"size:255;comment:chat endpoint path"`
	Temperature float64 `json:"temperature" gorm:"comment:temperature"`
	TimeoutSec  int     `json:"timeoutSec" gorm:"comment:timeout seconds"`
	Enabled     bool    `json:"enabled" gorm:"comment:enabled"`
}

func (AIAssistantConfig) TableName() string {
	return "ai_assistant_configs"
}

type AIAssistantPrompt struct {
	global.GVA_MODEL
	Title   string `json:"title" gorm:"size:128;comment:title"`
	Content string `json:"content" gorm:"type:text;comment:prompt content"`
	Sort    int    `json:"sort" gorm:"comment:sort"`
	Enabled bool   `json:"enabled" gorm:"comment:enabled"`
}

func (AIAssistantPrompt) TableName() string {
	return "ai_assistant_prompts"
}
