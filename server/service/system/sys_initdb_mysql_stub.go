//go:build driver_custom && !driver_mysql

package system

import (
	"context"

	"github.com/icosmos-space/iadmin/server/model/system/request"
)

type MysqlInitHandler struct{}

func NewMysqlInitHandler() *MysqlInitHandler {
	return &MysqlInitHandler{}
}

func (h MysqlInitHandler) WriteConfig(ctx context.Context) error {
	panic(driverOmittedMsg("driver_mysql"))
}

func (h MysqlInitHandler) EnsureDB(ctx context.Context, conf *request.InitDB) (context.Context, error) {
	panic(driverOmittedMsg("driver_mysql"))
}

func (h MysqlInitHandler) InitTables(ctx context.Context, inits initSlice) error {
	panic(driverOmittedMsg("driver_mysql"))
}

func (h MysqlInitHandler) InitData(ctx context.Context, inits initSlice) error {
	panic(driverOmittedMsg("driver_mysql"))
}
