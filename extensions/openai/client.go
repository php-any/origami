package openai

import (
	"context"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

// client 是 OpenAI 客户端的 Go 原生封装
type client struct {
	inner *openai.Client
}

// newClient 创建一个新的 OpenAI 客户端
func newClient(apiKey string, baseURL string) (*client, error) {
	opts := []option.RequestOption{
		option.WithAPIKey(apiKey),
	}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	c := openai.NewClient(opts...)
	return &client{inner: &c}, nil
}

// newClientFromEnv 从环境变量创建客户端
func newClientFromEnv() (*client, error) {
	c := openai.NewClient()
	return &client{inner: &c}, nil
}

// chat 发送 Chat Completion 请求
func (c *client) chat(model string, messages []openai.ChatCompletionMessageParamUnion, opts map[string]any) (*openai.ChatCompletion, error) {
	params := openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(model),
		Messages: messages,
	}
	applyChatOptions(&params, opts)
	return c.inner.Chat.Completions.New(context.Background(), params)
}

// embeddings 创建文本嵌入
func (c *client) embeddings(model string, input openai.EmbeddingNewParamsInputUnion, opts map[string]any) (*openai.CreateEmbeddingResponse, error) {
	params := openai.EmbeddingNewParams{
		Model: openai.EmbeddingModel(model),
		Input: input,
	}
	applyEmbeddingOptions(&params, opts)
	return c.inner.Embeddings.New(context.Background(), params)
}

// images 生成图片
func (c *client) images(model string, prompt string, opts map[string]any) (*openai.ImagesResponse, error) {
	params := openai.ImageGenerateParams{
		Prompt: prompt,
		Model:  model,
	}
	applyImageOptions(&params, opts)
	return c.inner.Images.Generate(context.Background(), params)
}

// speech 文本转语音并写入文件
func (c *client) speech(model string, input string, voice string, outputPath string, opts map[string]any) error {
	params := openai.AudioSpeechNewParams{
		Input: input,
		Model: openai.SpeechModel(model),
		Voice: openai.AudioSpeechNewParamsVoiceUnion{OfString: openai.String(voice)},
	}
	applySpeechOptions(&params, opts)
	resp, err := c.inner.Audio.Speech.New(context.Background(), params)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data := make([]byte, 0)
	buf := make([]byte, 32*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			data = append(data, buf[:n]...)
		}
		if readErr != nil {
			break
		}
	}
	return os.WriteFile(outputPath, data, 0644)
}

// transcription 语音转文字
func (c *client) transcription(model string, filePath string, opts map[string]any) (*openai.AudioTranscriptionNewResponseUnion, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	params := openai.AudioTranscriptionNewParams{
		Model: openai.AudioModel(model),
		File:  file,
	}
	applyTranscriptionOptions(&params, opts)
	return c.inner.Audio.Transcriptions.New(context.Background(), params)
}

// applyChatOptions 将 options map 应用到 ChatCompletionNewParams
func applyChatOptions(params *openai.ChatCompletionNewParams, opts map[string]any) {
	if opts == nil {
		return
	}
	if v, ok := opts["temperature"]; ok {
		if f, ok := toFloat64(v); ok {
			params.Temperature = openai.Float(f)
		}
	}
	if v, ok := opts["maxTokens"]; ok {
		if n, ok := toInt64(v); ok {
			params.MaxTokens = openai.Int(n)
		}
	}
	if v, ok := opts["topP"]; ok {
		if f, ok := toFloat64(v); ok {
			params.TopP = openai.Float(f)
		}
	}
	if v, ok := opts["stop"]; ok {
		if arr, ok := toSliceOfString(v); ok {
			params.Stop = openai.ChatCompletionNewParamsStopUnion{OfStringArray: arr}
		}
	}
	if v, ok := opts["seed"]; ok {
		if n, ok := toInt64(v); ok {
			params.Seed = openai.Int(n)
		}
	}
	if v, ok := opts["frequencyPenalty"]; ok {
		if f, ok := toFloat64(v); ok {
			params.FrequencyPenalty = openai.Float(f)
		}
	}
	if v, ok := opts["presencePenalty"]; ok {
		if f, ok := toFloat64(v); ok {
			params.PresencePenalty = openai.Float(f)
		}
	}
	if v, ok := opts["maxCompletionTokens"]; ok {
		if n, ok := toInt64(v); ok {
			params.MaxCompletionTokens = openai.Int(n)
		}
	}
	if v, ok := opts["n"]; ok {
		if n, ok := toInt64(v); ok {
			params.N = openai.Int(n)
		}
	}
	if v, ok := opts["responseFormat"]; ok {
		applyResponseFormat(params, v)
	}
}

