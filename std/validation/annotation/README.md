# Validation 注解模块

`Validation\Annotation` 提供请求 DTO / 参数上的**校验约束元数据**，设计参考 Jakarta Bean Validation（`@NotBlank`、`@Email`、`@Min` 等）。

> **当前范围**：注解定义与注册、运行时校验器 `Validation\Validator`、HTTP DTO 绑定后自动校验（422）。  
> **绑定层**：`Request::bind()` 负责 query/form/body → DTO 类型转换；**约束层**：属性上的 `#[Size]`、`#[Pattern]` 等由校验器执行。

## 注册方式

通过 `std/validation.Load(vm)` 加载，所有约束在 `specs.go` 中声明，由通用 `ConstraintClass` 统一注册。

新增约束：在 `specs.go` 的 `constraintSpecs` 追加一项即可，无需为每个注解单独写类文件。

## 注解目标

所有约束均标注在：

- **属性**（`TARGET_PROPERTY`）— DTO 字段
- **参数**（`TARGET_PARAMETER`）— 方法形参

除 `Name` 外，其余约束支持**重复使用**（`IS_REPEATABLE`），可在同一字段上叠加多个规则。

## 类型与约束的分工

| 层次 | 机制 | 说明 |
|------|------|------|
| **类型** | PHP 属性类型 | `public string $email`、`public int $age` 等，由绑定/反序列化层（如 `Request::bind()` → `JsonSerializer`）按类型转换 |
| **约束** | 本模块注解 | `NotBlank`、`Email`、`Min` 等，描述值的业务规则，由未来的校验器读取并执行 |

一般**不需要**单独的 `@Type` 注解；类型以属性声明为准，约束注解只表达格式、范围、非空等规则。

## 内置约束

### Name

字段别名 / 展示名（绑定与错误信息均可引用）。

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `value` | string | `""` | 外部字段名或展示名 |

不可重复。

```php
use Validation\Annotation\Name;

class LoginRequest {
    #[Name(value: "username")]
    public string $username;
}
```

### NotBlank

值不能为 `null`、空字符串或空数组。

| 参数 | 类型 | 默认值 |
|------|------|--------|
| `message` | string | `""` |

### Email

字符串须为合法邮箱格式（运行时按正则校验）。

| 参数 | 类型 | 默认值 |
|------|------|--------|
| `message` | string | `""` |

### Min / Max

数值下限 / 上限，适用于 `int`、`float` 等数值字段。值为空时通常跳过（由 `NotBlank` 负责必填）。

| 参数 | 类型 | 默认值 |
|------|------|--------|
| `value` | int | `0` |
| `message` | string | `""` |

### Size

字符串长度或数组长度的范围。`min` / `max` 为 `0` 表示不限制该侧。

| 参数 | 类型 | 默认值 |
|------|------|--------|
| `min` | int | `0` |
| `max` | int | `0` |
| `message` | string | `""` |

### Pattern

字符串须匹配给定正则。

| 参数 | 类型 | 默认值 |
|------|------|--------|
| `regexp` | string | `""` |
| `message` | string | `""` |

## 使用示例

```php
use Validation\Annotation\Name;
use Validation\Annotation\NotBlank;
use Validation\Annotation\Size;
use Validation\Annotation\Email;
use Validation\Annotation\Min;
use Validation\Annotation\Max;
use Validation\Annotation\Pattern;

class CreateUserRequest {
    #[Name(value: "user_name")]
    #[NotBlank(message: "用户名不能为空")]
    #[Size(min: 2, max: 32)]
    public string $username;

    #[NotBlank]
    #[Email]
    public string $email;

    #[Min(value: 18)]
    #[Max(value: 120)]
    public int $age;

    #[Pattern(regexp: "^1[3-9]\\d{9}$", message: "手机号格式不正确")]
    public string $phone;
}
```

## 实现结构

```
annotation/
├── constraint.go        # 通用 ConstraintClass 与构造函数
├── specs.go             # 各约束的参数与元数据定义
├── load.go              # vm.AddClass 注册
├── compile_bootstrap.go # 编译模式预构建实例
└── README.md
```

- 解析后约束挂在 `node.ClassProperty.Annotations`（`*data.ClassValue`，`Class` 为 `*ConstraintClass`）。
- 运行时可通过 `ConstraintClass.Spec()`、`State()` 读取约束名与构造参数。
- `gen-std` 根据命名空间含 `annotation` 自动生成 `#[\Attribute(...)]` 伪代码，无需在 `generate.go` 中硬编码类名。

## 运行时校验

### Validation\Validator

```php
$violations = \Validation\Validator::validate($dto);
// [ ['field' => 'username', 'message' => '...'], ... ]
```

### HTTP 集成

控制器 DTO 参数在 `Request::bind()` 之后自动校验（扫描阶段若检测到约束注解则 `Validate=true`）。失败时返回 **422**，body 含 `violations` 数组。

## 后续扩展
