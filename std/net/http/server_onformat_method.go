package http

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ServerOnFormatMethod struct {
	server *ServerClass
}

func (h *ServerOnFormatMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("需要传入格式化闭包"))
	}
	fv, ok := v.(*data.FuncValue)
	if !ok {
		return nil, utils.NewThrow(errors.New("onFormat 参数必须是闭包"))
	}
	if len(fv.Value.GetVariables()) < 3 {
		return nil, utils.NewThrow(errors.New("onFormat 闭包需要 3 个参数: ($code, $message, $data)"))
	}
	h.server.formatHandler = &formatHandlerSlot{
		fn:  fv.Value,
		ctx: ctx,
	}
	return nil, nil
}

func (h *ServerOnFormatMethod) GetName() string            { return "onFormat" }
func (h *ServerOnFormatMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ServerOnFormatMethod) GetIsStatic() bool          { return false }
func (h *ServerOnFormatMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "handler", 0, nil, nil),
	}
}
func (h *ServerOnFormatMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "handler", 0, nil),
	}
}
func (h *ServerOnFormatMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
