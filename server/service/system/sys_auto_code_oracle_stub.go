//go:build driver_custom && !driver_oracle

package system

import "github.com/icosmos-space/iadmin/server/model/system/response"

var AutoCodeOracle = new(autoCodeOracleStub)

type autoCodeOracleStub struct{}

func (s *autoCodeOracleStub) GetDB(businessDB string) (data []response.Db, err error) {
	panic(driverOmittedMsg("driver_oracle"))
}

func (s *autoCodeOracleStub) GetTables(businessDB string, dbName string) (data []response.Table, err error) {
	panic(driverOmittedMsg("driver_oracle"))
}

func (s *autoCodeOracleStub) GetColumn(businessDB string, tableName string, dbName string) (data []response.Column, err error) {
	panic(driverOmittedMsg("driver_oracle"))
}
