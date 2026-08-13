package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

type ResponseWriterFileMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterFileMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	path, err := utils.ConvertFromIndex[string](ctx, 0)
	if err != nil {
		return nil, utils.NewThrowf("file 方法缺少路径参数: %v", err)
	}

	downloadName := ""
	if _, ok := ctx.GetIndexValue(1); ok {
		downloadName, err = utils.ConvertFromIndex[string](ctx, 1)
		if err != nil {
			return nil, utils.NewThrowf("downloadName 参数转换失败: %v", err)
		}
	}

	if err := h.w.SendFile(path, downloadName); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}

func (h *ResponseWriterFileMethod) GetName() string            { return "file" }
func (h *ResponseWriterFileMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterFileMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterFileMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "path", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "downloadName", 1, nil, data.NewBaseType("string")),
	}
}
func (h *ResponseWriterFileMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "path", 0, nil),
		node.NewVariable(nil, "downloadName", 1, nil),
	}
}
func (h *ResponseWriterFileMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
