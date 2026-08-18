<?php

// Livewire bootstrap smoke: 验证 Livewire 服务提供者能正确注册和启动

require __DIR__.'/../../vendor/autoload.php';

/** @var Illuminate\Foundation\Application $app */
$app = require __DIR__.'/../../bootstrap/app.php';

$app->bootstrapWith([
    Illuminate\Foundation\Bootstrap\LoadEnvironmentVariables::class,
    Illuminate\Foundation\Bootstrap\LoadConfiguration::class,
    Illuminate\Foundation\Bootstrap\HandleExceptions::class,
    Illuminate\Foundation\Bootstrap\RegisterFacades::class,
    Illuminate\Foundation\Bootstrap\SetRequestForConsole::class,
    Illuminate\Foundation\Bootstrap\RegisterProviders::class,
    Illuminate\Foundation\Bootstrap\BootProviders::class,
]);

$checks = 0;

// 1. Livewire 服务提供者已注册
if (!$app->getProvider('Livewire\LivewireServiceProvider')) {
    echo "FAIL: Livewire 服务提供者未注册\n";
    exit(1);
}
$checks++;

// 2. Livewire 管理器可解析
try {
    $manager = $app->make('livewire');
    if (!$manager) {
        echo "FAIL: livewire 服务无法解析\n";
        exit(1);
    }
    $checks++;
} catch (Throwable $e) {
    echo "FAIL: livewire 服务解析失败: " . $e->getMessage() . "\n";
    exit(1);
}

// 3. Livewire 类可加载
if (!class_exists('Livewire\\Livewire')) {
    echo "FAIL: Livewire\\Livewire 类不存在\n";
    exit(1);
}
$checks++;

// 4. 核心机制类可加载
$mechanisms = [
    'Livewire\\Mechanisms\\Mechanism',
    'Livewire\\Mechanisms\\HandleComponents\\HandleComponents',
    'Livewire\\Mechanisms\\HandleRequests\\HandleRequests',
    'Livewire\\Mechanisms\\HandleRouting\\HandleRouting',
    'Livewire\\Mechanisms\\ExtendBlade\\ExtendBlade',
    'Livewire\\Mechanisms\\FrontendAssets\\FrontendAssets',
];
foreach ($mechanisms as $mech) {
    if (!class_exists($mech)) {
        echo "FAIL: $mech 类不存在\n";
        exit(1);
    }
    $checks++;
}

echo "OK: livewire bootstrap smoke passed ($checks checks)\n";
