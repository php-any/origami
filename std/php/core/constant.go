package core

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ConstantFunction 实现 constant(string $name): mixed
// 读取全局常量或 Class::CONST / Enum::Case。
type ConstantFunction struct{}

func NewConstantFunction() data.FuncStmt {
	return &ConstantFunction{}
}

func (f *ConstantFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	nameVal, ok := ctx.GetIndexValue(0)
	if !ok || nameVal == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("constant(): Argument #1 ($name) must be of type string"))
	}
	name := nameVal.AsString()
	vm := ctx.GetVM()
	if vm == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("constant(): 无法获取 VM"))
	}

	if v, exists := vm.GetConstant(name); exists {
		return v, nil
	}
	if v, ok := lookupClassConstant(vm, name); ok {
		return v, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Undefined constant \"%s\"", name))
}

func (f *ConstantFunction) GetName() string { return "constant" }

func (f *ConstantFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.String{}),
	}
}

func (f *ConstantFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.String{}),
	}
}
