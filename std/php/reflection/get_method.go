package reflection

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionClassGetMethodMethod 实现 ReflectionClass::getMethod
// 根据方法名获取被反射类的指定方法（含父类公开/受保护方法，对齐 PHP）。
type ReflectionClassGetMethodMethod struct{}

// GetName 返回方法名 "getMethod"
func (m *ReflectionClassGetMethodMethod) GetName() string { return "getMethod" }

// GetModifier 返回方法修饰符，公开方法
func (m *ReflectionClassGetMethodMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，非静态方法
func (m *ReflectionClassGetMethodMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - name: 方法名（字符串）
func (m *ReflectionClassGetMethodMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.String{}),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionClassGetMethodMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.String{}),
	}
}

// GetReturnType 返回返回类型
func (m *ReflectionClassGetMethodMethod) GetReturnType() data.Types {
	return data.NewBaseType("ReflectionMethod")
}

// Call 执行 getMethod 方法。
// 方法不存在时抛出 ReflectionException（供 catch (ReflectionException) 捕获）。
func (m *ReflectionClassGetMethodMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	methodNameValue, _ := ctx.GetIndexValue(0)
	if methodNameValue == nil {
		return nil, createReflectionException("ReflectionClass::getMethod(): Argument #1 ($name) must be of type string, null given", ctx, nil)
	}

	methodName := methodNameValue.AsString()
	className, classStmt := getReflectionClassInfo(ctx)
	if classStmt == nil {
		return nil, createReflectionException(fmt.Sprintf("Method %s does not exist", methodName), ctx, nil)
	}

	declaredClass := className
	current := classStmt
	for current != nil {
		if _, exists := current.GetMethod(methodName); exists {
			declaredClass = current.GetName()
			methodClass := &ReflectionMethodClass{}
			methodValue := data.NewClassValue(methodClass, ctx.CreateBaseContext())
			methodValue.ObjectValue.SetProperty("_className", data.NewStringValue(declaredClass))
			methodValue.ObjectValue.SetProperty("_methodName", data.NewStringValue(methodName))
			// 保留原始被反射类名，便于声明类与反射目标区分（PHP getDeclaringClass 语义）
			methodValue.ObjectValue.SetProperty("_reflectedClassName", data.NewStringValue(className))
			return methodValue, nil
		}
		if staticMethods, ok := current.(data.GetStaticMethod); ok {
			if _, exists := staticMethods.GetStaticMethod(methodName); exists {
				declaredClass = current.GetName()
				methodClass := &ReflectionMethodClass{}
				methodValue := data.NewClassValue(methodClass, ctx.CreateBaseContext())
				methodValue.ObjectValue.SetProperty("_className", data.NewStringValue(declaredClass))
				methodValue.ObjectValue.SetProperty("_methodName", data.NewStringValue(methodName))
				methodValue.ObjectValue.SetProperty("_reflectedClassName", data.NewStringValue(className))
				return methodValue, nil
			}
		}
		extend := current.GetExtend()
		if extend == nil || *extend == "" {
			break
		}
		parent, acl := ctx.GetVM().GetOrLoadClass(*extend)
		if acl != nil {
			return nil, acl
		}
		if parent == nil {
			break
		}
		current = parent
	}

	return nil, createReflectionException(fmt.Sprintf("Method %s does not exist", methodName), ctx, nil)
}
