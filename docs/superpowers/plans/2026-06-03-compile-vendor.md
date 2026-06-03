# `zy compile` — Vendor AOT 编译实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 `vendor/` 目录下的 PHP 文件预编译为 Go 源码（AST 结构体字面量），消除运行时对 vendor 代码的词法分析和解析开销。

**Architecture:** 构建时扫描 vendor 目录，逐个解析 PHP 文件为 AST，然后为每个 AST 节点生成对应的 Go 构造代码。生成的 Go 包通过 `Register(vm)` 函数将预构建的 AST 注册到 VM 的 `compiledFiles` 注册表中。运行时 `LoadAndRun` 优先查找注册表，命中则直接执行 AST，未命中则回退到原有解析。

**Tech Stack:** Go, cobra (CLI), origami parser/runtime

---

## 文件结构

| 操作 | 文件路径 | 职责 |
|------|----------|------|
| 修改 | `data/context.go` | VM 接口添加 `RegisterCompiledFile` 方法 |
| 修改 | `runtime/vm.go` | VM 实现 `compiledFiles` 注册表和 `RegisterCompiledFile`，修改 `LoadAndRun` |
| 修改 | `cmd/root.go` | 注册 `compileCmd` 子命令，添加到 `IsDirectScriptArg` 排除列表 |
| 修改 | `node/from.go` | 添加 `NewFrom(path string)` 辅助函数，供生成的代码创建 `TokenFrom` |
| 创建 | `cmd/compile.go` | `zy compile` 子命令入口，参数解析，调用编译流程 |
| 创建 | `cmd/compile_scan.go` | 扫描 vendor 目录收集所有 .php 文件 |
| 创建 | `cmd/compile_parse.go` | 解析所有 .php 文件为 AST |
| 创建 | `cmd/compile_gen.go` | AST→Go 代码生成器核心（Generator 结构体、分派函数、输出管理） |
| 创建 | `cmd/compile_gen_literal.go` | 字面量节点生成（IntLiteral、FloatLiteral、StringLiteral、BooleanLiteral、NullLiteral） |
| 创建 | `cmd/compile_gen_variable.go` | 变量节点生成（VariableExpression、VariableReference、ValueReference） |
| 创建 | `cmd/compile_gen_binary.go` | 二元运算节点生成（所有 Binary* 类型） |
| 创建 | `cmd/compile_gen_unary.go` | 一元/后缀运算节点生成 |
| 创建 | `cmd/compile_gen_control.go` | 控制流节点生成（If、For、While、Foreach、Switch、Match、Try 等） |
| 创建 | `cmd/compile_gen_call.go` | 函数/方法调用节点生成 |
| 创建 | `cmd/compile_gen_oop.go` | OOP 节点生成（Class、Interface、Trait、Enum、New 等） |
| 创建 | `cmd/compile_gen_other.go` | 其他节点生成（Array、Index、Ternary、Closure、Include 等） |
| 创建 | `cmd/compile_output.go` | 生成最终的 register.go 和 go.mod 文件 |

---

## Task 1: VM 注册表 — 添加 compiledFiles 支持

**Files:**
- Modify: `data/context.go:44-82` (VM 接口)
- Modify: `runtime/vm.go` (VM 实现)

### Step 1: 在 VM 接口添加 RegisterCompiledFile 方法

在 `data/context.go` 的 `VM` 接口中，在 `LoadAndRun` 方法后添加：

```go
// RegisterCompiledFile 注册预编译的文件 AST，LoadAndRun 时优先使用
RegisterCompiledFile(file string, fn func() (GetValue, []Variable))
```

### Step 2: 在 VM 实现中添加 compiledFiles 字段

在 `runtime/vm.go` 的 `VM` 结构体中添加字段：

```go
// 预编译文件注册表
compiledFiles map[string]func() (data.GetValue, []data.Variable)
```

### Step 3: 在 NewVM 中初始化 compiledFiles

在 `runtime/vm.go` 的 `NewVM` 函数中，初始化 map：

```go
compiledFiles: make(map[string]func() (data.GetValue, []data.Variable)),
```

### Step 4: 实现 RegisterCompiledFile 方法

在 `runtime/vm.go` 中添加方法：

```go
func (vm *VM) RegisterCompiledFile(file string, fn func() (data.GetValue, []data.Variable)) {
    vm.compiledFiles[file] = fn
}
```

### Step 5: 修改 LoadAndRun 支持预编译文件

修改 `runtime/vm.go` 的 `LoadAndRun` 方法，在最前面添加预编译检查：

```go
func (vm *VM) LoadAndRun(file string) (data.GetValue, data.Control) {
    // 检查预编译注册表
    if fn, ok := vm.compiledFiles[file]; ok {
        program, vars := fn()
        ctx := vm.CreateContext(vars)
        vm.RegisterGlobalContext(vars, ctx)
        result, ctrl := program.(data.GetValue).GetValue(ctx)
        if data.FlushAllBuffersFn != nil {
            data.FlushAllBuffersFn()
        }
        return result, ctrl
    }

    // 原有逻辑
    data.ResetUserOutput()
    p := vm.parser.Clone()
    // ... 保持不变 ...
}
```

### Step 6: 验证编译通过

```bash
cd D:/github.cocm/php-any/origami && go build ./...
```

### Step 7: 提交

```bash
git add data/context.go runtime/vm.go
git commit -m "feat: add compiledFiles registry to VM for AOT support"
```

---

## Task 2: CLI 子命令骨架

**Files:**
- Create: `cmd/compile.go`
- Modify: `cmd/root.go:38` (IsDirectScriptArg 排除列表)
- Modify: `cmd/root.go:60-63` (init 注册子命令)

### Step 1: 创建 cmd/compile.go 基本结构

```go
package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
    Use:   "compile [directory]",
    Short: "将 vendor 目录预编译为 Go 源码",
    Long: `将 vendor/ 目录下的 PHP 文件解析为 AST，
并生成对应的 Go 结构体字面量代码。

