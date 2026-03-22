//go:build driver_custom && !driver_mysql

package system

import "github.com/icosmos-space/iadmin/server/model/system/response"

var AutoCodeMysql = new(autoCodeMysqlStub)

type autoCodeMysqlStub struct{}

func (s *autoCodeMysqlStub) GetDB(businessDB string) (data []response.Db, err error) {
	panic(driverOmittedMsg("driver_mysql"))
}

func (s *autoCodeMysqlStub) GetTables(businessDB string, dbName string) (data []response.Table, err error) {
	panic(driverOmittedMsg("driver_mysql"))
}

func (s *autoCodeMysqlStub) GetColumn(businessDB string, tableName string, dbName string) (data []response.Column, err error) {
	panic(driverOmittedMsg("driver_mysql"))
}
