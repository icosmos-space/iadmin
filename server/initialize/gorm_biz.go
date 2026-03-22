package initialize

import (
	"github.com/icosmos-space/iadmin/server/global"
)

func bizModel() error {
	db := global.IADMIN_DB
	err := db.AutoMigrate()
	if err != nil {
		return err
	}
	return nil
}