生成的 Go 包可以在运行时注册到 VM，
使 LoadAndRun 跳过词法分析和解析阶段，
直接执行预构建的 AST。

示例:
  zy compile vendor/
  zy compile vendor/ -o .zy/build
  zy compile vendor/ --pkg myapp`,
    SilenceUsage: true,
    Args:         cobra.ExactArgs(1),
    RunE:         runCompileCommand,
}

var (
    compileOutput string
    compilePkg    string
)

func init() {
    compileCmd.Flags().StringVarP(&compileOutput, "output", "o", ".zy/build", "输出目录")
    compileCmd.Flags().StringVar(&compilePkg, "pkg", "build", "生成的 Go 包名")
}

func runCompileCommand(cmd *cobra.Command, args []string) error {
    vendorDir := args[0]

    // 检查目录是否存在
    info, err := os.Stat(vendorDir)
    if err != nil {
        return fmt.Errorf("目录不存在: %s", vendorDir)
    }
    if !info.IsDir() {
        return fmt.Errorf("不是目录: %s", vendorDir)
    }

    fmt.Printf("扫描目录: %s\n", vendorDir)
    fmt.Printf("输出目录: %s\n", compileOutput)
    fmt.Printf("包名: %s\n", compilePkg)

    // TODO: 后续任务实现
    return nil
}
```

### Step 2: 在 root.go 注册子命令

在 `cmd/root.go` 的 `init()` 中添加：

```go
rootCmd.AddCommand(compileCmd)
```

### Step 3: 在 IsDirectScriptArg 中添加排除

在 `cmd/root.go:38` 的 switch 中添加 `"compile"`:

```go
case "gen-std", "help", "completion", "phpt", "compile":
    return false
```

### Step 4: 验证编译和帮助输出

```bash
cd D:/github.cocm/php-any/origami && go build ./... && ./origami compile --help
```

### Step 5: 提交

```bash
git add cmd/compile.go cmd/root.go
git commit -m "feat: add zy compile subcommand skeleton"
```

---

## Task 3: 文件扫描器

**Files:**
- Create: `cmd/compile_scan.go`

### Step 1: 实现文件扫描器

```go
package cmd

import (
    "os"
    "path/filepath"
    "strings"
)

// collectPhpFiles 扫描目录下所有 .php 文件
func collectPhpFiles(dir string) ([]string, error) {
    var files []string
    err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }
        if d.IsDir() {
            return nil
        }
        if strings.HasSuffix(strings.ToLower(path), ".php") {
            files = append(files, path)
        }
        return nil
    })
    if err != nil {
        return nil, err
    }
    return files, nil
}
```

### Step 2: 在 runCompileCommand 中调用扫描

修改 `cmd/compile.go` 的 `runCompileCommand`：

```go
func runCompileCommand(cmd *cobra.Command, args []string) error {
    vendorDir := args[0]

    info, err := os.Stat(vendorDir)
    if err != nil {
        return fmt.Errorf("目录不存在: %s", vendorDir)
    }
    if !info.IsDir() {
        return fmt.Errorf("不是目录: %s", vendorDir)
    }

    files, err := collectPhpFiles(vendorDir)
    if err != nil {
        return fmt.Errorf("扫描失败: %w", err)
    }
    if len(files) == 0 {
        return fmt.Errorf("未找到 .php 文件: %s", vendorDir)
    }

    fmt.Printf("找到 %d 个 PHP 文件\n", len(files))
    for _, f := range files {
        fmt.Println(f)
    }
    return nil
}
```

### Step 3: 验证

```bash
cd D:/github.cocm/php-any/origami && go build ./... && ./origami compile vendor/
```

### Step 4: 提交

```bash
git add cmd/compile_scan.go cmd/compile.go
git commit -m "feat: add vendor PHP file scanner"
```

---

## Task 4: AST 解析器

**Files:**
- Create: `cmd/compile_parse.go`

### Step 1: 实现批量解析

```go
package cmd

import (
    "fmt"

    "github.com/php-any/origami/node"
    "github.com/php-any/origami/parser"
)

// ParsedFile 表示一个已解析的 PHP 文件
type ParsedFile struct {
    Path     string
    Program  *node.Program
    Variables []parser.Variable
}

// parseFiles 批量解析 PHP 文件为 AST
func parseFiles(files []string) ([]ParsedFile, []error) {
    p := parser.NewParser()
    var parsed []ParsedFile
    var errs []error

    for _, file := range files {
        clone := p.Clone()
        program, acl := clone.ParseFile(file)
        if acl != nil {
            errs = append(errs, fmt.Errorf("解析 %s 失败: %v", file, acl))
            continue
        }
        vars := clone.GetVariables()
        parsed = append(parsed, ParsedFile{
            Path:      file,
            Program:   program,
            Variables: vars,
        })
    }
    return parsed, errs
}
```

### Step 2: 需要确认 parser.Variable 类型

检查 `parser/parser.go` 中 `GetVariables()` 的返回类型，确认是 `[]data.Variable` 还是 `[]parser.Variable`。根据之前的探索，它返回 `[]data.Variable`。

修正 `ParsedFile` 的 Variables 字段类型为 `[]data.Variable`：

```go
type ParsedFile struct {
    Path      string
    Program   *node.Program
    Variables []data.Variable
}
```

### Step 3: 在 runCompileCommand 中集成解析

```go
func runCompileCommand(cmd *cobra.Command, args []string) error {
    vendorDir := args[0]

    info, err := os.Stat(vendorDir)
    if err != nil {
        return fmt.Errorf("目录不存在: %s", vendorDir)
    }
    if !info.IsDir() {
        return fmt.Errorf("不是目录: %s", vendorDir)
    }

    files, err := collectPhpFiles(vendorDir)
    if err != nil {
        return fmt.Errorf("扫描失败: %w", err)
    }
    if len(files) == 0 {
        return fmt.Errorf("未找到 .php 文件: %s", vendorDir)
    }

    fmt.Printf("找到 %d 个 PHP 文件，开始解析...\n", len(files))

    parsed, errs := parseFiles(files)
    if len(errs) > 0 {
        for _, e := range errs {
            fmt.Fprintf(os.Stderr, "警告: %v\n", e)
        }
    }
    if len(parsed) == 0 {
        return fmt.Errorf("没有文件解析成功")
    }

    fmt.Printf("成功解析 %d 个文件\n", len(parsed))
    return nil
}
```

