<?php

namespace tests\date;

/**
 * date：DateTime / DateTimeImmutable 格式化。
 */

$d = new \DateTime('2020-01-02 03:04:05', new \DateTimeZone('UTC'));
if ($d->format('Y-m-d') !== '2020-01-02') {
    Log::fatal('DateTime format 失败: ' . $d->format('Y-m-d'));
}
$i = new \DateTimeImmutable('2020-01-02', new \DateTimeZone('UTC'));
if ($i->format('Y-m-d') !== '2020-01-02') {
    Log::fatal('DateTimeImmutable format 失败');
}
if (!($d instanceof \DateTimeInterface) || !($i instanceof \DateTimeInterface)) {
    Log::fatal('DateTimeInterface instanceof 失败');
}

Log::info('date 扩展测试通过');
