package fyne

import (
	"errors"

	fyneLib "fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// DialogClass 是 Fyne\Dialog 类，提供静态方法
type DialogClass struct{}

func NewDialogClass() data.ClassStmt { return &DialogClass{} }

func (c *DialogClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *DialogClass) GetFrom() data.From                            { return nil }
func (c *DialogClass) GetName() string                               { return "Fyne\\Dialog" }
func (c *DialogClass) GetExtend() *string                            { return nil }
func (c *DialogClass) GetImplements() []string                       { return nil }
func (c *DialogClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *DialogClass) GetPropertyList() []data.Property              { return nil }
func (c *DialogClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *DialogClass) GetMethods() []data.Method                     { return nil }
func (c *DialogClass) GetConstruct() data.Method                     { return nil }

func (c *DialogClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "showInformation":
		return &dialogShowInformationMethod{}, true
	case "showConfirm":
		return &dialogShowConfirmMethod{}, true
	case "showError":
		return &dialogShowErrorMethod{}, true
	case "showCustom":
		return &dialogShowCustomMethod{}, true
	default:
		return nil, false
	}
}

// getFyneWindowFromArg 从参数中获取 fyne.Window
func getFyneWindowFromArg(v data.Value) fyneLib.Window {
	if cv, ok := v.(*data.ClassValue); ok {
		return getFyneWindow(cv)
	}
	return nil
}

// ====== showInformation ======

type dialogShowInformationMethod struct{}

func (m *dialogShowInformationMethod) GetName() string            { return "showInformation" }
func (m *dialogShowInformationMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *dialogShowInformationMethod) GetIsStatic() bool          { return true }
func (m *dialogShowInformationMethod) GetReturnType() data.Types  { return nil }
func (m *dialogShowInformationMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "title", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "message", 1, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "window", 2, nil, data.NewBaseType("Fyne\\Window")),
	}
}
func (m *dialogShowInformationMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "title", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "message", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "window", 2, data.NewBaseType("Fyne\\Window")),
	}
}
func (m *dialogShowInformationMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	title, message := "", ""
	var w fyneLib.Window
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			title = s.AsString()
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		if s, ok := v.(data.AsString); ok {
			message = s.AsString()
		}
	}
	if v, ok := ctx.GetIndexValue(2); ok {
		w = getFyneWindowFromArg(v)
	}
	if w == nil {
		return nil, utils.NewThrow(errors.New("Dialog::showInformation requires a valid Window"))
	}
	dialog.ShowInformation(title, message, w)
	return nil, nil
}

// ====== showConfirm ======

type dialogShowConfirmMethod struct{}

func (m *dialogShowConfirmMethod) GetName() string            { return "showConfirm" }
func (m *dialogShowConfirmMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *dialogShowConfirmMethod) GetIsStatic() bool          { return true }
func (m *dialogShowConfirmMethod) GetReturnType() data.Types  { return nil }
func (m *dialogShowConfirmMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "title", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "message", 1, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "callback", 2, nil, nil),
		node.NewParameter(nil, "window", 3, nil, data.NewBaseType("Fyne\\Window")),
	}
}
func (m *dialogShowConfirmMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "title", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "message", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "callback", 2, nil),
		node.NewVariable(nil, "window", 3, data.NewBaseType("Fyne\\Window")),
	}
}
func (m *dialogShowConfirmMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	title, message := "", ""
	var callback data.FuncStmt
	var w fyneLib.Window
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			title = s.AsString()
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		if s, ok := v.(data.AsString); ok {
			message = s.AsString()
		}
	}
	if v, ok := ctx.GetIndexValue(2); ok {
		if fv, ok := v.(*data.FuncValue); ok {
			callback = fv.Value
		}
	}
	if v, ok := ctx.GetIndexValue(3); ok {
		w = getFyneWindowFromArg(v)
	}
	if w == nil {
		return nil, utils.NewThrow(errors.New("Dialog::showConfirm requires a valid Window"))
	}
	dialog.ShowConfirm(title, message, func(confirmed bool) {
		callPHPCallbackWith(callback, ctx, data.NewBoolValue(confirmed))
	}, w)
	return nil, nil
}

// ====== showError ======

type dialogShowErrorMethod struct{}

func (m *dialogShowErrorMethod) GetName() string            { return "showError" }
func (m *dialogShowErrorMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *dialogShowErrorMethod) GetIsStatic() bool          { return true }
func (m *dialogShowErrorMethod) GetReturnType() data.Types  { return nil }
func (m *dialogShowErrorMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "message", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "window", 1, nil, data.NewBaseType("Fyne\\Window")),
	}
}
func (m *dialogShowErrorMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "message", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "window", 1, data.NewBaseType("Fyne\\Window")),
	}
}
func (m *dialogShowErrorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	message := ""
	var w fyneLib.Window
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			message = s.AsString()
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		w = getFyneWindowFromArg(v)
	}
	if w == nil {
		return nil, utils.NewThrow(errors.New("Dialog::showError requires a valid Window"))
	}
	dialog.ShowError(errors.New(message), w)
	return nil, nil
}

// ====== showCustom ======

type dialogShowCustomMethod struct{}

func (m *dialogShowCustomMethod) GetName() string            { return "showCustom" }
func (m *dialogShowCustomMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *dialogShowCustomMethod) GetIsStatic() bool          { return true }
func (m *dialogShowCustomMethod) GetReturnType() data.Types  { return nil }
func (m *dialogShowCustomMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "title", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "dismiss", 1, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "content", 2, nil, data.NewBaseType("Fyne\\CanvasObject")),
		node.NewParameter(nil, "window", 3, nil, data.NewBaseType("Fyne\\Window")),
	}
}
func (m *dialogShowCustomMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "title", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "dismiss", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "content", 2, data.NewBaseType("Fyne\\CanvasObject")),
		node.NewVariable(nil, "window", 3, data.NewBaseType("Fyne\\Window")),
	}
}
func (m *dialogShowCustomMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	title, dismiss := "", ""
	var content fyneLib.CanvasObject
	var w fyneLib.Window
	if v, ok := ctx.GetIndexValue(0); ok {
		if s, ok := v.(data.AsString); ok {
			title = s.AsString()
		}
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		if s, ok := v.(data.AsString); ok {
			dismiss = s.AsString()
		}
	}
	if v, ok := ctx.GetIndexValue(2); ok {
		if cv, ok := v.(*data.ClassValue); ok {
			content = getFyneObject(cv)
		}
	}
	if v, ok := ctx.GetIndexValue(3); ok {
		w = getFyneWindowFromArg(v)
	}
	if w == nil {
		return nil, utils.NewThrow(errors.New("Dialog::showCustom requires a valid Window"))
	}
	dialog.ShowCustom(title, dismiss, content, w)
	return nil, nil
}
