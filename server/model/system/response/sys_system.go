package response

import "github.com/icosmos-space/iadmin/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
