package fyne

import (
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// CheckClass 是 Fyne\Widget\Check 类
type CheckClass struct {
	*CanvasObjectClass
}

func NewCheckClass() data.ClassStmt {
	return &CheckClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\Check", nil),
	}
}

func (c *CheckClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *CheckClass) GetConstruct() data.Method { return &checkConstruct{} }

func (c *CheckClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setChecked":
		return &checkSetCheckedMethod{}, true
	case "isChecked":
		return &checkIsCheckedMethod{}, true
	case "onChanged":
		return &checkOnChangedMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *CheckClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&checkSetCheckedMethod{},
		&checkIsCheckedMethod{},
		&checkOnChangedMethod{},
	)
}

type checkConstruct struct{}

func (m *checkConstruct) GetName() string            { return token.ConstructName }
func (m *checkConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *checkConstruct) GetIsStatic() bool          { return false }
func (m *checkConstruct) GetReturnType() data.Types  { return nil }
func (m *checkConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "label", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "callback", 1, data.NewNullValue(), nil),
	}
}
func (m *checkConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "label", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "callback", 1, nil),
	}
}

func (m *checkConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	label := ""
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			label = s.AsString()
		}
	}
	var callback data.GetValue
	if v, ok := ctx.GetIndexValue(1); ok {
		callback = v
	}
	check := widget.NewCheck(label, func(checked bool) {
		if callback != nil {
			callback.GetValue(ctx)
		}
	})
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			setFyneObject(classVal, check)
			classVal.SetProperty("_check", data.NewAnyValue(check))
		}
	}
	return nil, nil
}

func getCheck(cv *data.ClassValue) *widget.Check {
	if v, _ := cv.GetProperty("_check"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if c, ok := av.Value.(*widget.Check); ok {
				return c
			}
		}
	}
	return nil
}

type checkSetCheckedMethod struct{}

func (m *checkSetCheckedMethod) GetName() string            { return "setChecked" }
func (m *checkSetCheckedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *checkSetCheckedMethod) GetIsStatic() bool          { return false }
func (m *checkSetCheckedMethod) GetReturnType() data.Types  { return nil }
func (m *checkSetCheckedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "checked", 0, nil, data.NewBaseType("bool")),
	}
}
func (m *checkSetCheckedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "checked", 0, data.NewBaseType("bool")),
	}
}
func (m *checkSetCheckedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if c := getCheck(classVal); c != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if b, ok := v.(data.AsBool); ok {
						checked, _ := b.AsBool()
						c.SetChecked(checked)
					}
				}
			}
		}
	}
	return nil, nil
}

type checkIsCheckedMethod struct{}

func (m *checkIsCheckedMethod) GetName() string               { return "isChecked" }
func (m *checkIsCheckedMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *checkIsCheckedMethod) GetIsStatic() bool             { return false }
func (m *checkIsCheckedMethod) GetReturnType() data.Types     { return data.NewBaseType("bool") }
func (m *checkIsCheckedMethod) GetParams() []data.GetValue    { return nil }
func (m *checkIsCheckedMethod) GetVariables() []data.Variable { return nil }
func (m *checkIsCheckedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if c := getCheck(classVal); c != nil {
				return data.NewBoolValue(c.Checked), nil
			}
		}
	}
	return data.NewBoolValue(false), nil
}

type checkOnChangedMethod struct{}

func (m *checkOnChangedMethod) GetName() string            { return "onChanged" }
func (m *checkOnChangedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *checkOnChangedMethod) GetIsStatic() bool          { return false }
func (m *checkOnChangedMethod) GetReturnType() data.Types  { return nil }
func (m *checkOnChangedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
	}
}
func (m *checkOnChangedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, nil),
	}
}
func (m *checkOnChangedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if c := getCheck(classVal); c != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					callback := v
					c.OnChanged = func(checked bool) {
						callback.GetValue(ctx)
					}
				}
			}
		}
	}
	return nil, nil
}
