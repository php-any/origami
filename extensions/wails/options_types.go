package wails

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ============================================================================
// Wails\Options\RGBA — 颜色值
// ============================================================================

type RGBAClass struct{}

func NewRGBAClass() data.ClassStmt { return &RGBAClass{} }

func (c *RGBAClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *RGBAClass) GetFrom() data.From                            { return nil }
func (c *RGBAClass) GetName() string                               { return "Wails\\Options\\RGBA" }
func (c *RGBAClass) GetExtend() *string                            { return nil }
func (c *RGBAClass) GetImplements() []string                       { return nil }
func (c *RGBAClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *RGBAClass) GetPropertyList() []data.Property              { return nil }
func (c *RGBAClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *RGBAClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *RGBAClass) GetMethods() []data.Method { return nil }
func (c *RGBAClass) GetConstruct() data.Method { return &rgbaConstruct{} }

type rgbaConstruct struct{}

func (m *rgbaConstruct) GetName() string            { return token.ConstructName }
func (m *rgbaConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *rgbaConstruct) GetIsStatic() bool          { return false }
func (m *rgbaConstruct) GetReturnType() data.Types  { return nil }

func (m *rgbaConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "r", 0, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "g", 1, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "b", 2, data.NewIntValue(0), data.NewBaseType("int")),
		node.NewParameter(nil, "a", 3, data.NewIntValue(255), data.NewBaseType("int")),
	}
}

func (m *rgbaConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "r", 0, data.NewBaseType("int")),
		node.NewVariable(nil, "g", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "b", 2, data.NewBaseType("int")),
		node.NewVariable(nil, "a", 3, data.NewBaseType("int")),
	}
}

func (m *rgbaConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("r", data.NewIntValue(toInt(v)))
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		cv.SetProperty("g", data.NewIntValue(toInt(v)))
	}
	if v, ok := ctx.GetIndexValue(2); ok {
		cv.SetProperty("b", data.NewIntValue(toInt(v)))
	}
	if v, ok := ctx.GetIndexValue(3); ok {
		cv.SetProperty("a", data.NewIntValue(toInt(v)))
	}
	return nil, nil
}

// ============================================================================
// Wails\Options\SingleInstanceLock — 单实例锁
// ============================================================================

type SingleInstanceLockClass struct{}

func NewSingleInstanceLockClass() data.ClassStmt { return &SingleInstanceLockClass{} }

func (c *SingleInstanceLockClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *SingleInstanceLockClass) GetFrom() data.From                            { return nil }
func (c *SingleInstanceLockClass) GetName() string                               { return "Wails\\Options\\SingleInstanceLock" }
func (c *SingleInstanceLockClass) GetExtend() *string                            { return nil }
func (c *SingleInstanceLockClass) GetImplements() []string                       { return nil }
func (c *SingleInstanceLockClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SingleInstanceLockClass) GetPropertyList() []data.Property              { return nil }
func (c *SingleInstanceLockClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *SingleInstanceLockClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *SingleInstanceLockClass) GetMethods() []data.Method { return nil }
func (c *SingleInstanceLockClass) GetConstruct() data.Method { return &singleInstanceLockConstruct{} }

type singleInstanceLockConstruct struct{}

func (m *singleInstanceLockConstruct) GetName() string            { return token.ConstructName }
func (m *singleInstanceLockConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *singleInstanceLockConstruct) GetIsStatic() bool          { return false }
func (m *singleInstanceLockConstruct) GetReturnType() data.Types  { return nil }

func (m *singleInstanceLockConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "uniqueId", 0, data.NewStringValue(""), data.NewBaseType("string")),
	}
}

func (m *singleInstanceLockConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "uniqueId", 0, data.NewBaseType("string")),
	}
}

func (m *singleInstanceLockConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	uniqueId := ""
	if v, ok := ctx.GetIndexValue(0); ok {
		uniqueId = toString(v)
	}
	cv.SetProperty("UniqueId", data.NewStringValue(uniqueId))
	return nil, nil
}

// ============================================================================
// Wails\Options\DragAndDrop — 拖放配置
// ============================================================================

type DragAndDropClass struct{}

func NewDragAndDropClass() data.ClassStmt { return &DragAndDropClass{} }

func (c *DragAndDropClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *DragAndDropClass) GetFrom() data.From                            { return nil }
func (c *DragAndDropClass) GetName() string                               { return "Wails\\Options\\DragAndDrop" }
func (c *DragAndDropClass) GetExtend() *string                            { return nil }
func (c *DragAndDropClass) GetImplements() []string                       { return nil }
func (c *DragAndDropClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *DragAndDropClass) GetPropertyList() []data.Property              { return nil }
func (c *DragAndDropClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *DragAndDropClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *DragAndDropClass) GetMethods() []data.Method { return nil }
func (c *DragAndDropClass) GetConstruct() data.Method { return &dragAndDropConstruct{} }

