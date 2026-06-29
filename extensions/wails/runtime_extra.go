package wails

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ============================================================================
// Wails\Runtime\Dialog — 对话框操作 (静态方法)
// ============================================================================

type RuntimeDialogClass struct{}

func NewRuntimeDialogClass() data.ClassStmt { return &RuntimeDialogClass{} }

func (c *RuntimeDialogClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *RuntimeDialogClass) GetFrom() data.From                            { return nil }
func (c *RuntimeDialogClass) GetName() string                               { return "Wails\\Runtime\\Dialog" }
func (c *RuntimeDialogClass) GetExtend() *string                            { return nil }
func (c *RuntimeDialogClass) GetImplements() []string                       { return nil }
func (c *RuntimeDialogClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *RuntimeDialogClass) GetPropertyList() []data.Property              { return nil }
func (c *RuntimeDialogClass) GetConstruct() data.Method                     { return nil }

var rdMethods = map[string]data.Method{
	"openFile":          &rdOpenFileMethod{},
	"openMultipleFiles": &rdOpenMultipleFilesMethod{},
	"openDirectory":     &rdOpenDirectoryMethod{},
	"saveFile":          &rdSaveFileMethod{},
	"message":           &rdMessageMethod{},
}

func (c *RuntimeDialogClass) GetMethod(name string) (data.Method, bool) {
	m, ok := rdMethods[name]; return m, ok
}
func (c *RuntimeDialogClass) GetStaticMethod(name string) (data.Method, bool) {
	m, ok := rdMethods[name]; return m, ok
}
func (c *RuntimeDialogClass) GetMethods() []data.Method {
	methods := make([]data.Method, 0, len(rdMethods))
	for _, m := range rdMethods { methods = append(methods, m) }
	return methods
}

type rdOpenFileMethod struct{}

func (m *rdOpenFileMethod) GetName() string { return "openFile" }
func (m *rdOpenFileMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rdOpenFileMethod) GetIsStatic() bool { return true }
func (m *rdOpenFileMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
func (m *rdOpenFileMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("Wails\\Dialog\\OpenDialogOptions"))}
}
func (m *rdOpenFileMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "options", 0, data.NewBaseType("Wails\\Dialog\\OpenDialogOptions"))}
}
func (m *rdOpenFileMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx == nil { return data.NewStringValue(""), nil }
	var opts runtime.OpenDialogOptions
	if v, ok := ctx.GetIndexValue(0); ok { opts = buildOpenDialogOptions(v) }
	result, err := runtime.OpenFileDialog(wailsCtx, opts)
	if err != nil { return data.NewStringValue(""), nil }
	return data.NewStringValue(result), nil
}

type rdOpenMultipleFilesMethod struct{}

func (m *rdOpenMultipleFilesMethod) GetName() string { return "openMultipleFiles" }
func (m *rdOpenMultipleFilesMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rdOpenMultipleFilesMethod) GetIsStatic() bool { return true }
func (m *rdOpenMultipleFilesMethod) GetReturnType() data.Types { return data.NewBaseType("array") }
func (m *rdOpenMultipleFilesMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("Wails\\Dialog\\OpenDialogOptions"))}
}
func (m *rdOpenMultipleFilesMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "options", 0, data.NewBaseType("Wails\\Dialog\\OpenDialogOptions"))}
}
func (m *rdOpenMultipleFilesMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx == nil { return data.NewArrayValue(nil), nil }
	var opts runtime.OpenDialogOptions
	if v, ok := ctx.GetIndexValue(0); ok { opts = buildOpenDialogOptions(v) }
	results, err := runtime.OpenMultipleFilesDialog(wailsCtx, opts)
	if err != nil { return data.NewArrayValue(nil), nil }
	vals := make([]data.Value, len(results))
	for i, r := range results { vals[i] = data.NewStringValue(r) }
	return data.NewArrayValue(vals), nil
}

type rdOpenDirectoryMethod struct{}

func (m *rdOpenDirectoryMethod) GetName() string { return "openDirectory" }
func (m *rdOpenDirectoryMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rdOpenDirectoryMethod) GetIsStatic() bool { return true }
func (m *rdOpenDirectoryMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
func (m *rdOpenDirectoryMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("Wails\\Dialog\\OpenDialogOptions"))}
}
func (m *rdOpenDirectoryMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "options", 0, data.NewBaseType("Wails\\Dialog\\OpenDialogOptions"))}
}
func (m *rdOpenDirectoryMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx == nil { return data.NewStringValue(""), nil }
	var opts runtime.OpenDialogOptions
	if v, ok := ctx.GetIndexValue(0); ok { opts = buildOpenDialogOptions(v) }
	result, err := runtime.OpenDirectoryDialog(wailsCtx, opts)
	if err != nil { return data.NewStringValue(""), nil }
	return data.NewStringValue(result), nil
}

