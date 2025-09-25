package http

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
	httpsrc "net/http"
)

type RequestCookieMethod struct {
	source *httpsrc.Request
}

func (h *RequestCookieMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}

	ret0, ret1 := h.source.Cookie(param0)
	return data.NewArrayValue([]data.Value{data.NewAnyValue(ret0), data.NewAnyValue(ret1)}), nil
}

func (h *RequestCookieMethod) GetName() string            { return "cookie" }
func (h *RequestCookieMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestCookieMethod) GetIsStatic() bool          { return true }
func (h *RequestCookieMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *RequestCookieMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *RequestCookieMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