### Step 4: 验证

```bash
cd D:/github.cocm/php-any/origami && go build ./... && ./origami compile vendor/
```

### Step 5: 提交

```bash
git add cmd/compile_parse.go cmd/compile.go
git commit -m "feat: add batch PHP file parser for compile command"
```

---

## Task 5: 代码生成器核心框架

**Files:**
- Create: `cmd/compile_gen.go`

### Step 1: 创建 Generator 结构体和核心方法

```go
package cmd

import (
    "fmt"
    "strings"

    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

// Generator 将 AST 节点转换为 Go 构造代码
type Generator struct {
    buf        strings.Builder
    indent     int
    imports    map[string]bool
    file       string // 当前处理的文件路径
}

// NewGenerator 创建新的代码生成器
func NewGenerator() *Generator {
    return &Generator{
        imports: make(map[string]bool),
    }
}

// Generate 为一个解析后的文件生成 Go 代码
func (g *Generator) Generate(pf ParsedFile) string {
    g.file = pf.Path
    g.buf.Reset()
    g.imports = make(map[string]bool)
    g.indent = 0

    funcName := g.funcNameForPath(pf.Path)

    g.printf("func %s() (data.GetValue, []data.Variable) {\n", funcName)
    g.indent++
    g.printf("from := &node.TokenFrom{}\n")
    g.printf("_ = from\n")
    g.printf("\n")
    g.printf("stmts := []data.GetValue{\n")
    g.indent++
    for _, stmt := range pf.Program.Statements {
        g.genGetValue(stmt)
        g.printf(",\n")
    }
    g.indent--
    g.printf("}\n")
    g.printf("\n")

    // 生成变量列表
    g.printf("vars := []data.Variable{\n")
    g.indent++
    for _, v := range pf.Variables {
        g.printf("{Index: %d, Name: %q},\n", v.GetIndex(), v.GetName())
    }
    g.indent--
    g.printf("}\n")
    g.printf("\n")
    g.printf("return node.NewProgram(from, stmts), vars\n")
    g.indent--
    g.printf("}\n")

    return g.buf.String()
}

// funcNameForPath 将文件路径转换为合法的 Go 函数名
func (g *Generator) funcNameForPath(path string) string {
    // 移除扩展名
    name := strings.TrimSuffix(path, ".php")
    name = strings.TrimSuffix(name, ".zy")
    // 替换非字母数字字符为下划线
    var b strings.Builder
    for _, c := range name {
        if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
            b.WriteRune(c)
        } else {
            b.WriteByte('_')
        }
    }
    // 确保以大写字母开头（导出函数）
    result := b.String()
    if len(result) > 0 && result[0] >= 'a' && result[0] <= 'z' {
        result = string(rune(result[0]-32)) + result[1:]
    }
    return "AST_" + result
}

// genGetValue 根据节点类型分派到对应的生成函数
func (g *Generator) genGetValue(v data.GetValue) {
    if v == nil {
        g.printf("nil")
        return
    }
    switch n := v.(type) {
    // 字面量
    case *node.IntLiteral:
        g.genIntLiteral(n)
    case *node.FloatLiteral:
        g.genFloatLiteral(n)
    case *node.StringLiteral:
        g.genStringLiteral(n)
    case *node.BooleanLiteral:
        g.genBooleanLiteral(n)
    case *node.NullLiteral:
        g.genNullLiteral(n)
    // 变量
    case *node.VariableExpression:
        g.genVariableExpression(n)
    case *node.VariableReference:
        g.genVariableReference(n)
    case *node.ValueReference:
        g.genValueReference(n)
    // 二元运算
    case *node.BinaryAdd:
        g.genBinaryOp("BinaryAdd", n.Left, n.Right)
    case *node.BinarySub:
        g.genBinaryOp("BinarySub", n.Left, n.Right)
    case *node.BinaryMul:
        g.genBinaryOp("BinaryMul", n.Left, n.Right)
    case *node.BinaryQuo:
        g.genBinaryOp("BinaryQuo", n.Left, n.Right)
    case *node.BinaryRem:
        g.genBinaryOp("BinaryRem", n.Left, n.Right)
    case *node.BinaryPow:
        g.genBinaryOp("BinaryPow", n.Left, n.Right)
    case *node.BinaryDot:
        g.genBinaryOp("BinaryDot", n.Left, n.Right)
    case *node.BinaryEq:
        g.genBinaryOp("BinaryEq", n.Left, n.Right)
    case *node.BinaryNe:
        g.genBinaryOp("BinaryNe", n.Left, n.Right)
    case *node.BinaryEqStrict:
        g.genBinaryOp("BinaryEqStrict", n.Left, n.Right)
    case *node.BinaryNeStrict:
        g.genBinaryOp("BinaryNeStrict", n.Left, n.Right)
    case *node.BinaryLt:
        g.genBinaryOp("BinaryLt", n.Left, n.Right)
    case *node.BinaryLe:
        g.genBinaryOp("BinaryLe", n.Left, n.Right)
    case *node.BinaryGt:
        g.genBinaryOp("BinaryGt", n.Left, n.Right)
    case *node.BinaryGe:
        g.genBinaryOp("BinaryGe", n.Left, n.Right)
    case *node.BinaryLand:
        g.genBinaryOp("BinaryLand", n.Left, n.Right)
    case *node.BinaryLor:
        g.genBinaryOp("BinaryLor", n.Left, n.Right)
    case *node.BinarySpaceship:
        g.genBinaryOp("BinarySpaceship", n.Left, n.Right)
    // 控制流
    case *node.IfStatement:
        g.genIfStatement(n)
    case *node.ReturnStatement:
        g.genReturnStatement(n)
    case *node.EchoStatement:
        g.genEchoStatement(n)
    // 未支持的类型
    default:
        g.printf("nil /* TODO: unsupported %T */", v)
    }
}

// genBinaryOp 生成二元运算的通用方法
func (g *Generator) genBinaryOp(typeName string, left, right data.GetValue) {
    g.printf("node.New%s(from,\n", typeName)
    g.indent++
    g.genGetValue(left)
    g.printf(",\n")
    g.genGetValue(right)
    g.printf(",\n")
    g.indent--
    g.printf(")")
}

// printf 带缩进的格式化输出
func (g *Generator) printf(format string, args ...interface{}) {
    // 在换行后添加缩进
    msg := fmt.Sprintf(format, args...)
    lines := strings.Split(msg, "\n")
    for i, line := range lines {
        if i > 0 {
            g.buf.WriteString("\n")
            if len(line) > 0 {
                for j := 0; j < g.indent; j++ {
                    g.buf.WriteString("\t")
                }
            }
        }
        g.buf.WriteString(line)
    }
}
```

