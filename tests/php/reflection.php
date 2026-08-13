<?php
namespace tests\php;

use ReflectionNamedType;
use ReflectionType;

echo "=== ReflectionClass 功能测试 ===\n";

// 定义一个测试类
class TestReflectionClass {
    public $publicProp = "public";
    protected $protectedProp = "protected";
    private $privateProp = "private";

    public function __construct($value = "default") {
        $this->publicProp = $value;
    }

    public function publicMethod() {
        return "public method";
    }

    public function methodWithParams($param1, $param2 = "default") {
        return $param1 . $param2;
    }

    public function methodWithTypedParams(string $strParam, int $intParam = 0) {
        return $strParam . $intParam;
    }

    public function methodWithClassParam(TestReflectionClass $objParam) {
        return $objParam->publicProp;
    }

    protected function protectedMethod() {
        return "protected method";
    }

    private function privateMethod() {
        return "private method";
    }
}

// 定义继承类
class ChildReflectionClass extends TestReflectionClass {
    public $childProp = "child";

    public function childMethod() {
        return "child method";
    }
}

// 先实例化类以确保类被加载到 VM 中
$tempObj = new TestReflectionClass("temp");
$tempChild = new ChildReflectionClass("temp");

// 测试 ReflectionClass::__construct 和 getName
echo "=== 测试 ReflectionClass::getName ===\n";
$reflector = new \ReflectionClass("tests\\php\\TestReflectionClass");
$className = $reflector->getName();
if ($className == "tests\\php\\TestReflectionClass") {
    Log::info("getName() 测试通过");
} else {
    Log::fatal("getName() 测试失败，期望: tests\\php\\TestReflectionClass, 实际: {$className}");
}

// 测试 ReflectionClass::hasMethod
echo "\n=== 测试 ReflectionClass::hasMethod ===\n";
$hasPublic = $reflector->hasMethod("publicMethod");
$hasProtected = $reflector->hasMethod("protectedMethod");
$hasPrivate = $reflector->hasMethod("privateMethod");
$hasNotExist = $reflector->hasMethod("notExistMethod");

if ($hasPublic && $hasProtected && $hasPrivate && !$hasNotExist) {
    Log::info("hasMethod() 测试通过");
} else {
    Log::fatal("hasMethod() 测试失败");
}

// 测试 ReflectionClass::hasProperty
echo "\n=== 测试 ReflectionClass::hasProperty ===\n";
$hasPublicProp = $reflector->hasProperty("publicProp");
$hasProtectedProp = $reflector->hasProperty("protectedProp");
$hasPrivateProp = $reflector->hasProperty("privateProp");
$hasNotExistProp = $reflector->hasProperty("notExistProp");

if ($hasPublicProp && $hasProtectedProp && $hasPrivateProp && !$hasNotExistProp) {
    Log::info("hasProperty() 测试通过");
} else {
    Log::fatal("hasProperty() 测试失败");
}

// 测试 ReflectionClass::getMethods
echo "\n=== 测试 ReflectionClass::getMethods ===\n";
$methods = $reflector->getMethods();
$methodCount = count($methods);
if ($methodCount >= 4) { // 至少应该有 __construct, publicMethod, protectedMethod, privateMethod
    Log::info("getMethods() 测试通过，方法数量: {$methodCount}");
} else {
    Log::fatal("getMethods() 测试失败，方法数量: {$methodCount}");
}

// 测试 ReflectionClass::getProperties
echo "\n=== 测试 ReflectionClass::getProperties ===\n";
$properties = $reflector->getProperties();
$propertyCount = count($properties);
if ($propertyCount >= 3) { // 至少应该有 publicProp, protectedProp, privateProp
    Log::info("getProperties() 测试通过，属性数量: {$propertyCount}");
} else {
    Log::fatal("getProperties() 测试失败，属性数量: {$propertyCount}");
}

// 测试 ReflectionClass::isSubclassOf
echo "\n=== 测试 ReflectionClass::isSubclassOf ===\n";
$childReflector = new \ReflectionClass("tests\\php\\ChildReflectionClass");
$isSubclass = $childReflector->isSubclassOf("tests\\php\\TestReflectionClass");
if ($isSubclass) {
    Log::info("isSubclassOf() 测试通过");
} else {
    Log::fatal("isSubclassOf() 测试失败");
}

