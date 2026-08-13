package fyne

import "github.com/php-any/origami/data"

// Load 将 Fyne 扩展注册到 VM 中
func Load(vm data.VM) {
	// 基础类型（必须先注册，其他组件 extend 它）
	vm.AddClass(newCanvasObjectClass("Fyne\\CanvasObject", nil))

	// 核心类型
	vm.AddClass(NewAppClass())
	vm.AddClass(NewWindowClass())
	vm.AddClass(NewSizeClass())
	vm.AddClass(NewPositionClass())
	vm.AddClass(NewColorClass())
	vm.AddClass(NewTextStyleClass())
	vm.AddClass(NewResourceClass())

	// 容器和布局
	vm.AddClass(NewContainerClass())
	vm.AddClass(NewLayoutClass())

	// Canvas 基元
	vm.AddClass(NewCanvasTextClass())
	vm.AddClass(NewCanvasRectangleClass())
	vm.AddClass(NewCanvasCircleClass())
	vm.AddClass(NewCanvasLineClass())
	vm.AddClass(NewCanvasImageClass())

	// Widget
	vm.AddClass(NewLabelClass())
	vm.AddClass(NewButtonClass())
	vm.AddClass(NewEntryClass())
	vm.AddClass(NewCheckClass())
	vm.AddClass(NewSelectClass())
	vm.AddClass(NewRadioGroupClass())
	vm.AddClass(NewSliderClass())
	vm.AddClass(NewProgressBarClass())
	vm.AddClass(NewProgressBarInfiniteClass())
	vm.AddClass(NewFormClass())
	vm.AddClass(NewCardClass())
	vm.AddClass(NewAccordionClass())
	vm.AddClass(NewTabsClass())
	vm.AddClass(NewToolbarClass())
	vm.AddClass(NewToolbarActionClass())
	vm.AddClass(NewToolbarSeparatorClass())
	vm.AddClass(NewToolbarSpacerClass())

	// BottomTabBar (iOS 风格)
	vm.AddClass(NewBottomTabBarClass())

	// Dialog
	vm.AddClass(NewDialogClass())

	// Theme 常量
	vm.AddClass(NewThemeColorNameClass())
	vm.AddClass(NewThemeSizeNameClass())
	vm.AddClass(NewThemeIconNameClass())
}