type dragAndDropConstruct struct{}

func (m *dragAndDropConstruct) GetName() string            { return token.ConstructName }
func (m *dragAndDropConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *dragAndDropConstruct) GetIsStatic() bool          { return false }
func (m *dragAndDropConstruct) GetReturnType() data.Types  { return nil }

func (m *dragAndDropConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "enableFileDrop", 0, data.NewBoolValue(false), data.NewBaseType("bool")),
		node.NewParameter(nil, "disableWebViewDrop", 1, data.NewBoolValue(false), data.NewBaseType("bool")),
		node.NewParameter(nil, "cssDropProperty", 2, data.NewStringValue("--wails-drop-target"), data.NewBaseType("string")),
		node.NewParameter(nil, "cssDropValue", 3, data.NewStringValue("drop"), data.NewBaseType("string")),
	}
}

func (m *dragAndDropConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "enableFileDrop", 0, data.NewBaseType("bool")),
		node.NewVariable(nil, "disableWebViewDrop", 1, data.NewBaseType("bool")),
		node.NewVariable(nil, "cssDropProperty", 2, data.NewBaseType("string")),
		node.NewVariable(nil, "cssDropValue", 3, data.NewBaseType("string")),
	}
}

func (m *dragAndDropConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("EnableFileDrop", data.NewBoolValue(toBool(v)))
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		cv.SetProperty("DisableWebViewDrop", data.NewBoolValue(toBool(v)))
	}
	if v, ok := ctx.GetIndexValue(2); ok {
		cv.SetProperty("CSSDropProperty", data.NewStringValue(toString(v)))
	}
	if v, ok := ctx.GetIndexValue(3); ok {
		cv.SetProperty("CSSDropValue", data.NewStringValue(toString(v)))
	}
	return nil, nil
}

// ============================================================================
// Wails\Options\Debug — 调试配置
// ============================================================================

type DebugClass struct{}

func NewDebugClass() data.ClassStmt { return &DebugClass{} }

func (c *DebugClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *DebugClass) GetFrom() data.From                            { return nil }
func (c *DebugClass) GetName() string                               { return "Wails\\Options\\Debug" }
func (c *DebugClass) GetExtend() *string                            { return nil }
func (c *DebugClass) GetImplements() []string                       { return nil }
func (c *DebugClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *DebugClass) GetPropertyList() []data.Property              { return nil }
func (c *DebugClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *DebugClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *DebugClass) GetMethods() []data.Method { return nil }
func (c *DebugClass) GetConstruct() data.Method { return &debugConstruct{} }

type debugConstruct struct{}

func (m *debugConstruct) GetName() string            { return token.ConstructName }
func (m *debugConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *debugConstruct) GetIsStatic() bool          { return false }
func (m *debugConstruct) GetReturnType() data.Types  { return nil }

func (m *debugConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "openInspectorOnStartup", 0, data.NewBoolValue(false), data.NewBaseType("bool")),
	}
}

func (m *debugConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "openInspectorOnStartup", 0, data.NewBaseType("bool")),
	}
}

func (m *debugConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("OpenInspectorOnStartup", data.NewBoolValue(toBool(v)))
	}
	return nil, nil
}

// ============================================================================
// Wails\Options\AssetServer — 资源服务器配置
// ============================================================================

type AssetServerClass struct{}

func NewAssetServerClass() data.ClassStmt { return &AssetServerClass{} }

func (c *AssetServerClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *AssetServerClass) GetFrom() data.From                            { return nil }
func (c *AssetServerClass) GetName() string                               { return "Wails\\Options\\AssetServer" }
func (c *AssetServerClass) GetExtend() *string                            { return nil }
func (c *AssetServerClass) GetImplements() []string                       { return nil }
func (c *AssetServerClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *AssetServerClass) GetPropertyList() []data.Property              { return nil }
func (c *AssetServerClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *AssetServerClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *AssetServerClass) GetMethods() []data.Method { return nil }
func (c *AssetServerClass) GetConstruct() data.Method { return &assetServerConstruct{} }

type assetServerConstruct struct{}

func (m *assetServerConstruct) GetName() string            { return token.ConstructName }
func (m *assetServerConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *assetServerConstruct) GetIsStatic() bool          { return false }
func (m *assetServerConstruct) GetReturnType() data.Types  { return nil }