### Step 2: 验证编译

```bash
cd D:/github.cocm/php-any/origami && go build ./...
```

### Step 3: 提交

```bash
git add cmd/compile_gen.go
git commit -m "feat: add AST-to-Go code generator core framework"
```

---

## Task 6: 字面量节点生成

**Files:**
- Create: `cmd/compile_gen_literal.go`

### Step 1: 实现字面量生成器

```go
package cmd

import (
    "fmt"
    "strconv"

    "github.com/php-any/origami/node"
)

func (g *Generator) genIntLiteral(n *node.IntLiteral) {
    // IntLiteral 存储的是 data.Value，需要提取 int 值
    if n.V != nil {
        s := n.V.AsString()
        g.printf("&node.IntLiteral{Node: node.NewNode(from), V: data.NewIntValue(%s)}", s)
    } else {
        g.printf("&node.IntLiteral{Node: node.NewNode(from), V: data.NewIntValue(0)}")
    }
}

func (g *Generator) genFloatLiteral(n *node.FloatLiteral) {
    if n.V != nil {
        s := n.V.AsString()
        g.printf("&node.FloatLiteral{Node: node.NewNode(from), V: data.NewFloatValue(%s)}", s)
    } else {
        g.printf("&node.FloatLiteral{Node: node.NewNode(from), V: data.NewFloatValue(0)}")
    }
}

func (g *Generator) genStringLiteral(n *node.StringLiteral) {
    g.printf("&node.StringLiteral{Node: node.NewNode(from), Value: %s}", strconv.Quote(n.Value))
}

func (g *Generator) genBooleanLiteral(n *node.BooleanLiteral) {
    if n.Value {
        g.printf("&node.BooleanLiteral{Node: node.NewNode(from), Value: true}")
    } else {
        g.printf("&node.BooleanLiteral{Node: node.NewNode(from), Value: false}")
    }
}

func (g *Generator) genNullLiteral(n *node.NullLiteral) {
    g.printf("&node.NullLiteral{Node: node.NewNode(from)}")
}
```

### Step 2: 需要确认字面量结构体字段

检查 `node/float_literal.go`、`node/boolean_literal.go`、`node/null_literal.go` 确认字段名。根据探索结果：
- `FloatLiteral`: `*Node`, `V data.Value`
- `BooleanLiteral`: 需要确认
- `NullLiteral`: 需要确认

读取这些文件确认字段，然后修正代码。

### Step 3: 验证编译

```bash
cd D:/github.cocm/php-any/origami && go build ./...
```

### Step 4: 提交

```bash
git add cmd/compile_gen_literal.go
git commit -m "feat: add literal node code generators"
```

---

## Task 7: 变量节点生成

**Files:**
- Create: `cmd/compile_gen_variable.go`

### Step 1: 实现变量生成器

```go
package cmd

import (
    "fmt"

    "github.com/php-any/origami/node"
)

func (g *Generator) genVariableExpression(n *node.VariableExpression) {
    g.printf("node.NewVariable(from, %q, %d, nil)", n.Name, n.Index)
}

func (g *Generator) genVariableReference(n *node.VariableReference) {
    g.printf("node.NewVariableReference(from, %q, %d, nil)", n.Name, n.Index)
}

func (g *Generator) genValueReference(n *node.ValueReference) {
    // ValueReference 包装一个变量表达式
    g.printf("node.NewValueReference(from, ")
    if n.Variable != nil {
        g.genGetValue(n.Variable)
    } else {
        g.printf("nil")
    }
    g.printf(")")
}
```

### Step 2: 确认 ValueReference 结构

读取 `node/value_reference.go` 确认字段名和构造函数签名。

### Step 3: 验证编译

```bash
cd D:/github.cocm/php-any/origami && go build ./...
```

### Step 4: 提交

```bash
git add cmd/compile_gen_variable.go
git commit -m "feat: add variable node code generators"
```

---

## Task 8: 二元运算节点生成（完善）

**Files:**
- Modify: `cmd/compile_gen.go` (genGetValue 中的 Binary* 分支)
- Create: `cmd/compile_gen_binary.go`

### Step 1: 将二元运算生成移到独立文件

将 `genBinaryOp` 方法和所有 Binary* 分支移到 `cmd/compile_gen_binary.go`。添加剩余的二元运算类型：

