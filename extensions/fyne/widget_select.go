package fyne

import (
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// SelectClass 是 Fyne\Widget\Select 类
type SelectClass struct {
	*CanvasObjectClass
}

func NewSelectClass() data.ClassStmt {
	return &SelectClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\Select", nil),
	}
}

func (c *SelectClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *SelectClass) GetConstruct() data.Method { return &selectConstruct{} }

func (c *SelectClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setSelected":
		return &selectSetSelectedMethod{}, true
	case "clearSelected":
		return &selectClearSelectedMethod{}, true
	case "getSelected":
		return &selectGetSelectedMethod{}, true
	case "onChanged":
		return &selectOnChangedMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *SelectClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&selectSetSelectedMethod{},
		&selectClearSelectedMethod{},
		&selectGetSelectedMethod{},
		&selectOnChangedMethod{},
	)
}

type selectConstruct struct{}

func (m *selectConstruct) GetName() string            { return token.ConstructName }
func (m *selectConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *selectConstruct) GetIsStatic() bool          { return false }
func (m *selectConstruct) GetReturnType() data.Types  { return nil }
func (m *selectConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, nil, data.NewBaseType("array")),
		node.NewParameter(nil, "callback", 1, data.NewNullValue(), nil),
	}
}
func (m *selectConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "callback", 1, nil),
	}
}

func (m *selectConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	var options []string
	if v, ok := ctx.GetIndexValue(0); ok {
		if arr, ok := v.(*data.ArrayValue); ok {
			for _, item := range arr.ToValueList() {
				if s, ok := item.(data.AsString); ok {
					options = append(options, s.AsString())
				}
			}
		}
	}
	var callback data.FuncStmt
	if v, ok := ctx.GetIndexValue(1); ok {
		if fv, ok := v.(*data.FuncValue); ok {
			callback = fv.Value
		}
	}
	s := widget.NewSelect(options, func(selected string) {
		callPHPCallbackWith(callback, ctx, data.NewStringValue(selected))
	})
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			setFyneObject(classVal, s)
			classVal.SetProperty("_select", data.NewAnyValue(s))
		}
	}
	return nil, nil
}

func getSelect(cv *data.ClassValue) *widget.Select {
	if v, _ := cv.GetProperty("_select"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if s, ok := av.Value.(*widget.Select); ok {
				return s
			}
		}
	}
	return nil
}

type selectSetSelectedMethod struct{}

func (m *selectSetSelectedMethod) GetName() string            { return "setSelected" }
func (m *selectSetSelectedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *selectSetSelectedMethod) GetIsStatic() bool          { return false }
func (m *selectSetSelectedMethod) GetReturnType() data.Types  { return nil }
func (m *selectSetSelectedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, data.NewBaseType("string")),
	}
}
func (m *selectSetSelectedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("string")),
	}
}
func (m *selectSetSelectedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if s := getSelect(classVal); s != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if str, ok := v.(data.AsString); ok {
						s.SetSelected(str.AsString())
					}
				}
			}
		}
	}
	return nil, nil
}

type selectClearSelectedMethod struct{}

func (m *selectClearSelectedMethod) GetName() string               { return "clearSelected" }
func (m *selectClearSelectedMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *selectClearSelectedMethod) GetIsStatic() bool             { return false }
func (m *selectClearSelectedMethod) GetReturnType() data.Types     { return nil }
func (m *selectClearSelectedMethod) GetParams() []data.GetValue    { return nil }
func (m *selectClearSelectedMethod) GetVariables() []data.Variable { return nil }
func (m *selectClearSelectedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if s := getSelect(classVal); s != nil {
				s.ClearSelected()
			}
		}
	}
	return nil, nil
}

type selectGetSelectedMethod struct{}

func (m *selectGetSelectedMethod) GetName() string               { return "getSelected" }
func (m *selectGetSelectedMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *selectGetSelectedMethod) GetIsStatic() bool             { return false }
func (m *selectGetSelectedMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
func (m *selectGetSelectedMethod) GetParams() []data.GetValue    { return nil }
func (m *selectGetSelectedMethod) GetVariables() []data.Variable { return nil }
func (m *selectGetSelectedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if s := getSelect(classVal); s != nil {
				return data.NewStringValue(s.Selected), nil
			}
		}
	}
	return data.NewStringValue(""), nil
}

type selectOnChangedMethod struct{}

func (m *selectOnChangedMethod) GetName() string            { return "onChanged" }
func (m *selectOnChangedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *selectOnChangedMethod) GetIsStatic() bool          { return false }
func (m *selectOnChangedMethod) GetReturnType() data.Types  { return nil }
func (m *selectOnChangedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
	}
}
func (m *selectOnChangedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, nil),
	}
}
func (m *selectOnChangedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if s := getSelect(classVal); s != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					var callback data.FuncStmt
					if fv, ok := v.(*data.FuncValue); ok {
						callback = fv.Value
					}
					s.OnChanged = func(selected string) {
						callPHPCallbackWith(callback, ctx, data.NewStringValue(selected))
					}
				}
			}
		}
	}
	return nil, nil
}
