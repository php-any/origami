<?php

namespace tests\php;

/**
 * PHP 双引号字符串 \u{XXXX} Unicode 转义。
 */

$zap = "\u{26A1}";
if (strlen($zap) < 1 || $zap === '\\u{26A1}') {
    Log::fatal('\\\\u{26A1} 未解析为 Unicode 字符, got=' . bin2hex($zap) . ' raw=' . var_export($zap, true));
}

$pattern = '/' . $zap . '[\x{FE0E}\x{FE0F}]?/u';
$out = preg_replace($pattern, '', 'admin.login');
if ($out !== 'admin.login') {
    Log::fatal('preg_replace with unicode pattern failed: ' . var_export($out, true) . ' type=' . gettype($out));
}

// 模拟 Livewire Finder::normalizeName 尾部
$name = preg_replace($pattern, '', 'admin.login');
$name = str_replace('/', '.', $name);
if ($name !== 'admin.login') {
    Log::fatal('normalize 模拟失败: ' . var_export($name, true));
}

Log::info('Unicode 字符串转义测试通过 hex=' . bin2hex($zap));
