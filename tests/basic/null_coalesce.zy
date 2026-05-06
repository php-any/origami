<?php

echo "=== null 合并运算符 ?? 和 ??= 测试 ===\n";

// ?? 基本用法: 左侧为 null
$unsetVar = null;
$result = $unsetVar ?? "default";
if ($result == "default") {
    Log::info("?? null 左侧返回默认值测试通过");
} else {
    Log::fatal("?? null 左侧测试失败, 实际 '{$result}'");
}

// ?? 左侧有值
$name = "origami";
$result2 = $name ?? "fallback";
if ($result2 == "origami") {
    Log::info("?? 有值左侧返回原值测试通过");
} else {
    Log::fatal("?? 有值左侧测试失败, 实际 '{$result2}'");
}

// ?? 链式使用
$a = null;
$b = null;
$c = "found";
$result3 = $a ?? $b ?? $c;
if ($result3 == "found") {
    Log::info("?? 链式使用测试通过");
} else {
    Log::fatal("?? 链式使用测试失败, 实际 '{$result3}'");
}

// ?? 链式全部为 null
$result4 = $a ?? $b ?? null;
if ($result4 === null) {
    Log::info("?? 链式全部为 null 测试通过");
} else {
    Log::fatal("?? 链式全部为 null 测试失败");
}

// ?? 在数组访问中的使用
$config = ["host" => "localhost"];
$port = $config["port"] ?? 3306;
if ($port == 3306) {
    Log::info("?? 数组访问默认值测试通过");
} else {
    Log::fatal("?? 数组访问默认值测试失败, 实际 {$port}");
}

// ?? 在数组访问中值存在
$host = $config["host"] ?? "127.0.0.1";
if ($host == "localhost") {
    Log::info("?? 数组访问存在值测试通过");
} else {
    Log::fatal("?? 数组访问存在值测试失败, 实际 '{$host}'");
}

// ??= null 合并赋值
$x = null;
$x ??= "assigned";
if ($x == "assigned") {
    Log::info("??= null 时赋值测试通过");
} else {
    Log::fatal("??= null 时赋值测试失败, 实际 '{$x}'");
}

// ??= 已有值时不赋值
$y = "original";
$y ??= "overwritten";
if ($y == "original") {
    Log::info("??= 已有值不覆盖测试通过");
} else {
    Log::fatal("??= 已有值不覆盖测试失败, 实际 '{$y}'");
}

// ?? 与 0 和空字符串 (这些不是 null, 应该返回原值)
$zero = 0;
$result5 = $zero ?? "default";
if ($result5 === 0) {
    Log::info("?? 0 不是 null 测试通过");
} else {
    Log::fatal("?? 0 不是 null 测试失败, 实际 {$result5}");
}

$empty = "";
$result6 = $empty ?? "default";
if ($result6 === "") {
    Log::info("?? 空字符串不是 null 测试通过");
} else {
    Log::fatal("?? 空字符串不是 null 测试失败, 实际 '{$result6}'");
}

echo "=== null 合并运算符 ?? 和 ??= 测试完成 ===\n";