```go
package cmd

import (
    "github.com/php-any/origami/data"
    "github.com/php-any/origami/node"
)

func (g *Generator) genBinaryOp(typeName string, left, right data.GetValue) {
    g.printf("node.New%s(from,\n", typeName)
    g.indent++
    g.genGetValue(left)
    g.printf(",\n")
    g.genGetValue(right)
    g.printf(",\n")
    g.indent--
    g.printf(")")
}

func (g *Generator) genBinaryAssign(n *node.BinaryAssign) {
    g.printf("node.NewBinaryAssign(from,\n")
    g.indent++
    g.genGetValue(n.Left)
    g.printf(",\n")
    g.genGetValue(n.Right)
    g.printf(",\n")
    g.indent--
    g.printf(")")
}

func (g *Generator) genBinaryAssignVariable(n *node.BinaryAssignVariable) {
    g.printf("&node.BinaryAssignVariable{\n")
    g.indent++
    g.printf("Node: node.NewNode(from),\n")
    g.printf("Left: ")
    g.genGetValue(n.Left)
    g.printf(",\nRight: ")
    g.genGetValue(n.Right)
    g.printf(",\n")
    g.indent--
    g.printf("}")
}

func (g *Generator) genBinaryBitwise(typeName string, left, right data.GetValue) {
    g.printf("node.New%s(from,\n", typeName)
    g.indent++
    g.genGetValue(left)
    g.printf(",\n")
    g.genGetValue(right)
    g.printf(",\n")
    g.indent--
    g.printf(")")
}

func (g *Generator) genBinaryShift(typeName string, left, right data.GetValue) {
    g.printf("node.New%s(from,\n", typeName)
    g.indent++
    g.genGetValue(left)
    g.printf(",\n")
    g.genGetValue(right)
    g.printf(",\n")
    g.indent--
    g.printf(")")
}
```

### Step 2: 在 genGetValue 中添加所有 Binary* 分支

在 `cmd/compile_gen.go` 的 `genGetValue` 中添加：

```go
case *node.BinaryAssign:
    g.genBinaryAssign(n)
case *node.BinaryAssignVariable:
    g.genBinaryAssignVariable(n)
case *node.BinaryBitAnd:
    g.genBinaryBitwise("BinaryBitAnd", n.Left, n.Right)
case *node.BinaryBitOr:
    g.genBinaryBitwise("BinaryBitOr", n.Left, n.Right)
case *node.BinaryBitXor:
    g.genBinaryBitwise("BinaryBitXor", n.Left, n.Right)
case *node.BinaryShl:
    g.genBinaryShift("BinaryShl", n.Left, n.Right)
case *node.BinaryShr:
    g.genBinaryShift("BinaryShr", n.Left, n.Right)
case *node.BinaryLink:
    g.genBinaryOp("BinaryLink", n.Left, n.Right)
```

### Step 3: 验证编译

```bash
cd D:/github.cocm/php-any/origami && go build ./...
```

### Step 4: 提交

```bash
git add cmd/compile_gen_binary.go cmd/compile_gen.go
git commit -m "feat: add all binary operation code generators"
```

---

## Task 9: 一元/后缀运算节点生成

**Files:**
- Create: `cmd/compile_gen_unary.go`

### Step 1: 确认一元运算结构体字段

读取以下文件确认字段：
- `node/unary_incr.go`, `node/unary_decr.go`
- `node/postfix_incr.go`, `node/postfix_decr.go`
- `node/error_suppress.go`
- `node/expression.go` (UnaryExpression)

### Step 2: 实现一元运算生成器

```go
package cmd

import (
    "github.com/php-any/origami/node"
)

func (g *Generator) genUnaryIncr(n *node.UnaryIncr) {
    g.printf("&node.UnaryIncr{Node: node.NewNode(from), Value: ")
    g.genGetValue(n.Value)
    g.printf("}")
}

func (g *Generator) genUnaryDecr(n *node.UnaryDecr) {
    g.printf("&node.UnaryDecr{Node: node.NewNode(from), Value: ")
    g.genGetValue(n.Value)
    g.printf("}")
}

func (g *Generator) genPostfixIncr(n *node.PostfixIncr) {
    g.printf("&node.PostfixIncr{Node: node.NewNode(from), Value: ")
    g.genGetValue(n.Value)
    g.printf("}")
}

func (g *Generator) genPostfixDecr(n *node.PostfixDecr) {
    g.printf("&node.PostfixDecr{Node: node.NewNode(from), Value: ")
    g.genGetValue(n.Value)
    g.printf("}")
}

func (g *Generator) genErrorSuppress(n *node.ErrorSuppress) {
    g.printf("&node.ErrorSuppress{Node: node.NewNode(from), Value: ")
    g.genGetValue(n.Value)
    g.printf("}")
}

func (g *Generator) genUnaryExpression(n *node.UnaryExpression) {
    g.printf("&node.UnaryExpression{Node: node.NewNode(from), Op: %q, Value: ", n.Op)
    g.genGetValue(n.Value)
    g.printf("}")
}
```

### Step 3: 在 genGetValue 中添加分支

```go
case *node.UnaryIncr:
    g.genUnaryIncr(n)
case *node.UnaryDecr:
    g.genUnaryDecr(n)
case *node.PostfixIncr:
    g.genPostfixIncr(n)
case *node.PostfixDecr:
    g.genPostfixDecr(n)
case *node.ErrorSuppress:
    g.genErrorSuppress(n)
case *node.UnaryExpression:
    g.genUnaryExpression(n)
```

### Step 4: 验证编译

```bash
cd D:/github.cocm/php-any/origami && go build ./...
```

### Step 5: 提交

```bash
git add cmd/compile_gen_unary.go cmd/compile_gen.go
git commit -m "feat: add unary/postfix operation code generators"
```

---

## Task 10: 控制流节点生成

**Files:**
- Create: `cmd/compile_gen_control.go`

### Step 1: 确认控制流结构体字段

读取以下文件确认字段：
- `node/if.go` — IfStatement, ElseIfBranch
- `node/for.go` — ForStatement
- `node/foreach.go` — ForeachStatement
- `node/while.go` — WhileStatement
- `node/do_while.go` — DoWhileStatement
- `node/switch.go` — SwitchStatement, SwitchCase
- `node/match.go` — MatchStatement, MatchArm
- `node/try.go` — TryStatement, CatchBlock
- `node/return.go` — ReturnStatement
- `node/break.go` — BreakStatement
- `node/continue.go` — ContinueStatement
- `node/throw.go` — ThrowStatement

### Step 2: 实现控制流生成器

