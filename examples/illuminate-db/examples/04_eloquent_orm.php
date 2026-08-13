<?php
/**
 * 示例 4：Eloquent ORM（对象关系映射）
 *
 * 演示使用 Eloquent ORM 定义模型、执行查询、更新与删除操作。
 *
 * 覆盖的能力：
 * - `$model->save()` 插入新记录
 * - `User::where()` 静态调用查询构建器
 * - `update()` / `delete()` 返回实际影响行数
 * - `User::query()` 链式查询、聚合、find()
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

echo "=== 插入数据（通过 Eloquent save()）===\n";
// 使用 Eloquent 模型实例 + save() 插入记录（Origami 已支持关联数组键传递）。
$alice = new User();
$alice->name  = 'Alice';
$alice->email = 'alice@example.com';
$alice->age   = 25;
$alice->save();
$aliceId = $alice->id;

$bob = new User();
$bob->name  = 'Bob';
$bob->email = 'bob@example.com';
$bob->age   = 30;
$bob->save();
$bobId = $bob->id;

// 通过 fill() + save() 插入
$post = new Post();
$post->fill([
    'user_id' => $aliceId,
    'title'   => 'Origami 介绍',
    'content' => 'Origami 是融合 PHP 与 Go 的混合脚本语言。',
]);
$post->save();

$post2 = new Post();
$post2->fill([
    'user_id' => $aliceId,
    'title'   => 'Illuminate Database 独立使用',
    'content' => '不依赖 Laravel 框架，独立使用 Eloquent ORM。',
]);
$post2->save();

$post3 = new Post();
$post3->fill([
    'user_id' => $bobId,
    'title'   => '查询构建器技巧',
    'content' => '高效使用查询构建器的技巧总结。',
]);
$post3->save();

echo "✓ 数据插入完成（Alice id={$aliceId}, Bob id={$bobId}）\n";

echo "\n=== Eloquent 查询 (Read) ===\n";
// 查询所有用户
$users = User::query()->get();
echo "所有用户 (" . count($users) . " 条):\n";
foreach ($users as $u) {
    echo "  [{$u->id}] {$u->name} <{$u->email}> (年龄 {$u->age})\n";
}

// 静态调用 where() 进行条件查询（Origami 已支持 __callStatic 参数传递）
$alice = User::where('name', 'Alice')->first();
echo "\nUser::where('name','Alice')->first(): " . ($alice ? $alice->email : '未找到') . "\n";

// 按主键查找
$found = User::where('id', $aliceId)->first();
echo "User::where('id', {$aliceId})->first(): " . ($found ? $found->name : '未找到') . "\n";

// 排序
$sorted = User::orderBy('age', 'desc')->get();
echo "\nUser::orderBy('age','desc')->get():\n";
foreach ($sorted as $u) {
    echo "  {$u->name} ({$u->age} 岁)\n";
}

// 聚合统计
$count   = User::where('age', '>', 20)->count();
$avgAge  = User::avg('age');
echo "\nUser::where('age','>',20)->count(): {$count}, User::avg('age'): " . number_format($avgAge, 1) . "\n";

echo "\n=== Eloquent 更新 (Update) ===\n";
// 通过静态调用更新，并返回实际影响行数
$affected = User::where('id', $aliceId)->update(['age' => 26]);
echo "User::where('id',{$aliceId})->update(['age'=>26]) 影响行数: {$affected}\n";

$updatedAlice = User::find($aliceId);
echo "Alice 更新后: {$updatedAlice->name}, 年龄 {$updatedAlice->age}\n";

$affected2 = User::where('id', $aliceId)->update(['name' => 'Alice Wonderland']);
echo "改名影响行数: {$affected2}, 新名字: " . User::find($aliceId)->name . "\n";

echo "\n=== Eloquent 删除 (Delete) ===\n";
// 删除 Bob 及其文章，并返回实际影响行数
$delPosts = Post::where('user_id', $bobId)->delete();
$delUsers = User::where('id', $bobId)->delete();
echo "删除文章影响行数: {$delPosts}, 删除用户影响行数: {$delUsers}\n";

$remainingUsers = User::count();
$remainingPosts = Post::count();
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
