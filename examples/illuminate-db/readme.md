# Illuminate Database 独立使用示例

本目录演示如何在 **不依赖 Laravel 框架** 的前提下，独立使用 [illuminate/database](https://github.com/illuminate/database)（Eloquent ORM + 查询构建器），并通过 **Origami** 的运行时来执行。

## 特性

- ✅ 独立使用 `illuminate/database`（仅依赖 Composer 包，不引入 `laravel/framework`）
- ✅ 使用 Origami 运行时（Go 实现）解释执行 PHP 代码
- ✅ 完整支持：
  - Capsule 连接管理
  - Schema Builder（表结构定义与修改）
  - 查询构建器（Query Builder）CRUD
  - Eloquent ORM 模型查询
  - 原生 SQL / PDO 预处理语句
  - 数据库事务（含回滚）

## 目录结构

```
examples/illuminate-db/
├── composer.json          # Composer 依赖（仅 illuminate/database）
├── go.mod / go.sum        # Go 模块依赖
├── main.go                # Origami 运行时入口
├── bootstrap.php          # Capsule 引导（数据库连接初始化）
├── examples/              # 示例脚本
│   ├── 01_basic_connection.php    # 基础连接与查询
│   ├── 02_query_builder_crud.php  # 查询构建器 CRUD
│   ├── 03_schema_builder.php      # Schema Builder 表结构
│   ├── 04_eloquent_orm.php        # Eloquent ORM
│   ├── 05_raw_sql.php             # 原生 SQL
│   ├── 06_transactions.php        # 数据库事务
│   ├── 07_eloquent_relationships.php # Eloquent 模型关联
│   ├── 08_eloquent_advanced.php   # Eloquent 高级特性
│   ├── 09_query_builder_advanced.php # 查询构建器高级用法
│   └── 10_multiple_connections.php  # 多数据库连接
└── readme.md              # 本文档
```

## 快速开始

### 1. 安装 Composer 依赖

```bash
cd examples/illuminate-db
composer install
```

### 2. 构建 Origami 运行时

```bash
go build -mod=mod -o illuminate-db .
```

### 3. 运行示例

```bash
# 运行单个示例
./illuminate-db examples/01_basic_connection.php
./illuminate-db examples/02_query_builder_crud.php
./illuminate-db examples/03_schema_builder.php
./illuminate-db examples/04_eloquent_orm.php
./illuminate-db examples/05_raw_sql.php
./illuminate-db examples/06_transactions.php
./illuminate-db examples/07_eloquent_relationships.php
./illuminate-db examples/08_eloquent_advanced.php
./illuminate-db examples/09_query_builder_advanced.php
./illuminate-db examples/10_multiple_connections.php
```

### 4. 使用文件型数据库（可选）

默认使用内存 SQLite（`:memory:`）。如需使用文件型数据库：

```bash
DB_DATABASE=/path/to/database.sqlite ./illuminate-db examples/01_basic_connection.php
```

## 示例说明

| 示例 | 说明 |
|------|------|
| `01_basic_connection.php` | 演示 Capsule 连接初始化、建表、基础增删查改、聚合查询 |
| `02_query_builder_crud.php` | 完整演示查询构建器的 CRUD：插入、批量插入、条件查询、排序、分页、聚合、自增自减 |
| `03_schema_builder.php` | 演示 Schema Builder：创建表、字段类型、外键、索引、修改表结构、删除表 |
| `04_eloquent_orm.php` | 演示 Eloquent ORM：`$model->save()` 插入、`User::where()` 静态调用、`find()`、条件查询、更新/删除影响行数、属性操作 |
| `05_raw_sql.php` | 演示原生 SQL：`select()`、`insert()`、`update()`、`delete()`、联表查询、PDO 预处理 |
| `06_transactions.php` | 演示事务：`beginTransaction()` / `commit()` / `rollBack()`、闭包事务、失败回滚、Capsule 静态事务调用 |
| `07_eloquent_relationships.php` | 演示 Eloquent 模型关联：`hasMany` / `belongsTo` / `hasOne`、关联属性与方法访问、`with()` 预加载、`withCount()`、`has()` / `doesntHave()`、关联创建与删除 |
| `08_eloquent_advanced.php` | 演示 Eloquent 高级特性：`Model::create()`、`firstOrCreate()` / `updateOrCreate()`、访问器/修改器、`$appends`、局部作用域、软删除（SoftDeletes）、模型事件与时间戳 |
| `09_query_builder_advanced.php` | 演示查询构建器高级用法：`join()` / `leftJoin()`、`whereIn()` / `whereBetween()` / `whereNull()`、`pluck()` / `value()`、分组聚合 + `havingRaw()`、自增自减、`updateOrInsert()`、`chunkById()`、`toSql()` |
| `10_multiple_connections.php` | 演示多数据库连接：`addConnection()` 多连接、`connection()` 切换、Schema 跨连接操作、`setDefaultConnection()` / `getDefaultConnection()`、断开与重连 |

## 使用说明

### 连接数据库

```php
use Illuminate\Database\Capsule\Manager as Capsule;

$capsule = new Capsule();
$capsule->addConnection([
    'driver'   => 'sqlite',
    'database' => ':memory:',
]);
$capsule->setAsGlobal();
$capsule->bootEloquent();
```

### 查询构建器

```php
// 查询
$users = Capsule::table('users')->where('age', '>', 18)->get();

// 插入
$id = Capsule::table('users')->insertGetId([
    'name' => 'Alice',
    'email' => 'alice@example.com',
]);

// 更新
Capsule::table('users')->where('id', $id)->update(['name' => 'Alice Updated']);

// 删除
Capsule::table('users')->where('id', $id)->delete();
```

### Eloquent ORM

```php
class User extends Model
{
    protected $table = 'users';
    protected $fillable = ['name', 'email'];
}

// 查询（支持静态调用）
$users = User::query()->get();
$user = User::find(1);
$filtered = User::where('name', 'Alice')->first();
$sorted = User::orderBy('age', 'desc')->get();

// 插入（支持 $model->save()）
$user = new User();
$user->name  = 'Alice';
$user->email = 'alice@example.com';
$user->save();

// 更新 / 删除（返回实际影响行数）
$affected = User::where('id', $user->id)->update(['name' => 'Alice Updated']);
$deleted  = User::where('id', $user->id)->delete();

// 事务（支持 Capsule 静态调用）
Capsule::beginTransaction();
Capsule::table('users')->where('id', $user->id)->delete();
Capsule::rollback(); // 回滚
Capsule::commit();   // 提交
```

## Origami 运行时兼容性说明

当前示例依赖 Origami 运行时对以下 PHP 特性的支持（已实现）：

1. **Eloquent 模型 `save()` 方法**：支持 `$model->save()` 通过模型属性插入新记录（关联数组键在 `__call` 魔术方法参数传递中得以保留）。

2. **`Model::where()` 静态调用**：支持 `User::where()->get()` 等静态调用（`__callStatic` 的实参展开与 `func_get_args()` 计数正确）。

3. **更新/删除影响行数**：`update()` 和 `delete()` 返回实际影响行数（PDO `rowCount()` 正确上报）。

4. **`Capsule::rollback()`**：支持通过 Capsule 静态调用事务控制方法（方法名大小写不敏感匹配）。

> 若运行示例时出现兼容性问题，请确保使用包含上述修复的 Origami 运行时版本。

### 新增示例（07-10）覆盖的额外特性

示例 `07`-`10` 进一步覆盖 Illuminate Database 的常用高级特性：

- **模型关联**：`hasMany` / `belongsTo` / `hasOne` 关系定义与访问、`with()` 预加载、`withCount()`、`has()` / `doesntHave()`、关联创建。
- **Eloquent 高级特性**：`create()` / `firstOrCreate()` / `updateOrCreate()`、访问器/修改器、局部作用域、软删除（SoftDeletes）。
- **查询构建器高级用法**：`join` / `leftJoin`、`whereIn` / `whereBetween` / `whereNull`、`pluck` / `value`、分组聚合 + `havingRaw`、`chunkById` 分块处理、`toSql()`。
- **多数据库连接**：多连接注册、连接切换、跨连接 Schema 与查询、默认连接动态切换、断开与重连。

> 若这些示例因运行时特性暂未实现而报错，可结合 `run_all.sh` 定位问题，并在 Origami 运行时中补充对应的 PHP 特性支持。

## 相关链接

- [illuminate/database 文档](https://laravel.com/docs/10.x/database)
- [Eloquent ORM 文档](https://laravel.com/docs/10.x/eloquent)
- [Origami 项目](https://github.com/php-any/origami)
