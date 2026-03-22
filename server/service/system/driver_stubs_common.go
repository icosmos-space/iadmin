//go:build driver_custom

package system

import "fmt"

func driverOmittedMsg(driverTag string) string {
	return fmt.Sprintf("iadmin: 未编译 %s 驱动，请使用: go build -tags \"driver_custom,%s\"", driverTag, driverTag)
}
