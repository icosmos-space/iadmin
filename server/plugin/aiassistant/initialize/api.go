package initialize

import (
	"context"

	model "github.com/icosmos-space/iadmin/server/model/system"
	"github.com/icosmos-space/iadmin/server/plugin/plugin-tool/utils"
)

func Api(ctx context.Context) {
	entities := []model.SysApi{
		{Path: "/aiAssistant/chat", Description: "AI assistant chat", ApiGroup: "AI Assistant", Method: "POST"},
		{Path: "/aiAssistant/getEnabledPrompts", Description: "Get enabled preset prompts", ApiGroup: "AI Assistant", Method: "GET"},
		{Path: "/aiAssistant/getConfig", Description: "Get AI assistant config", ApiGroup: "AI Assistant", Method: "GET"},
		{Path: "/aiAssistant/updateConfig", Description: "Update AI assistant config", ApiGroup: "AI Assistant", Method: "PUT"},
		{Path: "/aiAssistant/createPrompt", Description: "Create preset prompt", ApiGroup: "AI Assistant", Method: "POST"},
		{Path: "/aiAssistant/updatePrompt", Description: "Update preset prompt", ApiGroup: "AI Assistant", Method: "PUT"},
		{Path: "/aiAssistant/deletePrompt", Description: "Delete preset prompt", ApiGroup: "AI Assistant", Method: "DELETE"},
		{Path: "/aiAssistant/deletePromptByIds", Description: "Batch delete preset prompts", ApiGroup: "AI Assistant", Method: "DELETE"},
		{Path: "/aiAssistant/findPrompt", Description: "Find preset prompt by ID", ApiGroup: "AI Assistant", Method: "GET"},
		{Path: "/aiAssistant/getPromptList", Description: "Get preset prompt list", ApiGroup: "AI Assistant", Method: "GET"},
	}
	utils.RegisterApis(entities...)
}
