package wails

import (
	"testing"

	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

// TestLoad 验证 wails v3 扩展可以在 VM 中正常加载，所有类已注册
func TestLoad(t *testing.T) {
	vm := runtime.NewVM(parser.NewParser())
	Load(vm)

	classes := []string{
		"Wails\\Application",
		"Wails\\Options\\App",
		"Wails\\Options\\RGBA",
		"Wails\\Options\\Windows",
		"Wails\\Options\\Mac",
		"Wails\\Options\\Linux",
		"Wails\\Runtime\\Window",
		"Wails\\Runtime\\Dialog",
		"Wails\\Runtime\\Events",
		"Wails\\Runtime\\Log",
		"Wails\\Runtime\\Browser",
		"Wails\\Runtime\\Screen",
		"Wails\\Runtime\\Environment",
		"Wails\\Menu\\Menu",
		"Wails\\Menu\\MenuItem",
		"Wails\\Menu\\Keys",
		"Wails\\Dialog\\FileFilter",
		"Wails\\Dialog\\OpenDialogOptions",
		"Wails\\Dialog\\SaveDialogOptions",
		"Wails\\Dialog\\MessageDialogOptions",
		// 枚举
		"Wails\\WindowStartState",
		"Wails\\BackdropType",
		"Wails\\Theme",
		"Wails\\WebviewGpuPolicy",
		"Wails\\DialogType",
		"Wails\\LogLevel",
		"Wails\\MenuItemType",
		"Wails\\MacAppearance",
		"Wails\\ImagePosition",
	}

	missing := 0
	for _, name := range classes {
		_, ok := vm.GetClass(name)
		if !ok {
			t.Errorf("class not registered: %s", name)
			missing++
		}
	}

	if missing > 0 {
		t.Fatalf("%d classes missing from VM", missing)
	}
	t.Logf("All %d classes registered successfully", len(classes))
}
