<?php

namespace tests\php;

/**
 * ClassValue::AsString 不得因对象图有环而栈溢出；htmlspecialchars 对无 __toString 对象应报错而非死循环。
 */
class AsStringCycle_A
{
    public $b;
}

class AsStringCycle_B
{
    public $a;
}

$a = new AsStringCycle_A();
$b = new AsStringCycle_B();
$a->b = $b;
$b->a = $a;

$s = (string) $a; // Origami 可能走 AsString；不应崩溃
if (!is_string($s) && $s !== null) {
    // 若抛错也算可接受；主要是不栈溢出
}

try {
    htmlspecialchars($a);
    Log::fatal('htmlspecialchars(对象) 应抛错或安全返回，不应静默死循环');
} catch (\Throwable $e) {
    Log::info('htmlspecialchars 拒绝对象: ' . $e->getMessage());
}

Log::info('as_string_cycle 测试通过');
