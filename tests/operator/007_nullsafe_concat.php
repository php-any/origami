<?php
namespace tests\operator;

// 测试空安全调用 ?-> 和字符串连接优先级

// Test 1: 字符串连接 (.) 优先级
$r1 = "Hello" . " " . "World";
if ($r1 !== "Hello World") {
    Log::fatal("[FAIL] 字符串连接: 应得 'Hello World', 实际: ", $r1);
} else {
    Log::info("[PASS] 字符串连接 test1 正确");
}

// Test 2: . 与 + 比较（PHP 8+：. 优先级低于 +，故 "Number: " . 1 + 2 => "Number: 3"）
$r2 = "Number: " . 1 + 2;
if ($r2 !== "Number: 3") {
    Log::fatal("[FAIL] . 与 + 优先级: 应得 'Number: 3', 实际: ", $r2);
} else {
    Log::info("[PASS] . 与 + 优先级 test2 正确");
}
// Test 3: 点号与三目
$_ = function($v) { return $v ?? '-'; };
$r3b = "Value: " . $_(null);
// 期望: "Value: -"
$expected = "Value: -";
if ($r3b !== $expected) {
    Log::fatal("[FAIL] 字符串拼接: 应得 'Value: -', 实际: ", $r3b);
} else {
    Log::info("[PASS] 字符串拼接 null 合并 test3 正确");
}

// Test 4: 多级属性链
class ChainNode {
    public $next = null;
    public $value;
    public function __construct($val, $next=null) {
        $this->value = $val;
        $this->next = $next;
    }
}
$node = new ChainNode("a", new ChainNode("b", new ChainNode("c")));
$r4 = $node->next->next->value;
if ($r4 !== "c") {
    Log::fatal("[FAIL] 链式属性: 应得 'c', 实际: ", $r4);
} else {
    Log::info("[PASS] 链式属性 test4 正确");
}

Log::info("空安全和字符串连接测试完成");
