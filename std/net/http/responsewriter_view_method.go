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

	var layoutPath string
	if layoutValue, hasLayout := ctx.GetIndexValue(2); hasLayout {
		if layoutValue != nil {
			layoutPath = layoutValue.AsString()
		}
	}

	output := rendered
	if layoutPath != "" {
		content, ok := rendered.(data.Value)
		if !ok {
			return nil, utils.NewThrowf("view 渲染结果无法作为 layout content")
		}
		layoutData := viewDataWithContent(objectValue, content)
		layoutValue, ok := layoutData.(data.Value)
		if !ok {
			return nil, utils.NewThrowf("view layout 数据无法转换为 Value")
		}
		output, acl = ctx.GetVM().ParseFile(layoutPath, layoutValue)
		if acl != nil {
			return nil, acl
		}
	}

	h.w.SetHeader("Content-Type", "text/html; charset=utf-8")
	if output != nil {
		if val, ok := output.(data.Value); ok {
			if _, err := h.w.Write([]byte(val.AsString())); err != nil {
				return nil, utils.NewThrow(err)
			}
		}
	}
	return nil, nil
}

func viewDataWithContent(objectValue data.GetValue, content data.Value) data.GetValue {
	props := viewDataProps(objectValue)
	props["content"] = content

	items := make([]*data.ZVal, 0, len(props))
	for name, value := range props {
		items = append(items, data.NewNamedZVal(name, value))
	}
	return &data.ArrayValue{List: items}
}

func viewDataProps(objectValue data.GetValue) map[string]data.Value {
	switch v := objectValue.(type) {
	case *data.ObjectValue:
		return v.GetProperties()
	case *data.ClassValue:
		return v.GetProperties()
	case *data.ArrayValue:
		props := make(map[string]data.Value)
		for _, z := range v.List {
			if z == nil || z.Name == "" {
				continue
			}
			props[z.Name] = z.Value
		}
		return props
	default:
		return map[string]data.Value{}
	}
}

func (h *ResponseWriterViewMethod) GetName() string            { return "view" }
func (h *ResponseWriterViewMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *ResponseWriterViewMethod) GetIsStatic() bool          { return false }
func (h *ResponseWriterViewMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "templatePath", 0, nil, nil),
		node.NewParameter(nil, "data", 1, nil, data.Object{}),
		node.NewParameter(nil, "layoutPath", 2, data.NewNullValue(), data.NewBaseType("string")),
	}
}
func (h *ResponseWriterViewMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "templatePath", 0, nil),
		node.NewVariable(nil, "data", 1, nil),
		node.NewVariable(nil, "layoutPath", 2, nil),
	}
}
func (h *ResponseWriterViewMethod) GetReturnType() data.Types { return data.NewBaseType("void") }
