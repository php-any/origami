package wails

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// ============================================================================
// 全局 Wails v3 引用
// ============================================================================

var (
	wailsApp        *application.App
	wailsMainWindow *application.WebviewWindow
)

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

// ============================================================================
// 构建 Application 选项
// ============================================================================

func buildAppOptions(v data.Value) (*application.Options, error) {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return nil, errors.New("Wails\\Options\\App instance required")
	}

	name := getPropString(cv, "Title", "Origami App")

	return &application.Options{
		Name: name,
	}, nil
}

// ============================================================================
// 构建 Window 选项
// ============================================================================

func buildWindowOptions(v data.Value) *application.WebviewWindowOptions {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return &application.WebviewWindowOptions{
			Title:  "Origami App",
			Width:  800,
			Height: 600,
		}
	}

	opts := &application.WebviewWindowOptions{
		Title:                      getPropString(cv, "Title", ""),
		Width:                      getPropInt(cv, "Width", 800),
		Height:                     getPropInt(cv, "Height", 600),
		Frameless:                  getPropBool(cv, "Frameless", false),
		DisableResize:              getPropBool(cv, "DisableResize", false),
		MinWidth:                   getPropInt(cv, "MinWidth", 0),
		MinHeight:                  getPropInt(cv, "MinHeight", 0),
		MaxWidth:                   getPropInt(cv, "MaxWidth", 0),
		MaxHeight:                  getPropInt(cv, "MaxHeight", 0),
		AlwaysOnTop:                getPropBool(cv, "AlwaysOnTop", false),
		StartState:                 application.WindowState(getPropInt(cv, "WindowStartState", 0)),
		Hidden:                     getPropBool(cv, "StartHidden", false),
		DefaultContextMenuDisabled: !getPropBool(cv, "EnableDefaultContextMenu", false),
		ContentProtectionEnabled:   getPropBool(cv, "ContentProtection", false),
		EnableFileDrop:             getPropBool(cv, "EnableFileDrop", false),
	}

	// 背景色
	if v, ctrl := cv.GetProperty("BackgroundColour"); ctrl == nil && v != nil {
		if rgba := buildRGBA(v); rgba != nil {
			opts.BackgroundColour = *rgba
		}
	}

	// 平台选项
	if v, ctrl := cv.GetProperty("Windows"); ctrl == nil && v != nil {
		opts.Windows = buildWindowsWindow(v)
	}
	if v, ctrl := cv.GetProperty("Mac"); ctrl == nil && v != nil {
		opts.Mac = buildMacWindow(v)
	}
	if v, ctrl := cv.GetProperty("Linux"); ctrl == nil && v != nil {
		opts.Linux = buildLinuxWindow(v)
	}

	return opts
}

// ============================================================================
// 子选项构建器
// ============================================================================

func buildRGBA(v data.Value) *application.RGBA {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return nil
	}
	return &application.RGBA{
		Red:   uint8(getPropInt(cv, "r", 0)),
		Green: uint8(getPropInt(cv, "g", 0)),
		Blue:  uint8(getPropInt(cv, "b", 0)),
		Alpha: uint8(getPropInt(cv, "a", 255)),
	}
}

func buildWindowsWindow(v data.Value) application.WindowsWindow {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return application.WindowsWindow{}
	}
	return application.WindowsWindow{
		BackdropType:                      application.BackdropType(getPropInt(cv, "BackdropType", 0)),
		DisableIcon:                       getPropBool(cv, "DisableWindowIcon", false),
		Theme:                             application.Theme(getPropInt(cv, "Theme", 0)),
		DisableFramelessWindowDecorations: getPropBool(cv, "DisableFramelessWindowDecorations", false),
		ResizeDebounceMS:                  uint16(getPropInt(cv, "ResizeDebounceMS", 0)),
		EnableSwipeGestures:               getPropBool(cv, "EnableSwipeGestures", false),
	}
}

func buildMacWindow(v data.Value) application.MacWindow {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return application.MacWindow{}
	}
	opts := application.MacWindow{
		Appearance:                      application.MacAppearanceType(getPropString(cv, "Appearance", "")),
		EnableFraudulentWebsiteWarnings: getPropBool(cv, "EnableFraudulentWebsiteWarnings", false),
	}
	if v, ctrl := cv.GetProperty("TitleBar"); ctrl == nil && v != nil {
		if tbCV, ok := v.(*data.ClassValue); ok {
			opts.TitleBar = application.MacTitleBar{
				AppearsTransparent:   getPropBool(tbCV, "AppearsTransparent", false),
				HideTitle:            getPropBool(tbCV, "HideTitle", false),
				Hide:                 getPropBool(tbCV, "Hide", false),
				FullSizeContent:      getPropBool(tbCV, "FullSizeContent", false),
				UseToolbar:           getPropBool(tbCV, "UseToolbar", false),
				HideToolbarSeparator: getPropBool(tbCV, "HideToolbarSeparator", false),
			}
		}
	}
	return opts
}

