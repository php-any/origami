package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterCookieMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterCookieMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	c, err := cookieFromPHP(ctx)
	if err != nil {
		return nil, utils.NewThrowf("cookie 参数无效: %v", err)
	}
	h.w.SetCookie(c)
	return nil, nil
}

func (h *ResponseWriterCookieMethod) GetName() string            { return "cookie" }
func (h *ResponseWriterCookieMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterCookieMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterCookieMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "nameOrCookie", 0, nil, nil),
		node.NewParameter(nil, "value", 1, nil, nil),
		node.NewParameter(nil, "options", 2, nil, nil),
	}
}
func (h *ResponseWriterCookieMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "nameOrCookie", 0, nil),
		node.NewVariable(nil, "value", 1, nil),
		node.NewVariable(nil, "options", 2, nil),
	}
}
func (h *ResponseWriterCookieMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
