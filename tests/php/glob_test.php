<?php

namespace tests\php;

/**
 * glob 函数测试：通配符匹配、GLOB_BRACE 与 GLOB_ONLYDIR。
 */

$files = glob('tests/php/*.php');
if (!is_array($files)) {
    Log::fatal('glob 应返回数组');
}
if (count($files) === 0) {
    Log::fatal('glob tests/php/*.php 应匹配到文件');
}
$hasBin2hex = false;
foreach ($files as $file) {
    if (str_ends_with($file, 'bin2hex_test.php')) {
        $hasBin2hex = true;
        break;
    }
}
if (!$hasBin2hex) {
    Log::fatal('glob 未匹配到 bin2hex_test.php');
}

$none = glob('tests/php/__no_such_glob_pattern_12345__/*.php');
if (!is_array($none) || count($none) !== 0) {
    Log::fatal('glob 无匹配时应返回空数组');
}

$brace = glob('tests/php/{bin2hex,http_build_query}_test.php', GLOB_BRACE);
if (!is_array($brace) || count($brace) !== 2) {
    Log::fatal('glob GLOB_BRACE 展开失败，数量: ' . count($brace));
}

$dirs = glob('tests/*', GLOB_ONLYDIR);
if (!is_array($dirs) || count($dirs) === 0) {
    Log::fatal('glob GLOB_ONLYDIR 应返回目录');
}
foreach ($dirs as $dir) {
    if (!is_dir($dir)) {
        Log::fatal('glob GLOB_ONLYDIR 结果包含非目录: ' . $dir);
    }
}

Log::info('glob 函数测试通过');
