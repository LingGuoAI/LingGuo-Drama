package cache

import "time"

type Store interface {
	Set(key string, value string, expireTime time.Duration)
	Get(key string) string
	Has(key string) bool
	Forget(key string)
	Forever(key string, value string)
	Flush()

	IsAlive() error

	// Increment 当参数只有 1 个时，为 key，增加 1。
	// 当参数有 2 个时，第一个参数为 key ，第二个参数为要增加的值 int64 类型。
	Increment(parameters ...interface{})

	// Decrement 当参数只有 1 个时，为 key，减去 1。
	// 当参数有 2 个时，第一个参数为 key ，第二个参数为要减去的值 int64 类型。
	Decrement(parameters ...interface{})
}
