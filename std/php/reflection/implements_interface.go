package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassImplementsInterfaceMethod 实现 ReflectionClass::implementsInterface
// 检查被反射的类是否实现了指定接口
type ReflectionClassImplementsInterfaceMethod struct{}

func (m *ReflectionClassImplementsInterfaceMethod) GetName() string {
	return "implementsInterface"
}

func (m *ReflectionClassImplementsInterfaceMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *ReflectionClassImplementsInterfaceMethod) GetIsStatic() bool { return false }

func (m *ReflectionClassImplementsInterfaceMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "interface", 0, nil, data.Mixed{}),
	}
}

func (m *ReflectionClassImplementsInterfaceMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "interface", 0, data.Mixed{}),
	}
}

func (m *ReflectionClassImplementsInterfaceMethod) GetReturnType() data.Types {
	return data.Bool{}
}

// Call 执行 implementsInterface 方法
// 检查被反射的类是否实现了指定接口
func (m *ReflectionClassImplementsInterfaceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	ifaceValue, _ := ctx.GetIndexValue(0)
	if ifaceValue == nil {
		return data.NewBoolValue(false), nil
	}

	var targetInterfaceName string
	switch v := ifaceValue.(type) {
	case *data.ClassValue:
		targetInterfaceName = v.Class.GetName()
	case *data.StringValue:
		targetInterfaceName = v.Value
	default:
		targetInterfaceName = ifaceValue.AsString()
	}

	_, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return data.NewBoolValue(false), nil
	}

	vm := ctx.GetVM()

	// 递归检查：遍历类自身及其继承链，检查是否实现了目标接口
	var checkClass func(stmt data.ClassStmt) bool
	checkClass = func(stmt data.ClassStmt) bool {
		// 检查当前类直接实现的接口
		if implements := stmt.GetImplements(); implements != nil {
			for _, impl := range implements {
				if impl == targetInterfaceName {
					return true
				}
				// 递归检查接口的父接口（接口可以继承接口）
				if ifaceStmt, ok := vm.GetClass(impl); ok {
					if checkClass(ifaceStmt) {
						return true
					}
				}
			}
		}
		// 检查父类
		if stmt.GetExtend() != nil {
			if parentStmt, ok := vm.GetClass(*stmt.GetExtend()); ok {
				return checkClass(parentStmt)
			}
		}
		return false
	}

	return data.NewBoolValue(checkClass(classStmt)), nil
}
