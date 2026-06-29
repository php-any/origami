package wails

import (
	"context"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ============================================================================
// wailsContext — 全局 Wails 上下文存储
// ============================================================================

var wailsCtx context.Context

// SetWailsContext 设置当前 Wails 上下文 (在生命周期回调中调用)
func SetWailsContext(ctx context.Context) {
	wailsCtx = ctx
}

// GetWailsContext 获取当前 Wails 上下文
func GetWailsContext() context.Context {
	return wailsCtx
}

// ============================================================================
// Wails\Runtime\Window — 窗口操作 (静态方法)
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
	"setTitle":              &rwSetTitleMethod{},
	"show":                  &rwShowMethod{},
	"hide":                  &rwHideMethod{},
	"center":                &rwCenterMethod{},
	"maximise":              &rwMaximiseMethod{},
	"unmaximise":            &rwUnmaximiseMethod{},
	"toggleMaximise":        &rwToggleMaximiseMethod{},
	"minimise":              &rwMinimiseMethod{},
	"unminimise":            &rwUnminimiseMethod{},
	"setSize":               &rwSetSizeMethod{},
	"setMinSize":            &rwSetMinSizeMethod{},
	"setMaxSize":            &rwSetMaxSizeMethod{},
	"setPosition":           &rwSetPositionMethod{},
	"getPosition":           &rwGetPositionMethod{},
	"getSize":               &rwGetSizeMethod{},
	"fullscreen":            &rwFullscreenMethod{},
	"unfullscreen":          &rwUnfullscreenMethod{},
	"close":                 &rwCloseMethod{},
	"setBackgroundColour":   &rwSetBackgroundColourMethod{},
	"setAlwaysOnTop":        &rwSetAlwaysOnTopMethod{},
	"reload":                &rwReloadMethod{},
	"reloadApp":             &rwReloadAppMethod{},
	"execJS":                &rwExecJSMethod{},
	"setDarkTheme":          &rwSetDarkThemeMethod{},
	"setLightTheme":         &rwSetLightThemeMethod{},
	"setSystemDefaultTheme": &rwSetSystemDefaultThemeMethod{},
	"isMaximised":           &rwIsMaximisedMethod{},
	"isMinimised":           &rwIsMinimisedMethod{},
	"isFullscreen":          &rwIsFullscreenMethod{},
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

// --- Window method implementations ---

type rwSetTitleMethod struct{}

func (m *rwSetTitleMethod) GetName() string { return "setTitle" }
func (m *rwSetTitleMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetTitleMethod) GetIsStatic() bool { return true }
func (m *rwSetTitleMethod) GetReturnType() data.Types { return nil }
func (m *rwSetTitleMethod) GetParams() []data.GetValue {
	return []data.GetValue{data.NewParameter("title", 0)}
}
func (m *rwSetTitleMethod) GetVariables() []data.Variable {
	return []data.Variable{data.NewVariable("title", 0, data.NewBaseType("string"))}
}
func (m *rwSetTitleMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok {
			runtime.WindowSetTitle(wailsCtx, toString(v))
		}
	}
	return nil, nil
}

type rwShowMethod struct{}

func (m *rwShowMethod) GetName() string { return "show" }
func (m *rwShowMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwShowMethod) GetIsStatic() bool { return true }
func (m *rwShowMethod) GetReturnType() data.Types { return nil }
func (m *rwShowMethod) GetParams() []data.GetValue { return nil }
func (m *rwShowMethod) GetVariables() []data.Variable { return nil }
func (m *rwShowMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowShow(wailsCtx) }
	return nil, nil
}

type rwHideMethod struct{}

func (m *rwHideMethod) GetName() string { return "hide" }
func (m *rwHideMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwHideMethod) GetIsStatic() bool { return true }
func (m *rwHideMethod) GetReturnType() data.Types { return nil }
func (m *rwHideMethod) GetParams() []data.GetValue { return nil }
func (m *rwHideMethod) GetVariables() []data.Variable { return nil }
func (m *rwHideMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowHide(wailsCtx) }
	return nil, nil
}

type rwCenterMethod struct{}

