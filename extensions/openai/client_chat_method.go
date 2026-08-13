package openai

import (
	"errors"

	oai "github.com/openai/openai-go/v3"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ClientChatMethod 是 OpenAI\Client 的 chat 方法
type ClientChatMethod struct {
	source *client
}

func (m *ClientChatMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if m.source.inner == nil {
		return nil, utils.NewThrow(errors.New("OpenAI client not initialized, call __construct() first"))
	}

	// 参数 1: model (string)
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("chat() requires at least 2 arguments (model, messages)"))
	}
	model, ok := a0.(data.AsString)
	if !ok {
		return nil, utils.NewThrow(errors.New("chat() argument 1 (model) must be a string"))
	}

	// 参数 2: messages (array)
	a1, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("chat() requires at least 2 arguments (model, messages)"))
	}
	messages, ctl := convertMessages(a1)
	if ctl != nil {
		return nil, ctl
	}

	// 参数 3: options (可选关联数组)
	var opts map[string]any
	if v, ok := ctx.GetIndexValue(2); ok {
		if obj, ok := v.(*data.ObjectValue); ok {
			opts = objectToMap(obj)
		}
	}

	result, err := m.source.chat(model.AsString(), messages, opts)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	return NewChatCompletionClassValue(result, ctx), nil
}

func (m *ClientChatMethod) GetName() string            { return "chat" }
func (m *ClientChatMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *ClientChatMethod) GetIsStatic() bool          { return false }
func (m *ClientChatMethod) GetReturnType() data.Types {
	return data.NewBaseType("OpenAI\\ChatCompletion")
}

func (m *ClientChatMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "model", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "messages", 1, nil, data.NewBaseType("array")),
		node.NewParameter(nil, "options", 2, data.NewNullValue(), data.NewBaseType("array")),
	}
}

func (m *ClientChatMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "model", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "messages", 1, data.NewBaseType("array")),
		node.NewVariable(nil, "options", 2, data.NewBaseType("array")),
	}
}

// convertMessages 将消息列表转换为 OpenAI 消息参数列表
// 支持两种格式:
//   - 原始数组: [["role" => "user", "content" => "Hi"]]
//   - 消息对象: [new OpenAI\UserMessage("Hi")]
//
// 参考 openai-go README: 支持 system / user / assistant / developer / tool 角色
func convertMessages(v data.Value) ([]oai.ChatCompletionMessageParamUnion, data.Control) {
	arr, ok := v.(*data.ArrayValue)
	if !ok {
		return nil, utils.NewThrow(errors.New("messages must be an array"))
	}

	var messages []oai.ChatCompletionMessageParamUnion
	for _, item := range arr.ToValueList() {
		role, content, toolCallID, ctl := extractMessage(item)
		if ctl != nil {
			return nil, ctl
		}

		switch role {
		case "system":
			messages = append(messages, oai.SystemMessage(content))
		case "user":
			messages = append(messages, oai.UserMessage(content))
		case "assistant":
			messages = append(messages, oai.AssistantMessage(content))
		case "developer":
			messages = append(messages, oai.DeveloperMessage(content))
		case "tool":
			messages = append(messages, oai.ToolMessage(toolCallID, content))
		default:
			messages = append(messages, oai.UserMessage(content))
		}
	}

	return messages, nil
}

// extractMessage 从一条消息中提取 role / content / toolCallID
// 支持 ClassValue (消息对象) 和 ObjectValue (原始数组)
func extractMessage(item data.Value) (role, content, toolCallID string, ctl data.Control) {
	// 先尝试作为 ClassValue 处理（new OpenAI\UserMessage("...") 等消息对象）
	if cv, ok := item.(*data.ClassValue); ok {
		// 从类名推导 role
		role = roleFromClass(cv)
		// 从属性获取 content
		c, _ := cv.GetProperty("content")
		if c != nil {
			if _, isNull := c.(*data.NullValue); !isNull {
				content = c.AsString()
			}
		}
		// tool_call_id (仅 ToolMessage)
		if id, _ := cv.GetProperty("tool_call_id"); id != nil {
			if _, isNull := id.(*data.NullValue); !isNull {
				toolCallID = id.AsString()
			}
		}
		if content == "" {
			return "", "", "", utils.NewThrow(errors.New("message object must have 'content'"))
		}
		return
	}

	// 回退: 作为 ObjectValue 处理（["role" => "user", "content" => "Hi"] 格式）
	obj, ok := item.(*data.ObjectValue)
	if !ok {
		return "", "", "", utils.NewThrow(errors.New("each message must be a message object or an associative array with 'role' and 'content'"))
	}

	roleVal, _ := obj.GetProperty("role")
	contentVal, _ := obj.GetProperty("content")
	if roleVal == nil || contentVal == nil {
		return "", "", "", utils.NewThrow(errors.New("each message must have 'role' and 'content' fields"))
	}
	role = roleVal.AsString()
	content = contentVal.AsString()

	if idVal, _ := obj.GetProperty("tool_call_id"); idVal != nil {
		if _, isNull := idVal.(*data.NullValue); !isNull {
			toolCallID = idVal.AsString()
		}
	}
	return
}

// roleFromClass 从消息类的 ClassValue 推导 role 字符串
func roleFromClass(cv *data.ClassValue) string {
	switch cv.Class.GetName() {
	case "OpenAI\\SystemMessage":
		return "system"
	case "OpenAI\\UserMessage":
		return "user"
	case "OpenAI\\AssistantMessage":
		return "assistant"
	case "OpenAI\\DeveloperMessage":
		return "developer"
	case "OpenAI\\ToolMessage":
		return "tool"
	default:
		return "user"
	}
}

// objectToMap 将 ObjectValue 递归转换为 map[string]any，支持嵌套对象和数组
func objectToMap(obj *data.ObjectValue) map[string]any {
	result := make(map[string]any)
	obj.RangeProperties(func(key string, val data.Value) bool {
		result[key] = valueToAny(val)
		return true
	})
	return result
}

// valueToAny 将 data.Value 递归转换为 Go any 类型
func valueToAny(val data.Value) any {
	switch v := val.(type) {
	case *data.StringValue:
		return v.AsString()
	case *data.IntValue:
		n, _ := v.AsInt()
		return n
	case *data.FloatValue:
		f, _ := v.AsFloat()
		return f
	case *data.BoolValue:
		b, _ := v.AsBool()
		return b
	case *data.NullValue:
		return nil
	case *data.ObjectValue:
		return objectToMap(v)
	case *data.ArrayValue:
		items := make([]any, 0)
		for _, item := range v.ToValueList() {
			items = append(items, valueToAny(item))
		}
		return items
	default:
		return v.AsString()
	}
}
