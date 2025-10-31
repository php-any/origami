package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type HeaderDelMethod struct {
	source *httpsrc.Header
}

func (h *HeaderDelMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	h.source.Del(param0)
	return nil, nil
}

func (h *HeaderDelMethod) GetName() string            { return "del" }
func (h *HeaderDelMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *HeaderDelMethod) GetIsStatic() bool          { return false }
func (h *HeaderDelMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *HeaderDelMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *HeaderDelMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
