package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RunPhpFileFunction 执行通过 RegisterCompiledFile 注册的预编译 PHP 文件
type RunPhpFileFunction struct{}

func NewRunPhpFileFunction() data.FuncStmt {
	return &RunPhpFileFunction{}
}

func (f *RunPhpFileFunction) GetName() string {
	return "run_php_file"
}

func (f *RunPhpFileFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterRawAST(nil, "file", 0, data.Mixed{}),
	}
}

func (f *RunPhpFileFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "file", 0, data.Mixed{}),
	}
}

func (f *RunPhpFileFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	fileVal, ok := ctx.GetIndexValue(0)
	if !ok || fileVal == nil {
		return data.NewBoolValue(false), nil
	}
	sv, ok := fileVal.(data.AsString)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	file := sv.AsString()

	vm := ctx.GetVM()
	result, ctrl := vm.RunCompiledFile(file)
	if ctrl != nil {
		return nil, ctrl
	}
	if result == nil {
		return data.NewNullValue(), nil
	}
	if v, ok := result.(data.Value); ok {
		return v, nil
	}
	return data.NewNullValue(), nil
}
