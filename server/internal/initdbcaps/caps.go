package initdbcaps

import "sort"

var supportedInitDBTypes []string

// SupportedInitDBTypes 返回当前二进制实际包含的数据库初始化能力（与 driver_custom / driver_* 构建标签一致）。
func SupportedInitDBTypes() []string {
	out := make([]string, len(supportedInitDBTypes))
	copy(out, supportedInitDBTypes)
	sort.Strings(out)
	return out
}
