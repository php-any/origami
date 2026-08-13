<?php

namespace tests\php;

/**
 * DateTime / DateTimeInterface 基本测试：
 * - 全局 DateTime 类存在，可实例化；
 * - new DateTime() instanceof \DateTimeInterface 为 true；
 * - format/getTimestamp 方法可调用。
 */

$dt = new \DateTime();

if (!($dt instanceof \DateTimeInterface)) {
    Log::fatal('DateTime instanceof DateTimeInterface 测试失败');
}

// 调用 format / getTimestamp，至少应返回合理类型
$formatted = $dt->format('Y-m-d H:i:s');
if (!is_string($formatted) || strlen($formatted) < 8) {
    Log::fatal('DateTime::format 返回值异常: ' . var_export($formatted, true));
}

$ts = $dt->getTimestamp();
if (!is_int($ts)) {
    Log::fatal('DateTime::getTimestamp 应返回 int，实际: ' . gettype($ts));
}

Log::info('DateTime 基本测试通过');

