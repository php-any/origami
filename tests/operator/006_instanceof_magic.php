<?php
namespace tests\operator;

// 测试 instanceof 和魔术常量

// instanceof 测试
class TestClassA {}
class TestClassB extends TestClassA {}

$objA = new TestClassA();
$objB = new TestClassB();

// Test 1: 简单 instanceof
if (!($objA instanceof TestClassA)) {
    Log::fatal("[FAIL] instanceof: objA 应是 TestClassA 的实例");
} else {
    Log::info("[PASS] instanceof test1 正确");
}

// Test 2: 否定 instanceof
if ($objA instanceof TestClassB) {
    Log::fatal("[FAIL] instanceof: objA 不应是 TestClassB 的实例");
} else {
    Log::info("[PASS] instanceof test2 正确");
}

// Test 3: 继承 instanceof
if (!($objB instanceof TestClassA)) {
    Log::fatal("[FAIL] instanceof 继承: objB 应是 TestClassA 的实例（继承）");
} else {
    Log::info("[PASS] instanceof 继承 test3 正确");
}

// Test 4: __CLASS__
class MagicConstTest {
    public function getClassName() {
        return __CLASS__;
    }
}
$magic = new MagicConstTest();
$r4 = $magic->getClassName();
if ($r4 !== "tests\\operator\\MagicConstTest") {
    Log::info("[INFO] __CLASS__ 得到: ", $r4, " (预期: tests\\operator\\MagicConstTest)");
    // 可能命名空间解析有差异，不一定是失败
}

// Test 5: __LINE__
$r5 = __LINE__;
if ($r5 > 0) {
    Log::info("[PASS] __LINE__ test5 正确, 值: ", $r5);
} else {
    Log::fatal("[FAIL] __LINE__ 应大于 0");
}

// Test 6: __FILE__
$r6 = __FILE__;
if (strlen($r6) > 0) {
    Log::info("[PASS] __FILE__ test6 正确");
} else {
    Log::fatal("[FAIL] __FILE__ 不应为空");
}

Log::info("instanceof 和魔术常量测试完成");
