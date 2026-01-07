<?php

// 测试 ReflectionClass 功能

// 定义一个测试类
class TestClass {
    public $publicProp = "public";
    protected $protectedProp = "protected";
    private $privateProp = "private";
    
    public function __construct($value = "default") {
        $this->publicProp = $value;
    }
    
    public function publicMethod() {
        return "public method";
    }
    
    protected function protectedMethod() {
        return "protected method";
    }
    
    private function privateMethod() {
        return "private method";
    }
}

// 定义继承类
class ChildClass extends TestClass {
    public $childProp = "child";
    
    public function childMethod() {
        return "child method";
    }
}

// 测试 ReflectionClass::__construct 和 getName
echo "=== 测试 ReflectionClass::getName ===\n";
$reflector = new ReflectionClass("TestClass");
$className = $reflector->getName();
echo "类名: " . $className . "\n";
if ($className == "TestClass") {
    echo "✓ getName() 测试通过\n";
} else {
    echo "✗ getName() 测试失败\n";
}

// 测试 ReflectionClass::hasMethod
echo "\n=== 测试 ReflectionClass::hasMethod ===\n";
$hasPublic = $reflector->hasMethod("publicMethod");
$hasProtected = $reflector->hasMethod("protectedMethod");
$hasPrivate = $reflector->hasMethod("privateMethod");
$hasNotExist = $reflector->hasMethod("notExistMethod");

echo "hasMethod('publicMethod'): " . ($hasPublic ? "true" : "false") . "\n";
echo "hasMethod('protectedMethod'): " . ($hasProtected ? "true" : "false") . "\n";
echo "hasMethod('privateMethod'): " . ($hasPrivate ? "true" : "false") . "\n";
echo "hasMethod('notExistMethod'): " . ($hasNotExist ? "true" : "false") . "\n";

if ($hasPublic && $hasProtected && $hasPrivate && !$hasNotExist) {
    echo "✓ hasMethod() 测试通过\n";
} else {
    echo "✗ hasMethod() 测试失败\n";
}

// 测试 ReflectionClass::hasProperty
echo "\n=== 测试 ReflectionClass::hasProperty ===\n";
$hasPublicProp = $reflector->hasProperty("publicProp");
$hasProtectedProp = $reflector->hasProperty("protectedProp");
$hasPrivateProp = $reflector->hasProperty("privateProp");
$hasNotExistProp = $reflector->hasProperty("notExistProp");

echo "hasProperty('publicProp'): " . ($hasPublicProp ? "true" : "false") . "\n";
echo "hasProperty('protectedProp'): " . ($hasProtectedProp ? "true" : "false") . "\n";
echo "hasProperty('privateProp'): " . ($hasPrivateProp ? "true" : "false") . "\n";
echo "hasProperty('notExistProp'): " . ($hasNotExistProp ? "true" : "false") . "\n";

if ($hasPublicProp && $hasProtectedProp && $hasPrivateProp && !$hasNotExistProp) {
    echo "✓ hasProperty() 测试通过\n";
} else {
    echo "✗ hasProperty() 测试失败\n";
}

// 测试 ReflectionClass::getMethods
echo "\n=== 测试 ReflectionClass::getMethods ===\n";
$methods = $reflector->getMethods();
echo "方法数量: " . count($methods) . "\n";
foreach ($methods as $method) {
    echo "  - " . $method . "\n";
}

// 测试 ReflectionClass::getProperties
echo "\n=== 测试 ReflectionClass::getProperties ===\n";
$properties = $reflector->getProperties();
echo "属性数量: " . count($properties) . "\n";
foreach ($properties as $prop) {
    echo "  - " . $prop . "\n";
}

// 测试 ReflectionClass::isSubclassOf
echo "\n=== 测试 ReflectionClass::isSubclassOf ===\n";
$childReflector = new ReflectionClass("ChildClass");
$isSubclass = $childReflector->isSubclassOf("TestClass");
echo "ChildClass 是否是 TestClass 的子类: " . ($isSubclass ? "true" : "false") . "\n";
if ($isSubclass) {
    echo "✓ isSubclassOf() 测试通过\n";
} else {
    echo "✗ isSubclassOf() 测试失败\n";
}

// 测试 ReflectionClass::getParentClass
echo "\n=== 测试 ReflectionClass::getParentClass ===\n";
$parentClass = $childReflector->getParentClass();
echo "父类: " . $parentClass . "\n";
if ($parentClass == "TestClass") {
    echo "✓ getParentClass() 测试通过\n";
} else {
    echo "✗ getParentClass() 测试失败\n";
}

// 测试 ReflectionClass::isInstance
echo "\n=== 测试 ReflectionClass::isInstance ===\n";
$testObj = new TestClass("test");
$isInstance = $reflector->isInstance($testObj);
echo "testObj 是否是 TestClass 的实例: " . ($isInstance ? "true" : "false") . "\n";
if ($isInstance) {
    echo "✓ isInstance() 测试通过\n";
} else {
    echo "✗ isInstance() 测试失败\n";
}

// 测试 ReflectionClass::isInstantiable
echo "\n=== 测试 ReflectionClass::isInstantiable ===\n";
$isInstantiable = $reflector->isInstantiable();
echo "TestClass 是否可实例化: " . ($isInstantiable ? "true" : "false") . "\n";
if ($isInstantiable) {
    echo "✓ isInstantiable() 测试通过\n";
} else {
    echo "✗ isInstantiable() 测试失败\n";
}

// 测试 ReflectionClass::newInstance
echo "\n=== 测试 ReflectionClass::newInstance ===\n";
$newObj = $reflector->newInstance("new value");
echo "新实例的 publicProp: " . $newObj->publicProp . "\n";
if ($newObj->publicProp == "new value") {
    echo "✓ newInstance() 测试通过\n";
} else {
    echo "✗ newInstance() 测试失败\n";
}

// 测试 ReflectionClass::newInstanceWithoutConstructor
echo "\n=== 测试 ReflectionClass::newInstanceWithoutConstructor ===\n";
$newObjWithoutCtor = $reflector->newInstanceWithoutConstructor();
echo "新实例的 publicProp: " . $newObjWithoutCtor->publicProp . "\n";
if ($newObjWithoutCtor->publicProp == "public") {
    echo "✓ newInstanceWithoutConstructor() 测试通过\n";
} else {
    echo "✗ newInstanceWithoutConstructor() 测试失败\n";
}

// 测试使用对象创建 ReflectionClass
echo "\n=== 测试使用对象创建 ReflectionClass ===\n";
$obj = new TestClass("object test");
$reflectorFromObj = new ReflectionClass($obj);
$classNameFromObj = $reflectorFromObj->getName();
echo "从对象获取的类名: " . $classNameFromObj . "\n";
if ($classNameFromObj == "TestClass") {
    echo "✓ 从对象创建 ReflectionClass 测试通过\n";
} else {
    echo "✗ 从对象创建 ReflectionClass 测试失败\n";
}

// 测试 isInstantiable 是否能正确跟踪类信息
echo "\n=== 测试 isInstantiable 跟踪类信息 ===\n";
$reflector2 = new ReflectionClass("TestClass");
$isInstantiable2 = $reflector2->isInstantiable();
echo "isInstantiable() 结果: " . ($isInstantiable2 ? "true" : "false") . "\n";
if ($isInstantiable2) {
    echo "✓ isInstantiable() 类信息跟踪测试通过\n";
} else {
    echo "✗ isInstantiable() 类信息跟踪测试失败\n";
}

echo "\n=== 所有测试完成 ===\n";
