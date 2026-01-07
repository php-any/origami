package reflection

import (
	"github.com/php-any/origami/data"
)

// getReflectionClassInfo 从上下文中获取 ReflectionClass 的类信息
// 该函数从 ReflectionClass 实例的 _className 属性中获取被反射的类名，
// 然后从 VM 中加载对应的类语句
//
// 参数:
//   - ctx: 运行时上下文
//
// 返回:
//   - string: 类名，如果获取失败则返回空字符串
//   - data.ClassStmt: 类语句对象，如果获取失败则返回 nil
func getReflectionClassInfo(ctx data.Context) (string, data.ClassStmt) {
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// 从 ObjectValue 的 property 中获取类名（实例属性）
		// 使用 ObjectValue.GetProperty 来获取实例属性，而不是 ClassValue.GetProperty
		classNameVal, _ := objCtx.ObjectValue.GetProperty("_className")
		if strVal, ok := classNameVal.(*data.StringValue); ok {
			className := strVal.AsString()
			vm := ctx.GetVM()
			stmt, _ := vm.GetOrLoadClass(className)
			return className, stmt
		}
	}
	return "", nil
}
