# 施工单 00：基线与度量设施

## 任务卡

- **目标**：让"变快了多少"成为可复现的数字。建立 Go Benchmark 套件、进程外 pprof 采样、PHP 官方解释器对照、Laravel 端到端压测。
- **独占文件**：`bench/`（新建目录）、`scripts/bench/`（新建目录）、`cmd/profile.go`（新建）、`cmd/root.go`
- **前置依赖**：无。**这是 01–07 全部方案的前置**。
- **禁止触碰**：任何 `node/`、`data/`、`runtime/`、`parser/`、`lexer/` 的执行逻辑。本单只加度量，不改行为。
- **预计规模**：新增约 8 个文件，改动 `cmd/root.go` 约 10 行。

## 为什么这单必须先做

仓库当前的性能资产只有根目录 `performance_comparison.md` 一份一次性微基准（一百万次 `$value = $i; $result = $value * 2; $sum = $value + $result`，Origami 0.1085s vs Python 0.0701s）。它无法用来验证任何具体优化：没有机器信息、没有复跑协议、没有 PHP 官方对照、粒度只有"整个循环"。

全仓库 55 个 `*_test.go` 里**没有一个 `Benchmark*` 函数**，没有 `runtime/pprof`、`net/http/pprof`、`runtime/trace` 的任何引用，没有 wrk/ab/hey 脚本。也就是说，在这单完成前，后面 7 份方案的"收益"全都只能是猜测。

同时 AGENTS.md 明确禁止在 `Call()`、每方法/函数、`json_encode` 等热路径插追踪。所以度量必须是**进程外采样**（pprof）和**进程外计时**（Benchmark / 压测），不能是埋点计数。这一条决定了下面的技术选型。

## Step 1：Go Benchmark 套件

新建 `bench/` 包（不放在 `runtime/` 或 `node/` 里，避免和各方案的独占文件冲突）。

### 1.1 公共夹具 `bench/bench_test.go`

Benchmark 必须能在**不重建 VM** 的前提下反复执行同一段 PHP，否则测到的全是启动开销。关键点是把「建 VM + 解析」放在计时外，只对「执行」计时。

```go
package bench

import (
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
	"github.com/php-any/origami/std"
	"github.com/php-any/origami/std/php"
)

// newVM 建一个只加载 std + php 的 VM。不加 http/websocket，减少无关注册开销。
func newVM(tb testing.TB) (*runtime.VM, *parser.Parser) {
	tb.Helper()
	p := parser.NewParser()
	vm := runtime.NewVM(p).(*runtime.VM)
	std.Load(vm)
	php.Load(vm)
	return vm, p
}

// compileSnippet 把一段 PHP 源码解析成 AST，返回可重复执行的闭包。
// 解析在计时外完成，闭包内只做 GetValue，这样 Benchmark 测的是解释执行。
func compileSnippet(tb testing.TB, src string) (data.GetValue, []data.Variable, *runtime.VM) {
	tb.Helper()
	vm, p := newVM(tb)
	// 写临时文件后走正常解析入口，避免绕过 lexer 的模板/PHP 模式处理
	path := writeTempPHP(tb, src)
	program, vars, ctl := vm.ParseFileCached(path)
	if ctl != nil {
		tb.Fatalf("解析失败: %v", ctl)
	}
	return program, vars, vm
}

// runSnippet 执行一次，返回控制流以便断言没有异常。
func runSnippet(tb testing.TB, vm *runtime.VM, program data.GetValue, vars []data.Variable) {
	ctx := vm.CreateContext(vars)
	if _, ctl := program.GetValue(ctx); ctl != nil {
		if _, isReturn := ctl.(data.ReturnControl); !isReturn {
			tb.Fatalf("执行失败: %v", ctl)
		}
	}
}
```

`writeTempPHP` 用 `tb.TempDir()` + `os.WriteFile`，文件名带 `.php` 后缀（`.zy` 走不同的词法路径，基准要测 PHP 路径）。

**注意 `ParseFileCached` 会按路径缓存 AST**，所以每个 Benchmark 用独立临时文件名，避免互相污染。

### 1.2 必须覆盖的基准项

按施工单归属组织文件，每份方案有明确的"我该看哪几个数字"。

`bench/values_test.go`（服务施工单 01）：

