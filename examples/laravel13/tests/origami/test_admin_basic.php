<?php
require __DIR__.'/../../vendor/autoload.php';

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

use Illuminate\Support\Facades\DB;

$checks = 0;

// 1. Models load
foreach (['App\Models\Admin','App\Models\Role','App\Models\Permission','App\Models\Product','App\Models\Order','App\Models\OrderItem'] as $model) {
    if (!class_exists($model)) { echo "FAIL: $model\n"; exit(1); }
    $checks++;
}

// 2. Livewire components load
foreach (['App\Livewire\Admin\Login','App\Livewire\Admin\Dashboard','App\Livewire\Admin\Admins\Index','App\Livewire\Admin\Admins\Form','App\Livewire\Admin\Users\Index','App\Livewire\Admin\Roles\Index','App\Livewire\Admin\Permissions\Index','App\Livewire\Admin\Products\Index','App\Livewire\Admin\Orders\Index','App\Livewire\Admin\Orders\Detail','App\Livewire\Admin\Profile'] as $comp) {
    if (!class_exists($comp)) { echo "FAIL: $comp\n"; exit(1); }
    $checks++;
}

// 3. Verify database data via DB facade (bypasses Eloquent issues)
$adminCount = DB::table('admins')->count();
$roleCount = DB::table('roles')->count();
$permCount = DB::table('permissions')->count();
$userCount = DB::table('users')->count();
$productCount = DB::table('products')->count();
$orderCount = DB::table('orders')->count();

echo "DB data: admins=$adminCount roles=$roleCount perms=$permCount users=$userCount products=$productCount orders=$orderCount\n";
if ($adminCount < 1 || $roleCount < 3 || $permCount < 10 || $userCount < 1 || $productCount < 1 || $orderCount < 1) {
    echo "FAIL: 数据库数据不完整\n";
    exit(1);
}
$checks++;

// 4. Livewire is registered
if (!$app->getProvider('Livewire\LivewireServiceProvider')) {
    echo "FAIL: Livewire provider not registered\n";
    exit(1);
}
$checks++;

echo "OK: admin basic smoke passed ($checks checks)\n";