// 测试 ReflectionClass::getParentClass
echo "\n=== 测试 ReflectionClass::getParentClass ===\n";
$parentClass = $childReflector->getParentClass();
// getParentClass 返回的是 ReflectionClass 实例，需要取名字再比较
if ($parentClass && $parentClass->getName() == "tests\\php\\TestReflectionClass") {
    Log::info("getParentClass() 测试通过");
} else {
    Log::fatal("getParentClass() 测试失败，期望: tests\\php\\TestReflectionClass, 实际: {$parentClass}");
}

// 测试 ReflectionClass::isInstance
echo "\n=== 测试 ReflectionClass::isInstance ===\n";
$testObj = new TestReflectionClass("test");
$isInstance = $reflector->isInstance($testObj);
if ($isInstance) {
    Log::info("isInstance() 测试通过");
} else {
    Log::fatal("isInstance() 测试失败");
}

// 测试 ReflectionClass::isInstantiable
echo "\n=== 测试 ReflectionClass::isInstantiable ===\n";
$isInstantiable = $reflector->isInstantiable();
if ($isInstantiable) {
    Log::info("isInstantiable() 测试通过");
} else {
    Log::fatal("isInstantiable() 测试失败");
}

// 测试 ReflectionClass::newInstance
echo "\n=== 测试 ReflectionClass::newInstance ===\n";
$newObj = $reflector->newInstance("new value");
if ($newObj->publicProp == "new value") {
    Log::info("newInstance() 测试通过");
} else {
    Log::fatal("newInstance() 测试失败，期望: new value, 实际: {$newObj->publicProp}");
}

// 测试 ReflectionClass::newInstanceWithoutConstructor
echo "\n=== 测试 ReflectionClass::newInstanceWithoutConstructor ===\n";
$newObjWithoutCtor = $reflector->newInstanceWithoutConstructor();
if ($newObjWithoutCtor->publicProp == "public") {
    Log::info("newInstanceWithoutConstructor() 测试通过");
} else {
    Log::fatal("newInstanceWithoutConstructor() 测试失败，期望: public, 实际: {$newObjWithoutCtor->publicProp}");
}

// 测试 ReflectionClass::newInstanceArgs
echo "\n=== 测试 ReflectionClass::newInstanceArgs ===\n";
$newObjArgs = $reflector->newInstanceArgs(["args value"]);
if ($newObjArgs->publicProp == "args value") {
    Log::info("newInstanceArgs() 测试通过");
} else {
    Log::fatal("newInstanceArgs() 测试失败，期望: args value, 实际: {$newObjArgs->publicProp}");
}

// 测试使用对象创建 ReflectionClass
echo "\n=== 测试使用对象创建 ReflectionClass ===\n";
$obj = new TestReflectionClass("object test");
$reflectorFromObj = new \ReflectionClass($obj);
$classNameFromObj = $reflectorFromObj->getName();
if ($classNameFromObj == "tests\\php\\TestReflectionClass") {
    Log::info("从对象创建 ReflectionClass 测试通过");
} else {
    Log::fatal("从对象创建 ReflectionClass 测试失败，期望: tests\\php\\TestReflectionClass, 实际: {$classNameFromObj}");
}

// 测试 ReflectionClass::getConstructor
echo "\n=== 测试 ReflectionClass::getConstructor ===\n";
$constructor = $reflector->getConstructor();
if ($constructor != null) {
    $constructorName = $constructor->getName();
    if ($constructorName == "__construct") {
        Log::info("getConstructor() 测试通过");
    } else {
        Log::fatal("getConstructor() 测试失败，期望方法名: __construct, 实际: {$constructorName}");
    }
} else {
    Log::fatal("getConstructor() 测试失败，构造函数为 null");
}

// 测试 ReflectionClass::getMethod
echo "\n=== 测试 ReflectionClass::getMethod ===\n";
$method = $reflector->getMethod("publicMethod");
if ($method != null) {
    $methodName = $method->getName();
    if ($methodName == "publicMethod") {
        Log::info("getMethod() 测试通过");
    } else {
        Log::fatal("getMethod() 测试失败，期望方法名: publicMethod, 实际: {$methodName}");
    }
} else {
    Log::fatal("getMethod() 测试失败，方法为 null");
}

