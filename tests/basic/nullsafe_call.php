<?php

namespace tests\basic;

echo "=== 空安全调用操作符 ?-> 测试 ===\n";

// 定义一个测试类
class TestClass {
    public $property = "test property";
    
    public function method() {
        return "test method";
    }
    
    public function methodWithParam($param) {
        return "param: " . $param;
    }
    
    public function getNested() {
        // 返回自身作为嵌套对象
        return $this;
    }
    
    public function getNull() {
        return null;
    }
}

// 测试非空对象的方法调用
$obj = new TestClass();
$result1 = $obj?->method();
if($result1 == "test method") {
    Log::info("非空对象 ?-> 方法调用测试通过");
} else {
    Log::fatal("非空对象 ?-> 方法调用测试失败，期望: test method, 实际: " . $result1);
}

// 测试非空对象的属性访问
$result2 = $obj?->property;
if($result2 == "test property") {
    Log::info("非空对象 ?-> 属性访问测试通过");
} else {
    Log::fatal("非空对象 ?-> 属性访问测试失败，期望: test property, 实际: " . $result2);
}

// 测试 null 对象的方法调用
$nullObj = null;
$result3 = $nullObj?->method();
if($result3 === null) {
    Log::info("null 对象 ?-> 方法调用测试通过");
} else {
    Log::fatal("null 对象 ?-> 方法调用测试失败，期望: null, 实际: " . ($result3 === null ? "null" : $result3));
}

// 测试 null 对象的属性访问
$result4 = $nullObj?->property;
if($result4 === null) {
    Log::info("null 对象 ?-> 属性访问测试通过");
} else {
    Log::fatal("null 对象 ?-> 属性访问测试失败，期望: null, 实际: " . ($result4 === null ? "null" : $result4));
}

// 测试链式调用 - 非空链
$obj2 = new TestClass();
$result5 = $obj2?->getNested()?->method();
if($result5 == "test method") {
    Log::info("链式 ?-> 调用 - 非空链测试通过");
} else {
    Log::fatal("链式 ?-> 调用 - 非空链测试失败，期望: test method, 实际: " . $result5);
}

// 测试链式调用 - 中间为 null
$result6 = $obj2?->getNull()?->method();
if($result6 === null) {
    Log::info("链式 ?-> 调用 - 中间为 null 测试通过");
} else {
    Log::fatal("链式 ?-> 调用 - 中间为 null 测试失败，期望: null, 实际: " . ($result6 === null ? "null" : $result6));
}

// 测试链式调用 - 开头为 null
$nullObj2 = null;
$result7 = $nullObj2?->getNested()?->method();
if($result7 === null) {
    Log::info("链式 ?-> 调用 - 开头为 null 测试通过");
} else {
    Log::fatal("链式 ?-> 调用 - 开头为 null 测试失败，期望: null, 实际: " . ($result7 === null ? "null" : $result7));
}

// 测试带参数的方法调用
$obj3 = new TestClass();
$result8 = $obj3?->methodWithParam("test");
if($result8 == "param: test") {
    Log::info("?-> 带参数方法调用测试通过");
} else {
    Log::fatal("?-> 带参数方法调用测试失败，期望: param: test, 实际: " . $result8);
}

// 测试 null 对象的带参数方法调用
$nullObj3 = null;
$result9 = $nullObj3?->methodWithParam("test");
if($result9 === null) {
    Log::info("null 对象 ?-> 带参数方法调用测试通过");
} else {
    Log::fatal("null 对象 ?-> 带参数方法调用测试失败，期望: null, 实际: " . ($result9 === null ? "null" : $result9));
}

// 测试三链式调用（使用属性访问）
$obj4 = new TestClass();
$result10 = $obj4?->getNested()?->method();
if($result10 == "test method") {
    Log::info("三链式 ?-> 调用测试通过");
} else {
    Log::fatal("三链式 ?-> 调用测试失败，期望: test method, 实际: " . $result10);
}

// 测试混合使用 -> 和 ?->
$obj5 = new TestClass();
$result11 = $obj5->getNested()?->method();
if($result11 == "test method") {
    Log::info("混合使用 -> 和 ?-> 测试通过");
} else {
    Log::fatal("混合使用 -> 和 ?-> 测试失败，期望: test method, 实际: " . $result11);
}

// 测试在表达式中使用
$obj6 = new TestClass();
$result12 = ($obj6?->property) . " suffix";
if($result12 == "test property suffix") {
    Log::info("表达式中使用 ?-> 测试通过");
} else {
    Log::fatal("表达式中使用 ?-> 测试失败，期望: test property suffix, 实际: " . $result12);
}

echo "=== 空安全调用操作符 ?-> 测试完成 ===\n";

