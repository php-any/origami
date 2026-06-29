package openai

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ClientSpeechMethod 是 OpenAI\Client 的 speech 方法
type ClientSpeechMethod struct {
	source *client
}

func (m *ClientSpeechMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.inner == nil {
		return nil, utils.NewThrow(errors.New("OpenAI client not initialized, call __construct() first"))
	}

	// 参数 1: model (string)
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("speech() requires at least 4 arguments (model, input, voice, output)"))
	}
	model, ok := a0.(data.AsString)
	if !ok {
		return nil, utils.NewThrow(errors.New("speech() argument 1 (model) must be a string"))
	}

	// 参数 2: input (string)
	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("speech() requires at least 4 arguments (model, input, voice, output)"))
	}
	input, ok := a1.(data.AsString)
	if !ok {
		return nil, utils.NewThrow(errors.New("speech() argument 2 (input) must be a string"))
	}

	// 参数 3: voice (string)
	a2, ok := ctx.GetIndexValue(2)
	if !ok {
		return nil, utils.NewThrow(errors.New("speech() requires at least 4 arguments (model, input, voice, output)"))
	}
	voice, ok := a2.(data.AsString)
	if !ok {
		return nil, utils.NewThrow(errors.New("speech() argument 3 (voice) must be a string"))
	}

	// 参数 4: output (string - 文件路径)
	a3, ok := ctx.GetIndexValue(3)
	if !ok {
		return nil, utils.NewThrow(errors.New("speech() requires at least 4 arguments (model, input, voice, output)"))
	}
	output, ok := a3.(data.AsString)
	if !ok {
		return nil, utils.NewThrow(errors.New("speech() argument 4 (output) must be a string"))
	}

	// 参数 5: options (可选)
	var opts map[string]any
	if v, ok := ctx.GetIndexValue(4); ok {
		if obj, ok := v.(*data.ObjectValue); ok {
			opts = objectToMap(obj)
		}
	}

	err := m.source.speech(model.AsString(), input.AsString(), voice.AsString(), output.AsString(), opts)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	return nil, nil
}

func (m *ClientSpeechMethod) GetName() string            { return "speech" }
func (m *ClientSpeechMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ClientSpeechMethod) GetIsStatic() bool          { return false }
func (m *ClientSpeechMethod) GetReturnType() data.Types  { return data.NewBaseType("void") }

func (m *ClientSpeechMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "model", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "input", 1, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "voice", 2, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "output", 3, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "options", 4, data.NewNullValue(), data.NewBaseType("array")),
	}
}

func (m *ClientSpeechMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "model", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "input", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "voice", 2, data.NewBaseType("string")),
		node.NewVariable(nil, "output", 3, data.NewBaseType("string")),
		node.NewVariable(nil, "options", 4, data.NewBaseType("array")),
	}
}
