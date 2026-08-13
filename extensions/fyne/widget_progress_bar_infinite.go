package fyne

import (
	"fyne.io/fyne/v2/widget"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/token"
)

// ProgressBarInfiniteClass 是 Fyne\Widget\ProgressBarInfinite 类
type ProgressBarInfiniteClass struct {
	*CanvasObjectClass
}

func NewProgressBarInfiniteClass() data.ClassStmt {
	return &ProgressBarInfiniteClass{
		CanvasObjectClass: newCanvasObjectClass("Fyne\\Widget\\ProgressBarInfinite", nil),
	}
}

func (c *ProgressBarInfiniteClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *ProgressBarInfiniteClass) GetConstruct() data.Method { return &progressBarInfiniteConstruct{} }

func (c *ProgressBarInfiniteClass) GetMethod(name string) (data.Method, bool) {
	return c.CanvasObjectClass.GetMethod(name)
}

func (c *ProgressBarInfiniteClass) GetMethods() []data.Method {
	return c.CanvasObjectClass.GetMethods()
}

type progressBarInfiniteConstruct struct{}

func (m *progressBarInfiniteConstruct) GetName() string               { return token.ConstructName }
func (m *progressBarInfiniteConstruct) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *progressBarInfiniteConstruct) GetIsStatic() bool             { return false }
func (m *progressBarInfiniteConstruct) GetReturnType() data.Types     { return nil }
func (m *progressBarInfiniteConstruct) GetParams() []data.GetValue    { return nil }
func (m *progressBarInfiniteConstruct) GetVariables() []data.Variable { return nil }

func (m *progressBarInfiniteConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	pb := widget.NewProgressBarInfinite()
	if cv, ok := ctx.(*data.ClassMethodContext); ok {
		if classVal := cv.ClassValue; classVal != nil {
			setFyneObject(classVal, pb)
		}
	}
	return nil, nil
}
