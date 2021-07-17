package curr

import (
	"context"
	"github.com/go-kratos/kratos/v2/middleware"
	"sync/atomic"
)

var inflight *int32 = new(int32)

// 正在进行中的请求 middleware
func Inflight() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			defer func() {
				atomic.AddInt32(inflight, -1)
			}()
			atomic.AddInt32(inflight, 1)
			return handler(ctx, req)
		}
	}
}

// 获取正在进行中的请求
func CurrInflight() int32 {
	return atomic.LoadInt32(inflight)
}
