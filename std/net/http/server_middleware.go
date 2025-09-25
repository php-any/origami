package http

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type ServerMiddlewareMethod struct {
	server *ServerClass
}

func (h *ServerMiddlewareMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// middleware 接受一个闭包: function(r, w, next)
	v, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("需要传入中间件闭包"))
	}
	if fv, ok := v.(*data.FuncValue); ok {
		mw, err := newMiddleware(fv.Value, ctx)
		if err != nil {
			return nil, data.NewErrorThrow(nil, err)
		}
		h.server.Middlewares = append(h.server.Middlewares, mw)
		return nil, nil
	}
	return nil, data.NewErrorThrow(nil, errors.New("参数必须是中间件闭包"))
}

func (h *ServerMiddlewareMethod) GetName() string            { return "middleware" }
func (h *ServerMiddlewareMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerMiddlewareMethod) GetIsStatic() bool          { return false }
func (h *ServerMiddlewareMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "mid", 0, nil, nil),
	}
}
func (h *ServerMiddlewareMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "mid", 0, nil),
	}
}
func (h *ServerMiddlewareMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
