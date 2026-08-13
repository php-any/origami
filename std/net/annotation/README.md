# HTTP 控制器请求参数自动注入

本文说明 Origami 中基于注解的 HTTP 控制器**形参自动绑定**机制：从路由扫描到运行时注入、校验、调用控制器的完整流程。

## 总览

开发者只需在控制器方法上声明形参类型（及可选的 Validation 约束），框架在**扫描期**生成绑定计划 `HandlerSpec`，在**请求期**按该计划从 HTTP 请求中解析实参并调用方法。

```
应用启动 / 类加载
  └─ #[Controller] 扫描方法上的 @*Mapping
       └─ MacroExpander.Expand() → AnalyzeHandlerParams() → HandlerSpec
            └─ 写入 PendingRoute，最终注册到全局路由表

HTTP 请求命中路由
  └─ executeMiddlewareChain() → executeControllerMethod()
       └─ resolveHandlerArgs(HandlerSpec) → 按来源注入各形参
            └─ （可选）Validation 校验
                 └─ CallHTTPControllerMethod($receiver, $method, $args)
                      └─ writeHandlerReturnValue()（若方法 return array/object）
```

**核心设计**：绑定（类型转换 / 取值）与约束校验（`#[NotBlank]`、`#[Size]` 等）分层处理，互不混用。

| 阶段 | 包 / 文件 | 职责 |
|------|-----------|------|
| 扫描期计划 | `std/net/annotation/macro.go` | 分析形参 → `HandlerSpec` |
| 路由注册 | `std/net/annotation/controller_class.go` | `@Controller` + `@*Mapping` → `PendingRoute` |
| 运行时解析 | `std/net/http/argument_resolver.go` | 按 `HandlerSpec` 从请求取值 |
| DTO 绑定 | `std/net/http/request_bind_method.go` | `Request::bind()` query/form/body → 类实例 |
| 约束校验 | `std/validation/validator.go` | DTO 属性 / 形参约束 |
| 方法调用 | `std/net/http/route_dispatch.go` | 中间件链 + 控制器调用 |

---

## 扫描期：生成 HandlerSpec

### 触发时机

类带有 `#[Controller]` 时，其构造逻辑会遍历类中所有方法，读取方法上的 `#[GetMapping]`、`#[PostMapping]` 等映射注解。每个映射注解实现 **`MacroExpander`** 接口，在注册路由前调用 `Expand()`：

```go
// std/net/annotation/macro.go
type MacroExpander interface {
    Expand(target data.Method, ctx data.Context, routePath string) (effective data.Method, spec HandlerSpec, acl data.Control)
}
```

`@GetMapping` / `@PostMapping` / `@PutMapping` / `@DeleteMapping` 的 `Expand()` 均委托给 `AnalyzeHandlerParams()`，根据**方法形参**与**完整路由路径**生成 `HandlerSpec`。

### HandlerSpec 结构

```go
// std/net/data/handler_spec.go
type HandlerSpec struct {
    Params []ParamBinding
}

type ParamBinding struct {
    Name        string              // 形参名，如 name、request
    Label       string              // #[Name] 展示名，默认同 Name
    TypeFQN     string              // 类型全名，如 string、Spring\DTO\Request\LoginRequest
    Source      BindingSource       // 值来源（见下表）
    Index       int                 // 在形参列表中的下标
    Validate    bool                // 是否执行校验
    PathKey     string              // SourcePath 时路由模板 {key}
    QueryKey    string              // SourceQuery 时 ?key=
    Nullable    bool                // 是否可空类型 ?int 等
    Constraints []*data.ClassValue  // 形参上的 Validation 约束实例
}
```

### 形参来源推断规则

`AnalyzeHandlerParams()` 按形参**类型**与路由路径中的 `{变量}` **自动推断**绑定来源：

| 形参类型 | 条件 | Source | 运行时取值 |
|----------|------|--------|------------|
| `Net\Http\Request` | — | `SourceRequest` | 当前请求的 Request 代理 |
| `Net\Http\Response` | — | `SourceResponse` | 当前请求的 Response 代理 |
| `string` / `int` / `float` / `bool` | 路由含 `{形参名}` | `SourcePath` | `Request::pathValue(key)` |
| `string` / `int` / `float` / `bool` | 路由不含同名路径变量 | `SourceQuery` | `Request::input(key)` → `?key=` |
| 其他 class（DTO） | — | `SourceDTO` | `Request::bind(ClassName)` |
| `Net\Http\UploadedFile` | — | `SourceFormFile` | multipart 字段（名同形参） |
| **无类型** | — | — | **扫描期报错** |

路径变量从路由模板解析，例如 `/api/product/{id}` 得到 `id`：

```go
var pathVarPattern = regexp.MustCompile(`\{([^}/:]+)\}`)
```

**命名约定**：路径参数名须与形参名一致（`int $id` 对应 `{id}`）；查询参数名默认与形参名一致（`string $name` 对应 `?name=`）。

