package initialize

import (
	_ "github.com/icosmos-space/iadmin/server/source/example"
	_ "github.com/icosmos-space/iadmin/server/source/system"
)

func init() {
	// do nothing,only import source package so that inits can be registered
}
