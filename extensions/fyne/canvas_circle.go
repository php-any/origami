package fyne

import (
	"image/color"

	"fyne.io/fyne/v2/canvas"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// CanvasCircleClass 是 Fyne\Canvas\Circle 类
type CanvasCircleClass struct {
	*CanvasObjectClass
}

func NewCanvasCircleClass() data.ClassStmt {
	return &CanvasCircleClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Canvas\\Circle", nil),
	}
}

func (c *CanvasCircleClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *CanvasCircleClass) GetConstruct() data.Method { return &canvasCircleConstruct{} }

func (c *CanvasCircleClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setFillColor":
		return &canvasCircleSetFillColorMethod{}, true
	case "setStrokeColor":
		return &canvasCircleSetStrokeColorMethod{}, true
	case "setStrokeWidth":
		return &canvasCircleSetStrokeWidthMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *CanvasCircleClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&canvasCircleSetFillColorMethod{},
		&canvasCircleSetStrokeColorMethod{},
		&canvasCircleSetStrokeWidthMethod{},
	)
}

type canvasCircleConstruct struct{}

func (m *canvasCircleConstruct) GetName() string            { return token.ConstructName }
func (m *canvasCircleConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasCircleConstruct) GetIsStatic() bool          { return false }
func (m *canvasCircleConstruct) GetReturnType() data.Types  { return nil }
func (m *canvasCircleConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "color", 0, data.NewNullValue(), nil),
	}
}
func (m *canvasCircleConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "color", 0, nil),
	}
}

func (m *canvasCircleConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	c := color.Black
	if v, ok := ctx.GetIndexValue(0); ok {
		c = fyneColorToGo(v)
	}
	circle := canvas.NewCircle(c)
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			setFyneObject(classVal, circle)
			classVal.SetProperty("_circle", data.NewAnyValue(circle))
		}
	}
	return nil, nil
}

func getCanvasCircle(cv *data.ClassValue) *canvas.Circle {
	if v, _ := cv.GetProperty("_circle"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if c, ok := av.Value.(*canvas.Circle); ok {
				return c
			}
		}
	}
	return nil
}

type canvasCircleSetFillColorMethod struct{}

func (m *canvasCircleSetFillColorMethod) GetName() string            { return "setFillColor" }
func (m *canvasCircleSetFillColorMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasCircleSetFillColorMethod) GetIsStatic() bool          { return false }
func (m *canvasCircleSetFillColorMethod) GetReturnType() data.Types  { return nil }
func (m *canvasCircleSetFillColorMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "color", 0, nil, nil),
	}
}
func (m *canvasCircleSetFillColorMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "color", 0, nil),
	}
}
func (m *canvasCircleSetFillColorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if c := getCanvasCircle(classVal); c != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					c.FillColor = fyneColorToGo(v)
				}
			}
		}
	}
	return nil, nil
}

type canvasCircleSetStrokeColorMethod struct{}

func (m *canvasCircleSetStrokeColorMethod) GetName() string            { return "setStrokeColor" }
func (m *canvasCircleSetStrokeColorMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasCircleSetStrokeColorMethod) GetIsStatic() bool          { return false }
func (m *canvasCircleSetStrokeColorMethod) GetReturnType() data.Types  { return nil }
func (m *canvasCircleSetStrokeColorMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "color", 0, nil, nil),
	}
}
func (m *canvasCircleSetStrokeColorMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "color", 0, nil),
	}
}
func (m *canvasCircleSetStrokeColorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if c := getCanvasCircle(classVal); c != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					c.StrokeColor = fyneColorToGo(v)
				}
			}
		}
	}
	return nil, nil
}

type canvasCircleSetStrokeWidthMethod struct{}

func (m *canvasCircleSetStrokeWidthMethod) GetName() string            { return "setStrokeWidth" }
func (m *canvasCircleSetStrokeWidthMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *canvasCircleSetStrokeWidthMethod) GetIsStatic() bool          { return false }
func (m *canvasCircleSetStrokeWidthMethod) GetReturnType() data.Types  { return nil }
func (m *canvasCircleSetStrokeWidthMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "width", 0, nil, data.NewBaseType("float")),
	}
}
func (m *canvasCircleSetStrokeWidthMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "width", 0, data.NewBaseType("float")),
	}
}
func (m *canvasCircleSetStrokeWidthMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if c := getCanvasCircle(classVal); c != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if f, ok := v.(data.AsFloat); ok {
						w, _ := f.AsFloat()
						c.StrokeWidth = float32(w)
					}
				}
			}
		}
	}
	return nil, nil
}
