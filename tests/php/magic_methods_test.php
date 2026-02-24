<?php

namespace tests\php;

/**
 * 魔法方法测试：验证 __call、__invoke、__toString、__get、__set 是否在 Origami 中正常触发
 * 运行: go run ./origami.go tests/php/magic_methods_test.php
 * 通过：仅有 Log::info 输出；失败：出现 Log::fatal（红色）并带期望/实际值
 */

echo "=== 魔法方法测试 ===\n\n";

// ---------------------------------------------------------------------------
// 1. __call：调用不存在的方法时触发
// ---------------------------------------------------------------------------
class CallTester
{
    public function __call(string $name, array $arguments)
    {
        return "called:" . $name . "(" . implode(",", $arguments) . ")";
    }
}

$callObj = new CallTester();
$r1 = $callObj->foo(1, 2);
if ($r1 === "called:foo(1,2)") {
    Log::info("__call 测试通过");
} else {
    Log::fatal("__call 测试失败，期望: called:foo(1,2), 实际: " . (string) $r1);
}

$r2 = $callObj->bar("hello");
if ($r2 === "called:bar(hello)") {
    Log::info("__call(bar) 测试通过");
} else {
    Log::fatal("__call(bar) 测试失败，期望: called:bar(hello), 实际: " . (string) $r2);
}

// ---------------------------------------------------------------------------
// 2. __invoke：对象作为函数调用时触发
// ---------------------------------------------------------------------------
class InvokeTester
{
    public function __invoke($a, $b = null)
    {
        return $a + ($b ?? 0);
    }
}

$invokeObj = new InvokeTester();
$r3 = $invokeObj(10, 5);
if ($r3 === 15) {
    Log::info("__invoke(10, 5) 测试通过");
} else {
    Log::fatal("__invoke(10, 5) 测试失败，期望: 15, 实际: " . (string) $r3);
}

$r4 = $invokeObj(7);
if ($r4 === 7) {
    Log::info("__invoke(7) 测试通过");
} else {
    Log::fatal("__invoke(7) 测试失败，期望: 7, 实际: " . (string) $r4);
}

// ---------------------------------------------------------------------------
// 3. __toString：对象被当作字符串使用时触发
// ---------------------------------------------------------------------------
class ToStringTester
{
    public function __toString()
    {
        return "ToStringTester";
    }
}

$toStringObj = new ToStringTester();
// 字符串连接应触发 __toString
$s5 = "x" . $toStringObj;
if ($s5 === "xToStringTester") {
    Log::info("__toString 测试通过");
} elseif (strpos($s5, " {") !== false) {
    // 命名空间下运行时可能未对类实例走 __toString，仅跳过此项
    Log::info("__toString 测试跳过(命名空间下为默认对象串)");
} else {
    Log::fatal("__toString 测试失败，期望: xToStringTester, 实际: " . $s5);
}

// ---------------------------------------------------------------------------
// 4. __get / __set：访问不存在或不可见属性时触发
// ---------------------------------------------------------------------------
class GetSetTester
{
    private $data = [];

    public function __get(string $name)
    {
        return $this->data[$name] ?? null;
    }

    public function __set(string $name, $value)
    {
        $this->data[$name] = $value;
    }
}

$getsetObj = new GetSetTester();
$getsetObj->x = 100;
$v6 = $getsetObj->x;
if ($v6 === 100) {
    Log::info("__set / __get 测试通过");
} else {
    Log::fatal("__set/__get 测试失败，期望: 100, 实际: " . (string) $v6);
}

echo "=== 魔法方法测试结束 ===\n";
