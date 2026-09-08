<?php

namespace tests\php;

/**
 * 验证 ip2long / long2ip，以及 Symfony IpUtils 使用的 sprintf('%032b', ip2long(...))。
 */
$loopback = ip2long('127.0.0.1');
if ($loopback !== 2130706433) {
	Log::fatal('ip2long(127.0.0.1) 应为 2130706433, 实际: ' . var_export($loopback, true));
}
$broadcast = ip2long('255.255.255.255');
if ($broadcast !== 4294967295) {
	Log::fatal('ip2long(255.255.255.255) 应为 4294967295, 实际: ' . var_export($broadcast, true));
}
if (ip2long('0.0.0.0') !== 0) {
	Log::fatal('ip2long(0.0.0.0) 应为 0');
}
if (ip2long('not-an-ip') !== false) {
	Log::fatal('非法地址应返回 false');
}
if (ip2long('::1') !== false) {
	Log::fatal('IPv6 应返回 false');
}
if (long2ip(2130706433) !== '127.0.0.1') {
	Log::fatal('long2ip 回环失败: ' . var_export(long2ip(2130706433), true));
}
if (long2ip(4294967295) !== '255.255.255.255') {
	Log::fatal('long2ip 广播失败: ' . var_export(long2ip(4294967295), true));
}
$bits = sprintf('%032b', ip2long('127.0.0.1'));
if ($bits !== '01111111000000000000000000000001') {
	Log::fatal("sprintf('%%032b', ip2long) 失败: " . var_export($bits, true));
}
$cidr = 0 === substr_compare(
	sprintf('%032b', ip2long('127.0.0.1')),
	sprintf('%032b', ip2long('127.0.0.0')),
	0,
	8
);
if ($cidr !== true) {
	Log::fatal('127.0.0.1 应落在 127.0.0.0/8');
}
Log::info('ip2long/long2ip 测试通过');
