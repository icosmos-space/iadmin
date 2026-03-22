//go:build driver_custom && driver_mssql

package initdbcaps

func init() {
	supportedInitDBTypes = append(supportedInitDBTypes, "mssql")
}
