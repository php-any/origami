<?php

namespace tests\php;

/**
 * enum_exists：不存在的名字与普通类都应为 false；函数必须已注册。
 */
if (!function_exists('enum_exists')) {
	Log::fatal('enum_exists 未注册');
}
if (enum_exists('stdClass') !== false) {
	Log::fatal('stdClass 不是 enum');
}
if (enum_exists('EnumExists_NoSuchType') !== false) {
	Log::fatal('不存在的类型应为 false');
}
Log::info('enum_exists 测试通过');
