package fyne

import (
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// FormClass 是 Fyne\Widget\Form 类
type FormClass struct {
	*CanvasObjectClass
}

func NewFormClass() data.ClassStmt {
	return &FormClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\Form", nil),
	}
}

func (c *FormClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *FormClass) GetConstruct() data.Method { return &formConstruct{} }

func (c *FormClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "append":
		return &formAppendMethod{}, true
	case "onSubmit":
		return &formOnSubmitMethod{}, true
	case "onCancel":
		return &formOnCancelMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *FormClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&formAppendMethod{},
		&formOnSubmitMethod{},
		&formOnCancelMethod{},
	)
}

type formConstruct struct{}

func (m *formConstruct) GetName() string               { return token.ConstructName }
func (m *formConstruct) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *formConstruct) GetIsStatic() bool             { return false }
func (m *formConstruct) GetReturnType() data.Types     { return nil }
func (m *formConstruct) GetParams() []data.GetValue    { return nil }
func (m *formConstruct) GetVariables() []data.Variable { return nil }

func (m *formConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	form := widget.NewForm()
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			setFyneObject(classVal, form)
			classVal.SetProperty("_form", data.NewAnyValue(form))
		}
	}
	return nil, nil
}

func getForm(cv *data.ClassValue) *widget.Form {
	if v, _ := cv.GetProperty("_form"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if f, ok := av.Value.(*widget.Form); ok {
				return f
			}
		}
	}
	return nil
}

type formAppendMethod struct{}

func (m *formAppendMethod) GetName() string            { return "append" }
func (m *formAppendMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *formAppendMethod) GetIsStatic() bool          { return false }
func (m *formAppendMethod) GetReturnType() data.Types  { return nil }
func (m *formAppendMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "label", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "widget", 1, nil, data.NewBaseType("Fyne\\CanvasObject")),
	}
}
func (m *formAppendMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "label", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "widget", 1, data.NewBaseType("Fyne\\CanvasObject")),
	}
}

func (m *formAppendMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if f := getForm(classVal); f != nil {
				label := ""
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						label = s.AsString()
					}
				}
				if v, ok := ctx.GetIndexValue(1); ok {
					if wCV, ok := v.(*data.ClassValue); ok {
						if obj := getFyneObject(wCV); obj != nil {
							f.Append(label, obj)
						}
					}
				}
			}
		}
	}
	return nil, nil
}

type formOnSubmitMethod struct{}

func (m *formOnSubmitMethod) GetName() string            { return "onSubmit" }
func (m *formOnSubmitMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *formOnSubmitMethod) GetIsStatic() bool          { return false }
func (m *formOnSubmitMethod) GetReturnType() data.Types  { return nil }
func (m *formOnSubmitMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
	}
}
func (m *formOnSubmitMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, nil),
	}
}
func (m *formOnSubmitMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if f := getForm(classVal); f != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					callback := v
					f.OnSubmit = func() {
						callback.GetValue(ctx)
					}
				}
			}
		}
	}
	return nil, nil
}

type formOnCancelMethod struct{}

func (m *formOnCancelMethod) GetName() string            { return "onCancel" }
func (m *formOnCancelMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *formOnCancelMethod) GetIsStatic() bool          { return false }
func (m *formOnCancelMethod) GetReturnType() data.Types  { return nil }
func (m *formOnCancelMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
	}
}
func (m *formOnCancelMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, nil),
	}
}
func (m *formOnCancelMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if f := getForm(classVal); f != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					callback := v
					f.OnCancel = func() {
						callback.GetValue(ctx)
					}
				}
			}
		}
	}
	return nil, nil
}
