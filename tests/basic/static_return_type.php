<?php

namespace tests\basic;

echo "=== static 返回类型测试 ===\n";

// 测试函数返回类型为 static
class BaseClass {
    public function getInstance(): static {
        return $this;
    }
    
    public function create(): static {
        return new static();
    }
}

class DerivedClass extends BaseClass {
}

// 测试基本 static 返回类型
$base = new BaseClass();
$instance = $base->getInstance();
if($instance instanceof BaseClass) {
    Log::info("BaseClass::getInstance() static 返回类型测试通过");
} else {
    Log::fatal("BaseClass::getInstance() static 返回类型测试失败");
}

// 测试派生类的 static 返回类型
$derived = new DerivedClass();
$instance2 = $derived->getInstance();
if($instance2 instanceof DerivedClass) {
    Log::info("DerivedClass::getInstance() static 返回类型测试通过");
} else {
    Log::fatal("DerivedClass::getInstance() static 返回类型测试失败");
}

// 测试 create() 方法
$instance3 = $base->create();
if($instance3 instanceof BaseClass) {
    Log::info("BaseClass::create() static 返回类型测试通过");
} else {
    Log::fatal("BaseClass::create() static 返回类型测试失败");
}

$instance4 = $derived->create();
if($instance4 instanceof DerivedClass) {
    Log::info("DerivedClass::create() static 返回类型测试通过");
} else {
    Log::fatal("DerivedClass::create() static 返回类型测试失败");
}

// 测试可空 static 返回类型
class NullableStatic {
    public function getInstanceOrNull(): ?static {
        return $this;
    }
}

$nullable = new NullableStatic();
$result = $nullable->getInstanceOrNull();
if($result instanceof NullableStatic) {
    Log::info("可空 static 返回类型测试通过");
} else {
    Log::fatal("可空 static 返回类型测试失败");
}

// 测试联合类型中的 static
class UnionStatic {
    public function getInstanceOrString(): static|string {
        return $this;
    }
}

$union = new UnionStatic();
$result2 = $union->getInstanceOrString();
if($result2 instanceof UnionStatic) {
    Log::info("联合类型中的 static 返回类型测试通过");
} else {
    Log::fatal("联合类型中的 static 返回类型测试失败");
}

echo "=== static 返回类型测试完成 ===\n";

