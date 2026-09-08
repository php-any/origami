<?php

namespace tests\php;

/**
 * include 无 PHP 标签的资源文件应原样输出（PHP 从 HTML 模式开始）。
 * Symfony 异常页会 include favicon.png.base64 这类内容。
 */
ob_start();
include __DIR__ . '/include_raw_asset.txt';
$out = ob_get_clean();
if (!is_string($out) || !str_contains($out, 'iVBORw0KGgo') || !str_contains($out, ',')) {
    Log::fatal('include 无 PHP 标签文件未原样输出: ' . var_export($out, true));
}

Log::info('include 无 PHP 标签资源测试通过');
