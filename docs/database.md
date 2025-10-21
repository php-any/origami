# 数据库模块文档

本文档介绍 Origami 语言中的数据库模块，包括连接管理、查询构建器、CRUD 操作、注解支持等功能。

## 目录

- [快速开始](#快速开始)
- [连接管理](#连接管理)
- [注解支持](#注解支持)
- [查询构建器](#查询构建器)
- [CRUD 操作](#crud-操作)
- [原生 SQL 支持](#原生-sql-支持)
- [高级功能](#高级功能)
- [最佳实践](#最佳实践)

## 快速开始

### 基本使用

```zy
use database\DB;
use database\sql\open;

// 连接数据库
$db = open("mysql", "root:password@/database_name");
$db->ping();

// 注册为默认连接
database\registerDefaultConnection($db);

// 使用查询构建器
$users = DB<User>()->where("age > ?", 18)->get();
```

### 定义模型类

```zy
namespace App;

use database\annotation\Table;
use database\annotation\Column;

@Table("users")
class User {
    public int $id;

    @Column("name")
    public string $userName;

    public int $age;

    public float $coin;

    @Column("create_at")
    public string $createAt;
}
```

## 连接管理

### 连接注册

数据库模块支持多连接管理，可以注册多个数据库连接：

```zy
use database\sql\open;

// 注册默认连接
$defaultDb = open("mysql", "root:password@/main_db");
database\registerDefaultConnection($defaultDb);

// 注册命名连接
$userDb = open("mysql", "user:pass@/user_db");
database\registerConnection("users", $userDb);

$logDb = open("mysql", "log:pass@/log_db");
database\registerConnection("logs", $logDb);
```

### 连接使用

```zy
// 使用默认连接
$data = DB<User>();

// 使用指定连接
$data = DB<User>("users");
$logData = DB<Log>("logs");
```

### 连接管理函数

```zy
// 获取连接
$conn = database\getConnection("users");
$defaultConn = database\getDefaultConnection();

// 移除连接
database\removeConnection("users");

// 列出所有连接
$connections = database\listConnections();
```

## 注解支持

### @Table 注解

用于指定数据库表名：

```zy
@Table("user_profiles")
class UserProfile {
    // 类定义
}

// 如果没有 @Table 注解，将使用类名作为表名
class Product {
    // 对应 product 表
}
```

### @Column 注解

用于映射类属性到数据库列名：

```zy
@Table("users")
class User {
    public int $id;

    @Column("user_name")
    public string $userName;

    @Column("email_address")
    public string $email;

    @Column("created_at")
    public string $createAt;
}
```

### @Id 注解

标识主键字段：

```zy
class User {
    @Id
    public int $id;

    public string $name;
}
```

### @GeneratedValue 注解

标识自动生成的字段：

```zy
class User {
    @Id
    @GeneratedValue
    public int $id;

    public string $name;
}
```

## 查询构建器

### 基础查询

```zy
// 查询所有记录
$users = DB<User>()->get();

// 查询单条记录
$user = DB<User>()->first();

// 条件查询
$users = DB<User>()->where("age > ?", 18)->get();
```

### 字段选择

```zy
// 选择特定字段
$users = DB<User>()
    ->select("id, name, age")
    ->get();

// 使用别名
$users = DB<User>()
    ->select("id as user_id, name as user_name")
    ->get();
```

### 条件查询

```zy
// 单个条件
$users = DB<User>()->where("age > ?", 18)->get();

// 多个条件
$users = DB<User>()
    ->where("age > ?", 18)
    ->where("status = ?", "active")
    ->get();

// 复杂条件
$users = DB<User>()
    ->where("(age > ? OR age < ?) AND status = ?", [18, 65, "active"])
    ->get();
```

### 排序

```zy
// 升序排序
$users = DB<User>()->orderBy("age")->get();

// 降序排序
$users = DB<User>()->orderBy("age DESC")->get();

// 多字段排序
$users = DB<User>()
    ->orderBy("status ASC, age DESC")
    ->get();
```

### 分组

```zy
// 按年龄分组统计
$stats = DB<User>()
    ->select("age, COUNT(*) as count")
    ->groupBy("age")
    ->get();
```

### 限制和偏移

```zy
// 限制记录数
$users = DB<User>()->limit(10)->get();

// 分页查询
$users = DB<User>()
    ->offset(20)
    ->limit(10)
    ->get();
```

### 连接查询

```zy
// 内连接
$results = DB<User>()
    ->join("INNER JOIN user_profiles up ON users.id = up.user_id")
    ->select("users.*, up.bio, up.avatar")
    ->get();

// 左连接
$results = DB<User>()
    ->join("LEFT JOIN orders o ON users.id = o.user_id")
    ->where("o.status = ?", "pending")
    ->get();
```

## CRUD 操作

### 插入 (Create)

```zy
// 插入单个记录
$user = new User();
$user->userName = "张三";
$user->age = 25;
$user->coin = 100.0;

$result = DB<User>()->insert($user);
echo "插入成功，ID: " . $result->lastInsertId;

// 插入数组数据
$userData = [
    "userName" => "李四",
    "age" => 30,
    "coin" => 200.0
];

$result = DB<User>()->insert($userData);
```

### 查询 (Read)

```zy
// 查询所有记录
$users = DB<User>()->get();

// 查询单条记录
$user = DB<User>()->where("id = ?", 1)->first();

// 条件查询
$activeUsers = DB<User>()
    ->where("status = ?", "active")
    ->orderBy("created_at DESC")
    ->get();

// 统计查询
$count = DB<User>()->where("age > ?", 18)->get();
echo "成年用户数量: " . count($count);
```

### 更新 (Update)

```zy
// 更新记录
$updateData = [
    "coin" => 500.0,
    "age" => 26
];

$result = DB<User>()
    ->where("id = ?", 1)
    ->update($updateData);

echo "更新了 " . $result->rowsAffected . " 条记录";

// 更新类实例
$user = new User();
$user->coin = 1000.0;
$user->age = 30;

$result = DB<User>()
    ->where("name = ?", "张三")
    ->update($user);
```

### 删除 (Delete)

```zy
// 删除记录
$result = DB<User>()
    ->where("id = ?", 1)
    ->delete();

echo "删除了 " . $result->rowsAffected . " 条记录";

// 批量删除
$result = DB<User>()
    ->where("status = ?", "inactive")
    ->delete();
```

## 原生 SQL 支持

### 查询操作 (query)

```zy
// 简单查询
$results = DB<User>()->query("SELECT * FROM users WHERE age > ?", [18]);

// 复杂查询
$results = DB<User>()->query("
    SELECT u.name, p.title
    FROM users u
    JOIN posts p ON u.id = p.user_id
    WHERE u.age > ?
", [18]);

// 统计查询
$stats = DB<User>()->query("
    SELECT COUNT(*) as total, AVG(age) as avg_age
    FROM users
");
```

### 执行操作 (exec)

```zy
// 插入操作
$result = DB<User>()->exec("
    INSERT INTO users (name, age, coin)
    VALUES (?, ?, ?)
", ["新用户", 25, 100.0]);

echo "插入成功，ID: " . $result->lastInsertId;

// 更新操作
$result = DB<User>()->exec("
    UPDATE users
    SET age = ?
    WHERE name = ?
", [30, "新用户"]);

echo "更新了 " . $result->rowsAffected . " 条记录";

// 删除操作
$result = DB<User>()->exec("
    DELETE FROM users
    WHERE age < ?
", [18]);

// DDL 操作
$result = DB<User>()->exec("
    CREATE TABLE IF NOT EXISTS logs (
        id INT AUTO_INCREMENT PRIMARY KEY,
        message TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
");
```

## 高级功能

### 事务支持

```zy
// 开始事务
$db = database\getDefaultConnection();
$db->begin();

try {
    // 执行多个操作
    DB<User>()->insert($user1);
    DB<User>()->insert($user2);

    // 提交事务
    $db->commit();
    echo "事务提交成功";
} catch (Exception $e) {
    // 回滚事务
    $db->rollback();
    echo "事务回滚: " . $e->getMessage();
}
```

### 批量操作

```zy
// 批量插入
$users = [
    ["name" => "用户1", "age" => 25],
    ["name" => "用户2", "age" => 30],
    ["name" => "用户3", "age" => 28]
];

foreach ($users as $userData) {
    $user = new User();
    $user->userName = $userData["name"];
    $user->age = $userData["age"];

    DB<User>()->insert($user);
}
```

### 复杂查询示例

```zy
// 多表连接查询
$results = DB<User>()->query("
    SELECT
        u.id,
        u.name,
        u.age,
        COUNT(o.id) as order_count,
        SUM(o.amount) as total_amount
    FROM users u
    LEFT JOIN orders o ON u.id = o.user_id
    WHERE u.status = ?
    GROUP BY u.id, u.name, u.age
    HAVING COUNT(o.id) > 0
    ORDER BY total_amount DESC
    LIMIT 10
", ["active"]);

// 子查询
$results = DB<User>()->query("
    SELECT * FROM users
    WHERE id IN (
        SELECT user_id FROM orders
        WHERE amount > ?
    )
", [1000]);
```

## 最佳实践

### 1. 使用注解进行映射

```zy
@Table("user_profiles")
class UserProfile {
    @Id
    @GeneratedValue
    public int $id;

    @Column("user_id")
    public int $userId;

    @Column("full_name")
    public string $fullName;

    @Column("bio")
    public string $bio;

    @Column("created_at")
    public string $createdAt;
}
```

### 2. 合理使用连接

```zy
// 为不同模块使用不同连接
$userData = DB<User>("users");      // 用户数据库
$logData = DB<Log>("logs");         // 日志数据库
$cacheData = DB<Cache>("cache");     // 缓存数据库
```

### 3. 使用参数绑定防止 SQL 注入

```zy
// ✅ 正确：使用参数绑定
$users = DB<User>()->where("name = ?", $userName)->get();

// ❌ 错误：直接拼接 SQL
$users = DB<User>()->query("SELECT * FROM users WHERE name = '" . $userName . "'");
```

### 4. 合理使用索引

```zy
// 为常用查询字段创建索引
DB<User>()->exec("CREATE INDEX idx_users_age ON users(age)");
DB<User>()->exec("CREATE INDEX idx_users_status ON users(status)");
```

### 5. 错误处理

```zy
try {
    $user = DB<User>()->where("id = ?", $userId)->first();
    if ($user === null) {
        throw new Exception("用户不存在");
    }

    $result = DB<User>()
        ->where("id = ?", $userId)
        ->update(["last_login" => date("Y-m-d H:i:s")]);

} catch (Exception $e) {
    echo "数据库操作失败: " . $e->getMessage();
}
```

### 6. 性能优化

```zy
// 使用 select 限制字段
$users = DB<User>()
    ->select("id, name, email")
    ->where("status = ?", "active")
    ->get();

// 使用 limit 限制结果集
$recentUsers = DB<User>()
    ->orderBy("created_at DESC")
    ->limit(100)
    ->get();

// 使用分页
$page = 1;
$pageSize = 20;
$users = DB<User>()
    ->offset(($page - 1) * $pageSize)
    ->limit($pageSize)
    ->get();
```

## 常见问题

### Q: 如何处理数据库连接失败？

A: 使用 try-catch 块捕获异常：

```zy
try {
    $db = open("mysql", "root:password@/database");
    $db->ping();
    database\registerDefaultConnection($db);
} catch (Exception $e) {
    echo "数据库连接失败: " . $e->getMessage();
}
```

### Q: 如何调试 SQL 查询？

A: 可以启用查询日志或使用原生 SQL 进行调试：

```zy
// 使用原生 SQL 查看实际执行的查询
$results = DB<User>()->query("EXPLAIN SELECT * FROM users WHERE age > ?", [18]);
```

### Q: 如何处理大量数据？

A: 使用分页和批量处理：

```zy
$page = 1;
$pageSize = 1000;

do {
    $users = DB<User>()
        ->offset(($page - 1) * $pageSize)
        ->limit($pageSize)
        ->get();

    // 处理数据
    foreach ($users as $user) {
        // 处理逻辑
    }

    $page++;
} while (count($users) === $pageSize);
```

## 总结

数据库模块提供了完整的 ORM 功能，包括：

- ✅ 多连接管理
- ✅ 注解支持
- ✅ 查询构建器
- ✅ CRUD 操作
- ✅ 原生 SQL 支持
- ✅ 事务支持
- ✅ 性能优化

通过合理使用这些功能，可以构建高效、安全的数据库应用程序。
