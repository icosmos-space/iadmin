package initialize

import (
	"context"

	model "github.com/icosmos-space/iadmin/server/model/system"
	"github.com/icosmos-space/iadmin/server/plugin/plugin-tool/utils"
)

func Menu(ctx context.Context) {
	// 第一个菜单会作为父菜单，其他菜单会自动使用第一个菜单的 ID 作为 ParentId
	entities := []model.SysBaseMenu{
		// 父菜单 - AI 助手（路由容器）
		{
			MenuLevel: 0,
			ParentId:  0,
			Path:      "aiAssistant",
			Name:      "aiAssistant",
			Hidden:    false,
			Component: "plugin/aiassistant/view/routerHolder.vue",
			Sort:      0,
			Meta:      model.Meta{Title: "AI 助手", Icon: "chat-line-round"},
		},
		// 子菜单 - AI 聊天
		{
			MenuLevel: 1,
			ParentId:  0, // 会被自动更新为父菜单 ID
			Path:      "aiAssistantChat",
			Name:      "aiAssistantChat",
			Hidden:    false,
			Component: "plugin/aiassistant/view/chat.vue",
			Sort:      0,
			Meta:      model.Meta{Title: "AI 聊天", Icon: "chat-dot-round"},
		},
		// 子菜单 - AI 助手配置
		{
			MenuLevel: 1,
			ParentId:  0, // 会被自动更新为父菜单 ID
			Path:      "aiAssistantConfig",
			Name:      "aiAssistantConfig",
			Hidden:    false,
			Component: "plugin/aiassistant/view/config.vue",
			Sort:      1,
			Meta:      model.Meta{Title: "AI 助手配置", Icon: "cpu"},
		},
		// 子菜单 - AI 预置提示词
		{
			MenuLevel: 1,
			ParentId:  0, // 会被自动更新为父菜单 ID
			Path:      "aiAssistantPrompt",
			Name:      "aiAssistantPrompt",
			Hidden:    false,
			Component: "plugin/aiassistant/view/prompt.vue",
			Sort:      2,
			Meta:      model.Meta{Title: "AI 预置提示词", Icon: "chat-dot-round"},
		},
	}
	utils.RegisterMenus(entities...)
}
