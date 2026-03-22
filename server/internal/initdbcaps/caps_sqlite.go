//go:build driver_custom && driver_sqlite

package initdbcaps

func init() {
	supportedInitDBTypes = append(supportedInitDBTypes, "sqlite")
}
