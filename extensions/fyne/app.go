package fyne

import (
	"errors"

	fyneLib "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
	"github.com/php-any/origami/utils"
)

// AppClass 是 Fyne\App 类
type AppClass struct {
	inner fyneLib.App
}

func NewAppClass() data.ClassStmt {
	return &AppClass{}
}

func (c *AppClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *AppClass) GetFrom() data.From                              { return nil }
func (c *AppClass) GetName() string                                 { return "Fyne\\App" }
func (c *AppClass) GetExtend() *string                              { return nil }
func (c *AppClass) GetImplements() []string                         { return nil }
func (c *AppClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *AppClass) GetPropertyList() []data.Property                { return nil }
func (c *AppClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *AppClass) GetMethods() []data.Method {
	return []data.Method{
		&appNewWindowMethod{},
		&appRunMethod{},
		&appQuitMethod{},
	}
}
func (c *AppClass) GetConstruct() data.Method { return &appConstruct{} }

func (c *AppClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "newWindow":
		return &appNewWindowMethod{}, true
	case "run":
		return &appRunMethod{}, true
	case "quit":
		return &appQuitMethod{}, true
	default:
		return nil, false
	}
}

// getApp 从 ClassValue 中获取底层 fyne.App
func getApp(cv *data.ClassValue) fyneLib.App {
	if v, _ := cv.GetProperty("_app"); v != nil {
		if av, ok := v.(*data.AnyValue); ok {
			if a, ok := av.Value.(fyneLib.App); ok {
				return a
			}
		}
	}
	return nil
}

// ====== 构造函数 ======

type appConstruct struct{}

func (m *appConstruct) GetName() string            { return token.ConstructName }
func (m *appConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *appConstruct) GetIsStatic() bool          { return false }
func (m *appConstruct) GetReturnType() data.Types  { return nil }

func (m *appConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "id", 0, data.NewStringValue(""), data.NewBaseType("string")),
	}
}

func (m *appConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "id", 0, data.NewBaseType("string")),
	}
}

func (m *appConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	id := ""
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			id = s.AsString()
		}
	}
	var a fyneLib.App
	if id != "" {
		a = app.NewWithID(id)
	} else {
		a = app.New()
	}
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			classVal.SetProperty("_app", data.NewAnyValue(a))
		}
	}
	return nil, nil
}

// ====== newWindow 方法 ======

type appNewWindowMethod struct{}

func (m *appNewWindowMethod) GetName() string            { return "newWindow" }
func (m *appNewWindowMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *appNewWindowMethod) GetIsStatic() bool          { return false }
func (m *appNewWindowMethod) GetReturnType() data.Types  { return data.NewBaseType("Fyne\\Window") }

func (m *appNewWindowMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "title", 0, nil, data.NewBaseType("string")),
	}
}

func (m *appNewWindowMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "title", 0, data.NewBaseType("string")),
	}
}

func (m *appNewWindowMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	title := ""
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			title = s.AsString()
		}
	}
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			a := getApp(classVal)
			if a == nil {
				return nil, utils.NewThrow(errors.New("App not initialized"))
			}
			w := a.NewWindow(title)
			return NewWindowValue(w, ctx), nil
		}
	}
	return nil, nil
}

// ====== run 方法 ======

type appRunMethod struct{}

func (m *appRunMethod) GetName() string               { return "run" }
func (m *appRunMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *appRunMethod) GetIsStatic() bool             { return false }
func (m *appRunMethod) GetReturnType() data.Types     { return nil }
func (m *appRunMethod) GetParams() []data.GetValue    { return nil }
func (m *appRunMethod) GetVariables() []data.Variable { return nil }

func (m *appRunMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			a := getApp(classVal)
			if a != nil {
				a.Run()
			}
		}
	}
	return nil, nil
}

// ====== quit 方法 ======

type appQuitMethod struct{}

func (m *appQuitMethod) GetName() string               { return "quit" }
func (m *appQuitMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *appQuitMethod) GetIsStatic() bool             { return false }
func (m *appQuitMethod) GetReturnType() data.Types     { return nil }
func (m *appQuitMethod) GetParams() []data.GetValue    { return nil }
func (m *appQuitMethod) GetVariables() []data.Variable { return nil }

func (m *appQuitMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal, ok := cv.GetThis().(*data.ClassValue); ok {
			a := getApp(classVal)
			if a != nil {
				a.Quit()
			}
		}
	}
	return nil, nil
}
