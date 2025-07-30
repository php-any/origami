# Annotations (注解系统)

Origami 支持强大的注解系统，允许你在编译时和运行时对代码进行元编程。注解系统分为两种类型：特性注解和宏注解。

## 概述

注解是 Origami 语言中的元编程特性，允许你：

- 在编译时修改语法节点
- 在运行时通过反射获取注解信息
- 实现依赖注入、路由映射等功能
- 添加元数据到类、方法、属性等

## 注解类型

### 1. 特性注解 (Feature Annotations)

特性注解只接收注解声明的参数，主要用于标记和元数据存储。

**特点：**

- 只接收注解参数
- 不修改被注解的节点
- 主要用于反射和元数据

**示例：**

```php
@Controller(name: "UserController")
@Route(prefix: "/api/users")
@GetMapping(path: "/list")
```

### 2. 宏注解 (Macro Annotations)

宏注解可以接收注解参数和被注解的节点，能够修改语法节点。

**特点：**

- 接收注解参数
- 可以修改被注解的节点
- 能够添加默认值、修改属性等

**示例：**

```php
@Inject(service: "UserService")
public $userService;
```

## 内置注解

### Controller 注解

用于标记控制器类。

**语法：**

```php
@Controller(name?: string)
```

**参数：**

- `name` (可选): 控制器名称

**示例：**

```php
@Controller(name: "UserController")
class UserController {
    // 控制器逻辑
}
```

**方法：**

- `process()`: 处理注解
- `register()`: 注册控制器

### Route 注解

用于设置路由前缀。

**语法：**

```php
@Route(prefix: string)
```

**参数：**

- `prefix`: 路由前缀

**示例：**

```php
@Route(prefix: "/api/users")
class UserController {
    // 所有路由都会以 /api/users 开头
}
```

**方法：**

- `process()`: 处理注解
- `register()`: 注册路由

### Inject 注解 (宏注解)

用于依赖注入，可以修改被注解的属性。

**语法：**

```php
@Inject(service: string)
```

**参数：**

- `service`: 要注入的服务名称

**示例：**

```php
@Inject(service: "UserService")
public $userService;
```

**功能：**

- 为属性添加默认值
- 实现依赖注入
- 修改语法节点

**方法：**

- `process()`: 处理注解
- `inject()`: 执行注入
- `__construct()`: 构造函数，接收注解参数和被注解节点

### GetMapping 注解

用于标记 GET 请求映射。

**语法：**

```php
@GetMapping(path: string)
```

**参数：**

- `path`: 请求路径

**示例：**

```php
@GetMapping(path: "/list")
public function getUserList() {
    return $this->userService->getAllUsers();
}
```

**方法：**

- `process()`: 处理注解
- `mapping()`: 注册映射

## 完整示例

```php
namespace App\Controller;

use Annotation\Route;
use Annotation\Controller;
use Annotation\GetMapping;
use Annotation\Inject;

@Controller
@Route(prefix: "/api/users")
class UserController {
    // 宏注解允许编辑语法节点, 可以删除或者添加信息, 比如为$userService属性添加默认值, 宏注解就能实现
    @Inject(service: "UserService")
    public $userService;

    // 特性注解, 函数经过反射后, 能获取注解类信息
    @GetMapping(path: "/list")
    public function getUserList() {
        return $this->userService->getAllUsers();
    }
}

echo new UserController()->userService;
```

## 注解处理流程

### 编译时处理

1. **解析注解**: 解析注解名称和参数
2. **类型检查**: 根据注解类型（特性/宏）进行不同处理
3. **节点修改**: 宏注解可以修改被注解的节点
4. **元数据存储**: 特性注解存储元数据用于反射

### 运行时处理

1. **反射获取**: 通过反射获取注解信息
2. **注解执行**: 调用注解的处理方法
3. **功能实现**: 实现具体的业务逻辑

## 注解类型区别

| 特性     | 特性注解                | 宏注解                   |
| -------- | ----------------------- | ------------------------ |
| 参数接收 | 只接收注解参数          | 接收注解参数和被注解节点 |
| 节点修改 | 不修改节点              | 可以修改节点             |
| 主要用途 | 标记和元数据            | 代码生成和修改           |
| 示例     | `@Controller`, `@Route` | `@Inject`                |

## 自定义注解

你可以创建自定义注解类：

```php
// 特性注解示例
class MyFeatureAnnotation {
    public function __construct($param) {
        // 只接收注解参数
    }

    public function process() {
        // 处理逻辑
    }
}

// 宏注解示例
class MyMacroAnnotation {
    public function __construct($param, $target) {
        // 接收注解参数和被注解节点
    }

    public function process() {
        // 处理逻辑
    }
}
```

## 最佳实践

1. **命名规范**: 注解类名使用 PascalCase
2. **参数验证**: 在构造函数中验证必要参数
3. **错误处理**: 提供清晰的错误信息
4. **文档注释**: 为注解添加详细的文档说明
5. **测试覆盖**: 为自定义注解编写测试用例

## 注意事项

1. 宏注解可以修改语法节点，使用时要谨慎
2. 特性注解主要用于反射，不会修改代码结构
3. 注解参数支持命名参数和位置参数
4. 注解处理在编译时进行，确保性能
5. 复杂的注解逻辑应该在运行时处理

## 相关文档

- [Classes](./classes.md) - 类定义和继承
- [Functions](./functions.md) - 函数定义和调用
- [Reflection](./reflection.md) - 反射系统
- [Go Integration](./go-integration.md) - Go 语言集成
