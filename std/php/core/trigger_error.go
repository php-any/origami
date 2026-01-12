package core

import (
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

	// 第二个参数：错误级别（当前实现忽略，仅保留签名兼容）
	// PHP: E_USER_ERROR / E_USER_WARNING / E_USER_NOTICE 等
	level, _ := ctx.GetIndexValue(1)
	if i, ok := level.(data.AsInt); ok {
		if i, _ := i.AsInt(); i == 16384 {
			return nil, nil
		}
	}
	// 这里统一按致命错误处理，直接抛出异常

	return nil, utils.NewThrowf("trigger_error: %s", msg)
}

func (f *TriggerErrorFunction) GetName() string {
	return "trigger_error"
}

func (f *TriggerErrorFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "message", 0, nil, data.String{}),
		node.NewParameter(nil, "error_type", 1, nil, data.Int{}),
	}
}

func (f *TriggerErrorFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "message", 0, data.String{}),
		node.NewVariable(nil, "error_type", 1, data.Int{}),
	}
}
