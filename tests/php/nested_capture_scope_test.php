<?php

namespace tests\php;

/**
 * 方法内嵌套闭包：list 赋值的 $name 须与双引号 {$name} 插值为同一局部变量。
 */
class NestedCaptureScope_Host
{
    public function register()
    {
        return function (string $expression): string {
            [$name, $arguments] = str_contains($expression, ',') ?
                array_map('trim', explode(',', $expression, 2)) :
                [$expression, ''];

            return "<?php {$name} = (function (\$args) { return function ({$arguments}) use (\$args) {; ?>";
        };
    }
}

$host = new NestedCaptureScope_Host();
$fn = $host->register();
$out = $fn('$content, $logo, $isDarkMode = false');
echo "out=".json_encode($out)."\n";
if (!str_contains($out, '<?php $content =')) {
    Log::fatal('方法内闭包 {$name} 未绑定到 list 赋值: ' . json_encode($out));
}
Log::info('nested_capture_scope 测试通过');
