<?php
/**
 * ::class 编译期类名解析：与官方 PHP 一致，不依赖类已加载、不触发 autoload。
 */

namespace Foo {
    final class Probe {
        public static function short(): string {
            return Bar::class;
        }

        public static function qualified(): string {
            return \Foo\Bar::class;
        }

        public static function relative(): string {
            return Baz\Qux::class;
        }
    }
}

namespace Foo\Baz {
    final class Qux {}
}

namespace UseAliasProbe {
    use Foo\Baz\Qux as ImportedQux;

    final class Consumer {
        public static function viaUse(): string {
            return ImportedQux::class;
        }
    }
}

namespace {
    $expect = 'Foo\\Bar';

    $short = Foo\Probe::short();
    if ($short !== $expect) {
        Log::fatal('未限定 Bar::class 期望 ' . $expect . '，得到: ' . $short);
    }

    $qualified = Foo\Probe::qualified();
    if ($qualified !== $expect) {
        Log::fatal('\\Foo\\Bar::class 期望 ' . $expect . '，得到: ' . $qualified);
    }

    $relExpect = 'Foo\\Baz\\Qux';
    $rel = Foo\Probe::relative();
    if ($rel !== $relExpect) {
        Log::fatal('相对 Baz\\Qux::class 期望 ' . $relExpect . '，得到: ' . $rel);
    }

    class_alias('stdClass', 'Foo\\Bar');
    $after = Foo\Probe::short();
    if ($after !== $expect) {
        Log::fatal('加载后 short 仍应为编译期 FQCN，得到: ' . $after);
    }

    $viaUse = UseAliasProbe\Consumer::viaUse();
    if ($viaUse !== $relExpect) {
        Log::fatal('use 别名 ::class 期望 ' . $relExpect . '，得到: ' . $viaUse);
    }

    if (stdClass::class !== 'stdClass') {
        Log::fatal('全局 stdClass::class 期望 stdClass，得到: ' . stdClass::class);
    }

    Log::info('class_const_compile_time 测试通过');
}
