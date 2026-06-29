package fyne

import (
	fyneLib "fyne.io/fyne/v2"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// PositionClass 是 Fyne\Position 值对象类
type PositionClass struct{}

func NewPositionClass() data.ClassStmt { return &PositionClass{} }

func (c *PositionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *PositionClass) GetFrom() data.From                              { return nil }
func (c *PositionClass) GetName() string                                 { return "Fyne\\Position" }
func (c *PositionClass) GetExtend() *string                              { return nil }
func (c *PositionClass) GetImplements() []string                         { return nil }
func (c *PositionClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *PositionClass) GetPropertyList() []data.Property                { return nil }
func (c *PositionClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *PositionClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *PositionClass) GetMethods() []data.Method                       { return nil }
func (c *PositionClass) GetConstruct() data.Method                       { return &positionConstruct{} }

type positionConstruct struct{}

func (m *positionConstruct) GetName() string            { return token.ConstructName }
func (m *positionConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *positionConstruct) GetIsStatic() bool          { return false }
func (m *positionConstruct) GetReturnType() data.Types  { return nil }

func (m *positionConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "x", 0, data.NewFloatValue(0), data.NewBaseType("float")),
		node.NewParameter(nil, "y", 1, data.NewFloatValue(0), data.NewBaseType("float")),
	}
}

func (m *positionConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "x", 0, data.NewBaseType("float")),
		node.NewVariable(nil, "y", 1, data.NewBaseType("float")),
	}
}

func (m *positionConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	x := 0.0
	y := 0.0
	if v, ok := ctx.GetIndexValue(0); ok {
		if f, ok := v.(data.AsFloat); ok {
			x, _ = f.AsFloat()
		} else if i, ok := v.(data.AsInt); ok {
			n, _ := i.AsInt()
			x = float64(n)
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		if f, ok := v.(data.AsFloat); ok {
			y, _ = f.AsFloat()
		} else if i, ok := v.(data.AsInt); ok {
			n, _ := i.AsInt()
			y = float64(n)
		}
	}
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			classVal.SetProperty("_pos", data.NewAnyValue(fyneLib.NewPos(float32(x), float32(y))))
		}
	}
	return nil, nil
}

// NewPositionValue 创建一个 Fyne\Position 的 ClassValue
func NewPositionValue(pos fyneLib.Position, ctx data.Context) *data.ClassValue {
	cl := NewPositionClass()
	cv, _ := cl.GetValue(ctx)
	if classVal, ok := cv.(*data.ClassValue); ok {
		classVal.SetProperty("x", data.NewFloatValue(float64(pos.X)))
		classVal.SetProperty("y", data.NewFloatValue(float64(pos.Y)))
		classVal.SetProperty("_pos", data.NewAnyValue(pos))
		return classVal
	}
	return nil
}
