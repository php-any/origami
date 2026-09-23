package tinker

import "github.com/php-any/origami/data"

// Load 不注册 Go 类：依赖 PsySH 交互式 REPL，继续由 vendor PHP 提供。
func Load(_ data.VM) {}