// 测试 ReflectionMethod::getName
echo "\n=== 测试 ReflectionMethod::getName ===\n";
if ($method->getName() == "publicMethod") {
    Log::info("ReflectionMethod::getName() 测试通过");
} else {
    Log::fatal("ReflectionMethod::getName() 测试失败");
}

// 测试 ReflectionMethod::isPublic
echo "\n=== 测试 ReflectionMethod::isPublic ===\n";
if ($method->isPublic()) {
    Log::info("ReflectionMethod::isPublic() 测试通过");
} else {
    Log::fatal("ReflectionMethod::isPublic() 测试失败");
}

// 测试 ReflectionMethod::isStatic
echo "\n=== 测试 ReflectionMethod::isStatic ===\n";
if (!$method->isStatic()) {
    Log::info("ReflectionMethod::isStatic() 测试通过");
} else {
    Log::fatal("ReflectionMethod::isStatic() 测试失败");
}

// 测试 ReflectionMethod::getNumberOfParameters
echo "\n=== 测试 ReflectionMethod::getNumberOfParameters ===\n";
$paramCount = $method->getNumberOfParameters();
if ($paramCount == 0) {
    Log::info("ReflectionMethod::getNumberOfParameters() 测试通过");
} else {
    Log::fatal("ReflectionMethod::getNumberOfParameters() 测试失败，期望: 0, 实际: {$paramCount}");
}

// 测试 ReflectionMethod::getParameters
echo "\n=== 测试 ReflectionMethod::getParameters ===\n";
$parameters = $method->getParameters();
if (count($parameters) == 0) {
    Log::info("ReflectionMethod::getParameters() 测试通过（无参数方法）");
} else {
    Log::fatal("ReflectionMethod::getParameters() 测试失败，期望参数数量: 0, 实际: " . count($parameters));
}

// 测试有参数的方法
echo "\n=== 测试有参数的方法 getParameters ===\n";
$methodWithParams = $reflector->getMethod("methodWithParams");
if ($methodWithParams != null) {
    $params = $methodWithParams->getParameters();
    $paramCount = count($params);
    if ($paramCount == 2) {
        Log::info("ReflectionMethod::getParameters() 有参数方法测试通过，参数数量: {$paramCount}");
        // 检查返回的是 ReflectionParameter 对象
        if (count($params) >= 2) {
            $param1 = $params[0];
            $param2 = $params[1];

            // 验证返回的是对象
            if (is_object($param1) && is_object($param2)) {
                Log::info("ReflectionMethod::getParameters() 返回对象测试通过");

                // 检查参数名
                $param1Name = $param1->getName();
                $param2Name = $param2->getName();
                if ($param1Name == "param1" && $param2Name == "param2") {
                    Log::info("ReflectionMethod::getParameters() 参数名测试通过");
                } else {
                    Log::fatal("ReflectionMethod::getParameters() 参数名测试失败，期望: param1, param2, 实际: {$param1Name}, {$param2Name}");
                }

                // 检查参数位置
                $param1Pos = $param1->getPosition();
                $param2Pos = $param2->getPosition();
                if ($param1Pos == 0 && $param2Pos == 1) {
                    Log::info("ReflectionParameter::getPosition() 测试通过");
                } else {
                    Log::fatal("ReflectionParameter::getPosition() 测试失败，期望: 0, 1, 实际: {$param1Pos}, {$param2Pos}");
                }

                // 检查参数是否可选
                $param1Optional = $param1->isOptional();
                $param2Optional = $param2->isOptional();
                if (!$param1Optional && $param2Optional) {
                    Log::info("ReflectionParameter::isOptional() 测试通过");
                } else {
                    Log::fatal("ReflectionParameter::isOptional() 测试失败，期望: false, true, 实际: " . ($param1Optional ? "true" : "false") . ", " . ($param2Optional ? "true" : "false"));
                }

                // 测试 ReflectionParameter::getType
                // methodWithParams 的参数没有类型声明，所以应该返回 null
                $param1Type = $param1->getType();
                $param2Type = $param2->getType();
                if ($param1Type === null && $param2Type === null) {
                    Log::info("ReflectionParameter::getType() 测试通过（无类型声明返回 null）");
                } else {
                    Log::fatal("ReflectionParameter::getType() 测试失败，无类型声明应该返回 null");
                }

                // 测试用户代码逻辑：当 $type 为 null 时的处理
                // 验证用户代码：if (! $type instanceof ReflectionNamedType || $type->isBuiltin())
                // 当 $type 为 null 时，! $type instanceof ReflectionNamedType 为 true，所以会跳过
                if (!($param1Type instanceof ReflectionNamedType) || ($param1Type !== null && $param1Type->isBuiltin())) {
                    Log::info("用户代码逻辑测试通过：null 类型会被跳过");
                } else {
                    Log::fatal("用户代码逻辑测试失败：null 类型应该被跳过");
                }
            } else {
                Log::fatal("ReflectionMethod::getParameters() 返回对象测试失败，期望对象，实际类型: " . gettype($param1) . ", " . gettype($param2));
            }
        }
    } else {
        Log::fatal("ReflectionMethod::getParameters() 有参数方法测试失败，期望参数数量: 2, 实际: {$paramCount}");
    }
} else {
    Log::fatal("ReflectionMethod::getMethod() 获取 methodWithParams 失败");
}

