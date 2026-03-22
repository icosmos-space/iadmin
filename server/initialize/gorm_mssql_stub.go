//go:build driver_custom && !driver_mssql

package initialize

import (
	"github.com/icosmos-space/iadmin/server/config"
	"gorm.io/gorm"
)

func GormMssql() *gorm.DB {
	panic(driverOmittedMsg("driver_mssql"))
}

func GormMssqlByConfig(m config.Mssql) *gorm.DB {
	panic(driverOmittedMsg("driver_mssql"))
}