```go
package cmd

import "github.com/php-any/origami/node"

func (g *Generator) genIfStatement(n *node.IfStatement) {
    g.printf("node.NewIfStatement(from,\n")
    g.indent++
    // Condition
    g.genGetValue(n.Condition)
    g.printf(",\n")
    // ThenBranch
    g.printf("[]data.GetValue{\n")
    g.indent++
    for _, stmt := range n.ThenBranch {
        g.genGetValue(stmt)
        g.printf(",\n")
    }
    g.indent--
    g.printf("},\n")
    // ElseIf
    g.printf("[]node.ElseIfBranch{\n")
    g.indent++
    for _, elif := range n.ElseIf {
        g.printf("{Condition: ")
        g.genGetValue(elif.Condition)
        g.printf(", ThenBranch: []data.GetValue{\n")
        g.indent++
        for _, stmt := range elif.ThenBranch {
            g.genGetValue(stmt)
            g.printf(",\n")
        }
        g.indent--
        g.printf("}},\n")
    }
    g.indent--
    g.printf("},\n")
    // ElseBranch
    g.printf("[]data.GetValue{\n")
    g.indent++
    for _, stmt := range n.ElseBranch {
        g.genGetValue(stmt)
        g.printf(",\n")
    }
    g.indent--
    g.printf("},\n")
    g.indent--
    g.printf(")")
}

func (g *Generator) genReturnStatement(n *node.ReturnStatement) {
    g.printf("&node.ReturnStatement{Node: node.NewNode(from), Value: ")
    if n.Value != nil {
        g.genGetValue(n.Value)
    } else {
        g.printf("nil")
    }
    g.printf("}")
}

func (g *Generator) genEchoStatement(n *node.EchoStatement) {
    g.printf("node.NewEchoStatement(from, []data.GetValue{\n")
    g.indent++
    for _, expr := range n.Expressions {
        g.genGetValue(expr)
        g.printf(",\n")
    }
    g.indent--
    g.printf("})")
}

func (g *Generator) genBreakStatement(n *node.BreakStatement) {
    g.printf("&node.BreakStatement{Node: node.NewNode(from)")
    if n.Level != nil {
        g.printf(", Level: ")
        g.genGetValue(n.Level)
    }
    g.printf("}")
}

func (g *Generator) genContinueStatement(n *node.ContinueStatement) {
    g.printf("&node.ContinueStatement{Node: node.NewNode(from)")
    if n.Level != nil {
        g.printf(", Level: ")
        g.genGetValue(n.Level)
    }
    g.printf("}")
}

func (g *Generator) genThrowStatement(n *node.ThrowStatement) {
    g.printf("&node.ThrowStatement{Node: node.NewNode(from), Value: ")
    g.genGetValue(n.Value)
    g.printf("}")
}

func (g *Generator) genForStatement(n *node.ForStatement) {
    g.printf("&node.ForStatement{\n")
    g.indent++
    g.printf("Node: node.NewNode(from),\n")
    g.printf("Init: ")
    g.genGetValue(n.Init)
    g.printf(",\nCondition: ")
    g.genGetValue(n.Condition)
    g.printf(",\nUpdate: ")
    g.genGetValue(n.Update)
    g.printf(",\nBody: []data.GetValue{\n")
    g.indent++
    for _, stmt := range n.Body {
        g.genGetValue(stmt)
        g.printf(",\n")
    }
    g.indent--
    g.printf("},\n")
    g.indent--
    g.printf("}")
}

func (g *Generator) genWhileStatement(n *node.WhileStatement) {
    g.printf("&node.WhileStatement{\n")
    g.indent++
    g.printf("Node: node.NewNode(from),\n")
    g.printf("Condition: ")
    g.genGetValue(n.Condition)
    g.printf(",\nBody: []data.GetValue{\n")
    g.indent++
    for _, stmt := range n.Body {
        g.genGetValue(stmt)
        g.printf(",\n")
    }
    g.indent--
    g.printf("},\n")
    g.indent--
    g.printf("}")
}

func (g *Generator) genForeachStatement(n *node.ForeachStatement) {
    g.printf("&node.ForeachStatement{\n")
    g.indent++
    g.printf("Node: node.NewNode(from),\n")
    g.printf("Source: ")
    g.genGetValue(n.Source)
    g.printf(",\n")
    if n.Key != nil {
        g.printf("Key: ")
        g.genGetValue(n.Key)
        g.printf(",\n")
    }
    g.printf("Value: ")
    g.genGetValue(n.Value)
    g.printf(",\nBody: []data.GetValue{\n")
    g.indent++
    for _, stmt := range n.Body {
        g.genGetValue(stmt)
        g.printf(",\n")
    }
    g.indent--
    g.printf("},\n")
    g.indent--
    g.printf("}")
}

func (g *Generator) genDoWhileStatement(n *node.DoWhileStatement) {
    g.printf("&node.DoWhileStatement{\n")
    g.indent++
    g.printf("Node: node.NewNode(from),\n")
    g.printf("Condition: ")
    g.genGetValue(n.Condition)
    g.printf(",\nBody: []data.GetValue{\n")
    g.indent++
    for _, stmt := range n.Body {
        g.genGetValue(stmt)
        g.printf(",\n")
    }
    g.indent--
    g.printf("},\n")
    g.indent--
    g.printf("}")
}

func (g *Generator) genTryStatement(n *node.TryStatement) {
    g.printf("&node.TryStatement{\n")
    g.indent++
    g.printf("Node: node.NewNode(from),\n")
    g.printf("Body: []data.GetValue{\n")
    g.indent++
    for _, stmt := range n.Body {
        g.genGetValue(stmt)
        g.printf(",\n")
    }
    g.indent--
    g.printf("},\nCatches: []node.CatchBlock{\n")
    g.indent++
    for _, c := range n.Catches {
        g.printf("{Types: %q, Variable: %q, Body: []data.GetValue{\n", c.Types, c.Variable)
        g.indent++
        for _, stmt := range c.Body {
            g.genGetValue(stmt)
            g.printf(",\n")
        }
        g.indent--
        g.printf("}},\n")
    }
    g.indent--
    g.printf("},\n")
    if n.Finally != nil {
        g.printf("Finally: []data.GetValue{\n")
        g.indent++
        for _, stmt := range n.Finally {
            g.genGetValue(stmt)
            g.printf(",\n")
        }
        g.indent--
        g.printf("},\n")
    }
    g.indent--
    g.printf("}")
}
```

