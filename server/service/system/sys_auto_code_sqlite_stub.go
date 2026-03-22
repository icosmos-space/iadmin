//go:build driver_custom && !driver_sqlite

package system

import "github.com/icosmos-space/iadmin/server/model/system/response"

var AutoCodeSqlite = new(autoCodeSqliteStub)

type autoCodeSqliteStub struct{}

func (s *autoCodeSqliteStub) GetDB(businessDB string) (data []response.Db, err error) {
	panic(driverOmittedMsg("driver_sqlite"))
}

func (s *autoCodeSqliteStub) GetTables(businessDB string, dbName string) (data []response.Table, err error) {
	panic(driverOmittedMsg("driver_sqlite"))
}

func (s *autoCodeSqliteStub) GetColumn(businessDB string, tableName string, dbName string) (data []response.Column, err error) {
	panic(driverOmittedMsg("driver_sqlite"))
}
