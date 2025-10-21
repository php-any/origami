# 封装 DB 操作类

代码映射成脚本域的操作如下

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

// 使用默认连接（自动表名）
$data = DB<User>();
$data->where("name = ?", "测试用户")->first();

// 使用指定连接（自动表名）
$data = DB<User>("slave");
$data->where("name = ?", "测试用户")->first();

// 显式指定表名（覆盖注解）
$data = DB<User>();
$data->table("custom_table")->where("id = ?", 1)->first();
```

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
$data = DB<User>();
$data->where("name = ?", "测试用户")->first();
$data->insert($user);
$data->where("id = ?", 1)->update($updateUser);
$data->where("id = ?", 1)->delete();
```

### 显式指定表名

```php
// 显式调用 table() 方法会覆盖 @Table 注解
$data = DB<User>();
$data->table("custom_table")->where("id = ?", 1)->first();
```

## 原生 SQL 支持

数据库模块支持执行原生 SQL 语句，提供最大的灵活性：

### 查询操作 (query)

```php
// 执行 SELECT 查询
$results = $data->query("SELECT * FROM users WHERE age > ?", [25]);

// 复杂查询
$results = $data->query("SELECT u.name, p.title FROM users u JOIN posts p ON u.id = p.user_id WHERE u.age > ?", [18]);

// 统计查询
$stats = $data->query("SELECT COUNT(*) as total, AVG(age) as avg_age FROM users");
```

### 执行操作 (exec)

```php
// 执行 INSERT 语句
$result = $data->exec("INSERT INTO users (name, age, coin) VALUES (?, ?, ?)", ["新用户", 25, 100.0]);

// 执行 UPDATE 语句
$result = $data->exec("UPDATE users SET age = ? WHERE name = ?", [30, "新用户"]);

// 执行 DELETE 语句
$result = $data->exec("DELETE FROM users WHERE age < ?", [18]);

// 执行 DDL 语句
$result = $data->exec("CREATE TABLE IF NOT EXISTS logs (id INT AUTO_INCREMENT PRIMARY KEY, message TEXT)");
```

### 返回值说明

- **query()**: 返回数组，包含查询结果对象
- **exec()**: 返回对象，包含：
  - `rowsAffected`: 影响的行数
  - `lastInsertId`: 最后插入的 ID（仅 INSERT 语句）
  - `success`: 操作是否成功

## 连接管理

数据库连接管理器支持多连接管理：

- `DB<User>()` - 使用默认连接
- `DB<User>("连接名称")` - 使用指定名称的连接

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
$user = DB<User>()->where("id = ?", 1)->first();
$userFromSlave = DB<User>("slave")->where("id = ?", 1)->first();
```
