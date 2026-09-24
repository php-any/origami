package http

import (
	httpsrc "net/http"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// RequestIsMethodMethod 检查请求方法
type RequestIsMethodMethod struct {
	source *httpsrc.Request
}

func (h *RequestIsMethodMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewBoolValue(false), nil
	}

	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	method := strings.ToUpper(param0)
	return data.NewBoolValue(h.source.Method == method), nil
}

func (h *RequestIsMethodMethod) GetName() string            { return "isMethod" }
func (h *RequestIsMethodMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestIsMethodMethod) GetIsStatic() bool          { return false }
var requestIsMethodMethodGetParams = []data.GetValue{
	node.NewParameter(nil, "method", 0, nil, nil),
}

func (h *RequestIsMethodMethod) GetParams() []data.GetValue {
	return requestIsMethodMethodGetParams
}
var requestIsMethodMethodGetVariables = []data.Variable{
	node.NewVariable(nil, "method", 0, nil),
}

func (h *RequestIsMethodMethod) GetVariables() []data.Variable {
	return requestIsMethodMethodGetVariables
}
func (h *RequestIsMethodMethod) GetReturnType() data.Types { return data.NewBaseType("bool") }
