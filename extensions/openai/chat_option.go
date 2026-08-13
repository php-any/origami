package openai

import (
	"github.com/php-any/origami/data"
)

// ============================================================
// ChatOption — chat() 方法的选项键枚举
//
// 使用方式:
//
//	$client->chat($model, $messages, [
//	    OpenAI\ChatOption::TEMPERATURE   => 0,
//	    OpenAI\ChatOption::MAX_TOKENS    => 100,
//	    OpenAI\ChatOption::RESPONSE_FORMAT => "json_object",
//	])
//
// ============================================================

type ChatOptionClass struct{}

func NewChatOptionClass() data.ClassStmt { return &ChatOptionClass{} }

func (c *ChatOptionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *ChatOptionClass) GetFrom() data.From                              { return nil }
func (c *ChatOptionClass) GetName() string                                 { return "OpenAI\\ChatOption" }
func (c *ChatOptionClass) GetExtend() *string                              { return nil }
func (c *ChatOptionClass) GetImplements() []string                         { return nil }
func (c *ChatOptionClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *ChatOptionClass) GetPropertyList() []data.Property                { return nil }
func (c *ChatOptionClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *ChatOptionClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *ChatOptionClass) GetMethods() []data.Method                       { return nil }
func (c *ChatOptionClass) GetConstruct() data.Method                       { return nil }

// opts 定义所有支持的 chat 选项，对应 openai-go SDK ChatCompletionNewParams 字段
var opts = map[string]string{
	"TEMPERATURE":           "temperature",
	"TOP_P":                 "topP",
	"MAX_TOKENS":            "maxTokens",
	"MAX_COMPLETION_TOKENS": "maxCompletionTokens",
	"STOP":                  "stop",
	"SEED":                  "seed",
	"N":                     "n",
	"FREQUENCY_PENALTY":     "frequencyPenalty",
	"PRESENCE_PENALTY":      "presencePenalty",
	"RESPONSE_FORMAT":       "responseFormat",
}

// GetStaticProperty 实现静态属性访问: OpenAI\ChatOption::TEMPERATURE → "temperature"
func (c *ChatOptionClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := opts[name]; ok {
		return data.NewStringValue(v), true
	}
	return nil, false
}
