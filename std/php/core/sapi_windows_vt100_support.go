package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SapiWindowsVt100SupportFunction 实现 sapi_windows_vt100_support 函数
// sapi_windows_vt100_support(resource $stream, ?bool $enable = null): bool
// 检查或设置 Windows 控制台的 VT100 支持
type SapiWindowsVt100SupportFunction struct{}

func NewSapiWindowsVt100SupportFunction() data.FuncStmt {
	return &SapiWindowsVt100SupportFunction{}
}

func (f *SapiWindowsVt100SupportFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewBoolValue(false), nil
}

func (f *SapiWindowsVt100SupportFunction) GetName() string {
	return "sapi_windows_vt100_support"
}

func (f *SapiWindowsVt100SupportFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stream", 0, nil, nil),
		node.NewParameter(nil, "enable", 1, nil, data.Bool{}),
	}
}

func (f *SapiWindowsVt100SupportFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stream", 0, data.NewBaseType("resource")),
		node.NewVariable(nil, "enable", 1, data.Bool{}),
	}
}
