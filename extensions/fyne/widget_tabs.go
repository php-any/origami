package fyne

import (
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// TabsClass 是 Fyne\Widget\Tabs 类
type TabsClass struct {
	*CanvasObjectClass
}

func NewTabsClass() data.ClassStmt {
	return &TabsClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\Tabs", nil),
	}
}

func (c *TabsClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *TabsClass) GetConstruct() data.Method { return &tabsConstruct{} }

func (c *TabsClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "append":
		return &tabsAppendMethod{}, true
	case "select":
		return &tabsSelectMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *TabsClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&tabsAppendMethod{},
		&tabsSelectMethod{},
	)
}

type tabsConstruct struct{}

func (m *tabsConstruct) GetName() string               { return token.ConstructName }
func (m *tabsConstruct) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *tabsConstruct) GetIsStatic() bool             { return false }
func (m *tabsConstruct) GetReturnType() data.Types     { return nil }
func (m *tabsConstruct) GetParams() []data.GetValue    { return nil }
func (m *tabsConstruct) GetVariables() []data.Variable { return nil }

func (m *tabsConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	tabs := widget.NewTabs()
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			setFyneObject(classVal, tabs)
			classVal.SetProperty("_tabs", data.NewAnyValue(tabs))
		}
	}
	return nil, nil
}

func getTabs(cv *data.ClassValue) *widget.Tabs {
	if v, _ := cv.GetProperty("_tabs"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if t, ok := av.Value.(*widget.Tabs); ok {
				return t
			}
		}
	}
	return nil
}

type tabsAppendMethod struct{}

func (m *tabsAppendMethod) GetName() string            { return "append" }
func (m *tabsAppendMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *tabsAppendMethod) GetIsStatic() bool          { return false }
func (m *tabsAppendMethod) GetReturnType() data.Types  { return nil }
func (m *tabsAppendMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "text", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "content", 1, nil, data.NewBaseType("Fyne\\CanvasObject")),
	}
}
func (m *tabsAppendMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "text", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "content", 1, data.NewBaseType("Fyne\\CanvasObject")),
	}
}

func (m *tabsAppendMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if t := getTabs(classVal); t != nil {
				text := ""
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						text = s.AsString()
					}
				}
				if v, ok := ctx.GetIndexValue(1); ok {
					if contentCV, ok := v.(*data.ClassValue); ok {
						if obj := getFyneObject(contentCV); obj != nil {
							item := widget.NewTabItem(text, obj)
							t.Append(item)
						}
					}
				}
			}
		}
	}
	return nil, nil
}

type tabsSelectMethod struct{}

func (m *tabsSelectMethod) GetName() string            { return "select" }
func (m *tabsSelectMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *tabsSelectMethod) GetIsStatic() bool          { return false }
func (m *tabsSelectMethod) GetReturnType() data.Types  { return nil }
func (m *tabsSelectMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "index", 0, nil, data.NewBaseType("int")),
	}
}
func (m *tabsSelectMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "index", 0, data.NewBaseType("int")),
	}
}
func (m *tabsSelectMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if t := getTabs(classVal); t != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if i, ok := v.(data.AsInt); ok {
						idx, _ := i.AsInt()
						t.SelectIndex(idx)
					}
				}
			}
		}
	}
	return nil, nil
}
