package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ResponseWriterViewMethod 支持渲染 HTML 模板并可传入参数
type ResponseWriterViewMethod struct {
	w *bufferedWriter
}

func (h *ResponseWriterViewMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	pathValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrowf("view 方法缺少模板路径参数: %v", 0)
	}
	templatePath := pathValue.AsString()
	if templatePath == "" {
		return nil, utils.NewThrowf("view 方法模板路径为空")
	}

	objectValue, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrowf("view 方法缺少objectValue参数: %v", 1)
	}

	rendered, acl := ctx.GetVM().ParseFile(templatePath, objectValue)
	if acl != nil {
		return nil, acl
	}

	h.w.SetHeader("Content-Type", "text/html; charset=utf-8")
	if rendered != nil {
		if val, ok := rendered.(data.Value); ok {
			if _, err := h.w.Write([]byte(val.AsString())); err != nil {
				return nil, utils.NewThrow(err)
			}
		}
	}
	return nil, nil
}

func (h *ResponseWriterViewMethod) GetName() string            { return "view" }
func (h *ResponseWriterViewMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterViewMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterViewMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "templatePath", 0, nil, nil),
		node.NewParameter(nil, "data", 1, nil, data.Object{}),
	}
}
func (h *ResponseWriterViewMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "templatePath", 0, nil),
		node.NewVariable(nil, "data", 1, nil),
	}
}
func (h *ResponseWriterViewMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
