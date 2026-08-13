<?php

namespace tests\php;

/**
 * in_array 严格模式测试：
 * - 第三参 true 时按类型+值比较，空字符串 '' 应在 ['', '_global'] 中匹配
 */

$key = '';
$ok = \in_array($key, ['', '_global'], true);

if (!$ok) {
    Log::fatal('in_array 严格模式测试失败: in_array("", ["", "_global"], true) 应返回 true');
}

Log::info('in_array 严格模式测试通过');
