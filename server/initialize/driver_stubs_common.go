//go:build driver_custom

package initialize

import "fmt"

// driverOmittedMsg 供桩函数 panic 使用（含 driver_custom 与各 driver_* 标签说明）。
func driverOmittedMsg(driverTag string) string {
	return fmt.Sprintf("iadmin: 未编译 %s 驱动，请使用: go build -tags \"driver_custom,%s\"", driverTag, driverTag)
}
