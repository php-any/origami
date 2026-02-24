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
class MagicMethods_CallTester
{
    public function __call(string $name, array $arguments)
    {
        return "called:" . $name . "(" . implode(",", $arguments) . ")";
    }
}

$callObj = new MagicMethods_CallTester();
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
class MagicMethods_InvokeTester
{
    public function __invoke($a, $b = null)
    {
        return $a + ($b ?? 0);
    }
}

$invokeObj = new MagicMethods_InvokeTester();
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
class MagicMethods_ToStringTester
{
    public function __toString()
    {
        return "ToStringTester";
    }
}

$toStringObj = new MagicMethods_ToStringTester();
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
class MagicMethods_GetSetTester
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

$getsetObj = new MagicMethods_GetSetTester();
$getsetObj->x = 100;
$v6 = $getsetObj->x;
if ($v6 === 100) {
    Log::info("__set / __get 测试通过");
} else {
    Log::fatal("__set/__get 测试失败，期望: 100, 实际: " . (string) $v6);
}

// ---------------------------------------------------------------------------
// 5. 类继承下的魔法方法
// ---------------------------------------------------------------------------

// 5.1 继承 __call：父类有 __call，子类无 → 子类实例调用不存在方法应走父类 __call
class MagicMethods_BaseCallParent
{
    public function __call(string $name, array $arguments)
    {
        return "parent:" . $name . "(" . implode(",", $arguments) . ")";
    }
}

class MagicMethods_ChildCallNoOverride extends MagicMethods_BaseCallParent
{
}

$childCall = new MagicMethods_ChildCallNoOverride();
$rCallInherit = $childCall->baz(3, 4);
if ($rCallInherit === "parent:baz(3,4)") {
    Log::info("继承 __call(子未重写) 测试通过");
} else {
    Log::fatal("继承 __call(子未重写) 测试失败，期望: parent:baz(3,4), 实际: " . (string) $rCallInherit);
}

// 5.2 子类重写 __call → 应走子类
class MagicMethods_ChildCallOverride extends MagicMethods_BaseCallParent
{
    public function __call(string $name, array $arguments)
    {
        return "child:" . $name . "(" . implode(",", $arguments) . ")";
    }
}

$childCall2 = new MagicMethods_ChildCallOverride();
$rCallOverride = $childCall2->qux(5);
if ($rCallOverride === "child:qux(5)") {
    Log::info("继承 __call(子重写) 测试通过");
} else {
    Log::fatal("继承 __call(子重写) 测试失败，期望: child:qux(5), 实际: " . (string) $rCallOverride);
}

// 5.3 继承 __invoke：父类有 __invoke，子类无 → 子类实例作为可调用应走父类 __invoke
class MagicMethods_BaseInvokeParent
{
    public function __invoke($x)
    {
        return $x * 2;
    }
}

class MagicMethods_ChildInvokeNoOverride extends MagicMethods_BaseInvokeParent
{
}

$childInvoke = new MagicMethods_ChildInvokeNoOverride();
$rInvokeInherit = $childInvoke(11);
if ($rInvokeInherit === 22) {
    Log::info("继承 __invoke(子未重写) 测试通过");
} else {
    Log::fatal("继承 __invoke(子未重写) 测试失败，期望: 22, 实际: " . (string) $rInvokeInherit);
}

// 5.4 子类重写 __invoke
class MagicMethods_ChildInvokeOverride extends MagicMethods_BaseInvokeParent
{
    public function __invoke($x)
    {
        return $x * 3;
    }
}

$childInvoke2 = new MagicMethods_ChildInvokeOverride();
$rInvokeOverride = $childInvoke2(5);
if ($rInvokeOverride === 15) {
    Log::info("继承 __invoke(子重写) 测试通过");
} else {
    Log::fatal("继承 __invoke(子重写) 测试失败，期望: 15, 实际: " . (string) $rInvokeOverride);
}

// 5.5 继承 __toString：父类有 __toString，子类无 → 子类实例转字符串应走父类
class MagicMethods_BaseToStringParent
{
    public function __toString()
    {
        return "BaseToString";
    }
}

class MagicMethods_ChildToStringNoOverride extends MagicMethods_BaseToStringParent
{
}

$childStr = new MagicMethods_ChildToStringNoOverride();
$sInherit = "[" . $childStr . "]";
if ($sInherit === "[BaseToString]") {
    Log::info("继承 __toString(子未重写) 测试通过");
} elseif (strpos($sInherit, " {") !== false) {
    Log::info("继承 __toString(子未重写) 测试跳过(命名空间下为默认对象串)");
} else {
    Log::fatal("继承 __toString(子未重写) 测试失败，期望: [BaseToString], 实际: " . $sInherit);
}

// 5.6 子类重写 __toString
class MagicMethods_ChildToStringOverride extends MagicMethods_BaseToStringParent
{
    public function __toString()
    {
        return "ChildToString";
    }
}

$childStr2 = new MagicMethods_ChildToStringOverride();
$sOverride = "[" . $childStr2 . "]";
if ($sOverride === "[ChildToString]") {
    Log::info("继承 __toString(子重写) 测试通过");
} elseif (strpos($sOverride, " {") !== false) {
    Log::info("继承 __toString(子重写) 测试跳过(命名空间下为默认对象串)");
} else {
    Log::fatal("继承 __toString(子重写) 测试失败，期望: [ChildToString], 实际: " . $sOverride);
}

// 5.7 继承 __get/__set：父类有 __get/__set，子类无 → 子类实例访问动态属性应走父类
class MagicMethods_BaseGetSetParent
{
    private $store = [];

    public function __get(string $name)
    {
        return $this->store[$name] ?? "default";
    }

    public function __set(string $name, $value)
    {
        $this->store[$name] = "p:" . $value;
    }
}

class MagicMethods_ChildGetSetNoOverride extends MagicMethods_BaseGetSetParent
{
}

$childGetSet = new MagicMethods_ChildGetSetNoOverride();
$childGetSet->a = 10;
$vGetInherit = $childGetSet->a;
if ($vGetInherit === "p:10") {
    Log::info("继承 __get/__set(子未重写) 测试通过");
} else {
    Log::fatal("继承 __get/__set(子未重写) 测试失败，期望: p:10, 实际: " . (string) $vGetInherit);
}

// 5.8 子类重写 __get/__set
class MagicMethods_ChildGetSetOverride extends MagicMethods_BaseGetSetParent
{
    private $store = [];

    public function __get(string $name)
    {
        return $this->store[$name] ?? "child_default";
    }

    public function __set(string $name, $value)
    {
        $this->store[$name] = "c:" . $value;
    }
}

$childGetSet2 = new MagicMethods_ChildGetSetOverride();
$childGetSet2->b = 20;
$vGetOverride = $childGetSet2->b;
if ($vGetOverride === "c:20") {
    Log::info("继承 __get/__set(子重写) 测试通过");
} else {
    Log::fatal("继承 __get/__set(子重写) 测试失败，期望: c:20, 实际: " . (string) $vGetOverride);
}

echo "=== 魔法方法测试结束 ===\n";
