package core

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// DefinedFunction 实现 defined 函数
// defined(string $constant_name): bool
// 支持全局常量与 Class::CONST / Enum::Case。
type DefinedFunction struct{}

func NewDefinedFunction() data.FuncStmt {
	return &DefinedFunction{}
}

func (f *DefinedFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	constNameValue, ok := ctx.GetIndexValue(0)
	if !ok || constNameValue == nil {
		return data.NewBoolValue(false), nil
	}

	constName := constNameValue.AsString()
	vm := ctx.GetVM()
	if vm == nil {
		return data.NewBoolValue(false), nil
	}

	if _, exists := vm.GetConstant(constName); exists {
		return data.NewBoolValue(true), nil
	}

	if _, ok := lookupClassConstant(vm, constName); ok {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(false), nil
}

func (f *DefinedFunction) GetName() string {
	return "defined"
}

func (f *DefinedFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "constant_name", 0, nil, data.String{}),
	}
}

func (f *DefinedFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "constant_name", 0, data.String{}),
	}
}

// lookupClassConstant 解析 "Class::NAME" / "Namespace\Class::NAME"
func lookupClassConstant(vm data.VM, name string) (data.Value, bool) {
	idx := strings.LastIndex(name, "::")
	if idx <= 0 || idx+2 >= len(name) {
		return nil, false
	}
	className := strings.TrimPrefix(name[:idx], "\\")
	constName := name[idx+2:]
	if className == "" || constName == "" {
		return nil, false
	}

	cls, ok := vm.GetClass(className)
	if !ok {
		cls, _ = vm.GetOrLoadClass(className)
	}
	if cls == nil {
		return nil, false
	}

	if gsp, ok := cls.(data.GetStaticProperty); ok {
		if v, found := gsp.GetStaticProperty(constName); found {
			return v, true
		}
	}
	return nil, false
}