### Step 3: 在 genGetValue 中添加所有控制流分支

```go
case *node.ReturnStatement:
    g.genReturnStatement(n)
case *node.EchoStatement:
    g.genEchoStatement(n)
case *node.BreakStatement:
    g.genBreakStatement(n)
case *node.ContinueStatement:
    g.genContinueStatement(n)
case *node.ThrowStatement:
    g.genThrowStatement(n)
case *node.ForStatement:
    g.genForStatement(n)
case *node.WhileStatement:
    g.genWhileStatement(n)
case *node.ForeachStatement:
    g.genForeachStatement(n)
case *node.DoWhileStatement:
    g.genDoWhileStatement(n)
case *node.TryStatement:
    g.genTryStatement(n)
case *node.SwitchStatement:
    g.genSwitchStatement(n)
case *node.MatchStatement:
    g.genMatchStatement(n)
```

### Step 4: 验证编译

```bash
cd D:/github.cocm/php-any/origami && go build ./...
```

### Step 5: 提交

```bash
git add cmd/compile_gen_control.go cmd/compile_gen.go
git commit -m "feat: add control flow code generators"
```

---

## Task 11: 函数/方法调用节点生成

**Files:**
- Create: `cmd/compile_gen_call.go`

### Step 1: 确认调用结构体字段

读取以下文件确认字段：
- `node/call.go` — CallExpression
- `node/call_method.go` — CallMethod
- `node/call_static_method.go` — CallStaticMethod
- `node/call_object_method.go` — CallObjectMethod
- `node/call_parent_method.go` — CallParentMethod
- `node/call_self_method.go` — CallSelfMethod
- `node/call_nullsafe.go` — NullsafeCall

### Step 2: 实现调用生成器

```go
package cmd

import "github.com/php-any/origami/node"

func (g *Generator) genCallExpression(n *node.CallExpression) {
    g.printf("node.NewCall(from, %q, []data.GetValue{\n", n.FunName)
    g.indent++
    for _, arg := range n.Args {
        g.genGetValue(arg)
        g.printf(",\n")
    }
    g.indent--
    g.printf("})")
}

func (g *Generator) genCallMethod(n *node.CallMethod) {
    g.printf("&node.CallMethod{\n")
    g.indent++
    g.printf("Node: node.NewNode(from),\n")
    g.printf("Object: ")
    g.genGetValue(n.Object)
    g.printf(",\nMethod: %q,\n", n.Method)
    g.printf("Args: []data.GetValue{\n")
    g.indent++
    for _, arg := range n.Args {
        g.genGetValue(arg)
        g.printf(",\n")
    }
    g.indent--
    g.printf("},\n")
    g.indent--
    g.printf("}")
}

func (g *Generator) genCallStaticMethod(n *node.CallStaticMethod) {
    g.printf("&node.CallStaticMethod{\n")
    g.indent++
    g.printf("Node: node.NewNode(from),\n")
    g.printf("Class: %q,\n", n.Class)
    g.printf("Method: %q,\n", n.Method)
    g.printf("Args: []data.GetValue{\n")
    g.indent++
    for _, arg := range n.Args {
        g.genGetValue(arg)
        g.printf(",\n")
    }
    g.indent--
    g.printf("},\n")
    g.indent--
    g.printf("}")
}
```

### Step 3: 在 genGetValue 中添加调用分支

```go
case *node.CallExpression:
    g.genCallExpression(n)
case *node.CallMethod:
    g.genCallMethod(n)
case *node.CallStaticMethod:
    g.genCallStaticMethod(n)
case *node.CallObjectMethod:
    g.genCallObjectMethod(n)
case *node.CallParentMethod:
    g.genCallParentMethod(n)
case *node.CallSelfMethod:
    g.genCallSelfMethod(n)
case *node.NullsafeCall:
    g.genNullsafeCall(n)
```

### Step 4: 验证编译

```bash
cd D:/github.cocm/php-any/origami && go build ./...
```

### Step 5: 提交

```bash
git add cmd/compile_gen_call.go cmd/compile_gen.go
git commit -m "feat: add function/method call code generators"
```

---

## Task 12: OOP 节点生成

**Files:**
- Create: `cmd/compile_gen_oop.go`

### Step 1: 确认 OOP 结构体字段

读取以下文件确认字段：
- `node/class.go` — ClassStatement, ClassProperty, ClassMethod
- `node/interface.go` — InterfaceStatement
- `node/new.go` — NewExpression
- `node/instanceof.go` — InstanceOfExpression
- `node/clone.go` — CloneExpression
- `node/init_class.go` — InitClass
- `node/class_constant.go` — ClassConstant

### Step 2: 实现 OOP 生成器

根据字段确认结果，为每个 OOP 节点类型实现生成函数。ClassStatement 是最复杂的，需要处理：
- Name, Extends, Implements
- Properties (ClassProperty 列表)
- Methods (ClassMethod 列表)
- Constructor
- Annotations

### Step 3: 在 genGetValue 中添加 OOP 分支

### Step 4: 验证编译

### Step 5: 提交

```bash
git add cmd/compile_gen_oop.go cmd/compile_gen.go
git commit -m "feat: add OOP node code generators"
```

---

## Task 13: 其他节点生成

**Files:**
- Create: `cmd/compile_gen_other.go`

### Step 1: 实现剩余节点类型

包括：
- Array, ArraySpread, Index, KV, Compact
- Ternary, NullCoalesce, Range
- Function, Closure, Lambda
- IncludeStatement, Namespace, Use, Const
- GlobalStatement, VarStatement
- IssetStatement, UnsetStatement
- InstanceOfExpression, LikeExpression
- NamedArgument, SpreadArgument
- Annotation, CallAnn

### Step 2: 在 genGetValue 中添加所有剩余分支

