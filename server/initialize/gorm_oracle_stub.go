//go:build driver_custom && !driver_oracle

package initialize

import (
	"github.com/icosmos-space/iadmin/server/config"
	"gorm.io/gorm"
)

func GormOracle() *gorm.DB {
	panic(driverOmittedMsg("driver_oracle"))
}

func GormOracleByConfig(m config.Oracle) *gorm.DB {
	panic(driverOmittedMsg("driver_oracle"))
}
