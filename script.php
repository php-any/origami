<?php
// 反射功能演示

echo "=== 反射功能演示 ===\n";

// 定义一个测试类
class User {
    public string $name = "John";
    public int $age = 30;
    private string $email = "john@example.com";
    public static string $version = "1.0";
    
    public function getName() {
        return $this->name;
    }
    
    public function setName($name) {
        $this->name = $name;
    }
    
    private function getEmail() {
        return $this->email;
    }
}

// 创建反射对象
$reflect = new Reflect();

// 获取类信息
echo "=== 类信息 ===\n";
$classInfo = $reflect->getClassInfo("User");
echo "类信息: " . $classInfo . "\n";

// 列出所有方法
echo "\n=== 方法列表 ===\n";
$methods = $reflect->listMethods("User");
echo "方法: " . $methods . "\n";

// 列出所有属性
echo "\n=== 属性列表 ===\n";
$properties = $reflect->listProperties("User");
echo "属性详细信息:\n";
foreach ($properties as $property) {
    $info = $reflect->getPropertyInfo("User", $property);
    echo "  " . $info->name . ":\n";
    echo "    类型: " . $info->type . "\n";
    echo "    权限: " . $info->modifier . "\n";
    echo "    静态: " . ($info->isStatic ? "是" : "否") . "\n";
    echo "    默认值: " . $info->defaultValue . "\n";
}

// 显示每个属性的详细信息
echo "\n=== 属性详细信息 ===\n";
echo "name 属性: " . $reflect->getPropertyInfo("User", "name") . "\n";
echo "age 属性: " . $reflect->getPropertyInfo("User", "age") . "\n";
echo "email 属性: " . $reflect->getPropertyInfo("User", "email") . "\n";
echo "version 属性: " . $reflect->getPropertyInfo("User", "version") . "\n";

// 获取特定方法信息
echo "\n=== 方法详细信息 ===\n";
$methodInfo = $reflect->getMethodInfo("User", "getName");
echo "getName 方法: " . $methodInfo . "\n";

// 获取特定属性信息
echo "\n=== 属性详细信息 ===\n";
$propertyInfo = $reflect->getPropertyInfo("User", "name");
echo "name 属性: " . $propertyInfo . "\n";

echo "\n=== 反射演示完成 ===\n";