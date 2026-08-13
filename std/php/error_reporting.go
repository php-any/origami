package php

import (
	"github.com/php-any/origami/data"
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

// 全局错误报告级别（简化实现，实际应该存储在 VM 或上下文中）
var errorReportingLevel int = E_ALL

// E_ALL 常量值
const E_ALL = 32767

func (f *ErrorReportingFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	oldLevel := errorReportingLevel

	// 检查是否有参数
	levelVal, ok := ctx.GetIndexValue(0)
	if !ok {
		// 无参数，返回当前级别
		return data.NewIntValue(oldLevel), nil
	}

	// 有参数，转换为整数并设置新级别
	var newLevel int
	switch v := levelVal.(type) {
	case *data.IntValue:
		newLevel = v.Value
	case *data.FloatValue:
		newLevel = int(v.Value)
	case *data.StringValue:
		// 尝试从字符串解析整数（简化处理）
		newLevel = E_ALL
	case *data.NullValue:
		newLevel = E_ALL
	case *data.BoolValue:
		if v.Value {
			newLevel = E_ALL
		} else {
			newLevel = 0
		}
	default:
		newLevel = E_ALL
	}

	errorReportingLevel = newLevel
	return data.NewIntValue(oldLevel), nil
}

func (f *ErrorReportingFunction) GetName() string {
	return "error_reporting"
}

func (f *ErrorReportingFunction) GetParams() []data.GetValue {
	return nil
}

func (f *ErrorReportingFunction) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("level", 0, data.NewBaseType("int")),
	}
}
