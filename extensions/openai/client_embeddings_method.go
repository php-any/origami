package openai

import (
	"errors"

	oai "github.com/openai/openai-go/v3"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ClientEmbeddingsMethod 是 OpenAI\Client 的 embeddings 方法
type ClientEmbeddingsMethod struct {
	source *client
}

func (m *ClientEmbeddingsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.inner == nil {
		return nil, utils.NewThrow(errors.New("OpenAI client not initialized, call __construct() first"))
	}

	// 参数 1: model (string)
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("embeddings() requires at least 2 arguments (model, input)"))
	}
	model, ok := a0.(data.AsString)
	if !ok {
		return nil, utils.NewThrow(errors.New("embeddings() argument 1 (model) must be a string"))
	}

	// 参数 2: input (string 或 array)
	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("embeddings() requires at least 2 arguments (model, input)"))
	}

	var input oai.EmbeddingNewParamsInputUnion
	switch v := a1.(type) {
	case *data.StringValue:
		input = oai.EmbeddingNewParamsInputUnion{OfString: oai.String(v.AsString())}
	case *data.ArrayValue:
		strs := make([]string, 0)
		for _, item := range v.ToValueList() {
			strs = append(strs, item.AsString())
		}
		input = oai.EmbeddingNewParamsInputUnion{OfArrayOfStrings: strs}
	default:
		input = oai.EmbeddingNewParamsInputUnion{OfString: oai.String(a1.AsString())}
	}

	// 参数 3: options (可选)
	var opts map[string]any
	if v, ok := ctx.GetIndexValue(2); ok {
		if obj, ok := v.(*data.ObjectValue); ok {
			opts = objectToMap(obj)
		}
	}

	result, err := m.source.embeddings(model.AsString(), input, opts)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	return NewEmbeddingResultClassValue(result, ctx), nil
}

func (m *ClientEmbeddingsMethod) GetName() string            { return "embeddings" }
func (m *ClientEmbeddingsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ClientEmbeddingsMethod) GetIsStatic() bool          { return false }
func (m *ClientEmbeddingsMethod) GetReturnType() data.Types {
	return data.NewBaseType("OpenAI\\EmbeddingResult")
}

func (m *ClientEmbeddingsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "model", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "input", 1, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "options", 2, data.NewNullValue(), data.NewBaseType("array")),
	}
}

func (m *ClientEmbeddingsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "model", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "input", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "options", 2, data.NewBaseType("array")),
	}
}
