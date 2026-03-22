package ast

import (
	"path/filepath"

	"github.com/icosmos-space/iadmin/server/global"
)

func init() {
	global.IADMIN_CONFIG.AutoCode.Root, _ = filepath.Abs("../../../")
	global.IADMIN_CONFIG.AutoCode.Server = "server"
}
