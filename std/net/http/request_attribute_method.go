package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type RequestAttributeMethod struct {
	source *httpsrc.Request
}

func (h *RequestAttributeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	key, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	bag := requestAttrs(h.source)
	if bag == nil {
		return nil, utils.NewThrowf("无法访问 request attributes")
	}

	value, hasValue := ctx.GetIndexValue(1)
	if !hasValue {
		val, ok := bag[key]
		if !ok || val == nil {
			return data.NewNullValue(), nil
		}
		return val, nil
	}

	if _, isNull := value.(*data.NullValue); isNull {
		delete(bag, key)
		return nil, nil
	}
	bag[key] = value.(data.Value)
	return nil, nil
}

func (h *RequestAttributeMethod) GetName() string            { return "attribute" }
func (h *RequestAttributeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestAttributeMethod) GetIsStatic() bool          { return false }
func (h *RequestAttributeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "key", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "value", 1, nil, nil),
	}
}
func (h *RequestAttributeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "key", 0, nil),
		node.NewVariable(nil, "value", 1, nil),
	}
}
func (h *RequestAttributeMethod) GetReturnType() data.Types { return data.NewBaseType("mixed") }
