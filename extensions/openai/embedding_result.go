package openai

import (
	oai "github.com/openai/openai-go/v3"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewEmbeddingResultClass 创建 OpenAI\EmbeddingResult 类定义
func NewEmbeddingResultClass() data.ClassStmt {
	return &EmbeddingResultClass{}
}

// NewEmbeddingResultClassValue 从 OpenAI 响应创建 EmbeddingResult 类实例
func NewEmbeddingResultClassValue(resp *oai.CreateEmbeddingResponse, ctx data.Context) *data.ClassValue {
	class := &EmbeddingResultClass{resp: resp}
	cv := data.NewClassValue(class, ctx)

	if len(resp.Data) > 0 {
		embedding := resp.Data[0].Embedding
		// 将 []float64 转换为 data.ArrayValue
		vals := make([]data.Value, len(embedding))
		for i, f := range embedding {
			vals[i] = data.NewFloatValue(f)
		}
		cv.SetProperty("embedding", data.NewArrayValue(vals))
		cv.SetProperty("dimensions", data.NewIntValue(len(embedding)))
	} else {
		cv.SetProperty("embedding", data.NewArrayValue([]data.Value{}))
		cv.SetProperty("dimensions", data.NewIntValue(0))
	}

	// usage
	usageObj := data.NewObjectValue()
	usageObj.SetProperty("promptTokens", data.NewIntValue(int(resp.Usage.PromptTokens)))
	usageObj.SetProperty("totalTokens", data.NewIntValue(int(resp.Usage.TotalTokens)))
	cv.SetProperty("usage", usageObj)

	return cv
}

// EmbeddingResultClass 是 OpenAI\EmbeddingResult 的类定义
type EmbeddingResultClass struct {
	resp *oai.CreateEmbeddingResponse
}

func (c *EmbeddingResultClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if c.resp != nil {
		return NewEmbeddingResultClassValue(c.resp, ctx), nil
	}
	return data.NewClassValue(c, ctx), nil
}

func (c *EmbeddingResultClass) GetFrom() data.From      { return nil }
func (c *EmbeddingResultClass) GetName() string         { return "OpenAI\\EmbeddingResult" }
func (c *EmbeddingResultClass) GetExtend() *string      { return nil }
func (c *EmbeddingResultClass) GetImplements() []string { return nil }

func (c *EmbeddingResultClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (c *EmbeddingResultClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "embedding", "public", false, nil, data.NewBaseType("array")),
		node.NewProperty(nil, "dimensions", "public", false, nil, data.NewBaseType("int")),
		node.NewProperty(nil, "usage", "public", false, nil, data.NewBaseType("object")),
	}
}

func (c *EmbeddingResultClass) GetMethod(name string) (data.Method, bool) {
	return nil, false
}

func (c *EmbeddingResultClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}

func (c *EmbeddingResultClass) GetMethods() []data.Method { return nil }

func (c *EmbeddingResultClass) GetConstruct() data.Method { return nil }
