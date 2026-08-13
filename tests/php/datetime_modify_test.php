<?php

namespace tests\php;

/**
 * DateTime::modify 相对修饰符：Carbon/Telescope 会用 "4 second"、"+4 seconds" 等
 */

$dt = new \DateTime('2024-01-01 12:00:00');
$dt->modify('4 second');
if ($dt->format('Y-m-d H:i:s') !== '2024-01-01 12:00:04') {
    Log::fatal('modify("4 second") 失败: ' . $dt->format('Y-m-d H:i:s'));
}

$dt2 = new \DateTime('2024-01-01 12:00:00');
$dt2->modify('+4 seconds');
if ($dt2->format('Y-m-d H:i:s') !== '2024-01-01 12:00:04') {
    Log::fatal('modify("+4 seconds") 失败: ' . $dt2->format('Y-m-d H:i:s'));
}

$dt3 = new \DateTime('2024-01-01 12:00:00');
$dt3->modify('-1 day');
if ($dt3->format('Y-m-d H:i:s') !== '2023-12-31 12:00:00') {
    Log::fatal('modify("-1 day") 失败: ' . $dt3->format('Y-m-d H:i:s'));
}

$ts = strtotime('4 second', strtotime('2024-01-01 12:00:00'));
if ($ts !== strtotime('2024-01-01 12:00:04')) {
    Log::fatal('strtotime("4 second", base) 失败');
}

Log::info('DateTime::modify 相对修饰符测试通过');
