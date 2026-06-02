# Spring 示例 - 快速参考

## 🚀 快速启动

```bash
# 1. 启动服务器
./origami run examples/spring/index.php

# 2. 测试 API（新终端）
curl http://127.0.0.1:8080/api/hello
```

## 📁 项目结构

```
examples/spring/
├── index.php                 # 入口文件
├── README.md                 # 详细文档
├── test_api.sh              # Linux/Mac 测试脚本
├── test_api.ps1             # Windows 测试脚本
└── src/
    ├── main.php             # 应用主配置
    ├── controllers/         # 控制器层
    │   ├── HelloController.php
    │   ├── UserController.php
    │   ├── ProductController.php
    │   └── AuthController.php
    ├── services/            # 服务层
    │   ├── UserService.php
    │   ├── ProductService.php
    │   └── AuthService.php
    ├── models/              # 模型层
    │   ├── User.php
    │   ├── Product.php
    │   └── Users.php
    ├── dto/                 # 数据传输对象
    │   └── ResponseDTO.php
    ├── middleware/          # 中间件
    │   ├── AuthMiddleware.php
    │   └── CorsMiddleware.php
    └── config/              # 配置
        └── AppConfig.php
```

## 🔗 API 端点速查

### 基础接口
```
GET  /api/hello              # 问候
GET  /api/info               # 应用信息
GET  /api/status             # 状态检查
```

### 用户管理
```
GET  /api/users              # 用户列表
GET  /api/user/{id}          # 用户详情
POST /api/users              # 创建用户
```

### 商品管理
```
GET  /api/products           # 商品列表
GET  /api/product/{id}       # 商品详情
POST /api/products           # 创建商品
PUT  /api/product/{id}       # 更新商品
DELETE /api/product/{id}     # 删除商品
GET  /api/products/search    # 搜索商品
```

### 认证相关
```
POST /api/auth/login         # 登录
POST /api/auth/register      # 注册
GET  /api/auth/profile       # 用户信息（需认证）
POST /api/auth/logout        # 登出
```

## 💡 常用注解

```php
// 类级别
#[Controller]                // 标记控制器
#[Route(prefix: "/api")]     // 路由前缀

// 方法级别
#[GetMapping(path: "/users")]    // GET 请求
#[PostMapping(path: "/users")]   // POST 请求
#[PutMapping(path: "/user/{id}")] // PUT 请求
#[DeleteMapping(path: "/user/{id}")] // DELETE 请求

// 应用级别
#[Application(name: "AppName")]  // 应用入口
```

## 📝 代码示例

### 控制器
```php
#[Controller]
#[Route(prefix: "/api")]
class UserController {
    
    private $userService;
    
    public function __construct() {
        $this->userService = new UserService();
    }
    
    #[GetMapping(path: "/users")]
    public function list($request, $response) {
        $users = $this->userService->findAll();
        $response->json([
            "code" => 200,
            "message" => "success",
            "data" => $users
        ]);
    }
}
```

### 服务层
```php
class UserService {
    public function findAll() {
        // 业务逻辑
        return $users;
    }
    
    public function findById($id) {
        // 查找逻辑
        return $user;
    }
}
```

### 中间件
```php
// 在 index.php 中使用
$server->middleware(function ($request, $response, $next) {
    Log::info("请求: " . $request->path());
    $next($request, $response);
});
```

## 🧪 测试命令

```bash
# 获取用户列表
curl http://127.0.0.1:8080/api/users

# 获取商品详情
curl http://127.0.0.1:8080/api/product/1

# 创建商品
curl -X POST http://127.0.0.1:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"商品","price":99.99,"category":"分类"}'

# 用户登录
curl -X POST http://127.0.0.1:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"123456"}'
```

## ⚙️ 配置项

```php
AppConfig::get('app.name')           // 应用名称
AppConfig::get('app.version')        // 应用版本
AppConfig::get('server.port')        // 服务器端口
AppConfig::get('api.prefix')         // API 前缀
AppConfig::get('auth.token_expiry')  // Token 过期时间
```

## 🔑 默认账号

```
用户名: admin
密码: 123456

用户名: user1
密码: password
```

## 📊 响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {},
  "timestamp": 1234567890
}
```

## 🎯 关键特性

- ✅ 分层架构 (Controller → Service → Model)
- ✅ RESTful API 设计
- ✅ 统一响应格式
- ✅ 注解驱动路由
- ✅ 中间件支持
- ✅ Token 认证
- ✅ CORS 支持
- ✅ 配置管理

## 📖 更多信息

查看完整文档：[README.md](README.md)
查看项目总结：[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)
