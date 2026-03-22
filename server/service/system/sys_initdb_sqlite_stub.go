//go:build driver_custom && !driver_sqlite

package system

import (
	"context"

	"github.com/icosmos-space/iadmin/server/model/system/request"
)

type SqliteInitHandler struct{}

func NewSqliteInitHandler() *SqliteInitHandler {
	return &SqliteInitHandler{}
}

func (h SqliteInitHandler) WriteConfig(ctx context.Context) error {
	panic(driverOmittedMsg("driver_sqlite"))
}

func (h SqliteInitHandler) EnsureDB(ctx context.Context, conf *request.InitDB) (context.Context, error) {
	panic(driverOmittedMsg("driver_sqlite"))
}

func (h SqliteInitHandler) InitTables(ctx context.Context, inits initSlice) error {
	panic(driverOmittedMsg("driver_sqlite"))
}

func (h SqliteInitHandler) InitData(ctx context.Context, inits initSlice) error {
	panic(driverOmittedMsg("driver_sqlite"))
}