func buildLinuxWindow(v data.Value) application.LinuxWindow {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return application.LinuxWindow{}
	}
	return application.LinuxWindow{
		WindowIsTranslucent: getPropBool(cv, "WindowIsTranslucent", false),
		WebviewGpuPolicy:    application.WebviewGpuPolicy(getPropInt(cv, "WebviewGpuPolicy", 0)),
	}
}

// ============================================================================
// 对话框选项构建
// ============================================================================

func buildOpenFileDialogBuilder(v data.Value) *application.OpenFileDialogStruct {
	if wailsApp == nil {
		return nil
	}
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return wailsApp.Dialog.OpenFile()
	}
	d := wailsApp.Dialog.OpenFile().
		SetTitle(getPropString(cv, "Title", "")).
		SetDirectory(getPropString(cv, "DefaultDirectory", ""))

	if v, ctrl := cv.GetProperty("Filters"); ctrl == nil && v != nil {
		if av, ok := v.(*data.ArrayValue); ok {
			for _, z := range av.List {
				if z != nil {
					if fcv, ok := z.Value.(*data.ClassValue); ok {
						d.AddFilter(
							getPropString(fcv, "DisplayName", ""),
							getPropString(fcv, "Pattern", "*"),
						)
					}
				}
			}
		}
	}
	return d
}

func buildOpenDirectoryDialogBuilder(v data.Value) *application.OpenFileDialogStruct {
	if wailsApp == nil {
		return nil
	}
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return wailsApp.Dialog.OpenFile().CanChooseDirectories(true).CanChooseFiles(false)
	}
	return wailsApp.Dialog.OpenFile().
		SetTitle(getPropString(cv, "Title", "")).
		SetDirectory(getPropString(cv, "DefaultDirectory", "")).
		CanChooseDirectories(true).
		CanChooseFiles(false)
}

func buildSaveFileDialogBuilder(v data.Value) *application.SaveFileDialogStruct {
	if wailsApp == nil {
		return nil
	}
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return wailsApp.Dialog.SaveFile()
	}
	d := wailsApp.Dialog.SaveFile().
		SetMessage(getPropString(cv, "Title", "")).
		SetDirectory(getPropString(cv, "DefaultDirectory", "")).
		SetFilename(getPropString(cv, "DefaultFilename", ""))

	if v, ctrl := cv.GetProperty("Filters"); ctrl == nil && v != nil {
		if av, ok := v.(*data.ArrayValue); ok {
			for _, z := range av.List {
				if z != nil {
					if fcv, ok := z.Value.(*data.ClassValue); ok {
						d.AddFilter(
							getPropString(fcv, "DisplayName", ""),
							getPropString(fcv, "Pattern", "*"),
						)
					}
				}
			}
		}
	}
	return d
}

func buildMessageDialog(v data.Value) *application.MessageDialog {
	if wailsApp == nil {
		return nil
	}
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return wailsApp.Dialog.Info()
	}

	var d *application.MessageDialog
	dialogType := getPropString(cv, "Type", "info")
	switch dialogType {
	case "warning":
		d = wailsApp.Dialog.Warning()
	case "error":
		d = wailsApp.Dialog.Error()
	case "question":
		d = wailsApp.Dialog.Question()
	default:
		d = wailsApp.Dialog.Info()
	}

	d.SetTitle(getPropString(cv, "Title", "")).
		SetMessage(getPropString(cv, "Message", ""))

	// 添加按钮
	var defBtn, cancelBtn *application.Button
	if v, ctrl := cv.GetProperty("Buttons"); ctrl == nil && v != nil {
		if av, ok := v.(*data.ArrayValue); ok {
			buttons := make([]*application.Button, 0, len(av.List))
			for _, z := range av.List {
				if z != nil {
					label := toString(z.Value)
					btn := d.AddButton(label)
					buttons = append(buttons, btn)
					if label == getPropString(cv, "DefaultButton", "") {
						defBtn = btn
					}
					if label == getPropString(cv, "CancelButton", "") {
						cancelBtn = btn
					}
				}
			}
			if len(buttons) > 0 {
				d.AddButtons(buttons)
			}
		}
	}

	if defBtn != nil {
		d.SetDefaultButton(defBtn)
	}
	if cancelBtn != nil {
		d.SetCancelButton(cancelBtn)
	}

	return d
}

// ============================================================================
// Wails Application Entry Point
// ============================================================================

func RunApp(v data.Value) error {
	appOpts, err := buildAppOptions(v)
	if err != nil {
		return err
	}

	wailsApp = application.New(*appOpts)

	winOpts := buildWindowOptions(v)
	if winOpts == nil {
		winOpts = &application.WebviewWindowOptions{
			Title:  "Origami App",
			Width:  800,
			Height: 600,
		}
	}

	wailsMainWindow = wailsApp.Window.NewWithOptions(*winOpts)

	return wailsApp.Run()
}
