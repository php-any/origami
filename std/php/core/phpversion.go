package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"strings"
)

type PhpVersionFunction struct{}

func NewPhpVersionFunction() data.FuncStmt {
	return &PhpVersionFunction{}
}

func (f *PhpVersionFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	if extension, ok := ctx.GetIndexValue(0); ok && extension != nil {
		if _, null := extension.(*data.NullValue); !null && !extensionAvailable(vm, strings.ToLower(extension.AsString())) {
			return data.NewBoolValue(false), nil
		}
	}
	if v, ok := vm.GetConstant("PHP_VERSION"); ok {
		return v, nil
	}
	return data.NewStringValue("8.2.0"), nil
}

func (f *PhpVersionFunction) GetName() string { return "phpversion" }

var phpVersionFunctionGetParams = []data.GetValue{node.NewParameter(nil, "extension", 0, data.NewNullValue(), nil)}

func (f *PhpVersionFunction) GetParams() []data.GetValue {
	return phpVersionFunctionGetParams
}

var phpVersionFunctionGetVariables = []data.Variable{node.NewVariable(nil, "extension", 0, nil)}

func (f *PhpVersionFunction) GetVariables() []data.Variable {
	return phpVersionFunctionGetVariables
}
