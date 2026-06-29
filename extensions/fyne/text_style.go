package fyne

import (
	fyneLib "fyne.io/fyne/v2"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// TextStyleClass 是 Fyne\TextStyle 类
type TextStyleClass struct{}

func NewTextStyleClass() data.ClassStmt { return &TextStyleClass{} }

func (c *TextStyleClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *TextStyleClass) GetFrom() data.From                              { return nil }
func (c *TextStyleClass) GetName() string                                 { return "Fyne\\TextStyle" }
func (c *TextStyleClass) GetExtend() *string                              { return nil }
func (c *TextStyleClass) GetImplements() []string                         { return nil }
func (c *TextStyleClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *TextStyleClass) GetPropertyList() []data.Property                { return nil }
func (c *TextStyleClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *TextStyleClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *TextStyleClass) GetMethods() []data.Method                       { return nil }
func (c *TextStyleClass) GetConstruct() data.Method                       { return &textStyleConstruct{} }

type textStyleConstruct struct{}

func (m *textStyleConstruct) GetName() string            { return token.ConstructName }
func (m *textStyleConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *textStyleConstruct) GetIsStatic() bool          { return false }
func (m *textStyleConstruct) GetReturnType() data.Types  { return nil }

func (m *textStyleConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "bold", 0, data.NewBoolValue(false), data.NewBaseType("bool")),
		node.NewParameter(nil, "italic", 1, data.NewBoolValue(false), data.NewBaseType("bool")),
		node.NewParameter(nil, "monospace", 2, data.NewBoolValue(false), data.NewBaseType("bool")),
	}
}

func (m *textStyleConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "bold", 0, data.NewBaseType("bool")),
		node.NewVariable(nil, "italic", 1, data.NewBaseType("bool")),
		node.NewVariable(nil, "monospace", 2, data.NewBaseType("bool")),
	}
}

func (m *textStyleConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	style := fyneLib.TextStyle{}
	if v, ok := ctx.GetIndexValue(0); ok {
		if b, ok := v.(data.AsBool); ok {
			style.Bold, _ = b.AsBool()
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		if b, ok := v.(data.AsBool); ok {
			style.Italic, _ = b.AsBool()
		}
	}
	if v, ok := ctx.GetIndexValue(2); ok {
		if b, ok := v.(data.AsBool); ok {
			style.Monospace, _ = b.AsBool()
		}
	}
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			classVal.SetProperty("_style", data.NewAnyValue(style))
		}
	}
	return nil, nil
}

// getFyneTextStyle 从 data.Value 中提取 fyne.TextStyle
func getFyneTextStyle(v data.Value) fyneLib.TextStyle {
	if cv, ok := v.(*data.ClassValue); ok {
		if s, _ := cv.GetProperty("_style"); s != nil {
			if av, ok := s.(*data.AnyValue); ok {
				if ts, ok := av.Value.(fyneLib.TextStyle); ok {
					return ts
				}
			}
		}
	}
	return fyneLib.TextStyle{}
}
