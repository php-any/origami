<?php

// 测试 get_debug_type 函数

echo "=== get_debug_type() 函数测试 ===\n\n";

// 1. 基本类型
echo "1. 基本类型:\n";
echo "int: " . get_debug_type(42) . "\n";
echo "float: " . get_debug_type(3.14) . "\n";
echo "string: " . get_debug_type("hello") . "\n";
echo "bool: " . get_debug_type(true) . "\n";
echo "null: " . get_debug_type(null) . "\n";
echo "array: " . get_debug_type([1, 2, 3]) . "\n";

// 2. 对象类型（返回类名）
echo "\n2. 对象类型:\n";
class TestClass {
    public $prop = "value";
}
$obj = new TestClass();
echo "object: " . get_debug_type($obj) . "\n";

// 3. 资源类型
echo "\n3. 资源类型:\n";
$file = fopen("/tmp/test_get_debug_type.txt", "w");
if ($file !== false) {
    echo "resource (open): " . get_debug_type($file) . "\n";
    fclose($file);
    echo "resource (closed): " . get_debug_type($file) . "\n";
    unlink("/tmp/test_get_debug_type.txt");
}

// 4. 与 gettype() 对比
echo "\n4. 与 gettype() 对比:\n";
$int = 42;
echo "gettype(42): " . gettype($int) . "\n";
echo "get_debug_type(42): " . get_debug_type($int) . "\n";

$obj2 = new TestClass();
echo "gettype(object): " . gettype($obj2) . "\n";
echo "get_debug_type(object): " . get_debug_type($obj2) . "\n";

echo "\n=== get_debug_type() 测试完成 ===\n";
