# 设计：`zy compile` — Vendor 目录 AOT 编译为 Go 包

## 目标

将 `vendor/` 目录下的 PHP 文件预编译为 Go 源码（AST 结构体字面量），消除运行时对 vendor 代码的词法分析和解析开销，提升启动性能。

**核心原则**：预编译是透明优化，不改变任何语义。无法预编译的场景自动回退到原有解释器。

## 架构

```
┌─────────────────────────────────────────────────────────┐
│                    构建时 (zy compile vendor/)            │
│                                                          │
│  vendor/                                                 │
│  ├─ monolog/monolog/src/Logger.php                      │
│  ├─ laravel/framework/src/.../*.php                     │
│  └─ ...                                                  │
│         │                                                │
│         ▼                                                │
│  遍历 vendor/**/*.php → 逐个 Parser.ParseFile() → AST   │
│         │                                                │
│         ▼                                                │
│  AST → Go 代码生成器 (每个节点类型一个生成函数)           │
│         │                                                │
│         ▼                                                │
│  .zy/build/                                              │
│  ├─ go.mod                                               │
│  ├─ vendor_ast.go    (所有 vendor 文件的预构建 AST)       │
│  └─ register.go      (Register(vm) 注册函数)             │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│                    运行时                                │
│                                                          │
│  用户代码 app.php                                        │
│         │                                                │
│         ▼                                                │
│  Parser.ParseFile("app.php") → AST (正常解析)            │
│         │                                                │
│         ▼  执行过程中遇到 vendor 类/函数                  │
│  ClassPathManager → vendor/xxx/Class.php                 │
│         │                                                │
│         ▼                                                │
│  VM.LoadAndRun() → 检查 compiledFiles 注册表             │
│  ├─ 命中 → 直接执行预构建 AST (跳过词法+解析)             │
│  └─ 未命中 → 回退到正常解析                               │
└─────────────────────────────────────────────────────────┘
```

## 核心组件

### 1. 文件扫描器 (`cmd/compile_scan.go`)

扫描 vendor 目录下所有 `.php` 文件：

```go
func collectVendorFiles(vendorDir string) ([]string, error) {
    filepath.WalkDir(vendorDir, func(path string, d fs.DirEntry, err error) error {
        if strings.HasSuffix(path, ".php") {
            files = append(files, path)
        }
        return nil
    })
    return files, nil
}
```

### 2. AST→Go 代码生成器 (`cmd/compile_gen.go`)

为每种节点类型写一个生成函数，输出 Go 构造代码：

```go
type Generator struct {
    buf     strings.Builder
    indent  int
    imports map[string]bool
}

func (g *Generator) Generate(file string, program *node.Program) string {
    // 生成 package 声明、import、函数
    // 遍历 program.Statements，递归生成每个节点
}

// 核心分派函数
func (g *Generator) genGetValue(v data.GetValue) {
    switch n := v.(type) {
    case *node.IntLiteral:       g.genIntLit(n)
    case *node.StringLiteral:    g.genStringLit(n)
    case *node.VariableExpression: g.genVar(n)
    case *node.BinaryAdd:        g.genBinaryAdd(n)
    case *node.IfStatement:      g.genIf(n)
    case *node.Call:             g.genCall(n)
    // ... 130+ 类型，分阶段实现
    default:
        g.printf("nil // TODO: unsupported node type %T", n)
    }
}
```

**变量索引处理**：AST 节点中已包含 ScopeManager 分配的变量索引（整数字面量），生成时直接嵌入：

```go
// PHP: $x = 1 + 2; echo $x;
// 生成:
node.NewBinaryAssign(from,
    node.NewVariableExpression(from, 0),  // $x → 索引 0
    node.NewBinaryAdd(from,
        node.NewIntLiteral(from, 1),
        node.NewIntLiteral(from, 2)))
node.NewCall(from, "echo",
    []data.GetValue{node.NewVariableExpression(from, 0)})
```

### 3. VM 注册表扩展 (`runtime/vm.go`)

给 VM 添加预编译文件注册表：

