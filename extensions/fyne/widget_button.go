package fyne

import (
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ButtonClass 是 Fyne\Widget\Button 类
type ButtonClass struct {
	*CanvasObjectClass
}

func NewButtonClass() data.ClassStmt {
	return &ButtonClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\Button", nil),
	}
}

func (c *ButtonClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *ButtonClass) GetConstruct() data.Method { return &buttonConstruct{} }

func (c *ButtonClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setText":
		return &buttonSetTextMethod{}, true
	case "getText":
		return &buttonGetTextMethod{}, true
	case "setImportance":
		return &buttonSetImportanceMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *ButtonClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&buttonSetTextMethod{},
		&buttonGetTextMethod{},
		&buttonSetImportanceMethod{},
	)
}

type buttonConstruct struct{}

func (m *buttonConstruct) GetName() string            { return token.ConstructName }
func (m *buttonConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *buttonConstruct) GetIsStatic() bool          { return false }
func (m *buttonConstruct) GetReturnType() data.Types  { return nil }

func (m *buttonConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "text", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "callback", 1, data.NewNullValue(), nil),
	}
}
func (m *buttonConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "text", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "callback", 1, nil),
	}
}

func (m *buttonConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	text := ""
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			text = s.AsString()
		}
	}

	// 获取回调闭包
	var callback data.FuncStmt
	if v, ok := ctx.GetIndexValue(1); ok {
		if fv, ok := v.(*data.FuncValue); ok {
			callback = fv.Value
		}
	}

	btn := widget.NewButton(text, func() {
		callPHPCallback(callback, ctx)
	})

	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			setFyneObject(classVal, btn)
			classVal.SetProperty("_button", data.NewAnyValue(btn))
		}
	}
	return nil, nil
}

func getButton(cv *data.ClassValue) *widget.Button {
	if v, _ := cv.GetProperty("_button"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if b, ok := av.Value.(*widget.Button); ok {
				return b
			}
		}
	}
	return nil
}

type buttonSetTextMethod struct{}

func (m *buttonSetTextMethod) GetName() string            { return "setText" }
func (m *buttonSetTextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *buttonSetTextMethod) GetIsStatic() bool          { return false }
func (m *buttonSetTextMethod) GetReturnType() data.Types  { return nil }
func (m *buttonSetTextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "text", 0, nil, data.NewBaseType("string")),
	}
}
func (m *buttonSetTextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "text", 0, data.NewBaseType("string")),
	}
}
func (m *buttonSetTextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if b := getButton(classVal); b != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						b.SetText(s.AsString())
					}
				}
			}
		}
	}
	return nil, nil
}

type buttonGetTextMethod struct{}

func (m *buttonGetTextMethod) GetName() string               { return "getText" }
func (m *buttonGetTextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *buttonGetTextMethod) GetIsStatic() bool             { return false }
func (m *buttonGetTextMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
func (m *buttonGetTextMethod) GetParams() []data.GetValue    { return nil }
func (m *buttonGetTextMethod) GetVariables() []data.Variable { return nil }
func (m *buttonGetTextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if b := getButton(classVal); b != nil {
				return data.NewStringValue(b.Text), nil
			}
		}
	}
	return data.NewStringValue(""), nil
}

type buttonSetImportanceMethod struct{}

func (m *buttonSetImportanceMethod) GetName() string            { return "setImportance" }
func (m *buttonSetImportanceMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *buttonSetImportanceMethod) GetIsStatic() bool          { return false }
func (m *buttonSetImportanceMethod) GetReturnType() data.Types  { return nil }
func (m *buttonSetImportanceMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "importance", 0, nil, data.NewBaseType("string")),
	}
}
func (m *buttonSetImportanceMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "importance", 0, data.NewBaseType("string")),
	}
}
func (m *buttonSetImportanceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if b := getButton(classVal); b != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						switch s.AsString() {
						case "low":
							b.Importance = widget.LowImportance
						case "high", "danger":
							b.Importance = widget.HighImportance
						default:
							b.Importance = widget.MediumImportance
						}
					}
				}
			}
		}
	}
	return nil, nil
}
