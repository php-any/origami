package wails

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"testing/fstest"

	"github.com/php-any/origami/data"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
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

	opts := &application.Options{
		Name: name,
	}
	// macOS 默认关闭最后一个窗口后 app 不退出（停留在 Dock）。
	// 对单窗口示例而言这并不符合预期，这里让关闭最后窗口即退出程序。
	opts.Mac.ApplicationShouldTerminateAfterLastWindowClosed = true
	return opts, nil
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
				AppearsTransparent: getPropBool(tbCV, "TitlebarAppearsTransparent",
					getPropBool(tbCV, "AppearsTransparent", false)),
				Hide: getPropBool(tbCV, "HideTitleBar",
					getPropBool(tbCV, "Hide", false)),
				HideTitle:            getPropBool(tbCV, "HideTitle", false),
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
// 资源服务辅助
// ============================================================================

// htmlAssetServer 把一段内联 HTML 封装为内存文件系统并交给 Wails 的
// BundledAssetFileServer 提供服务。这样 /wails/runtime.js 也会被自动提供。
// 注意：runtime.js 是 ES 模块，前端须用 <script type="module"> + import 引入。
func htmlAssetServer(html string) application.AssetOptions {
	fsys := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte(html)},
	}
	return application.AssetOptions{
		Handler:        application.BundledAssetFileServer(fsys),
		DisableLogging: true,
	}
}

// defaultHTML 在用户未提供 HTML / 资源目录 / URL 时生成一个友好的默认页面，
// 避免窗口显示 Wails 的内置默认页。
func defaultHTML(title string) string {
	if title == "" {
		title = "Origami + Wails v3"
	}
	return fmt.Sprintf(`<!doctype html>
<html lang="zh">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%[1]s</title>
<style>
  :root { color-scheme: light dark; }
  body { margin:0; font-family: -apple-system, "Segoe UI", system-ui, sans-serif;
         display:flex; min-height:100vh; align-items:center; justify-content:center;
         background: linear-gradient(135deg,#1e293b,#0f172a); color:#e2e8f0; }
  .card { text-align:center; padding:48px 64px; border-radius:20px;
          background:rgba(255,255,255,.05); box-shadow:0 20px 60px rgba(0,0,0,.4); }
  h1 { margin:0 0 12px; font-size:28px; }
  p { margin:0; opacity:.7; }
  .logo { font-size:56px; margin-bottom:16px; }
</style>
</head>
<body>
  <div class="card">
    <div class="logo">🚀</div>
    <h1>%[1]s</h1>
    <p>Powered by Origami + Wails v3</p>
  </div>
</body>
</html>`, title)
}

// getCallbackProp 读取 ClassValue 上的一个可调用闭包属性。
func getCallbackProp(cv *data.ClassValue, name string) data.Value {
	if cv == nil {
		return nil
	}
	if v, ctrl := cv.GetProperty(name); ctrl == nil && isCallable(v) {
		return v
	}
	return nil
}

// ============================================================================
// Wails Application Entry Point
// ============================================================================

// RunApp 构建并启动 Wails 应用，阻塞直到应用退出。
// ctx 用于在 Wails 回调线程中执行 PHP 闭包（生命周期 / 事件回调）。
func RunApp(ctx data.Context, v data.Value) error {
	wailsRootCtx = ctx

	cv, _ := v.(*data.ClassValue)

	appOpts, err := buildAppOptions(v)
	if err != nil {
		return err
	}

	// ── 资源服务：资源目录 / 内联 HTML / 外部 URL ──
	html, assetDir, customURL := "", "", ""
	if cv != nil {
		html = getPropString(cv, "HTML", "")
		assetDir = getPropString(cv, "AssetDir", "")
		customURL = getPropString(cv, "URL", "")
	}
	switch {
	case assetDir != "":
		appOpts.Assets = application.AssetOptions{
			Handler: application.BundledAssetFileServer(os.DirFS(assetDir)),
		}
	case customURL == "":
		if html == "" {
			html = defaultHTML(appOpts.Name)
		}
		appOpts.Assets = htmlAssetServer(html)
	}

	// ── 生命周期：onShutdown / onBeforeClose ──
	if shutdownCB := getCallbackProp(cv, "_onShutdown"); shutdownCB != nil {
		appOpts.OnShutdown = func() { invokeCallback(shutdownCB) }
	}
	if beforeCloseCB := getCallbackProp(cv, "_onBeforeClose"); beforeCloseCB != nil {
		appOpts.ShouldQuit = func() bool {
			// onBeforeClose 返回 true 表示阻止关闭 → ShouldQuit 取反
			return !toBool(invokeCallback(beforeCloseCB))
		}
	}

	wailsApp = application.New(*appOpts)

	// 注册脚本里在 run 之前排队的事件监听器
	flushPendingEventListeners()

	winOpts := buildWindowOptions(v)
	if winOpts == nil {
		winOpts = &application.WebviewWindowOptions{Title: "Origami App", Width: 800, Height: 600}
	}
	if customURL != "" {
		winOpts.URL = customURL
	} else {
		winOpts.URL = "/"
	}

	wailsMainWindow = wailsApp.Window.NewWithOptions(*winOpts)

	// ── 菜单 ──
	// macOS 上 WebviewWindow.SetMenu 是空操作，必须通过 Application Menu 设置，
	// 否则自定义菜单项和快捷键都不会生效。
	if cv != nil {
		if mv, ctrl := cv.GetProperty("_menu"); ctrl == nil && mv != nil {
			if menu := buildApplicationMenu(mv); menu != nil {
				switch runtime.GOOS {
				case "darwin":
					wailsApp.Menu.Set(menu)
				default:
					wailsMainWindow.SetMenu(menu)
				}
			}
		}
	}

	// ── 生命周期：onStartup / onDomReady ──
	if startupCB := getCallbackProp(cv, "_onStartup"); startupCB != nil {
		wailsApp.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) {
			invokeCallback(startupCB)
		})
	}
	if domReadyCB := getCallbackProp(cv, "_onDomReady"); domReadyCB != nil {
		wailsMainWindow.OnWindowEvent(events.Common.WindowRuntimeReady, func(*application.WindowEvent) {
			invokeCallback(domReadyCB)
		})
	}

	return wailsApp.Run()
}