func (m *rwCenterMethod) GetName() string { return "center" }
func (m *rwCenterMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwCenterMethod) GetIsStatic() bool { return true }
func (m *rwCenterMethod) GetReturnType() data.Types { return nil }
func (m *rwCenterMethod) GetParams() []data.GetValue { return nil }
func (m *rwCenterMethod) GetVariables() []data.Variable { return nil }
func (m *rwCenterMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowCenter(wailsCtx) }
	return nil, nil
}

type rwMaximiseMethod struct{}

func (m *rwMaximiseMethod) GetName() string { return "maximise" }
func (m *rwMaximiseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwMaximiseMethod) GetIsStatic() bool { return true }
func (m *rwMaximiseMethod) GetReturnType() data.Types { return nil }
func (m *rwMaximiseMethod) GetParams() []data.GetValue { return nil }
func (m *rwMaximiseMethod) GetVariables() []data.Variable { return nil }
func (m *rwMaximiseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowMaximise(wailsCtx) }
	return nil, nil
}

type rwUnmaximiseMethod struct{}

func (m *rwUnmaximiseMethod) GetName() string { return "unmaximise" }
func (m *rwUnmaximiseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwUnmaximiseMethod) GetIsStatic() bool { return true }
func (m *rwUnmaximiseMethod) GetReturnType() data.Types { return nil }
func (m *rwUnmaximiseMethod) GetParams() []data.GetValue { return nil }
func (m *rwUnmaximiseMethod) GetVariables() []data.Variable { return nil }
func (m *rwUnmaximiseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowUnmaximise(wailsCtx) }
	return nil, nil
}

type rwToggleMaximiseMethod struct{}

func (m *rwToggleMaximiseMethod) GetName() string { return "toggleMaximise" }
func (m *rwToggleMaximiseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwToggleMaximiseMethod) GetIsStatic() bool { return true }
func (m *rwToggleMaximiseMethod) GetReturnType() data.Types { return nil }
func (m *rwToggleMaximiseMethod) GetParams() []data.GetValue { return nil }
func (m *rwToggleMaximiseMethod) GetVariables() []data.Variable { return nil }
func (m *rwToggleMaximiseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowToggleMaximise(wailsCtx) }
	return nil, nil
}

type rwMinimiseMethod struct{}

func (m *rwMinimiseMethod) GetName() string { return "minimise" }
func (m *rwMinimiseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwMinimiseMethod) GetIsStatic() bool { return true }
func (m *rwMinimiseMethod) GetReturnType() data.Types { return nil }
func (m *rwMinimiseMethod) GetParams() []data.GetValue { return nil }
func (m *rwMinimiseMethod) GetVariables() []data.Variable { return nil }
func (m *rwMinimiseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowMinimise(wailsCtx) }
	return nil, nil
}

type rwUnminimiseMethod struct{}

func (m *rwUnminimiseMethod) GetName() string { return "unminimise" }
func (m *rwUnminimiseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwUnminimiseMethod) GetIsStatic() bool { return true }
func (m *rwUnminimiseMethod) GetReturnType() data.Types { return nil }
func (m *rwUnminimiseMethod) GetParams() []data.GetValue { return nil }
func (m *rwUnminimiseMethod) GetVariables() []data.Variable { return nil }
func (m *rwUnminimiseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowUnminimise(wailsCtx) }
	return nil, nil
}

type rwSetSizeMethod struct{}

func (m *rwSetSizeMethod) GetName() string { return "setSize" }
func (m *rwSetSizeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetSizeMethod) GetIsStatic() bool { return true }
func (m *rwSetSizeMethod) GetReturnType() data.Types { return nil }
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
	if wailsCtx != nil {
		var w, h int
		if v, ok := ctx.GetIndexValue(0); ok { w = toInt(v) }
		if v, ok := ctx.GetIndexValue(1); ok { h = toInt(v) }
		runtime.WindowSetSize(wailsCtx, w, h)
	}
	return nil, nil
}

type rwSetMinSizeMethod struct{}

