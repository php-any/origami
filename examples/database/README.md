# Database 扩展示例

本目录包含了 Origami 语言中数据库功能的使用示例，展示了如何使用 DB<M>泛型类进行数据库操作。

## 项目结构

```
examples/database/
├── go.mod                    # Go模块定义
├── main.go                   # 独立的Go程序入口
├── database.zy               # 数据库操作演示脚本
└── README.md                 # 说明文档
```

## 功能特性

### 数据库操作功能

- 使用 DB<M>泛型类进行类型安全的数据库操作
- 支持@Table 和@Column 注解自动映射
- 完整的 CRUD 操作（Create、Read、Update、Delete）
- 查询构建器支持
- 原生 SQL 查询支持
- 关联查询和复杂查询

### 支持的操作

- 基本查询：get()、first()、where()
- 条件查询：where()、orderBy()、groupBy()
- 分页查询：limit()、offset()
- 插入操作：insert()
- 更新操作：update()
- 删除操作：delete()
- 原生 SQL：query()、exec()
- 关联查询：join()

## 使用方法

### 1. 进入 Database 例子目录

```bash
cd examples/database
```

### 2. 运行数据库操作示例

```bash
go run main.go
```

## 演示内容

### 1. 基本查询操作

- 获取所有用户
- 条件查询（年龄小于 30 的用户）
- 排序查询（按年龄降序）
- 限制查询（限制 5 个用户）

### 2. 插入操作

- 创建新用户对象
- 插入到数据库
- 批量插入示例

### 3. 更新操作

- 条件更新用户信息
- 批量更新示例

### 4. 删除操作

- 条件删除用户
- 清理数据示例

### 5. 复杂查询

- 分组查询
- 聚合查询（平均年龄）
- 统计查询

### 6. 原生 SQL 查询

- 执行原生 SELECT 查询
- 执行原生 UPDATE 语句
- 复杂 SQL 操作

### 7. 关联查询

- 用户和文章的关联查询
- JOIN 操作示例

### 8. 分页查询

- 分页获取数据
- 偏移量查询

## 技术特点

1. **独立 Go 程序**: 有自己的 go.mod 和 main.go 文件
2. **引入主项目库**: 通过 replace 指令引入主项目的功能
3. **完整功能**: 包含解析器、VM、标准库等完整功能
4. **功能单一**: 直接运行 database.zy 脚本，无需命令行参数
5. **信号处理**: 支持 Ctrl+C 优雅停止程序
6. **数据库支持**: 使用 SQLite 数据库进行演示

## 开发说明

### Go 模块配置

```go
module github.com/php-any/origami/examples/database

go 1.24

require (
	github.com/php-any/origami v0.0.0
	github.com/go-sql-driver/mysql v1.9.3
	github.com/mattn/go-sqlite3 v1.14.17
)

require filippo.io/edwards25519 v1.1.0 // indirect

replace github.com/php-any/origami => ../../
```

### 主程序功能

- 创建解析器和 VM 环境
- 加载所有标准库（std、php、database、system）
- 直接运行 database.zy 脚本
- 简洁的错误处理

### 脚本功能

- 在脚本中设置 SQLite 数据库连接
- 自动创建示例表（users、posts）
- 演示完整的数据库操作流程

### 数据库模型

**用户模型 (User)**:

```php
@Table("users")
class User {
    @Column("name")
    public string $userName;

    public string $email;
    public int $age;
    public int $id;
}
```

**文章模型 (Post)**:

```php
@Table("posts")
class Post {
    public string $title;
    public string $content;
    public int $user_id;
    public int $id;
}
```

## 注意事项

1. 确保安装了 SQLite3 驱动
2. 示例使用 SQLite 数据库，文件名为 example.db
3. 程序会自动创建 users 和 posts 表
4. 演示完成后会清理测试数据
5. 可以根据需要修改数据库连接和表结构
6. 在生产环境中应该添加更多的安全措施和错误处理
7. 使用 `Ctrl+C` 停止程序
