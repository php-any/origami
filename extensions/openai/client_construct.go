package openai

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
	"github.com/php-any/origami/utils"
)

// ClientConstructMethod 是 OpenAI\Client 的构造函数
type ClientConstructMethod struct {
	source *client
}

func (m *ClientConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取 apiKey 参数
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("OpenAI\\Client.__construct() requires at least 1 argument (apiKey)"))
	}
	apiKey, ok := a0.(data.AsString)
	if !ok {
		return nil, utils.NewThrow(errors.New("OpenAI\\Client.__construct() argument 1 must be a string"))
	}

	// 获取可选的 baseURL 参数
	baseURL := ""
	if v, ok := ctx.GetIndexValue(1); ok {
		if s, ok := v.(data.AsString); ok {
			baseURL = s.AsString()
		}
	}

	c, err := newClient(apiKey.AsString(), baseURL)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	// 将创建的客户端赋值到共享的 source 上
	m.source.inner = c.inner

	return nil, nil
}

func (m *ClientConstructMethod) GetName() string            { return token.ConstructName }
func (m *ClientConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ClientConstructMethod) GetIsStatic() bool          { return false }
func (m *ClientConstructMethod) GetReturnType() data.Types  { return nil }

func (m *ClientConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "apiKey", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "baseURL", 1, data.NewStringValue(""), data.NewBaseType("string")),
	}
}

func (m *ClientConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "apiKey", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "baseURL", 1, data.NewBaseType("string")),
	}
}