func (m *rwSetMinSizeMethod) GetName() string { return "setMinSize" }
func (m *rwSetMinSizeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetMinSizeMethod) GetIsStatic() bool { return true }
func (m *rwSetMinSizeMethod) GetReturnType() data.Types { return nil }
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
	if wailsCtx != nil {
		var w, h int
		if v, ok := ctx.GetIndexValue(0); ok { w = toInt(v) }
		if v, ok := ctx.GetIndexValue(1); ok { h = toInt(v) }
		runtime.WindowSetMinSize(wailsCtx, w, h)
	}
	return nil, nil
}

type rwSetMaxSizeMethod struct{}

func (m *rwSetMaxSizeMethod) GetName() string { return "setMaxSize" }
func (m *rwSetMaxSizeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetMaxSizeMethod) GetIsStatic() bool { return true }
func (m *rwSetMaxSizeMethod) GetReturnType() data.Types { return nil }
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
	if wailsCtx != nil {
		var w, h int
		if v, ok := ctx.GetIndexValue(0); ok { w = toInt(v) }
		if v, ok := ctx.GetIndexValue(1); ok { h = toInt(v) }
		runtime.WindowSetMaxSize(wailsCtx, w, h)
	}
	return nil, nil
}

type rwSetPositionMethod struct{}

func (m *rwSetPositionMethod) GetName() string { return "setPosition" }
func (m *rwSetPositionMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetPositionMethod) GetIsStatic() bool { return true }
func (m *rwSetPositionMethod) GetReturnType() data.Types { return nil }
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
	if wailsCtx != nil {
		var x, y int
		if v, ok := ctx.GetIndexValue(0); ok { x = toInt(v) }
		if v, ok := ctx.GetIndexValue(1); ok { y = toInt(v) }
		runtime.WindowSetPosition(wailsCtx, x, y)
	}
	return nil, nil
}

type rwGetPositionMethod struct{}

func (m *rwGetPositionMethod) GetName() string { return "getPosition" }
func (m *rwGetPositionMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwGetPositionMethod) GetIsStatic() bool { return true }
func (m *rwGetPositionMethod) GetReturnType() data.Types { return data.NewBaseType("array") }
func (m *rwGetPositionMethod) GetParams() []data.GetValue { return nil }
func (m *rwGetPositionMethod) GetVariables() []data.Variable { return nil }
func (m *rwGetPositionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		x, y := runtime.WindowGetPosition(wailsCtx)
		return data.NewArrayValue([]data.Value{data.NewIntValue(x), data.NewIntValue(y)}), nil
	}
	return data.NewArrayValue([]data.Value{data.NewIntValue(0), data.NewIntValue(0)}), nil
}

type rwGetSizeMethod struct{}

func (m *rwGetSizeMethod) GetName() string { return "getSize" }
func (m *rwGetSizeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwGetSizeMethod) GetIsStatic() bool { return true }
func (m *rwGetSizeMethod) GetReturnType() data.Types { return data.NewBaseType("array") }
func (m *rwGetSizeMethod) GetParams() []data.GetValue { return nil }
func (m *rwGetSizeMethod) GetVariables() []data.Variable { return nil }
func (m *rwGetSizeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		w, h := runtime.WindowGetSize(wailsCtx)
		return data.NewArrayValue([]data.Value{data.NewIntValue(w), data.NewIntValue(h)}), nil
	}
	return data.NewArrayValue([]data.Value{data.NewIntValue(0), data.NewIntValue(0)}), nil
}

type rwFullscreenMethod struct{}

func (m *rwFullscreenMethod) GetName() string { return "fullscreen" }
func (m *rwFullscreenMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwFullscreenMethod) GetIsStatic() bool { return true }
func (m *rwFullscreenMethod) GetReturnType() data.Types { return nil }
func (m *rwFullscreenMethod) GetParams() []data.GetValue { return nil }
func (m *rwFullscreenMethod) GetVariables() []data.Variable { return nil }
func (m *rwFullscreenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowFullscreen(wailsCtx) }
	return nil, nil
}

type rwUnfullscreenMethod struct{}

