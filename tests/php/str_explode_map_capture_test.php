<?php

namespace tests\php;

/**
 * Str::of()->explode()->map()->toArray() 对齐，供 blade-capture-directive @capture 使用。
 */
use Illuminate\Support\Str;

// Need laravel autoload for Str - skip if unavailable; use minimal mock path via direct test in examples

$expression = '$content, $logo, $isDarkMode = false';

// Simulate without full Laravel if needed - this test runs under examples
require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

[$name, $args] = Str::contains($expression, ',') ?
    Str::of($expression)->trim()->explode(',', 2)->map(fn ($part) => trim($part))->toArray() :
    [$expression, ''];

echo "name=".json_encode($name)."\n";
echo "args=".json_encode($args)."\n";
echo "name_type=".gettype($name)."\n";
echo "args_type=".gettype($args)."\n";

if ($name !== '$content') {
    Log::fatal('Str explode map toArray name 错误: ' . json_encode($name));
}
if (!is_string($args) || !str_contains($args, '$logo')) {
    Log::fatal('Str explode map toArray args 错误: ' . json_encode($args));
}

$out = "<?php {$name} = (function (\$a) { return function ({$args}) use (\$a) {";
if (!str_contains($out, '<?php $content =')) {
    Log::fatal('插值失败: ' . json_encode($out));
}

Log::info('str_explode_map_capture 测试通过');
