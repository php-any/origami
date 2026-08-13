package openai

import (
	oai "github.com/openai/openai-go/v3"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewChatCompletionClass 创建 OpenAI\ChatCompletion 类定义
func NewChatCompletionClass() data.ClassStmt {
	return &ChatCompletionClass{}
}

// NewChatCompletionClassValue 从 OpenAI 响应创建 ChatCompletion 类实例
func NewChatCompletionClassValue(resp *oai.ChatCompletion, ctx data.Context) *data.ClassValue {
	class := &ChatCompletionClass{resp: resp}
	cv := data.NewClassValue(class, ctx)

	// 设置属性
	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		cv.SetProperty("content", data.NewStringValue(choice.Message.Content))
		cv.SetProperty("role", data.NewStringValue(string(choice.Message.Role)))
		cv.SetProperty("finishReason", data.NewStringValue(string(choice.FinishReason)))
	} else {
		cv.SetProperty("content", data.NewStringValue(""))
		cv.SetProperty("role", data.NewStringValue(""))
		cv.SetProperty("finishReason", data.NewStringValue(""))
	}

	// usage 信息
	usageObj := data.NewObjectValue()
	usageObj.SetProperty("promptTokens", data.NewIntValue(int(resp.Usage.PromptTokens)))
	usageObj.SetProperty("completionTokens", data.NewIntValue(int(resp.Usage.CompletionTokens)))
	usageObj.SetProperty("totalTokens", data.NewIntValue(int(resp.Usage.TotalTokens)))
	cv.SetProperty("usage", usageObj)

	return cv
}

// ChatCompletionClass 是 OpenAI\ChatCompletion 的类定义
type ChatCompletionClass struct {
	resp *oai.ChatCompletion
}

func (c *ChatCompletionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if c.resp != nil {
		return NewChatCompletionClassValue(c.resp, ctx), nil
	}
	return data.NewClassValue(c, ctx), nil
}

func (c *ChatCompletionClass) GetFrom() data.From      { return nil }
func (c *ChatCompletionClass) GetName() string         { return "OpenAI\\ChatCompletion" }
func (c *ChatCompletionClass) GetExtend() *string      { return nil }
func (c *ChatCompletionClass) GetImplements() []string { return nil }

func (c *ChatCompletionClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (c *ChatCompletionClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "content", "public", false, nil, data.NewBaseType("string")),
		node.NewProperty(nil, "role", "public", false, nil, data.NewBaseType("string")),
		node.NewProperty(nil, "finishReason", "public", false, nil, data.NewBaseType("string")),
		node.NewProperty(nil, "usage", "public", false, nil, data.NewBaseType("object")),
	}
}

func (c *ChatCompletionClass) GetMethod(name string) (data.Method, bool) {
	return nil, false
}

func (c *ChatCompletionClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}

func (c *ChatCompletionClass) GetMethods() []data.Method { return nil }

func (c *ChatCompletionClass) GetConstruct() data.Method { return nil }
