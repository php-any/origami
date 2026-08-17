<?php
/**
 * 示例 9：查询构建器高级用法
 *
 * 演示 Illuminate Database 查询构建器的进阶能力：
 * - 多表关联查询 `join()` / `leftJoin()`
 * - `whereIn()` / `whereNotIn()` / `whereBetween()` / `whereNull()`
 * - `pluck()` 提取单列、`value()` 取单值
 * - 分组聚合 `groupBy()` + `having()`
 * - 自增自减 `increment()` / `decrement()`
 * - 批量更新 `updateOrInsert()`
 * - 大规模分块处理 `chunkById()`
 * - 锁 / `toSql()` 打印生成 SQL
 *
 * 运行方式: ./illuminate-db examples/09_query_builder_advanced.php
 */

require __DIR__ . '/../vendor/autoload.php';
require __DIR__ . '/../bootstrap.php';

use Illuminate\Database\Capsule\Manager as Capsule;

// 初始化 Capsule
bootstrap_capsule();

// 创建用户表
Capsule::schema()->create('users', function ($table) {
    $table->increments('id');
    $table->string('name');
    $table->integer('age');
    $table->string('city');
});

// 创建订单表（外键关联用户）
Capsule::schema()->create('orders', function ($table) {
    $table->increments('id');
    $table->unsignedInteger('user_id');
    $table->decimal('amount', 10, 2);
    $table->string('status')->default('pending');
    $table->foreign('user_id')->references('id')->on('users');
});

// 插入示例用户
Capsule::table('users')->insert([
    ['name' => 'Alice', 'age' => 25, 'city' => '北京'],
    ['name' => 'Bob',   'age' => 30, 'city' => '上海'],
    ['name' => 'Carol', 'age' => 28, 'city' => '广州'],
    ['name' => 'David', 'age' => 35, 'city' => '北京'],
    ['name' => 'Eve',   'age' => 22, 'city' => '深圳'],
]);

// 插入示例订单
Capsule::table('orders')->insert([
    ['user_id' => 1, 'amount' => 199.00,  'status' => 'completed'],
    ['user_id' => 1, 'amount' => 299.00,  'status' => 'pending'],
    ['user_id' => 2, 'amount' => 499.00,  'status' => 'completed'],
    ['user_id' => 2, 'amount' => 1299.00, 'status' => 'completed'],
    ['user_id' => 3, 'amount' => 99.00,   'status' => 'cancelled'],
    ['user_id' => 4, 'amount' => 899.00,  'status' => 'completed'],
]);

echo "=== 多表关联查询 join ===\n";
// 内连接：查询所有有订单的用户及订单金额
$joined = Capsule::table('orders')
    ->join('users', 'users.id', '=', 'orders.user_id')
    ->select('users.name', 'orders.amount', 'orders.status')
    ->get();
echo "内连接 (orders JOIN users):\n";
foreach ($joined as $row) {
    echo "  {$row->name} 订单 ¥{$row->amount} [{$row->status}]\n";
}

echo "\n=== leftJoin 左连接 ===\n";
// 左连接：查出所有用户及其订单金额（无订单的用户也出现）
$left = Capsule::table('users')
    ->leftJoin('orders', 'users.id', '=', 'orders.user_id')
    ->select('users.name', 'orders.amount')
    ->orderBy('users.id')
    ->get();
echo "左连接 (users LEFT JOIN orders):\n";
foreach ($left as $row) {
    echo "  {$row->name}: " . ($row->amount !== null ? '¥' . $row->amount : '无订单') . "\n";
}

echo "\n=== whereIn / whereNotIn / whereBetween / whereNull ===\n";
$inCity = Capsule::table('users')->whereIn('city', ['北京', '上海'])->get();
echo "whereIn(北京, 上海): " . count($inCity) . " 人\n";
foreach ($inCity as $u) {
    echo "  - {$u->name} ({$u->city})\n";
}

$notIn = Capsule::table('users')->whereNotIn('city', ['北京'])->get();
echo "whereNotIn(北京): " . count($notIn) . " 人\n";

