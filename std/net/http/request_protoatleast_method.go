package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type RequestProtoAtLeastMethod struct {
	source *httpsrc.Request
}

func (h *RequestProtoAtLeastMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[int](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}
	param1, err := utils.ConvertFromIndex[int](ctx, 1)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	ret0 := h.source.ProtoAtLeast(param0, param1)
	return data.NewAnyValue(ret0), nil
}

func (h *RequestProtoAtLeastMethod) GetName() string            { return "protoAtLeast" }
func (h *RequestProtoAtLeastMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestProtoAtLeastMethod) GetIsStatic() bool          { return false }
func (h *RequestProtoAtLeastMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
		node.NewParameter(nil, "param1", 1, nil, nil),
	}
}
func (h *RequestProtoAtLeastMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
		node.NewVariable(nil, "param1", 1, nil),
	}
}
func (h *RequestProtoAtLeastMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
