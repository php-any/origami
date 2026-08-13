package openai

import (
	oai "github.com/openai/openai-go/v3"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewImageResultClass 创建 OpenAI\ImageResult 类定义
func NewImageResultClass() data.ClassStmt {
	return &ImageResultClass{}
}

// NewImageResultClassValue 从 OpenAI 响应创建 ImageResult 类实例
func NewImageResultClassValue(resp *oai.ImagesResponse, ctx data.Context) *data.ClassValue {
	class := &ImageResultClass{resp: resp}
	cv := data.NewClassValue(class, ctx)

	if len(resp.Data) > 0 {
		img := resp.Data[0]
		cv.SetProperty("url", data.NewStringValue(img.URL))
		if img.B64JSON != "" {
			cv.SetProperty("b64Json", data.NewStringValue(img.B64JSON))
		} else {
			cv.SetProperty("b64Json", data.NewStringValue(""))
		}
		if img.RevisedPrompt != "" {
			cv.SetProperty("revisedPrompt", data.NewStringValue(img.RevisedPrompt))
		} else {
			cv.SetProperty("revisedPrompt", data.NewStringValue(""))
		}
	} else {
		cv.SetProperty("url", data.NewStringValue(""))
		cv.SetProperty("b64Json", data.NewStringValue(""))
		cv.SetProperty("revisedPrompt", data.NewStringValue(""))
	}

	return cv
}

// ImageResultClass 是 OpenAI\ImageResult 的类定义
type ImageResultClass struct {
	resp *oai.ImagesResponse
}

func (c *ImageResultClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if c.resp != nil {
		return NewImageResultClassValue(c.resp, ctx), nil
	}
	return data.NewClassValue(c, ctx), nil
}

func (c *ImageResultClass) GetFrom() data.From      { return nil }
func (c *ImageResultClass) GetName() string         { return "OpenAI\\ImageResult" }
func (c *ImageResultClass) GetExtend() *string      { return nil }
func (c *ImageResultClass) GetImplements() []string { return nil }

func (c *ImageResultClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (c *ImageResultClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "url", "public", false, nil, data.NewBaseType("string")),
		node.NewProperty(nil, "b64Json", "public", false, nil, data.NewBaseType("string")),
		node.NewProperty(nil, "revisedPrompt", "public", false, nil, data.NewBaseType("string")),
	}
}

func (c *ImageResultClass) GetMethod(name string) (data.Method, bool) {
	return nil, false
}

func (c *ImageResultClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}

func (c *ImageResultClass) GetMethods() []data.Method { return nil }

func (c *ImageResultClass) GetConstruct() data.Method { return nil }
