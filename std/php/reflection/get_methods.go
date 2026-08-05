package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassGetMethodsMethod 实现 ReflectionClass::getMethods
// 返回被反射的类的所有方法列表
type ReflectionClassGetMethodsMethod struct{}

// GetName 返回方法名 "getMethods"
func (m *ReflectionClassGetMethodsMethod) GetName() string { return "getMethods" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassGetMethodsMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassGetMethodsMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - filter: 可选的过滤器标志（?int），用于过滤方法的可见性，默认为 null
func (m *ReflectionClassGetMethodsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filter", 0, node.NewNullLiteral(nil), data.NewBaseType("?int")),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionClassGetMethodsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filter", 0, data.NewBaseType("?int")),
	}
}

// GetReturnType 返回返回类型，返回数组类型
func (m *ReflectionClassGetMethodsMethod) GetReturnType() data.Types {
	return data.Arrays{}
}

// Call 执行 getMethods 方法
// 返回被反射的类的所有方法列表（当前实现返回方法名数组）
// TODO: 实现完整的过滤器逻辑，当前忽略 filter 参数
func (m *ReflectionClassGetMethodsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	filter := 0
	if filterValue, ok := ctx.GetIndexValue(0); ok && filterValue != nil {
		if _, isNull := filterValue.(*data.NullValue); !isNull {
			if asInt, ok := filterValue.(data.AsInt); ok {
				filter, _ = asInt.AsInt()
			}
		}
	}

	result := make([]data.Value, 0)
	seen := make(map[string]bool)
	currentName := className
	current := classStmt
	inherited := false

	for current != nil {
		for _, method := range current.GetMethods() {
			if method == nil || seen[method.GetName()] {
				continue
			}
			if inherited && method.GetModifier() == data.ModifierPrivate {
				continue
			}
			seen[method.GetName()] = true
			if filter != 0 && reflectionMethodModifiers(method)&filter == 0 {
				continue
			}
			result = append(result, newReflectionMethod(ctx, currentName, method.GetName()))
		}

		extend := current.GetExtend()
		if extend == nil || *extend == "" {
			break
		}
		parent, acl := ctx.GetVM().GetOrLoadClass(*extend)
		if acl != nil || parent == nil {
			break
		}
		currentName = *extend
		current = parent
		inherited = true
	}

	return data.NewArrayValue(result), nil
}

func reflectionMethodModifiers(method data.Method) int {
	modifiers := 0
	switch method.GetModifier() {
	case data.ModifierPublic:
		modifiers |= 1
	case data.ModifierProtected:
		modifiers |= 2
	case data.ModifierPrivate:
		modifiers |= 4
	}
	if method.GetIsStatic() {
		modifiers |= 16
	}
	return modifiers
}