```go
// BenchmarkIntArithmetic 纯整数算术，测标量装箱成本
// $sum = 0; for ($i = 0; $i < 1000; $i++) { $sum = $sum + $i * 2; }
func BenchmarkIntArithmetic(b *testing.B) { ... }

// BenchmarkStringConcat 字符串拼接，测 BinaryDot 的 reflect + fmt 成本
// $s = ""; for ($i = 0; $i < 500; $i++) { $s = $s . "x" . $i; }
func BenchmarkStringConcat(b *testing.B) { ... }

// BenchmarkStringInterpolation 插值，测 BinaryLink 链
// for (...) { $s = "prefix {$a} mid {$b} suffix"; }
func BenchmarkStringInterpolation(b *testing.B) { ... }

// BenchmarkStringLiteralRead 循环里反复读同一字面量，测 StringLiteral.GetValue 的每次分配
func BenchmarkStringLiteralRead(b *testing.B) { ... }
```

`bench/array_test.go`（服务施工单 02）：

```go
// BenchmarkArrayDenseIndexRead 密集整数键读取，当前是 O(n)，是本项要打掉的目标
// $a = [...1000 项...]; for ($i = 0; $i < 1000; $i++) { $x = $a[$i]; }
func BenchmarkArrayDenseIndexRead(b *testing.B) { ... }

// BenchmarkArrayStringKeyRead 字符串键读取，当前 O(n)
func BenchmarkArrayStringKeyRead(b *testing.B) { ... }

// BenchmarkArrayAppend 尾部追加
func BenchmarkArrayAppend(b *testing.B) { ... }

// BenchmarkArrayAssignLarge 大数组赋值，测 eager CloneArrayValue
// $a = [...1000 项...]; for (...) { $b = $a; }
func BenchmarkArrayAssignLarge(b *testing.B) { ... }

// BenchmarkArrayPassByValue 大数组作参数传递
func BenchmarkArrayPassByValue(b *testing.B) { ... }

// BenchmarkForeachLarge 遍历，测快照分配
func BenchmarkForeachLarge(b *testing.B) { ... }
```

`bench/call_test.go`（服务施工单 03）：

```go
// BenchmarkFuncCallSimple 无参函数调用，测建帧成本
func BenchmarkFuncCallSimple(b *testing.B) { ... }

// BenchmarkFuncCallManyLocals 函数体内声明 20 个局部变量但只用 1 个，
// 放大 resetVariables 的「每槽一个新 ZVal」成本
func BenchmarkFuncCallManyLocals(b *testing.B) { ... }

// BenchmarkFuncCallTyped 带标量类型声明的参数，测 PrepareTypedValue
func BenchmarkFuncCallTyped(b *testing.B) { ... }

// BenchmarkFuncCallTypedClass 带类类型声明的参数，测 Class.Is 的继承链字符串比较
func BenchmarkFuncCallTypedClass(b *testing.B) { ... }

// BenchmarkMethodCall 实例方法调用
func BenchmarkMethodCall(b *testing.B) { ... }

// BenchmarkStaticCall 静态方法调用，当前走双帧，是本项要打掉的目标
func BenchmarkStaticCall(b *testing.B) { ... }

// BenchmarkClosureCall 闭包调用，当前双重 CreateContext
func BenchmarkClosureCall(b *testing.B) { ... }

// BenchmarkBuiltinCall 内置函数调用（strlen），测 GetParams 每次分配
func BenchmarkBuiltinCall(b *testing.B) { ... }

// BenchmarkIssetEmpty 测 empty/isset 因 ParameterRawAST 永远走慢路径
func BenchmarkIssetEmpty(b *testing.B) { ... }

// BenchmarkReturnValue 大量 return，测 ReturnValue 堆分配
func BenchmarkReturnValue(b *testing.B) { ... }
```

`bench/class_test.go`（服务施工单 04、05）：

```go
// BenchmarkNewSimple 实例化无继承类
func BenchmarkNewSimple(b *testing.B) { ... }

// BenchmarkNewDeepInherit 实例化 5 层继承的类，
// 放大「每次重跑属性默认值 AST + 整链抽象校验 + 造丢弃父对象」
func BenchmarkNewDeepInherit(b *testing.B) { ... }

// BenchmarkPropertyRead $this->x 读取，测 OrderedMap 双重加锁
func BenchmarkPropertyRead(b *testing.B) { ... }

// BenchmarkPropertyWrite
func BenchmarkPropertyWrite(b *testing.B) { ... }

// BenchmarkMethodLookupDeep 调用继承链深处的方法
func BenchmarkMethodLookupDeep(b *testing.B) { ... }

// BenchmarkStaticPropertyRead Foo::$x 与 Foo::CONST，无 lookup cache
func BenchmarkStaticPropertyRead(b *testing.B) { ... }

// BenchmarkInstanceOfDeep instanceof 深继承 + 接口
func BenchmarkInstanceOfDeep(b *testing.B) { ... }

// BenchmarkClassExistsMiss 反复探测不存在的类，
// 直接命中 findClassCaseInsensitive 的全表 Range + FindClassFile 无负缓存
func BenchmarkClassExistsMiss(b *testing.B) { ... }
```

