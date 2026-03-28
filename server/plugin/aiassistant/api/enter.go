package api

import "github.com/icosmos-space/iadmin/server/plugin/aiassistant/service"

var (
	Api              = new(api)
	serviceAssistant = service.Service.Assistant
)

type api struct{ Assistant assistant }
