package wails

import "github.com/php-any/origami/data"

// ============================================================================
// Wails\WindowStartState — 窗口启动状态枚举
// ============================================================================

type WindowStartStateClass struct{}

func NewWindowStartStateClass() data.ClassStmt { return &WindowStartStateClass{} }

func (c *WindowStartStateClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *WindowStartStateClass) GetFrom() data.From                            { return nil }
func (c *WindowStartStateClass) GetName() string                               { return "Wails\\WindowStartState" }
func (c *WindowStartStateClass) GetExtend() *string                            { return nil }
func (c *WindowStartStateClass) GetImplements() []string                       { return nil }
func (c *WindowStartStateClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *WindowStartStateClass) GetPropertyList() []data.Property              { return nil }
func (c *WindowStartStateClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *WindowStartStateClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *WindowStartStateClass) GetMethods() []data.Method { return nil }
func (c *WindowStartStateClass) GetConstruct() data.Method { return nil }

var windowStartStates = map[string]int{
	"NORMAL":     0,
	"MAXIMISED":  1,
	"MINIMISED":  2,
	"FULLSCREEN": 3,
}

func (c *WindowStartStateClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := windowStartStates[name]; ok {
		return data.NewIntValue(v), true
	}
	return nil, false
}

// ============================================================================
// Wails\BackdropType — Windows 背景材质类型枚举
// ============================================================================

type BackdropTypeClass struct{}

func NewBackdropTypeClass() data.ClassStmt { return &BackdropTypeClass{} }

func (c *BackdropTypeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *BackdropTypeClass) GetFrom() data.From                            { return nil }
func (c *BackdropTypeClass) GetName() string                               { return "Wails\\BackdropType" }
func (c *BackdropTypeClass) GetExtend() *string                            { return nil }
func (c *BackdropTypeClass) GetImplements() []string                       { return nil }
func (c *BackdropTypeClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *BackdropTypeClass) GetPropertyList() []data.Property              { return nil }
func (c *BackdropTypeClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *BackdropTypeClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *BackdropTypeClass) GetMethods() []data.Method { return nil }
func (c *BackdropTypeClass) GetConstruct() data.Method { return nil }

var backdropTypes = map[string]int32{
	"AUTO":    0,
	"NONE":    1,
	"MICA":    2,
	"ACRYLIC": 3,
	"TABBED":  4,
}

func (c *BackdropTypeClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := backdropTypes[name]; ok {
		return data.NewIntValue(int(v)), true
	}
	return nil, false
}

// ============================================================================
// Wails\Theme — Windows 主题枚举
// ============================================================================

type ThemeClass struct{}

func NewThemeClass() data.ClassStmt { return &ThemeClass{} }

func (c *ThemeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *ThemeClass) GetFrom() data.From                            { return nil }
func (c *ThemeClass) GetName() string                               { return "Wails\\Theme" }
func (c *ThemeClass) GetExtend() *string                            { return nil }
func (c *ThemeClass) GetImplements() []string                       { return nil }
func (c *ThemeClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *ThemeClass) GetPropertyList() []data.Property              { return nil }
func (c *ThemeClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *ThemeClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *ThemeClass) GetMethods() []data.Method { return nil }
func (c *ThemeClass) GetConstruct() data.Method { return nil }

var themes = map[string]int{
	"SYSTEM_DEFAULT": 0,
	"DARK":           1,
	"LIGHT":          2,
}

func (c *ThemeClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := themes[name]; ok {
		return data.NewIntValue(v), true
	}
	return nil, false
}

// ============================================================================
// Wails\WebviewGpuPolicy — Linux WebView GPU 策略枚举
// ============================================================================

type WebviewGpuPolicyClass struct{}

func NewWebviewGpuPolicyClass() data.ClassStmt { return &WebviewGpuPolicyClass{} }

func (c *WebviewGpuPolicyClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *WebviewGpuPolicyClass) GetFrom() data.From                            { return nil }
func (c *WebviewGpuPolicyClass) GetName() string                               { return "Wails\\WebviewGpuPolicy" }
func (c *WebviewGpuPolicyClass) GetExtend() *string                            { return nil }
func (c *WebviewGpuPolicyClass) GetImplements() []string                       { return nil }
func (c *WebviewGpuPolicyClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *WebviewGpuPolicyClass) GetPropertyList() []data.Property              { return nil }
func (c *WebviewGpuPolicyClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *WebviewGpuPolicyClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *WebviewGpuPolicyClass) GetMethods() []data.Method { return nil }
func (c *WebviewGpuPolicyClass) GetConstruct() data.Method { return nil }

var gpuPolicies = map[string]int{
	"ALWAYS":    0,
	"ON_DEMAND": 1,
	"NEVER":     2,
}

func (c *WebviewGpuPolicyClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := gpuPolicies[name]; ok {
		return data.NewIntValue(v), true
	}
	return nil, false
}

// ============================================================================
// Wails\DialogType — 对话框类型枚举
// ============================================================================

type DialogTypeClass struct{}

func NewDialogTypeClass() data.ClassStmt { return &DialogTypeClass{} }

func (c *DialogTypeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *DialogTypeClass) GetFrom() data.From                            { return nil }
func (c *DialogTypeClass) GetName() string                               { return "Wails\\DialogType" }
func (c *DialogTypeClass) GetExtend() *string                            { return nil }
func (c *DialogTypeClass) GetImplements() []string                       { return nil }
func (c *DialogTypeClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *DialogTypeClass) GetPropertyList() []data.Property              { return nil }
func (c *DialogTypeClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *DialogTypeClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *DialogTypeClass) GetMethods() []data.Method { return nil }
func (c *DialogTypeClass) GetConstruct() data.Method { return nil }

