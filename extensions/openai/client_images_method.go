package openai

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ClientImagesMethod 是 OpenAI\Client 的 images 方法
type ClientImagesMethod struct {
	source *client
}

func (m *ClientImagesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.inner == nil {
		return nil, utils.NewThrow(errors.New("OpenAI client not initialized, call __construct() first"))
	}

	// 参数 1: model (string)
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("images() requires at least 2 arguments (model, prompt)"))
	}
	model, ok := a0.(data.AsString)
	if !ok {
		return nil, utils.NewThrow(errors.New("images() argument 1 (model) must be a string"))
	}

	// 参数 2: prompt (string)
	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("images() requires at least 2 arguments (model, prompt)"))
	}
	prompt, ok := a1.(data.AsString)
	if !ok {
		return nil, utils.NewThrow(errors.New("images() argument 2 (prompt) must be a string"))
	}

	// 参数 3: options (可选)
	var opts map[string]any
	if v, ok := ctx.GetIndexValue(2); ok {
		if obj, ok := v.(*data.ObjectValue); ok {
			opts = objectToMap(obj)
		}
	}

	result, err := m.source.images(model.AsString(), prompt.AsString(), opts)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	return NewImageResultClassValue(result, ctx), nil
}

func (m *ClientImagesMethod) GetName() string            { return "images" }
func (m *ClientImagesMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ClientImagesMethod) GetIsStatic() bool          { return false }
func (m *ClientImagesMethod) GetReturnType() data.Types {
	return data.NewBaseType("OpenAI\\ImageResult")
}

func (m *ClientImagesMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "model", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "prompt", 1, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "options", 2, data.NewNullValue(), data.NewBaseType("array")),
	}
}

func (m *ClientImagesMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "model", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "prompt", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "options", 2, data.NewBaseType("array")),
	}
}
