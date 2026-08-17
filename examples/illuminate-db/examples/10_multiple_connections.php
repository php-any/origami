<?php
/**
 * 示例 10：多数据库连接管理
 *
 * 演示使用 Capsule 同时管理多个数据库连接：
 * - 添加多个连接（默认连接 + 命名连接）
 * - `Capsule::connection()` 切换连接
 * - 跨连接查询
 * - `$capsule->setAsGlobal()` 后通过 `Capsule::connection('xxx')` 访问
 * - Schema 在不同连接上操作
 * - `getDefaultConnection()` / `setDefaultConnection()`
 *
 * 运行方式: ./illuminate-db examples/10_multiple_connections.php
 */

require __DIR__ . '/../vendor/autoload.php';
require __DIR__ . '/../bootstrap.php';

use Illuminate\Database\Capsule\Manager as Capsule;
use Illuminate\Database\Eloquent\Model;

// 手动创建 Capsule 实例，添加两个连接
$capsule = new Capsule();

// 默认连接（sqlite 内存库）
$capsule->addConnection([
    'driver'   => 'sqlite',
    'database' => ':memory:',
], 'default');

// 第二个命名连接：主数据库
$capsule->addConnection([
    'driver'   => 'sqlite',
    'database' => ':memory:',
], 'primary');

// 第三个命名连接：报表/日志库
$capsule->addConnection([
    'driver'   => 'sqlite',
    'database' => ':memory:',
], 'report');

$capsule->setAsGlobal();
$capsule->bootEloquent();

echo "=== 连接概览 ===\n";
echo "默认连接: " . Capsule::connection()->getName() . "\n";
echo "当前连接 PDO: " . get_class(Capsule::connection()->getPdo()) . "\n";

echo "\n=== 在不同连接上创建表 ===\n";
// 默认连接：创建 users 表
Capsule::schema('default')->create('users', function ($table) {
    $table->increments('id');
    $table->string('name');
});

// primary 连接：创建 users 表
Capsule::schema('primary')->create('users', function ($table) {
    $table->increments('id');
    $table->string('name');
});

// report 连接：创建 logs 表
Capsule::schema('report')->create('logs', function ($table) {
    $table->increments('id');
    $table->string('message');
});

echo "✓ 三个连接均已完成建表\n";

echo "\n=== 跨连接写入数据 ===\n";
// 默认连接写入
Capsule::table('users')->insert(['name' => 'default-user']);
// primary 连接写入
Capsule::connection('primary')->table('users')->insert(['name' => 'primary-user']);
// report 连接写入
Capsule::connection('report')->table('logs')->insert(['message' => 'report-log-1']);

echo "✓ 各连接数据写入完成\n";

echo "\n=== 读取各连接数据 ===\n";
$defaultUsers = Capsule::table('users')->get();
echo "默认连接 users: " . count($defaultUsers) . " 条 (";
foreach ($defaultUsers as $u) { echo $u->name . " "; }
echo ")\n";

$primaryUsers = Capsule::connection('primary')->table('users')->get();
echo "primary 连接 users: " . count($primaryUsers) . " 条 (";
foreach ($primaryUsers as $u) { echo $u->name . " "; }
echo ")\n";

$reportLogs = Capsule::connection('report')->table('logs')->get();
echo "report 连接 logs: " . count($reportLogs) . " 条 (";
foreach ($reportLogs as $l) { echo $l->message . " "; }
echo ")\n";

echo "\n=== 动态切换默认连接 ===\n";
// 通过 DatabaseManager 访问默认连接的读取与设置（Capsule 静态代理不直接暴露）
$dbManager = $capsule->getDatabaseManager();
$previous = $dbManager->getDefaultConnection();
$dbManager->setDefaultConnection('primary');
echo "切换后默认连接: " . $dbManager->getDefaultConnection() . "\n";
$users = Capsule::table('users')->get();
echo "不指定连接查询（此时走 primary）: " . count($users) . " 条\n";

// 恢复默认连接
$dbManager->setDefaultConnection($previous);
echo "恢复默认连接: " . $dbManager->getDefaultConnection() . "\n";

echo "\n=== 使用 raw 在指定连接执行 SQL ===\n";
$reportRows = Capsule::connection('report')->select('SELECT * FROM logs WHERE message = ?', ['report-log-1']);
echo "report 连接原生 SQL 查询: " . count($reportRows) . " 条\n";

echo "\n=== 查看连接配置 ===\n";
$config = Capsule::connection('primary')->getConfig();
echo "primary 连接的 driver: " . $config['driver'] . "\n";

echo "\n=== 获取所有已注册连接 ===\n";
$pdo = Capsule::connection()->getPdo();
echo "默认连接 PDO 已建立\n";

echo "\n=== 断开与重连 ===\n";
$pdoDisconnected = Capsule::connection('report');
// 断开连接
Capsule::connection('report')->disconnect();
echo "已断开 report 连接\n";
// 重新连接（懒加载）；注意：内存 SQLite 断开后数据会清空，需重建表
Capsule::schema('report')->create('logs', function ($table) {
    $table->increments('id');
    $table->string('message');
});
Capsule::connection('report')->table('logs')->insert(['message' => 'reconnect-log']);
$again = Capsule::connection('report')->select('SELECT * FROM logs');
echo "重新连接后查询 logs: " . count($again) . " 条\n";

echo "\n✓ 多数据库连接示例执行成功\n";
