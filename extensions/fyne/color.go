package fyne

import (
	"image/color"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ColorClass 是 Fyne\Color 值对象类
type ColorClass struct{}

func NewColorClass() data.ClassStmt { return &ColorClass{} }

func (c *ColorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *ColorClass) GetFrom() data.From                              { return nil }
func (c *ColorClass) GetName() string                                 { return "Fyne\\Color" }
func (c *ColorClass) GetExtend() *string                              { return nil }
func (c *ColorClass) GetImplements() []string                         { return nil }
func (c *ColorClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *ColorClass) GetPropertyList() []data.Property                { return nil }
func (c *ColorClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *ColorClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *ColorClass) GetMethods() []data.Method                       { return nil }
func (c *ColorClass) GetConstruct() data.Method                       { return &colorConstruct{} }

type colorConstruct struct{}

func (m *colorConstruct) GetName() string            { return token.ConstructName }
func (m *colorConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *colorConstruct) GetIsStatic() bool          { return false }
func (m *colorConstruct) GetReturnType() data.Types  { return nil }

func (m *colorConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "r", 0, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "g", 1, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "b", 2, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "a", 3, data.NewIntValue(255), data.NewBaseType("int")),
	}
}

func (m *colorConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "r", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "g", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "b", 2, data.NewBaseType("int")),
		node.NewVariable(nil, "a", 3, data.NewBaseType("int")),
	}
}

func (m *colorConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	r, g, b, a := 0, 0, 0, 255
	if v, ok := ctx.GetIndexValue(0); ok {
		if i, ok := v.(data.AsInt); ok {
			r, _ = i.AsInt()
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		if i, ok := v.(data.AsInt); ok {
			g, _ = i.AsInt()
		}
	}
	if v, ok := ctx.GetIndexValue(2); ok {
		if i, ok := v.(data.AsInt); ok {
			b, _ = i.AsInt()
		}
	}
	if v, ok := ctx.GetIndexValue(3); ok {
		if i, ok := v.(data.AsInt); ok {
			a, _ = i.AsInt()
		}
	}
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			classVal.SetProperty("r", data.NewIntValue(r))
			classVal.SetProperty("g", data.NewIntValue(g))
			classVal.SetProperty("b", data.NewIntValue(b))
			classVal.SetProperty("a", data.NewIntValue(a))
			classVal.SetProperty("_color", data.NewAnyValue(color.NRGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}))
		}
	}
	return nil, nil
}

// NewColorValue 创建一个 Fyne\Color 的 ClassValue
func NewColorValue(c color.Color, ctx data.Context) *data.ClassValue {
	r, g, b, a := c.RGBA()
	cl := NewColorClass()
	cv, _ := cl.GetValue(ctx)
	if classVal, ok := cv.(*data.ClassValue); ok {
		classVal.SetProperty("r", data.NewIntValue(int(r>>8)))
		classVal.SetProperty("g", data.NewIntValue(int(g>>8)))
		classVal.SetProperty("b", data.NewIntValue(int(b>>8)))
		classVal.SetProperty("a", data.NewIntValue(int(a>>8)))
		classVal.SetProperty("_color", data.NewAnyValue(c))
		return classVal
	}
	return nil
}
