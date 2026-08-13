package wails

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// ============================================================================
// Wails\Runtime\Window — 窗口操作 (静态方法)
//
// 所有方法委托给全局 wailsMainWindow (*application.WebviewWindow)
// ============================================================================

type RuntimeWindowClass struct{}

func NewRuntimeWindowClass() data.ClassStmt { return &RuntimeWindowClass{} }

func (c *RuntimeWindowClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *RuntimeWindowClass) GetFrom() data.From                            { return nil }
func (c *RuntimeWindowClass) GetName() string                               { return "Wails\\Runtime\\Window" }
func (c *RuntimeWindowClass) GetExtend() *string                            { return nil }
func (c *RuntimeWindowClass) GetImplements() []string                       { return nil }
func (c *RuntimeWindowClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *RuntimeWindowClass) GetPropertyList() []data.Property              { return nil }
func (c *RuntimeWindowClass) GetConstruct() data.Method                     { return nil }

var rwMethods = map[string]data.Method{
	"setTitle":            &rwSetTitleMethod{},
	"show":                &rwShowMethod{},
	"hide":                &rwHideMethod{},
	"center":              &rwCenterMethod{},
	"maximise":            &rwMaximiseMethod{},
	"unmaximise":          &rwUnmaximiseMethod{},
	"toggleMaximise":      &rwToggleMaximiseMethod{},
	"minimise":            &rwMinimiseMethod{},
	"unminimise":          &rwUnminimiseMethod{},
	"setSize":             &rwSetSizeMethod{},
	"setMinSize":          &rwSetMinSizeMethod{},
	"setMaxSize":          &rwSetMaxSizeMethod{},
	"setPosition":         &rwSetPositionMethod{},
	"getPosition":         &rwGetPositionMethod{},
	"getSize":             &rwGetSizeMethod{},
	"fullscreen":          &rwFullscreenMethod{},
	"unfullscreen":        &rwUnfullscreenMethod{},
	"close":               &rwCloseMethod{},
	"setBackgroundColour": &rwSetBackgroundColourMethod{},
	"setAlwaysOnTop":      &rwSetAlwaysOnTopMethod{},
	"reload":              &rwReloadMethod{},
	"reloadApp":           &rwReloadAppMethod{},
	"execJS":              &rwExecJSMethod{},
	"isMaximised":         &rwIsMaximisedMethod{},
	"isMinimised":         &rwIsMinimisedMethod{},
	"isFullscreen":        &rwIsFullscreenMethod{},
}

func (c *RuntimeWindowClass) GetMethod(name string) (data.Method, bool) {
	m, ok := rwMethods[name]
	return m, ok
}
func (c *RuntimeWindowClass) GetStaticMethod(name string) (data.Method, bool) {
	m, ok := rwMethods[name]
	return m, ok
}
func (c *RuntimeWindowClass) GetMethods() []data.Method {
	methods := make([]data.Method, 0, len(rwMethods))
	for _, m := range rwMethods {
		methods = append(methods, m)
	}
	return methods
}

// ====== setTitle ======
type rwSetTitleMethod struct{}

func (m *rwSetTitleMethod) GetName() string            { return "setTitle" }
func (m *rwSetTitleMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetTitleMethod) GetIsStatic() bool          { return true }
func (m *rwSetTitleMethod) GetReturnType() data.Types  { return nil }
func (m *rwSetTitleMethod) GetParams() []data.GetValue {
	return []data.GetValue{data.NewParameter("title", 0)}
}
func (m *rwSetTitleMethod) GetVariables() []data.Variable {
	return []data.Variable{data.NewVariable("title", 0, data.NewBaseType("string"))}
}
func (m *rwSetTitleMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		if v, ok := ctx.GetIndexValue(0); ok {
			wailsMainWindow.SetTitle(toString(v))
		}
	}
	return nil, nil
}

// ====== show ======
type rwShowMethod struct{}

func (m *rwShowMethod) GetName() string               { return "show" }
func (m *rwShowMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwShowMethod) GetIsStatic() bool             { return true }
func (m *rwShowMethod) GetReturnType() data.Types     { return nil }
func (m *rwShowMethod) GetParams() []data.GetValue    { return nil }
func (m *rwShowMethod) GetVariables() []data.Variable { return nil }
func (m *rwShowMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.Show()
	}
	return nil, nil
}

// ====== hide ======
type rwHideMethod struct{}

