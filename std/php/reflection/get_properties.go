package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassGetPropertiesMethod 实现 ReflectionClass::getProperties
// 返回被反射的类的所有属性列表
type ReflectionClassGetPropertiesMethod struct{}

// GetName 返回方法名 "getProperties"
func (m *ReflectionClassGetPropertiesMethod) GetName() string { return "getProperties" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassGetPropertiesMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassGetPropertiesMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - filter: 可选的过滤器标志（?int），用于过滤属性的可见性，默认为 null
func (m *ReflectionClassGetPropertiesMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filter", 0, node.NewNullLiteral(nil), data.NewBaseType("?int")),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionClassGetPropertiesMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filter", 0, data.NewBaseType("?int")),
	}
}

// GetReturnType 返回返回类型，返回数组类型
func (m *ReflectionClassGetPropertiesMethod) GetReturnType() data.Types {
	return data.Arrays{}
}

// Call 执行 getProperties 方法
// 返回被反射的类的所有属性列表（ReflectionProperty 对象数组）
// TODO: 实现完整的过滤器逻辑，当前忽略 filter 参数
func (m *ReflectionClassGetPropertiesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	className, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// PHP：无 filter 时包含继承链上的声明属性（Livewire 靠此拿到父类 public Collection $notifications）。
	seen := make(map[string]bool)
	result := make([]data.Value, 0)
	currentName := className
	current := classStmt
	vm := ctx.GetVM()
	for current != nil {
		for _, prop := range current.GetPropertyList() {
			if prop == nil {
				continue
			}
			name := prop.GetName()
			if seen[name] {
				continue
			}
			seen[name] = true
			result = append(result, newReflectionProperty(ctx, currentName, name))
		}
		extend := current.GetExtend()
		if extend == nil || *extend == "" || vm == nil {
			break
		}
		parent, acl := vm.GetOrLoadClass(*extend)
		if acl != nil || parent == nil {
			break
		}
		currentName = parent.GetName()
		current = parent
	}

	return data.NewArrayValue(result), nil
}
