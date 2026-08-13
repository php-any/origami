<?php

echo "=== extract() 闭包 & flags 测试 ===\n";

// -------------------------------------------------------------------
// 测试1：在立即执行闭包中使用 extract（模拟 Laravel Filesystem::getRequire）
// -------------------------------------------------------------------
$__data = ["title" => "Hello", "author" => "Laravel"];

$fn1 = static function () use ($__data) {
    extract($__data);
    return $title . " by " . $author;
};
$result = $fn1();

if ($result == "Hello by Laravel") {
    Log::info("测试1 通过：闭包内 extract 基本场景: " . $result);
} else {
    Log::fatal("测试1 失败：期望 Hello by Laravel, 实际: " . $result);
}

// -------------------------------------------------------------------
// 测试2：extract() 配合 EXTR_SKIP —— 不覆盖已有变量
// （模拟 Laravel Filesystem::getRequire 使用的 EXTR_SKIP 模式）
// -------------------------------------------------------------------
$fn2 = static function () {
    $title = "Original";
    $data = ["title" => "New", "extra" => "Appended"];
    extract($data, EXTR_SKIP);
    // title 已存在应跳过保持 Original，extra 注入到上下文（若已有槽位）
    return $title;
};
$result2 = $fn2();

if ($result2 == "Original") {
    Log::info("测试2 通过：EXTR_SKIP 跳过已有变量, title=' " . $result2 . "'");
} else {
    Log::fatal("测试2 失败：期望 Original, 实际: " . $result2);
}

// -------------------------------------------------------------------
// 测试3：extract() 配合 EXTR_OVERWRITE（默认）—— 覆盖已有变量
// -------------------------------------------------------------------
$fn3 = static function () {
    $name = "OldName";
    $data = ["name" => "NewName", "role" => "admin"];
    extract($data, EXTR_OVERWRITE);
    return $name;
};
$result3 = $fn3();

if ($result3 == "NewName") {
    Log::info("测试3 通过：EXTR_OVERWRITE 覆盖已有变量: " . $result3);
} else {
    Log::fatal("测试3 失败：期望 NewName, 实际: " . $result3);
}

// -------------------------------------------------------------------
// 测试4：extract() 配合 EXTR_PREFIX_ALL —— 所有变量加前缀
// -------------------------------------------------------------------
$fn4 = static function () {
    $data = ["foo" => "bar"];
    $pfx_foo = "before";
    extract($data, EXTR_PREFIX_ALL, "pfx");
    return $pfx_foo;
};
$result4 = $fn4();

if ($result4 == "bar") {
    Log::info("测试4 通过：EXTR_PREFIX_ALL 加前缀: " . $result4);
} else {
    Log::fatal("测试4 失败：期望 bar, 实际: " . $result4);
}

// -------------------------------------------------------------------
// 测试5：EXTR_SKIP 常量值验证（与整数 1 等价）
// -------------------------------------------------------------------
$__data2 = ["key1" => "val1", "key2" => "val2"];
$fn5 = static function () use ($__data2) {
    $key1 = "existing";
    extract($__data2, EXTR_SKIP);
    return $key1; // key1 已存在，EXTR_SKIP 不覆盖，应仍为 existing
};
$result5 = $fn5();

if ($result5 == "existing") {
    Log::info("测试5 通过：EXTR_SKIP use 闭包中不覆盖已有变量: " . $result5);
} else {
    Log::fatal("测试5 失败：期望 existing, 实际: " . $result5);
}

echo "=== extract() 闭包 & flags 测试完成 ===\n";
