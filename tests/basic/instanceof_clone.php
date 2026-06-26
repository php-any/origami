<?php

echo "=== instanceof 和 clone 测试 ===\n";

// 定义测试类
class Animal {
    public $name;
    public function __construct($name) {
        $this->name = $name;
    }
    public function speak() {
        return $this->name . " speaks";
    }
}

class Dog extends Animal {
    public $breed;
    public function __construct($name, $breed) {
        $this->name = $name;
        $this->breed = $breed;
    }
    public function fetch() {
        return $this->name . " fetches";
    }
}

interface Walker {
    public function walk();
}

class Cat extends Animal implements Walker {
    public function walk() {
        return $this->name . " walks";
    }
}

// instanceof 基本用法
$dog = new Dog("Buddy", "Labrador");
if ($dog instanceof Dog) {
    Log::info("instanceof 同类测试通过");
} else {
    Log::fatal("instanceof 同类测试失败");
}

// instanceof 父类检查
if ($dog instanceof Animal) {
    Log::info("instanceof 父类测试通过");
} else {
    Log::fatal("instanceof 父类测试失败");
}

// instanceof 负面测试
$animal = new Animal("Generic");
if (!($animal instanceof Dog)) {
    Log::info("instanceof 非子类测试通过");
} else {
    Log::fatal("instanceof 非子类测试失败");
}

// instanceof 接口检查
$cat = new Cat("Whiskers");
if ($cat instanceof Walker) {
    Log::info("instanceof 接口测试通过");
} else {
    Log::fatal("instanceof 接口测试失败");
}

// instanceof 字面量检查
if ($cat instanceof Animal) {
    Log::info("instanceof 继承+接口 测试通过");
} else {
    Log::fatal("instanceof 继承+接口 测试失败");
}

// instanceof 变量
$obj = new Dog("Rex", "Husky");
$className = "Animal";
// instanceof 需要类名直接写
if ($obj instanceof Animal) {
    Log::info("instanceof 变量对象测试通过");
} else {
    Log::fatal("instanceof 变量对象测试失败");
}

// clone 基本用法
$original = new Animal("Original");
$cloned = clone $original;
if ($cloned->name == "Original" && $cloned !== $original) {
    Log::info("clone 基本测试通过");
} else {
    Log::fatal("clone 基本测试失败");
}

// clone 后修改不影响原对象
$cloned->name = "Cloned";
if ($original->name == "Original" && $cloned->name == "Cloned") {
    Log::info("clone 独立副本测试通过");
} else {
    Log::fatal("clone 独立副本测试失败, original='{$original->name}', cloned='{$cloned->name}'");
}

// clone 对象的方法可用
$cloned2 = clone $original;
$result = $cloned2->speak();
if ($result == "Original speaks") {
    Log::info("clone 对象方法测试通过");
} else {
    Log::fatal("clone 对象方法测试失败, 实际 '{$result}'");
}

echo "=== instanceof 和 clone 测试完成 ===\n";
