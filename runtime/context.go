package runtime

import (
	"context"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
)

// Context 表示运行时上下文
type Context struct {
	vm data.VM
	// 命名空间
	namespace string

	// 变量存储符号表
	variables []data.Value

	// 其他字段...
}

// NewContext 创建一个新的运行时上下文
func NewContext(vm *VM) data.Context {
	return &Context{
		vm: vm,
	}
}

// SetNamespace 设置命名空间
func (c *Context) SetNamespace(name string) data.Context {
	c.namespace = name
	return c
}

// GetNamespace 获取命名空间
func (c *Context) GetNamespace() string {
	return c.namespace
}

// GetVariableValue 获取变量值
func (c *Context) GetVariableValue(variable data.Variable) (data.Value, data.Control) {
	// 实现获取变量值的逻辑
	return c.variables[variable.GetIndex()], nil
}

func (c *Context) GetIndexValue(index int) (data.Value, bool) {
	if index < 0 || index >= len(c.variables) {
		return nil, false
	}
	return c.variables[index], true
}

// SetVariableValue 设置变量值
func (c *Context) SetVariableValue(variable data.Variable, value data.Value) data.Control {
	//if len(c.variables) <= variable.GetIndex() {
	//	c.variables = append(c.variables, data.NewNullValue())
	//}
	c.variables[variable.GetIndex()] = value
	return nil
}

// CreateContext 创建函数上下文
func (c *Context) CreateContext(vars []data.Variable) data.Context {
	return &Context{
		vm:        c.vm,
		variables: makeSliceVariable(len(vars)),
	}
}

func (c *Context) CreateBaseContext() data.Context {
	return &Context{
		vm:        c.vm,
		variables: nil,
	}
}

func (c *Context) GetVM() data.VM {
	return c.vm
}

func (c *Context) GoContext() context.Context {
	return context.Background()
}

func makeSliceVariable(i int) []data.Value {
	l := make([]data.Value, i)
	for i := range l {
		l[i] = data.NewNullValue()
	}
	return l
}

// NewContextToDo 不实现具体功能的上下文
func NewContextToDo() data.Context {
	vm := NewVM(&parser.Parser{})
	return vm.CreateContext([]data.Variable{})
}
