//go:build driver_custom && !driver_sqlite

package initialize

import (
	"github.com/icosmos-space/iadmin/server/config"
	"gorm.io/gorm"
)

func GormSqlite() *gorm.DB {
	panic(driverOmittedMsg("driver_sqlite"))
}

func GormSqliteByConfig(s config.Sqlite) *gorm.DB {
	panic(driverOmittedMsg("driver_sqlite"))
}