type rdSaveFileMethod struct{}

func (m *rdSaveFileMethod) GetName() string { return "saveFile" }
func (m *rdSaveFileMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rdSaveFileMethod) GetIsStatic() bool { return true }
func (m *rdSaveFileMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
func (m *rdSaveFileMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("Wails\\Dialog\\SaveDialogOptions"))}
}
func (m *rdSaveFileMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "options", 0, data.NewBaseType("Wails\\Dialog\\SaveDialogOptions"))}
}
func (m *rdSaveFileMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx == nil { return data.NewStringValue(""), nil }
	var opts runtime.SaveDialogOptions
	if v, ok := ctx.GetIndexValue(0); ok { opts = buildSaveDialogOptions(v) }
	result, err := runtime.SaveFileDialog(wailsCtx, opts)
	if err != nil { return data.NewStringValue(""), nil }
	return data.NewStringValue(result), nil
}

type rdMessageMethod struct{}

func (m *rdMessageMethod) GetName() string { return "message" }
func (m *rdMessageMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rdMessageMethod) GetIsStatic() bool { return true }
func (m *rdMessageMethod) GetReturnType() data.Types { return data.NewBaseType("string") }
func (m *rdMessageMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("Wails\\Dialog\\MessageDialogOptions"))}
}
func (m *rdMessageMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "options", 0, data.NewBaseType("Wails\\Dialog\\MessageDialogOptions"))}
}
func (m *rdMessageMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx == nil { return data.NewStringValue(""), nil }
	var opts runtime.MessageDialogOptions
	if v, ok := ctx.GetIndexValue(0); ok { opts = buildMessageDialogOptions(v) }
	result, err := runtime.MessageDialog(wailsCtx, opts)
	if err != nil { return data.NewStringValue(""), nil }
	return data.NewStringValue(result), nil
}

// ============================================================================
// Wails\Runtime\Events — 事件系统 (静态方法)
// ============================================================================

type RuntimeEventsClass struct{}

func NewRuntimeEventsClass() data.ClassStmt { return &RuntimeEventsClass{} }

func (c *RuntimeEventsClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *RuntimeEventsClass) GetFrom() data.From                            { return nil }
func (c *RuntimeEventsClass) GetName() string                               { return "Wails\\Runtime\\Events" }
func (c *RuntimeEventsClass) GetExtend() *string                            { return nil }
func (c *RuntimeEventsClass) GetImplements() []string                       { return nil }
func (c *RuntimeEventsClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *RuntimeEventsClass) GetPropertyList() []data.Property              { return nil }
func (c *RuntimeEventsClass) GetConstruct() data.Method                     { return nil }

var reMethods = map[string]data.Method{
	"on":   &reOnMethod{},
	"once": &reOnceMethod{},
	"emit": &reEmitMethod{},
	"off":  &reOffMethod{},
}

func (c *RuntimeEventsClass) GetMethod(name string) (data.Method, bool) {
	m, ok := reMethods[name]; return m, ok
}
func (c *RuntimeEventsClass) GetStaticMethod(name string) (data.Method, bool) {
	m, ok := reMethods[name]; return m, ok
}
func (c *RuntimeEventsClass) GetMethods() []data.Method {
	return []data.Method{&reOnMethod{}, &reOnceMethod{}, &reEmitMethod{}, &reOffMethod{}}
}

type reOnMethod struct{}

func (m *reOnMethod) GetName() string { return "on" }
func (m *reOnMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *reOnMethod) GetIsStatic() bool { return true }
func (m *reOnMethod) GetReturnType() data.Types { return nil }
func (m *reOnMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "eventName", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "callback", 1, nil, data.NewBaseType("callable")),
	}
}
func (m *reOnMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "eventName", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "callback", 1, data.NewBaseType("callable")),
	}
}
func (m *reOnMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok {
			runtime.EventsOn(wailsCtx, toString(v), func(data ...interface{}) {})
		}
	}
	return nil, nil
}

type reOnceMethod struct{}

