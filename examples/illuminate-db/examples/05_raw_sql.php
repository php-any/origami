<?php
/**
 * 示例 5：原生 SQL 查询
 *
 * 演示使用 Capsule 直接执行原生 SQL 语句，
 * 包含查询、写入、更新与删除操作。
 *
 * 运行方式: ./illuminate-db examples/05_raw_sql.php
 */

require __DIR__ . '/../vendor/autoload.php';
require __DIR__ . '/../bootstrap.php';

use Illuminate\Database\Capsule\Manager as Capsule;

// 初始化 Capsule
bootstrap_capsule();

// 创建 users 表
Capsule::schema()->create('users', function ($table) {
    $table->increments('id');
    $table->string('name');
    $table->string('email')->unique();
    $table->integer('age');
});

// 插入示例数据
Capsule::table('users')->insert([
    ['name' => 'Alice', 'email' => 'alice@example.com', 'age' => 25],
    ['name' => 'Bob',   'email' => 'bob@example.com',   'age' => 30],
    ['name' => 'Carol', 'email' => 'carol@example.com', 'age' => 28],
    ['name' => 'David', 'email' => 'david@example.com', 'age' => 35],
]);

echo "=== 原生 SELECT ===\n";
// 基础查询
$rows = Capsule::select('SELECT * FROM users WHERE age >= ?', [28]);
echo "年龄 ≥ 28 的用户 (" . count($rows) . " 条):\n";
foreach ($rows as $row) {
    echo "  [{$row->id}] {$row->name} ({$row->age} 岁) <{$row->email}>\n";
}

// 聚合查询
$stats = Capsule::select('SELECT COUNT(*) as total, AVG(age) as avg_age FROM users');
if (count($stats) > 0) {
    echo "\n统计: 共 {$stats[0]->total} 人, 平均年龄 " . number_format($stats[0]->avg_age, 1) . "\n";
}

// 带参数的多条件查询
$young = Capsule::select(
    'SELECT * FROM users WHERE age BETWEEN ? AND ? ORDER BY age ASC',
    [20, 30]
);
echo "\n年龄 20-30 之间的用户 (" . count($young) . " 条):\n";
foreach ($young as $row) {
    echo "  {$row->name} ({$row->age} 岁)\n";
}

// 联表查询（演示 JOIN）
Capsule::schema()->create('posts', function ($table) {
    $table->increments('id');
    $table->unsignedInteger('user_id');
    $table->string('title');
});

Capsule::table('posts')->insert([
    ['user_id' => 1, 'title' => 'Alice 的文章 1'],
    ['user_id' => 1, 'title' => 'Alice 的文章 2'],
    ['user_id' => 2, 'title' => 'Bob 的文章'],
]);

$joined = Capsule::select(
    'SELECT u.name, p.title FROM users u JOIN posts p ON u.id = p.user_id ORDER BY u.name'
);
echo "\n联表查询 (用户 + 文章):\n";
foreach ($joined as $row) {
    echo "  {$row->name} -> {$row->title}\n";
}

echo "\n=== 原生 INSERT ===\n";
Capsule::insert('INSERT INTO users (name, email, age) VALUES (?, ?, ?)', ['Eve', 'eve@example.com', 22]);
echo "✓ 插入新用户 Eve\n";

$eve = Capsule::select('SELECT * FROM users WHERE email = ?', ['eve@example.com']);
if (count($eve) > 0) {
    echo "验证: {$eve[0]->name}, id={$eve[0]->id}\n";
}

echo "\n=== 原生 UPDATE ===\n";
Capsule::update('UPDATE users SET age = ? WHERE name = ?', [23, 'Eve']);
$eve = Capsule::select('SELECT * FROM users WHERE email = ?', ['eve@example.com']);
echo "Eve 年龄更新为: {$eve[0]->age}\n";

echo "\n=== 原生 DELETE ===\n";
Capsule::delete('DELETE FROM users WHERE name = ?', ['David']);
$david = Capsule::select('SELECT * FROM users WHERE name = ?', ['David']);
echo "删除 David 后剩余记录: " . (count($david) === 0 ? '已删除' : '仍存在') . "\n";

$total = Capsule::select('SELECT COUNT(*) as total FROM users');
echo "当前用户总数: {$total[0]->total}\n";

echo "\n=== 预处理语句复用 ===\n";
// 使用预处理语句重复执行相同查询
$statement = Capsule::connection()->getPdo()->prepare(
    'SELECT * FROM users WHERE age >= ? ORDER BY age'
);
$statement->execute([28]);
$matureUsers = $statement->fetchAll();
echo "使用 PDO 预处理查询年龄 ≥ 28 的用户: " . count($matureUsers) . " 条\n";

echo "\n✓ 原生 SQL 示例执行成功\n";