func (m *assetServerConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "assetsDir", 0, data.NewStringValue(""), data.NewBaseType("string")),
	}
}

func (m *assetServerConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "assetsDir", 0, data.NewBaseType("string")),
	}
}

func (m *assetServerConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	assetsDir := ""
	if v, ok := ctx.GetIndexValue(0); ok {
		assetsDir = toString(v)
	}
	cv.SetProperty("AssetsDir", data.NewStringValue(assetsDir))
	return nil, nil
}

// ============================================================================
// Wails\Options\Windows — Windows 平台选项
// ============================================================================

type WindowsOptionsClass struct{}

func NewWindowsOptionsClass() data.ClassStmt { return &WindowsOptionsClass{} }

func (c *WindowsOptionsClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *WindowsOptionsClass) GetFrom() data.From                            { return nil }
func (c *WindowsOptionsClass) GetName() string                               { return "Wails\\Options\\Windows" }
func (c *WindowsOptionsClass) GetExtend() *string                            { return nil }
func (c *WindowsOptionsClass) GetImplements() []string                       { return nil }
func (c *WindowsOptionsClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *WindowsOptionsClass) GetPropertyList() []data.Property              { return nil }
func (c *WindowsOptionsClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *WindowsOptionsClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *WindowsOptionsClass) GetMethods() []data.Method { return nil }
func (c *WindowsOptionsClass) GetConstruct() data.Method { return &windowsOptionsConstruct{} }

type windowsOptionsConstruct struct{}

func (m *windowsOptionsConstruct) GetName() string            { return token.ConstructName }
func (m *windowsOptionsConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *windowsOptionsConstruct) GetIsStatic() bool          { return false }
func (m *windowsOptionsConstruct) GetReturnType() data.Types  { return nil }

func (m *windowsOptionsConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("array")),
	}
}

func (m *windowsOptionsConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, data.NewBaseType("array")),
	}
}

func (m *windowsOptionsConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	setDefaultBoolProperty(cv, "WebviewIsTransparent", false)
	setDefaultBoolProperty(cv, "WindowIsTranslucent", false)
	setDefaultBoolProperty(cv, "ContentProtection", false)
	setDefaultBoolProperty(cv, "DisablePinchZoom", false)
	setDefaultBoolProperty(cv, "DisableWindowIcon", false)
	setDefaultBoolProperty(cv, "DisableFramelessWindowDecorations", false)
	setDefaultBoolProperty(cv, "WebviewGpuDisabled", false)
	setDefaultBoolProperty(cv, "IsZoomControlEnabled", false)
	setDefaultBoolProperty(cv, "EnableSwipeGestures", false)
	setDefaultBoolProperty(cv, "WebviewDisableRendererCodeIntegrity", false)
	setDefaultIntProperty(cv, "Theme", 0)         // SystemDefault
	setDefaultIntProperty(cv, "BackdropType", 0)  // Auto
	setDefaultFloatProperty(cv, "ZoomFactor", 1.0)
	setDefaultIntProperty(cv, "ResizeDebounceMS", 0)
	setDefaultStringProperty(cv, "WebviewUserDataPath", "")
	setDefaultStringProperty(cv, "WebviewBrowserPath", "")
	setDefaultStringProperty(cv, "WindowClassName", "wailsWindow")

	if v, ok := ctx.GetIndexValue(0); ok {
		if av, ok := v.(*data.ArrayValue); ok {
			applyArrayToClassValue(cv, av, []string{
				"WebviewIsTransparent", "WindowIsTranslucent", "ContentProtection",
				"DisablePinchZoom", "DisableWindowIcon", "DisableFramelessWindowDecorations",
				"WebviewGpuDisabled", "IsZoomControlEnabled", "EnableSwipeGestures",
				"WebviewDisableRendererCodeIntegrity", "Theme", "BackdropType",
				"ZoomFactor", "ResizeDebounceMS", "WebviewUserDataPath",
				"WebviewBrowserPath", "WindowClassName",
			})
		}
	}
	return nil, nil
}

// ============================================================================
// Wails\Options\Mac — macOS 平台选项
// ============================================================================

type MacOptionsClass struct{}

func NewMacOptionsClass() data.ClassStmt { return &MacOptionsClass{} }

func (c *MacOptionsClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *MacOptionsClass) GetFrom() data.From                            { return nil }
func (c *MacOptionsClass) GetName() string                               { return "Wails\\Options\\Mac" }
func (c *MacOptionsClass) GetExtend() *string                            { return nil }
func (c *MacOptionsClass) GetImplements() []string                       { return nil }
func (c *MacOptionsClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *MacOptionsClass) GetPropertyList() []data.Property              { return nil }
func (c *MacOptionsClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *MacOptionsClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *MacOptionsClass) GetMethods() []data.Method { return nil }
func (c *MacOptionsClass) GetConstruct() data.Method { return &macOptionsConstruct{} }

