package wails

import (
	"context"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// buildAppOptions 从 Origami ClassValue 构建 Wails options.App
func buildAppOptions(v data.Value) (*options.App, error) {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return nil, errors.New("Wails\\Options\\App instance required")
	}

	app := &options.App{}

	// 核心窗口属性
	app.Title = getPropString(cv, "Title", "")
	app.Width = getPropInt(cv, "Width", 1024)
	app.Height = getPropInt(cv, "Height", 768)
	app.DisableResize = getPropBool(cv, "DisableResize", false)
	app.Frameless = getPropBool(cv, "Frameless", false)
	app.MinWidth = getPropInt(cv, "MinWidth", 0)
	app.MinHeight = getPropInt(cv, "MinHeight", 0)
	app.MaxWidth = getPropInt(cv, "MaxWidth", 0)
	app.MaxHeight = getPropInt(cv, "MaxHeight", 0)
	app.StartHidden = getPropBool(cv, "StartHidden", false)
	app.HideWindowOnClose = getPropBool(cv, "HideWindowOnClose", false)
	app.AlwaysOnTop = getPropBool(cv, "AlwaysOnTop", false)
	app.CSSDragProperty = getPropString(cv, "CSSDragProperty", "--wails-draggable")
	app.CSSDragValue = getPropString(cv, "CSSDragValue", "drag")
	app.EnableFraudulentWebsiteDetection = getPropBool(cv, "EnableFraudulentWebsiteDetection", false)
	app.WindowStartState = options.WindowStartState(getPropInt(cv, "WindowStartState", 0))

	// 背景色
	if v, ctrl := cv.GetProperty("BackgroundColour"); ctrl == nil && v != nil {
		if rgba := buildRGBA(v); rgba != nil {
			app.BackgroundColour = rgba
		}
	}

	// 平台选项
	if v, ctrl := cv.GetProperty("Windows"); ctrl == nil && v != nil {
		app.Windows = buildWindowsOptions(v)
	}
	if v, ctrl := cv.GetProperty("Mac"); ctrl == nil && v != nil {
		app.Mac = buildMacOptions(v)
	}
	if v, ctrl := cv.GetProperty("Linux"); ctrl == nil && v != nil {
		app.Linux = buildLinuxOptions(v)
	}

	// Debug
	if v, ctrl := cv.GetProperty("Debug"); ctrl == nil && v != nil {
		if debugCV, ok := v.(*data.ClassValue); ok {
			app.Debug = options.Debug{
				OpenInspectorOnStartup: getPropBool(debugCV, "OpenInspectorOnStartup", false),
			}
		}
	}

	// 生命周期回调 — 在回调中设置 Wails 上下文
	app.OnStartup = func(ctx context.Context) {
		SetWailsContext(ctx)
		// invokePHP callback if set
	}
	app.OnDomReady = func(ctx context.Context) {
		SetWailsContext(ctx)
	}
	app.OnShutdown = func(ctx context.Context) {
		// invokePHP callback if set
	}
	app.OnBeforeClose = func(ctx context.Context) bool {
		return false
	}

	// 绑定 — 暂只支持预注册的 Go struct 对象
	if v, ctrl := cv.GetProperty("_bind"); ctrl == nil && v != nil {
		if av, ok := v.(*data.ArrayValue); ok {
			bindings := make([]interface{}, len(av.List))
			for i, z := range av.List {
				if z != nil {
					bindings[i] = z.Value
				}
			}
			app.Bind = bindings
		}
	}

	// 日志级别
	if v, ctrl := cv.GetProperty("LogLevel"); ctrl == nil && v != nil {
		app.LogLevel = logger.LogLevel(toInt(v))
	}

	// 应用默认值
	options.MergeDefaults(app)

	return app, nil
}

// buildRGBA 构建 options.RGBA
func buildRGBA(v data.Value) *options.RGBA {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return nil
	}
	return &options.RGBA{
		R: uint8(getPropInt(cv, "r", 0)),
		G: uint8(getPropInt(cv, "g", 0)),
		B: uint8(getPropInt(cv, "b", 0)),
		A: uint8(getPropInt(cv, "a", 255)),
	}
}

