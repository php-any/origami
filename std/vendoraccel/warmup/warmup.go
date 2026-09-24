// Package warmup 负责 vendor classmap / polyfill 预热。
//
// 单独成包是为了让 std/laravel 侧（serve 启动预热）能复用，而 std/vendoraccel
// 又依赖 std/laravel —— 放在同一包会形成 import cycle。
package warmup

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/runtime"
	polyfillctype "github.com/php-any/origami/std/symfony/polyfill-ctype"
	polyfillintlgrapheme "github.com/php-any/origami/std/symfony/polyfill-intl-grapheme"
	polyfillintlidn "github.com/php-any/origami/std/symfony/polyfill-intl-idn"
	polyfillintlnormalizer "github.com/php-any/origami/std/symfony/polyfill-intl-normalizer"
	polyfillmbstring "github.com/php-any/origami/std/symfony/polyfill-mbstring"
	polyfillphp84 "github.com/php-any/origami/std/symfony/polyfill-php84"
	polyfillphp85 "github.com/php-any/origami/std/symfony/polyfill-php85"
	polyfillphp86 "github.com/php-any/origami/std/symfony/polyfill-php86"
	polyfilluuid "github.com/php-any/origami/std/symfony/polyfill-uuid"
)

// MarkPolyfills 把 vendor/symfony/polyfill-*/bootstrap.php 标进 phpFileCache，跳过重复解析。
func MarkPolyfills(vm data.VM, vendorRoot string) {
	if vendorRoot == "" {
		return
	}
	polyfillmbstring.MarkLoaded(vm, vendorRoot)
	polyfillctype.MarkLoaded(vm, vendorRoot)
	polyfilluuid.MarkLoaded(vm, vendorRoot)
	polyfillintlgrapheme.MarkLoaded(vm, vendorRoot)
	polyfillintlidn.MarkLoaded(vm, vendorRoot)
	polyfillintlnormalizer.MarkLoaded(vm, vendorRoot)
	polyfillphp84.MarkLoaded(vm, vendorRoot)
	polyfillphp85.MarkLoaded(vm, vendorRoot)
	polyfillphp86.MarkLoaded(vm, vendorRoot)
}

// WarmupVendorClassmap 仅预热 vendor/ 下 classmap 中尚未 native 的类文件。
// artisan 一次性命令默认关闭；serve 或 ORIGAMI_LARAVEL_PRELOAD=1 时打开。
func WarmupVendorClassmap(vm data.VM, projectRoot string) {
	if projectRoot == "" {
		return
	}
	vendorRoot := filepath.Join(projectRoot, "vendor")
	MarkPolyfills(vm, vendorRoot)

	classmapPath := filepath.Join(vendorRoot, "composer", "autoload_classmap.php")
	if _, err := os.Stat(classmapPath); err != nil {
		return
	}

	rtm, ok := vm.(*runtime.VM)
	if !ok {
		return
	}

	// 解析 classmap：轻量扫描 `'FQCN' => $baseDir . '/rel'` 或绝对路径字符串
	raw, err := os.ReadFile(classmapPath)
	if err != nil {
		return
	}
	content := string(raw)
	vendorAbs, _ := filepath.Abs(vendorRoot)

	// 粗解析：找 => 后面的路径拼接；完整 PHP 执行太重，这里用启发式提取相对路径片段
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if !strings.Contains(line, "=>") {
			continue
		}
		// 提取 FQCN
		fqcn := extractQuoted(line)
		if fqcn == "" || (!strings.HasPrefix(fqcn, "Illuminate\\") && !strings.HasPrefix(fqcn, "Symfony\\")) {
			continue
		}
		if skipClassmapWarmup(fqcn, "") {
			continue
		}
		// 已有原生类则跳过
		if _, exists := vm.GetClass(fqcn); exists {
			continue
		}
		// 提取路径中 vendor 之后的片段
		rel := extractClassmapRelPath(line)
		if rel == "" {
			continue
		}
		if skipClassmapWarmup(fqcn, rel) {
			continue
		}
		full := filepath.Clean(filepath.Join(vendorRoot, rel))
		fullAbs, err := filepath.Abs(full)
		if err != nil {
			continue
		}
		if !strings.HasPrefix(fullAbs, vendorAbs) {
			continue
		}
		if _, err := os.Stat(fullAbs); err != nil {
			continue
		}
		if rtm.GetPhpFileCache(fullAbs) {
			continue
		}
		_, _, ctl := rtm.ParseFileCached(fullAbs)
		if ctl != nil {
			continue
		}
		// 执行类定义（LoadAndRun 会标记 phpFileCache）
		_, _ = rtm.LoadAndRun(fullAbs)
	}
}

// skipClassmapWarmup 跳过测试辅助类：serve 预热不应加载 PHPUnit 才能解析的 Illuminate\Testing。
func skipClassmapWarmup(fqcn, rel string) bool {
	n := strings.ToLower(fqcn)
	p := strings.ToLower(filepath.ToSlash(rel))
	if strings.Contains(n, "\\testing\\") || strings.Contains(n, "\\tests\\") {
		return true
	}
	if strings.Contains(p, "/testing/") || strings.Contains(p, "/tests/") {
		return true
	}
	return strings.Contains(n, "phpunit")
}

func extractQuoted(line string) string {
	i := strings.Index(line, "'")
	if i < 0 {
		i = strings.Index(line, `"`)
		if i < 0 {
			return ""
		}
		j := strings.Index(line[i+1:], `"`)
		if j < 0 {
			return ""
		}
		return line[i+1 : i+1+j]
	}
	j := strings.Index(line[i+1:], "'")
	if j < 0 {
		return ""
	}
	return line[i+1 : i+1+j]
}

func extractClassmapRelPath(line string) string {
	// 常见：'Foo\\Bar' => $baseDir . '/symfony/http-foundation/Request.php',
	idx := strings.Index(line, "/symfony/")
	if idx < 0 {
		idx = strings.Index(line, "/laravel/")
	}
	if idx < 0 {
		idx = strings.Index(line, "/illuminate/")
	}
	if idx < 0 {
		return ""
	}
	rest := line[idx+1:] // drop leading /
	end := strings.IndexAny(rest, "'\"")
	if end < 0 {
		return ""
	}
	return filepath.FromSlash(rest[:end])
}

// ShouldWarmup 判断是否应在当前进程做 vendor 预热。
func ShouldWarmup() bool {
	v := strings.TrimSpace(os.Getenv("ORIGAMI_LARAVEL_PRELOAD"))
	return v == "1" || strings.EqualFold(v, "true") || strings.EqualFold(v, "yes")
}