func (m *rwHideMethod) GetName() string               { return "hide" }
func (m *rwHideMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwHideMethod) GetIsStatic() bool             { return true }
func (m *rwHideMethod) GetReturnType() data.Types     { return nil }
func (m *rwHideMethod) GetParams() []data.GetValue    { return nil }
func (m *rwHideMethod) GetVariables() []data.Variable { return nil }
func (m *rwHideMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.Hide()
	}
	return nil, nil
}

// ====== center ======
type rwCenterMethod struct{}

func (m *rwCenterMethod) GetName() string               { return "center" }
func (m *rwCenterMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwCenterMethod) GetIsStatic() bool             { return true }
func (m *rwCenterMethod) GetReturnType() data.Types     { return nil }
func (m *rwCenterMethod) GetParams() []data.GetValue    { return nil }
func (m *rwCenterMethod) GetVariables() []data.Variable { return nil }
func (m *rwCenterMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.Center()
	}
	return nil, nil
}

// ====== maximise ======
type rwMaximiseMethod struct{}

func (m *rwMaximiseMethod) GetName() string               { return "maximise" }
func (m *rwMaximiseMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwMaximiseMethod) GetIsStatic() bool             { return true }
func (m *rwMaximiseMethod) GetReturnType() data.Types     { return nil }
func (m *rwMaximiseMethod) GetParams() []data.GetValue    { return nil }
func (m *rwMaximiseMethod) GetVariables() []data.Variable { return nil }
func (m *rwMaximiseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.Maximise()
	}
	return nil, nil
}

// ====== unmaximise ======
type rwUnmaximiseMethod struct{}

func (m *rwUnmaximiseMethod) GetName() string               { return "unmaximise" }
func (m *rwUnmaximiseMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwUnmaximiseMethod) GetIsStatic() bool             { return true }
func (m *rwUnmaximiseMethod) GetReturnType() data.Types     { return nil }
func (m *rwUnmaximiseMethod) GetParams() []data.GetValue    { return nil }
func (m *rwUnmaximiseMethod) GetVariables() []data.Variable { return nil }
func (m *rwUnmaximiseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.UnMaximise()
	}
	return nil, nil
}

// ====== toggleMaximise ======
type rwToggleMaximiseMethod struct{}

func (m *rwToggleMaximiseMethod) GetName() string               { return "toggleMaximise" }
func (m *rwToggleMaximiseMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwToggleMaximiseMethod) GetIsStatic() bool             { return true }
func (m *rwToggleMaximiseMethod) GetReturnType() data.Types     { return nil }
func (m *rwToggleMaximiseMethod) GetParams() []data.GetValue    { return nil }
func (m *rwToggleMaximiseMethod) GetVariables() []data.Variable { return nil }
func (m *rwToggleMaximiseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.ToggleMaximise()
	}
	return nil, nil
}

// ====== minimise ======
type rwMinimiseMethod struct{}

func (m *rwMinimiseMethod) GetName() string               { return "minimise" }
func (m *rwMinimiseMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwMinimiseMethod) GetIsStatic() bool             { return true }
func (m *rwMinimiseMethod) GetReturnType() data.Types     { return nil }
func (m *rwMinimiseMethod) GetParams() []data.GetValue    { return nil }
func (m *rwMinimiseMethod) GetVariables() []data.Variable { return nil }
func (m *rwMinimiseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.Minimise()
	}
	return nil, nil
}

// ====== unminimise ======
type rwUnminimiseMethod struct{}

func (m *rwUnminimiseMethod) GetName() string               { return "unminimise" }
func (m *rwUnminimiseMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwUnminimiseMethod) GetIsStatic() bool             { return true }
func (m *rwUnminimiseMethod) GetReturnType() data.Types     { return nil }
func (m *rwUnminimiseMethod) GetParams() []data.GetValue    { return nil }
func (m *rwUnminimiseMethod) GetVariables() []data.Variable { return nil }
func (m *rwUnminimiseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.UnMinimise()
	}
	return nil, nil
}

// ====== setSize ======
type rwSetSizeMethod struct{}

