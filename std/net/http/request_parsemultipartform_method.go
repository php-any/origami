package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type RequestParseMultipartFormMethod struct {
	source *httpsrc.Request
}

func (h *RequestParseMultipartFormMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[int64](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	ret0 := h.source.ParseMultipartForm(param0)
	return data.NewAnyValue(ret0), nil
}

func (h *RequestParseMultipartFormMethod) GetName() string            { return "parseMultipartForm" }
func (h *RequestParseMultipartFormMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestParseMultipartFormMethod) GetIsStatic() bool          { return false }
func (h *RequestParseMultipartFormMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *RequestParseMultipartFormMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *RequestParseMultipartFormMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
