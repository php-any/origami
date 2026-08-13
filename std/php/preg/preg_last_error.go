package preg

import (
	"github.com/php-any/origami/data"
)

// PregLastErrorFunction 实现 preg_last_error 函数
// 返回最后一次 preg 函数调用的错误码
// 由于我们的实现不跟踪 preg 错误状态，始终返回 PREG_NO_ERROR (0)
type PregLastErrorFunction struct{}

func NewPregLastErrorFunction() data.FuncStmt {
	return &PregLastErrorFunction{}
}

func (f *PregLastErrorFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewIntValue(0), nil
}

func (f *PregLastErrorFunction) GetName() string { return "preg_last_error" }

func (f *PregLastErrorFunction) GetParams() []data.GetValue { return nil }

func (f *PregLastErrorFunction) GetVariables() []data.Variable { return nil }

// PregLastErrorMsgFunction 实现 preg_last_error_msg 函数
// 返回最后一次 preg 函数调用的错误消息
// 由于我们的实现不跟踪 preg 错误状态，始终返回空字符串
type PregLastErrorMsgFunction struct{}

func NewPregLastErrorMsgFunction() data.FuncStmt {
	return &PregLastErrorMsgFunction{}
}

func (f *PregLastErrorMsgFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue("No error"), nil
}

func (f *PregLastErrorMsgFunction) GetName() string { return "preg_last_error_msg" }

func (f *PregLastErrorMsgFunction) GetParams() []data.GetValue { return nil }

func (f *PregLastErrorMsgFunction) GetVariables() []data.Variable { return nil }