type macOptionsConstruct struct{}

func (m *macOptionsConstruct) GetName() string            { return token.ConstructName }
func (m *macOptionsConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *macOptionsConstruct) GetIsStatic() bool          { return false }
func (m *macOptionsConstruct) GetReturnType() data.Types  { return nil }

func (m *macOptionsConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("array")),
	}
}

func (m *macOptionsConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, data.NewBaseType("array")),
	}
}

func (m *macOptionsConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	setDefaultBoolProperty(cv, "WebviewIsTransparent", false)
	setDefaultBoolProperty(cv, "WindowIsTranslucent", false)
	setDefaultBoolProperty(cv, "ContentProtection", false)
	setDefaultBoolProperty(cv, "DisableZoom", false)
	setDefaultStringProperty(cv, "Appearance", "")

	if v, ok := ctx.GetIndexValue(0); ok {
		if av, ok := v.(*data.ArrayValue); ok {
			applyArrayToClassValue(cv, av, []string{
				"WebviewIsTransparent", "WindowIsTranslucent",
				"ContentProtection", "DisableZoom", "Appearance",
			})
		}
	}
	return nil, nil
}

// ============================================================================
// Wails\Options\Mac\TitleBar — macOS 标题栏配置
// ============================================================================

type MacTitleBarClass struct{}

func NewMacTitleBarClass() data.ClassStmt { return &MacTitleBarClass{} }

func (c *MacTitleBarClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *MacTitleBarClass) GetFrom() data.From                            { return nil }
func (c *MacTitleBarClass) GetName() string                               { return "Wails\\Options\\Mac\\TitleBar" }
func (c *MacTitleBarClass) GetExtend() *string                            { return nil }
func (c *MacTitleBarClass) GetImplements() []string                       { return nil }
func (c *MacTitleBarClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *MacTitleBarClass) GetPropertyList() []data.Property              { return nil }
func (c *MacTitleBarClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *MacTitleBarClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *MacTitleBarClass) GetMethods() []data.Method { return nil }
func (c *MacTitleBarClass) GetConstruct() data.Method { return &macTitleBarConstruct{} }

type macTitleBarConstruct struct{}

func (m *macTitleBarConstruct) GetName() string            { return token.ConstructName }
func (m *macTitleBarConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *macTitleBarConstruct) GetIsStatic() bool          { return false }
func (m *macTitleBarConstruct) GetReturnType() data.Types  { return nil }

func (m *macTitleBarConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("array")),
	}
}

func (m *macTitleBarConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, data.NewBaseType("array")),
	}
}

func (m *macTitleBarConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	setDefaultBoolProperty(cv, "TitlebarAppearsTransparent", false)
	setDefaultBoolProperty(cv, "HideTitle", false)
	setDefaultBoolProperty(cv, "HideTitleBar", false)
	setDefaultBoolProperty(cv, "FullSizeContent", false)
	setDefaultBoolProperty(cv, "UseToolbar", false)
	setDefaultBoolProperty(cv, "HideToolbarSeparator", false)

	if v, ok := ctx.GetIndexValue(0); ok {
		if av, ok := v.(*data.ArrayValue); ok {
			applyArrayToClassValue(cv, av, []string{
				"TitlebarAppearsTransparent", "HideTitle", "HideTitleBar",
				"FullSizeContent", "UseToolbar", "HideToolbarSeparator",
			})
		}
	}
	return nil, nil
}

// ============================================================================
// Wails\Options\Mac\AboutInfo — macOS 关于信息
// ============================================================================

type MacAboutInfoClass struct{}

func NewMacAboutInfoClass() data.ClassStmt { return &MacAboutInfoClass{} }

func (c *MacAboutInfoClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *MacAboutInfoClass) GetFrom() data.From                            { return nil }
func (c *MacAboutInfoClass) GetName() string                               { return "Wails\\Options\\Mac\\AboutInfo" }
func (c *MacAboutInfoClass) GetExtend() *string                            { return nil }
func (c *MacAboutInfoClass) GetImplements() []string                       { return nil }
func (c *MacAboutInfoClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *MacAboutInfoClass) GetPropertyList() []data.Property              { return nil }
func (c *MacAboutInfoClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *MacAboutInfoClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *MacAboutInfoClass) GetMethods() []data.Method { return nil }
func (c *MacAboutInfoClass) GetConstruct() data.Method { return &macAboutInfoConstruct{} }

