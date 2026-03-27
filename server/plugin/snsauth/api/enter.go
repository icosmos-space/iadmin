package api

import "github.com/icosmos-space/iadmin/server/plugin/snsauth/service"

var (
	Api         = new(api)
	serviceAuth = service.Service.Auth
)

type api struct{ Auth auth }
