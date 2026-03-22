//go:build driver_custom && !driver_mysql

package initialize

import (
	"github.com/icosmos-space/iadmin/server/config"
	"gorm.io/gorm"
)

func GormMysql() *gorm.DB {
	panic(driverOmittedMsg("driver_mysql"))
}

func GormMysqlByConfig(m config.Mysql) *gorm.DB {
	panic(driverOmittedMsg("driver_mysql"))
}
