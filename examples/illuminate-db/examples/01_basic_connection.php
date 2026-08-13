<?php
/**
 * 示例 1：基础连接与查询
 *
 * 演示如何使用 Capsule 建立独立于 Laravel 框架的数据库连接，
 * 并通过查询构建器执行基础查询。
 *
 * 运行方式: ./illuminate-db examples/01_basic_connection.php
 */

require __DIR__ . '/../vendor/autoload.php';
require __DIR__ . '/../bootstrap.php';

use Illuminate\Database\Capsule\Manager as Capsule;

// 初始化 Capsule（内存 SQLite；也可通过 DB_DATABASE 环境变量指定文件路径）
$capsule = bootstrap_capsule();

// 创建 users 表
Capsule::schema()->create('users', function ($table) {
    $table->increments('id');
    $table->string('name');
    $table->string('email')->unique();
    $table->timestamp('created_at')->nullable();
});

// 插入数据
Capsule::table('users')->insert([
    'name'       => 'John Doe',
    'email'      => 'john@example.com',
    'created_at' => '2025-01-01 10:00:00',
]);

Capsule::table('users')->insert([
    'name'       => 'Jane Smith',
    'email'      => 'jane@example.com',
    'created_at' => '2025-02-15 14:30:00',
]);

// 查询所有用户
$users = Capsule::table('users')->get();
echo "所有用户 (" . count($users) . " 条):\n";
foreach ($users as $user) {
    echo "  [{$user->id}] {$user->name} <{$user->email}>\n";
}

// 条件查询：取第一条匹配记录
$john = Capsule::table('users')
    ->where('name', 'John Doe')
    ->first();
echo "\n条件查询 first(): " . ($john ? $john->name : '未找到') . "\n";

// 条件查询：where  + orderBy
$recent = Capsule::table('users')
    ->orderBy('created_at', 'desc')
    ->get();
echo "按时间倒序:\n";
foreach ($recent as $user) {
    echo "  {$user->name} ({$user->created_at})\n";
}

// 聚合统计
$count = Capsule::table('users')->count();
echo "\n用户总数: {$count}\n";

echo "\n✓ 基础连接示例执行成功\n";