var dialogTypes = map[string]string{
	"INFO":     "info",
	"WARNING":  "warning",
	"ERROR":    "error",
	"QUESTION": "question",
}

func (c *DialogTypeClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := dialogTypes[name]; ok {
		return data.NewStringValue(v), true
	}
	return nil, false
}

// ============================================================================
// Wails\LogLevel — 日志级别枚举
// ============================================================================

type LogLevelClass struct{}

func NewLogLevelClass() data.ClassStmt { return &LogLevelClass{} }

func (c *LogLevelClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *LogLevelClass) GetFrom() data.From                            { return nil }
func (c *LogLevelClass) GetName() string                               { return "Wails\\LogLevel" }
func (c *LogLevelClass) GetExtend() *string                            { return nil }
func (c *LogLevelClass) GetImplements() []string                       { return nil }
func (c *LogLevelClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *LogLevelClass) GetPropertyList() []data.Property              { return nil }
func (c *LogLevelClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *LogLevelClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *LogLevelClass) GetMethods() []data.Method { return nil }
func (c *LogLevelClass) GetConstruct() data.Method { return nil }

var logLevels = map[string]int{
	"TRACE":   1,
	"DEBUG":   2,
	"INFO":    3,
	"WARNING": 4,
	"ERROR":   5,
}

func (c *LogLevelClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := logLevels[name]; ok {
		return data.NewIntValue(v), true
	}
	return nil, false
}

// ============================================================================
// Wails\MenuItemType — 菜单项类型常量
// ============================================================================

type MenuItemTypeClass struct{}

func NewMenuItemTypeClass() data.ClassStmt { return &MenuItemTypeClass{} }

func (c *MenuItemTypeClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *MenuItemTypeClass) GetFrom() data.From                            { return nil }
func (c *MenuItemTypeClass) GetName() string                               { return "Wails\\MenuItemType" }
func (c *MenuItemTypeClass) GetExtend() *string                            { return nil }
func (c *MenuItemTypeClass) GetImplements() []string                       { return nil }
func (c *MenuItemTypeClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *MenuItemTypeClass) GetPropertyList() []data.Property              { return nil }
func (c *MenuItemTypeClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *MenuItemTypeClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *MenuItemTypeClass) GetMethods() []data.Method { return nil }
func (c *MenuItemTypeClass) GetConstruct() data.Method { return nil }

var menuItemTypes = map[string]string{
	"TEXT":      "Text",
	"SEPARATOR": "Separator",
	"SUBMENU":   "Submenu",
	"CHECKBOX":  "Checkbox",
	"RADIO":     "Radio",
}

func (c *MenuItemTypeClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := menuItemTypes[name]; ok {
		return data.NewStringValue(v), true
	}
	return nil, false
}

// ============================================================================
// Wails\MacAppearance — macOS 外观常量
// ============================================================================

type MacAppearanceClass struct{}

func NewMacAppearanceClass() data.ClassStmt { return &MacAppearanceClass{} }

func (c *MacAppearanceClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *MacAppearanceClass) GetFrom() data.From                            { return nil }
func (c *MacAppearanceClass) GetName() string                               { return "Wails\\MacAppearance" }
func (c *MacAppearanceClass) GetExtend() *string                            { return nil }
func (c *MacAppearanceClass) GetImplements() []string                       { return nil }
func (c *MacAppearanceClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *MacAppearanceClass) GetPropertyList() []data.Property              { return nil }
func (c *MacAppearanceClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *MacAppearanceClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *MacAppearanceClass) GetMethods() []data.Method { return nil }
func (c *MacAppearanceClass) GetConstruct() data.Method { return nil }

var macAppearances = map[string]string{
	"DEFAULT":          "",
	"AQUA":             "NSAppearanceNameAqua",
	"DARK_AQUA":        "NSAppearanceNameDarkAqua",
	"VIBRANT_LIGHT":    "NSAppearanceNameVibrantLight",
	"VIBRANT_DARK":     "NSAppearanceNameVibrantDark",
	"ACCESSIBILITY_HC": "NSAppearanceNameAccessibilityHighContrast",
}

func (c *MacAppearanceClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := macAppearances[name]; ok {
		return data.NewStringValue(v), true
	}
	return nil, false
}

// ============================================================================
// Wails\ImagePosition — macOS 托盘图标位置枚举
// ============================================================================

type ImagePositionClass struct{}

func NewImagePositionClass() data.ClassStmt { return &ImagePositionClass{} }

func (c *ImagePositionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *ImagePositionClass) GetFrom() data.From                            { return nil }
func (c *ImagePositionClass) GetName() string                               { return "Wails\\ImagePosition" }
func (c *ImagePositionClass) GetExtend() *string                            { return nil }
func (c *ImagePositionClass) GetImplements() []string                       { return nil }
func (c *ImagePositionClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *ImagePositionClass) GetPropertyList() []data.Property              { return nil }
func (c *ImagePositionClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *ImagePositionClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *ImagePositionClass) GetMethods() []data.Method { return nil }
func (c *ImagePositionClass) GetConstruct() data.Method { return nil }

var imagePositions = map[string]int{
	"NO_IMAGE":         0,
	"IMAGE_ONLY":       1,
	"IMAGE_LEFT":       2,
	"IMAGE_RIGHT":      3,
	"IMAGE_BELOW":      4,
	"IMAGE_ABOVE":      5,
	"IMAGE_OVERLAPS":   6,
	"IMAGE_LEADING":    7,
	"IMAGE_TRAILING":   8,
}

func (c *ImagePositionClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := imagePositions[name]; ok {
		return data.NewIntValue(v), true
	}
	return nil, false
}
