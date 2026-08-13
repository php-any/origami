# Schema Migration

根据 Entity 类的 `Database\Annotation` 注解（`#[Table]`、`#[Column]`、`#[Id]`、`#[GeneratedValue]`）自动同步数据库 Schema。

## 用法

```php
use Database\Sql\open;

$db = open("sqlite", "app.db");
$db->ping();
\Database\registerDefaultConnection($db);

$result = \Database\migrate($db, __DIR__ . "/Model/Entity");
// $result->createdCount  — 新建的表数量
// $result->alteredCount  — 新增列数量
// $result->created       — [{ table, class }, ...]
// $result->altered       — [{ table, column, class }, ...]
// $result->tables        — 所有同步的表
```

也支持传入连接名称（需先 `registerDefaultConnection` / `registerConnection`）：

```php
\Database\migrate("default", __DIR__ . "/Model/Entity");
```

## Entity 要求

模型类需使用 `std/database/annotation` 提供的注解：

```php
use Database\Annotation\Column;
use Database\Annotation\GeneratedValue;
use Database\Annotation\Id;
use Database\Annotation\Table;

#[Table("users")]
class UserEntity {
    #[Id]
    #[GeneratedValue("AUTO")]
    #[Column("id", nullable: false)]
    public int $id;

    #[Column("name", nullable: false, length: 100)]
    public string $name;
}
```

## 同步行为

1. 扫描模型目录下所有 `.php` 文件并加载类定义
2. 筛选带 `#[Table]` 注解、且源文件位于该目录下的实体类
3. 表不存在 → `CREATE TABLE IF NOT EXISTS`
4. 表已存在但缺少列 → `ALTER TABLE ... ADD COLUMN`
5. 不会删除已有表或列（安全增量同步）
6. 所有 DDL 在单个事务中执行，失败自动回滚（SQLite 完整支持；MySQL 的 DDL 会隐式提交，见下方说明）

## 数据库方言

自动根据连接驱动选择方言：

| 驱动 | 元数据查询 | 自增主键 |
|------|-----------|---------|
| SQLite | `sqlite_master` / `PRAGMA table_info` | `AUTOINCREMENT` |
| MySQL | `information_schema` | `AUTO_INCREMENT` |

MySQL 类型映射：`int`→`INT`，`float`→`DOUBLE`，`bool`→`TINYINT(1)`，`string`→`VARCHAR(n)`。

> **注意**：MySQL 中 `CREATE TABLE` / `ALTER TABLE` 会触发隐式提交，事务无法保证多步 DDL 的原子性；SQLite 可完整回滚。

## 类型映射

| PHP 类型 | SQL 类型 |
|----------|----------|
| int | INTEGER |
| float | REAL |
| bool | INTEGER |
| string | VARCHAR(length)，length 来自 `#[Column]` |
| ?T | 允许 NULL |

`#[Id]` + `#[GeneratedValue("AUTO")]` 在 SQLite 下生成 `PRIMARY KEY AUTOINCREMENT`，MySQL 下生成 `AUTO_INCREMENT PRIMARY KEY`。

## 与 Entity 的职责划分

| 层 | 职责 |
|----|------|
| `Model/Entity/` | ORM 映射（`#[Table]` / `#[Column]` 等注解） |
| `Database\migrate()` | 读取注解，同步 DDL |

无需手写迁移类，也无需额外的 Migration 注解。

示例见 [`example.php`](example.php)。
