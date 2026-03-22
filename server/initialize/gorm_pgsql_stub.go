//go:build driver_custom && !driver_pgsql

package initialize

import (
	"github.com/icosmos-space/iadmin/server/config"
	"gorm.io/gorm"
)

func GormPgSql() *gorm.DB {
	panic(driverOmittedMsg("driver_pgsql"))
}

func GormPgSqlByConfig(p config.Pgsql) *gorm.DB {
	panic(driverOmittedMsg("driver_pgsql"))
}