func (m *rwSetSizeMethod) GetName() string            { return "setSize" }
func (m *rwSetSizeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetSizeMethod) GetIsStatic() bool          { return true }
func (m *rwSetSizeMethod) GetReturnType() data.Types  { return nil }
func (m *rwSetSizeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "width", 0, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "height", 1, data.NewIntValue(0), data.NewBaseType("int")),
	}
}
func (m *rwSetSizeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "width", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "height", 1, data.NewBaseType("int")),
	}
}
func (m *rwSetSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		var w, h int
		if v, ok := ctx.GetIndexValue(0); ok {
			w = toInt(v)
		}
		if v, ok := ctx.GetIndexValue(1); ok {
			h = toInt(v)
		}
		wailsMainWindow.SetSize(w, h)
	}
	return nil, nil
}

// ====== setMinSize ======
type rwSetMinSizeMethod struct{}

func (m *rwSetMinSizeMethod) GetName() string            { return "setMinSize" }
func (m *rwSetMinSizeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetMinSizeMethod) GetIsStatic() bool          { return true }
func (m *rwSetMinSizeMethod) GetReturnType() data.Types  { return nil }
func (m *rwSetMinSizeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "width", 0, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "height", 1, data.NewIntValue(0), data.NewBaseType("int")),
	}
}
func (m *rwSetMinSizeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "width", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "height", 1, data.NewBaseType("int")),
	}
}
func (m *rwSetMinSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		var w, h int
		if v, ok := ctx.GetIndexValue(0); ok {
			w = toInt(v)
		}
		if v, ok := ctx.GetIndexValue(1); ok {
			h = toInt(v)
		}
		wailsMainWindow.SetMinSize(w, h)
	}
	return nil, nil
}

// ====== setMaxSize ======
type rwSetMaxSizeMethod struct{}

func (m *rwSetMaxSizeMethod) GetName() string            { return "setMaxSize" }
func (m *rwSetMaxSizeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetMaxSizeMethod) GetIsStatic() bool          { return true }
func (m *rwSetMaxSizeMethod) GetReturnType() data.Types  { return nil }
func (m *rwSetMaxSizeMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "width", 0, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "height", 1, data.NewIntValue(0), data.NewBaseType("int")),
	}
}
func (m *rwSetMaxSizeMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "width", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "height", 1, data.NewBaseType("int")),
	}
}
func (m *rwSetMaxSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		var w, h int
		if v, ok := ctx.GetIndexValue(0); ok {
			w = toInt(v)
		}
		if v, ok := ctx.GetIndexValue(1); ok {
			h = toInt(v)
		}
		wailsMainWindow.SetMaxSize(w, h)
	}
	return nil, nil
}

// ====== setPosition ======
type rwSetPositionMethod struct{}

func (m *rwSetPositionMethod) GetName() string            { return "setPosition" }
func (m *rwSetPositionMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetPositionMethod) GetIsStatic() bool          { return true }
func (m *rwSetPositionMethod) GetReturnType() data.Types  { return nil }
func (m *rwSetPositionMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "x", 0, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "y", 1, data.NewIntValue(0), data.NewBaseType("int")),
	}
}
func (m *rwSetPositionMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "x", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "y", 1, data.NewBaseType("int")),
	}
}
func (m *rwSetPositionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		var x, y int
		if v, ok := ctx.GetIndexValue(0); ok {
			x = toInt(v)
		}
		if v, ok := ctx.GetIndexValue(1); ok {
			y = toInt(v)
		}
		wailsMainWindow.SetPosition(x, y)
	}
	return nil, nil
}

// ====== getPosition ======
type rwGetPositionMethod struct{}

func (m *rwGetPositionMethod) GetName() string               { return "getPosition" }
func (m *rwGetPositionMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwGetPositionMethod) GetIsStatic() bool             { return true }
func (m *rwGetPositionMethod) GetReturnType() data.Types     { return data.NewBaseType("array") }
func (m *rwGetPositionMethod) GetParams() []data.GetValue    { return nil }
func (m *rwGetPositionMethod) GetVariables() []data.Variable { return nil }
func (m *rwGetPositionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		x, y := wailsMainWindow.Position()
		return data.NewArrayValue([]data.Value{data.NewIntValue(x), data.NewIntValue(y)}), nil
	}
	return data.NewArrayValue([]data.Value{data.NewIntValue(0), data.NewIntValue(0)}), nil
}

// ====== getSize ======
type rwGetSizeMethod struct{}

