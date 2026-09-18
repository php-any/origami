<?php

namespace tests\php;

/**
 * Finder 递归扫描嵌套目录时，应返回 Symfony SplFileInfo，且 getRelativePathname 含目录前缀。
 */
if (!is_file(dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php')) {
    Log::info('skip: 缺少 vendor 依赖，跳过测试');
    return;
}

require dirname(__DIR__, 2) . '/examples/laravel13/vendor/autoload.php';

use Symfony\Component\Finder\Finder;
use Symfony\Component\Finder\SplFileInfo as FinderSplFileInfo;

$dir = dirname(__DIR__, 2) . '/examples/laravel13/app/Filament/Resources';
if (!is_dir($dir)) {
    Log::fatal('Resources 目录不存在');
}

$found = null;
foreach (Finder::create()->files()->in($dir)->name('AdminResource.php')->sortByName() as $file) {
    $found = $file;
    break;
}

if ($found === null) {
    Log::fatal('未找到 AdminResource.php');
}

if (!($found instanceof FinderSplFileInfo)) {
    Log::fatal('期望 Symfony Finder SplFileInfo，实际: ' . get_class($found));
}

$rel = $found->getRelativePathname();
if ($rel === '' || !str_contains(str_replace('\\', '/', $rel), 'Admins/')) {
    Log::fatal('getRelativePathname 应含 Admins/ 前缀，实际: [' . $rel . ']');
}

$path = $found->getPathname();
if ($path === '' || !str_contains($path, 'AdminResource.php')) {
    Log::fatal('getPathname 无效: [' . $path . ']');
}

Log::info('finder_nested_relative_pathname 测试通过 rel=' . $rel);
