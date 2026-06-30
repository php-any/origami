package fyne

import (
	"image/color"

	"fyne.io/fyne/v2/canvas"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// CanvasRectangleClass 是 Fyne\Canvas\Rectangle 类
type CanvasRectangleClass struct {
	*CanvasObjectClass
}

func NewCanvasRectangleClass() data.ClassStmt {
	return &CanvasRectangleClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Canvas\\Rectangle", nil),
	}
}

func (c *CanvasRectangleClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *CanvasRectangleClass) GetConstruct() data.Method { return &canvasRectConstruct{} }

func (c *CanvasRectangleClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setFillColor":
		return &canvasRectSetFillColorMethod{}, true
	case "setStrokeColor":
		return &canvasRectSetStrokeColorMethod{}, true
	case "setStrokeWidth":
		return &canvasRectSetStrokeWidthMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *CanvasRectangleClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&canvasRectSetFillColorMethod{},
		&canvasRectSetStrokeColorMethod{},
		&canvasRectSetStrokeWidthMethod{},
	)
}

type canvasRectConstruct struct{}

func (m *canvasRectConstruct) GetName() string            { return token.ConstructName }
func (m *canvasRectConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasRectConstruct) GetIsStatic() bool          { return false }
func (m *canvasRectConstruct) GetReturnType() data.Types  { return nil }
func (m *canvasRectConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "color", 0, data.NewNullValue(), nil),
	}
}
func (m *canvasRectConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "color", 1, nil),
	}
}

func (m *canvasRectConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	var c color.Color = color.Black
	if v, ok := ctx.GetIndexValue(0); ok {
		c = fyneColorToGo(v)
	}
	rect := canvas.NewRectangle(c)
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			setFyneObject(classVal, rect)
			classVal.SetProperty("_rect", data.NewAnyValue(rect))
		}
	}
	return nil, nil
}

func getCanvasRect(cv *data.ClassValue) *canvas.Rectangle {
	if v, _ := cv.GetProperty("_rect"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if r, ok := av.Value.(*canvas.Rectangle); ok {
				return r
			}
		}
	}
	return nil
}

type canvasRectSetFillColorMethod struct{}

func (m *canvasRectSetFillColorMethod) GetName() string            { return "setFillColor" }
func (m *canvasRectSetFillColorMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasRectSetFillColorMethod) GetIsStatic() bool          { return false }
func (m *canvasRectSetFillColorMethod) GetReturnType() data.Types  { return nil }
func (m *canvasRectSetFillColorMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "color", 0, nil, nil),
	}
}
func (m *canvasRectSetFillColorMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "color", 0, nil),
	}
}
func (m *canvasRectSetFillColorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if r := getCanvasRect(classVal); r != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					r.FillColor = fyneColorToGo(v)
				}
			}
		}
	}
	return nil, nil
}

type canvasRectSetStrokeColorMethod struct{}

func (m *canvasRectSetStrokeColorMethod) GetName() string            { return "setStrokeColor" }
func (m *canvasRectSetStrokeColorMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasRectSetStrokeColorMethod) GetIsStatic() bool          { return false }
func (m *canvasRectSetStrokeColorMethod) GetReturnType() data.Types  { return nil }
func (m *canvasRectSetStrokeColorMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "color", 0, nil, nil),
	}
}
func (m *canvasRectSetStrokeColorMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "color", 0, nil),
	}
}
func (m *canvasRectSetStrokeColorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if r := getCanvasRect(classVal); r != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					r.StrokeColor = fyneColorToGo(v)
				}
			}
		}
	}
	return nil, nil
}

type canvasRectSetStrokeWidthMethod struct{}

func (m *canvasRectSetStrokeWidthMethod) GetName() string            { return "setStrokeWidth" }
func (m *canvasRectSetStrokeWidthMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasRectSetStrokeWidthMethod) GetIsStatic() bool          { return false }
func (m *canvasRectSetStrokeWidthMethod) GetReturnType() data.Types  { return nil }
func (m *canvasRectSetStrokeWidthMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "width", 0, nil, data.NewBaseType("float")),
	}
}
func (m *canvasRectSetStrokeWidthMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "width", 0, data.NewBaseType("float")),
	}
}
func (m *canvasRectSetStrokeWidthMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if r := getCanvasRect(classVal); r != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if f, ok := v.(data.AsFloat); ok {
						w, _ := f.AsFloat()
						r.StrokeWidth = float32(w)
					}
				}
			}
		}
	}
	return nil, nil
}
