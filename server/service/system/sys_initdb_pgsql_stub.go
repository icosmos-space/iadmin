//go:build driver_custom && !driver_pgsql

package system

import (
	"context"

	"github.com/icosmos-space/iadmin/server/model/system/request"
)

type PgsqlInitHandler struct{}

func NewPgsqlInitHandler() *PgsqlInitHandler {
	return &PgsqlInitHandler{}
}

func (h PgsqlInitHandler) WriteConfig(ctx context.Context) error {
	panic(driverOmittedMsg("driver_pgsql"))
}

func (h PgsqlInitHandler) EnsureDB(ctx context.Context, conf *request.InitDB) (context.Context, error) {
	panic(driverOmittedMsg("driver_pgsql"))
}

func (h PgsqlInitHandler) InitTables(ctx context.Context, inits initSlice) error {
	panic(driverOmittedMsg("driver_pgsql"))
}

func (h PgsqlInitHandler) InitData(ctx context.Context, inits initSlice) error {
	panic(driverOmittedMsg("driver_pgsql"))
}
