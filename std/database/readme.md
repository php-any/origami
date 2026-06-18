# 封装 DB 操作类

## 入口一览

| 场景 | 推荐写法 | 说明 |
|------|---------|------|
| ORM 链式查询 | `DB::model(User::class)` | 通过类名绑定模型并自动映射实体 |
| ORM 链式查询（动态类名） | `DB::model($class)` | 运行时类名字符串 |
| 原生 SELECT | `DB::query($sql, ...$args)` | 返回行对象，配合 `toEntity` 映射 |
| 原生写操作 | `DB::execute($sql, ...$args)` | 返回 `{ rowsAffected, lastInsertId, success }` |
| 插入实体 | `DB::insert($entity)` | 从实体推断模型类并插入 |
| 指定连接 | `->connection($name)` / `DB::connection($name)` | 链式或静态切换命名连接 |

```
namespace App;

use database\DB;
use database\annotation\Table;
use database\annotation\Column;

@Table("users")
class User {
    @Column("name")
    public string $userName;

    public int $age;

    public float $coin;
}

// 原生 SQL：查询 + 映射分两步
$rows = DB::query("SELECT * FROM users WHERE age > ?", 18);
$users = DB::toEntity(User::class, $rows);

DB::insert($user);
DB::execute("UPDATE users SET age = ? WHERE id = ?", 30, 1);

// ORM 链式（类名入口，自动映射实体）
DB::model(User::class)->where("name = ?", "测试用户")->first();
DB::model(User::class)->select("id, name")->get();

// 动态类名
DB::model(User::class)->where("name = ?", "测试用户")->first();
```

### 切换连接

```php
// 链式切换（推荐，可与其他条件任意组合）
DB::model(User::class)->connection("slave")->where("id = ?", 1)->first();

// 原生 SQL 走从库
$rows = DB::connection("slave")->query("SELECT * FROM users WHERE age > ?", 18);

// model / insert 第二参数也可指定连接
DB::model(User::class, "slave")->get();
DB::insert($user, "master");
```

## 原生 SQL API

`Database\DB` 提供统一的原生 SQL 入口，静态方法与实例方法语义一致：

| 操作 | 静态方法 | 实例方法 | 返回值 |
|------|---------|---------|--------|
| SELECT 查询 | `DB::query($sql, ...$args)` | `->query($sql, ...$args)` | 行对象数组 |
| 写操作 / DDL | `DB::execute($sql, ...$args)` | `->execute($sql, ...$args)` | `{ rowsAffected, lastInsertId, success }` |

参数支持可变参数（与 `where` 一致）或数组传参：`DB::query($sql, [$a, $b])`。

底层连接封装（`database\sql`）仍使用 Go `database/sql` 命名：`Query` / `Exec`。

## 自动表名功能

数据库模块支持通过 `@Table` 注解自动获取表名，无需显式调用 `table()` 方法：

### 使用 @Table 注解

```php
@Table("users")
class User {
    @Column("name")
    public string $userName;

    public int $age;
}

// 自动使用 @Table 注解中的表名
$data = DB::model(User::class);
$data->where("name = ?", "测试用户")->first();
$data->insert($user);
$data->where("id = ?", 1)->update($updateUser);
$data->where("id = ?", 1)->delete();
```

### 显式指定表名

```php
// 显式调用 table() 方法会覆盖 @Table 注解
$data = DB::model(User::class);
$data->table("custom_table")->where("id = ?", 1)->first();
```

## 连接管理

数据库连接管理器支持多连接管理：

- `DB::model(User::class)` - 使用默认连接
- `DB::model(User::class, "连接名称")` - 使用指定名称的连接

连接需要在 Go 代码中预先注册到连接管理器中。

## 脚本域连接管理函数

数据库模块提供了以下脚本域函数来管理数据库连接：

### 注册连接

```php
// 注册默认连接
database\registerDefaultConnection($db);

// 注册命名连接
database\registerConnection("slave", $db);
database\registerConnection("master", $db);
```

### 获取连接

```php
// 获取默认连接
$defaultConn = database\getDefaultConnection();

// 获取指定连接
$slaveConn = database\getConnection("slave");
```

### 连接管理

```php
// 列出所有连接
$connections = database\listConnections();

// 移除连接
database\removeConnection("slave");
```

### 完整使用示例

```php
use database\sql\open;

// 打开数据库连接
$db = open("mysql", "root:root@/temp");
$db->ping();

// 注册到连接管理器
database\registerDefaultConnection($db);
database\registerConnection("slave", $db);

// 使用连接查询
$user = DB::model(User::class)->where("id = ?", 1)->first();
$userFromSlave = DB::model(User::class, "slave")->where("id = ?", 1)->first();
```