`bench/parse_test.go`（服务施工单 06）：

```go
// BenchmarkLexLargeFile 对一个真实 vendor 文件做词法，测 token 分配
func BenchmarkLexLargeFile(b *testing.B) { ... }

// BenchmarkParserClone 只测 Parser.Clone，暴露每次重建词法 DAG 的成本
func BenchmarkParserClone(b *testing.B) { ... }

// BenchmarkParseLargeFile 词法+语法全量
func BenchmarkParseLargeFile(b *testing.B) { ... }

// BenchmarkAutoloadColdClass 冷加载一个类（含 FindClassFile 的 Stat/ReadDir）
func BenchmarkAutoloadColdClass(b *testing.B) { ... }
```

`BenchmarkLexLargeFile` / `BenchmarkParseLargeFile` 需要一个稳定的大文件输入。**不要指向 `examples/laravel13/vendor/`**（体积巨大且可能不在检出中）；在 `bench/testdata/` 放一个自造的、约 2000 行的 PHP 文件，覆盖类、方法、注解、heredoc、插值、alt 语法。

### 1.3 运行与归档

```bash
# 全套，10 次取样，含分配统计
go test ./bench/ -bench=. -benchmem -count=10 > bench/results/<日期>-<git短哈希>.txt

# 单项，改动前后对比
go test ./bench/ -bench=BenchmarkArrayDenseIndexRead -benchmem -count=10
```

对比用 `benchstat`（`golang.org/x/perf/cmd/benchstat`）：

```bash
benchstat bench/results/before.txt bench/results/after.txt
```

`bench/results/` 加入版本控制，文件名含日期与 git 短哈希。**每份施工单的收尾必须提交一份 after 结果。**

`bench/results/README.md` 记录机器口径：CPU 型号、核数、Go 版本、OS、是否插电、是否有其他负载。跨机器的数字不可比，只比同机 before/after。

## Step 2：进程外 pprof 采样

### 2.1 为什么用环境变量而不是命令行 flag

`zy.go` 的 `main` 有一条**绕过 cobra 的直接脚本路径**：

```30:36:zy.go
func main() {
	if len(os.Args) > 1 && cmd.IsDirectScriptArg(os.Args[1]) {
		if err := cmd.RunScriptFile(os.Args[1]); err != nil {
			os.Exit(1)
		}
		return
	}
```

`IsDirectScriptArg` 对任何不以 `-` 开头、且不在子命令白名单里的参数都返回 true。也就是说 `./origami script.php --cpuprofile=x` 里的 flag **根本不会被解析**。同理 `examples/laravel13/main.go` 是独立 main，也不走 `cmd` 的 flag。

所以采样开关用**环境变量**，这样三个入口（`zy.go` 直接脚本、`zy.go` 子命令、`examples/laravel13`）可以共用同一套实现。

### 2.2 `cmd/profile.go`

```go
package cmd

import (
	"os"
	"runtime"
	"runtime/pprof"
	"runtime/trace"
)

// StartProfiling 按环境变量开启采样，返回的 stop 必须在进程退出前调用。
// 环境变量：
//   ORIGAMI_CPUPROFILE=<path>  CPU profile
//   ORIGAMI_MEMPROFILE=<path>  退出时的堆快照
//   ORIGAMI_TRACE=<path>       execution trace
// 未设置时零开销：不注册任何 hook，解释器热路径完全不受影响。
func StartProfiling() (stop func()) {
	var stops []func()

	if path := os.Getenv("ORIGAMI_CPUPROFILE"); path != "" {
		f, err := os.Create(path)
		if err == nil {
			if err := pprof.StartCPUProfile(f); err == nil {
				stops = append(stops, func() {
					pprof.StopCPUProfile()
					_ = f.Close()
				})
			}
		}
	}

	if path := os.Getenv("ORIGAMI_TRACE"); path != "" { /* trace.Start / trace.Stop */ }

	if path := os.Getenv("ORIGAMI_MEMPROFILE"); path != "" {
		stops = append(stops, func() {
			f, err := os.Create(path)
			if err != nil {
				return
			}
			defer f.Close()
			runtime.GC() // 取 in-use 前先 GC，否则读到的是垃圾
			_ = pprof.WriteHeapProfile(f)
		})
	}

	return func() {
		for i := len(stops) - 1; i >= 0; i-- {
			stops[i]()
		}
	}
}
```