// 测试带类型声明的方法
echo "\n=== 测试带类型声明的方法 getType 和 isBuiltin ===\n";
$methodWithTypedParams = $reflector->getMethod("methodWithTypedParams");
if ($methodWithTypedParams != null) {
    $typedParams = $methodWithTypedParams->getParameters();
    if (count($typedParams) >= 2) {
        $strParam = $typedParams[0];
        $intParam = $typedParams[1];

        $strType = $strParam->getType();
        $intType = $intParam->getType();

        if ($strType !== null && $intType !== null) {
            // 验证返回的是 ReflectionNamedType 对象
            if (is_object($strType) && is_object($intType)) {
                Log::info("ReflectionParameter::getType() 返回对象测试通过");

                // 测试 instanceof ReflectionNamedType
                if ($strType instanceof ReflectionNamedType && $intType instanceof ReflectionNamedType) {
                    Log::info("instanceof ReflectionNamedType 测试通过");
                } else {
                    Log::fatal("instanceof ReflectionNamedType 测试失败");
                }

                // 测试 instanceof ReflectionType（应该也返回 true，因为 ReflectionNamedType 继承 ReflectionType）
                if ($strType instanceof ReflectionType && $intType instanceof ReflectionType) {
                    Log::info("instanceof ReflectionType 测试通过（继承关系正确）");
                } else {
                    Log::fatal("instanceof ReflectionType 测试失败");
                }

                // 测试 ReflectionType::getName()
                $strTypeName = $strType->getName();
                $intTypeName = $intType->getName();
                if ($strTypeName == "string" && $intTypeName == "int") {
                    Log::info("ReflectionType::getName() 测试通过");
                } else {
                    Log::fatal("ReflectionType::getName() 测试失败，期望: string, int, 实际: {$strTypeName}, {$intTypeName}");
                }

                // 测试 ReflectionType::isBuiltin()
                // 验证内置类型的 isBuiltin() 返回 true
                $strIsBuiltin = $strType->isBuiltin();
                $intIsBuiltin = $intType->isBuiltin();
                if ($strIsBuiltin && $intIsBuiltin) {
                    Log::info("ReflectionType::isBuiltin() 测试通过（string 和 int 都是内置类型，返回 true）");
                } else {
                    Log::fatal("ReflectionType::isBuiltin() 测试失败，期望: true, true, 实际: " . ($strIsBuiltin ? "true" : "false") . ", " . ($intIsBuiltin ? "true" : "false"));
                }

                // 测试用户代码中的逻辑：! $type instanceof ReflectionNamedType || $type->isBuiltin()
                // 对于内置类型，这个条件应该为 true（因为 isBuiltin() 返回 true）
                if (!($strType instanceof ReflectionNamedType) || $strType->isBuiltin()) {
                    Log::info("内置类型检查逻辑测试通过：! instanceof ReflectionNamedType || isBuiltin() 返回 true");
                } else {
                    Log::fatal("内置类型检查逻辑测试失败");
                }

                if (!($intType instanceof ReflectionNamedType) || $intType->isBuiltin()) {
                    Log::info("内置类型检查逻辑测试通过：int 类型");
                } else {
                    Log::fatal("内置类型检查逻辑测试失败：int 类型");
                }
            } else {
                Log::fatal("ReflectionParameter::getType() 返回对象测试失败，期望 ReflectionNamedType 对象");
            }
        } else {
            Log::fatal("ReflectionParameter::getType() 测试失败，期望非 null");
        }
    } else {
        Log::fatal("ReflectionMethod::getParameters() 测试失败，期望至少 2 个参数");
    }
} else {
    Log::fatal("ReflectionMethod::getMethod() 获取 methodWithTypedParams 失败");
}

