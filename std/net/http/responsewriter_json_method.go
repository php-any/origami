package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/std/serializer/json"
	"github.com/php-any/origami/utils"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type GetProperties interface {
	GetProperties() map[string]data.Value
}

type ResponseWriterJsonMethod struct {
	source httpsrc.ResponseWriter
}

func (h *ResponseWriterJsonMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrowf("json 方法缺少参数: %v", 0)
	}

	switch msg := param0.(type) {
	case data.ValueSerializer:
		bytes, err := msg.Marshal(json.NewJsonSerializer())
		if err != nil {
			return nil, utils.NewThrow(err)
		}
		h.source.Header().Set("Content-Type", "application/json; charset=utf-8")
		ret0, ret1 := h.source.Write(bytes)
		if ret1 != nil {
			return nil, utils.NewThrow(ret1)
		}
		return data.NewIntValue(ret0), nil
	}

	return nil, utils.NewThrowf("使用未支持json序列化的结构%#v", param0)
}

func (h *ResponseWriterJsonMethod) GetName() string            { return "json" }
func (h *ResponseWriterJsonMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterJsonMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterJsonMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "data", 0, nil, data.Object{}),
	}
}
func (h *ResponseWriterJsonMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "data", 0, data.Object{}),
	}
}
func (h *ResponseWriterJsonMethod) GetReturnType() data.Types { return data.NewBaseType("void") }

func propertiesToJsonString(properties map[string]data.Value) []byte {
	return []byte("")
}
