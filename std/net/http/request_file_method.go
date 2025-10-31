package http

import (
	httpsrc "net/http"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// RequestFileMethod 获取上传的文件
type RequestFileMethod struct {
	source *httpsrc.Request
}

func (h *RequestFileMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if h.source == nil {
		return data.NewAnyValue(nil), nil
	}

	// 检查是否有参数
	_, hasKey := ctx.GetIndexValue(0)

	// 如果没有参数，返回所有文件
	if !hasKey {
		if h.source.MultipartForm == nil {
			return data.NewObjectValue(), nil
		}
		result := data.NewObjectValue()
		for key, files := range h.source.MultipartForm.File {
			// 将文件信息转换为字符串
			fileInfo := ""
			for _, file := range files {
				fileInfo += file.Filename + ":" + string(rune(file.Size)) + ";"
			}
			result.SetProperty(key, data.NewStringValue(fileInfo))
		}
		return result, nil
	}

	// 如果有参数，返回指定键的文件
	param0, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("参数转换失败: %v", err)
	}

	if h.source.MultipartForm == nil {
		return data.NewStringValue(""), nil
	}

	if files, exists := h.source.MultipartForm.File[param0]; exists {
		// 将文件信息转换为字符串
		fileInfo := ""
		for _, file := range files {
			fileInfo += file.Filename + ":" + string(rune(file.Size)) + ";"
		}
		return data.NewStringValue(fileInfo), nil
	}

	return data.NewStringValue(""), nil
}

func (h *RequestFileMethod) GetName() string            { return "file" }
func (h *RequestFileMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *RequestFileMethod) GetIsStatic() bool          { return false }
func (h *RequestFileMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "key", 0, nil, nil),
	}
}
func (h *RequestFileMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "key", 0, nil),
	}
}
func (h *RequestFileMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
