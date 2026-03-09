<?php

namespace tests\php;

$tz = timezone_open('Europe/Paris');

if (!($tz instanceof \DateTimeZone)) {
    Log::fatal('not instance of DateTimeZone, got ' . get_class($tz));
}

// 暂时只验证 timezone_open 返回对象是否正确，name 由 timezone_name_get 专门测试
Log::info('timezone_open debug: class=' . get_class($tz));

