package fyne

import (
	fyneLib "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// LabelClass 是 Fyne\Widget\Label 类
type LabelClass struct {
	*CanvasObjectClass
}

func NewLabelClass() data.ClassStmt {
	return &LabelClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\Label", nil),
	}
}

func (c *LabelClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *LabelClass) GetConstruct() data.Method { return &labelConstruct{} }

func (c *LabelClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setText":
		return &labelSetTextMethod{}, true
	case "getText":
		return &labelGetTextMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *LabelClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&labelSetTextMethod{},
		&labelGetTextMethod{},
	)
}

type labelConstruct struct{}

func (m *labelConstruct) GetName() string            { return token.ConstructName }
func (m *labelConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *labelConstruct) GetIsStatic() bool          { return false }
func (m *labelConstruct) GetReturnType() data.Types  { return nil }

func (m *labelConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "text", 0, data.NewStringValue(""), data.NewBaseType("string")),
	}
}
func (m *labelConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "text", 0, data.NewBaseType("string")),
	}
}

func (m *labelConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	text := ""
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			text = s.AsString()
		}
	}
	label := widget.NewLabel(text)
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			setFyneObject(classVal, label)
			classVal.SetProperty("_label", data.NewAnyValue(label))
		}
	}
	return nil, nil
}

func getLabel(cv *data.ClassValue) *widget.Label {
	if v, _ := cv.GetProperty("_label"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if l, ok := av.Value.(*widget.Label); ok {
				return l
			}
		}
	}
	return nil
}

type labelSetTextMethod struct{}

func (m *labelSetTextMethod) GetName() string            { return "setText" }
func (m *labelSetTextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *labelSetTextMethod) GetIsStatic() bool          { return false }
func (m *labelSetTextMethod) GetReturnType() data.Types  { return nil }
func (m *labelSetTextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "text", 0, nil, data.NewBaseType("string")),
	}
}
func (m *labelSetTextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "text", 0, data.NewBaseType("string")),
	}
}
func (m *labelSetTextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if l := getLabel(classVal); l != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						l.SetText(s.AsString())
					}
				}
			}
		}
	}
	return nil, nil
}

type labelGetTextMethod struct{}

func (m *labelGetTextMethod) GetName() string               { return "getText" }
func (m *labelGetTextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *labelGetTextMethod) GetIsStatic() bool             { return false }
func (m *labelGetTextMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
func (m *labelGetTextMethod) GetParams() []data.GetValue    { return nil }
func (m *labelGetTextMethod) GetVariables() []data.Variable { return nil }
func (m *labelGetTextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if l := getLabel(classVal); l != nil {
				return data.NewStringValue(l.Text), nil
			}
		}
	}
	return data.NewStringValue(""), nil
}