func (m *rwGetSizeMethod) GetName() string               { return "getSize" }
func (m *rwGetSizeMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwGetSizeMethod) GetIsStatic() bool             { return true }
func (m *rwGetSizeMethod) GetReturnType() data.Types     { return data.NewBaseType("array") }
func (m *rwGetSizeMethod) GetParams() []data.GetValue    { return nil }
func (m *rwGetSizeMethod) GetVariables() []data.Variable { return nil }
func (m *rwGetSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		w, h := wailsMainWindow.Size()
		return data.NewArrayValue([]data.Value{data.NewIntValue(w), data.NewIntValue(h)}), nil
	}
	return data.NewArrayValue([]data.Value{data.NewIntValue(0), data.NewIntValue(0)}), nil
}

// ====== fullscreen ======
type rwFullscreenMethod struct{}

func (m *rwFullscreenMethod) GetName() string               { return "fullscreen" }
func (m *rwFullscreenMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwFullscreenMethod) GetIsStatic() bool             { return true }
func (m *rwFullscreenMethod) GetReturnType() data.Types     { return nil }
func (m *rwFullscreenMethod) GetParams() []data.GetValue    { return nil }
func (m *rwFullscreenMethod) GetVariables() []data.Variable { return nil }
func (m *rwFullscreenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.Fullscreen()
	}
	return nil, nil
}

// ====== unfullscreen ======
type rwUnfullscreenMethod struct{}

func (m *rwUnfullscreenMethod) GetName() string               { return "unfullscreen" }
func (m *rwUnfullscreenMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwUnfullscreenMethod) GetIsStatic() bool             { return true }
func (m *rwUnfullscreenMethod) GetReturnType() data.Types     { return nil }
func (m *rwUnfullscreenMethod) GetParams() []data.GetValue    { return nil }
func (m *rwUnfullscreenMethod) GetVariables() []data.Variable { return nil }
func (m *rwUnfullscreenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.UnFullscreen()
	}
	return nil, nil
}

// ====== close ======
type rwCloseMethod struct{}

func (m *rwCloseMethod) GetName() string               { return "close" }
func (m *rwCloseMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwCloseMethod) GetIsStatic() bool             { return true }
func (m *rwCloseMethod) GetReturnType() data.Types     { return nil }
func (m *rwCloseMethod) GetParams() []data.GetValue    { return nil }
func (m *rwCloseMethod) GetVariables() []data.Variable { return nil }
func (m *rwCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsApp != nil {
		wailsApp.Quit()
	}
	return nil, nil
}

// ====== setBackgroundColour ======
type rwSetBackgroundColourMethod struct{}

func (m *rwSetBackgroundColourMethod) GetName() string            { return "setBackgroundColour" }
func (m *rwSetBackgroundColourMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetBackgroundColourMethod) GetIsStatic() bool          { return true }
func (m *rwSetBackgroundColourMethod) GetReturnType() data.Types  { return nil }
func (m *rwSetBackgroundColourMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "r", 0, data.NewIntValue(255), data.NewBaseType("int")),
		node.NewParameter(nil, "g", 1, data.NewIntValue(255), data.NewBaseType("int")),
		node.NewParameter(nil, "b", 2, data.NewIntValue(255), data.NewBaseType("int")),
		node.NewParameter(nil, "a", 3, data.NewIntValue(255), data.NewBaseType("int")),
	}
}
func (m *rwSetBackgroundColourMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "r", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "g", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "b", 2, data.NewBaseType("int")),
		node.NewVariable(nil, "a", 3, data.NewBaseType("int")),
	}
}
func (m *rwSetBackgroundColourMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		r, g, b, a := 255, 255, 255, 255
		if v, ok := ctx.GetIndexValue(0); ok {
			r = toInt(v)
		}
		if v, ok := ctx.GetIndexValue(1); ok {
			g = toInt(v)
		}
		if v, ok := ctx.GetIndexValue(2); ok {
			b = toInt(v)
		}
		if v, ok := ctx.GetIndexValue(3); ok {
			a = toInt(v)
		}
		wailsMainWindow.SetBackgroundColour(application.RGBA{
			Red:   uint8(r),
			Green: uint8(g),
			Blue:  uint8(b),
			Alpha: uint8(a),
		})
	}
	return nil, nil
}

// ====== setAlwaysOnTop ======
type rwSetAlwaysOnTopMethod struct{}

