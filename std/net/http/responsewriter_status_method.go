package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterStatusMethod struct {
	source httpsrc.ResponseWriter
}

func (h *ResponseWriterStatusMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, err := utils.ConvertFromIndex[int](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	h.source.WriteHeader(param0)
	// 返回 Response 自身以支持链式调用
	return data.NewProxyValue(NewResponseWriterClassFrom(h.source), ctx.CreateBaseContext()), nil
}

func (h *ResponseWriterStatusMethod) GetName() string            { return "status" }
func (h *ResponseWriterStatusMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterStatusMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterStatusMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "statusCode", 0, nil, data.NewBaseType("int")),
	}
}
func (h *ResponseWriterStatusMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "statusCode", 0, nil),
	}
}
func (h *ResponseWriterStatusMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
