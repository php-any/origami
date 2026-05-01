# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

**语言偏好：所有对话使用中文回答。**

## Build & Run

```bash
go build -o origami .        # Build the interpreter
./origami <script.zy>        # Run a script
./origami tests/run_tests.zy # Run all test suites
```

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
| `runtime/` | `VM` struct (`vm.go`) is the global runtime container: class/interface/function/constant registry. `TempVM` (`vm_temp.go`) provides request-level isolation (analogous to php-fpm's per-request model)—temporary class/function registrations don't leak between requests. |
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

Classes are auto-discovered from the filesystem. `ClassPathManager` maps class names to file paths. A class `Foo\Bar` must be in a file named `Bar.zy` within a directory structure mirroring the namespace. The convention mirrors PSR-4.

### TempVM (request-level isolation)

`runtime.NewTempVM(baseVM)` creates a temporary VM that delegates reads to the base VM but stores new class/interface/function registrations locally. Used to simulate php-fpm's per-request model—discard after the request completes. The LSP server uses this to parse files without polluting the global VM.

## Tests
    
Tests are `.zy` script files in `tests/` subdirectories (`tests/basic/`, `tests/func/`, `tests/obj/`, `tests/php/`, etc.). The test runner is `tests/run_tests.zy` which scans subdirectories and `include()`s each `.zy` file. Each test file should print output for manual inspection—there is no assertion framework. Red output in the console indicates failures.

Run all tests: `go run origami.go tests/run_tests.zy`

Go-level unit tests exist only in `lexer/` (lexer_test.go, preprocessor_test.go, special_test.go). There are no Go-level tests for the parser or runtime.
