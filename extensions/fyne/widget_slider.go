package fyne

import (
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// SliderClass 是 Fyne\Widget\Slider 类
type SliderClass struct {
	*CanvasObjectClass
}

func NewSliderClass() data.ClassStmt {
	return &SliderClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\Slider", nil),
	}
}

func (c *SliderClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *SliderClass) GetConstruct() data.Method { return &sliderConstruct{} }

func (c *SliderClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setValue":
		return &sliderSetValueMethod{}, true
	case "getValue":
		return &sliderGetValueMethod{}, true
	case "onChanged":
		return &sliderOnChangedMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *SliderClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&sliderSetValueMethod{},
		&sliderGetValueMethod{},
		&sliderOnChangedMethod{},
	)
}

type sliderConstruct struct{}

func (m *sliderConstruct) GetName() string            { return token.ConstructName }
func (m *sliderConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *sliderConstruct) GetIsStatic() bool          { return false }
func (m *sliderConstruct) GetReturnType() data.Types  { return nil }
func (m *sliderConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "min", 0, data.NewFloatValue(0), data.NewBaseType("float")),
		node.NewParameter(nil, "max", 1, data.NewFloatValue(1), data.NewBaseType("float")),
	}
}
func (m *sliderConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "min", 0, data.NewBaseType("float")),
		node.NewVariable(nil, "max", 1, data.NewBaseType("float")),
	}
}

func (m *sliderConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	min, max := 0.0, 1.0
	if v, ok := ctx.GetIndexValue(0); ok {
		if f, ok := v.(data.AsFloat); ok {
			min, _ = f.AsFloat()
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		if f, ok := v.(data.AsFloat); ok {
			max, _ = f.AsFloat()
		}
	}
	slider := widget.NewSlider(min, max)
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			setFyneObject(classVal, slider)
			classVal.SetProperty("_slider", data.NewAnyValue(slider))
		}
	}
	return nil, nil
}

func getSlider(cv *data.ClassValue) *widget.Slider {
	if v, _ := cv.GetProperty("_slider"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if s, ok := av.Value.(*widget.Slider); ok {
				return s
			}
		}
	}
	return nil
}

type sliderSetValueMethod struct{}

func (m *sliderSetValueMethod) GetName() string            { return "setValue" }
func (m *sliderSetValueMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *sliderSetValueMethod) GetIsStatic() bool          { return false }
func (m *sliderSetValueMethod) GetReturnType() data.Types  { return nil }
func (m *sliderSetValueMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, data.NewBaseType("float")),
	}
}
func (m *sliderSetValueMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("float")),
	}
}
func (m *sliderSetValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if s := getSlider(classVal); s != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if f, ok := v.(data.AsFloat); ok {
						val, _ := f.AsFloat()
						s.SetValue(val)
					}
				}
			}
		}
	}
	return nil, nil
}

type sliderGetValueMethod struct{}

func (m *sliderGetValueMethod) GetName() string               { return "getValue" }
func (m *sliderGetValueMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *sliderGetValueMethod) GetIsStatic() bool             { return false }
func (m *sliderGetValueMethod) GetReturnType() data.Types     { return data.NewBaseType("float") }
func (m *sliderGetValueMethod) GetParams() []data.GetValue    { return nil }
func (m *sliderGetValueMethod) GetVariables() []data.Variable { return nil }
func (m *sliderGetValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if s := getSlider(classVal); s != nil {
				return data.NewFloatValue(s.Value), nil
			}
		}
	}
	return data.NewFloatValue(0), nil
}

type sliderOnChangedMethod struct{}

func (m *sliderOnChangedMethod) GetName() string            { return "onChanged" }
func (m *sliderOnChangedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *sliderOnChangedMethod) GetIsStatic() bool          { return false }
func (m *sliderOnChangedMethod) GetReturnType() data.Types  { return nil }
func (m *sliderOnChangedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
	}
}
func (m *sliderOnChangedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, nil),
	}
}
func (m *sliderOnChangedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if s := getSlider(classVal); s != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					var callback data.FuncStmt
					if fv, ok := v.(*data.FuncValue); ok {
						callback = fv.Value
					}
					s.OnChanged = func(value float64) {
						callPHPCallbackWith(callback, ctx, data.NewFloatValue(value))
					}
				}
			}
		}
	}
	return nil, nil
}
