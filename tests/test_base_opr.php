<?php

echo "=== 运算符测试 ===\n";

// 基本算术运算符测试
echo "--- 基本算术运算符 ---\n";
$a = 10;
$b = 3;

echo "加法: {$a} + {$b} = " . ($a + $b) . "\n";
echo "减法: {$a} - {$b} = " . ($a - $b) . "\n";
echo "乘法: {$a} * {$b} = " . ($a * $b) . "\n";
echo "除法: {$a} / {$b} = " . ($a / $b) . "\n";
echo "取余: {$a} % {$b} = " . ($a % $b) . "\n";

// 比较运算符测试
echo "\n--- 比较运算符 ---\n";
echo "大于: {$a} > {$b} = " . ($a > $b ? 'true' : 'false') . "\n";
echo "小于: {$a} < {$b} = " . ($a < $b ? 'true' : 'false') . "\n";
echo "大于等于: {$a} >= {$b} = " . ($a >= $b ? 'true' : 'false') . "\n";
echo "小于等于: {$a} <= {$b} = " . ($a <= $b ? 'true' : 'false') . "\n";
echo "等于: {$a} == {$b} = " . ($a == $b ? 'true' : 'false') . "\n";
echo "不等于: {$a} != {$b} = " . ($a != $b ? 'true' : 'false') . "\n";

// 逻辑运算符测试
echo "\n--- 逻辑运算符 ---\n";
$c = true;
$d = false;
echo "逻辑与: {$c} && {$d} = " . ($c && $d ? 'true' : 'false') . "\n";
echo "逻辑或: {$c} || {$d} = " . ($c || $d ? 'true' : 'false') . "\n";
echo "逻辑非: !{$c} = " . (!$c ? 'true' : 'false') . "\n";

// 字符串连接运算符测试
echo "\n--- 字符串连接运算符 ---\n";
$str1 = "Hello";
$str2 = "World";
echo "字符串连接: {$str1} . {$str2} = " . ($str1 . $str2) . "\n";
echo "字符串连接数字: {$str1} . {$a} = " . ($str1 . $a) . "\n";

// 运算符优先级测试
echo "\n--- 运算符优先级测试 ---\n";

echo "1 + 2 * 3 = " . (1 + 2 * 3) . "\n";
echo "(1 + 2) * 3 = " . ((1 + 2) * 3) . "\n";
echo "10 - 3 * 2 = " . (10 - 3 * 2) . "\n";
echo "(10 - 3) * 2 = " . ((10 - 3) * 2) . "\n";

echo "1 + 2 > 3 && 4 < 5 = " . (1 + 2 > 3 && 4 < 5 ? 'true' : 'false') . "\n";
echo "(1 + 2 > 3) && (4 < 5) = " . ((1 + 2 > 3) && (4 < 5) ? 'true' : 'false') . "\n";

echo "'a' . 1 + 2 = " . ('a' . 1 + 2) . "\n";
echo "'a' . (1 + 2) = " . ('a' . (1 + 2)) . "\n";

// 复杂表达式测试
echo "\n--- 复杂表达式测试 ---\n";
$x = 5;
$y = 3;
$z = 2;

$result1 = $x * $y + $z;
echo "x * y + z = {$x} * {$y} + {$z} = {$result1}\n";

$result2 = $x * ($y + $z);
echo "x * (y + z) = {$x} * ({$y} + {$z}) = {$result2}\n";

$result3 = $x > $y && $y > $z;
echo "x > y && y > z = {$x} > {$y} && {$y} > {$z} = " . ($result3 ? 'true' : 'false') . "\n";

$result4 = ($x > $y) && ($y > $z);
echo "(x > y) && (y > z) = ({$x} > {$y}) && ({$y} > {$z}) = " . ($result4 ? 'true' : 'false') . "\n";

// 一元运算符测试
echo "\n--- 一元运算符测试 ---\n";
$num = 5;
echo "原始值: {$num}\n";
echo "负值: -{$num} = " . (-$num) . "\n";
echo "逻辑非: !true = " . (!true ? 'true' : 'false') . "\n";
echo "逻辑非: !false = " . (!false ? 'true' : 'false') . "\n";

// 自增自减运算符测试
echo "\n--- 自增自减运算符测试 ---\n";
$counter = 1;
echo "原始值: {$counter}\n";
$counter++;
echo "后自增: {$counter}\n";
++$counter;
echo "前自增: {$counter}\n";
$counter--;
echo "后自减: {$counter}\n";
--$counter;
echo "前自减: {$counter}\n";

// 赋值运算符测试
echo "\n--- 赋值运算符测试 ---\n";
$value = 10;
echo "原始值: {$value}\n";
$value += 5;
echo "+= 5: {$value}\n";
$value -= 3;
echo "-= 3: {$value}\n";
$value *= 2;
echo "*= 2: {$value}\n";
$value /= 4;
echo "/= 4: {$value}\n";
$value %= 3;
echo "%= 3: {$value}\n";

// 浮点数运算测试
echo "\n--- 浮点数运算测试 ---\n";
$float1 = 3.14;
$float2 = 2.86;
echo "浮点数加法: {$float1} + {$float2} = " . ($float1 + $float2) . "\n";
echo "浮点数乘法: {$float1} * {$float2} = " . ($float1 * $float2) . "\n";
echo "浮点数除法: {$float1} / {$float2} = " . ($float1 / $float2) . "\n";

// 混合类型运算测试
echo "\n--- 混合类型运算测试 ---\n";
$int_val = 10;
$float_val = 3.5;
$str_val = "5";

echo "整数 + 浮点数: {$int_val} + {$float_val} = " . ($int_val + $float_val) . "\n";
echo "字符串 + 整数: {$str_val} + {$int_val} = " . ($str_val + $int_val) . "\n";
echo "字符串 . 整数: {$str_val} . {$int_val} = " . ($str_val . $int_val) . "\n";

echo "\n=== 运算符测试完成 ===\n"; 