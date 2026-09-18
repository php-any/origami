package reflection

import (
	"github.com/php-any/origami/data"
)

func setReflectionClassIdentity(classValue *data.ClassValue, className string) {
	if classValue == nil || classValue.ObjectValue == nil {
		return
	}
	s := data.NewStringValue(className)
	classValue.ObjectValue.SetProperty("_className", s)
	// PHP ReflectionClass 公开属性 $name
	classValue.ObjectValue.SetProperty("name", s)
}

func newReflectionClassValue(ctx data.Context, className string) *data.ClassValue {
	classValue := data.NewClassValue(&ReflectionClassClass{}, ctx.CreateBaseContext())
	setReflectionClassIdentity(classValue, className)
	return classValue
}

// getReflectionClassInfo 从上下文中获取 ReflectionClass 的类信息
// 该函数从 ReflectionClass 实例的 _className / name 属性中获取被反射的类名，
// 然后从 VM 中加载对应的类语句
func getReflectionClassInfo(ctx data.Context) (string, data.ClassStmt) {
	objCtx, ok := ctx.(*data.ClassMethodContext)
	if !ok || objCtx.ObjectValue == nil {
		return "", nil
	}
	props := objCtx.ObjectValue.GetProperties()
	className := phpValueAsString(props["_className"])
	if className == "" {
		className = phpValueAsString(props["name"])
	}
	if className == "" {
		return "", nil
	}
	vm := ctx.GetVM()
	v, acl := vm.LoadPkg(className)
	if acl != nil {
		return "", nil
	}
	if v != nil {
		if stmt, ok := v.(data.ClassStmt); ok {
			return className, stmt
		}
	}
	return className, nil
}
