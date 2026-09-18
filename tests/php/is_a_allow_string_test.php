<?php
namespace tests\php;

/**
 * is_a 第三个参数 allow_string 命名参数必须生效。
 */
class IsAAllowString_Base {}
class IsAAllowString_Child extends IsAAllowString_Base {}

$a = is_a(IsAAllowString_Child::class, IsAAllowString_Base::class, true);
$b = is_a(IsAAllowString_Child::class, IsAAllowString_Base::class, allow_string: true);
$c = is_a(IsAAllowString_Child::class, IsAAllowString_Child::class, allow_string: true);
$d = is_a(IsAAllowString_Child::class, IsAAllowString_Base::class); // false without allow

if (!$a) {
    Log::fatal('位置参数 allow_string=true 失败');
}
if (!$b) {
    Log::fatal('命名参数 allow_string: true 失败');
}
if (!$c) {
    Log::fatal('同类名 allow_string 失败');
}
if ($d) {
    Log::fatal('未允许字符串时应为 false');
}

Log::info('is_a_allow_string 测试通过');
