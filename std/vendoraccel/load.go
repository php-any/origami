package vendoraccel

import (
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel"
	"github.com/php-any/origami/std/symfony/clock"
	"github.com/php-any/origami/std/symfony/console"
	eventdispatcher "github.com/php-any/origami/std/symfony/event-dispatcher"
	eventdispatchercontracts "github.com/php-any/origami/std/symfony/event-dispatcher-contracts"
	"github.com/php-any/origami/std/symfony/finder"
	httpfoundation "github.com/php-any/origami/std/symfony/http-foundation"
	httpkernel "github.com/php-any/origami/std/symfony/http-kernel"
	polyfillctype "github.com/php-any/origami/std/symfony/polyfill-ctype"
	polyfillintlgrapheme "github.com/php-any/origami/std/symfony/polyfill-intl-grapheme"
	polyfillintlidn "github.com/php-any/origami/std/symfony/polyfill-intl-idn"
	polyfillintlnormalizer "github.com/php-any/origami/std/symfony/polyfill-intl-normalizer"
	polyfillmbstring "github.com/php-any/origami/std/symfony/polyfill-mbstring"
	polyfillphp84 "github.com/php-any/origami/std/symfony/polyfill-php84"
	polyfillphp85 "github.com/php-any/origami/std/symfony/polyfill-php85"
	polyfillphp86 "github.com/php-any/origami/std/symfony/polyfill-php86"
	polyfilluuid "github.com/php-any/origami/std/symfony/polyfill-uuid"
	"github.com/php-any/origami/std/symfony/process"
	"github.com/php-any/origami/std/symfony/routing"
	sfstring "github.com/php-any/origami/std/symfony/string"
	"github.com/php-any/origami/std/symfony/uid"
	vardumper "github.com/php-any/origami/std/symfony/var-dumper"
)

// Load 仅由 examples/laravel13 调用：按官方 Composer 包逐个注册 vendor 原生层。
// 禁止加入 zy.go / php.Load。
func Load(vm data.VM) {
	// Symfony：一包一子模块
	httpfoundation.Load(vm)
	finder.Load(vm)
	console.Load(vm)
	sfstring.Load(vm)
	routing.Load(vm)
	httpkernel.Load(vm)
	eventdispatchercontracts.Load(vm)
	eventdispatcher.Load(vm)
	process.Load(vm)
	vardumper.Load(vm)
	uid.Load(vm)
	clock.Load(vm)

	polyfillmbstring.Load(vm)
	polyfillctype.Load(vm)
	polyfilluuid.Load(vm)
	polyfillintlgrapheme.Load(vm)
	polyfillintlidn.Load(vm)
	polyfillintlnormalizer.Load(vm)
	polyfillphp84.Load(vm)
	polyfillphp85.Load(vm)
	polyfillphp86.Load(vm)

	// laravel/framework（Illuminate 组件）+ 后续 laravel/* 包
	laravel.Load(vm)

	// 若能定位 vendor，立即标记 polyfill files（serve 预热会再做一遍）
	if wd, err := os.Getwd(); err == nil {
		MarkPolyfills(vm, filepath.Join(wd, "vendor"))
	}
}