func (m *reOnceMethod) GetName() string { return "once" }
func (m *reOnceMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *reOnceMethod) GetIsStatic() bool { return true }
func (m *reOnceMethod) GetReturnType() data.Types { return nil }
func (m *reOnceMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "eventName", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "callback", 1, nil, data.NewBaseType("callable")),
	}
}
func (m *reOnceMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "eventName", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "callback", 1, data.NewBaseType("callable")),
	}
}
func (m *reOnceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok {
			runtime.EventsOnce(wailsCtx, toString(v), func(data ...interface{}) {})
		}
	}
	return nil, nil
}

type reEmitMethod struct{}

func (m *reEmitMethod) GetName() string { return "emit" }
func (m *reEmitMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *reEmitMethod) GetIsStatic() bool { return true }
func (m *reEmitMethod) GetReturnType() data.Types { return nil }
func (m *reEmitMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "eventName", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *reEmitMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "eventName", 0, data.NewBaseType("string"))}
}
func (m *reEmitMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok {
			runtime.EventsEmit(wailsCtx, toString(v))
		}
	}
	return nil, nil
}

type reOffMethod struct{}

func (m *reOffMethod) GetName() string { return "off" }
func (m *reOffMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *reOffMethod) GetIsStatic() bool { return true }
func (m *reOffMethod) GetReturnType() data.Types { return nil }
func (m *reOffMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "eventName", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *reOffMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "eventName", 0, data.NewBaseType("string"))}
}
func (m *reOffMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok {
			runtime.EventsOff(wailsCtx, toString(v))
		}
	}
	return nil, nil
}

// ============================================================================
// Wails\Runtime\Log — 日志操作 (静态方法)
// ============================================================================

type RuntimeLogClass struct{}

func NewRuntimeLogClass() data.ClassStmt { return &RuntimeLogClass{} }

func (c *RuntimeLogClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *RuntimeLogClass) GetFrom() data.From                            { return nil }
func (c *RuntimeLogClass) GetName() string                               { return "Wails\\Runtime\\Log" }
func (c *RuntimeLogClass) GetExtend() *string                            { return nil }
func (c *RuntimeLogClass) GetImplements() []string                       { return nil }
func (c *RuntimeLogClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *RuntimeLogClass) GetPropertyList() []data.Property              { return nil }
func (c *RuntimeLogClass) GetConstruct() data.Method                     { return nil }

var rlMethods = map[string]data.Method{
	"print":   &rlPrintMethod{},
	"trace":   &rlTraceMethod{},
	"debug":   &rlDebugMethod{},
	"info":    &rlInfoMethod{},
	"warning": &rlWarningMethod{},
	"error":   &rlErrorMethod{},
	"fatal":   &rlFatalMethod{},
}

func (c *RuntimeLogClass) GetMethod(name string) (data.Method, bool) {
	m, ok := rlMethods[name]; return m, ok
}
func (c *RuntimeLogClass) GetStaticMethod(name string) (data.Method, bool) {
	m, ok := rlMethods[name]; return m, ok
}
func (c *RuntimeLogClass) GetMethods() []data.Method {
	methods := make([]data.Method, 0, len(rlMethods))
	for _, m := range rlMethods { methods = append(methods, m) }
	return methods
}

type rlPrintMethod struct{}

func (m *rlPrintMethod) GetName() string { return "print" }
func (m *rlPrintMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rlPrintMethod) GetIsStatic() bool { return true }
func (m *rlPrintMethod) GetReturnType() data.Types { return nil }
func (m *rlPrintMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "message", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *rlPrintMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "message", 0, data.NewBaseType("string"))}
}
func (m *rlPrintMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok { runtime.LogPrint(wailsCtx, toString(v)) }
	}
	return nil, nil
}

type rlTraceMethod struct{}

func (m *rlTraceMethod) GetName() string { return "trace" }
func (m *rlTraceMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rlTraceMethod) GetIsStatic() bool { return true }
func (m *rlTraceMethod) GetReturnType() data.Types { return nil }
func (m *rlTraceMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "message", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *rlTraceMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "message", 0, data.NewBaseType("string"))}
}
func (m *rlTraceMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok { runtime.LogTrace(wailsCtx, toString(v)) }
	}
	return nil, nil
}

type rlDebugMethod struct{}

