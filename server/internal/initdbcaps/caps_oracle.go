//go:build driver_custom && driver_oracle

package initdbcaps

func init() {
	supportedInitDBTypes = append(supportedInitDBTypes, "oracle")
}
