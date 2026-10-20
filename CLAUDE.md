# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

**语言偏好：所有对话使用中文回答。**

代理必须遵守的项目目标与「禁止改 PHP 跳过错误」见仓库根目录 **[AGENTS.md](AGENTS.md)**。运行时改动必须考虑热路径成本（赋值/调用/属性读），禁止给无 `static` 的快路径加固定税。

## Build & Run

优先 `go run`，不要先编 exe 再执行：

```bash
go run ./zy.go <script.php>              # 跑脚本（.php 或 .zy）
go run ./zy.go tests/run_tests.php       # 全量 PHP 回归
cd examples/laravel13 && go run -mod=mod . serve --port=18086
```

仅发布 / 交叉编译 / 正式基准才 `go build`。

The LSP server is a separate module:

```bash
cd tools/lsp && go build -o zy-lsp . && ./zy-lsp
```

## Architecture

Origami-lang is a Go-based interpreter for a PHP-like scripting language. It tokenizes source, parses into AST nodes, then executes those nodes against a runtime VM.

### Pipeline

```
Source → lexer/ → tokens → parser/ → node/ (AST) → runtime/ (VM) executes nodes
```

### Core packages

| Package | Role |
|---|---|
| `token/` | Token type enum and keyword/operator config (`token.go`, `type_config.go`) |
| `lexer/` | Lexer: tokenizes `.zy`/`.php` source into tokens. Entry point: `lexer.go`. Also handles HTML blocks (`html_lexer.go`), PHP mode (`php_lexer.go`), heredoc strings. |
| `parser/` | Recursive-descent parser. `Parser` struct in `parser.go` holds the token stream, scope manager, expression parser, and class-path manager. Grammar rules are sub-parsers registered via `all_parser.go`'s `parserRouter` map—each token type maps to a constructor like `NewIfParser`, `NewFunctionParser`, `NewClassParser`. |
| `node/` | AST nodes. Every language construct (if, for, class, binary ops, calls, etc.) implements `data.GetValue`. This is the largest package—nodes do both structure and execution. |
| `data/` | Core interfaces and value types. Defines `VM`, `Context`, `GetValue`, `Types`, `ClassStmt`, `InterfaceStmt`, `FuncStmt`, and all value types (`ZVal`, `ArrayValue`, `ObjectValue`, `StringValue`, etc.). The `GetValue` interface is the fundamental contract: `GetValue(ctx Context) (GetValue, Control)`. |
| `runtime/` | Process `VM` (`vm.go`) holds the parser, startup declarations, shared parse cache, and compiled programs. `RequestVM` (`request_vm.go`) is the per-execution container: new declarations, constants, globals, include state, static locals, handlers, and shutdown stay there. |
| `std/` | Standard library bundles loaded into the VM: `php/` (PHP builtins like `empty`, `isset`, string/array functions, reflection), `net/http/`, `system/`, `database/`, `channel/`, `context/`. |

### Key interfaces (`data/` package)

- **`GetValue`** — single-method interface `GetValue(ctx Context) (GetValue, Control)`. Every AST node implements this. The return value is itself a `GetValue` (it can be a `Value` like `StringValue`, `ArrayValue`, or another node for lazy evaluation).
- **`Context`** — variable scope chain. Created via `VM.CreateContext(vars)`. Supports nested contexts for function calls.
- **`VM`** — global runtime: `AddClass`, `AddFunc`, `GetOrLoadClass`, `CreateContext`, `LoadAndRun`.
- **`Types`** — type-checking: `Is(value Value) bool`. Implementations include `BaseType`, `UnionType`, `NullableType`, `ClassType`.

### Execution model

1. `main()` creates a `Parser` and a `VM`, loads standard libraries (`std.Load`, `php.Load`, `http.Load`, `system.Load`).
2. `VM.LoadAndRun(path)` clones the parser to tokenize/parse the file, creates a `Context`, then calls `program.GetValue(ctx)` on the top-level AST node.
3. A `Control` return value (second return from `GetValue`) signals non-local control flow: `ReturnControl`, `ThrowControl` (exceptions), `BreakControl`, `ContinueControl`. A nil `Control` means normal execution.

### Class autoloading

Classes are auto-discovered from the filesystem. `ClassPathManager` maps class names to file paths. A class `Foo\Bar` must be in a file named `Bar.php` (or `Bar.zy`) within a directory structure mirroring the namespace. The convention mirrors PSR-4.

### VM and RequestVM

`runtime.NewVM` lives for the process. `runtime.NewRequestVM(base)` keeps a `Base *VM` for shared facilities and copies the startup class-registry snapshot, environment, and executed-file snapshot. Later includes and declarations write only the request maps. Parsing shares the AST; a throwaway `RequestVM` catches declarations so they are not published back to `Base`.

`RequestVM.Activate` publishes that VM's call stack and output into the current goroutine's slot and installs the static-property overlay. Laravel `serve` and `std/net/http` use it, so deadline checks and startup frames see the same objects. `BeginRequestOutput` still allocates a separate anonymous state for callers that only have the base `VM`. `SetOutputWriter` (`fpm.New`) binds output on the `RequestVM` without replacing an existing goroutine slot. Discard the `RequestVM` when the execution ends. There is no `TempVM`.

## Tests
    
Tests are `.php` script files in `tests/` subdirectories (`tests/basic/`, `tests/func/`, `tests/obj/`, `tests/php/`, etc.). The test runner is `tests/run_tests.php` which scans subdirectories and `include()`s each `.php` file. Each test file should print output for manual inspection—there is no assertion framework. Red output in the console indicates failures.

Run all tests: `go run zy.go tests/run_tests.php`

Go-level unit tests exist only in `lexer/` (lexer_test.go, preprocessor_test.go, special_test.go). There are no Go-level tests for the parser or runtime.
