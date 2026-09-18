package polyfilluuid

import (
	"path/filepath"

	"github.com/php-any/origami/data"
)

// VendorRelFiles 是 Composer files autoload 会加载的相对路径（相对 vendor/symfony/polyfill-uuid）。
var VendorRelFiles = []string{"bootstrap.php"}

// Load 将 polyfill bootstrap 标为已加载，跳过重复解析（函数已由 php.Load 提供）。
// 实际路径由 vendoraccel.MarkPolyfillFiles 在知道 vendor 根后写入 phpFileCache。
func Load(vm data.VM) {
	_ = vm
}

// MarkLoaded 把 vendor/symfony/polyfill-uuid/bootstrap.php 标进 phpFileCache。
func MarkLoaded(vm data.VM, vendorRoot string) {
	if vendorRoot == "" {
		return
	}
	base := filepath.Join(vendorRoot, "symfony", "polyfill-uuid")
	for _, rel := range VendorRelFiles {
		vm.SetPhpFileCache(filepath.Join(base, rel))
	}
}
