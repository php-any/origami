package fyne

import (
	fyneLib "fyne.io/fyne/v2"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// SizeClass 是 Fyne\Size 值对象类
type SizeClass struct{}

func NewSizeClass() data.ClassStmt { return &SizeClass{} }

func (c *SizeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *SizeClass) GetFrom() data.From                              { return nil }
func (c *SizeClass) GetName() string                                 { return "Fyne\\Size" }
func (c *SizeClass) GetExtend() *string                              { return nil }
func (c *SizeClass) GetImplements() []string                         { return nil }
func (c *SizeClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *SizeClass) GetPropertyList() []data.Property                { return nil }
func (c *SizeClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *SizeClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *SizeClass) GetMethods() []data.Method                       { return nil }
func (c *SizeClass) GetConstruct() data.Method                       { return &sizeConstruct{} }

type sizeConstruct struct{}

func (m *sizeConstruct) GetName() string            { return token.ConstructName }
func (m *sizeConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *sizeConstruct) GetIsStatic() bool          { return false }
func (m *sizeConstruct) GetReturnType() data.Types  { return nil }

func (m *sizeConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "width", 0, data.NewFloatValue(0), data.NewBaseType("float")),
		node.NewParameter(nil, "height", 1, data.NewFloatValue(0), data.NewBaseType("float")),
	}
}

func (m *sizeConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "width", 0, data.NewBaseType("float")),
		node.NewVariable(nil, "height", 1, data.NewBaseType("float")),
	}
}

func (m *sizeConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	width := 0.0
	height := 0.0
	if v, ok := ctx.GetIndexValue(0); ok {
		if f, ok := v.(data.AsFloat); ok {
			width, _ = f.AsFloat()
		} else if i, ok := v.(data.AsInt); ok {
			n, _ := i.AsInt()
			width = float64(n)
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		if f, ok := v.(data.AsFloat); ok {
			height, _ = f.AsFloat()
		} else if i, ok := v.(data.AsInt); ok {
			n, _ := i.AsInt()
			height = float64(n)
		}
	}
	// 将 fyne.Size 存储到对象的隐藏属性
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			classVal.SetProperty("_size", data.NewAnyValue(fyneLib.NewSize(float32(width), float32(height))))
		}
	}
	return nil, nil
}

// NewSizeValue 创建一个 Fyne\Size 的 ClassValue
func NewSizeValue(size fyneLib.Size, ctx data.Context) *data.ClassValue {
	cl := NewSizeClass()
	cv, _ := cl.GetValue(ctx)
	if classVal, ok := cv.(*data.ClassValue); ok {
		classVal.SetProperty("width", data.NewFloatValue(float64(size.Width)))
		classVal.SetProperty("height", data.NewFloatValue(float64(size.Height)))
		classVal.SetProperty("_size", data.NewAnyValue(size))
		return classVal
	}
	return nil
}