```go
type CompiledFileFunc func() (*node.Program, []parser.Variable)

type VM struct {
    // ... 现有字段 ...
    compiledFiles map[string]CompiledFileFunc
}

func (vm *VM) RegisterCompiledFile(path string, fn CompiledFileFunc) {
    vm.compiledFiles[path] = fn
}
```

修改 `LoadAndRun`，优先查找预编译注册表：

```go
func (vm *VM) LoadAndRun(file string) (data.GetValue, data.Control) {
    // 1. 检查预编译注册表（用相对路径匹配）
    relFile := toRelative(file)
    if fn, ok := vm.compiledFiles[relFile]; ok {
        program, vars := fn()
        ctx := vm.CreateContext(vars)
        vm.RegisterGlobalContext(vars, ctx)
        result, ctrl := program.GetValue(ctx)
        if data.FlushAllBuffersFn != nil {
            data.FlushAllBuffersFn()
        }
        return result, ctrl
    }

    // 2. 回退到原有解析逻辑（保持不变）
    data.ResetUserOutput()
    p := vm.parser.Clone()
    program, acl := p.ParseFile(file)
    // ... 现有逻辑 ...
}
```

### 4. 生成的 Go 包结构

```
.zy/build/
├── go.mod                  # module build
├── vendor_ast.go           # 所有 vendor 文件的预构建 AST
└── register.go             # Register(vm) 注册函数
```

`vendor_ast.go` 示例：

```go
package build

import (
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
    "github.com/php-any/origami/parser"
)

func MonologLoggerAST() (*node.Program, []parser.Variable) {
    from := data.From{File: "vendor/monolog/monolog/src/Logger.php"}
    stmts := []data.GetValue{
        node.NewClass(from, "Logger", nil, nil, []data.GetValue{
            // methods...
        }),
    }
    return node.NewProgram(from, stmts), nil
}
```

`register.go` 示例：

```go
package build

import "github.com/php-any/origami/data"

func Register(vm data.VM) {
    vm.RegisterCompiledFile("vendor/monolog/monolog/src/Logger.php", MonologLoggerAST)
    vm.RegisterCompiledFile("vendor/laravel/.../App.php", LaravelAppAST)
    // ... 所有 vendor 文件
}
```

### 5. CLI 接口

```bash
# 基本用法
zy compile vendor/

# 指定输出目录
zy compile vendor/ -o ./dist

# 指定包名（默认 "build"）
zy compile vendor/ --pkg myapp
```

用户集成方式：

```go
import build ".zy/build"

func main() {
    vm := getRuntimeVM()
    build.Register(vm)
    vm.LoadAndRun("app.php")
}
```

## 节点类型覆盖策略

### P0 — 首批实现（覆盖 vendor 代码 90%+）

| 分类 | 节点类型 |
|------|----------|
| 字面量 | `IntLiteral`, `FloatLiteral`, `StringLiteral`, `BooleanLiteral`, `NullLiteral` |
| 变量 | `VariableExpression`, `ValueReference` |
| 二元运算 | `BinaryAdd/Sub/Mul/Quo/Rem/Pow`, `BinaryDot`, `BinaryEq/Ne/Lt/Le/Gt/Ge`, `BinaryEqStrict/NeStrict`, `BinaryLand/Lor`, `BinaryBitwise/Shift`, `BinaryAssign`（所有复合赋值）, `BinarySpaceship`, `BinaryNullCoalesce` |
| 一元运算 | `UnaryIncr/Decr`, `PostfixIncr/Decr`, `ErrorSuppress`, `UnaryNot/BitwiseNot` |
| 控制流 | `If`, `For`, `Foreach`, `While`, `DoWhile`, `Switch`, `Match`, `Break`, `Continue`, `Return`, `Throw`, `Try` |
| 函数调用 | `Call`, `CallMethod`, `CallStaticMethod`, `CallObjectMethod`, `CallParentMethod`, `CallSelfMethod`, `CallNullsafe` |
| OOP | `Class`, `ClassAbstract`, `Interface`, `Trait`, `Enum`, `New`, `Clone`, `InstanceOf`, `StaticClass`, `SelfClass`, `Parent`, `InitClass`, `ClassConstant` |
| 数组 | `Array`, `ArraySpread`, `Index`, `KV`, `Compact` |
| 表达式 | `Ternary`, `NullCoalesce`, `Range`, `SpreadArgument`, `NameArgument` |
| 声明 | `Function`, `Closure`, `Lambda`, `Const`, `Namespace`, `Use`, `IncludeStatement` |

