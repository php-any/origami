package macroable

import "github.com/php-any/origami/data"

// Load 注册 illuminate/macroable。
// Macroable 为 trait；具体类的 macro/__call 由 kit.RegisterMacroable 在各类 Load 中挂载，此处不重复注册。
func Load(_ data.VM) {}
