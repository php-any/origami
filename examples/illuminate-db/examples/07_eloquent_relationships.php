<?php
/**
 * 示例 7：Eloquent 模型关联（Relationships）
 *
 * 演示使用 Eloquent ORM 定义模型关联并执行查询：
 * - `hasMany()` 一对多
 * - `belongsTo()` 反向一对多
 * - `hasOne()` 一对一
 * - 关联属性访问（`$user->posts`）
 * - 关联方法链式查询（`$user->posts()->where(...)->get()`）
 * - 预加载 eager loading（`with('posts')` 避免 N+1 查询）
 * - 关联创建（`$user->posts()->create([...])`）
 *
 * 运行方式: ./illuminate-db examples/07_eloquent_relationships.php
 */

require __DIR__ . '/../vendor/autoload.php';
require __DIR__ . '/../bootstrap.php';

use Illuminate\Database\Capsule\Manager as Capsule;
use Illuminate\Database\Eloquent\Model;

// 初始化 Capsule（内存 SQLite）
bootstrap_capsule();

/**
 * 用户模型（一对多 -> 文章）
 */
class User extends Model
{
    public $timestamps = false;
    protected $table = 'users';
    protected $fillable = ['name', 'email', 'age'];

    /**
     * 一对多：一个用户拥有多篇文章
     */
    public function posts()
    {
        return $this->hasMany(Post::class, 'user_id');
    }

    /**
     * 一对一：一个用户拥有一条个人资料
     */
    public function profile()
    {
        return $this->hasOne(Profile::class, 'user_id');
    }
}

/**
 * 文章模型（反向一对多 -> 用户）
 */
class Post extends Model
{
    public $timestamps = false;
    protected $table = 'posts';
    protected $fillable = ['user_id', 'title', 'content'];

    /**
     * 反向一对多：一篇文章属于一个用户
     */
    public function user()
    {
        return $this->belongsTo(User::class, 'user_id');
    }
}

/**
 * 个人资料模型（一对一 -> 用户）
 */
class Profile extends Model
{
    public $timestamps = false;
    protected $table = 'profiles';
    protected $fillable = ['user_id', 'bio', 'website'];

    public function user()
    {
        return $this->belongsTo(User::class, 'user_id');
    }
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

Capsule::schema()->create('profiles', function ($table) {
    $table->increments('id');
    $table->unsignedInteger('user_id');
    $table->string('bio')->nullable();
    $table->string('website')->nullable();
    $table->foreign('user_id')->references('id')->on('users');
});

echo "=== 准备数据 ===\n";
// 使用查询构建器插入用户（简化演示）
$aliceId = Capsule::table('users')->insertGetId([
    'name' => 'Alice', 'email' => 'alice@example.com', 'age' => 25,
]);
$bobId = Capsule::table('users')->insertGetId([
    'name' => 'Bob', 'email' => 'bob@example.com', 'age' => 30,
]);

// 使用关联方法创建文章
$alice = User::find($aliceId);
$alice->posts()->create(['title' => 'Alice 的文章 1', 'content' => '第一篇']);
$alice->posts()->create(['title' => 'Alice 的文章 2', 'content' => '第二篇']);
$alice->posts()->create(['title' => 'Alice 的文章 3', 'content' => '第三篇']);

$bob = User::find($bobId);
$bob->posts()->create(['title' => 'Bob 的文章', 'content' => 'Bob 的唯一一篇']);

// 使用关联方法创建个人资料
$alice->profile()->create(['bio' => 'Origami 开发者', 'website' => 'https://alice.example']);
$bob->profile()->create(['bio' => 'PHP 爱好者', 'website' => 'https://bob.example']);

echo "✓ 数据准备完成\n";

echo "\n=== 一对多 hasMany ===\n";
$alice = User::find($aliceId);
$alicePosts = $alice->posts; // 关联属性，返回集合
echo "Alice 的文章数: " . count($alicePosts) . "\n";
foreach ($alicePosts as $p) {
    echo "  - {$p->title}\n";
}

// 关联方法返回查询构建器，可链式过滤
$longPosts = $alice->posts()->where('title', 'like', '%2%')->get();
echo "Alice 标题含 '2' 的文章: " . count($longPosts) . " 篇\n";

echo "\n=== 反向关联 belongsTo ===\n";
$firstPost = Post::first();
echo "第一篇文章《{$firstPost->title}》作者: " . $firstPost->user->name . "\n";

echo "\n=== 一对一 hasOne ===\n";
$alice = User::find($aliceId);
echo "Alice 的简介: " . $alice->profile->bio . "\n";
echo "Alice 的网站: " . $alice->profile->website . "\n";

echo "\n=== 关联计数 withCount ===\n";
$usersWithCount = User::withCount('posts')->get();
foreach ($usersWithCount as $u) {
    echo "  {$u->name}: {$u->posts_count} 篇文章\n";
}

echo "\n=== 预加载 eager loading (with) ===\n";
// 一次性取出用户及其文章，避免 N+1 查询
$users = User::with('posts')->get();
foreach ($users as $u) {
    echo "  {$u->name} 的文章:\n";
    foreach ($u->posts as $p) {
        echo "    - {$p->title}\n";
    }
}

echo "\n=== 多层预加载 (with 嵌套) ===\n";
$posts = Post::with('user.profile')->get();
foreach ($posts as $p) {
    $profile = $p->user->profile;
    echo "  《{$p->title}》 作者 {$p->user->name} (网站: " . ($profile ? $profile->website : '无') . ")\n";
}

echo "\n=== 关联条件查询 has / doesntHave ===\n";
$hasPosts = User::has('posts')->get();
echo "有文章的用户: " . count($hasPosts) . " 个\n";

$noPosts = User::doesntHave('posts')->get();
echo "没有文章的用户: " . count($noPosts) . " 个\n";

echo "\n=== 关联删除 ===\n";
// 删除 Bob 及其关联文章（通过模型）
$bob = User::find($bobId);
$deletedPosts = $bob->posts()->delete();
echo "删除 Bob 的文章数: {$deletedPosts}\n";
$bob->delete();
echo "删除 Bob 后用户数: " . User::count() . "\n";

echo "\n=== 关联创建完整校验 ===\n";
$alice = User::find($aliceId);
$newPost = $alice->posts()->create([
    'title'   => '关联创建校验',
    'content' => '通过 $user->posts()->create() 创建的新文章',
]);
echo "新文章 id: {$newPost->id}, 标题: {$newPost->title}, user_id: {$newPost->user_id}\n";
echo "关联校验: 新文章属于 Alice? " . ($newPost->user->name === 'Alice' ? '是' : '否') . "\n";

echo "\n✓ Eloquent 模型关联示例执行成功\n";
