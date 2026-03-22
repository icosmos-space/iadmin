//go:build driver_custom && !driver_mssql

package system

import (
	"context"

	"github.com/icosmos-space/iadmin/server/model/system/request"
)

type MssqlInitHandler struct{}

func NewMssqlInitHandler() *MssqlInitHandler {
	return &MssqlInitHandler{}
}

func (h MssqlInitHandler) WriteConfig(ctx context.Context) error {
	panic(driverOmittedMsg("driver_mssql"))
}

func (h MssqlInitHandler) EnsureDB(ctx context.Context, conf *request.InitDB) (context.Context, error) {
	panic(driverOmittedMsg("driver_mssql"))
}

func (h MssqlInitHandler) InitTables(ctx context.Context, inits initSlice) error {
	panic(driverOmittedMsg("driver_mssql"))
}

func (h MssqlInitHandler) InitData(ctx context.Context, inits initSlice) error {
	panic(driverOmittedMsg("driver_mssql"))
}
