package fyne

import (
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// CardClass 是 Fyne\Widget\Card 类
type CardClass struct {
	*CanvasObjectClass
}

func NewCardClass() data.ClassStmt {
	return &CardClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\Card", nil),
	}
}

func (c *CardClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *CardClass) GetConstruct() data.Method { return &cardConstruct{} }

func (c *CardClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "setTitle":
		return &cardSetTitleMethod{}, true
	case "setSubTitle":
		return &cardSetSubTitleMethod{}, true
	case "setContent":
		return &cardSetContentMethod{}, true
	default:
		return c.CanvasObjectClass.GetMethod(name)
	}
}

func (c *CardClass) GetMethods() []data.Method {
	return append(c.CanvasObjectClass.GetMethods(),
		&cardSetTitleMethod{},
		&cardSetSubTitleMethod{},
		&cardSetContentMethod{},
	)
}

type cardConstruct struct{}

func (m *cardConstruct) GetName() string            { return token.ConstructName }
func (m *cardConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *cardConstruct) GetIsStatic() bool          { return false }
func (m *cardConstruct) GetReturnType() data.Types  { return nil }
func (m *cardConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "title", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "subtitle", 1, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "content", 2, data.NewNullValue(), data.NewBaseType("Fyne\\CanvasObject")),
	}
}
func (m *cardConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "title", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "subtitle", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "content", 2, data.NewBaseType("Fyne\\CanvasObject")),
	}
}

func (m *cardConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	title, subtitle := "", ""
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			title = s.AsString()
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		if s, ok := v.(data.AsString); ok {
			subtitle = s.AsString()
		}
	}
	var content fyneLib.CanvasObject
	if v, ok := ctx.GetIndexValue(2); ok {
		if cv, ok := v.(*data.ClassValue); ok {
			content = getFyneObject(cv)
		}
	}
	card := widget.NewCard(title, subtitle, content)
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			setFyneObject(classVal, card)
			classVal.SetProperty("_card", data.NewAnyValue(card))
		}
	}
	return nil, nil
}

func getCard(cv *data.ClassValue) *widget.Card {
	if v, _ := cv.GetProperty("_card"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if c, ok := av.Value.(*widget.Card); ok {
				return c
			}
		}
	}
	return nil
}

type cardSetTitleMethod struct{}

func (m *cardSetTitleMethod) GetName() string            { return "setTitle" }
func (m *cardSetTitleMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *cardSetTitleMethod) GetIsStatic() bool          { return false }
func (m *cardSetTitleMethod) GetReturnType() data.Types  { return nil }
func (m *cardSetTitleMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "title", 0, nil, data.NewBaseType("string")),
	}
}
func (m *cardSetTitleMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "title", 0, data.NewBaseType("string")),
	}
}
func (m *cardSetTitleMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if c := getCard(classVal); c != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						c.SetTitle(s.AsString())
					}
				}
			}
		}
	}
	return nil, nil
}

type cardSetSubTitleMethod struct{}

func (m *cardSetSubTitleMethod) GetName() string            { return "setSubTitle" }
func (m *cardSetSubTitleMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *cardSetSubTitleMethod) GetIsStatic() bool          { return false }
func (m *cardSetSubTitleMethod) GetReturnType() data.Types  { return nil }
func (m *cardSetSubTitleMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "subtitle", 0, nil, data.NewBaseType("string")),
	}
}
func (m *cardSetSubTitleMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "subtitle", 0, data.NewBaseType("string")),
	}
}
func (m *cardSetSubTitleMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if c := getCard(classVal); c != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if s, ok := v.(data.AsString); ok {
						c.SetSubTitle(s.AsString())
					}
				}
			}
		}
	}
	return nil, nil
}

type cardSetContentMethod struct{}

func (m *cardSetContentMethod) GetName() string            { return "setContent" }
func (m *cardSetContentMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *cardSetContentMethod) GetIsStatic() bool          { return false }
func (m *cardSetContentMethod) GetReturnType() data.Types  { return nil }
func (m *cardSetContentMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "content", 0, nil, data.NewBaseType("Fyne\\CanvasObject")),
	}
}
func (m *cardSetContentMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "content", 0, data.NewBaseType("Fyne\\CanvasObject")),
	}
}
func (m *cardSetContentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			if c := getCard(classVal); c != nil {
				if v, ok := ctx.GetIndexValue(0); ok {
					if contentCV, ok := v.(*data.ClassValue); ok {
						if obj := getFyneObject(contentCV); obj != nil {
							c.SetContent(obj)
						}
					}
				}
			}
		}
	}
	return nil, nil
}
