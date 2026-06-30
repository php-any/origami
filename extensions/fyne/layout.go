package fyne

import (
	fyneLib "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/layout"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// LayoutClass 是 Fyne\Layout 基类，提供静态工厂方法
type LayoutClass struct{}

func NewLayoutClass() data.ClassStmt { return &LayoutClass{} }

func (c *LayoutClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *LayoutClass) GetFrom() data.From                            { return nil }
func (c *LayoutClass) GetName() string                               { return "Fyne\\Layout" }
func (c *LayoutClass) GetExtend() *string                            { return nil }
func (c *LayoutClass) GetImplements() []string                       { return nil }
func (c *LayoutClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *LayoutClass) GetPropertyList() []data.Property              { return nil }
func (c *LayoutClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *LayoutClass) GetMethods() []data.Method                     { return nil }
func (c *LayoutClass) GetConstruct() data.Method                     { return nil }

func (c *LayoutClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "newHBoxLayout":
		return &layoutNewHBoxLayoutMethod{}, true
	case "newVBoxLayout":
		return &layoutNewVBoxLayoutMethod{}, true
	case "newGridLayout":
		return &layoutNewGridLayoutMethod{}, true
	case "newCenterLayout":
		return &layoutNewCenterLayoutMethod{}, true
	case "newMaxLayout":
		return &layoutNewMaxLayoutMethod{}, true
	case "newFormLayout":
		return &layoutNewFormLayoutMethod{}, true
	case "newPaddedLayout":
		return &layoutNewPaddedLayoutMethod{}, true
	default:
		return nil, false
	}
}

// wrapLayout 将 fyne.Layout 包装为 ClassValue
func wrapLayout(l fyneLib.Layout, ctx data.Context) *data.ClassValue {
	cl := NewLayoutClass()
	cv, _ := cl.GetValue(ctx)
	if classVal, ok := cv.(*data.ClassValue); ok {
		classVal.SetProperty("_layout", data.NewAnyValue(l))
		return classVal
	}
	return nil
}

// ====== HBox ======

type layoutNewHBoxLayoutMethod struct{}

func (m *layoutNewHBoxLayoutMethod) GetName() string            { return "newHBoxLayout" }
func (m *layoutNewHBoxLayoutMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *layoutNewHBoxLayoutMethod) GetIsStatic() bool          { return true }
func (m *layoutNewHBoxLayoutMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Layout")
}
func (m *layoutNewHBoxLayoutMethod) GetParams() []data.GetValue    { return nil }
func (m *layoutNewHBoxLayoutMethod) GetVariables() []data.Variable { return nil }
func (m *layoutNewHBoxLayoutMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return wrapLayout(layout.NewHBoxLayout(), ctx), nil
}

// ====== VBox ======

type layoutNewVBoxLayoutMethod struct{}

func (m *layoutNewVBoxLayoutMethod) GetName() string            { return "newVBoxLayout" }
func (m *layoutNewVBoxLayoutMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *layoutNewVBoxLayoutMethod) GetIsStatic() bool          { return true }
func (m *layoutNewVBoxLayoutMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Layout")
}
func (m *layoutNewVBoxLayoutMethod) GetParams() []data.GetValue    { return nil }
func (m *layoutNewVBoxLayoutMethod) GetVariables() []data.Variable { return nil }
func (m *layoutNewVBoxLayoutMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return wrapLayout(layout.NewVBoxLayout(), ctx), nil
}

// ====== Grid ======

type layoutNewGridLayoutMethod struct{}

func (m *layoutNewGridLayoutMethod) GetName() string            { return "newGridLayout" }
func (m *layoutNewGridLayoutMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *layoutNewGridLayoutMethod) GetIsStatic() bool          { return true }
func (m *layoutNewGridLayoutMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Layout")
}
func (m *layoutNewGridLayoutMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "cols", 0, nil, data.NewBaseType("int")),
	}
}
func (m *layoutNewGridLayoutMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "cols", 0, data.NewBaseType("int")),
	}
}
func (m *layoutNewGridLayoutMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cols := 2
	if v, ok := ctx.GetIndexValue(0); ok {
		if i, ok := v.(data.AsInt); ok {
			cols, _ = i.AsInt()
		}
	}
	return wrapLayout(layout.NewGridLayout(cols), ctx), nil
}

// ====== Center ======

type layoutNewCenterLayoutMethod struct{}

func (m *layoutNewCenterLayoutMethod) GetName() string            { return "newCenterLayout" }
func (m *layoutNewCenterLayoutMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *layoutNewCenterLayoutMethod) GetIsStatic() bool          { return true }
func (m *layoutNewCenterLayoutMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Layout")
}
func (m *layoutNewCenterLayoutMethod) GetParams() []data.GetValue    { return nil }
func (m *layoutNewCenterLayoutMethod) GetVariables() []data.Variable { return nil }
func (m *layoutNewCenterLayoutMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return wrapLayout(layout.NewCenterLayout(), ctx), nil
}

// ====== Max ======

type layoutNewMaxLayoutMethod struct{}

func (m *layoutNewMaxLayoutMethod) GetName() string            { return "newMaxLayout" }
func (m *layoutNewMaxLayoutMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *layoutNewMaxLayoutMethod) GetIsStatic() bool          { return true }
func (m *layoutNewMaxLayoutMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Layout")
}
func (m *layoutNewMaxLayoutMethod) GetParams() []data.GetValue    { return nil }
func (m *layoutNewMaxLayoutMethod) GetVariables() []data.Variable { return nil }
func (m *layoutNewMaxLayoutMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return wrapLayout(layout.NewStackLayout(), ctx), nil
}

// ====== Form ======

type layoutNewFormLayoutMethod struct{}

func (m *layoutNewFormLayoutMethod) GetName() string            { return "newFormLayout" }
func (m *layoutNewFormLayoutMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *layoutNewFormLayoutMethod) GetIsStatic() bool          { return true }
func (m *layoutNewFormLayoutMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Layout")
}
func (m *layoutNewFormLayoutMethod) GetParams() []data.GetValue    { return nil }
func (m *layoutNewFormLayoutMethod) GetVariables() []data.Variable { return nil }
func (m *layoutNewFormLayoutMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return wrapLayout(layout.NewFormLayout(), ctx), nil
}

// ====== Padded ======

type layoutNewPaddedLayoutMethod struct{}

func (m *layoutNewPaddedLayoutMethod) GetName() string            { return "newPaddedLayout" }
func (m *layoutNewPaddedLayoutMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *layoutNewPaddedLayoutMethod) GetIsStatic() bool          { return true }
func (m *layoutNewPaddedLayoutMethod) GetReturnType() data.Types {
	return data.NewBaseType("Fyne\\Layout")
}
func (m *layoutNewPaddedLayoutMethod) GetParams() []data.GetValue    { return nil }
func (m *layoutNewPaddedLayoutMethod) GetVariables() []data.Variable { return nil }
func (m *layoutNewPaddedLayoutMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return wrapLayout(layout.NewPaddedLayout(), ctx), nil
}
