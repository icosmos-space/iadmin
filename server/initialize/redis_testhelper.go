//go:build testhelper

package initialize

// ResetRedisSingletonForTest 重置 Redis 单例，仅在 testhelper 标签下可用。
func ResetRedisSingletonForTest() {
	resetRedisSingleton()
}