### 2.3 接进 `RunScriptFile`：必须处理 `os.Exit`

`cmd/root.go` 的 `RunScriptFile` 在脚本调用 `exit()` 时会直接 `os.Exit(code)`：

```53:64:cmd/root.go
	if err != nil {
		if exit, ok := err.(interface {
			IsExit() bool
			GetCode() int
		}); ok && exit.IsExit() {
			vm.RunShutdownCallbacks()
			if code := exit.GetCode(); code != 0 {
				os.Exit(code)
			}
			return nil
		}
```

`os.Exit` **不跑 defer**，profile 文件会是空的或截断的。这是本步最容易踩的坑。改法是在 `os.Exit` 前显式 stop：

```go
func RunScriptFile(scriptPath string) error {
	// ...文件存在性检查...
	stopProfiling := StartProfiling()
	defer stopProfiling()

	vm, p := getRuntimeVM()
	_, err := vm.LoadAndRun(scriptPath)
	if err != nil {
		if exit, ok := err.(interface{ IsExit() bool; GetCode() int }); ok && exit.IsExit() {
			vm.RunShutdownCallbacks()
			if code := exit.GetCode(); code != 0 {
				stopProfiling() // os.Exit 不跑 defer，必须先落盘
				os.Exit(code)
			}
			return nil
		}
		p.ShowControl(err)
	}
	vm.RunShutdownCallbacks()
	return nil
}
```

`stopProfiling` 要做成幂等的（`sync.Once` 包一层），因为 defer 和显式调用可能都执行。

同样的两行也要加进 `examples/laravel13/main.go` 的 `main`，让 `serve` 能被采样。

### 2.4 使用

```bash
# 采一个脚本
$env:ORIGAMI_CPUPROFILE="cpu.out"; ./origami.exe tests/php/some_test.php
go tool pprof -http=:8080 cpu.out

# 采 Laravel serve：起服务、压测、Ctrl+C（需要 serve 侧接信号后调 stop）
$env:ORIGAMI_CPUPROFILE="serve-cpu.out"; ./laravel13.exe serve

# 看分配来源
$env:ORIGAMI_MEMPROFILE="mem.out"; ./origami.exe bench/testdata/heavy.php
go tool pprof -sample_index=alloc_objects mem.out
```

`serve` 是长期运行的，Ctrl+C 走信号处理，要确认 `stopProfiling` 在信号分支里被调用，否则拿不到数据。若 serve 侧信号处理复杂，退路是给 serve 加一个只在 `ORIGAMI_CPUPROFILE` 非空时启用的定时停止（如采样 30 秒后自动落盘）。

## Step 3：PHP 官方解释器对照

现有 `performance_comparison.md` 只跟 Python 比，这对"对齐 PHP"的项目没有指导意义——需要知道的是**离 PHP 官方还有多远**。

新建 `scripts/bench/compare.ps1`（Windows 主用）与 `scripts/bench/compare.sh`（CI/Linux），对同一批 `.php` 脚本分别跑 `php` 和 `origami`，输出对照表。

对照脚本放 `scripts/bench/cases/`，每个是**纯 PHP、无 Origami 方言**（必须能被官方 `php` 直接执行）：

- `loop_arith.php`：整数算术循环
- `string_concat.php`：字符串拼接
- `array_dense.php`：密集数组读写
- `array_assoc.php`：关联数组读写
- `func_call.php`：函数调用
- `method_call.php`：实例方法调用
- `static_call.php`：静态方法调用
- `new_object.php`：对象实例化
- `sort_json.php`：`sort` + `json_encode` 混合

每个脚本自己用 `microtime(true)` 打印耗时，脚本参数化迭代次数（`argv[1]`），这样两个解释器跑的是同一份代码、同一个次数。

输出格式（写入 `scripts/bench/results/<日期>.md`）：

```
用例              PHP 8.x     Origami     倍数
loop_arith        0.0123s     0.0891s     7.2x
string_concat     ...
```

**倍数是本项目最有意义的单一指标**，比绝对时间更稳定。目标是让每个用例的倍数随施工单推进单调下降。

