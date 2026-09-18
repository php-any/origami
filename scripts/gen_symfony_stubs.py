#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""Generate Symfony one-package-one-module stubs."""
from pathlib import Path

ROOT = Path("std/symfony")

# name -> (go_package, composer_name, target_version, notes)
MODULES = {
    "console": ("console", "symfony/console", "v8.1.1", "ArgvInput/Output stubs"),
    "finder": ("finder", "symfony/finder", "v8.1.1", "Finder stub"),
    "string": ("string", "symfony/string", "v8.1.0", "UnicodeString stub"),
    "routing": ("routing", "symfony/routing", "v8.1.0", "Route stubs"),
    "http-kernel": ("httpkernel", "symfony/http-kernel", "v8.1.1", "KernelEvents stubs"),
    "event-dispatcher": ("eventdispatcher", "symfony/event-dispatcher", "v8.1.1", "EventDispatcher stub"),
    "event-dispatcher-contracts": ("eventdispatchercontracts", "symfony/event-dispatcher-contracts", "v3.7.1", "contracts"),
    "process": ("process", "symfony/process", "v8.1.0", "Process stub"),
    "var-dumper": ("vardumper", "symfony/var-dumper", "v8.1.1", "VarDumper stub"),
    "uid": ("uid", "symfony/uid", "v8.1.0", "Uuid stub"),
    "clock": ("clock", "symfony/clock", "v8.1.0", "Clock stub"),
    "polyfill-mbstring": ("polyfillmbstring", "symfony/polyfill-mbstring", "v1.38.2", "skip vendor files"),
    "polyfill-ctype": ("polyfillctype", "symfony/polyfill-ctype", "v1.37.0", "skip vendor files"),
    "polyfill-uuid": ("polyfilluuid", "symfony/polyfill-uuid", "v1.37.0", "skip vendor files"),
    "polyfill-intl-grapheme": ("polyfillintlgrapheme", "symfony/polyfill-intl-grapheme", "v1.41.0", "skip vendor files"),
    "polyfill-intl-idn": ("polyfillintlidn", "symfony/polyfill-intl-idn", "v1.38.1", "skip vendor files"),
    "polyfill-intl-normalizer": ("polyfillintlnormalizer", "symfony/polyfill-intl-normalizer", "v1.38.0", "skip vendor files"),
    "polyfill-php84": ("polyfillphp84", "symfony/polyfill-php84", "v1.38.1", "skip vendor files"),
    "polyfill-php85": ("polyfillphp85", "symfony/polyfill-php85", "v1.41.0", "skip vendor files"),
    "polyfill-php86": ("polyfillphp86", "symfony/polyfill-php86", "v1.41.0", "skip vendor files"),
}

POLYFILL_DIRS = {k for k in MODULES if k.startswith("polyfill-")}

VERSION_TMPL = '''package {pkg}

const ComposerName = "{composer}"

// TargetVersion 对齐 examples/laravel13/composer.lock 中该包的 version。
const TargetVersion = "{version}"
'''

LOAD_STUB = '''package {pkg}

import "github.com/php-any/origami/data"

// Load 注册 {composer} 原生加速（增量实现；未覆盖的 FQCN 仍走 vendor PHP）。
func Load(vm data.VM) {{
	// {notes}
}}
'''

LOAD_POLYFILL = '''package {pkg}

import (
	"path/filepath"

	"github.com/php-any/origami/data"
)

// VendorRelFiles 是 Composer files autoload 会加载的相对路径（相对 vendor/symfony/{dirname}）。
var VendorRelFiles = []string{{"bootstrap.php"}}

// Load 将 polyfill bootstrap 标为已加载，跳过重复解析（函数已由 php.Load 提供）。
// 实际路径由 vendoraccel.MarkPolyfillFiles 在知道 vendor 根后写入 phpFileCache。
func Load(vm data.VM) {{
	_ = vm
}}

// MarkLoaded 把 vendor/symfony/{dirname}/bootstrap.php 标进 phpFileCache。
func MarkLoaded(vm data.VM, vendorRoot string) {{
	if vendorRoot == "" {{
		return
	}}
	base := filepath.Join(vendorRoot, "symfony", "{dirname}")
	for _, rel := range VendorRelFiles {{
		vm.SetPhpFileCache(filepath.Join(base, rel))
	}}
}}
'''


def main() -> None:
    for dirname, (pkg, composer, version, notes) in MODULES.items():
        d = ROOT / dirname
        d.mkdir(parents=True, exist_ok=True)
        (d / "version.go").write_text(
            VERSION_TMPL.format(pkg=pkg, composer=composer, version=version),
            encoding="utf-8",
        )
        if dirname in POLYFILL_DIRS:
            (d / "load.go").write_text(
                LOAD_POLYFILL.format(pkg=pkg, dirname=dirname),
                encoding="utf-8",
            )
        else:
            (d / "load.go").write_text(
                LOAD_STUB.format(pkg=pkg, composer=composer, notes=notes),
                encoding="utf-8",
            )
        print("wrote", d)


if __name__ == "__main__":
    main()