// applyResponseFormat 设置 ResponseFormat
// 参考 DeepSeek docs: https://api-docs.deepseek.com/zh-cn/guides/json_mode
// 支持三种形式:
//   - "json_object" 字符串简写
//   - ["type" => "json_object"] 对象形式 (DeepSeek/OpenAI 标准写法)
//   - ["type" => "json_schema", "name" => ..., "schema" => ...] JSON Schema 结构化输出
func applyResponseFormat(params *openai.ChatCompletionNewParams, v any) {
	switch val := v.(type) {
	case string:
		if val == "json_object" {
			obj := shared.NewResponseFormatJSONObjectParam()
			params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONObject: &obj,
			}
		}
	case map[string]any:
		t, _ := val["type"].(string)
		switch t {
		case "json_object":
			obj := shared.NewResponseFormatJSONObjectParam()
			params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONObject: &obj,
			}
		case "json_schema":
			name := "schema"
			if n, ok := val["name"].(string); ok {
				name = n
			}
			schema := shared.ResponseFormatJSONSchemaJSONSchemaParam{
				Name: name,
			}
			if desc, ok := val["description"].(string); ok {
				schema.Description = openai.String(desc)
			}
			if strict, ok := val["strict"].(bool); ok {
				schema.Strict = openai.Bool(strict)
			}
			if s, ok := val["schema"].(map[string]any); ok {
				schema.Schema = s
			}
			params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
				OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
					JSONSchema: schema,
				},
			}
		}
	}
}

func applyEmbeddingOptions(params *openai.EmbeddingNewParams, opts map[string]any) {
	if opts == nil {
		return
	}
	if v, ok := opts["dimensions"]; ok {
		if n, ok := toInt64(v); ok {
			params.Dimensions = openai.Int(n)
		}
	}
}

func applyImageOptions(params *openai.ImageGenerateParams, opts map[string]any) {
	if opts == nil {
		return
	}
	if v, ok := opts["size"]; ok {
		if s, ok := v.(string); ok {
			params.Size = openai.ImageGenerateParamsSize(s)
		}
	}
	if v, ok := opts["quality"]; ok {
		if s, ok := v.(string); ok {
			params.Quality = openai.ImageGenerateParamsQuality(s)
		}
	}
	if v, ok := opts["n"]; ok {
		if n, ok := toInt64(v); ok {
			params.N = openai.Int(n)
		}
	}
}

func applySpeechOptions(params *openai.AudioSpeechNewParams, opts map[string]any) {
	if opts == nil {
		return
	}
	if v, ok := opts["speed"]; ok {
		if f, ok := toFloat64(v); ok {
			params.Speed = openai.Float(f)
		}
	}
}

func applyTranscriptionOptions(params *openai.AudioTranscriptionNewParams, opts map[string]any) {
	if opts == nil {
		return
	}
	if v, ok := opts["language"]; ok {
		if s, ok := v.(string); ok {
			params.Language = openai.String(s)
		}
	}
}

// 辅助类型转换函数
func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	default:
		return 0, false
	}
}

func toInt64(v any) (int64, bool) {
	switch val := v.(type) {
	case int:
		return int64(val), true
	case int64:
		return val, true
	case float64:
		return int64(val), true
	default:
		return 0, false
	}
}

func toSliceOfString(v any) ([]string, bool) {
	switch val := v.(type) {
	case []string:
		return val, true
	case []any:
		result := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result, true
	default:
		return nil, false
	}
}
