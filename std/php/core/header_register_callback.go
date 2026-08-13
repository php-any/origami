package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

var (
	headerCallbacks     []data.Value
	headerOutputStarted bool
)

// MarkHeaderOutputStarted 标记已有输出（此后注册的 header 回调无效）。
func MarkHeaderOutputStarted() {
	headerOutputStarted = true
}

// RunHeaderCallbacks 在请求/脚本结束时执行已注册的 header 回调。
func RunHeaderCallbacks(vm data.VM) {
	for _, cb := range headerCallbacks {
		if fv, ok := cb.(*data.FuncValue); ok {
			vars := fv.Value.GetVariables()
			ctx := vm.CreateContext(vars)
			_, _ = fv.Call(ctx)
		}
	}
	headerCallbacks = nil
}

type HeaderRegisterCallbackFunction struct{}

func NewHeaderRegisterCallbackFunction() data.FuncStmt {
	return &HeaderRegisterCallbackFunction{}
}

func (f *HeaderRegisterCallbackFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	if headerOutputStarted {
		return data.NewNullValue(), nil
	}
	cb, _ := ctx.GetIndexValue(0)
	if cb != nil {
		headerCallbacks = append(headerCallbacks, cb.(data.Value))
	}
	return data.NewNullValue(), nil
}

func (f *HeaderRegisterCallbackFunction) GetName() string { return "header_register_callback" }

func (f *HeaderRegisterCallbackFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
	}
}

func (f *HeaderRegisterCallbackFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, data.NewBaseType("callable")),
	}
}
