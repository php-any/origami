<?php

namespace tests\php;

/**
 * Str::substr 负 length（Blade callCustomDirective 剥括号依赖）。
 */

$value = '($content, $logo, $isDarkMode = false)';
$inner = \Illuminate\Support\Str::substr($value, 1, -1);
if ($inner !== '$content, $logo, $isDarkMode = false') {
    Log::fatal('Str::substr neg length: [' . $inner . ']');
}

$c = collect(['a', 'bb', 'ccc'])->sortBy(fn ($v) => strlen($v), descending: true)->values()->all();
if ($c !== ['ccc', 'bb', 'a']) {
    Log::fatal('sortBy named: ' . json_encode($c));
}

Log::info('illuminate_str_substr_sortby 测试通过');
