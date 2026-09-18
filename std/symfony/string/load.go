package sfstring

import (
	"github.com/php-any/origami/data"
)

// VendorRelFiles 故意为空：不要把 Resources/functions.php 写入 phpFileCache，
// 否则 Composer files autoload 会跳过该文件，u()/b()/s() 将不会被定义。
var VendorRelFiles []string

// Load 注册 Symfony String 原生层：抽象类桩（instanceof / 常量）+ 三个具体类。
// 必须五类都 AddClass；u()/b()/s() 仍由 vendor Resources/functions.php 定义。
func Load(vm data.VM) {
	vm.AddClass(NewAbstractStringClass())
	vm.AddClass(NewAbstractUnicodeStringClass())
	unicodeStringClass = NewUnicodeStringClass().(*strClass)
	byteStringClass = NewByteStringClass().(*strClass)
	codePointStringClass = NewCodePointStringClass().(*strClass)
	vm.AddClass(unicodeStringClass)
	vm.AddClass(byteStringClass)
	vm.AddClass(codePointStringClass)
}
