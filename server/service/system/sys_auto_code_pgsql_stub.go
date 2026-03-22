//go:build driver_custom && !driver_pgsql

package system

import "github.com/icosmos-space/iadmin/server/model/system/response"

var AutoCodePgsql = new(autoCodePgsqlStub)

type autoCodePgsqlStub struct{}

func (s *autoCodePgsqlStub) GetDB(businessDB string) (data []response.Db, err error) {
	panic(driverOmittedMsg("driver_pgsql"))
}

func (s *autoCodePgsqlStub) GetTables(businessDB string, dbName string) (data []response.Table, err error) {
	panic(driverOmittedMsg("driver_pgsql"))
}

func (s *autoCodePgsqlStub) GetColumn(businessDB string, tableName string, dbName string) (data []response.Column, err error) {
	panic(driverOmittedMsg("driver_pgsql"))
}
