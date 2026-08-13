package wails

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

func TestBuildMacWindowMapsTitleBarProperties(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	Load(vm)

	ctx := vm.CreateContext(nil)
	tbGV, _ := NewMacTitleBarClass().GetValue(ctx)
	tb, _ := tbGV.(*data.ClassValue)
	tb.SetVM(vm)
	tb.SetProperty("TitlebarAppearsTransparent", data.NewBoolValue(true))
	tb.SetProperty("FullSizeContent", data.NewBoolValue(true))
	tb.SetProperty("HideTitleBar", data.NewBoolValue(false))

	macGV, _ := NewMacOptionsClass().GetValue(ctx)
	mac, _ := macGV.(*data.ClassValue)
	mac.SetVM(vm)
	mac.SetProperty("TitleBar", tb)

	win := buildMacWindow(mac)
	if !win.TitleBar.AppearsTransparent {
		t.Fatal("AppearsTransparent should be true")
	}
	if !win.TitleBar.FullSizeContent {
		t.Fatal("FullSizeContent should be true")
	}
	if win.TitleBar.Hide {
		t.Fatal("Hide should be false")
	}
}
