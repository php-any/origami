package wails

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// buildApplicationMenu 根据 PHP 端的 Wails\Menu\Menu 构建一个 Wails 菜单。
func buildApplicationMenu(v data.Value) *application.Menu {
	cv, ok := v.(*data.ClassValue)
	if !ok {
		return nil
	}
	menu := application.NewMenu()
	appendMenuItems(menu, cv)
	return menu
}

// appendMenuItems 把一个 Menu ClassValue 的 _items 全部追加到 Wails 菜单中。
func appendMenuItems(menu *application.Menu, menuCV *data.ClassValue) {
	itemsV, _ := menuCV.GetProperty("_items")
	av, ok := itemsV.(*data.ArrayValue)
	if !ok {
		return
	}
	for _, z := range av.List {
		if z == nil {
			continue
		}
		if miCV, ok := z.Value.(*data.ClassValue); ok {
			addMenuItem(menu, miCV)
		}
	}
}

// addMenuItem 把单个 MenuItem ClassValue 加入到 Wails 菜单。
func addMenuItem(menu *application.Menu, mi *data.ClassValue) {
	typ := getPropString(mi, "Type", "Text")
	label := getPropString(mi, "Label", "")
	accel := getPropString(mi, "Accelerator", "")
	disabled := getPropBool(mi, "Disabled", false)
	hidden := getPropBool(mi, "Hidden", false)
	checked := getPropBool(mi, "Checked", false)

	var onClick data.Value
	if v, ctrl := mi.GetProperty("_onClick"); ctrl == nil {
		onClick = v
	}

	switch typ {
	case "Separator":
		menu.AddSeparator()
	case "Submenu":
		sub := menu.AddSubmenu(label)
		if v, ctrl := mi.GetProperty("_subMenu"); ctrl == nil && v != nil {
			if subCV, ok := v.(*data.ClassValue); ok {
				appendMenuItems(sub, subCV)
			}
		}
	case "Checkbox":
		it := menu.AddCheckbox(label, checked)
		applyMenuItemCommon(it, accel, disabled, hidden, onClick)
	case "Radio":
		it := menu.AddRadio(label, checked)
		applyMenuItemCommon(it, accel, disabled, hidden, onClick)
	default:
		it := menu.Add(label)
		applyMenuItemCommon(it, accel, disabled, hidden, onClick)
	}
}

func applyMenuItemCommon(it *application.MenuItem, accel string, disabled, hidden bool, onClick data.Value) {
	accel = canonicalAccelerator(accel)
	if accel != "" {
		it.SetAccelerator(accel)
	}
	if disabled {
		it.SetEnabled(false)
	}
	if hidden {
		it.SetHidden(true)
	}
	if isCallable(onClick) {
		cb := onClick
		it.OnClick(func(_ *application.Context) {
			invokeCallback(cb)
		})
	}
}

// canonicalAccelerator 把 PHP Keys::* 生成的字符串规范为 Wails 可识别的快捷键格式。
func canonicalAccelerator(accel string) string {
	if strings.TrimSpace(accel) == "" {
		return ""
	}
	probe := application.NewMenuItem("")
	probe.SetAccelerator(accel)
	if canon := probe.GetAccelerator(); canon != "" {
		return canon
	}
	return accel
}
