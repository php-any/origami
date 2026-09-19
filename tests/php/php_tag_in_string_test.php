<?php

namespace tests\php;

/**
 * Blade storePhpBlocks 使用 "<?php{$matches[1]}?>"。
 * 双引号字符串里的 <?php / ?> 必须当普通字符，不能切模式或改写成 hp。
 */

$inner = "\n\$a = 1;\n";
$s = "<?php{$inner}?>";
$want = "<?php\n\$a = 1;\n?>";
if ($s !== $want) {
    \Log::fatal('双引号内 PHP 标签被改写: ' . var_export($s, true));
}

$blade = "@php\n\$a = 1;\n@endphp";
$compiled = preg_replace_callback('/(?<!@)@php(.*?)@endphp/s', function ($m) {
    return "<?php{$m[1]}?>";
}, $blade);
$wantCompiled = "<?php\n\$a = 1;\n?>";
if ($compiled !== $wantCompiled) {
    \Log::fatal('preg @php 块编译错误: ' . var_export($compiled, true));
}

$open = '<?php';
$close = '?>';
if ($open !== '<' . '?php') {
    \Log::fatal('单引号 <?php 错误: ' . var_export($open, true));
}
if ($close !== '?' . '>') {
    \Log::fatal('单引号 ?> 错误: ' . var_export($close, true));
}

\Log::info('php_tag_in_string 测试通过');