### 校验元数据收集

- **标量形参**（Query / Path）：从形参上的 `Validation\Annotation\*` 收集约束（`#[NotBlank]`、`#[Size]` 等），写入 `ParamBinding.Constraints`，并设 `Validate = true`。
- **DTO 形参**：扫描 DTO 类**属性**上是否有 Validation 注解，有则 `Validate = true`（校验在 `ValidateObject()` 中按属性执行）。

形参注解在解析阶段挂载：`node.Parameter.Annotations`，由 `parser/parameter_parser.go` + `parser/annotation_apply.go` 完成。

---

## 运行时：解析并注入实参

路由命中后，`executeControllerMethod()` 调用 `resolveHandlerArgs()`，按 `HandlerSpec.Params` 顺序填充实参数组，再交给 `CallHTTPControllerMethod()`。

```
resolveHandlerArgs(spec, reqProxy, resProxy)
  ├─ SourceRequest   → reqProxy
  ├─ SourceResponse  → resProxy
  ├─ SourcePath      → pathValue → 类型转换 → （可选）ValidateConstraints
  ├─ SourceQuery     → input     → 类型转换 → （可选）ValidateConstraints
  └─ SourceDTO       → Request::bind() → （可选）ValidateObject
```

### Request / Response

直接注入框架创建的代理对象，与手写 `$request`、`$response` 相同。

### 路径参数（SourcePath）

1. 调用 `Request::pathValue($pathKey)` 取原始字符串
2. `resolveScalarValue()` 做类型转换（`strconv`）
3. 若形参有约束注解 → `ValidateConstraints()`，失败返回 **422**

### 查询参数（SourceQuery）

1. 调用 `Request::input($queryKey)` 取值
2. 缺失时：`string` → `""`；`?int` 等可空 → `null`；有 Validation 约束的非 string → `null`（留给 `#[NotBlank]` 等处理，避免绑定阶段 Fatal）
3. 类型转换 + 约束校验（同路径参数）

### DTO（SourceDTO）

1. 调用 `Request::bind($className)` 实例化并填充：
   - `Content-Type: application/json` 且有 body → JSON 反序列化
   - 否则有 form → 按属性名绑定 form 字段
   - 否则有 query → 按属性名绑定查询参数（`bindFlatMapToClass`，按属性 PHP 类型转换，**不经过** JsonSerializer 的字符串 hack）
2. 若 `Validate == true` → `validation.ValidateObject($dto)`，失败 **422**

### 校验失败响应

统一通过 `Response::error('validation failed', 422, $violations)` 写出 JSON：

```json
{
  "code": 422,
  "message": "validation failed",
  "data": [
    { "field": "name", "message": "name 不能为空" }
  ],
  "timestamp": ...
}
```

校验失败时 `resolveHandlerArgs` 返回 `nil` 实参，控制器**不会**被执行。

---

## 使用示例

### 1. 查询参数 + 形参约束（标量）

适合少量简单查询参数。

```php
#[GetMapping(path: "/hello")]
public function hello(
    #[NotBlank(message: "name 不能为空")]
    #[Size(min: 2, max: 64)]
    string $name,
    #[NotBlank(message: "age 不能为空")]
    int $age
): array {
    return [
        "code" => 200,
        "message" => "success",
        "data" => ["hello" => $name, "age" => $age],
    ];
}
```

- `GET /hello?name=ab&age=25` → 200，自动 JSON 写出返回值
- `GET /hello?name=ab` → 422，`age` 字段违规
- `GET /hello` → 422，`name` 字段违规

### 2. 路径参数

```php
#[GetMapping(path: "/api/product/{id}")]
public function show(int $id, Response $response): void {
    $response->success(['id' => $id]);
}
```

`$id` 来自路径 `/api/product/42`，不占用 query。

### 3. DTO + 属性约束（推荐用于多字段查询 / Body）

```php
// Spring\DTO\Request\UserListQuery
class UserListQuery {
    #[Min(value: 0)]
    public int $min_age = 0;

    #[Min(value: 1)]
    #[Max(value: 100)]
    public int $limit = 10;
}

#[GetMapping(path: "/simple/users")]
public function singleTableUsers(UserListQuery $query, Response $response): void {
    // $query 已从 ?min_age=&limit= 绑定并校验
    $response->success([...]);
}
```

### 4. JSON Body（登录等）

```php
#[PostMapping(path: "/login")]
public function login(LoginRequest $request, Response $response): void {
    // LoginRequest 从 application/json body 绑定
    // 属性上的 #[Size]、#[Pattern] 等自动校验
}
```

### 5. 文件上传（multipart/form-data）

单文件形参（字段名与形参名一致）：

