package request

import (
	"github.com/icosmos-space/iadmin/server/model/common/request"
	"github.com/icosmos-space/iadmin/server/model/system"
)

type SysApiTokenSearch struct {
	system.SysApiToken
	request.PageInfo
    Status *bool `json:"status" form:"status"`
}