func (m *rwUnfullscreenMethod) GetName() string { return "unfullscreen" }
func (m *rwUnfullscreenMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwUnfullscreenMethod) GetIsStatic() bool { return true }
func (m *rwUnfullscreenMethod) GetReturnType() data.Types { return nil }
func (m *rwUnfullscreenMethod) GetParams() []data.GetValue { return nil }
func (m *rwUnfullscreenMethod) GetVariables() []data.Variable { return nil }
func (m *rwUnfullscreenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowUnfullscreen(wailsCtx) }
	return nil, nil
}

type rwCloseMethod struct{}

func (m *rwCloseMethod) GetName() string { return "close" }
func (m *rwCloseMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwCloseMethod) GetIsStatic() bool { return true }
func (m *rwCloseMethod) GetReturnType() data.Types { return nil }
func (m *rwCloseMethod) GetParams() []data.GetValue { return nil }
func (m *rwCloseMethod) GetVariables() []data.Variable { return nil }
func (m *rwCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.Quit(wailsCtx) }
	return nil, nil
}

type rwSetBackgroundColourMethod struct{}

func (m *rwSetBackgroundColourMethod) GetName() string { return "setBackgroundColour" }
func (m *rwSetBackgroundColourMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetBackgroundColourMethod) GetIsStatic() bool { return true }
func (m *rwSetBackgroundColourMethod) GetReturnType() data.Types { return nil }
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
	if wailsCtx != nil {
		r, g, b, a := 255, 255, 255, 255
		if v, ok := ctx.GetIndexValue(0); ok { r = toInt(v) }
		if v, ok := ctx.GetIndexValue(1); ok { g = toInt(v) }
		if v, ok := ctx.GetIndexValue(2); ok { b = toInt(v) }
		if v, ok := ctx.GetIndexValue(3); ok { a = toInt(v) }
		runtime.WindowSetBackgroundColour(wailsCtx, uint8(r), uint8(g), uint8(b), uint8(a))
	}
	return nil, nil
}

type rwSetAlwaysOnTopMethod struct{}

func (m *rwSetAlwaysOnTopMethod) GetName() string { return "setAlwaysOnTop" }
func (m *rwSetAlwaysOnTopMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetAlwaysOnTopMethod) GetIsStatic() bool { return true }
func (m *rwSetAlwaysOnTopMethod) GetReturnType() data.Types { return nil }
func (m *rwSetAlwaysOnTopMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "alwaysOnTop", 0, data.NewBoolValue(false), data.NewBaseType("bool")),
	}
}
func (m *rwSetAlwaysOnTopMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "alwaysOnTop", 0, data.NewBaseType("bool")),
	}
}
func (m *rwSetAlwaysOnTopMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok {
			runtime.WindowSetAlwaysOnTop(wailsCtx, toBool(v))
		}
	}
	return nil, nil
}

type rwReloadMethod struct{}

func (m *rwReloadMethod) GetName() string { return "reload" }
func (m *rwReloadMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwReloadMethod) GetIsStatic() bool { return true }
func (m *rwReloadMethod) GetReturnType() data.Types { return nil }
func (m *rwReloadMethod) GetParams() []data.GetValue { return nil }
func (m *rwReloadMethod) GetVariables() []data.Variable { return nil }
func (m *rwReloadMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowReload(wailsCtx) }
	return nil, nil
}

type rwReloadAppMethod struct{}

func (m *rwReloadAppMethod) GetName() string { return "reloadApp" }
func (m *rwReloadAppMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwReloadAppMethod) GetIsStatic() bool { return true }
func (m *rwReloadAppMethod) GetReturnType() data.Types { return nil }
func (m *rwReloadAppMethod) GetParams() []data.GetValue { return nil }
func (m *rwReloadAppMethod) GetVariables() []data.Variable { return nil }
func (m *rwReloadAppMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowReloadApp(wailsCtx) }
	return nil, nil
}

type rwExecJSMethod struct{}

func (m *rwExecJSMethod) GetName() string { return "execJS" }
func (m *rwExecJSMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwExecJSMethod) GetIsStatic() bool { return true }
func (m *rwExecJSMethod) GetReturnType() data.Types { return nil }
func (m *rwExecJSMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "js", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *rwExecJSMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "js", 0, data.NewBaseType("string"))}
}
func (m *rwExecJSMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok {
			runtime.WindowExecJS(wailsCtx, toString(v))
		}
	}
	return nil, nil
}