$between = Capsule::table('users')->whereBetween('age', [25, 30])->get();
echo "whereBetween(age, 25-30): " . count($between) . " 人\n";

$nullTest = Capsule::table('users')
    ->leftJoin('orders', 'users.id', '=', 'orders.user_id')
    ->whereNull('orders.id')
    ->select('users.name')
    ->get();
echo "whereNull（无订单的用户）: " . count($nullTest) . " 人\n";

echo "\n=== pluck / value ===\n";
$names = Capsule::table('users')->pluck('name');
echo "pluck('name') 提取所有用户名: " . $names->implode(', ') . "\n";

$cityValue = Capsule::table('users')->where('name', 'Alice')->value('city');
echo "value('city') 取 Alice 所在城市: {$cityValue}\n";

// pluck 指定 key/value
$idNameMap = Capsule::table('users')->pluck('name', 'id');
echo "pluck('name', 'id') 生成 id=>name 映射:\n";
foreach ($idNameMap as $id => $name) {
    echo "  {$id} => {$name}\n";
}

echo "\n=== 分组聚合 groupBy + having ===\n";
$byCity = Capsule::table('users')
    ->select('city', Capsule::raw('COUNT(*) as cnt'))
    ->groupBy('city')
    ->get();
echo "按城市分组统计:\n";
foreach ($byCity as $row) {
    echo "  {$row->city}: {$row->cnt} 人\n";
}

// having 过滤分组结果（统计订单金额总和 > 500 的用户）
$bigSpenders = Capsule::table('orders')
    ->select('user_id', Capsule::raw('SUM(amount) as total'))
    ->groupBy('user_id')
    ->havingRaw('SUM(amount) > ?', [500])
    ->get();
echo "\n消费总额 > 500 的用户 (havingRaw):\n";
foreach ($bigSpenders as $row) {
    echo "  user_id={$row->user_id} 总额 ¥{$row->total}\n";
}

echo "\n=== 自增自减 increment / decrement ===\n";
// 给所有北京用户年龄 +1
$affected = Capsule::table('users')->where('city', '北京')->increment('age', 1);
echo "北京用户年龄 +1，影响 {$affected} 行\n";

$alice = Capsule::table('users')->where('name', 'Alice')->first();
echo "Alice 现在年龄: {$alice->age}（25 -> 26）\n";

echo "\n=== updateOrInsert ===\n";
// 存在则更新，不存在则插入
$result = Capsule::table('users')->updateOrInsert(
    ['name' => 'Frank'],
    ['age' => 40, 'city' => '杭州']
);
echo "updateOrInsert(Frank，不存在则插入): " . ($result ? '已插入' : '已更新') . "\n";

$result = Capsule::table('users')->updateOrInsert(
    ['name' => 'Alice'],
    ['age' => 27, 'city' => '北京']
);
echo "updateOrInsert(Alice，存在则更新): " . ($result ? '已插入' : '已更新') . "\n";
$alice = Capsule::table('users')->where('name', 'Alice')->first();
echo "Alice 更新后年龄: {$alice->age}\n";

echo "\n=== chunkById 分块处理 ===\n";
// 分批处理大量数据，避免一次性加载
$processed = 0;
Capsule::table('users')
    ->orderBy('id')
    ->chunkById(2, function ($users) use (&$processed) {
        foreach ($users as $u) {
            $processed++;
        }
        echo "  处理了一批 " . count($users) . " 条\n";
    });
echo "chunkById 共处理: {$processed} 条\n";

echo "\n=== toSql 查看生成的 SQL ===\n";
$sql = Capsule::table('users')
    ->where('age', '>', 20)
    ->orderBy('age', 'desc')
    ->toSql();
echo "生成的 SQL: {$sql}\n";

echo "\n=== 结果转为数组 ===\n";
$arrayResult = (array) Capsule::table('users')->where('name', 'Alice')->first();
echo "first() 转数组: ";
var_export($arrayResult);
echo "\n";

echo "\n✓ 查询构建器高级用法示例执行成功\n";
