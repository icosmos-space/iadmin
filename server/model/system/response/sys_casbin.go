package response

import (
	"github.com/icosmos-space/iadmin/server/model/system/request"
)

type PolicyPathResponse struct {
	Paths []request.CasbinInfo `json:"paths"`
}
