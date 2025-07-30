package node

import (
	"github.com/php-any/origami/data"
)

// InterfaceStatement 表示接口定义语句
type InterfaceStatement struct {
	*Node
	Name    string        // 接口名
	Extends *string       // 父接口名
	Methods []data.Method // 方法列表
}

// NewInterfaceStatement 创建一个新的接口定义语句
func NewInterfaceStatement(from data.From, name string, extends *string, methods []data.Method) *InterfaceStatement {
	return &InterfaceStatement{
		Node:    NewNode(from),
		Name:    name,
		Extends: extends,
		Methods: methods,
	}
}

// GetValue 获取接口定义语句的值
func (i *InterfaceStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {

	return i, nil
}

// GetName 返回接口名
func (i *InterfaceStatement) GetName() string {
	return i.Name
}

// GetExtend 返回父接口名
func (i *InterfaceStatement) GetExtend() *string {
	return i.Extends
}

// GetMethod 返回指定名称的方法
func (i *InterfaceStatement) GetMethod(name string) (data.Method, bool) {
	for _, method := range i.Methods {
		if method.GetName() == name {
			return method, true
		}
	}
	return nil, false
}

// GetMethods 返回所有方法
func (i *InterfaceStatement) GetMethods() []data.Method {
	return i.Methods
}

// InterfaceMethod 表示接口方法
type InterfaceMethod struct {
	*Node
	Name       string          // 方法名
	Modifier   data.Modifier   // 访问修饰符
	Params     []data.GetValue // 参数列表
	ReturnType data.Types      // 返回类型
}

// NewInterfaceMethod 创建一个新的接口方法
func NewInterfaceMethod(from data.From, name string, modifier string, params []data.GetValue, returnType data.Types) data.Method {
	return &InterfaceMethod{
		Node:       NewNode(from),
		Name:       name,
		Modifier:   data.NewModifier(modifier),
		Params:     params,
		ReturnType: returnType,
	}
}

// GetName 返回方法名
func (m *InterfaceMethod) GetName() string {
	return m.Name
}

// GetModifier 返回访问修饰符
func (m *InterfaceMethod) GetModifier() data.Modifier {
	return m.Modifier
}

// GetIsStatic 返回是否是静态方法
func (m *InterfaceMethod) GetIsStatic() bool {
	return false
}

// GetParams 返回参数列表
func (m *InterfaceMethod) GetParams() []data.GetValue {
	return m.Params
}

// GetVariables 返回变量列表
func (m *InterfaceMethod) GetVariables() []data.Variable {
	// 接口方法没有变量
	return nil
}

// Call 调用接口方法（接口方法不能直接调用）
func (m *InterfaceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 接口方法不能直接调用，应该抛出错误
	return nil, data.NewErrorThrow(m.from, data.NewError(m.from, "接口方法不能直接调用", nil))
}
