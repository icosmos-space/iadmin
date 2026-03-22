//go:build driver_custom && driver_mysql

package initdbcaps

func init() {
	supportedInitDBTypes = append(supportedInitDBTypes, "mysql")
}
