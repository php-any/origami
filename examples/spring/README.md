# Spring 风格示例

这是一个展示如何使用 Origami 框架的 Spring Boot 风格注解开发 Web 应用的示例项目。

## 功能特性

- ✅ **注解驱动路由** - 使用 `@Controller`、`@GetMapping`、`@PostMapping` 等注解定义路由
- ✅ **路径参数** - 支持 RESTful 风格的路径参数 `{id}`
- ✅ **请求处理** - 完整的请求/响应处理能力
- ✅ **JSON 响应** - 便捷的 JSON 数据返回
- ✅ **中间件支持** - 日志、认证等中间件
- ✅ **分层架构** - Controller → Service → Model 清晰的分层结构
- ✅ **依赖注入** - 通过构造函数和属性注入依赖

## 项目结构

```
spring/
├── index.php              # 应用入口，启动 HTTP 服务器
├── src/
│   ├── main.php          # 主应用配置，使用 @Application 注解
│   ├── controllers/      # 控制器层
│   │   ├── HelloController.php    # 基础示例控制器
│   │   ├── UserController.php     # 用户管理控制器
│   │   ├── ProductController.php  # 商品管理控制器
│   │   └── AuthController.php     # 认证控制器
│   ├── services/         # 服务层
│   │   ├── UserService.php        # 用户服务
│   │   ├── ProductService.php     # 商品服务
│   │   └── AuthService.php        # 认证服务
│   ├── models/           # 数据模型层
│   │   ├── User.php               # 用户模型
│   │   └── Product.php            # 商品模型
│   ├── dto/              # 数据传输对象
│   │   ├── UserDTO.php            # 用户 DTO
│   │   └── ResponseDTO.php        # 统一响应 DTO
│   └── middleware/       # 中间件
│       ├── AuthMiddleware.php     # 认证中间件
│       └── LogMiddleware.php      # 日志中间件
```

## 快速开始

### 1. 运行示例

```bash
# 在项目根目录运行
./origami run examples/spring/index.php
```

或者：

```bash
php examples/spring/index.php
```

### 2. 访问 API

服务将在 `http://127.0.0.1:8080` 启动，可以访问以下端点：

#### 基础示例
- `GET /api/hello` - 简单的问候接口
- `GET /api/users` - 获取用户列表
- `GET /api/user/{id}` - 获取指定用户详情

#### 商品管理
- `GET /api/products` - 获取商品列表
- `GET /api/product/{id}` - 获取指定商品详情
- `POST /api/products` - 创建新商品
- `PUT /api/product/{id}` - 更新商品信息
- `DELETE /api/product/{id}` - 删除商品

#### 认证相关
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册
- `GET /api/auth/profile` - 获取当前用户信息（需要认证）

## 代码示例

### 控制器示例

```php
<?php

namespace App\Controller;

use Net\Annotation\Controller;
use Net\Annotation\Route;
use Net\Annotation\GetMapping;
use Net\Annotation\PostMapping;

#[Controller]
#[Route(prefix: "/api")]
class UserController {
    
    #[GetMapping(path: "/users")]
    public function listUsers($request, $response) {
        $users = [
            ["id" => 1, "name" => "张三", "email" => "zhangsan@example.com"],
            ["id" => 2, "name" => "李四", "email" => "lisi@example.com"]
        ];
        
        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => $users
        ]);
    }
    
    #[GetMapping(path: "/user/{id}")]
    public function getUser($request, $response) {
        $id = $request->pathValue('id');
        
        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => [
                "id" => $id,
                "name" => "张三",
                "email" => "zhangsan@example.com"
            ]
        ]);
    }
    
    #[PostMapping(path: "/users")]
    public function createUser($request, $response) {
        $body = $request->body();
        
        // 处理创建逻辑...
        
        $response->status(201)->json([
            "code" => 201,
            "message" => "created",
            "data" => $body
        ]);
    }
}
```

### 服务层示例

```php
<?php

namespace App\Service;

class UserService {
    
    private $users = [];
    
    public function __construct() {
        // 初始化数据或连接数据库
        $this->users = [
            1 => ["id" => 1, "name" => "张三", "email" => "zhangsan@example.com"],
            2 => ["id" => 2, "name" => "李四", "email" => "lisi@example.com"]
        ];
    }
    
    public function findAll() {
        return array_values($this->users);
    }
    
    public function findById($id) {
        return $this->users[$id] ?? null;
    }
    
    public function create($data) {
        $id = count($this->users) + 1;
        $user = array_merge(["id" => $id], $data);
        $this->users[$id] = $user;
        return $user;
    }
    
    public function update($id, $data) {
        if (!isset($this->users[$id])) {
            return null;
        }
        $this->users[$id] = array_merge($this->users[$id], $data);
        return $this->users[$id];
    }
    
    public function delete($id) {
        if (!isset($this->users[$id])) {
            return false;
        }
        unset($this->users[$id]);
        return true;
    }
}
```

### 中间件示例

```php
<?php

// 在 index.php 中使用中间件
$server->middleware(function ($request, $response, $next) {
    // 记录请求日志
    $method = $request->method();
    $path = $request->path();
    Log::info("HTTP " . $method . " " . $path);
    
    // 执行下一个中间件或路由处理器
    $next($request, $response);
});
```

## 可用的注解

### 类级别注解

- `#[Controller]` - 标记一个类为控制器
- `#[Route(prefix: "/path")]` - 设置路由前缀

### 方法级别注解

- `#[GetMapping(path: "/path")]` - 处理 GET 请求
- `#[PostMapping(path: "/path")]` - 处理 POST 请求
- `#[PutMapping(path: "/path")]` - 处理 PUT 请求
- `#[DeleteMapping(path: "/path")]` - 处理 DELETE 请求

### 应用级别注解

- `#[Application(name: "AppName")]` - 标记应用入口函数

## 最佳实践

1. **分层架构** - 保持 Controller → Service → Model 的清晰分层
2. **统一响应格式** - 使用统一的 JSON 响应结构
3. **错误处理** - 合理使用 HTTP 状态码
4. **RESTful 设计** - 遵循 RESTful API 设计规范
5. **输入验证** - 在 Service 层进行数据验证
6. **日志记录** - 关键操作记录日志

## 注意事项

- 所有控制器必须使用 `#[Controller]` 注解标记
- 路由路径支持路径参数，使用 `{paramName}` 语法
- 可以通过 `$request->pathValue('paramName')` 获取路径参数
- 使用 `$response->json()` 方法返回 JSON 数据
- 使用 `$response->status(code)` 设置 HTTP 状态码

## 扩展阅读

- [Origami 文档](../../docs/README.md)
- [注解系统文档](../../docs/annotations.md)
- [HTTP 服务器文档](../../docs/std/Net/Http/)

## 许可证

MIT License
