//go:build driver_custom && driver_pgsql

package initdbcaps

func init() {
	supportedInitDBTypes = append(supportedInitDBTypes, "pgsql")
}