// buildWindowsOptions 构建 windows.Options
func buildWindowsOptions(v data.Value) *windows.Options {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return nil
	}
	return &windows.Options{
		WebviewIsTransparent:               getPropBool(cv, "WebviewIsTransparent", false),
		WindowIsTranslucent:                getPropBool(cv, "WindowIsTranslucent", false),
		DisableWindowIcon:                  getPropBool(cv, "DisableWindowIcon", false),
		DisableFramelessWindowDecorations:  getPropBool(cv, "DisableFramelessWindowDecorations", false),
		WebviewGpuIsDisabled:               getPropBool(cv, "WebviewGpuDisabled", false),
		IsZoomControlEnabled:               getPropBool(cv, "IsZoomControlEnabled", false),
		WebviewDisableRendererCodeIntegrity: getPropBool(cv, "WebviewDisableRendererCodeIntegrity", false),
		Theme:                              windows.Theme(getPropInt(cv, "Theme", 0)),
		BackdropType:                       windows.BackdropType(getPropInt(cv, "BackdropType", 0)),
		ZoomFactor:                         getPropFloat(cv, "ZoomFactor", 1.0),
		ResizeDebounceMS:                   uint16(getPropInt(cv, "ResizeDebounceMS", 0)),
		WebviewUserDataPath:                getPropString(cv, "WebviewUserDataPath", ""),
		WebviewBrowserPath:                 getPropString(cv, "WebviewBrowserPath", ""),
	}
}

// buildMacOptions 构建 mac.Options
func buildMacOptions(v data.Value) *mac.Options {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return nil
	}
	opts := &mac.Options{
		WebviewIsTransparent: getPropBool(cv, "WebviewIsTransparent", false),
		WindowIsTranslucent:  getPropBool(cv, "WindowIsTranslucent", false),
		Appearance:           mac.AppearanceType(getPropString(cv, "Appearance", "")),
	}
	if v, ctrl := cv.GetProperty("TitleBar"); ctrl == nil && v != nil {
		opts.TitleBar = buildMacTitleBar(v)
	}
	if v, ctrl := cv.GetProperty("About"); ctrl == nil && v != nil {
		opts.About = buildMacAboutInfo(v)
	}
	return opts
}

func buildMacTitleBar(v data.Value) *mac.TitleBar {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return nil
	}
	return &mac.TitleBar{
		TitlebarAppearsTransparent: getPropBool(cv, "TitlebarAppearsTransparent", false),
		HideTitle:                  getPropBool(cv, "HideTitle", false),
		HideTitleBar:               getPropBool(cv, "HideTitleBar", false),
		FullSizeContent:            getPropBool(cv, "FullSizeContent", false),
		UseToolbar:                 getPropBool(cv, "UseToolbar", false),
		HideToolbarSeparator:       getPropBool(cv, "HideToolbarSeparator", false),
	}
}

func buildMacAboutInfo(v data.Value) *mac.AboutInfo {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return nil
	}
	return &mac.AboutInfo{
		Title:   getPropString(cv, "Title", ""),
		Message: getPropString(cv, "Message", ""),
	}
}

// buildLinuxOptions 构建 linux.Options
func buildLinuxOptions(v data.Value) *linux.Options {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return nil
	}
	return &linux.Options{
		WindowIsTranslucent: getPropBool(cv, "WindowIsTranslucent", false),
		WebviewGpuPolicy:    linux.WebviewGpuPolicy(getPropInt(cv, "WebviewGpuPolicy", 0)),
	}
}

// ============================================================================
// ClassValue 属性读取辅助
// ============================================================================

func getPropString(cv *data.ClassValue, name string, defaultVal string) string {
	if v, ctrl := cv.GetProperty(name); ctrl == nil && v != nil {
		return toString(v)
	}
	return defaultVal
}

func getPropInt(cv *data.ClassValue, name string, defaultVal int) int {
	if v, ctrl := cv.GetProperty(name); ctrl == nil && v != nil {
		return toInt(v)
	}
	return defaultVal
}

func getPropBool(cv *data.ClassValue, name string, defaultVal bool) bool {
	if v, ctrl := cv.GetProperty(name); ctrl == nil && v != nil {
		return toBool(v)
	}
	return defaultVal
}

func getPropFloat(cv *data.ClassValue, name string, defaultVal float64) float64 {
	if v, ctrl := cv.GetProperty(name); ctrl == nil && v != nil {
		return toFloat(v)
	}
	return defaultVal
}

