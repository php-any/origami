package reflection

import (
	"github.com/php-any/origami/data"
)

// ReflectionClassGetConstructorMethod 实现 ReflectionClass::getConstructor
// 返回被反射类的构造函数
type ReflectionClassGetConstructorMethod struct{}

// GetName 返回方法名 "getConstructor"
func (m *ReflectionClassGetConstructorMethod) GetName() string { return "getConstructor" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassGetConstructorMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassGetConstructorMethod) GetIsStatic() bool { return false }

var reflectionClassGetConstructorMethodGetParams = []data.GetValue{}

// GetParams 返回参数列表，该方法无参数
func (m *ReflectionClassGetConstructorMethod) GetParams() []data.GetValue {
	return reflectionClassGetConstructorMethodGetParams
}

var reflectionClassGetConstructorMethodGetVariables = []data.Variable{}

// GetVariables 返回变量列表，该方法无变量
func (m *ReflectionClassGetConstructorMethod) GetVariables() []data.Variable {
	return reflectionClassGetConstructorMethodGetVariables
}

// GetReturnType 返回返回类型，返回混合类型（ReflectionMethod 对象或 null）
func (m *ReflectionClassGetConstructorMethod) GetReturnType() data.Types {
	return data.Mixed{}
}

// Call 执行 getConstructor 方法
// 返回被反射类的构造函数（含从父类继承的构造函数，对齐 PHP ReflectionClass::getConstructor）。
// 如果类与祖先都没有构造函数，返回 null
func (m *ReflectionClassGetConstructorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewNullValue(), nil
	}

	constructor, declaringName, acl := findInheritedConstructor(ctx, classStmt, className)
	if acl != nil {
		return nil, acl
	}
	if constructor == nil {
		return data.NewNullValue(), nil
	}

	methodValue := newReflectionMethod(ctx, declaringName, constructor.GetName())

	return methodValue, nil
}

// findInheritedConstructor 沿继承链查找 __construct（PHP：子类未声明则使用父类构造函数）。
func findInheritedConstructor(ctx data.Context, classStmt data.ClassStmt, className string) (data.Method, string, data.Control) {
	if ctor := classStmt.GetConstruct(); ctor != nil {
		return ctor, className, nil
	}
	vm := ctx.GetVM()
	if vm == nil {
		return nil, "", nil
	}
	last := classStmt
	for last.GetExtend() != nil {
		ext := last.GetExtend()
		next, acl := vm.GetOrLoadClass(*ext)
		if acl != nil {
			return nil, "", acl
		}
		if next == nil {
			break
		}
		if ctor := next.GetConstruct(); ctor != nil {
			return ctor, next.GetName(), nil
		}
		last = next
	}
	return nil, "", nil
}
