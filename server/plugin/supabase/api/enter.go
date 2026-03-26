package api

import "github.com/icosmos-space/iadmin/server/plugin/supabase/service"

var (
	Api         = new(api)
	serviceUser = service.Service.User
)

type api struct{ User user }
