package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type RequestAddCookieMethod struct {
	source *httpsrc.Request
}

func (h *RequestAddCookieMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[*httpsrc.Cookie](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	h.source.AddCookie(param0)
	return nil, nil
}

func (h *RequestAddCookieMethod) GetName() string            { return "addCookie" }
func (h *RequestAddCookieMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestAddCookieMethod) GetIsStatic() bool          { return false }
func (h *RequestAddCookieMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *RequestAddCookieMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *RequestAddCookieMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
