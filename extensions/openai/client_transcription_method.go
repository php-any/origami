package openai

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ClientTranscriptionMethod 是 OpenAI\Client 的 transcription 方法
type ClientTranscriptionMethod struct {
	source *client
}

func (m *ClientTranscriptionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.inner == nil {
		return nil, utils.NewThrow(errors.New("OpenAI client not initialized, call __construct() first"))
	}

	// 参数 1: model (string)
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("transcription() requires at least 2 arguments (model, file)"))
	}
	model, ok := a0.(data.AsString)
	if !ok {
		return nil, utils.NewThrow(errors.New("transcription() argument 1 (model) must be a string"))
	}

	// 参数 2: file (string - 文件路径)
	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("transcription() requires at least 2 arguments (model, file)"))
	}
	filePath, ok := a1.(data.AsString)
	if !ok {
		return nil, utils.NewThrow(errors.New("transcription() argument 2 (file) must be a string"))
	}

	// 参数 3: options (可选)
	var opts map[string]any
	if v, ok := ctx.GetIndexValue(2); ok {
		if obj, ok := v.(*data.ObjectValue); ok {
			opts = objectToMap(obj)
		}
	}

	result, err := m.source.transcription(model.AsString(), filePath.AsString(), opts)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	return NewTranscriptionResultClassValue(result, ctx), nil
}

func (m *ClientTranscriptionMethod) GetName() string            { return "transcription" }
func (m *ClientTranscriptionMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ClientTranscriptionMethod) GetIsStatic() bool          { return false }
func (m *ClientTranscriptionMethod) GetReturnType() data.Types {
	return data.NewBaseType("OpenAI\\TranscriptionResult")
}

func (m *ClientTranscriptionMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "model", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "file", 1, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "options", 2, data.NewNullValue(), data.NewBaseType("array")),
	}
}

func (m *ClientTranscriptionMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "model", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "file", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "options", 2, data.NewBaseType("array")),
	}
}
