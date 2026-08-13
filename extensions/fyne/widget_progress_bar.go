package fyne

import (
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ProgressBarClass 是 Fyne\Widget\ProgressBar 类
type ProgressBarClass struct {
	*CanvasObjectClass
}

func NewProgressBarClass() data.ClassStmt {
	return &ProgressBarClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\ProgressBar", nil),
	}
}

func (c *ProgressBarClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *ProgressBarClass) GetConstruct() data.Method { return &progressBarConstruct{} }

func (c *ProgressBarClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setValue":
		return &progressBarSetValueMethod{}, true
	case "getValue":
		return &progressBarGetValueMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *ProgressBarClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&progressBarSetValueMethod{},
		&progressBarGetValueMethod{},
	)
}

type progressBarConstruct struct{}

func (m *progressBarConstruct) GetName() string               { return token.ConstructName }
func (m *progressBarConstruct) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *progressBarConstruct) GetIsStatic() bool             { return false }
func (m *progressBarConstruct) GetReturnType() data.Types     { return nil }
func (m *progressBarConstruct) GetParams() []data.GetValue    { return nil }
func (m *progressBarConstruct) GetVariables() []data.Variable { return nil }

func (m *progressBarConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	pb := widget.NewProgressBar()
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			setFyneObject(classVal, pb)
			classVal.SetProperty("_progressBar", data.NewAnyValue(pb))
		}
	}
	return nil, nil
}

func getProgressBar(cv *data.ClassValue) *widget.ProgressBar {
	if v, _ := cv.GetProperty("_progressBar"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if pb, ok := av.Value.(*widget.ProgressBar); ok {
				return pb
			}
		}
	}
	return nil
}

type progressBarSetValueMethod struct{}

func (m *progressBarSetValueMethod) GetName() string            { return "setValue" }
func (m *progressBarSetValueMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *progressBarSetValueMethod) GetIsStatic() bool          { return false }
func (m *progressBarSetValueMethod) GetReturnType() data.Types  { return nil }
func (m *progressBarSetValueMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, data.NewBaseType("float")),
	}
}
func (m *progressBarSetValueMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("float")),
	}
}
func (m *progressBarSetValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if pb := getProgressBar(classVal); pb != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if f, ok := v.(data.AsFloat); ok {
						val, _ := f.AsFloat()
						pb.SetValue(val)
					}
				}
			}
		}
	}
	return nil, nil
}

type progressBarGetValueMethod struct{}

func (m *progressBarGetValueMethod) GetName() string               { return "getValue" }
func (m *progressBarGetValueMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *progressBarGetValueMethod) GetIsStatic() bool             { return false }
func (m *progressBarGetValueMethod) GetReturnType() data.Types     { return data.NewBaseType("float") }
func (m *progressBarGetValueMethod) GetParams() []data.GetValue    { return nil }
func (m *progressBarGetValueMethod) GetVariables() []data.Variable { return nil }
func (m *progressBarGetValueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			if pb := getProgressBar(classVal); pb != nil {
				return data.NewFloatValue(pb.Value), nil
			}
		}
	}
	return data.NewFloatValue(0), nil
}
