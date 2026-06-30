package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// RequestRouteMethod 获取路由参数（path values）
// 不带参数时返回所有路由参数，带参数时返回指定键的值
type RequestRouteMethod struct {
	source *httpsrc.Request
}

func (h *RequestRouteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewAnyValue(nil), nil
	}

	pathVals := collectPathValues(h.source)

	// 检查是否有参数
	_, hasKey := ctx.GetIndexValue(0)

	// 如果没有参数，返回所有路由参数
	if !hasKey {
		result := data.NewObjectValue()
		for key, val := range pathVals {
			result.SetProperty(key, data.NewStringValue(val))
		}
		return result, nil
	}

	// 如果有参数，返回指定键的值
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	if pathVals != nil {
		if val, exists := pathVals[param0]; exists {
			return data.NewStringValue(val), nil
		}
	}

	return data.NewAnyValue(nil), nil
}

func (h *RequestRouteMethod) GetName() string            { return "route" }
func (h *RequestRouteMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestRouteMethod) GetIsStatic() bool          { return false }
func (h *RequestRouteMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "key", 0, nil, nil),
	}
}
func (h *RequestRouteMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "key", 0, nil),
	}
}
func (h *RequestRouteMethod) GetReturnType() data.Types {
	return data.NewUnionType([]data.Types{
		data.NewBaseType("array"),
		data.NewBaseType("string"),
		data.NewBaseType("null"),
	})
}
