package router

import "github.com/icosmos-space/iadmin/server/plugin/supabase/api"

var (
	Router  = new(router)
	apiUser = api.Api.User
)

type router struct{ User user }
