package fyne

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// BottomTabBarClass ZY 类: Fyne\Widget\BottomTabBar
type BottomTabBarClass struct {
	*CanvasObjectClass
}

func NewBottomTabBarClass() data.ClassStmt {
	return &BottomTabBarClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\BottomTabBar", nil),
	}
}

func (c *BottomTabBarClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *BottomTabBarClass) GetConstruct() data.Method { return &bottomTabBarConstruct{} }

func (c *BottomTabBarClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "append":
		return &bottomTabBarAppendMethod{}, true
	case "setSelected":
		return &bottomTabBarSetSelectedMethod{}, true
	case "getSelected":
		return &bottomTabBarGetSelectedMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *BottomTabBarClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&bottomTabBarAppendMethod{},
		&bottomTabBarSetSelectedMethod{},
		&bottomTabBarGetSelectedMethod{},
	)
}

// ── 构造函数 ──

type bottomTabBarConstruct struct{}

func (m *bottomTabBarConstruct) GetName() string               { return token.ConstructName }
func (m *bottomTabBarConstruct) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *bottomTabBarConstruct) GetIsStatic() bool             { return false }
func (m *bottomTabBarConstruct) GetReturnType() data.Types     { return nil }
func (m *bottomTabBarConstruct) GetParams() []data.GetValue    { return nil }
func (m *bottomTabBarConstruct) GetVariables() []data.Variable { return nil }

func (m *bottomTabBarConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	bar := NewBottomTabBar()
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			setFyneObject(classVal, bar)
			classVal.SetProperty("_bottomTabBar", data.NewAnyValue(bar))
		}
	}
	return nil, nil
}

func getBottomTabBar(cv *data.ClassValue) *BottomTabBar {
	if v, _ := cv.GetProperty("_bottomTabBar"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if bar, ok := av.Value.(*BottomTabBar); ok {
				return bar
			}
		}
	}
	return nil
}

// ── append 方法 ──

type bottomTabBarAppendMethod struct{}

func (m *bottomTabBarAppendMethod) GetName() string            { return "append" }
func (m *bottomTabBarAppendMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *bottomTabBarAppendMethod) GetIsStatic() bool          { return false }
func (m *bottomTabBarAppendMethod) GetReturnType() data.Types  { return nil }
func (m *bottomTabBarAppendMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "title", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "icon", 1, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "callback", 2, nil, nil),
	}
}
func (m *bottomTabBarAppendMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "title", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "icon", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "callback", 2, nil),
	}
}

func (m *bottomTabBarAppendMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if bar := getBottomTabBar(classVal); bar != nil {
				title, icon := "", ""
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						title = s.AsString()
					}
				}
				if v, ok := ctx.GetIndexValue(1); ok {
					if s, ok := v.(data.AsString); ok {
						icon = s.AsString()
					}
				}
				// 捕获回调
				var callback data.FuncStmt
				if v, ok := ctx.GetIndexValue(2); ok {
					if fv, ok := v.(*data.FuncValue); ok {
						callback = fv.Value
					}
				}
				bar.Append(title, icon, func() {
					callPHPCallback(callback, ctx)
				})
			}
		}
	}
	return nil, nil
}

// ── setSelected 方法 ──

type bottomTabBarSetSelectedMethod struct{}

func (m *bottomTabBarSetSelectedMethod) GetName() string            { return "setSelected" }
func (m *bottomTabBarSetSelectedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *bottomTabBarSetSelectedMethod) GetIsStatic() bool          { return false }
func (m *bottomTabBarSetSelectedMethod) GetReturnType() data.Types  { return nil }
func (m *bottomTabBarSetSelectedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "index", 0, nil, data.NewBaseType("int")),
	}
}
func (m *bottomTabBarSetSelectedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "index", 0, data.NewBaseType("int")),
	}
}

func (m *bottomTabBarSetSelectedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if bar := getBottomTabBar(classVal); bar != nil {
				idx := 0
				if v, ok := ctx.GetIndexValue(0); ok {
					if i, ok := v.(data.AsInt); ok {
						idx, _ = i.AsInt()
					}
				}
				bar.SetSelected(idx)
			}
		}
	}
	return nil, nil
}

// ── getSelected 方法 ──

type bottomTabBarGetSelectedMethod struct{}

func (m *bottomTabBarGetSelectedMethod) GetName() string               { return "getSelected" }
func (m *bottomTabBarGetSelectedMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *bottomTabBarGetSelectedMethod) GetIsStatic() bool             { return false }
func (m *bottomTabBarGetSelectedMethod) GetReturnType() data.Types     { return data.NewBaseType("int") }
func (m *bottomTabBarGetSelectedMethod) GetParams() []data.GetValue    { return nil }
func (m *bottomTabBarGetSelectedMethod) GetVariables() []data.Variable { return nil }

func (m *bottomTabBarGetSelectedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if bar := getBottomTabBar(classVal); bar != nil {
				return data.NewIntValue(bar.Selected), nil
			}
		}
	}
	return data.NewIntValue(0), nil
}
