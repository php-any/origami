package fyne

import (
	"image/color"

	"fyne.io/fyne/v2/canvas"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// CanvasLineClass 是 Fyne\Canvas\Line 类
type CanvasLineClass struct {
	*CanvasObjectClass
}

func NewCanvasLineClass() data.ClassStmt {
	return &CanvasLineClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Canvas\\Line", nil),
	}
}

func (c *CanvasLineClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *CanvasLineClass) GetConstruct() data.Method { return &canvasLineConstruct{} }

func (c *CanvasLineClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setStrokeColor":
		return &canvasLineSetStrokeColorMethod{}, true
	case "setStrokeWidth":
		return &canvasLineSetStrokeWidthMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *CanvasLineClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&canvasLineSetStrokeColorMethod{},
		&canvasLineSetStrokeWidthMethod{},
	)
}

type canvasLineConstruct struct{}

func (m *canvasLineConstruct) GetName() string            { return token.ConstructName }
func (m *canvasLineConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasLineConstruct) GetIsStatic() bool          { return false }
func (m *canvasLineConstruct) GetReturnType() data.Types  { return nil }
func (m *canvasLineConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "color", 0, data.NewNullValue(), nil),
	}
}
func (m *canvasLineConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "color", 0, nil),
	}
}

func (m *canvasLineConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	c := color.Black
	if v, ok := ctx.GetIndexValue(0); ok {
		c = fyneColorToGo(v)
	}
	line := canvas.NewLine(c)
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			setFyneObject(classVal, line)
			classVal.SetProperty("_line", data.NewAnyValue(line))
		}
	}
	return nil, nil
}

func getCanvasLine(cv *data.ClassValue) *canvas.Line {
	if v, _ := cv.GetProperty("_line"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if l, ok := av.Value.(*canvas.Line); ok {
				return l
			}
		}
	}
	return nil
}

type canvasLineSetStrokeColorMethod struct{}

func (m *canvasLineSetStrokeColorMethod) GetName() string            { return "setStrokeColor" }
func (m *canvasLineSetStrokeColorMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasLineSetStrokeColorMethod) GetIsStatic() bool          { return false }
func (m *canvasLineSetStrokeColorMethod) GetReturnType() data.Types  { return nil }
func (m *canvasLineSetStrokeColorMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "color", 0, nil, nil),
	}
}
func (m *canvasLineSetStrokeColorMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "color", 0, nil),
	}
}
func (m *canvasLineSetStrokeColorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if l := getCanvasLine(classVal); l != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					l.StrokeColor = fyneColorToGo(v)
				}
			}
		}
	}
	return nil, nil
}

type canvasLineSetStrokeWidthMethod struct{}

func (m *canvasLineSetStrokeWidthMethod) GetName() string            { return "setStrokeWidth" }
func (m *canvasLineSetStrokeWidthMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasLineSetStrokeWidthMethod) GetIsStatic() bool          { return false }
func (m *canvasLineSetStrokeWidthMethod) GetReturnType() data.Types  { return nil }
func (m *canvasLineSetStrokeWidthMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "width", 0, nil, data.NewBaseType("float")),
	}
}
func (m *canvasLineSetStrokeWidthMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "width", 0, data.NewBaseType("float")),
	}
}
func (m *canvasLineSetStrokeWidthMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if l := getCanvasLine(classVal); l != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if f, ok := v.(data.AsFloat); ok {
						w, _ := f.AsFloat()
						l.StrokeWidth = float32(w)
					}
				}
			}
		}
	}
	return nil, nil
}