func (m *rlDebugMethod) GetName() string { return "debug" }
func (m *rlDebugMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rlDebugMethod) GetIsStatic() bool { return true }
func (m *rlDebugMethod) GetReturnType() data.Types { return nil }
func (m *rlDebugMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "message", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *rlDebugMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "message", 0, data.NewBaseType("string"))}
}
func (m *rlDebugMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok { runtime.LogDebug(wailsCtx, toString(v)) }
	}
	return nil, nil
}

type rlInfoMethod struct{}

func (m *rlInfoMethod) GetName() string { return "info" }
func (m *rlInfoMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rlInfoMethod) GetIsStatic() bool { return true }
func (m *rlInfoMethod) GetReturnType() data.Types { return nil }
func (m *rlInfoMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "message", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *rlInfoMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "message", 0, data.NewBaseType("string"))}
}
func (m *rlInfoMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok { runtime.LogInfo(wailsCtx, toString(v)) }
	}
	return nil, nil
}

type rlWarningMethod struct{}

func (m *rlWarningMethod) GetName() string { return "warning" }
func (m *rlWarningMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rlWarningMethod) GetIsStatic() bool { return true }
func (m *rlWarningMethod) GetReturnType() data.Types { return nil }
func (m *rlWarningMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "message", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *rlWarningMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "message", 0, data.NewBaseType("string"))}
}
func (m *rlWarningMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok { runtime.LogWarning(wailsCtx, toString(v)) }
	}
	return nil, nil
}

type rlErrorMethod struct{}

func (m *rlErrorMethod) GetName() string { return "error" }
func (m *rlErrorMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rlErrorMethod) GetIsStatic() bool { return true }
func (m *rlErrorMethod) GetReturnType() data.Types { return nil }
func (m *rlErrorMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "message", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *rlErrorMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "message", 0, data.NewBaseType("string"))}
}
func (m *rlErrorMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok { runtime.LogError(wailsCtx, toString(v)) }
	}
	return nil, nil
}

type rlFatalMethod struct{}

func (m *rlFatalMethod) GetName() string { return "fatal" }
func (m *rlFatalMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rlFatalMethod) GetIsStatic() bool { return true }
func (m *rlFatalMethod) GetReturnType() data.Types { return nil }
func (m *rlFatalMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "message", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *rlFatalMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "message", 0, data.NewBaseType("string"))}
}
func (m *rlFatalMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok { runtime.LogFatal(wailsCtx, toString(v)) }
	}
	return nil, nil
}

// ============================================================================
// Wails\Runtime\Browser — 浏览器操作 (静态方法)
// ============================================================================

type RuntimeBrowserClass struct{}

func NewRuntimeBrowserClass() data.ClassStmt { return &RuntimeBrowserClass{} }

func (c *RuntimeBrowserClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *RuntimeBrowserClass) GetFrom() data.From                            { return nil }
func (c *RuntimeBrowserClass) GetName() string                               { return "Wails\\Runtime\\Browser" }
func (c *RuntimeBrowserClass) GetExtend() *string                            { return nil }
func (c *RuntimeBrowserClass) GetImplements() []string                       { return nil }
func (c *RuntimeBrowserClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *RuntimeBrowserClass) GetPropertyList() []data.Property              { return nil }
func (c *RuntimeBrowserClass) GetConstruct() data.Method                     { return nil }

func (c *RuntimeBrowserClass) GetMethod(name string) (data.Method, bool) {
	if name == "openURL" { return &rbOpenURLMethod{}, true }
	return nil, false
}
func (c *RuntimeBrowserClass) GetStaticMethod(name string) (data.Method, bool) {
	if name == "openURL" { return &rbOpenURLMethod{}, true }
	return nil, false
}
func (c *RuntimeBrowserClass) GetMethods() []data.Method {
	return []data.Method{&rbOpenURLMethod{}}
}

type rbOpenURLMethod struct{}

func (m *rbOpenURLMethod) GetName() string { return "openURL" }
func (m *rbOpenURLMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rbOpenURLMethod) GetIsStatic() bool { return true }
func (m *rbOpenURLMethod) GetReturnType() data.Types { return nil }
func (m *rbOpenURLMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "url", 0, data.NewStringValue(""), data.NewBaseType("string"))}
}
func (m *rbOpenURLMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "url", 0, data.NewBaseType("string"))}
}
func (m *rbOpenURLMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx != nil {
		if v, ok := ctx.GetIndexValue(0); ok { runtime.BrowserOpenURL(wailsCtx, toString(v)) }
	}
	return nil, nil
}

// ============================================================================
// Wails\Runtime\Screen — 屏幕信息 (静态方法)
// ============================================================================

