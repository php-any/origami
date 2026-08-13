package reflection

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionAttributeConstructMethod 实现 ReflectionAttribute::__construct
// 构造函数用于初始化 ReflectionAttribute 实例，接收注解对象作为参数
type ReflectionAttributeConstructMethod struct{}

// GetName 返回方法名 "__construct"
func (m *ReflectionAttributeConstructMethod) GetName() string { return "__construct" }

// GetModifier 返回方法修饰符，构造函数是公开的
func (m *ReflectionAttributeConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }

// GetIsStatic 返回是否为静态方法，构造函数不是静态方法
func (m *ReflectionAttributeConstructMethod) GetIsStatic() bool { return false }

// GetParams 返回参数列表
// 参数:
//   - name: 属性名称（字符串），类型为 String
func (m *ReflectionAttributeConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.String{}),
	}
}

// GetVariables 返回变量列表
func (m *ReflectionAttributeConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.String{}),
	}
}

// GetReturnType 返回返回类型，构造函数无返回值
func (m *ReflectionAttributeConstructMethod) GetReturnType() data.Types { return nil }

// Call 执行构造函数
// 注意：PHP 的 ReflectionAttribute 构造函数实际上不接受参数
// 它是在内部创建的，但为了兼容性，我们接受一个注解对象作为参数
func (m *ReflectionAttributeConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：注解对象（ClassValue）
	annotationValue, _ := ctx.GetIndexValue(0)
	if annotationValue == nil {
		return nil, data.NewErrorThrow(nil, errors.New("ReflectionAttribute::__construct() expects parameter 1 to be object"))
	}

	// 验证参数是 ClassValue（注解对象）
	if _, ok := annotationValue.(*data.ClassValue); !ok {
		return nil, data.NewErrorThrow(nil, errors.New("ReflectionAttribute::__construct() expects parameter 1 to be object"))
	}

	// 将注解对象存储到当前对象的属性中
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// 存储注解对象到 ObjectValue 的实例属性中
		objCtx.ObjectValue.SetProperty("_annotation", annotationValue)
	}

	return nil, nil
}
