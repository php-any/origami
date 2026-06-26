<?php

echo "=== 可变变量 \$\$var 测试 ===\n";

// 基本可变变量
$varName = "greeting";
$$varName = "Hello World";
if ($greeting == "Hello World") {
    Log::info("基本可变变量测试通过");
} else {
    Log::fatal("基本可变变量测试失败, 实际 '{$greeting}'");
}

// 可变变量读取
$color = "red";
$attr = "color";
$result = $$attr;
if ($result == "red") {
    Log::info("可变变量读取测试通过");
} else {
    Log::fatal("可变变量读取测试失败, 实际 '{$result}'");
}

// 多层可变变量
$a = "b";
$b = "c";
$c = "found";
$result2 = $$$$a;
if ($result2 == "found") {
    Log::info("多层可变变量测试通过");
} else {
    Log::fatal("多层可变变量测试失败");
}

// 可变变量用于动态属性名
$key = "name";
$dynamic = "Origami";
$$key = $dynamic;
if ($name == "Origami") {
    Log::info("可变变量动态属性测试通过");
} else {
    Log::fatal("可变变量动态属性测试失败, 实际 '{$name}'");
}

// 可变变量在循环中
$prefix = "item";
for ($i = 1; $i <= 3; $i++) {
    $varN = $prefix . $i;
    $$varN = $i * 10;
}
if ($item1 == 10 && $item2 == 20 && $item3 == 30) {
    Log::info("循环中可变变量测试通过");
} else {
    Log::fatal("循环中可变变量测试失败, item1={$item1}, item2={$item2}, item3={$item3}");
}

echo "=== 可变变量 \$\$var 测试完成 ===\n";
