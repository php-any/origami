# 数据库模块快速参考

## 连接管理

```zy
use database\DB;
use database\sql\open;

// 连接数据库
$db = open("mysql", "root:password@/database");
database\registerDefaultConnection($db);

// 使用默认连接
$data = DB<User>();

// 使用指定连接
$data = DB<User>("connection_name");
```

## 注解

```zy
@Table("users")
class User {
    @Id
    public int $id;

    @Column("user_name")
    public string $userName;

    public int $age;
}
```

## 查询构建器

```zy
// 基础查询
$users = DB<User>()->get();
$user = DB<User>()->first();

// 条件查询
$users = DB<User>()->where("age > ?", 18)->get();

// 字段选择
$users = DB<User>()->select("id, name, age")->get();

// 排序
$users = DB<User>()->orderBy("age DESC")->get();

// 限制
$users = DB<User>()->limit(10)->get();

// 分页
$users = DB<User>()->offset(20)->limit(10)->get();

// 分组
$stats = DB<User>()->select("age, COUNT(*) as count")->groupBy("age")->get();

// 连接
$results = DB<User>()->join("INNER JOIN profiles p ON users.id = p.user_id")->get();
```

## CRUD 操作

```zy
// 插入
$user = new User();
$user->userName = "张三";
$user->age = 25;
$result = DB<User>()->insert($user);

// 查询
$user = DB<User>()->where("id = ?", 1)->first();
$users = DB<User>()->where("age > ?", 18)->get();

// 更新
$result = DB<User>()->where("id = ?", 1)->update(["age" => 26]);

// 删除
$result = DB<User>()->where("id = ?", 1)->delete();
```

## 原生 SQL

```zy
// 查询
$results = DB<User>()->query("SELECT * FROM users WHERE age > ?", [18]);

// 执行
$result = DB<User>()->exec("INSERT INTO users (name, age) VALUES (?, ?)", ["张三", 25]);
```

## 连接管理函数

```zy
// 注册连接
database\registerConnection("name", $db);
database\registerDefaultConnection($db);

// 获取连接
$conn = database\getConnection("name");
$default = database\getDefaultConnection();

// 移除连接
database\removeConnection("name");

// 列出连接
$connections = database\listConnections();
```

## 返回值

### 查询结果

- `get()`: 返回数组
- `first()`: 返回单个对象或 null
- `query()`: 返回数组

### 执行结果

```zy
$result = [
    "success" => true,
    "rowsAffected" => 1,
    "lastInsertId" => 123
];
```
