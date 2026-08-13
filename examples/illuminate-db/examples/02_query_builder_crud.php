<?php
/**
 * 示例 2：查询构建器 CRUD
 *
 * 演示使用 Illuminate Database 的查询构建器完成增删改查全套操作。
 *
 * 运行方式: ./illuminate-db examples/02_query_builder_crud.php
 */

require __DIR__ . '/../vendor/autoload.php';
require __DIR__ . '/../bootstrap.php';

use Illuminate\Database\Capsule\Manager as Capsule;

// 初始化 Capsule
bootstrap_capsule();

// 创建 products 表
Capsule::schema()->create('products', function ($table) {
    $table->increments('id');
    $table->string('name');
    $table->decimal('price', 10, 2);
    $table->integer('stock')->default(0);
    $table->boolean('active')->default(true);
});

echo "=== 插入 (Create) ===\n";
// 插入单条记录，获取自增 ID
$laptopId = Capsule::table('products')->insertGetId([
    'name'  => '笔记本电脑',
    'price' => 5999.00,
    'stock' => 10,
]);
echo "插入笔记本电脑: id={$laptopId}\n";

$phoneId = Capsule::table('products')->insertGetId([
    'name'  => '智能手机',
    'price' => 3999.00,
    'stock' => 20,
]);
echo "插入智能手机: id={$phoneId}\n";

// 批量插入多条记录
Capsule::table('products')->insert([
    ['name' => '蓝牙耳机', 'price' => 299.00,  'stock' => 50],
    ['name' => '机械键盘', 'price' => 699.00,  'stock' => 30],
    ['name' => '显示器',   'price' => 1999.00, 'stock' => 15],
]);
echo "批量插入 3 个商品\n";

echo "\n=== 查询 (Read) ===\n";
$all = Capsule::table('products')->get();
echo "全部商品 (" . count($all) . " 条):\n";
foreach ($all as $p) {
    $status = $p->active ? '在售' : '停售';
    echo "  [{$p->id}] {$p->name} ¥{$p->price} (库存 {$p->stock}) [{$status}]\n";
}

// 条件 + 排序 + 限制
$expensive = Capsule::table('products')
    ->where('price', '>=', 1000)
    ->orderBy('price', 'desc')
    ->get();
echo "\n价格 ≥ ¥1000 的商品 (按价格降序):\n";
foreach ($expensive as $p) {
    echo "  {$p->name} ¥{$p->price}\n";
}

// 分页
$page1 = Capsule::table('products')
    ->orderBy('id')
    ->offset(0)
    ->limit(2)
    ->get();
echo "\n第 1 页 (每页 2 条):\n";
foreach ($page1 as $p) {
    echo "  {$p->name}\n";
}

// 聚合函数
$avgPrice  = Capsule::table('products')->avg('price');
$totalStock = Capsule::table('products')->sum('stock');
$maxPrice  = Capsule::table('products')->max('price');
echo "\n平均价格: " . number_format($avgPrice, 2) . "\n";
echo "总库存: {$totalStock}\n";
echo "最高价格: {$maxPrice}\n";

echo "\n=== 更新 (Update) ===\n";
// 更新指定记录
Capsule::table('products')
    ->where('id', $laptopId)
    ->update(['price' => 5499.00, 'stock' => 8]);

$laptop = Capsule::table('products')->where('id', $laptopId)->first();
echo "更新后: {$laptop->name} ¥{$laptop->price} (库存 {$laptop->stock})\n";

// 自增 / 自减
Capsule::table('products')->where('id', $phoneId)->decrement('stock', 2);
$phone = Capsule::table('products')->where('id', $phoneId)->first();
echo "手机减库存 2 后: 库存 {$phone->stock}\n";

Capsule::table('products')->where('id', $phoneId)->increment('stock', 5);
$phone = Capsule::table('products')->where('id', $phoneId)->first();
echo "手机加库存 5 后: 库存 {$phone->stock}\n";

echo "\n=== 删除 (Delete) ===\n";
Capsule::table('products')->where('name', '显示器')->delete();

$count = Capsule::table('products')->count();
echo "删除后商品总数: {$count}\n";

echo "\n✓ 查询构建器 CRUD 示例执行成功\n";
