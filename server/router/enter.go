package router

import (
	"github.com/icosmos-space/iadmin/server/router/example"
	"github.com/icosmos-space/iadmin/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
}
