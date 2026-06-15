package http

import (
	"github.com/php-any/origami/std/serializer/json"
	"github.com/php-any/origami/utils"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type GetProperties interface {
	GetProperties() map[string]data.Value
}

type ResponseWriterJsonMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterJsonMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	param0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrowf("json 方法缺少参数: %v", 0)
	}

	msg, ok := param0.(data.ValueSerializer)
	if !ok {
		return nil, utils.NewThrowf("使用未支持json序列化的结构%#v", param0)
	}

	bytes, err := msg.Marshal(json.NewJsonSerializer())
	if err != nil {
		return nil, utils.NewThrow(err)
	}
	if err := h.w.WriteJSON(bytes); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
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
