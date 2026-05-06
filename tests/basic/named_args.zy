<?php

echo "=== 命名参数测试 ===\n";

// 基本命名参数
function createUser($name, $email, $role = "user") {
    return $name . "|" . $email . "|" . $role;
}

// 使用命名参数跳过可选参数
$result = createUser(name: "Alice", email: "alice@test.com", role: "admin");
if ($result == "Alice|alice@test.com|admin") {
    Log::info("命名参数基本测试通过");
} else {
    Log::fatal("命名参数基本测试失败, 实际 '{$result}'");
}

// 使用命名参数且不传可选参数
$result2 = createUser(email: "bob@test.com", name: "Bob");
if ($result2 == "Bob|bob@test.com|user") {
    Log::info("命名参数使用默认值测试通过");
} else {
    Log::fatal("命名参数使用默认值测试失败, 实际 '{$result2}'");
}

// 命名参数顺序无关
$result3 = createUser(role: "mod", email: "carol@test.com", name: "Carol");
if ($result3 == "Carol|carol@test.com|mod") {
    Log::info("命名参数顺序无关测试通过");
} else {
    Log::fatal("命名参数顺序无关测试失败, 实际 '{$result3}'");
}

// 命名参数与位置参数混合
function format($prefix, $value, $suffix = "") {
    return $prefix . $value . $suffix;
}
$result4 = format("[", "hello", suffix: "]");
if ($result4 == "[hello]") {
    Log::info("命名参数与位置参数混合测试通过");
} else {
    Log::fatal("命名参数与位置参数混合测试失败, 实际 '{$result4}'");
}

// 命名参数用于数学函数
function power($base, $exp = 2) {
    $result = 1;
    for ($i = 0; $i < $exp; $i++) {
        $result = $result * $base;
    }
    return $result;
}
$result5 = power(exp: 3, base: 2);
if ($result5 == 8) {
    Log::info("命名参数数学函数测试通过");
} else {
    Log::fatal("命名参数数学函数测试失败, 实际 {$result5}");
}

echo "=== 命名参数测试完成 ===\n";
