package fyne

import (
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// EntryClass 是 Fyne\Widget\Entry 类
type EntryClass struct {
	*CanvasObjectClass
}

func NewEntryClass() data.ClassStmt {
	return &EntryClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\Entry", nil),
	}
}

func (c *EntryClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *EntryClass) GetConstruct() data.Method { return &entryConstruct{} }

func (c *EntryClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setText":
		return &entrySetTextMethod{}, true
	case "getText":
		return &entryGetTextMethod{}, true
	case "setPlaceHolder":
		return &entrySetPlaceHolderMethod{}, true
	case "setPassword":
		return &entrySetPasswordMethod{}, true
	case "setMultiLine":
		return &entrySetMultiLineMethod{}, true
	case "onChanged":
		return &entryOnChangedMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *EntryClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&entrySetTextMethod{},
		&entryGetTextMethod{},
		&entrySetPlaceHolderMethod{},
		&entrySetPasswordMethod{},
		&entrySetMultiLineMethod{},
		&entryOnChangedMethod{},
	)
}

type entryConstruct struct{}

func (m *entryConstruct) GetName() string               { return token.ConstructName }
func (m *entryConstruct) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *entryConstruct) GetIsStatic() bool             { return false }
func (m *entryConstruct) GetReturnType() data.Types     { return nil }
func (m *entryConstruct) GetParams() []data.GetValue    { return nil }
func (m *entryConstruct) GetVariables() []data.Variable { return nil }

func (m *entryConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	entry := widget.NewEntry()
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			setFyneObject(classVal, entry)
			classVal.SetProperty("_entry", data.NewAnyValue(entry))
		}
	}
	return nil, nil
}

func getEntry(cv *data.ClassValue) *widget.Entry {
	if v, _ := cv.GetProperty("_entry"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if e, ok := av.Value.(*widget.Entry); ok {
				return e
			}
		}
	}
	return nil
}

type entrySetTextMethod struct{}

func (m *entrySetTextMethod) GetName() string            { return "setText" }
func (m *entrySetTextMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *entrySetTextMethod) GetIsStatic() bool          { return false }
func (m *entrySetTextMethod) GetReturnType() data.Types  { return nil }
func (m *entrySetTextMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "text", 0, nil, data.NewBaseType("string")),
	}
}
func (m *entrySetTextMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "text", 0, data.NewBaseType("string")),
	}
}
func (m *entrySetTextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if e := getEntry(classVal); e != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						e.SetText(s.AsString())
					}
				}
			}
		}
	}
	return nil, nil
}

type entryGetTextMethod struct{}

func (m *entryGetTextMethod) GetName() string               { return "getText" }
func (m *entryGetTextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *entryGetTextMethod) GetIsStatic() bool             { return false }
func (m *entryGetTextMethod) GetReturnType() data.Types     { return data.NewBaseType("string") }
func (m *entryGetTextMethod) GetParams() []data.GetValue    { return nil }
func (m *entryGetTextMethod) GetVariables() []data.Variable { return nil }
func (m *entryGetTextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if e := getEntry(classVal); e != nil {
				return data.NewStringValue(e.Text), nil
			}
		}
	}
	return data.NewStringValue(""), nil
}

type entrySetPlaceHolderMethod struct{}

func (m *entrySetPlaceHolderMethod) GetName() string            { return "setPlaceHolder" }
func (m *entrySetPlaceHolderMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *entrySetPlaceHolderMethod) GetIsStatic() bool          { return false }
func (m *entrySetPlaceHolderMethod) GetReturnType() data.Types  { return nil }
func (m *entrySetPlaceHolderMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "placeholder", 0, nil, data.NewBaseType("string")),
	}
}
func (m *entrySetPlaceHolderMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "placeholder", 0, data.NewBaseType("string")),
	}
}
func (m *entrySetPlaceHolderMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if e := getEntry(classVal); e != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						e.SetPlaceHolder(s.AsString())
					}
				}
			}
		}
	}
	return nil, nil
}

type entrySetPasswordMethod struct{}

func (m *entrySetPasswordMethod) GetName() string            { return "setPassword" }
func (m *entrySetPasswordMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *entrySetPasswordMethod) GetIsStatic() bool          { return false }
func (m *entrySetPasswordMethod) GetReturnType() data.Types  { return nil }
func (m *entrySetPasswordMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "password", 0, data.NewBoolValue(true), data.NewBaseType("bool")),
	}
}
func (m *entrySetPasswordMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "password", 0, data.NewBaseType("bool")),
	}
}
func (m *entrySetPasswordMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if e := getEntry(classVal); e != nil {
				password := true
				if v, ok := ctx.GetIndexValue(0); ok {
					if b, ok := v.(data.AsBool); ok {
						password, _ = b.AsBool()
					}
				}
				e.Password = password
			}
		}
	}
	return nil, nil
}

type entrySetMultiLineMethod struct{}

func (m *entrySetMultiLineMethod) GetName() string            { return "setMultiLine" }
func (m *entrySetMultiLineMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *entrySetMultiLineMethod) GetIsStatic() bool          { return false }
func (m *entrySetMultiLineMethod) GetReturnType() data.Types  { return nil }
func (m *entrySetMultiLineMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "multiLine", 0, data.NewBoolValue(true), data.NewBaseType("bool")),
	}
}
func (m *entrySetMultiLineMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "multiLine", 0, data.NewBaseType("bool")),
	}
}
func (m *entrySetMultiLineMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if e := getEntry(classVal); e != nil {
				multiLine := true
				if v, ok := ctx.GetIndexValue(0); ok {
					if b, ok := v.(data.AsBool); ok {
						multiLine, _ = b.AsBool()
					}
				}
				e.MultiLine = multiLine
			}
		}
	}
	return nil, nil
}

type entryOnChangedMethod struct{}

func (m *entryOnChangedMethod) GetName() string            { return "onChanged" }
func (m *entryOnChangedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *entryOnChangedMethod) GetIsStatic() bool          { return false }
func (m *entryOnChangedMethod) GetReturnType() data.Types  { return nil }
func (m *entryOnChangedMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
	}
}
func (m *entryOnChangedMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, nil),
	}
}
func (m *entryOnChangedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if e := getEntry(classVal); e != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					var callback data.FuncStmt
					if fv, ok := v.(*data.FuncValue); ok {
						callback = fv.Value
					}
					e.OnChanged = func(text string) {
						callPHPCallbackWith(callback, ctx, data.NewStringValue(text))
					}
				}
			}
		}
	}
	return nil, nil
}
