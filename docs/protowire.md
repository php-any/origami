# Protowire - Protobuf 二进制解析标准库

`Protowire` 是一个通用的 Protocol Buffers 二进制格式读写类，完全基于 `google.golang.org/protobuf/encoding/protowire` 底层包，**不依赖 `.proto` 文件或代码生成**。支持通过 `@Field` 注解声明字段位置，实现对象级别的序列化与反序列化。

## Wire Type 常量

| 常量 | 值 | 说明 |
|------|----|------|
| `PROTOWIRE_VARINT` | `0` | Varint 编码 |
| `PROTOWIRE_FIXED64` | `1` | 64 位定长 |
| `PROTOWIRE_LENGTH_DELIMITED` | `2` | 长度前缀 |
| `PROTOWIRE_START_GROUP` | `3` | 分组开始 |
| `PROTOWIRE_END_GROUP` | `4` | 分组结束 |
| `PROTOWIRE_FIXED32` | `5` | 32 位定长 |

## @Field 注解

声明类属性对应的 protobuf 字段编号和 wire type。

```php
use Protowire\Annotation\Field;

#[Field(number: 1, type: PROTOWIRE_VARINT)]
public int $id;

#[Field(number: 2, type: PROTOWIRE_LENGTH_DELIMITED)]
public string $name;
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `number` | `int` | 字段编号（即 protobuf field number） |
| `type` | `int` | Wire type 常量 |

## 静态方法

### Protowire::parse — 反序列化

解析 protobuf 二进制数据。

```php
// 原始字段解析
Protowire::parse(string $data, array $options = []): array

// 注解驱动：反序列化到指定类
Protowire::parse(string $data, string $className, array $options = []): object
```

- **只传 data**：返回原始字段数组（每个字段是一个关联数组 `[number, wire_type, value]`）
- **传 className**：读取类的 `@Field` 注解，反序列化为对象实例

**$options 配置项：**

| 键 | 类型 | 说明 |
|----|------|------|
| `message_fields` | `array<int, bool>` | 递归解析嵌套消息（仅原始解析模式） |
| `packed_fields` | `array<int, bool>` | 指定 packed 字段 |
| `packed_element_type` | `array<int, int>` | packed 元素的 wire type |
| `max_depth` | `int` | 最大递归深度，默认 64 |

### Protowire::serialize — 序列化

将对象序列化为 protobuf 二进制数据。

```php
Protowire::serialize(object $instance): string
```

读取实例类的 `@Field` 注解，按字段编号和 wire type 编码为二进制。

### Protowire::encodeVarint / encodeTag / encodeBytes / encodeFixed32 / encodeFixed64

编码辅助方法，用于手工构造 protobuf 二进制数据。

```php
Protowire::encodeVarint(int $value): string
Protowire::encodeTag(int $number, int $wireType): string
Protowire::encodeBytes(string $value): string
Protowire::encodeFixed32(int $value): string
Protowire::encodeFixed64(int $value): string
```

---

## 使用示例

### 注解驱动：定义 + 序列化/反序列化

```php
<?php

use Protowire\Annotation\Field;

class User {
    #[Field(number: 1, type: PROTOWIRE_VARINT)]
    public int $id;

    #[Field(number: 2, type: PROTOWIRE_LENGTH_DELIMITED)]
    public string $name;
}

// 序列化
$user = new User();
$user->id = 42;
$user->name = 'Alice';
$data = Protowire::serialize($user);

// 反序列化
$user2 = Protowire::parse($data, 'User');
// $user2->id === 42, $user2->name === 'Alice'
```

### 原始字段解析（不带注解）

```php
$data = Protowire::encodeTag(1, PROTOWIRE_VARINT) . Protowire::encodeVarint(42);
$fields = Protowire::parse($data);
// [['number' => 1, 'wire_type' => 0, 'value' => 42]]
```

### Packed 字段

```php
$packed = Protowire::encodeVarint(1) . Protowire::encodeVarint(2) . Protowire::encodeVarint(3);
$data = Protowire::encodeTag(7, PROTOWIRE_LENGTH_DELIMITED) . Protowire::encodeBytes($packed);

$fields = Protowire::parse($data, [
    'packed_fields' => [7 => true],
    'packed_element_type' => [7 => PROTOWIRE_VARINT],
]);
// [['number' => 7, 'value' => [1, 2, 3]]]
```

### 嵌套消息

```php
$inner = Protowire::encodeTag(1, PROTOWIRE_VARINT) . Protowire::encodeVarint(7);
$data = Protowire::encodeTag(5, PROTOWIRE_LENGTH_DELIMITED) . Protowire::encodeBytes($inner);

$fields = Protowire::parse($data, ['message_fields' => [5 => true]]);
// field 5 的 value 被递归解析为字段数组
```

### Group 分组

```php
$data = Protowire::encodeTag(10, PROTOWIRE_START_GROUP)
      . Protowire::encodeTag(1, PROTOWIRE_VARINT) . Protowire::encodeVarint(123)
      . Protowire::encodeTag(2, PROTOWIRE_FIXED32) . Protowire::encodeFixed32(456)
      . Protowire::encodeTag(10, PROTOWIRE_END_GROUP);

$fields = Protowire::parse($data);
// field 10 的 value 自动解析为内部字段数组
```

### 递归深度限制

```php
$fields = Protowire::parse($data, [
    'message_fields' => [1 => true],
    'max_depth' => 10,
]);
```

---

## 错误处理

解析错误通过抛出 `\Exception` 报告，可使用 `try/catch` 捕获：

```php
try {
    $fields = Protowire::parse($invalidData);
} catch (\Exception $e) {
    echo $e->getMessage();
}
```

常见错误信息包含字段编号和 wire type，便于定位问题：

```
Protowire::parse: field 1 (wire type 0): invalid varint encoding
Protowire::parse: class "XXX" not found
Protowire::parse: nested message: maximum recursion depth exceeded
Protowire::parse: group 10: mismatched end group tag number 99
Protowire::parse: packed field 7: PackedElementType not configured
```
