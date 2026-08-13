<?php

// autoload_smoke.php
// 阶段 1：验证 vendor/autoload 能正常加载，核心类可解析。

require __DIR__.'/../../vendor/autoload.php';

$checks = 0;

// 1. Composer autoload 注册后核心类可解析
if (!class_exists('Composer\Autoload\ClassLoader')) {
    echo "FAIL: Composer ClassLoader 未注册\n";
    exit(1);
}
$checks++;

// 2. Illuminate 核心类可加载
$coreClasses = [
    'Illuminate\\Foundation\\Application',
    'Illuminate\\Container\\Container',
    'Illuminate\\Support\\ServiceProvider',
    'Illuminate\\Config\\Repository',
];
foreach ($coreClasses as $class) {
    if (!class_exists($class)) {
        echo "FAIL: 类 $class 无法加载\n";
        exit(1);
    }
    $checks++;
}

// 3. Laravel helper 函数存在
if (!function_exists('app')) {
    echo "FAIL: app() helper 不存在\n";
    exit(1);
}
$checks++;

echo "OK: autoload smoke passed ($checks checks)\n";
