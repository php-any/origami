package sentinel

import "github.com/php-any/origami/data"

// Load 不注册 Go 类：SentinelManager / Facade / Drivers 依赖 Illuminate Manager 与容器，继续由 vendor PHP 提供。
func Load(_ data.VM) {}