type rwSetDarkThemeMethod struct{}

func (m *rwSetDarkThemeMethod) GetName() string { return "setDarkTheme" }
func (m *rwSetDarkThemeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetDarkThemeMethod) GetIsStatic() bool { return true }
func (m *rwSetDarkThemeMethod) GetReturnType() data.Types { return nil }
func (m *rwSetDarkThemeMethod) GetParams() []data.GetValue { return nil }
func (m *rwSetDarkThemeMethod) GetVariables() []data.Variable { return nil }
func (m *rwSetDarkThemeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowSetDarkTheme(wailsCtx) }
	return nil, nil
}

type rwSetLightThemeMethod struct{}

func (m *rwSetLightThemeMethod) GetName() string { return "setLightTheme" }
func (m *rwSetLightThemeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetLightThemeMethod) GetIsStatic() bool { return true }
func (m *rwSetLightThemeMethod) GetReturnType() data.Types { return nil }
func (m *rwSetLightThemeMethod) GetParams() []data.GetValue { return nil }
func (m *rwSetLightThemeMethod) GetVariables() []data.Variable { return nil }
func (m *rwSetLightThemeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowSetLightTheme(wailsCtx) }
	return nil, nil
}

type rwSetSystemDefaultThemeMethod struct{}

func (m *rwSetSystemDefaultThemeMethod) GetName() string { return "setSystemDefaultTheme" }
func (m *rwSetSystemDefaultThemeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwSetSystemDefaultThemeMethod) GetIsStatic() bool { return true }
func (m *rwSetSystemDefaultThemeMethod) GetReturnType() data.Types { return nil }
func (m *rwSetSystemDefaultThemeMethod) GetParams() []data.GetValue { return nil }
func (m *rwSetSystemDefaultThemeMethod) GetVariables() []data.Variable { return nil }
func (m *rwSetSystemDefaultThemeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil { runtime.WindowSetSystemDefaultTheme(wailsCtx) }
	return nil, nil
}

type rwIsMaximisedMethod struct{}

func (m *rwIsMaximisedMethod) GetName() string { return "isMaximised" }
func (m *rwIsMaximisedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwIsMaximisedMethod) GetIsStatic() bool { return true }
func (m *rwIsMaximisedMethod) GetReturnType() data.Types { return data.NewBaseType("bool") }
func (m *rwIsMaximisedMethod) GetParams() []data.GetValue { return nil }
func (m *rwIsMaximisedMethod) GetVariables() []data.Variable { return nil }
func (m *rwIsMaximisedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		return data.NewBoolValue(runtime.WindowIsMaximised(wailsCtx)), nil
	}
	return data.NewBoolValue(false), nil
}

type rwIsMinimisedMethod struct{}

func (m *rwIsMinimisedMethod) GetName() string { return "isMinimised" }
func (m *rwIsMinimisedMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwIsMinimisedMethod) GetIsStatic() bool { return true }
func (m *rwIsMinimisedMethod) GetReturnType() data.Types { return data.NewBaseType("bool") }
func (m *rwIsMinimisedMethod) GetParams() []data.GetValue { return nil }
func (m *rwIsMinimisedMethod) GetVariables() []data.Variable { return nil }
func (m *rwIsMinimisedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		return data.NewBoolValue(runtime.WindowIsMinimised(wailsCtx)), nil
	}
	return data.NewBoolValue(false), nil
}

type rwIsFullscreenMethod struct{}

func (m *rwIsFullscreenMethod) GetName() string { return "isFullscreen" }
func (m *rwIsFullscreenMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rwIsFullscreenMethod) GetIsStatic() bool { return true }
func (m *rwIsFullscreenMethod) GetReturnType() data.Types { return data.NewBaseType("bool") }
func (m *rwIsFullscreenMethod) GetParams() []data.GetValue { return nil }
func (m *rwIsFullscreenMethod) GetVariables() []data.Variable { return nil }
func (m *rwIsFullscreenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		return data.NewBoolValue(runtime.WindowIsFullscreen(wailsCtx)), nil
	}
	return data.NewBoolValue(false), nil
}
