<?php
/**
 * 示例 3：Schema Builder（表结构定义）
 *
 * 演示使用 Schema Builder 创建表、修改表结构、定义字段类型与约束。
 *
 * 运行方式: ./illuminate-db examples/03_schema_builder.php
 */

require __DIR__ . '/../vendor/autoload.php';
require __DIR__ . '/../bootstrap.php';

use Illuminate\Database\Capsule\Manager as Capsule;

// 初始化 Capsule
bootstrap_capsule();

echo "=== 创建表 ===\n";

// 创建 users 表（包含常用字段类型）
Capsule::schema()->create('users', function ($table) {
    $table->increments('id');           // 自增主键
    $table->string('name', 100);        // 字符串（可指定长度）
    $table->string('email')->unique();  // 唯一索引
    $table->text('bio')->nullable();    // 文本，可空
    $table->integer('age')->unsigned(); // 无符号整数
    $table->decimal('balance', 10, 2)->default(0); // 十进制数
    $table->boolean('is_active')->default(true);   // 布尔值
    $table->date('birthday')->nullable();          // 日期
    $table->dateTime('last_login')->nullable();    // 日期时间
    $table->json('preferences')->nullable();       // JSON
    $table->timestamps();                          // created_at + updated_at
});
echo "✓ 已创建 users 表\n";

// 创建 posts 表（带外键关联）
Capsule::schema()->create('posts', function ($table) {
    $table->increments('id');
    $table->unsignedInteger('user_id');
    $table->string('title');
    $table->text('content');
    $table->enum('status', ['draft', 'published', 'archived'])->default('draft');
    $table->timestamps();

    // 外键关联 users 表
    $table->foreign('user_id')->references('id')->on('users')->onDelete('cascade');
});
echo "✓ 已创建 posts 表（含外键）\n";

echo "\n=== 检查表是否存在 ===\n";
echo "users 表存在: " . (Capsule::schema()->hasTable('users') ? '是' : '否') . "\n";
echo "posts 表存在: " . (Capsule::schema()->hasTable('posts') ? '是' : '否') . "\n";
echo "products 表存在: " . (Capsule::schema()->hasTable('products') ? '是' : '否') . "\n";

echo "\n=== 检查字段是否存在 ===\n";
echo "users.name 存在: " . (Capsule::schema()->hasColumn('users', 'name') ? '是' : '否') . "\n";
echo "users.email 存在: " . (Capsule::schema()->hasColumn('users', 'email') ? '是' : '否') . "\n";

echo "\n=== 修改表结构 ===\n";
// 新增字段
Capsule::schema()->table('users', function ($table) {
    $table->string('phone')->nullable()->after('email');
});
echo "✓ 已添加 users.phone 字段\n";

// 检查新字段
echo "users.phone 存在: " . (Capsule::schema()->hasColumn('users', 'phone') ? '是' : '否') . "\n";

echo "\n=== 插入数据验证 Schema ===\n";
$userId = Capsule::table('users')->insertGetId([
    'name'     => '张三',
    'email'    => 'zhangsan@example.com',
    'bio'      => 'Origami 开发者',
    'age'      => 28,
    'balance'  => 1000.50,
    'is_active'=> true,
    'birthday' => '1997-03-15',
    'phone'    => '13800138000',
]);

Capsule::table('posts')->insert([
    'user_id' => $userId,
    'title'   => 'Origami 使用指南',
    'content' => '本文介绍如何在 Origami 中集成 Illuminate Database。',
    'status'  => 'published',
]);

$user = Capsule::table('users')->where('id', $userId)->first();
echo "用户: {$user->name}, 邮箱: {$user->email}, 余额: ¥{$user->balance}\n";
echo "活跃: " . ($user->is_active ? '是' : '否') . ", 生日: {$user->birthday}\n";
echo "简介: {$user->bio}\n";

$post = Capsule::table('posts')->where('user_id', $userId)->first();
echo "文章: {$post->title} [{$post->status}]\n";

echo "\n=== 删除表 ===\n";
// 删除 posts 表（先删除有外键依赖的表）
Capsule::schema()->drop('posts');
echo "✓ 已删除 posts 表\n";
echo "posts 表存在: " . (Capsule::schema()->hasTable('posts') ? '是' : '否') . "\n";

// 使用 dropIfExists
Capsule::schema()->dropIfExists('users');
echo "✓ 已删除 users 表\n";
echo "users 表存在: " . (Capsule::schema()->hasTable('users') ? '是' : '否') . "\n";

echo "\n✓ Schema Builder 示例执行成功\n";
