package response

import "github.com/icosmos-space/iadmin/server/model/system"

type SnsAuthResult struct {
	Mode     string         `json:"mode"`
	Provider string         `json:"provider"`
	Token    string         `json:"token,omitempty"`
	User     system.SysUser `json:"user,omitempty"`
	Msg      string         `json:"msg"`
}
