package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type RequestCookiesNamedMethod struct {
	source *httpsrc.Request
}

func (h *RequestCookiesNamedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	ret0 := h.source.CookiesNamed(param0)
	return data.NewAnyValue(ret0), nil
}

func (h *RequestCookiesNamedMethod) GetName() string            { return "cookiesNamed" }
func (h *RequestCookiesNamedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestCookiesNamedMethod) GetIsStatic() bool          { return false }
func (h *RequestCookiesNamedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *RequestCookiesNamedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *RequestCookiesNamedMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
