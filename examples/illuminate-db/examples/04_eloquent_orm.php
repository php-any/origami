<?php
/**
 * 示例 4：Eloquent ORM（对象关系映射）
 *
 * 演示使用 Eloquent ORM 定义模型、执行查询、更新与删除操作。
 *
 * ⚠️ 注意：由于 Origami 运行时对 __callStatic 魔术方法中数组参数传递的
 * 支持仍在完善中，建议通过 `Model::query()` 获取查询构建器后再进行
 * 链式操作，避免使用 `User::where()` 这类直接静态调用的方式。
 *
 * 运行方式: ./illuminate-db examples/04_eloquent_orm.php
 */

require __DIR__ . '/../vendor/autoload.php';
require __DIR__ . '/../bootstrap.php';

use Illuminate\Database\Capsule\Manager as Capsule;
use Illuminate\Database\Eloquent\Model;

// 初始化 Capsule（内存 SQLite）
bootstrap_capsule();

/**
 * 用户模型
 */
class User extends Model
{
    public $timestamps = false;
    protected $table = 'users';
    protected $fillable = ['name', 'email', 'age'];
}

/**
 * 文章模型
 */
class Post extends Model
{
    public $timestamps = false;
    protected $table = 'posts';
    protected $fillable = ['user_id', 'title', 'content'];
}

// 创建表
Capsule::schema()->create('users', function ($table) {
    $table->increments('id');
    $table->string('name');
    $table->string('email')->unique();
    $table->integer('age')->nullable();
});

Capsule::schema()->create('posts', function ($table) {
    $table->increments('id');
    $table->unsignedInteger('user_id');
    $table->string('title');
    $table->text('content');
    $table->foreign('user_id')->references('id')->on('users');
});

echo "=== 插入数据（通过查询构建器）===\n";
// 由于 Eloquent save() 方法在 Origami 运行时存在兼容性问题，
// 这里使用查询构建器插入数据，随后用 Eloquent 进行读取。
$aliceId = Capsule::table('users')->insertGetId([
    'name'  => 'Alice',
    'email' => 'alice@example.com',
    'age'   => 25,
]);
$bobId = Capsule::table('users')->insertGetId([
    'name'  => 'Bob',
    'email' => 'bob@example.com',
    'age'   => 30,
]);

Capsule::table('posts')->insert([
    'user_id' => $aliceId,
    'title'   => 'Origami 介绍',
    'content' => 'Origami 是融合 PHP 与 Go 的混合脚本语言。',
]);
Capsule::table('posts')->insert([
    'user_id' => $aliceId,
    'title'   => 'Illuminate Database 独立使用',
    'content' => '不依赖 Laravel 框架，独立使用 Eloquent ORM。',
]);
Capsule::table('posts')->insert([
    'user_id' => $bobId,
    'title'   => '查询构建器技巧',
    'content' => '高效使用查询构建器的技巧总结。',
]);
echo "✓ 数据插入完成\n";

echo "\n=== Eloquent 查询 (Read) ===\n";
// 查询所有用户（推荐通过 query() 获取构建器）
$users = User::query()->get();
echo "所有用户 (" . count($users) . " 条):\n";
foreach ($users as $u) {
    echo "  [{$u->id}] {$u->name} <{$u->email}> (年龄 {$u->age})\n";
}

// 条件查询：匹配特定名称
$alice = User::query()->where('name', 'Alice')->first();
echo "\n条件查询 Alice: " . ($alice ? $alice->email : '未找到') . "\n";

// 按主键查找
$found = User::query()->find(2);
echo "find(2): " . ($found ? $found->name : '未找到') . "\n";

// 排序
$sorted = User::query()->orderBy('age', 'desc')->get();
echo "\n按年龄降序:\n";
foreach ($sorted as $u) {
    echo "  {$u->name} ({$u->age} 岁)\n";
}

// 聚合统计
$count   = User::query()->count();
$avgAge  = User::query()->avg('age');
echo "\n用户数: {$count}, 平均年龄: " . number_format($avgAge, 1) . "\n";

echo "\n=== Eloquent 更新 (Update) ===\n";
// 通过 Eloquent 查询构建器更新
User::query()->where('id', $aliceId)->update(['age' => 26]);

$updatedAlice = User::query()->find($aliceId);
echo "Alice 更新后: {$updatedAlice->name}, 年龄 {$updatedAlice->age}\n";

// 通过查询构建器更新模型数据
User::query()->where('id', $aliceId)->update(['name' => 'Alice Wonderland']);
echo "Alice 改名后: " . User::query()->find($aliceId)->name . "\n";

echo "\n=== Eloquent 删除 (Delete) ===\n";
// 删除 Bob 及其文章
Post::query()->where('user_id', $bobId)->delete();
User::query()->where('id', $bobId)->delete();

$remainingUsers = User::query()->count();
$remainingPosts = Post::query()->count();
echo "删除后用户数: {$remainingUsers}, 文章数: {$remainingPosts}\n";

echo "\n=== 模型属性操作 ===\n";
$model = new User();
$model->fill(['name' => 'Carol', 'email' => 'carol@example.com', 'age' => 28]);
echo "通过 fill() 设置属性:\n";
echo "  name: {$model->name}\n";
echo "  email: {$model->email}\n";
echo "  age: {$model->age}\n";
echo "  isDirty(): " . ($model->isDirty() ? '是' : '否') . "\n";
echo "  getAttributes():\n";
$attrs = $model->getAttributes();
foreach ($attrs as $k => $v) {
    echo "    {$k} => {$v}\n";
}

echo "\n✓ Eloquent ORM 示例执行成功\n";
