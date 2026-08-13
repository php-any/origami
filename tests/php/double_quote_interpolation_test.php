<?php

namespace tests\php;

/**
 * 双引号字符串裸变量插值测试：
 * - "$var" 应把变量值插入字符串
 */

$help = "usage";
$out = "$help\n\n";

if ($out !== "usage\n\n") {
    Log::fatal('双引号插值测试失败: expected usage+newlines, got ' . var_export($out, true));
}

Log::info('双引号插值测试通过');
