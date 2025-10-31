package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type RequestPostFormValueMethod struct {
	source *httpsrc.Request
}

func (h *RequestPostFormValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	ret0 := h.source.PostFormValue(param0)
	return data.NewAnyValue(ret0), nil
}

func (h *RequestPostFormValueMethod) GetName() string            { return "postFormValue" }
func (h *RequestPostFormValueMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestPostFormValueMethod) GetIsStatic() bool          { return false }
func (h *RequestPostFormValueMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "param0", 0, nil, nil),
	}
}
func (h *RequestPostFormValueMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "param0", 0, nil),
	}
}
func (h *RequestPostFormValueMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
