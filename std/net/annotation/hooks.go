package annotation

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// OnApplicationScanStart 在 #[Application] 扫描目录前调用；返回的 cleanup 在路由注册完成后执行。
// 由 std/container 等标准库注册，用于建立应用级 IoC 作用域。
var OnApplicationScanStart func(ctx data.Context) (cleanup func(), acl data.Control)

// ControllerInstantiator 在路由注册阶段实例化控制器；默认直接 new，不经过容器。
var ControllerInstantiator = node.InstantiateController
