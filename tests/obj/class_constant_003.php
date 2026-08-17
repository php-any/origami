<?php
namespace tests\obj;

// 测试 self::class 和 static::class 在继承中的行为
class ParentClass {
    public function getSelfClass() {
        return self::class;
    }
    
    public function getStaticClass() {
        return static::class;
    }
}

class ChildClass extends ParentClass {
}

$parent = new ParentClass();
$child = new ChildClass();

$parentSelf = $parent->getSelfClass();
$parentStatic = $parent->getStaticClass();
$childSelf = $child->getSelfClass();
$childStatic = $child->getStaticClass();

if($parentSelf == "tests\\obj\\ParentClass") {
    Log::info("父类 self::class 正确");
} else {
    Log::fatal("父类 self::class 错误: " . $parentSelf);
}

if($parentStatic == "tests\\obj\\ParentClass") {
    Log::info("父类 static::class 正确");
} else {
    Log::fatal("父类 static::class 错误: " . $parentStatic);
}

// self::class 是词法绑定，返回定义方法的类（ParentClass），即使由子类实例调用
if($childSelf == "tests\\obj\\ParentClass") {
    Log::info("子类 self::class 正确（返回定义方法的类 ParentClass）");
} else {
    Log::fatal("子类 self::class 错误: " . $childSelf);
}

// static::class 使用后期静态绑定，返回实际调用时的类（ChildClass）
if($childStatic == "tests\\obj\\ChildClass") {
    Log::info("子类 static::class 正确（返回实际调用时的类 ChildClass）");
} else {
    Log::fatal("子类 static::class 错误: " . $childStatic);
}
