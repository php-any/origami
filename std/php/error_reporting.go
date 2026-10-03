package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewErrorReportingFunction 创建 error_reporting 函数。
// PHP 语义（简化版）：
//
//	error_reporting(level?: int): int
//
// - 不带参数时返回当前错误报告级别；
// - 带参数时设置新的错误报告级别并返回旧值。
func NewErrorReportingFunction() data.FuncStmt {
	return &ErrorReportingFunction{}
}

type ErrorReportingFunction struct {
	data.Function
}

// E_ALL 常量值
const E_ALL = 32767

func (f *ErrorReportingFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	state := data.ErrorState(ctx)
	if level, ok := ctx.GetIndexValue(0); ok {
		if level, ok := level.(*data.IntValue); ok {
			return data.NewIntValue(state.SetReporting(level.Value)), nil
		}
	}
	return data.NewIntValue(state.Reporting()), nil
}

func (f *ErrorReportingFunction) GetName() string {
	return "error_reporting"
}

func (f *ErrorReportingFunction) GetParams() []data.GetValue {
	return errorReportingParams
}

var errorReportingParams = []data.GetValue{node.NewParameter(nil, "level", 0, data.NewNullValue(), data.NewNullableType(data.Int{}))}

var errorReportingFunctionGetVariables = []data.Variable{
	data.NewVariable("level", 0, data.NewBaseType("int")),
}

func (f *ErrorReportingFunction) GetVariables() []data.Variable {
	return errorReportingFunctionGetVariables
}
