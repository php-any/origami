package core

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// TriggerErrorFunction 实现 trigger_error 函数
// 用于在运行时触发一个用户级错误，目前实现为直接抛出异常终止执行。
type TriggerErrorFunction struct{}

func NewTriggerErrorFunction() data.FuncStmt {
	return &TriggerErrorFunction{}
}

func (f *TriggerErrorFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 第一个参数：错误消息
	msgVal, _ := ctx.GetIndexValue(0)
	if msgVal == nil {
		return nil, utils.NewThrowf("trigger_error: message is required")
	}
	msg := msgVal.AsString()

	level := 1024
	if v, _ := ctx.GetIndexValue(1); v != nil {
		if i, ok := v.(*data.IntValue); ok {
			level = i.Value
		}
	}
	if level != 256 && level != 512 && level != 1024 && level != 16384 {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("trigger_error(): Argument #2 ($error_level) must be one of E_USER_ERROR, E_USER_WARNING, E_USER_NOTICE, or E_USER_DEPRECATED"), "ValueError")
	}
	return data.NewBoolValue(true), data.EmitPHPError(ctx, level, msg, nil)
}

func (f *TriggerErrorFunction) GetName() string {
	return "trigger_error"
}

var triggerErrorFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "message", 0, nil, data.String{}),
	node.NewParameter(nil, "error_type", 1, data.NewIntValue(1024), data.Int{}),
}

func (f *TriggerErrorFunction) GetParams() []data.GetValue {
	return triggerErrorFunctionGetParams
}

var triggerErrorFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "message", 0, data.String{}),
	node.NewVariable(nil, "error_type", 1, data.Int{}),
}

func (f *TriggerErrorFunction) GetVariables() []data.Variable {
	return triggerErrorFunctionGetVariables
}
