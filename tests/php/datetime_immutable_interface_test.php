<?php

namespace tests\php;

/**
 * DateTimeImmutable 应实现 DateTimeInterface，子类扩展时不应报抽象方法缺失。
 */

class DTI_Child extends \DateTimeImmutable implements \JsonSerializable
{
    public function jsonSerialize(): string
    {
        return $this->format('c');
    }
}

$d = new DTI_Child('now');
if (!($d instanceof \DateTimeInterface)) {
    Log::fatal('子类应 instanceof DateTimeInterface');
}
$s = $d->format('Y-m-d');
if (!is_string($s) || strlen($s) < 8) {
    Log::fatal('format 应返回日期字符串');
}
$ts = $d->getTimestamp();
if (!is_int($ts) || $ts <= 0) {
    Log::fatal('getTimestamp 应返回正整数');
}

Log::info('datetime_immutable_interface 测试通过');