// 测试非内置类型（自定义类）的 isBuiltin()
echo "\n=== 测试非内置类型的 isBuiltin() ===\n";
$methodWithClassParam = $reflector->getMethod("methodWithClassParam");
if ($methodWithClassParam != null) {
    $classParams = $methodWithClassParam->getParameters();
    if (count($classParams) >= 1) {
        $objParam = $classParams[0];
        $objType = $objParam->getType();

        if ($objType !== null) {
            if (is_object($objType)) {
                // 验证是 ReflectionNamedType
                if ($objType instanceof ReflectionNamedType) {
                    Log::info("自定义类类型返回 ReflectionNamedType 测试通过");

                    // 测试 isBuiltin()，应该返回 false（因为 TestReflectionClass 不是内置类型）
                    $objIsBuiltin = $objType->isBuiltin();
                    if (!$objIsBuiltin) {
                        Log::info("ReflectionType::isBuiltin() 测试通过（自定义类返回 false）");
                    } else {
                        Log::fatal("ReflectionType::isBuiltin() 测试失败，自定义类应该返回 false，实际: " . ($objIsBuiltin ? "true" : "false"));
                    }

                    // 测试用户代码中的逻辑：! $type instanceof ReflectionNamedType || $type->isBuiltin()
                    // 对于非内置类型，这个条件应该为 false（因为 instanceof 为 true 且 isBuiltin() 为 false）
                    if (!($objType instanceof ReflectionNamedType) || $objType->isBuiltin()) {
                        Log::fatal("非内置类型检查逻辑测试失败：应该返回 false");
                    } else {
                        Log::info("非内置类型检查逻辑测试通过：! instanceof ReflectionNamedType || isBuiltin() 返回 false");
                    }

                    // 验证类型名称
                    $objTypeName = $objType->getName();
                    if ($objTypeName == "tests\\php\\TestReflectionClass" || $objTypeName == "TestReflectionClass") {
                        Log::info("ReflectionType::getName() 测试通过，自定义类类型名: {$objTypeName}");
                    } else {
                        Log::fatal("ReflectionType::getName() 测试失败，期望: TestReflectionClass, 实际: {$objTypeName}");
                    }
                } else {
                    Log::fatal("自定义类类型应该返回 ReflectionNamedType");
                }
            } else {
                Log::fatal("ReflectionParameter::getType() 返回对象测试失败");
            }
        } else {
            Log::fatal("ReflectionParameter::getType() 测试失败，期望非 null");
        }
    } else {
        Log::fatal("ReflectionMethod::getParameters() 测试失败，期望至少 1 个参数");
    }
} else {
    Log::fatal("ReflectionMethod::getMethod() 获取 methodWithClassParam 失败");
}

// 测试 isInstantiable 是否能正确跟踪类信息
echo "\n=== 测试 isInstantiable 跟踪类信息 ===\n";
$reflector2 = new \ReflectionClass("tests\\php\\TestReflectionClass");
$isInstantiable2 = $reflector2->isInstantiable();
if ($isInstantiable2) {
    Log::info("isInstantiable() 类信息跟踪测试通过");
} else {
    Log::fatal("isInstantiable() 类信息跟踪测试失败");
}

echo "\n=== ReflectionClass 所有测试完成 ===\n";