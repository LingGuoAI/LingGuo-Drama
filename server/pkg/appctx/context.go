// pkg/appctx/context.go
package appctx

import (
	"context"
	"sync"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
	once   sync.Once
)

// Initialize 初始化全局 context
func Initialize() (context.Context, context.CancelFunc) {
	once.Do(func() {
		ctx, cancel = context.WithCancel(context.Background())
	})
	return ctx, cancel
}

// GetContext 获取全局 context
func GetContext() context.Context {
	return ctx
}

// Cancel 取消全局 context
func Cancel() {
	if cancel != nil {
		cancel()
	}
}
