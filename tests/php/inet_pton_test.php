<?php

namespace tests\php;

/**
 * 验证 inet_pton / inet_ntop。
 */
$packed = inet_pton('127.0.0.1');
if ($packed === false || strlen($packed) !== 4) {
	Log::fatal('inet_pton(127.0.0.1) 应返回 4 字节');
}
$back = inet_ntop($packed);
if ($back !== '127.0.0.1') {
	Log::fatal('inet_ntop 回环失败: ' . var_export($back, true));
}
if (inet_pton('not-an-ip') !== false) {
	Log::fatal('非法 IP 应返回 false');
}
$v6 = inet_pton('::1');
if ($v6 === false || strlen($v6) !== 16) {
	Log::fatal('inet_pton(::1) 应返回 16 字节');
}
$v6back = inet_ntop($v6);
if ($v6back !== '::1') {
	Log::fatal('inet_ntop(::1) 失败: ' . var_export($v6back, true));
}
Log::info('inet_pton/ntop 测试通过');