// buildOpenDialogOptions 从 Origami ClassValue 构建 runtime.OpenDialogOptions
func buildOpenDialogOptions(v data.Value) runtime.OpenDialogOptions {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return runtime.OpenDialogOptions{}
	}
	opts := runtime.OpenDialogOptions{
		DefaultDirectory:           getPropString(cv, "DefaultDirectory", ""),
		DefaultFilename:            getPropString(cv, "DefaultFilename", ""),
		Title:                      getPropString(cv, "Title", ""),
		ShowHiddenFiles:            getPropBool(cv, "ShowHiddenFiles", false),
		CanCreateDirectories:       getPropBool(cv, "CanCreateDirectories", false),
		ResolvesAliases:            getPropBool(cv, "ResolvesAliases", false),
		TreatPackagesAsDirectories: getPropBool(cv, "TreatPackagesAsDirectories", false),
	}
	if v, ctrl := cv.GetProperty("Filters"); ctrl == nil && v != nil {
		if av, ok := v.(*data.ArrayValue); ok {
			filters := make([]runtime.FileFilter, 0, len(av.List))
			for _, z := range av.List {
				if z != nil {
					if fcv, ok := z.Value.(*data.ClassValue); ok {
						filters = append(filters, runtime.FileFilter{
							DisplayName: getPropString(fcv, "DisplayName", ""),
							Pattern:     getPropString(fcv, "Pattern", "*"),
						})
					}
				}
			}
			opts.Filters = filters
		}
	}
	return opts
}

// buildSaveDialogOptions 从 Origami ClassValue 构建 runtime.SaveDialogOptions
func buildSaveDialogOptions(v data.Value) runtime.SaveDialogOptions {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return runtime.SaveDialogOptions{}
	}
	opts := runtime.SaveDialogOptions{
		DefaultDirectory:           getPropString(cv, "DefaultDirectory", ""),
		DefaultFilename:            getPropString(cv, "DefaultFilename", ""),
		Title:                      getPropString(cv, "Title", ""),
		ShowHiddenFiles:            getPropBool(cv, "ShowHiddenFiles", false),
		CanCreateDirectories:       getPropBool(cv, "CanCreateDirectories", false),
		TreatPackagesAsDirectories: getPropBool(cv, "TreatPackagesAsDirectories", false),
	}
	if v, ctrl := cv.GetProperty("Filters"); ctrl == nil && v != nil {
		if av, ok := v.(*data.ArrayValue); ok {
			filters := make([]runtime.FileFilter, 0, len(av.List))
			for _, z := range av.List {
				if z != nil {
					if fcv, ok := z.Value.(*data.ClassValue); ok {
						filters = append(filters, runtime.FileFilter{
							DisplayName: getPropString(fcv, "DisplayName", ""),
							Pattern:     getPropString(fcv, "Pattern", "*"),
						})
					}
				}
			}
			opts.Filters = filters
		}
	}
	return opts
}

// buildMessageDialogOptions 从 Origami ClassValue 构建 runtime.MessageDialogOptions
func buildMessageDialogOptions(v data.Value) runtime.MessageDialogOptions {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return runtime.MessageDialogOptions{}
	}
	opts := runtime.MessageDialogOptions{
		Type:          runtime.DialogType(getPropString(cv, "Type", "info")),
		Title:         getPropString(cv, "Title", ""),
		Message:       getPropString(cv, "Message", ""),
		DefaultButton: getPropString(cv, "DefaultButton", ""),
		CancelButton:  getPropString(cv, "CancelButton", ""),
	}
	if v, ctrl := cv.GetProperty("Buttons"); ctrl == nil && v != nil {
		if av, ok := v.(*data.ArrayValue); ok {
			buttons := make([]string, 0, len(av.List))
			for _, z := range av.List {
				if z != nil {
					buttons = append(buttons, toString(z.Value))
				}
			}
			opts.Buttons = buttons
		}
	}
	return opts
}

// ============================================================================
// Wails Application Entry Point
// ============================================================================

// RunApp 从 Origami options ClassValue 启动 Wails 应用
func RunApp(v data.Value) error {
	appOptions, err := buildAppOptions(v)
	if err != nil {
		return err
	}
	return wails.Run(appOptions)
}
