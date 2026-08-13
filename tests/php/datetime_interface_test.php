<?php

namespace tests\php;

/**
 * DateTimeInterface 接口测试：
 * - 全局 DateTime 实现 \DateTimeInterface
 * - instanceof \DateTimeInterface 为 true
 * - DateTimeInterface 类型提示可以接受 DateTime 实例
 */

class DateTimeInterface_Acceptor
{
    public static function accept(\DateTimeInterface $dt): void
    {
        Log::info('DateTimeInterface 参数类型检查通过: ' . get_class($dt));
    }
}

$dt = new \DateTime();

if (!($dt instanceof \DateTimeInterface)) {
    Log::fatal('DateTimeInterface instanceof 测试失败：期望 true');
}

DateTimeInterface_Acceptor::accept($dt);

Log::info('DateTimeInterface 测试通过');

