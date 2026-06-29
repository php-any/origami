package fyne

import (
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// AccordionClass 是 Fyne\Widget\Accordion 类
type AccordionClass struct {
	*CanvasObjectClass
}

func NewAccordionClass() data.ClassStmt {
	return &AccordionClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\Accordion", nil),
	}
}

func (c *AccordionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *AccordionClass) GetConstruct() data.Method { return &accordionConstruct{} }

func (c *AccordionClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "append":
		return &accordionAppendMethod{}, true
	case "open":
		return &accordionOpenMethod{}, true
	case "close":
		return &accordionCloseMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *AccordionClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&accordionAppendMethod{},
		&accordionOpenMethod{},
		&accordionCloseMethod{},
	)
}

type accordionConstruct struct{}

func (m *accordionConstruct) GetName() string               { return token.ConstructName }
func (m *accordionConstruct) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *accordionConstruct) GetIsStatic() bool             { return false }
func (m *accordionConstruct) GetReturnType() data.Types     { return nil }
func (m *accordionConstruct) GetParams() []data.GetValue    { return nil }
func (m *accordionConstruct) GetVariables() []data.Variable { return nil }

func (m *accordionConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	accordion := widget.NewAccordion()
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			setFyneObject(classVal, accordion)
			classVal.SetProperty("_accordion", data.NewAnyValue(accordion))
		}
	}
	return nil, nil
}

func getAccordion(cv *data.ClassValue) *widget.Accordion {
	if v, _ := cv.GetProperty("_accordion"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if a, ok := av.Value.(*widget.Accordion); ok {
				return a
			}
		}
	}
	return nil
}

type accordionAppendMethod struct{}

func (m *accordionAppendMethod) GetName() string            { return "append" }
func (m *accordionAppendMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *accordionAppendMethod) GetIsStatic() bool          { return false }
func (m *accordionAppendMethod) GetReturnType() data.Types  { return nil }
func (m *accordionAppendMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "title", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "content", 1, nil, data.NewBaseType("Fyne\\CanvasObject")),
	}
}
func (m *accordionAppendMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "title", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "content", 1, data.NewBaseType("Fyne\\CanvasObject")),
	}
}

func (m *accordionAppendMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if a := getAccordion(classVal); a != nil {
				title := ""
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						title = s.AsString()
					}
				}
				if v, ok := ctx.GetIndexValue(1); ok {
					if contentCV, ok := v.(*data.ClassValue); ok {
						if obj := getFyneObject(contentCV); obj != nil {
							item := widget.NewAccordionItem(title, obj)
							a.Append(item)
						}
					}
				}
			}
		}
	}
	return nil, nil
}

type accordionOpenMethod struct{}

func (m *accordionOpenMethod) GetName() string            { return "open" }
func (m *accordionOpenMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *accordionOpenMethod) GetIsStatic() bool          { return false }
func (m *accordionOpenMethod) GetReturnType() data.Types  { return nil }
func (m *accordionOpenMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "index", 0, nil, data.NewBaseType("int")),
	}
}
func (m *accordionOpenMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "index", 0, data.NewBaseType("int")),
	}
}
func (m *accordionOpenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if a := getAccordion(classVal); a != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if i, ok := v.(data.AsInt); ok {
						idx, _ := i.AsInt()
						a.Open(idx)
					}
				}
			}
		}
	}
	return nil, nil
}

type accordionCloseMethod struct{}

func (m *accordionCloseMethod) GetName() string            { return "close" }
func (m *accordionCloseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *accordionCloseMethod) GetIsStatic() bool          { return false }
func (m *accordionCloseMethod) GetReturnType() data.Types  { return nil }
func (m *accordionCloseMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "index", 0, nil, data.NewBaseType("int")),
	}
}
func (m *accordionCloseMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "index", 0, data.NewBaseType("int")),
	}
}
func (m *accordionCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if a := getAccordion(classVal); a != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if i, ok := v.(data.AsInt); ok {
						idx, _ := i.AsInt()
						a.Close(idx)
					}
				}
			}
		}
	}
	return nil, nil
}
