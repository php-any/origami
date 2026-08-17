<?php

// Livewire Admin smoke: 验证 Livewire 管理后台的模型、路由和组件

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

// 1. Models exist and can be loaded
$models = [
    'App\\Models\\Admin',
    'App\\Models\\User',
    'App\\Models\\Role',
    'App\\Models\\Permission',
    'App\\Models\\Product',
    'App\\Models\\Order',
    'App\\Models\\OrderItem',
];
foreach ($models as $model) {
    if (!class_exists($model)) {
        echo "FAIL: $model 不存在\n";
        exit(1);
    }
    $checks++;
}

// 2. Livewire components exist
$components = [
    'App\\Livewire\\Admin\\Login',
    'App\\Livewire\\Admin\\Dashboard',
    'App\\Livewire\\Admin\\Admins\\Index',
    'App\\Livewire\\Admin\\Admins\\Form',
    'App\\Livewire\\Admin\\Users\\Index',
    'App\\Livewire\\Admin\\Users\\Form',
    'App\\Livewire\\Admin\\Roles\\Index',
    'App\\Livewire\\Admin\\Roles\\Form',
    'App\\Livewire\\Admin\\Permissions\\Index',
    'App\\Livewire\\Admin\\Permissions\\Form',
    'App\\Livewire\\Admin\\Products\\Index',
    'App\\Livewire\\Admin\\Products\\Form',
    'App\\Livewire\\Admin\\Orders\\Index',
    'App\\Livewire\\Admin\\Orders\\Detail',
    'App\\Livewire\\Admin\\Profile',
];
foreach ($components as $comp) {
    if (!class_exists($comp)) {
        echo "FAIL: $comp 不存在\n";
        exit(1);
    }
    $checks++;
}

// 3. Verify database has seeded data
$adminCount = \App\Models\Admin::count();
if ($adminCount < 1) {
    echo "FAIL: 管理员数量错误\n";
    exit(1);
}
$checks++;

$roleCount = \App\Models\Role::count();
if ($roleCount < 3) {
    echo "FAIL: 角色数量错误\n";
    exit(1);
}
$checks++;

$userCount = \App\Models\User::count();
if ($userCount < 1) {
    echo "FAIL: 用户数量错误\n";
    exit(1);
}
$checks++;

// 4. Verify auth guard exists
$guard = $app->make('auth')->guard('admin');
if (!$guard) {
    echo "FAIL: admin guard 不存在\n";
    exit(1);
}
$checks++;

// 5. Verify admin model can query roles
$admin = \App\Models\Admin::first();
$roles = $admin->roles()->get();
if ($roles->isEmpty()) {
    echo "FAIL: 管理员没有角色\n";
    exit(1);
}
$checks++;

echo "OK: livewire admin smoke passed ($checks checks)\n";
