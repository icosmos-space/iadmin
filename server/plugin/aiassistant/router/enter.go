package router

import "github.com/icosmos-space/iadmin/server/plugin/aiassistant/api"

var (
	Router       = new(router)
	apiAssistant = api.Api.Assistant
)

type router struct{ Assistant assistant }
