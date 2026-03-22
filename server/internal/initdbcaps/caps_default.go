//go:build !driver_custom

package initdbcaps

func init() {
	supportedInitDBTypes = []string{"mysql", "pgsql", "oracle", "mssql", "sqlite"}
}
