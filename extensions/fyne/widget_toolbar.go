package fyne

import (
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ToolbarClass 是 Fyne\Widget\Toolbar 类
type ToolbarClass struct {
	*CanvasObjectClass
}

func NewToolbarClass() data.ClassStmt {
	return &ToolbarClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\Toolbar", nil),
	}
}

func (c *ToolbarClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *ToolbarClass) GetConstruct() data.Method { return &toolbarConstruct{} }

func (c *ToolbarClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "append":
		return &toolbarAppendMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *ToolbarClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&toolbarAppendMethod{},
	)
}

type toolbarConstruct struct{}

func (m *toolbarConstruct) GetName() string               { return token.ConstructName }
func (m *toolbarConstruct) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *toolbarConstruct) GetIsStatic() bool             { return false }
func (m *toolbarConstruct) GetReturnType() data.Types     { return nil }
func (m *toolbarConstruct) GetParams() []data.GetValue    { return nil }
func (m *toolbarConstruct) GetVariables() []data.Variable { return nil }

func (m *toolbarConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	toolbar := widget.NewToolbar()
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			setFyneObject(classVal, toolbar)
			classVal.SetProperty("_toolbar", data.NewAnyValue(toolbar))
		}
	}
	return nil, nil
}

func getToolbar(cv *data.ClassValue) *widget.Toolbar {
	if v, _ := cv.GetProperty("_toolbar"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if t, ok := av.Value.(*widget.Toolbar); ok {
				return t
			}
		}
	}
	return nil
}

type toolbarAppendMethod struct{}

func (m *toolbarAppendMethod) GetName() string            { return "append" }
func (m *toolbarAppendMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *toolbarAppendMethod) GetIsStatic() bool          { return false }
func (m *toolbarAppendMethod) GetReturnType() data.Types  { return nil }
func (m *toolbarAppendMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "item", 0, nil, data.NewBaseType("Fyne\\Widget\\ToolbarItem")),
	}
}
func (m *toolbarAppendMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "item", 0, data.NewBaseType("Fyne\\Widget\\ToolbarItem")),
	}
}

func (m *toolbarAppendMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if t := getToolbar(classVal); t != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if itemCV, ok := v.(*data.ClassValue); ok {
						if item, _ := itemCV.GetProperty("_toolbarItem"); item != nil {
							if av, ok := item.(*data.AnyValue); ok {
								if ti, ok := av.Value.(widget.ToolbarItem); ok {
									t.Append(ti)
								}
							}
						}
					}
				}
			}
		}
	}
	return nil, nil
}

// ToolbarActionClass 是 Fyne\Widget\ToolbarAction 类
type ToolbarActionClass struct{}

func NewToolbarActionClass() data.ClassStmt { return &ToolbarActionClass{} }

func (c *ToolbarActionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *ToolbarActionClass) GetFrom() data.From                              { return nil }
func (c *ToolbarActionClass) GetName() string                                 { return "Fyne\\Widget\\ToolbarAction" }
func (c *ToolbarActionClass) GetExtend() *string                              { return nil }
func (c *ToolbarActionClass) GetImplements() []string                         { return nil }
func (c *ToolbarActionClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *ToolbarActionClass) GetPropertyList() []data.Property                { return nil }
func (c *ToolbarActionClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *ToolbarActionClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *ToolbarActionClass) GetMethods() []data.Method                       { return nil }
func (c *ToolbarActionClass) GetConstruct() data.Method                       { return &toolbarActionConstruct{} }

type toolbarActionConstruct struct{}

func (m *toolbarActionConstruct) GetName() string            { return token.ConstructName }
func (m *toolbarActionConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *toolbarActionConstruct) GetIsStatic() bool          { return false }
func (m *toolbarActionConstruct) GetReturnType() data.Types  { return nil }
func (m *toolbarActionConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "icon", 0, nil, nil),
		node.NewParameter(nil, "callback", 1, nil, nil),
	}
}
func (m *toolbarActionConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "icon", 0, nil),
		node.NewVariable(nil, "callback", 1, nil),
	}
}

func (m *toolbarActionConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	// icon 暂时忽略，使用 nil
	var callback data.FuncStmt
	if v, ok := ctx.GetIndexValue(1); ok {
		if fv, ok := v.(*data.FuncValue); ok {
			callback = fv.Value
		}
	}
	action := widget.NewToolbarAction(nil, func() {
		callPHPCallback(callback, ctx)
	})
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			classVal.SetProperty("_toolbarItem", data.NewAnyValue(action))
		}
	}
	return nil, nil
}

// ToolbarSeparatorClass 是 Fyne\Widget\ToolbarSeparator 类
type ToolbarSeparatorClass struct{}

func NewToolbarSeparatorClass() data.ClassStmt { return &ToolbarSeparatorClass{} }

func (c *ToolbarSeparatorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx)
	cv.SetProperty("_toolbarItem", data.NewAnyValue(widget.NewToolbarSeparator()))
	return cv, nil
}
func (c *ToolbarSeparatorClass) GetFrom() data.From                              { return nil }
func (c *ToolbarSeparatorClass) GetName() string                                 { return "Fyne\\Widget\\ToolbarSeparator" }
func (c *ToolbarSeparatorClass) GetExtend() *string                              { return nil }
func (c *ToolbarSeparatorClass) GetImplements() []string                         { return nil }
func (c *ToolbarSeparatorClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *ToolbarSeparatorClass) GetPropertyList() []data.Property                { return nil }
func (c *ToolbarSeparatorClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *ToolbarSeparatorClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *ToolbarSeparatorClass) GetMethods() []data.Method                       { return nil }
func (c *ToolbarSeparatorClass) GetConstruct() data.Method                       { return nil }

// ToolbarSpacerClass 是 Fyne\Widget\ToolbarSpacer 类
type ToolbarSpacerClass struct{}

func NewToolbarSpacerClass() data.ClassStmt { return &ToolbarSpacerClass{} }

func (c *ToolbarSpacerClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx)
	cv.SetProperty("_toolbarItem", data.NewAnyValue(widget.NewToolbarSpacer()))
	return cv, nil
}
func (c *ToolbarSpacerClass) GetFrom() data.From                              { return nil }
func (c *ToolbarSpacerClass) GetName() string                                 { return "Fyne\\Widget\\ToolbarSpacer" }
func (c *ToolbarSpacerClass) GetExtend() *string                              { return nil }
func (c *ToolbarSpacerClass) GetImplements() []string                         { return nil }
func (c *ToolbarSpacerClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *ToolbarSpacerClass) GetPropertyList() []data.Property                { return nil }
func (c *ToolbarSpacerClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *ToolbarSpacerClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *ToolbarSpacerClass) GetMethods() []data.Method                       { return nil }
func (c *ToolbarSpacerClass) GetConstruct() data.Method                       { return nil }
