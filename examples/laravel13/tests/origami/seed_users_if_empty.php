<?php
/**
 * 在 class_basename($this) 修好后补种 users（完整 db:seed 会撞已有 admins）。
 */
$_SERVER['HTTP_HOST'] = '127.0.0.1';
$_SERVER['SERVER_NAME'] = '127.0.0.1';
$_SERVER['REQUEST_URI'] = '/';
$_SERVER['REQUEST_METHOD'] = 'GET';
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\User;
use Illuminate\Support\Facades\Hash;

echo 'table='.var_export((new User)->getTable(), true)."\n";
echo 'users_before='.User::query()->count()."\n";

$seeds = [
    ['name' => '测试用户', 'email' => 'user@example.com'],
    ['name' => '用户1', 'email' => 'user1@example.com'],
    ['name' => '用户2', 'email' => 'user2@example.com'],
    ['name' => '用户3', 'email' => 'user3@example.com'],
    ['name' => '用户4', 'email' => 'user4@example.com'],
    ['name' => '用户5', 'email' => 'user5@example.com'],
];
foreach ($seeds as $row) {
    User::query()->firstOrCreate(
        ['email' => $row['email']],
        ['name' => $row['name'], 'password' => Hash::make('password')]
    );
}
echo 'users_after='.User::query()->count()."\n";