如果机器上没装官方 PHP，脚本要能优雅跳过 PHP 侧并只输出 Origami 数字，不能直接失败。

## Step 4：Laravel 13 端到端压测

`examples/laravel13/go-support/serve_command.go` 用 Go 的 `net/http` 起服务，默认 `127.0.0.1:8000`，启动时无条件 `WarmupVendorClassmap`。

新建 `scripts/bench/laravel.ps1` / `.sh`：

1. 构建：`cd examples/laravel13 && go build -mod=mod -o laravel13 .`
2. 后台启动 `./laravel13 serve`，轮询 `http://127.0.0.1:8000/` 直到返回 200（**记录这段冷启动时间，它本身就是 06 的关键指标**）。
3. 压测：优先 `wrk`，退路 `hey`，再退路 `ab`。固定 `-t2 -c10 -d30s`，压首页与一个 Livewire 端点。
4. 记录 RPS、p50/p99 延迟、进程 RSS 峰值。
5. 关服务、归档到 `scripts/bench/results/laravel-<日期>.md`。

Windows 上没有 wrk 时，退路是 PowerShell 并发 `Invoke-WebRequest` 循环——精度差，但足够看出数量级变化。目录里已有的 `examples/laravel13/_profile_dash.ps1` / `_profile_stats.ps1` 是探活脚本（看 CPU/工作集、超时杀进程），可以复用它们的进程监控部分。

**冷启动时间要单独记录。** Laravel 首请求要 autoload 数百个 vendor 文件，这个数字直接反映 04（类路径查找）和 06（AST 缓存统一）的效果，而稳定态 RPS 主要反映 01/02/03/05/07。两者混在一起看会误判。

## Step 5：把基线数字写进文档

跑完上面四步，在 `bench/results/` 和 `scripts/bench/results/` 各留一份初始快照，并在本文件末尾补一节「初始基线」，记录：

- `go test ./bench/ -bench=. -benchmem` 的完整输出
- PHP 对照倍数表
- Laravel 冷启动时间 + 稳定态 RPS
- CPU profile 的 top 20（这会直接告诉后续施工单该先做哪一份）

`performance_comparison.md`（仓库根）在本单完成后应标注为历史文档，并指向本目录。**不要删除**，它是唯一的历史数据点。

## 语义风险

本单不改执行逻辑，风险很低，但有两点要注意：

1. **`StartProfiling` 在环境变量未设置时必须是零开销**。不要为了"方便"注册全局 hook、包装 `Call`、或者启动后台 goroutine。这会直接违反 AGENTS.md 的热路径禁令。
2. **Benchmark 里的 `ParseFileCached` 会污染 VM 状态**。每个 Benchmark 用独立 VM + 独立临时文件，不要为了"快"共享 VM，否则类注册表会互相影响，测出来的数字不可复现。

## 回归门槛

```bash
go build ./... && go test ./...
go test ./bench/ -bench=. -benchmem -count=1   # 先确认能跑通
go run ./zy.go tests/run_tests.php             # 确认 cmd/root.go 的改动没破坏正常执行
```

再手工确认三条采样路径都能落盘非空文件：

```bash
$env:ORIGAMI_CPUPROFILE="t-cpu.out"; ./origami.exe tests/php/strlen_test.php
$env:ORIGAMI_MEMPROFILE="t-mem.out"; ./origami.exe tests/php/strlen_test.php
# 再测一个会调 exit() 的脚本，确认 os.Exit 路径也落盘了
```

无需新增 `tests/php/` 测试——本单不改 PHP 语义。这是 8 份施工单里**唯一**不要求补 PHP 回归的一份。

## 完成定义

- [ ] `bench/` 下 5 个 `*_test.go`，覆盖上面列出的全部 Benchmark 名
- [ ] `bench/testdata/` 有稳定的大 PHP 文件输入
- [ ] `bench/results/` 有初始快照 + `README.md` 记录机器口径
- [ ] `cmd/profile.go` 实现三种采样，未设环境变量时零开销
- [ ] `cmd/root.go` 与 `examples/laravel13/main.go` 接入，且 `os.Exit` 路径已处理
- [ ] `scripts/bench/compare.*` + `cases/` 能产出 PHP 对照倍数表
- [ ] `scripts/bench/laravel.*` 能产出冷启动时间 + RPS
- [ ] 本文件补上「初始基线」一节，含 CPU profile top 20
- [ ] `performance_comparison.md` 标注为历史文档
