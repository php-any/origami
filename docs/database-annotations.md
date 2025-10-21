# 数据库注解文档

本文档详细介绍数据库模块中使用的各种注解。

## 目录

- [@Table 注解](#table-注解)
- [@Column 注解](#column-注解)
- [@Id 注解](#id-注解)
- [@GeneratedValue 注解](#generatedvalue-注解)
- [注解组合使用](#注解组合使用)
- [最佳实践](#最佳实践)

## @Table 注解

### 功能

指定数据库表名，用于类与数据库表的映射。

### 语法

```zy
@Table("table_name")
class ClassName {
    // 类定义
}
```

### 示例

```zy
namespace App\Models;

@Table("user_profiles")
class UserProfile {
    public int $id;
    public string $name;
}

@Table("order_items")
class OrderItem {
    public int $id;
    public int $orderId;
    public string $productName;
}
```

### 规则

- 如果没有 `@Table` 注解，将使用类名作为表名
- 表名区分大小写
- 支持数据库中的实际表名

## @Column 注解

### 功能

映射类属性到数据库列名，支持属性名与列名不一致的情况。

### 语法

```zy
@Column("column_name")
public type $propertyName;
```

### 示例

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

    @Column("is_active")
    public bool $isActive;
}
```

### 规则

- 如果没有 `@Column` 注解，将使用属性名作为列名
- 列名区分大小写
- 支持数据库中的实际列名

## @Id 注解

### 功能

标识主键字段，用于数据库操作中的主键识别。

### 语法

```zy
@Id
public type $propertyName;
```

### 示例

```zy
@Table("users")
class User {
    @Id
    public int $id;

    public string $name;
    public string $email;
}
```

### 规则

- 每个类只能有一个 `@Id` 注解
- 主键字段在插入时会被特殊处理
- 支持复合主键（多个字段都有 `@Id` 注解）

## @GeneratedValue 注解

### 功能

标识自动生成的字段，如自增主键、时间戳等。

### 语法

```zy
@GeneratedValue
public type $propertyName;
```

### 示例

```zy
@Table("users")
class User {
    @Id
    @GeneratedValue
    public int $id;

    public string $name;

    @GeneratedValue
    public string $createdAt;
}
```

### 规则

- 通常与 `@Id` 注解一起使用
- 自动生成的字段在插入时会被忽略
- 支持数据库的自增字段、默认值等

## 注解组合使用

### 完整示例

```zy
namespace App\Models;

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

    @Column("avatar_url")
    public string $avatarUrl;

    @Column("created_at")
    @GeneratedValue
    public string $createdAt;

    @Column("updated_at")
    @GeneratedValue
    public string $updatedAt;
}
```

### 复杂映射示例

```zy
@Table("order_items")
class OrderItem {
    @Id
    public int $orderId;

    @Id
    public int $productId;

    @Column("item_quantity")
    public int $quantity;

    @Column("unit_price")
    public float $unitPrice;

    @Column("total_amount")
    public float $totalAmount;

    @Column("created_at")
    @GeneratedValue
    public string $createdAt;
}
```

## 最佳实践

### 1. 命名规范

```zy
// ✅ 推荐：使用清晰的表名和列名
@Table("user_accounts")
class UserAccount {
    @Id
    @GeneratedValue
    public int $id;

    @Column("user_id")
    public int $userId;

    @Column("account_balance")
    public float $balance;

    @Column("created_at")
    public string $createdAt;
}

// ❌ 避免：使用模糊的命名
@Table("t1")
class T1 {
    public int $a;
    public string $b;
}
```

### 2. 注解顺序

```zy
// ✅ 推荐：注解顺序清晰
@Table("users")
class User {
    @Id
    @GeneratedValue
    public int $id;

    @Column("user_name")
    public string $userName;

    public int $age;
}
```

### 3. 类型安全

```zy
// ✅ 推荐：使用明确的类型
@Table("products")
class Product {
    @Id
    @GeneratedValue
    public int $id;

    public string $name;
    public float $price;
    public bool $isActive;
    public string $description;
}

// ❌ 避免：使用模糊的类型
class Product {
    public $id;
    public $name;
    public $price;
}
```

### 4. 数据库设计对应

```zy
// 数据库表结构
// CREATE TABLE user_sessions (
//     id INT AUTO_INCREMENT PRIMARY KEY,
//     user_id INT NOT NULL,
//     session_token VARCHAR(255) NOT NULL,
//     expires_at TIMESTAMP NOT NULL,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );

@Table("user_sessions")
class UserSession {
    @Id
    @GeneratedValue
    public int $id;

    @Column("user_id")
    public int $userId;

    @Column("session_token")
    public string $sessionToken;

    @Column("expires_at")
    public string $expiresAt;

    @Column("created_at")
    @GeneratedValue
    public string $createdAt;
}
```

### 5. 继承和组合

```zy
// 基础实体类
@Table("base_entities")
abstract class BaseEntity {
    @Id
    @GeneratedValue
    public int $id;

    @Column("created_at")
    @GeneratedValue
    public string $createdAt;

    @Column("updated_at")
    public string $updatedAt;
}

// 具体实体类
@Table("users")
class User extends BaseEntity {
    @Column("user_name")
    public string $userName;

    public string $email;
}

@Table("products")
class Product extends BaseEntity {
    public string $name;
    public float $price;
    public string $description;
}
```

## 常见问题

### Q: 如何处理数据库字段名与类属性名不一致？

A: 使用 `@Column` 注解进行映射：

```zy
@Table("users")
class User {
    @Column("user_name")  // 数据库列名
    public string $userName;  // 类属性名
}
```

### Q: 如何处理复合主键？

A: 在多个字段上使用 `@Id` 注解：

```zy
@Table("order_items")
class OrderItem {
    @Id
    public int $orderId;

    @Id
    public int $productId;

    public int $quantity;
}
```

### Q: 如何处理自动生成的时间戳？

A: 使用 `@GeneratedValue` 注解：

```zy
@Table("users")
class User {
    @Id
    @GeneratedValue
    public int $id;

    @Column("created_at")
    @GeneratedValue
    public string $createdAt;
}
```

### Q: 注解是否区分大小写？

A: 是的，注解和参数都区分大小写：

```zy
// ✅ 正确
@Table("users")
@Column("user_name")

// ❌ 错误
@table("users")
@column("user_name")
```

## 总结

数据库注解提供了强大的 ORM 映射功能：

- ✅ `@Table`: 指定数据库表名
- ✅ `@Column`: 映射属性到列名
- ✅ `@Id`: 标识主键字段
- ✅ `@GeneratedValue`: 标识自动生成字段
- ✅ 支持复杂映射和继承
- ✅ 类型安全和命名规范

通过合理使用这些注解，可以构建清晰、可维护的数据库模型。
