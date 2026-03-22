//go:build driver_custom && !driver_mssql

package system

import "github.com/icosmos-space/iadmin/server/model/system/response"

var AutoCodeMssql = new(autoCodeMssqlStub)

type autoCodeMssqlStub struct{}

func (s *autoCodeMssqlStub) GetDB(businessDB string) (data []response.Db, err error) {
	panic(driverOmittedMsg("driver_mssql"))
}

func (s *autoCodeMssqlStub) GetTables(businessDB string, dbName string) (data []response.Table, err error) {
	panic(driverOmittedMsg("driver_mssql"))
}

func (s *autoCodeMssqlStub) GetColumn(businessDB string, tableName string, dbName string) (data []response.Column, err error) {
	panic(driverOmittedMsg("driver_mssql"))
}
