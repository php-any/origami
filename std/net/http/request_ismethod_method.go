package http

import (
	"fmt"
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
		return nil, data.NewErrorThrow(nil, fmt.Errorf("参数转换失败: %v", err))
	}

	method := strings.ToUpper(param0)
	return data.NewBoolValue(h.source.Method == method), nil
}

func (h *RequestIsMethodMethod) GetName() string            { return "isMethod" }
func (h *RequestIsMethodMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestIsMethodMethod) GetIsStatic() bool          { return false }
func (h *RequestIsMethodMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "method", 0, nil, nil),
	}
}
func (h *RequestIsMethodMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "method", 0, nil),
	}
}
func (h *RequestIsMethodMethod) GetReturnType() data.Types { return data.NewBaseType("bool") }
