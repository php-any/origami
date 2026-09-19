<?php

namespace tests\php;

/**
 * 诊断 Filament tables index 上 @php...@endphp 的 preg 捕获，
 * 确认是否被切成 @p + hp 导致 compile 出 <?phphp。
 */

$src = file_get_contents(__DIR__ . '/../../examples/laravel13/vendor/filament/tables/resources/views/index.blade.php');
if ($src === false || $src === '') {
    \Log::fatal('读不到 index.blade.php');
}

$heads = [];
preg_replace_callback('/(?<!@)@php(.*?)@endphp/s', function ($m) use (&$heads) {
    $heads[] = [
        'm0' => substr($m[0], 0, 12),
        'm1' => substr($m[1], 0, 16),
        'm0end' => substr($m[0], -10),
    ];
    return $m[0];
}, $src);

if (count($heads) < 2) {
    \Log::fatal('应匹配多个 @php 块, 实际 '.count($heads));
}

foreach ($heads as $i => $h) {
    if (!str_starts_with($h['m0'], '@php')) {
        \Log::fatal("块 {$i} m0 不以 @php 开头: ".var_export($h, true));
    }
    if (str_starts_with($h['m1'], 'hp')) {
        \Log::fatal("块 {$i} 捕获以 hp 开头: ".var_export($h, true));
    }
}

$stored = [];
$mid = preg_replace_callback('/(?<!@)@php(.*?)@endphp/s', function ($m) use (&$stored) {
    $stored[] = "<?php{$m[1]}?>";
    return '@__raw_block_'.(count($stored) - 1).'__@';
}, $src);

foreach ($stored as $i => $s) {
    if (str_contains($s, '<?phphp') || str_contains($s, '?>hp')) {
        \Log::fatal("存储块 {$i} 已含 hp: ".substr($s, 0, 40).' ... '.substr($s, -20));
    }
    if (!str_starts_with($s, '<?php') || !str_ends_with($s, '?>')) {
        \Log::fatal("存储块 {$i} 标签不对: ".substr($s, 0, 20).' / '.substr($s, -10));
    }
}

\Log::info('blade_php_blocks_preg 测试通过, n='.count($heads));
