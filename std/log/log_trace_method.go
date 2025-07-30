package log

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type LogTraceMethod struct {
	source *Log
}

func (h *LogTraceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {

	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数, index: 1"))
	}

	h.source.Trace(a0.(*data.StringValue).AsString(), *a1.(*data.ArrayValue))
	return nil, nil
}

func (h *LogTraceMethod) GetName() string {
	return "trace"
}

func (h *LogTraceMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *LogTraceMethod) GetIsStatic() bool {
	return true
}

func (h *LogTraceMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "msg", 0, nil, nil),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}

func (h *LogTraceMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "msg", 0, nil),
		node.NewVariable(nil, "args", 1, nil),
	}
}

// GetReturnType 返回方法返回类型
func (h *LogTraceMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}
