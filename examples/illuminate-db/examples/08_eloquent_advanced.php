<?php
/**
 * 示例 8：Eloquent 高级特性
 *
 * 演示 Eloquent ORM 的高级能力：
 * - `Model::create()` 批量创建（使用 $fillable）
 * - `firstOrCreate()` / `firstOrNew()` 存在则取、不存在则创建
 * - `updateOrCreate()` 存在则更新、不存在则创建
 * - `firstOrNew()` 结合 `save()` 实现自定义创建流程
 * - 访问器（Accessor）`getXxxAttribute()`
 * - 修改器（Mutator）`setXxxAttribute()`
 * - 追加访问器到序列化数组（$appends）
 * - 局部作用域（Local Scope）`scopeXxx()`
 * - 软删除（SoftDeletes）
 * - 模型事件（booted / static::creating 等）
 *
 * 运行方式: ./illuminate-db examples/08_eloquent_advanced.php
 */

require __DIR__ . '/../vendor/autoload.php';
require __DIR__ . '/../bootstrap.php';

use Illuminate\Database\Capsule\Manager as Capsule;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\SoftDeletes;

// 初始化 Capsule（内存 SQLite）
bootstrap_capsule();

/**
 * 产品模型
 */
class Product extends Model
{
    use SoftDeletes; // 启用软删除

    public $timestamps = true; // 自动维护 created_at / updated_at
    protected $table = 'products';
    protected $fillable = ['name', 'price', 'stock', 'is_active'];

    // 追加访问器字段到数组/JSON 序列化
    protected $appends = ['price_label'];

    /**
     * 访问器：格式化价格显示
     */
    public function getPriceLabelAttribute()
    {
        return '¥' . number_format($this->price, 2);
    }

    /**
     * 修改器：入库时统一格式化价格（保留两位小数）
     */
    public function setPriceAttribute($value)
    {
        $this->attributes['price'] = round((float) $value, 2);
    }

    /**
     * 访问器：判断是否在售
     */
    public function getActiveLabelAttribute()
    {
        return $this->is_active ? '在售' : '停售';
    }

    /**
     * 局部作用域：仅查询在售商品
     */
    public function scopeActive($query)
    {
        return $query->where('is_active', true);
    }

    /**
     * 局部作用域：价格区间
     */
    public function scopePriceBetween($query, $min, $max)
    {
        return $query->whereBetween('price', [$min, $max]);
    }

    /**
     * 模型事件：创建前自动记录
     */
    protected static function booted()
    {
        static::creating(function ($product) {
            // 创建前回调（此处仅演示，不实际做操作）
        });
    }
}

// 创建表（含 deleted_at 软删除字段）
Capsule::schema()->create('products', function ($table) {
    $table->increments('id');
    $table->string('name');
    $table->decimal('price', 10, 2);
    $table->integer('stock')->default(0);
    $table->boolean('is_active')->default(true);
    $table->timestamps();
    $table->softDeletes(); // deleted_at 字段，软删除支持
});

echo "=== Model::create() 批量创建 ===\n";
$laptop = Product::create([
    'name'  => '笔记本电脑',
    'price' => 5999.00,
    'stock' => 10,
]);
echo "✓ 创建商品: {$laptop->name}, 价格 " . $laptop->price_label . "\n";

$phone = Product::create([
    'name'  => '智能手机',
    'price' => 3999.00,
    'stock' => 20,
]);
echo "✓ 创建商品: {$phone->name}, 价格 " . $phone->price_label . "\n";

echo "\n=== 访问器 / 修改器 ===\n";
// 修改器：设置价格会自动 round 到两位小数
$laptop->price = 5499.9999;
$laptop->save();
echo "修改器验证: 设置 5499.9999 后实际存储 " . $laptop->price . "\n";

// 访问器：price_label
echo "访问器验证: " . $laptop->price_label . "\n";

// 追加字段出现在数组/JSON 中
$laptopArray = $laptop->toArray();
echo "追加字段 toArray()['price_label']: " . $laptopArray['price_label'] . "\n";

echo "\n=== 局部作用域 (Scopes) ===\n";
// 停售一个商品
$phone->is_active = false;
$phone->save();

$activeProducts = Product::active()->get();
echo "在售商品数: " . count($activeProducts) . "\n";

$midRange = Product::priceBetween(3000, 6000)->get();
echo "价格 3000-6000 的商品: " . count($midRange) . " 个\n";
foreach ($midRange as $p) {
    echo "  - {$p->name} " . $p->price_label . " [" . $p->active_label . "]\n";
}

echo "\n=== firstOrCreate / updateOrCreate ===\n";
// 已存在则返回现有记录，不存在则创建
$found = Product::firstOrCreate(['name' => '笔记本电脑'], ['price' => 1, 'stock' => 1]);
echo "firstOrCreate（已存在）: id={$found->id}, 价格 " . $found->price_label . "（保持原数据）\n";

$created = Product::firstOrCreate(
    ['name' => '机械键盘'],
    ['price' => 699.00, 'stock' => 30]
);
echo "firstOrCreate（不存在，新创建）: id={$created->id}, 名称={$created->name}\n";

// 存在则更新，不存在则创建
$updated = Product::updateOrCreate(
    ['name' => '智能手机'],
    ['price' => 3799.00, 'stock' => 25]
);
echo "updateOrCreate（已存在则更新）: {$updated->name} 新价格 " . $updated->price_label . "\n";

// firstOrNew 不自动保存，需手动 save()，适合自定义创建流程
$fn = Product::firstOrNew(['name' => '平板电脑']);
if (!$fn->exists) {
    $fn->fill(['price' => 2999.00, 'stock' => 5])->save();
    echo "firstOrNew（不存在则手动 save() 创建）: id={$fn->id}, 名称={$fn->name}\n";
} else {
    echo "firstOrNew（已存在则直接复用）: id={$fn->id}\n";
}

echo "\n=== 软删除 (SoftDeletes) ===\n";
$before = Product::count();
echo "删除前商品总数: {$before}\n";

// 软删除：不会真正从数据库删除，而是设置 deleted_at
$keyboard = Product::where('name', '机械键盘')->first();
$keyboard->delete();

echo "软删除后正常查询数: " . Product::count() . "\n";
echo "包含已删除(withTrashed): " . Product::withTrashed()->count() . "\n";

$trashed = Product::onlyTrashed()->get();
echo "仅查已删除(onlyTrashed): " . count($trashed) . " 个\n";
foreach ($trashed as $p) {
    echo "  - {$p->name} (deleted_at={$p->deleted_at})\n";
}

// 恢复软删除
$trashed->first()->restore();
echo "恢复后商品总数: " . Product::count() . "\n";

echo "\n=== 强制删除 (forceDelete) ===\n";
$tmp = Product::create(['name' => '临时商品', 'price' => 1.00, 'stock' => 1]);
$tmp->forceDelete();
echo "forceDelete 后（正常+已删除）: " . Product::withTrashed()->where('name', '临时商品')->count() . "\n";

echo "\n=== 模型时间戳 ===\n";
$first = Product::first();
echo "created_at: {$first->created_at}\n";
echo "updated_at: {$first->updated_at}\n";

echo "\n✓ Eloquent 高级特性示例执行成功\n";
