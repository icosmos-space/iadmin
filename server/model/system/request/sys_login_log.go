package request

import (
	"github.com/icosmos-space/iadmin/server/model/common/request"
	"github.com/icosmos-space/iadmin/server/model/system"
)

type SysLoginLogSearch struct {
	system.SysLoginLog
	request.PageInfo
}
