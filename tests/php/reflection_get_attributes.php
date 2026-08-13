<?php

namespace tests\php;

echo "=== ReflectionClass::getAttributes 功能测试 ===\n";

// 定义一个测试类（无注解）
class TestReflectionAttributes {
    public $publicProp = "public";
    
    public function publicMethod() {
        return "public method";
    }
}

// 定义另一个测试类
class TestMultipleAttributes {
    public $id;
}

// 先实例化类以确保类被加载到 VM 中
$tempObj = new TestReflectionAttributes();
$tempObj2 = new TestMultipleAttributes();

// 测试 ReflectionClass::getAttributes 基本功能
echo "\n=== 测试 ReflectionClass::getAttributes ===\n";
$reflector = new ReflectionClass("tests\\php\\TestReflectionAttributes");
$attributes = $reflector->getAttributes();

if(is_array($attributes)) {
    Log::info("getAttributes() 返回数组测试通过");
} else {
    Log::fatal("getAttributes() 返回数组测试失败");
}

// 测试无注解的类
echo "\n=== 测试无注解类的 getAttributes ===\n";
class TestNoAttributes {
    public $prop;
}

$tempObj3 = new TestNoAttributes();
$reflector2 = new ReflectionClass("tests\\php\\TestNoAttributes");
$attributes2 = $reflector2->getAttributes();

if(is_array($attributes2) && count($attributes2) == 0) {
    Log::info("无注解类的 getAttributes() 测试通过");
} else {
    Log::fatal("无注解类的 getAttributes() 测试失败，期望空数组，实际: " . count($attributes2));
}

// 测试带 name 参数的 getAttributes（传入 null）
echo "\n=== 测试 getAttributes(name) ===\n";
$reflector3 = new ReflectionClass("tests\\php\\TestMultipleAttributes");
$attributes3 = $reflector3->getAttributes(null);

if(is_array($attributes3)) {
    Log::info("getAttributes(name) 测试通过");
} else {
    Log::fatal("getAttributes(name) 测试失败");
}

// 测试带 flags 参数的 getAttributes
echo "\n=== 测试 getAttributes(name, flags) ===\n";
$attributes4 = $reflector3->getAttributes(null, 0);

if(is_array($attributes4)) {
    Log::info("getAttributes(name, flags) 测试通过");
} else {
    Log::fatal("getAttributes(name, flags) 测试失败");
}

// 测试默认参数（无参数调用）
echo "\n=== 测试 getAttributes() 默认参数 ===\n";
$attributes5 = $reflector3->getAttributes();

if(is_array($attributes5)) {
    Log::info("getAttributes() 默认参数测试通过");
} else {
    Log::fatal("getAttributes() 默认参数测试失败");
}

// 测试只传入 flags 参数（使用默认的 name=null）
echo "\n=== 测试 getAttributes(flags) 使用默认 name ===\n";
$attributes6 = $reflector3->getAttributes(null, 0);

if(is_array($attributes6)) {
    Log::info("getAttributes(flags) 使用默认 name 测试通过");
} else {
    Log::fatal("getAttributes(flags) 使用默认 name 测试失败");
}

echo "\n=== ReflectionClass::getAttributes 功能测试完成 ===\n";
