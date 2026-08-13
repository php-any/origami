package fyne

import (
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// RadioGroupClass 是 Fyne\Widget\RadioGroup 类
type RadioGroupClass struct {
	*CanvasObjectClass
}

func NewRadioGroupClass() data.ClassStmt {
	return &RadioGroupClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\RadioGroup", nil),
	}
}

func (c *RadioGroupClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *RadioGroupClass) GetConstruct() data.Method { return &radioGroupConstruct{} }

func (c *RadioGroupClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setSelected":
		return &radioGroupSetSelectedMethod{}, true
	case "getSelected":
		return &radioGroupGetSelectedMethod{}, true
	case "onChanged":
		return &radioGroupOnChangedMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *RadioGroupClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&radioGroupSetSelectedMethod{},
		&radioGroupGetSelectedMethod{},
		&radioGroupOnChangedMethod{},
	)
}

type radioGroupConstruct struct{}

func (m *radioGroupConstruct) GetName() string            { return token.ConstructName }
func (m *radioGroupConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *radioGroupConstruct) GetIsStatic() bool          { return false }
func (m *radioGroupConstruct) GetReturnType() data.Types  { return nil }
func (m *radioGroupConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, nil, data.NewBaseType("array")),
		node.NewParameter(nil, "callback", 1, data.NewNullValue(), nil),
	}
}
func (m *radioGroupConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "callback", 1, nil),
	}
}

func (m *radioGroupConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
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
	rg := widget.NewRadioGroup(options, func(selected string) {
		callPHPCallbackWith(callback, ctx, data.NewStringValue(selected))
	})
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			setFyneObject(classVal, rg)
			classVal.SetProperty("_radioGroup", data.NewAnyValue(rg))
		}
	}
	return nil, nil
}

func getRadioGroup(cv *data.ClassValue) *widget.RadioGroup {
	if v, _ := cv.GetProperty("_radioGroup"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if rg, ok := av.Value.(*widget.RadioGroup); ok {
				return rg
			}
		}
	}
	return nil
}

type radioGroupSetSelectedMethod struct{}

func (m *radioGroupSetSelectedMethod) GetName() string            { return "setSelected" }
func (m *radioGroupSetSelectedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *radioGroupSetSelectedMethod) GetIsStatic() bool          { return false }
func (m *radioGroupSetSelectedMethod) GetReturnType() data.Types  { return nil }
func (m *radioGroupSetSelectedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, data.NewBaseType("string")),
	}
}
func (m *radioGroupSetSelectedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("string")),
	}
}
func (m *radioGroupSetSelectedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if rg := getRadioGroup(classVal); rg != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						rg.SetSelected(s.AsString())
					}
				}
			}
		}
	}
	return nil, nil
}

type radioGroupGetSelectedMethod struct{}

func (m *radioGroupGetSelectedMethod) GetName() string               { return "getSelected" }
func (m *radioGroupGetSelectedMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *radioGroupGetSelectedMethod) GetIsStatic() bool             { return false }
func (m *radioGroupGetSelectedMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
func (m *radioGroupGetSelectedMethod) GetParams() []data.GetValue    { return nil }
func (m *radioGroupGetSelectedMethod) GetVariables() []data.Variable { return nil }
func (m *radioGroupGetSelectedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if rg := getRadioGroup(classVal); rg != nil {
				return data.NewStringValue(rg.Selected), nil
			}
		}
	}
	return data.NewStringValue(""), nil
}

type radioGroupOnChangedMethod struct{}

func (m *radioGroupOnChangedMethod) GetName() string            { return "onChanged" }
func (m *radioGroupOnChangedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *radioGroupOnChangedMethod) GetIsStatic() bool          { return false }
func (m *radioGroupOnChangedMethod) GetReturnType() data.Types  { return nil }
func (m *radioGroupOnChangedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
	}
}
func (m *radioGroupOnChangedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, nil),
	}
}
func (m *radioGroupOnChangedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if rg := getRadioGroup(classVal); rg != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					var callback data.FuncStmt
					if fv, ok := v.(*data.FuncValue); ok {
						callback = fv.Value
					}
					rg.OnChanged = func(selected string) {
						callPHPCallbackWith(callback, ctx, data.NewStringValue(selected))
					}
				}
			}
		}
	}
	return nil, nil
}