type RuntimeScreenClass struct{}

func NewRuntimeScreenClass() data.ClassStmt { return &RuntimeScreenClass{} }

func (c *RuntimeScreenClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *RuntimeScreenClass) GetFrom() data.From                            { return nil }
func (c *RuntimeScreenClass) GetName() string                               { return "Wails\\Runtime\\Screen" }
func (c *RuntimeScreenClass) GetExtend() *string                            { return nil }
func (c *RuntimeScreenClass) GetImplements() []string                       { return nil }
func (c *RuntimeScreenClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *RuntimeScreenClass) GetPropertyList() []data.Property              { return nil }
func (c *RuntimeScreenClass) GetConstruct() data.Method                     { return nil }

func (c *RuntimeScreenClass) GetMethod(name string) (data.Method, bool) {
	if name == "getAll" { return &rsGetAllMethod{}, true }
	return nil, false
}
func (c *RuntimeScreenClass) GetStaticMethod(name string) (data.Method, bool) {
	if name == "getAll" { return &rsGetAllMethod{}, true }
	return nil, false
}
func (c *RuntimeScreenClass) GetMethods() []data.Method {
	return []data.Method{&rsGetAllMethod{}}
}

type rsGetAllMethod struct{}

func (m *rsGetAllMethod) GetName() string { return "getAll" }
func (m *rsGetAllMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rsGetAllMethod) GetIsStatic() bool { return true }
func (m *rsGetAllMethod) GetReturnType() data.Types { return data.NewBaseType("array") }
func (m *rsGetAllMethod) GetParams() []data.GetValue { return nil }
func (m *rsGetAllMethod) GetVariables() []data.Variable { return nil }
func (m *rsGetAllMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx == nil { return data.NewArrayValue(nil), nil }
	screens, err := runtime.ScreenGetAll(wailsCtx)
	if err != nil { return data.NewArrayValue(nil), nil }
	vals := make([]data.Value, 0, len(screens))
	for _, s := range screens {
		vals = append(vals, data.NewArrayValue([]data.Value{
			data.NewBoolValue(s.IsCurrent),
			data.NewBoolValue(s.IsPrimary),
			data.NewIntValue(s.Width),
			data.NewIntValue(s.Height),
		}))
	}
	return data.NewArrayValue(vals), nil
}

// ============================================================================
// Wails\Runtime\Environment — 环境信息 (静态方法)
// ============================================================================

type RuntimeEnvironmentClass struct{}

func NewRuntimeEnvironmentClass() data.ClassStmt { return &RuntimeEnvironmentClass{} }

func (c *RuntimeEnvironmentClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *RuntimeEnvironmentClass) GetFrom() data.From                            { return nil }
func (c *RuntimeEnvironmentClass) GetName() string                               { return "Wails\\Runtime\\Environment" }
func (c *RuntimeEnvironmentClass) GetExtend() *string                            { return nil }
func (c *RuntimeEnvironmentClass) GetImplements() []string                       { return nil }
func (c *RuntimeEnvironmentClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *RuntimeEnvironmentClass) GetPropertyList() []data.Property              { return nil }
func (c *RuntimeEnvironmentClass) GetConstruct() data.Method                     { return nil }

func (c *RuntimeEnvironmentClass) GetMethod(name string) (data.Method, bool) {
	if name == "get" { return &reGetMethod{}, true }
	return nil, false
}
func (c *RuntimeEnvironmentClass) GetStaticMethod(name string) (data.Method, bool) {
	if name == "get" { return &reGetMethod{}, true }
	return nil, false
}
func (c *RuntimeEnvironmentClass) GetMethods() []data.Method {
	return []data.Method{&reGetMethod{}}
}

type reGetMethod struct{}

func (m *reGetMethod) GetName() string { return "get" }
func (m *reGetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *reGetMethod) GetIsStatic() bool { return true }
func (m *reGetMethod) GetReturnType() data.Types { return data.NewBaseType("array") }
func (m *reGetMethod) GetParams() []data.GetValue { return nil }
func (m *reGetMethod) GetVariables() []data.Variable { return nil }
func (m *reGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if wailsCtx == nil { return data.NewArrayValue(nil), nil }
	info := runtime.Environment(wailsCtx)
	return data.NewArrayValue([]data.Value{
		data.NewStringValue(info.BuildType),
		data.NewStringValue(info.Platform),
		data.NewStringValue(info.Arch),
	}), nil
}
