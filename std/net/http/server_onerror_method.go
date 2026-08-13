package http

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ServerOnErrorMethod struct {
	server *ServerClass
}

func (h *ServerOnErrorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("需要传入错误处理闭包"))
	}
	fv, ok := v.(*data.FuncValue)
	if !ok {
		return nil, utils.NewThrow(errors.New("onError 参数必须是闭包"))
	}
	if len(fv.Value.GetVariables()) < 3 {
		return nil, utils.NewThrow(errors.New("onError 闭包需要 3 个参数: ($request, $response, $error)"))
	}
	h.server.errorHandler = &errorHandlerSlot{
		fn:  fv.Value,
		ctx: ctx,
	}
	return nil, nil
}

func (h *ServerOnErrorMethod) GetName() string            { return "onError" }
func (h *ServerOnErrorMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerOnErrorMethod) GetIsStatic() bool          { return false }
func (h *ServerOnErrorMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "handler", 0, nil, nil),
	}
}
func (h *ServerOnErrorMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "handler", 0, nil),
	}
}
func (h *ServerOnErrorMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
