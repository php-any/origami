<?php

echo "=== global 关键字和引用测试 ===\n";

// 全局变量测试
$globalVar = 100;

function readGlobal() {
    global $globalVar;
    return $globalVar;
}

$result = readGlobal();
if ($result == 100) {
    Log::info("global 读取全局变量测试通过");
} else {
    Log::fatal("global 读取全局变量测试失败, 实际 {$result}");
}

// 全局变量修改
function modifyGlobal() {
    global $globalVar;
    $globalVar = 200;
}

modifyGlobal();
if ($globalVar == 200) {
    Log::info("global 修改全局变量测试通过");
} else {
    Log::fatal("global 修改全局变量测试失败, 实际 {$globalVar}");
}
$globalVar = 100;

// 引用传递测试
$value = 10;
function addTen(&$ref) {
    $ref = $ref + 10;
}
addTen($value);
if ($value == 20) {
    Log::info("引用传递修改测试通过");
} else {
    Log::fatal("引用传递修改测试失败, 实际 {$value}");
}

// 引用赋值
$a = 5;
$b = &$a;
$b = 10;
if ($a == 10) {
    Log::info("引用赋值测试通过");
} else {
    Log::fatal("引用赋值测试失败, 实际 {$a}");
}

// 引用数组元素
$arr = [1, 2, 3];
$ref = &$arr[1];
$ref = 99;
if ($arr[1] == 99) {
    Log::info("引用数组元素测试通过");
} else {
    Log::fatal("引用数组元素测试失败, 实际 {$arr[1]}");
}

// 函数返回引用修改
function &getRef(&$val) {
    return $val;
}
$x = 42;
$r = &getRef($x);
$r = 100;
if ($x == 100) {
    Log::info("函数返回引用测试通过");
} else {
    Log::fatal("函数返回引用测试失败, 实际 {$x}");
}

// 多个 global 声明
$g1 = "one";
$g2 = "two";
function readMultipleGlobals() {
    global $g1, $g2;
    return $g1 . "-" . $g2;
}
$result2 = readMultipleGlobals();
if ($result2 == "one-two") {
    Log::info("多个 global 声明测试通过");
} else {
    Log::fatal("多个 global 声明测试失败, 实际 '{$result2}'");
}

echo "=== global 关键字和引用测试完成 ===\n";
