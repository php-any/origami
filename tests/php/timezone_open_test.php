<?php

namespace tests\php;

/**
 * timezone_open / timezone_name_get 完整语义测试（当前支持的子集）：
 * - timezone_open 返回 DateTimeZone 实例
 * - timezone_name_get 返回该实例的时区 ID
 * - 非法时区名时 timezone_open 返回 false
 */

$tz = timezone_open('Europe/Paris');
if (!($tz instanceof \DateTimeZone)) {
    Log::fatal('timezone_open("Europe/Paris") 应返回 DateTimeZone 实例');
}

$name = timezone_name_get($tz);
if ($name !== 'Europe/Paris') {
    Log::fatal('timezone_name_get($tz) 期望 "Europe/Paris"，实际: ' . $name);
}

$bad = timezone_open('Not/Exists');
if ($bad !== false) {
    Log::fatal('timezone_open("Not/Exists") 期望 false');
}

Log::info('timezone_open / timezone_name_get 测试通过');

