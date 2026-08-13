package annotation

import (
	"sync"

	"github.com/php-any/origami/data"
)

// OnApplicationScanStart 在 #[Application] 扫描目录前调用；返回的 cleanup 在路由注册完成后执行。
// 由 std/container 等标准库注册，用于建立应用级 IoC 作用域。
var (
	applicationScanHookMu  sync.RWMutex
	onApplicationScanStart func(ctx data.Context) (cleanup func(), acl data.Control)
)

func SetApplicationScanStartHook(hook func(ctx data.Context) (cleanup func(), acl data.Control)) {
	applicationScanHookMu.Lock()
	defer applicationScanHookMu.Unlock()
	onApplicationScanStart = hook
}

func applicationScanStartHook() func(ctx data.Context) (cleanup func(), acl data.Control) {
	applicationScanHookMu.RLock()
	defer applicationScanHookMu.RUnlock()
	return onApplicationScanStart
}
