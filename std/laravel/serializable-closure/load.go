package serializableclosure

import "github.com/php-any/origami/data"

// Load 不注册 Go 类：SerializableClosure 依赖 Native/Signed 序列化与 ReflectionClosure，
// 无法在运行时热路径外安全复刻；继续由 vendor PHP 提供。
func Load(_ data.VM) {}
