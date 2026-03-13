package runtime

import (
	"context"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
)

// Context 表示运行时上下文
type Context struct {
	vm data.VM
	// 命名空间
	namespace string

	// 变量存储符号表
	variables []*data.ZVal

	// 记录本次函数/方法调用时的实参表达式列表（用于 func_get_args 等）
	callArgs []data.GetValue
}

// NewContext 创建一个新的运行时上下文
func NewContext(vm data.VM) data.Context {
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
	if variable.GetIndex() >= len(c.variables) {
		return nil, data.NewErrorThrow(nil, errors.New("Variable does not exist"))
	}
	return c.variables[variable.GetIndex()].Value, nil
}

func (c *Context) GetIndexValue(index int) (data.Value, bool) {
	if index < 0 || index >= len(c.variables) {
		return nil, false
	}
	return c.variables[index].Value, true
}

func (c *Context) SetIndexZVal(index int, v *data.ZVal) {
	c.variables[index] = v
}

func (c *Context) GetIndexZVal(index int) *data.ZVal {
	return c.variables[index]
}

// SetVariableValue 设置变量值
func (c *Context) SetVariableValue(variable data.Variable, value data.Value) data.Control {
	switch v := value.(type) {
	case *data.ReferenceValue:
		c.variables[variable.GetIndex()] = v.Ctx.GetIndexZVal(v.Val.GetIndex())
	case *data.ArrayValue:
		c.variables[variable.GetIndex()].Value = data.CloneArrayValue(v)
	case *data.ObjectValue:
		// PHP 中 array 是按值赋值 + copy-on-write。
		// 在 Origami 里，关联数组可能由 ObjectValue 表示，这里也做一次结构级克隆，
		// 避免 `$b = $this->a; $b['k']=...` 反向修改到 `$this->a`（Symfony InputDefinition::$arguments 等场景）。
		c.variables[variable.GetIndex()].Value = data.CloneObjectValue(v)
	default:
		if len(c.variables) <= variable.GetIndex() {
			return data.NewErrorThrow(variable.(node.GetFrom).GetFrom(), errors.New("index out of range"))
		}
		c.variables[variable.GetIndex()].Value = value
	}

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

// SetVM 替换当前 Context 所绑定的 VM
func (c *Context) SetVM(vm data.VM) {
	c.vm = vm
}

// SetCallArgs 记录本次调用时传入的参数表达式列表
func (c *Context) SetCallArgs(args []data.GetValue) {
	c.callArgs = args
}

// GetCallArgs 获取本次调用时传入的参数表达式列表
func (c *Context) GetCallArgs() []data.GetValue {
	return c.callArgs
}

func makeSliceVariable(i int) []*data.ZVal {
	l := make([]*data.ZVal, i)
	for i := range l {
		l[i] = data.NewZVal(data.NewNullValue())
	}
	return l
}

// NewContextToDo 不实现具体功能的上下文
func NewContextToDo() data.Context {
	vm := NewVM(&parser.Parser{})
	return vm.CreateContext([]data.Variable{})
}