type macAboutInfoConstruct struct{}

func (m *macAboutInfoConstruct) GetName() string            { return token.ConstructName }
func (m *macAboutInfoConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *macAboutInfoConstruct) GetIsStatic() bool          { return false }
func (m *macAboutInfoConstruct) GetReturnType() data.Types  { return nil }

func (m *macAboutInfoConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "title", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "message", 1, data.NewStringValue(""), data.NewBaseType("string")),
	}
}

func (m *macAboutInfoConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "title", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "message", 1, data.NewBaseType("string")),
	}
}

func (m *macAboutInfoConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("Title", data.NewStringValue(toString(v)))
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		cv.SetProperty("Message", data.NewStringValue(toString(v)))
	}
	return nil, nil
}

// ============================================================================
// Wails\Options\Linux — Linux 平台选项
// ============================================================================

type LinuxOptionsClass struct{}

func NewLinuxOptionsClass() data.ClassStmt { return &LinuxOptionsClass{} }

func (c *LinuxOptionsClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *LinuxOptionsClass) GetFrom() data.From                            { return nil }
func (c *LinuxOptionsClass) GetName() string                               { return "Wails\\Options\\Linux" }
func (c *LinuxOptionsClass) GetExtend() *string                            { return nil }
func (c *LinuxOptionsClass) GetImplements() []string                       { return nil }
func (c *LinuxOptionsClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *LinuxOptionsClass) GetPropertyList() []data.Property              { return nil }
func (c *LinuxOptionsClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *LinuxOptionsClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *LinuxOptionsClass) GetMethods() []data.Method { return nil }
func (c *LinuxOptionsClass) GetConstruct() data.Method { return &linuxOptionsConstruct{} }

type linuxOptionsConstruct struct{}

func (m *linuxOptionsConstruct) GetName() string            { return token.ConstructName }
func (m *linuxOptionsConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *linuxOptionsConstruct) GetIsStatic() bool          { return false }
func (m *linuxOptionsConstruct) GetReturnType() data.Types  { return nil }

func (m *linuxOptionsConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("array")),
	}
}

func (m *linuxOptionsConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, data.NewBaseType("array")),
	}
}

func (m *linuxOptionsConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	setDefaultBoolProperty(cv, "WindowIsTranslucent", false)
	setDefaultIntProperty(cv, "WebviewGpuPolicy", 0)
	setDefaultStringProperty(cv, "ProgramName", "")

	if v, ok := ctx.GetIndexValue(0); ok {
		if av, ok := v.(*data.ArrayValue); ok {
			applyArrayToClassValue(cv, av, []string{
				"WindowIsTranslucent", "WebviewGpuPolicy", "ProgramName",
			})
		}
	}
	return nil, nil
}

// ============================================================================
// Wails\Options\SystemTray — 系统托盘选项
// ============================================================================

type SystemTrayClass struct{}

func NewSystemTrayClass() data.ClassStmt { return &SystemTrayClass{} }

func (c *SystemTrayClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *SystemTrayClass) GetFrom() data.From                            { return nil }
func (c *SystemTrayClass) GetName() string                               { return "Wails\\Options\\SystemTray" }
func (c *SystemTrayClass) GetExtend() *string                            { return nil }
func (c *SystemTrayClass) GetImplements() []string                       { return nil }
func (c *SystemTrayClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SystemTrayClass) GetPropertyList() []data.Property              { return nil }
func (c *SystemTrayClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *SystemTrayClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *SystemTrayClass) GetMethods() []data.Method { return nil }
func (c *SystemTrayClass) GetConstruct() data.Method { return &systemTrayConstruct{} }

type systemTrayConstruct struct{}

func (m *systemTrayConstruct) GetName() string            { return token.ConstructName }
func (m *systemTrayConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *systemTrayConstruct) GetIsStatic() bool          { return false }
func (m *systemTrayConstruct) GetReturnType() data.Types  { return nil }

func (m *systemTrayConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("array")),
	}
}

func (m *systemTrayConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, data.NewBaseType("array")),
	}
}

func (m *systemTrayConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	setDefaultStringProperty(cv, "Title", "")
	setDefaultStringProperty(cv, "Tooltip", "")
	setDefaultBoolProperty(cv, "StartHidden", false)

	if v, ok := ctx.GetIndexValue(0); ok {
		if av, ok := v.(*data.ArrayValue); ok {
			applyArrayToClassValue(cv, av, []string{
				"Title", "Tooltip", "StartHidden",
			})
		}
	}
	return nil, nil
}
