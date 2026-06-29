package wails

import "github.com/php-any/origami/data"

// Load 将 Wails 扩展注册到 VM 中
//
// 注册所有 Wails 相关的类，包括:
//   - 配置类 (Options\App, Options\RGBA, 平台选项, 等)
//   - 枚举类 (WindowStartState, Theme, DialogType, LogLevel, 等)
//   - 运行时类 (Runtime\Window, Runtime\Dialog, Runtime\Events, Runtime\Log, 等)
//   - 菜单类 (Menu\Menu, Menu\MenuItem, Menu\Keys)
//   - 对话框数据类型 (Dialog\FileFilter, Dialog\OpenDialogOptions, 等)
//   - 应用程序类 (Application)
//
// 用法:
//
//	import wailsExt "github.com/php-any/origami-wails"
//	wailsExt.Load(vm)
func Load(vm data.VM) {
	// ── Application ──
	vm.AddClass(NewApplicationClass())

	// ── Options ──
	vm.AddClass(NewOptionsAppClass())
	vm.AddClass(NewRGBAClass())
	vm.AddClass(NewSingleInstanceLockClass())
	vm.AddClass(NewDragAndDropClass())
	vm.AddClass(NewDebugClass())
	vm.AddClass(NewAssetServerClass())
	vm.AddClass(NewWindowsOptionsClass())
	vm.AddClass(NewMacOptionsClass())
	vm.AddClass(NewMacTitleBarClass())
	vm.AddClass(NewMacAboutInfoClass())
	vm.AddClass(NewLinuxOptionsClass())
	vm.AddClass(NewSystemTrayClass())

	// ── 枚举 ──
	vm.AddClass(NewWindowStartStateClass())
	vm.AddClass(NewBackdropTypeClass())
	vm.AddClass(NewThemeClass())
	vm.AddClass(NewWebviewGpuPolicyClass())
	vm.AddClass(NewDialogTypeClass())
	vm.AddClass(NewLogLevelClass())
	vm.AddClass(NewMenuItemTypeClass())
	vm.AddClass(NewMacAppearanceClass())
	vm.AddClass(NewImagePositionClass())

	// ── Runtime ──
	vm.AddClass(NewRuntimeWindowClass())
	vm.AddClass(NewRuntimeDialogClass())
	vm.AddClass(NewRuntimeEventsClass())
	vm.AddClass(NewRuntimeLogClass())
	vm.AddClass(NewRuntimeBrowserClass())
	vm.AddClass(NewRuntimeScreenClass())
	vm.AddClass(NewRuntimeEnvironmentClass())

	// ── Menu ──
	vm.AddClass(NewMenuClass())
	vm.AddClass(NewMenuItemClass())
	vm.AddClass(NewMenuKeysClass())

	// ── Dialog 类型 ──
	vm.AddClass(NewFileFilterClass())
	vm.AddClass(NewOpenDialogOptionsClass())
	vm.AddClass(NewSaveDialogOptionsClass())
	vm.AddClass(NewMessageDialogOptionsClass())
}
