package openai

import (
	"github.com/php-any/origami/data"
)

// ClientClass 是 OpenAI\Client 的类定义
type ClientClass struct {
	source              *client
	construct           data.Method
	chatMethod          data.Method
	embeddingsMethod    data.Method
	imagesMethod        data.Method
	speechMethod        data.Method
	transcriptionMethod data.Method
}

// NewClientClass 创建 OpenAI\Client 类
func NewClientClass() data.ClassStmt {
	source := &client{}
	return &ClientClass{
		source:              source,
		construct:           &ClientConstructMethod{source},
		chatMethod:          &ClientChatMethod{source},
		embeddingsMethod:    &ClientEmbeddingsMethod{source},
		imagesMethod:        &ClientImagesMethod{source},
		speechMethod:        &ClientSpeechMethod{source},
		transcriptionMethod: &ClientTranscriptionMethod{source},
	}
}

func (c *ClientClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *ClientClass) GetFrom() data.From { return nil }

func (c *ClientClass) GetName() string { return "OpenAI\\Client" }

func (c *ClientClass) GetExtend() *string { return nil }

func (c *ClientClass) GetImplements() []string { return nil }

func (c *ClientClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (c *ClientClass) GetPropertyList() []data.Property {
	return nil
}

func (c *ClientClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "chat":
		return c.chatMethod, true
	case "embeddings":
		return c.embeddingsMethod, true
	case "images":
		return c.imagesMethod, true
	case "speech":
		return c.speechMethod, true
	case "transcription":
		return c.transcriptionMethod, true
	default:
		return nil, false
	}
}

func (c *ClientClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}

func (c *ClientClass) GetMethods() []data.Method {
	return []data.Method{
		c.chatMethod,
		c.embeddingsMethod,
		c.imagesMethod,
		c.speechMethod,
		c.transcriptionMethod,
	}
}

func (c *ClientClass) GetConstruct() data.Method {
	return c.construct
}
