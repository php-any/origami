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
│   └── 06_transactions.php        # 数据库事务
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
| `04_eloquent_orm.php` | 演示 Eloquent ORM：模型定义、`query()` 查询、`find()`、条件查询、更新、删除、属性操作 |
| `05_raw_sql.php` | 演示原生 SQL：`select()`、`insert()`、`update()`、`delete()`、联表查询、PDO 预处理 |
| `06_transactions.php` | 演示事务：`beginTransaction()` / `commit()` / `rollBack()`、闭包事务、失败回滚 |

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

// 查询（推荐使用 query() 方法）
$users = User::query()->get();
$user = User::query()->find(1);
$filtered = User::query()->where('name', 'Alice')->first();
```

## 已知限制（Origami 运行时兼容性）

当前 Origami 运行时对以下 PHP 特性支持仍在完善中：

1. **Eloquent 模型 `save()` 方法**：`$model->save()` 在插入新记录时可能因数组参数在魔术方法 `__call` 中的传递问题而失败。建议通过查询构建器进行插入操作。

2. **`Model::where()` 静态调用**：`User::where()->get()` 可能无法正确传递查询条件。建议使用 `User::query()->where()->get()`。

3. **更新/删除影响行数**：`update()` 和 `delete()` 返回的影响行数为 `-1`，但操作本身会正常执行。

4. **`Capsule::rollback()`**：通过 Capsule 静态调用 `rollback()` 可能失败。建议使用 `Capsule::connection()->rollBack()` 或 `Capsule::connection()->beginTransaction()`。

## 相关链接

- [illuminate/database 文档](https://laravel.com/docs/10.x/database)
- [Eloquent ORM 文档](https://laravel.com/docs/10.x/eloquent)
- [Origami 项目](https://github.com/php-any/origami)
