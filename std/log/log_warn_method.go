package log

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type LogWarnMethod struct {
	source *Log
}

func (h *LogWarnMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 0"))
	}

	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少参数, index: 1"))
	}

	h.source.Warn(a0.(*data.StringValue).AsString(), *a1.(*data.ArrayValue))
	return nil, nil
}

func (h *LogWarnMethod) GetName() string {
	return "warn"
}

func (h *LogWarnMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (h *LogWarnMethod) GetIsStatic() bool {
	return true
}

func (h *LogWarnMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "msg", 0, nil, nil),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}

func (h *LogWarnMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "msg", 0, nil),
		node.NewVariable(nil, "args", 1, nil),
	}
}

// GetReturnType 返回方法返回类型
func (h *LogWarnMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}
