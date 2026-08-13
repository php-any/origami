package openai

import (
	oai "github.com/openai/openai-go/v3"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewTranscriptionResultClass 创建 OpenAI\TranscriptionResult 类定义
func NewTranscriptionResultClass() data.ClassStmt {
	return &TranscriptionResultClass{}
}

// NewTranscriptionResultClassValue 从 OpenAI 响应创建 TranscriptionResult 类实例
func NewTranscriptionResultClassValue(resp *oai.AudioTranscriptionNewResponseUnion, ctx data.Context) *data.ClassValue {
	class := &TranscriptionResultClass{resp: resp}
	cv := data.NewClassValue(class, ctx)

	cv.SetProperty("text", data.NewStringValue(resp.Text))
	cv.SetProperty("language", data.NewStringValue(resp.Language))
	cv.SetProperty("duration", data.NewFloatValue(resp.Duration))

	return cv
}

// TranscriptionResultClass 是 OpenAI\TranscriptionResult 的类定义
type TranscriptionResultClass struct {
	resp *oai.AudioTranscriptionNewResponseUnion
}

func (c *TranscriptionResultClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if c.resp != nil {
		return NewTranscriptionResultClassValue(c.resp, ctx), nil
	}
	return data.NewClassValue(c, ctx), nil
}

func (c *TranscriptionResultClass) GetFrom() data.From      { return nil }
func (c *TranscriptionResultClass) GetName() string         { return "OpenAI\\TranscriptionResult" }
func (c *TranscriptionResultClass) GetExtend() *string      { return nil }
func (c *TranscriptionResultClass) GetImplements() []string { return nil }

func (c *TranscriptionResultClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (c *TranscriptionResultClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "text", "public", false, nil, data.NewBaseType("string")),
		node.NewProperty(nil, "language", "public", false, nil, data.NewBaseType("string")),
		node.NewProperty(nil, "duration", "public", false, nil, data.NewBaseType("float")),
	}
}

func (c *TranscriptionResultClass) GetMethod(name string) (data.Method, bool) {
	return nil, false
}

func (c *TranscriptionResultClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}

func (c *TranscriptionResultClass) GetMethods() []data.Method { return nil }

func (c *TranscriptionResultClass) GetConstruct() data.Method { return nil }
