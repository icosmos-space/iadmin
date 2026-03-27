package router

import "github.com/icosmos-space/iadmin/server/plugin/snsauth/api"

var (
	Router  = new(router)
	apiAuth = api.Api.Auth
)

type router struct{ Auth auth }
