package telescope

import "github.com/php-any/origami/data"

// Load 不注册 Go 类：Telescope 为完整 Laravel 包（迁移、UI、采集），继续由 vendor PHP 提供。
func Load(_ data.VM) {}