func (m *rwSetAlwaysOnTopMethod) GetName() string            { return "setAlwaysOnTop" }
func (m *rwSetAlwaysOnTopMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetAlwaysOnTopMethod) GetIsStatic() bool          { return true }
func (m *rwSetAlwaysOnTopMethod) GetReturnType() data.Types  { return nil }
func (m *rwSetAlwaysOnTopMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "alwaysOnTop", 0, data.NewBoolValue(false), data.NewBaseType("bool"))}
}
func (m *rwSetAlwaysOnTopMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "alwaysOnTop", 0, data.NewBaseType("bool"))}
}
func (m *rwSetAlwaysOnTopMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		if v, ok := ctx.GetIndexValue(0); ok {
			wailsMainWindow.SetAlwaysOnTop(toBool(v))
		}
	}
	return nil, nil
}

// ====== reload ======
type rwReloadMethod struct{}

func (m *rwReloadMethod) GetName() string               { return "reload" }
func (m *rwReloadMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwReloadMethod) GetIsStatic() bool             { return true }
func (m *rwReloadMethod) GetReturnType() data.Types     { return nil }
func (m *rwReloadMethod) GetParams() []data.GetValue    { return nil }
func (m *rwReloadMethod) GetVariables() []data.Variable { return nil }
func (m *rwReloadMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.Reload()
	}
	return nil, nil
}

// ====== reloadApp ======
type rwReloadAppMethod struct{}

func (m *rwReloadAppMethod) GetName() string               { return "reloadApp" }
func (m *rwReloadAppMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwReloadAppMethod) GetIsStatic() bool             { return true }
func (m *rwReloadAppMethod) GetReturnType() data.Types     { return nil }
func (m *rwReloadAppMethod) GetParams() []data.GetValue    { return nil }
func (m *rwReloadAppMethod) GetVariables() []data.Variable { return nil }
func (m *rwReloadAppMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		wailsMainWindow.ForceReload()
	}
	return nil, nil
}

// ====== execJS ======
type rwExecJSMethod struct{}

func (m *rwExecJSMethod) GetName() string            { return "execJS" }
func (m *rwExecJSMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwExecJSMethod) GetIsStatic() bool          { return true }
func (m *rwExecJSMethod) GetReturnType() data.Types  { return nil }
func (m *rwExecJSMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "js", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *rwExecJSMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "js", 0, data.NewBaseType("string"))}
}
func (m *rwExecJSMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		if v, ok := ctx.GetIndexValue(0); ok {
			wailsMainWindow.ExecJS(toString(v))
		}
	}
	return nil, nil
}

// ====== isMaximised ======
type rwIsMaximisedMethod struct{}

func (m *rwIsMaximisedMethod) GetName() string               { return "isMaximised" }
func (m *rwIsMaximisedMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwIsMaximisedMethod) GetIsStatic() bool             { return true }
func (m *rwIsMaximisedMethod) GetReturnType() data.Types     { return data.NewBaseType("bool") }
func (m *rwIsMaximisedMethod) GetParams() []data.GetValue    { return nil }
func (m *rwIsMaximisedMethod) GetVariables() []data.Variable { return nil }
func (m *rwIsMaximisedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		return data.NewBoolValue(wailsMainWindow.IsMaximised()), nil
	}
	return data.NewBoolValue(false), nil
}

// ====== isMinimised ======
type rwIsMinimisedMethod struct{}

func (m *rwIsMinimisedMethod) GetName() string               { return "isMinimised" }
func (m *rwIsMinimisedMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwIsMinimisedMethod) GetIsStatic() bool             { return true }
func (m *rwIsMinimisedMethod) GetReturnType() data.Types     { return data.NewBaseType("bool") }
func (m *rwIsMinimisedMethod) GetParams() []data.GetValue    { return nil }
func (m *rwIsMinimisedMethod) GetVariables() []data.Variable { return nil }
func (m *rwIsMinimisedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		return data.NewBoolValue(wailsMainWindow.IsMinimised()), nil
	}
	return data.NewBoolValue(false), nil
}

// ====== isFullscreen ======
type rwIsFullscreenMethod struct{}

func (m *rwIsFullscreenMethod) GetName() string               { return "isFullscreen" }
func (m *rwIsFullscreenMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *rwIsFullscreenMethod) GetIsStatic() bool             { return true }
func (m *rwIsFullscreenMethod) GetReturnType() data.Types     { return data.NewBaseType("bool") }
func (m *rwIsFullscreenMethod) GetParams() []data.GetValue    { return nil }
func (m *rwIsFullscreenMethod) GetVariables() []data.Variable { return nil }
func (m *rwIsFullscreenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsMainWindow != nil {
		return data.NewBoolValue(wailsMainWindow.IsFullscreen()), nil
	}
	return data.NewBoolValue(false), nil
}
