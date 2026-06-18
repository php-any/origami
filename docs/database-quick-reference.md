# 数据库模块快速参考

## 连接管理

```zy
use database\DB;
use database\sql\open;

// 连接数据库
$db = open("mysql", "root:password@/database");
database\registerDefaultConnection($db);

// 使用默认连接
$data = DB::model(User::class);

// 使用指定连接
$data = DB::model(User::class, "connection_name");
$data = DB::model(User::class)->connection("connection_name");
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
$users = DB::model(User::class)->get();
$user = DB::model(User::class)->first();

// 条件查询
$users = DB::model(User::class)->where("age > ?", 18)->get();

// 字段选择
$users = DB::model(User::class)->select("id, name, age")->get();

// 排序
$users = DB::model(User::class)->orderBy("age DESC")->get();

// 限制
$users = DB::model(User::class)->limit(10)->get();

// 分页
$users = DB::model(User::class)->offset(20)->limit(10)->get();

// 分组
$stats = DB::model(User::class)->select("age, COUNT(*) as count")->groupBy("age")->get();

// 连接
$results = DB::model(User::class)->join("INNER JOIN profiles p ON users.id = p.user_id")->get();
```

## CRUD 操作

```zy
// 插入
$user = new User();
$user->userName = "张三";
$user->age = 25;
$result = DB::model(User::class)->insert($user);

// 查询
$user = DB::model(User::class)->where("id = ?", 1)->first();
$users = DB::model(User::class)->where("age > ?", 18)->get();

// 更新
$result = DB::model(User::class)->where("id = ?", 1)->update(["age" => 26]);

// 删除
$result = DB::model(User::class)->where("id = ?", 1)->delete();
```

## 原生 SQL

```zy
// 静态查询（无需绑定模型）
$rows = DB::query("SELECT * FROM users WHERE age > ?", 18);
$users = DB::toEntity(User::class, $rows);

// 构建器实例查询
$results = DB::model(User::class)->query("SELECT * FROM users WHERE age > ?", 18);

// 写操作
$result = DB::execute("INSERT INTO users (name, age) VALUES (?, ?)", "张三", 25);
$result = DB::model(User::class)->execute("UPDATE users SET age = ? WHERE id = ?", 26, 1);
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
