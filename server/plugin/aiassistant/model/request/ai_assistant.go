package request

import commonReq "github.com/icosmos-space/iadmin/server/model/common/request"

type PromptSearch struct {
	commonReq.PageInfo
	Enabled *bool `json:"enabled" form:"enabled"`
}

type UpdateConfigReq struct {
	BaseURL     string  `json:"baseURL"`
	Token       string  `json:"token"`
	Model       string  `json:"model"`
	ChatPath    string  `json:"chatPath"`
	Temperature float64 `json:"temperature"`
	TimeoutSec  int     `json:"timeoutSec"`
	Enabled     bool    `json:"enabled"`
}

type ChatMessage struct {
	Role        string           `json:"role"`
	Content     string           `json:"content"`
	Attachments []ChatAttachment `json:"attachments"`
}

type ChatAttachment struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

type ChatReq struct {
	Message     string           `json:"message"`
	Question    string           `json:"question"`
	Prompt      string           `json:"prompt"`
	Model       string           `json:"model"`
	Stream      bool             `json:"stream"`
	Messages    []ChatMessage    `json:"messages"`
	Attachments []ChatAttachment `json:"attachments"`
}
