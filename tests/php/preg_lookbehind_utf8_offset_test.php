<?php

namespace tests\php;

/**
 * regexp2 的 Index 是 rune，preg 必须转成字节。
 * UTF-8 字符之后 lookbehind 捕获错位会把 @php 编成 <?phphp。
 */

$src = "中@php\n\$a = 1;\n@endphp";
$got = null;
preg_replace_callback('/(?<!@)@php(.*?)@endphp/s', function ($m) use (&$got) {
    $got = "<?php{$m[1]}?>";
    return $got;
}, $src);

if ($got === null) {
    \Log::fatal('未匹配 @php 块');
}
if (str_contains($got, '<?phphp') || str_contains($got, '?>hp')) {
    \Log::fatal('UTF-8 后 @php 捕获错位: '.$got);
}
if ($got !== "<?php\n\$a = 1;\n?>") {
    \Log::fatal('捕获内容错误: '.var_export($got, true));
}

\Log::info('preg_lookbehind_utf8_offset 测试通过');