```php
#[PostMapping(path: "/api/upload/avatar")]
public function uploadAvatar(
    #[NotBlank(message: "请上传 avatar 文件")]
    UploadedFile $avatar,
    Response $response
): void {
    $path = $avatar->store('/tmp/uploads');
    $response->success(['saved_path' => $path], '上传成功', 201);
}
```

`curl -F "avatar=@photo.png" http://localhost:8080/api/upload/avatar`

DTO 混合文本字段与文件（`examples/spring` 中的 `UploadFileRequest`）：

```php
class UploadFileRequest {
    #[NotBlank]
    public string $title;
    #[NotBlank]
    public UploadedFile $file;
}

#[PostMapping(path: "/api/upload")]
public function upload(UploadFileRequest $request, Response $response): void {
    $request->file->store($uploadDir);
}
```

`UploadedFile` 方法：`originalName()`、`size()`、`mimeType()`、`extension()`、`isValid()`、`getContent()`、`store($directory, $name = null)`。

### 6. 经典 Request / Response 写法（仍支持）

```php
#[GetMapping(path: "/api/hello")]
public function hello(Request $request, Response $response): void {
    $response->success(['greeting' => 'Hello World!']);
}
```

未使用自动注入时，行为与改造前一致。

---

## 绑定 vs 校验：分层说明

| 层次 | 机制 | 示例 |
|------|------|------|
| **绑定** | 按 PHP 类型从 HTTP 取值并转换 | `?age=25` → `int 25`；缺失的 `string` → `""` |
| **约束** | `Validation\Annotation\*` | `#[NotBlank]`、`#[Min]`、`#[Size]` |

注意：

- `#[Size]` 对标量 `int` 按**字符串长度**语义工作（与 DTO 字符串字段一致）；数值范围请用 `#[Min]` / `#[Max]`。
- 无约束的必填 `int` 查询参数缺失时，绑定阶段仍会抛出转换错误（Fatal）；需要友好 422 时请添加 `#[NotBlank]`。
- DTO 属性默认值（如 `public int $limit = 10`）在 query 未传该字段时保留默认值。

更完整的约束列表见 [`std/validation/annotation/README.md`](../../validation/annotation/README.md)。

---

## 控制器返回值

若方法签名返回 `array` / 对象，且未通过 `$response` 写出内容，框架在请求末尾调用 `writeHandlerReturnValue()` 自动 JSON 序列化。

若使用 `#[Middleware]` 拦截器，中间件应 **`return $next($request, $response)`**（或在 `$next` 之后 `return $result`），以便返回值穿透中间件链。框架在 Go 层对 void 中间件做了兜底，但显式 `return` 仍是推荐写法。

---

## 相关源码索引

| 文件 | 说明 |
|------|------|
| `std/net/data/handler_spec.go` | `HandlerSpec` / `ParamBinding` / `BindingSource` |
| `std/net/data/pending.go` | 待注册路由与 `HandlerSpec` 持久化 |
| `std/net/annotation/macro.go` | `AnalyzeHandlerParams`、来源推断 |
| `std/net/annotation/controller_class.go` | `@Controller` 扫描与路由注册 |
| `std/net/annotation/*_mapping_class.go` | 各 HTTP 映射宏注解 |
| `std/net/http/argument_resolver.go` | 运行时参数解析与校验 |
| `std/net/http/route_dispatch.go` | 中间件链 + 控制器调用 |
| `std/net/http/request_bind_method.go` | `Request::bind` DTO 绑定 |
| `std/net/http/controller_return.go` | 返回值自动 JSON |
| `std/validation/validator.go` | `ValidateObject` / `ValidateConstraints` |
| `parser/parameter_parser.go` | 形参注解解析 |
| `node/function.go` | `Parameter.Annotations` |

## 测试参考

| 测试文件 | 覆盖点 |
|----------|--------|
| `std/net/annotation/macro_test.go` | HandlerSpec 生成、路径/查询推断 |
| `std/net/http/argument_resolver_test.go` | 查询/路径解析、标量校验 422 |
| `std/net/http/home_controller_integration_test.go` | 端到端 HomeController 绑定与返回 |
| `tests/net/handler_args_test.zy` | Zy 脚本级集成测试 |

---

## 快速对照表

| 你想注入… | 写法 | 来源 |
|-----------|------|------|
| 整个请求对象 | `Request $request` | 框架代理 |
| 响应对象 | `Response $response` | 框架代理 |
| `?page=1` | `int $page` | Query |
| `/user/{id}` | `int $id` | Path |
| 多个查询字段 | `MyQuery $query`（DTO） | Query → bind |
| POST JSON | `LoginRequest $dto`（DTO） | Body → bind |
| 上传文件 | `UploadedFile $avatar` | multipart 字段 `avatar` |
| 文本 + 文件 | `UploadFileRequest $dto` | multipart → bind |
| 简单字段校验 | 形参上 `#[NotBlank]` 等 | ValidateConstraints |
| 多字段校验 | DTO 属性上 `#[Size]` 等 | ValidateObject |