### Step 3: 验证编译

### Step 4: 提交

```bash
git add cmd/compile_gen_other.go cmd/compile_gen.go
git commit -m "feat: add remaining node code generators"
```

---

## Task 14: 输出文件生成

**Files:**
- Create: `cmd/compile_output.go`
- Modify: `cmd/compile.go` (集成完整流程)

### Step 1: 实现输出文件生成

```go
package cmd

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

// generateOutput 生成最终的 Go 源码文件
func generateOutput(parsed []ParsedFile, outputDir, pkgName string) error {
    // 创建输出目录
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return fmt.Errorf("创建输出目录失败: %w", err)
    }

    // 生成 register.go
    if err := generateRegisterFile(parsed, outputDir, pkgName); err != nil {
        return err
    }

    // 生成 vendor_ast.go
    if err := generateASTFile(parsed, outputDir, pkgName); err != nil {
        return err
    }

    // 生成 go.mod
    if err := generateGoMod(outputDir, pkgName); err != nil {
        return err
    }

    return nil
}

func generateRegisterFile(parsed []ParsedFile, outputDir, pkgName string) error {
    var b strings.Builder
    b.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
    b.WriteString("import (\n")
    b.WriteString("\t\"github.com/php-any/origami/data\"\n")
    b.WriteString("\t\"github.com/php-any/origami/node\"\n")
    b.WriteString(")\n\n")
    b.WriteString("// Register 将预编译的 vendor AST 注册到 VM\n")
    b.WriteString("func Register(vm data.VM) {\n")
    for _, pf := range parsed {
        funcName := NewGenerator().funcNameForPath(pf.Path)
        b.WriteString(fmt.Sprintf("\tvm.RegisterCompiledFile(%q, func() (data.GetValue, []data.Variable) {\n", pf.Path))
        b.WriteString(fmt.Sprintf("\t\treturn %s()\n", funcName))
        b.WriteString("\t})\n")
    }
    b.WriteString("}\n")

    return os.WriteFile(filepath.Join(outputDir, "register.go"), []byte(b.String()), 0644)
}

func generateASTFile(parsed []ParsedFile, outputDir, pkgName string) error {
    gen := NewGenerator()
    var b strings.Builder
    b.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
    b.WriteString("import (\n")
    b.WriteString("\t\"github.com/php-any/origami/data\"\n")
    b.WriteString("\t\"github.com/php-any/origami/node\"\n")
    b.WriteString(")\n\n")

    for _, pf := range parsed {
        code := gen.Generate(pf)
        b.WriteString(code)
        b.WriteString("\n\n")
    }

    return os.WriteFile(filepath.Join(outputDir, "vendor_ast.go"), []byte(b.String()), 0644)
}

func generateGoMod(outputDir, pkgName string) error {
    content := fmt.Sprintf("module %s\n\ngo 1.21\n\nrequire github.com/php-any/origami v0.0.0\n", pkgName)
    return os.WriteFile(filepath.Join(outputDir, "go.mod"), []byte(content), 0644)
}
```

### Step 2: 完善 runCompileCommand

```go
func runCompileCommand(cmd *cobra.Command, args []string) error {
    vendorDir := args[0]

    info, err := os.Stat(vendorDir)
    if err != nil {
        return fmt.Errorf("目录不存在: %s", vendorDir)
    }
    if !info.IsDir() {
        return fmt.Errorf("不是目录: %s", vendorDir)
    }

    files, err := collectPhpFiles(vendorDir)
    if err != nil {
        return fmt.Errorf("扫描失败: %w", err)
    }
    if len(files) == 0 {
        return fmt.Errorf("未找到 .php 文件: %s", vendorDir)
    }

    fmt.Printf("找到 %d 个 PHP 文件\n", len(files))

    parsed, parseErrs := parseFiles(files)
    if len(parseErrs) > 0 {
        for _, e := range parseErrs {
            fmt.Fprintf(os.Stderr, "警告: %v\n", e)
        }
    }
    if len(parsed) == 0 {
        return fmt.Errorf("没有文件解析成功")
    }

    fmt.Printf("成功解析 %d 个文件\n", len(parsed))

    if err := generateOutput(parsed, compileOutput, compilePkg); err != nil {
        return fmt.Errorf("生成失败: %w", err)
    }

    fmt.Printf("已生成 Go 包到 %s\n", compileOutput)
    return nil
}
```

### Step 3: 验证完整流程

```bash
cd D:/github.cocm/php-any/origami && go build ./... && ./origami compile vendor/
```

### Step 4: 提交

```bash
git add cmd/compile_output.go cmd/compile.go
git commit -m "feat: add output file generation for compile command"
```

---

## Task 15: 集成测试

**Files:**
- Create: `tests/compile_test/` 目录和测试文件

### Step 1: 创建测试用 vendor 目录

创建一个小型测试 vendor 目录 `tests/compile_test/vendor/test/lib.php`：

```php
<?php
class TestLib {
    public $name;
    
    public function __construct($name) {
        $this->name = $name;
    }
    
    public function greet() {
        return "Hello, " . $this->name;
    }
}

function test_add($a, $b) {
    return $a + $b;
}
```

### Step 2: 运行编译

```bash
cd D:/github.cocm/php-any/origami && go build ./... && ./origami compile tests/compile_test/vendor/ -o tests/compile_test/build
```

### Step 3: 检查生成的代码

检查 `tests/compile_test/build/` 下的文件是否正确生成。

### Step 4: 提交

```bash
git add tests/compile_test/
git commit -m "test: add compile command integration test"
```

---

## Task 16: 持续完善节点覆盖

随着实际 vendor 代码测试，发现未覆盖的节点类型时，逐步添加到 `genGetValue` 分派函数中。每个新节点类型：

1. 读取对应的 `node/*.go` 文件，确认结构体字段
2. 在对应的 `cmd/compile_gen_*.go` 文件中添加生成函数
3. 在 `genGetValue` 中添加 case 分支
4. 编译验证
5. 提交

目标：覆盖 P0 节点类型的 90%+。