### P1 — 第二批

| 分类 | 节点类型 |
|------|----------|
| 高级调用 | `CallObjectProperty`, `CallStaticKeywordMethod/Property`, `CallObjectDynamicMethod/Property` |
| 高级 OOP | `ClassGeneric`, `NewAnonymousClass` |
| 生成器 | `Yield`, `YieldFrom`, `GeneratorClass`, `FunctionYield` |
| 其他 | `Goto`, `LabelStatement`, `VarVar`, `Like` |

### P2 — 第三批

| 分类 | 节点类型 |
|------|----------|
| PHP 超全局 | `GlobalsGetVariable`, `GlobalsPostVariable`, `GlobalsServerVariable` 等 |
| HTML 模板 | `Html`, `HtmlAttrs`, `InlineHTML`, `HtmlDocType` |
| 特殊 | `Spawn`, `JsServer`, `Hooks`, `Todo` |

**未覆盖的节点类型**：生成 `nil` + `// TODO` 注释，运行时自动回退到解析器处理。

## 动态特性处理

| 场景 | 处理方式 |
|------|----------|
| eval($code) | AST 中保留 eval 节点，运行时由解释器动态解析执行 |
| $$x 变量变量 | AST 中保留 VarVar 节点，运行时动态解析变量名 |
| call_user_func() | AST 中保留调用节点，运行时动态分派 |
| 动态 include ($path) | 文件路径是变量，无法预编译 → LoadAndRun 未命中 → 回退解析 |
| new $class() | 类名是变量，运行时动态查找 |
| 反射 (ReflectionClass) | 正常工作，类/方法定义已通过 AST 注册到 VM |
| interface/trait 继承 | AST 中保留完整继承关系，注册时按顺序执行 |
| 全局变量 global | 变量索引在 AST 中已固化，RegisterGlobalContext 正确处理 |
| 常量 define() | AST 中保留 const 节点，执行时注册到 VM |
| 闭包 use($var) | AST 中保留捕获变量绑定，运行时正确闭包 |

## 路径匹配策略

生成时用相对路径（相对于项目根目录），不依赖绝对路径：

```go
// 生成时
vm.RegisterCompiledFile("vendor/monolog/src/Logger.php", fn)

// LoadAndRun 时统一转为相对路径再匹配
func (vm *VM) LoadAndRun(file string) (data.GetValue, data.Control) {
    relFile := toRelative(file)
    if fn, ok := vm.compiledFiles[relFile]; ok {
        // 命中预编译
    }
    // 回退解析
}
```

## 错误处理

| 场景 | 处理方式 |
|------|----------|
| vendor 目录不存在 | 报错退出：`vendor/ 目录不存在，请先运行 composer install` |
| PHP 文件解析失败 | 跳过该文件，打印警告，继续处理其他文件 |
| 不支持的节点类型 | 生成 `nil` + `// TODO: unsupported` 注释，运行时回退 |
| Go 代码生成失败 | 报错退出，提示具体文件和节点类型 |
| 生成的 Go 代码编译失败 | 用户手动检查 `.zy/build/` 中的代码 |

## 测试策略

| 层级 | 方式 |
|------|------|
| 生成器单元测试 | 每种节点类型一个测试：构造 AST → 生成 Go 代码 → 验证输出格式 |
| 集成测试 | 准备小型 vendor 目录 → `zy compile` → `go build` → 运行验证 |
| 回归测试 | 现有 `zy phpt` 测试套件，验证预编译后行为不变 |
| 性能测试 | 对比有/无预编译的启动时间 |

## 实现阶段

1. **阶段 1**：P0 节点类型 + VM 注册表 + 基本 CLI
2. **阶段 2**：P1 节点类型 + autoload 命名空间扫描
3. **阶段 3**：P2 节点类型 + 完整覆盖
4. **阶段 4**：优化（常量折叠、死代码消除）
