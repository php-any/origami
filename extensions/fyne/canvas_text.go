package fyne

import (
	"image/color"

	"fyne.io/fyne/v2/canvas"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// CanvasTextClass 是 Fyne\Canvas\Text 类
type CanvasTextClass struct {
	*CanvasObjectClass
}

func NewCanvasTextClass() data.ClassStmt {
	return &CanvasTextClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Canvas\\Text", nil),
	}
}

func (c *CanvasTextClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *CanvasTextClass) GetConstruct() data.Method { return &canvasTextConstruct{} }

func (c *CanvasTextClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setText":
		return &canvasTextSetTextMethod{}, true
	case "getText":
		return &canvasTextGetTextMethod{}, true
	case "setTextSize":
		return &canvasTextSetSizeMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *CanvasTextClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&canvasTextSetTextMethod{},
		&canvasTextGetTextMethod{},
		&canvasTextSetSizeMethod{},
	)
}

type canvasTextConstruct struct{}

func (m *canvasTextConstruct) GetName() string            { return token.ConstructName }
func (m *canvasTextConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasTextConstruct) GetIsStatic() bool          { return false }
func (m *canvasTextConstruct) GetReturnType() data.Types  { return nil }
func (m *canvasTextConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "text", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "color", 1, data.NewNullValue(), nil),
	}
}
func (m *canvasTextConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "text", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "color", 1, nil),
	}
}

func (m *canvasTextConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	text := ""
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			text = s.AsString()
		}
	}
	c := color.Black
	if v, ok := ctx.GetIndexValue(1); ok {
		c = fyneColorToGo(v)
	}
	t := canvas.NewText(text, c)
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			setFyneObject(classVal, t)
			classVal.SetProperty("_text", data.NewAnyValue(t))
		}
	}
	return nil, nil
}

func getCanvasText(cv *data.ClassValue) *canvas.Text {
	if v, _ := cv.GetProperty("_text"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if t, ok := av.Value.(*canvas.Text); ok {
				return t
			}
		}
	}
	return nil
}

type canvasTextSetTextMethod struct{}

func (m *canvasTextSetTextMethod) GetName() string            { return "setText" }
func (m *canvasTextSetTextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasTextSetTextMethod) GetIsStatic() bool          { return false }
func (m *canvasTextSetTextMethod) GetReturnType() data.Types  { return nil }
func (m *canvasTextSetTextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "text", 0, nil, data.NewBaseType("string")),
	}
}
func (m *canvasTextSetTextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "text", 0, data.NewBaseType("string")),
	}
}
func (m *canvasTextSetTextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if t := getCanvasText(classVal); t != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						t.Text = s.AsString()
					}
				}
			}
		}
	}
	return nil, nil
}

type canvasTextGetTextMethod struct{}

func (m *canvasTextGetTextMethod) GetName() string               { return "getText" }
func (m *canvasTextGetTextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *canvasTextGetTextMethod) GetIsStatic() bool             { return false }
func (m *canvasTextGetTextMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
func (m *canvasTextGetTextMethod) GetParams() []data.GetValue    { return nil }
func (m *canvasTextGetTextMethod) GetVariables() []data.Variable { return nil }
func (m *canvasTextGetTextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if t := getCanvasText(classVal); t != nil {
				return data.NewStringValue(t.Text), nil
			}
		}
	}
	return data.NewStringValue(""), nil
}

type canvasTextSetSizeMethod struct{}

func (m *canvasTextSetSizeMethod) GetName() string            { return "setTextSize" }
func (m *canvasTextSetSizeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasTextSetSizeMethod) GetIsStatic() bool          { return false }
func (m *canvasTextSetSizeMethod) GetReturnType() data.Types  { return nil }
func (m *canvasTextSetSizeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "size", 0, nil, data.NewBaseType("float")),
	}
}
func (m *canvasTextSetSizeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "size", 0, data.NewBaseType("float")),
	}
}
func (m *canvasTextSetSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if t := getCanvasText(classVal); t != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if f, ok := v.(data.AsFloat); ok {
						s, _ := f.AsFloat()
						t.TextSize = float32(s)
					}
				}
			}
		}
	}
	return nil, nil
}
