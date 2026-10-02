# Laravel 13 运行时设计审查

本文记录 `examples/laravel13` 当前常驻 HTTP 运行模型中已确认的语义缺口和高风险设计，作为后续对齐官方 Laravel 13 / PHP 请求生命周期的修复清单。

下文保留审查时的问题描述；实际实施范围以本节执行记录为准，不能将设计建议视为已经实现。判断标准是：官方 Laravel 应用不应为了适配 Origami 而修改 `vendor/`、业务代码或隐藏错误；通用差异应回推到 Origami 核心、标准库或 Laravel/Symfony 原生兼容层。

## 执行记录：2026-10-01，统一 VM / RequestVM 与异常边界

语言实现统一为 **`runtime.VM` 和 `runtime.RequestVM`**。删除旧 TempVM 实现；`std/php/fpm.RequestVM` 只保留核心类型别名，`fpm.New` 只绑定宿主输出。Worker 是宿主进程的调度概念，不是一种语言 VM。下文历史阶段中的类型名称已更新到现有实现。

VM 提供标准库、启动期声明、共享解析和编译程序缓存。RequestVM 持有请求新增的类/接口/函数、常量、全局变量、会话、include 状态、函数 static 局部变量、错误/异常处理器栈、调用栈、shutdown 和输出。普通及热 HTTP handler 与中间件共用一个 RequestVM；Laravel 与 FPM 入口复用这一实现。显式绑定输出的 RequestVM 在已有 HTTP 作用域中仍使用自己的输出、调用栈和处理器状态。

本阶段修复：

- 请求解析缓存只共享程序，不把解析时发现的请求声明发布到 VM；模板、compile-only 和预编译入口都在 RequestVM 中注册或执行。启动期已执行文件集合使用不可变的代际快照，后续请求 include 状态只写入请求映射。
- include 的作用域依据当前执行栈判断，防止保留的旧 Context 把 Blade 函数局部变量登记为全局。
- 请求创建时发现嵌套闭包和捕获数组中的 PHP 对象，将 ViewFinder、Filament 资源对象等保留句柄映射到请求对象。捕获图复制保留对象别名、循环、数组下一整数键，以及不同回调之间的引用别名；请求之间使用独立引用槽。原生实例状态继续使用明确的服务适配策略。
- CLI 在未捕获异常边界执行 `set_exception_handler`；支持公开函数、实例/静态方法数组、类方法字符串、invokable 和绑定闭包。处理器保持原始 callable 返回身份，支持 set/restore 栈、null、typed/variadic 参数，并传播处理器自身异常和 exit；编译 fatal 不由此处理器吞掉。错误处理器也使用独立请求栈。

验收：

- 扩展 PHP 回归 **100/100** 通过，包含此前失败的两个 CLI 异常处理器用例和新增捕获图用例。这是选择性语义回归，不代表全仓 PHP 测试全部通过。四个新增捕获/异常处理器脚本在宿主 PHP 8.1.34 也通过。
- `runtime`、`std/php/fpm`、`data`、`node` 的完整包 `go test -race` 通过；网络 HTTP 请求边界的相关 race 用例通过。
- 官方 Laravel / Livewire 生命周期、serve 相关回归及 8 个并发登录请求的整组 race 检查通过，最终记录为 36.238 秒；客户端均有超时。
- 通过官方 `go run -mod=mod . serve --port=18088` 验收，登录页和 Livewire 脚本均返回 200；4 个客户端共 8 次并发登录请求全部 200 / 31076 字节，均含 Livewire snapshot，约 94–125 ms。该数字不作为吞吐提升结论。测试服务已停止。
- 全仓 `go test ./... -run '^$'` 编译通过。完整 `std/net/http` 测试仍存在 `TestHomeControllerHelloQueryReturn` 注解路径的 null/string 参数类型失败，未计为通过；本文前述其他全量测试限制仍保留。
- 调用微基准的分配保持不变：普通/精确 int 函数为 48 B、1 次；弱 string 函数为 64 B、2 次；精确类型方法为 160 B、3 次；对象方法 Context 为 112 B、2 次。耗时有样本波动，本次未据此宣称整站吞吐提升。没有在 Call、方法调用或 JSON 热路径添加诊断追踪；对象图诊断在请求结束后进行，诊断代码已移除。

日志位于 `storage/origami-debug/`：`request-vm-final-php-regression.log`、`request-vm-host-php.json`、`request-vm-final-core-complete.log`、`request-vm-final-scope-race.log`、`request-vm-final-laravel-verification.log`、`request-vm-final-repository-compile.log`、`request-vm-final-runtime-bench.log` 和 `request-vm-final-http.json`。

**剩余边界**：统一不可变 ClassDescriptor Registry、完整原生对象图策略、启动期 include 返回值中可变对象的策略、全部 callable 可见性/magic 方法矩阵，以及原生替换类契约门槛仍未完成。统一两种 VM 不构成任意第三方 singleton 或所有 PHP-FPM 隔离语义的完整证明。

## 执行记录：2026-10-01，HTTP 生命周期阶段

本阶段落实请求执行、响应发送和结束清理，未完成本文全部架构重构。没有修改 vendor、应用业务或 Blade 模板。

### 已实施

- **P0-1 / P2-10**：通过官方 `Illuminate\Routing\Pipeline` 执行全局中间件，再进入 Router 的分组和路由 Pipeline；保留跳过中间件的容器配置，支持全局短路。
- **P1-5**：派发 `Terminating` 事件，按路由中间件、全局中间件的顺序重新解析并调用 `terminate()`，最后执行 Application 终止回调。冒号参数只用于 handle，不传给 terminate。
- **P0-4**：shutdown 队列归属于请求 CallState / RequestVM；常驻闭包通过基础 VM 注册时仍进入当前请求。正常结束、exit、PHP 异常、Go panic 和客户端取消均执行一次清理。支持回调参数、invokable 对象、shutdown 期间继续注册及重入防护。清理使用独立的 5 秒期限，并在输出和静态作用域释放前完成。
- **P0-3**：发送器按继承关系识别 StreamedResponse / BinaryFileResponse，执行流式 callback / chunks，绑定请求输出和 flush，不预设流长度。chunks 支持数组、IteratorAggregate、Iterator 和 Generator。普通对象必须先按 PHP Iterator 方法分派，避免被 ObjectValue 的数组迭代器接口误接收。
- **响应边界**：HEAD 和 204/304 抑制响应及 terminate/shutdown 的 body；文件发送使用真实请求处理 Range / 条件请求。响应提交后发生流异常时保留已发送状态和内容，继续执行 shutdown。
- **P1-6（部分）**：执行期限继承 `r.Context()`；Context、PDO 的连接/查询/事务操作、HTTP 文件流和 sleep 感知请求取消。sleep 不再在断开后继续等待原时长。
- **P1-7 / P1-8**：在监听并宣告 ready 前 resolve、bootstrap 和 warm Kernel；删除监听后的异步 vendor classmap 预热。可选的 classmap 预热仍由主入口在执行 artisan 前完成。
- **P2-11 / P2-12 / P2-13**：正常异常交给 Laravel Handler report/render，fallback 只记录一次诊断；删除未实现的 `--tries` / `--no-reload` 选项声明；删除 Livewire hashed dist 的 Server 特判，由官方路由返回 BinaryFileResponse。

### 请求隔离的实际边界

宿主复用启动期 Laravel 应用和程序缓存，语言层统一为 **VM / RequestVM**；应用目前仍未按每个请求重新 bootstrap。

- Application / Router / Kernel 使用请求实例，Container / Facade 静态状态进入请求 overlay。
- 每请求执行官方 `forgetScopedInstances()`，启动期已解析的 scoped binding 也会重新创建。
- Events Dispatcher 的 PHP 属性和 Go 原生 listener/cache/deferred 状态同时隔离；仅复制 PHP 对象不能隔离 AnyValue 指向的原生状态。
- Livewire EventBus 的监听器在请求开始时绑定到当前机制对象，包括仅由闭包 `$this` 保留的实例。Vite 的渲染重置操作写入请求对象。
- Kernel 共享类元数据，原生实例状态通过 `ClassValue.InstanceSource` 暴露，实际随其 ObjectValue 身份存储；方法帧不复制原生状态字段。
- 方法帧只借用对象存储，返回 `$this` 或原生 fluent 对象句柄时保留稳定的实例 Context。回收的 pooled Context 不再成为后续方法调用的对象上下文。
- 从常驻 Context 创建请求方法帧时，调用栈跟随当前输出作用域；header 回调状态也随该请求的 OutputState 隔离。

**P0-2 仍未整体完成**：上述回归证明已覆盖服务和 Laravel scoped binding 的隔离，不能证明任意自定义 mutable singleton、第三方 Manager、捕获引用的监听器或原生扩展都具备 PHP-FPM 等价隔离。现有服务策略仍有专用逻辑，统一策略注册表及通用隔离证明留待后续阶段；不可把当前白名单或 Livewire 策略描述为通用解决方案。

### 验收证据

真实 HTTP 验收入口：`runtime_lifecycle_test.go`，PHP fixture：`tests/origami/runtime_lifecycle.php`。所有客户端带 35 秒超时，取消验收另有 5 秒同步期限。

覆盖全局/分组/路由中间件顺序和短路、终止中间件、shutdown 输出及隔离、callback/chunks、HEAD、204、304、exit、PHP 异常、流中异常、Go panic、文件 Range/HEAD、客户端取消、预解析 scoped binding、Events 注册及 8 个并发请求。真实生命周期 `go test -race` 已通过；执行需可用的 C 编译器。

```powershell
cd examples/laravel13
go test -mod=mod -run TestRuntimeHTTPLifecycle -count=1 -timeout 180s .
# Windows race detector 需要 CGO_ENABLED=1 和已安装的 gcc/clang。
go test -race -mod=mod -run TestRuntimeHTTPLifecycle -count=1 -timeout 180s .
```

新增最小 PHP 回归：

```powershell
go run ./zy.go tests/php/shutdown_callbacks_order_test.php
go run ./zy.go tests/php/streamed_response_lifecycle_test.php
go run ./zy.go tests/php/method_object_identity_test.php
```

上述回归和修改涉及的 Go 包均通过。另运行 45 个对象身份、闭包 this、Reflection、COW、输出缓冲相关 PHP 用例，44 个通过；`coalesce_closure_test.php` 含反引号形式的变量语法，修改前 HEAD 同样解析失败。全量 `go test ./...` 未通过：根模块整套脚本退出失败，parser 的 `TestAltConvertRealFile` 依赖不存在的 `/workspace/...` 文件，`TestReproClassNamedArg` 在未安装 VM 时 panic；均在修改前 HEAD 快照复现，不计为本阶段通过项。

通过官方 `go run -mod=mod -tags origamidebug . serve --port=18087` 实测 `/admin/login` 和 `/livewire-6dd39ca7/livewire.js` 返回 200。四个并发客户端共 16 次登录页请求全部 200，单次约 94–117 ms；该数字包含调试构建成本，不作为吞吐基准或修复前后加速结论。测试服务结束后已停止。

### 性能检查

`runtime/object_identity_test.go` 提供对象方法 Context 和请求对象复制的基准。Windows amd64 / Go 1.27.1，同一机器对比修改前 `e630e8d` 快照与最终实现，三次测量：

- 方法 Context：修改前约 53–58 ns/op，最终约 52–54 ns/op；均为 **112 B/op、2 allocs/op**。
- 初版将原生状态字段复制到每个 ClassValue 帧，曾升至 128 B/op；该布局已撤回，状态改为随 ObjectValue 存储。
- 40 个数组属性的样例中，请求 overlay 的建立约 0.16 μs、7 次分配，全量 clone 约 8 μs、188 次分配。此基准只比较创建成本，不包含首次属性访问，也不证明所有服务迁移均有净收益。

CPU profile 和 pprof top、HTTP 响应、测试日志保存在 `storage/origami-debug/`。profile 含 bootstrap 和请求样本，不能当作纯请求热点排名。没有在生产 Call / 方法 / json_encode 热路径增加追踪。

### 仍需实施的后续阶段

1. **P0-2 完整策略体系**：统一 worker 请求策略、对象图和捕获引用的隔离、通用第三方 singleton/Manager/Facade 并发验收。
2. **P1-6 剩余边界**：文件/流/SPL/直接子进程和加载等待的取消已继续实施，详见后续执行记录；仍需检查其他扩展、不可中断的 OS 调用、共享标准流及 shell 派生进程树。
3. **P1-9 契约门槛**：现有 TargetVersion 仍锁定 Laravel `v13.23.0` / HttpFoundation `v8.1.1`，但完整方法签名、可见性、Reflection、序列化差异矩阵和未覆盖类禁止注册的门槛尚未建立；Livewire 本次验收版本为 `v4.4.0`。
4. **数组体系**：全仓 `.List` 迁移、私有 FlatArrayStore 和结构 Span 已完成，详见后续记录；OverlayArrayStore、键与引用值的解耦、嵌套路径代理、请求级引用提升及完整引用矩阵仍待实施。
5. **类型体系**：名义继承谓词、类别名身份、ASCII 大小写规则、弱标量存储、null 返回和 strict_types 已实施相应阶段；TypeArena / TypeRef、ValueKind 与集中声明检查已开始接入 Parser / Reflection / 代码生成，详见后续记录。声明字段仍使用 Types 适配，ClassDescriptor Registry、完整引用约束和旧 Types 删除仍待实施。
6. **Go 基础设施**：结构 Span、解析 single-flight、解析缓存代际和加载等待取消已实施；arena Span、Set/IDMap/Result、ClassDescriptor 不可变 registry、按 SymbolID 的 autoload 状态机和完整类型基准矩阵仍待实施。

这些项目没有被空壳接口或通过少量 HTTP 用例替代，本文整体执行仍未完成。

## 执行记录：2026-10-01，数组访问 API 阶段

本阶段完成数组迁移的优先入口和相应语义回归，仍未引入 ArrayStore / OverlayArrayStore。没有修改 vendor、应用业务或模板。

- `data/array_access.go` 封装 `Len`、按插入位置读取的 `At`、`Range`、浅槽位 `Snapshot`、`AppendSlotsTo`、`EditSlots`、`ReplaceAll` 和 `SetKey`；复用现有的整数/字符串键查找、Append 和 Unset API。`At` 的位置不等于 PHP 整数键，Snapshot 不代替 PHP 值复制。
- `node/index.go`、`node/value_reference.go`、`node/foreach.go`、`runtime/context.go` 和 Container/Events 指定路径已移除直接 `.List` 访问。可变数组方法经 ReplaceAll / PopSlot / ShiftSlot 更新索引；只读旧数组方法暂时保留 data 包内部的槽位视图。
- `ArraySlotRef` 保存实际 ZVal，替代数组加物理下标的句柄。引用绑定复用 `GetOrCreateZVal`，正确处理 packed 整数键、规范数字字符串、空字符串键、null、稀疏键和省略下标的追加；嵌套绑定先分离值副本，再写回。对象数组的键表达式不再因延迟绑定求值两次。
- 自动整数键状态独立于查找缓存，unset 后仍保留历史下一键，数组值复制及 __call 参数复制保留该状态。支持负整数键以及 PHP 8.3+ 的负数追加规则，PHP_INT_MAX 已占用时追加/array_push 抛出 Error，不覆盖原值。删除字符串键也保留其他整数键身份。
- `array_pop` 仅在移除紧邻下一键的整数时回退计数；shift、unshift、splice、usort、multisort 的重排通过集中 API 重建整数键及缓存。批量编辑复用原切片或已有目标缓冲区，没有为了封装接口额外复制整张数组。
- 修复 `$a[PHP_INT_MAX] = ...` 的解析错误：类型声明识别必须真的跟着变量标记或标识符，不能消费 `]`。null 作为数组键转换为 `''`，删除原先错误发出的 Deprecated 诊断；依据 [PHP 数组手册](https://www.php.net/manual/en/language.types.array.php)。pop 的下一键回退条件对照 [PHP 8.4 array.c](https://github.com/php/php-src/blob/PHP-8.4/ext/standard/array.c)。

### 本阶段验证

新增 PHP 最小回归：

```powershell
go run ./zy.go tests/php/array_reference_slot_keys_test.php
go run ./zy.go tests/php/array_automatic_key_state_test.php
go run ./zy.go tests/php/array_bulk_api_keys_test.php
```

三个用例均通过，并在修改前 `e630e8d` 源码快照分别复现 packed 引用、unset 后追加及单元素 usort 的失败。宿主 PHP 8.1.34 也通过三个用例；负数追加断言明确区分 PHP 8.3 前后的规则，不能把宿主 8.1 的结果视为完整的 PHP 8.4 差分验证。

另有 28 个既有 PHP 用例通过，覆盖数组键/COW、引用参数和构造器、foreach、重排、ArrayAccess、Laravel Arr::set；合计 **31 个相关 PHP 回归通过**。Go 测试覆盖脱离数组后的引用身份、重排后的键缓存和常量索引 AST。data/node/runtime/数组扩展/Container/Events/Serve/HttpFoundation 相关包的 race 检查通过；parser 排除上节已在 HEAD 复现的两个环境基线失败后全部通过。未把这次选择性回归描述为全量 PHP 或全量仓库验收。

官方 `go run -mod=mod . serve --port=18087` 验收 `/admin/login` 和 `/livewire-6dd39ca7/livewire.js` 返回 200；4 个并发客户端共 8 次登录页请求均为 200 / 30912 字节，约 96–121 ms。所有客户端设置 35 秒超时，服务日志没有 Warning / Fatal。Laravel HTTP 生命周期及 8 请求隔离的 `go test -race` 再次通过（9.710 s）。本次测试服务已停止，没有操作已有的 18086 服务。

### 本阶段性能与剩余边界

同一 Windows amd64 / Go 1.27.1，以 `-cpu=1 -count=3` 对比相同基准和修改前 HEAD 快照：

- 128 键数组的现有键引用：修改前约 159–172 ns/op、16 B/op；当前约 35–45 ns/op、8 B/op，均 1 次分配。收益来自哈希索引替代线性扫描，以及直接保存 ZVal 的较小引用句柄。
- packed 键查找约 2.3–2.7 ns/op，已有键替换约 2.6 ns/op，前后均 0 分配。
- 128 次 packed 追加：修改前约 3.02–3.39 μs，当前约 2.95–3.05 μs；前后均 8408 B/op、137 次分配。没有把这些微基准视为整站吞吐提升。
- Range 基准为 0 堆分配，但小数组的 closure iterator 有常数成本；foreach 仍保留已有的快照后切片循环，不强制所有热路径改用 closure iterator。

基准、PHP 失败对照、HTTP 响应和 race 日志位于 `storage/origami-debug/`。

**本节记录数组阶段 1 结束时的边界**：当时其他 node/std 路径仍有直接 `.List` 读写，data 旧只读方法也尚未迁移到 Store。这些访问与 FlatArrayStore 的后续实施见下节；Overlay、引用 ZVal 与键元数据的解耦、完整 foreach 引用/返回引用/捕获引用矩阵、ObjectValue 关联数组统一、子路径代理和请求级 ZVal 提升仍需独立实现与验收。

## 执行记录：2026-10-01，FlatArrayStore 与 PHP 数组语义

本阶段已将全部 Go 调用点迁移出公开 `ArrayValue.List`，包括 node/runtime/std、扩展和示例的原有调用点；不是向旧 Laravel 示例增加桥接。全仓搜索不再有 `.List` 访问，底层字段已删除。

### 已实施

- `FlatArrayStore` 私有嵌入 ArrayValue，持有顺序槽位、查找缓存、自动整数键状态与内部指针。平坦路径不增加 interface 分派、锁或额外存储对象分配。
- `View()` 返回 `Span[*ZVal]`，供连续读取和位置遍历；保留 Snapshot 的明确复制语义。这里是切片结构视图，尚不是 TypeArena 的数值区间 Span，也不是不可变 PHP 值或跨请求写入同步机制。
- 批量修改集中到 EditPreservingKeys / EditReindexing / RemovePositions / FilterSlots 等 API；修改期间 panic 也会失效索引。SPL heap/queue、Collection 和 RequestStack 不再通过复制整个切片后截断实现删除。
- 关联数组字面量统一使用 ArrayValue / SetKey；稀疏整数、空字符串、规范数字字符串、null/bool 键不再依赖 ObjectValue 代替数组。数组展开重新编号整数键，保留字符串键；严格相等检查键的类型、身份及顺序。
- 修复排序稳定性、sort/rsort 重排、asort/ksort/krsort 的键与自动追加状态，以及 SORT_FLAG_CASE。usort/uasort 在副本上执行比较器，再发布结果，比较器捕获的原数组不会观察到中间重排；异常会传播，usort 抛出异常后仍重新编号。实现对照 [PHP 8.4 array.c](https://github.com/php/php-src/blob/PHP-8.4/ext/standard/array.c)。
- 修复 splice 负长度、返回键和 replacement 键处理；reverse、slice、diff/intersect、flip/search、array_is_list、首尾键、array_walk 的稀疏/空键及 COW 行为。
- 修复 ARRAY_FILTER_USE_KEY=2、ARRAY_FILTER_USE_BOTH=1；array_filter 和 array_intersect 返回真正的数组并保留原键。纠正旧 `array_filter_basic_truthy_test.php` 错误的 `[0,1]` 预期为 PHP 的 `[1,3]`。
- 数组 current/key/reset/end/next/prev 使用独立内部指针；无效位置返回 false/null，不能靠 next/prev 恢复。值复制保留指针，foreach 不移动该指针。

### 验证与性能

新增 `tests/php/array_store_mutations_test.php`、`array_store_pointer_test.php`，均通过 Origami 与宿主 PHP 8.1.34；修改前 HEAD 快照分别在稀疏首键类型和初始指针键类型处失败。宿主版本不构成完整 PHP 8.4 差分验收。

扩大回归共 129 个数组/引用/foreach/排序/SPL 脚本，纠正上述错误测试后 126 个通过。三个未通过项明确保留：`array_pointer.php` 要求 fresh 数组 prev 返回倒数第二项、空数组 current/reset 返回 null，并把数组指针函数当作 Iterator 方法调用，预期与 PHP 不符；`proc_open_array_cmd_test.php` 使用 `/bin/echo` 和 stty，`spl_file_object_test.php` 的文件搜索在 Windows 失败，后两项在修改前 HEAD 同样失败。不能将这组选择性回归视为全量仓库通过。

核心、PHP、Symfony、Laravel 和网络相关包的 Go 回归与 race 检查通过；根模块 `go test ./... -run '^$'` 编译通过。原全量测试的已知基线限制仍见 HTTP 阶段记录。

Windows amd64 / Go 1.27.1，在其他回归结束后，以 `-cpu=1 -count=3` 测量并对照修改前 HEAD：

- 已有数组键引用：HEAD 159–174 ns/op、16 B/op；当前 36–38 ns/op、8 B/op，均 1 次分配。
- packed 查找：HEAD 2.35–2.39 ns/op，当前 2.54–2.58 ns/op；替换分别约 2.61–2.72 与 2.80–2.87 ns/op，均 0 分配。报告实际常数差异，不宣称所有操作都加速。
- 128 次追加：HEAD 3.02–3.62 μs，当前 3.04–3.12 μs；均 8408 B/op、137 次分配。
- Span 遍历 8 项约 3.45–3.48 ns，128 项约 59.9–60.2 ns，均 0 分配；对应直接切片遍历约 3.52–3.72 / 60.1–66.6 ns。closure Range 仍有常数成本，热路径使用 Span 循环。
- `unsafe.Sizeof` 回归约束 ZVal 不超过 48 B、ArrayValue 不超过 88 B。以上为微基准，不能替代请求吞吐测量。

**仍未完成 Overlay 和引用实体化**：共享引用槽位仍携带 Name/EmptyStrKey；重排引用数组的副本时，键元数据和引用值尚未完全解耦。部分旧内置仍返回 ObjectValue 关联数组；真实对象的指针/属性行为、完整回调与引用矩阵也未整体迁移。不得将 Flat 的 API 封装宣称为完整请求级数组代理。

## 执行记录：2026-10-01，阻塞取消与并发解析

### 已实施

- 请求拥有的 io.Closer 通过 context.AfterFunc 关闭；回调只捕获真实资源和请求 context，不捕获 pooled VM Context。手动关闭会注销回调；原生 SPL 文件对象也接入同一机制。
- StreamInfo 不再持元数据读锁等待 OS Read/Write；Close 可中断阻塞管道读。PHP fread/stream_get_contents/fwrite 和哈希流的取消会结束 PHP 执行，不返回空串后继续执行。
- proc_open / shell_exec / Symfony Process 使用 CommandContext；proc_close 可取消等待。部分创建失败会关闭已经创建的管道；后台 Wait 回收直接子进程，完成后注销取消回调，done 只关闭一次。WaitDelay 约束继承输出管道的等待，尚不等同于杀死完整进程树。
- proc_open 的 `$pipes` 返回 ArrayValue，保留整数描述符键。shell_exec 无输出返回 null；Windows 的 CRLF 和 Ctrl-Z 按文本管道处理，依据 [PHP shell_exec 手册](https://www.php.net/manual/en/function.shell-exec.php)。
- sleep/usleep/Symfony Clock 的等待可取消，usleep 负数抛出 ValueError。Windows select 超时等待可取消；Unix select 用最多 50 ms 的等待片段感知取消并恢复每轮 fd 集合。
- 本地读取/写入/追加、copy、hash_file、md5_file、finfo、Symfony/Finder 文件读取及 HTTP 文件发送绑定 owned handle 关闭。Background 读取保留 os.ReadFile 路径；HTTP file_get_contents 显式使用调用者的请求 context。
- SPL 行读取检查取消；READ_AHEAD/跳过空行遇到读取错误会停止，不再持续追加空行。请求取消会关闭 SPL 句柄，并删除其创建的临时文件；用户原文件保留。
- PHP 文件加载/编译等待检查请求取消。ParseFileCached 冷缓存采用 single-flight，共享同一解析结果；重入报错，leader panic 或取消会释放等待者并允许重试。热重载通过 atomic.Pointer 切换缓存代际，旧解析不能重新填入新缓存。

### 验证与边界

Go 测试覆盖阻塞管道读关闭、晚注册资源、PHP 读取消传播、HTTP 头部与 body 等待取消、直接子进程和输出管道回收、SPL 临时文件清理与错误退出，以及解析共享、等待者取消、panic/retry、重入和缓存代际。上述相关包的 race 检查通过；Unix stream 包仅做 Linux 交叉编译，未在 Linux 主机执行。

新增 `tests/php/request_resource_basic_test.php` 对照正常文件/流/哈希/copy/进程及参数错误；Origami 与宿主 PHP 均通过，修改前 HEAD 在 shell 输出处失败。

解析缓存命中：HEAD 25.4–27.6 ns/op，当前 25.8–30.3 ns/op，均 0 分配。32 KiB 文件读取：Background 约 63–72 μs、41433 B/op、5 次分配；可取消请求约 64 μs、41601 B/op、9 次分配。资源取消登记存在固定分配成本，不声称取消处理免费，也没有把检查插入所有 Call 热路径。

最终官方 Laravel 生命周期 `go test -race` 通过（9.491 s）；`go run -mod=mod . serve --port=18087` 的登录页和 Livewire JS 均返回 200。4 个并发客户端共 8 次登录页请求全部 200 / 31076 字节，约 82–120 ms；全部设置 35 秒超时，日志没有 Warning/Notice/Fatal。测试服务已停止，已有 18086 服务未操作。日志、响应和微基准均在 `storage/origami-debug/`。

**P1-6 和加载体系仍有边界**：未证明所有阻塞扩展已覆盖；OS open/stat 和不可中断磁盘调用、共享 stdin、shell 的完整派生进程树、Windows select 的真实管道 readiness 仍需实现与验收。解析 single-flight 不代替按 SymbolID 的 autoload 声明状态机，也不代替请求级 autoload 回调隔离或 ClassDescriptor registry。

## 执行记录：2026-10-01，名义类型与异常继承

### 已实施

- `data.NominalIsA` / `InterfaceIsA` 为声明类型、`instanceof`、`is_a`、`is_subclass_of`、迭代及 ArrayAccess 提供共同的已加载元数据判断。支持深层接口、父类接口、循环检测与类别名身份；比较阶段不运行 autoload。
- 动态 instanceof 字符串按完整类名解析，不追加调用点 namespace；对象 RHS 使用对象真实类名。非法 RHS 按 PHP 抛 `Error`。
- `is_subclass_of` 补齐第三参数 `allow_string = true` 和接口关系。只有允许的来源字符串可触发 autoload；目标名不加载。加载器将未找到符号与回调异常分开，回调异常向 PHP 传播。
- 内部 Error 不再被 Exception 捕获，不再把命名空间后缀当作内置类型。修正 ParseError → CompileError 和 BadMethodCallException → BadFunctionCallException；后者同步修正原生类元数据。
- 统一 PHP 声明失败的 `TypeError` 生成，包括参数、变量、属性及函数/方法返回失败。Closure 类型拒绝普通 callable 字符串/数组；iterable 不再接受命名空间中恰好叫 Iterator/Generator 的类。
- 类/接口注册和查找只折叠 ASCII 字母；请求接口支持大小写与前导分隔符，请求代理共享原类元数据并保留独立对象身份。移除 ArrayAccess 判断为调用 instanceof 构造临时对象的路径。

### 验证与性能

- 新增 `nominal_type_relations_test.php`、`nominal_identifier_case_test.php`，宿主 PHP 与 Origami 均通过；原 HEAD 在深层接口关系上失败。另有 20 个相关正常 PHP 回归通过；`throw_aborts_following_statements_test.php` 按设计在未捕获调用处退出 1，日志确认没有执行后续语句。
- data、runtime、node、std/php、std/exception 的普通/race 验收通过；全仓 Go 包编译通过。官方 Laravel bootstrap 和真实并发生命周期 race 回归通过。
- 官方 `go run -mod=mod . serve --port=18087`：登录页 200 / 31076 bytes，HTML 实际引用的 Livewire hashed JS 200 / 564841 bytes；4 并发共 8 次登录页请求全部 200。请求均有 35 秒超时，服务日志无 Warning / Notice / Fatal。
- 同一 Class.Is 基准在原 HEAD 与当前实现各跑 3 次：直接类约 4.15–4.23 → 3.42–3.55 ns；直接父类约 8.51–8.61 → 7.62–7.67 ns。深层接口约 102–103 → 100–116 ns，未命中约 120–121 → 123–126 ns；两种遍历路径从 32 B / 2 allocs 降为 0 B / 0 allocs。直接路径无新增分配，不宣称所有路径均加速。

日志保存在 `storage/origami-debug/nominal-*`。这一阶段仍使用 `ClassStmt` / `Types` 字符串元数据，尚不是 ClassDescriptor / TypeID；标量精确匹配与转换、null/void/never、strict_types、完整引用约束和原生类契约仍需后续迁移。

## 执行记录：2026-10-01，标量存储与 strict_types

### 已实施

- `String` / `Bool` / `Float` 的 `Is()` 只做精确匹配；声明边界使用 `PrepareTypedValueInContext` 返回实际存储值与转换控制流。弱模式参数、属性、函数/方法/闭包返回值保存转换后的标量，union 优先保留精确类型，再按 PHP 的标量转换顺序处理。数字字符串检查拒绝尾随文本及 NaN/INF 文本。
- 非 nullable 声明不再无条件接受 null；显式 null 与省略参数分开处理，只有未定义槽位才使用默认值。普通 `T $value = null` 参数保留隐式 nullable 元数据。声明返回类型的函数、方法和闭包没有执行 return 时抛 TypeError，包括 nullable 返回声明。
- 严格模式属于编译单位：Program 在 include 期间保存/恢复标志，声明保存函数体和返回值的模式，参数采用调用点模式；严格模式仅允许 int 向 float 扩宽。合法性检查覆盖值必须为整数 0/1、声明位置和禁止 block mode；真实 PHP 标签与词法预处理的换行有独立回归。
- direct/named/spread 参数、引用参数初次绑定、variadic、invokable、call_user_func 和新增 call_user_func_array 的已测试入口传播类型错误。array_map/array_filter/preg_replace_callback 使用内部弱参数模式；array_map 保留绑定闭包和对象方法上下文。统一原生绑定器执行默认值并返回失败，所有调用点处理该控制流，缺少必需参数抛 ArgumentCountError。
- Stringable 转换的异常向参数/属性赋值传播；返回转换异常按 PHP 规则包装 TypeError 并保留 previous。依据 [PHP 类型声明手册](https://www.php.net/manual/en/language.types.declarations.php) 与 [PHP 8.4 zend_execute.c](https://github.com/php/php-src/blob/PHP-8.4/Zend/zend_execute.c)。
- 输出引用参数单独标记，不把 matches/result/count 的旧值当作输入类型检查；补 str_ireplace 的 count 写回。PCRE 支持未引用的水平/垂直空白类及其补集、字符类内形式和 UTF 模式，集合对照 [PCRE2 pattern 文档](https://pcre.org/current/doc/html/pcre2pattern.html)。解构赋值返回实际 RHS 并清空缺失项，Carbon 的解构 while 能正确结束；DOMElement::getAttribute 缺失属性返回空串。
- 生成器保留编译单位/声明的严格模式，以及闭包返回类型、static 与引用返回标记。回归实际编译运行生成的 Go 构造代码，覆盖函数、方法及强/弱闭包返回行为。

### 验证与性能

新增 PHP 回归：`typed_scalar_storage_test.php`、`typed_null_return_test.php`、`strict_types_unit_test.php`（含强/弱文件 fixture）、`typed_native_callback_test.php`、`native_output_reference_test.php`、`pcre_whitespace_classes_test.php`、`destructure_assignment_result_test.php`、`dom_attribute_return_test.php`。新增用例在 Origami 和宿主 PHP 8.1.34 做差分验证；宿主不是完整 PHP 8.4 验收。另运行类型、引用、闭包、回调、PCRE、DOM 相关脚本，78/78 通过，排除了前文已复现的 coalesce_closure 语法基线限制。

核心和 PHP 标准库普通测试及 race 检查通过；生成代码执行回归通过；真实 Laravel HTTP 生命周期及 bootstrap 的 race 检查通过（17.734 s）。`go run -mod=mod . serve --port=18087` 下登录页返回 200 / 31076 B，页面实际引用的 Livewire JS 返回 200 / 564841 B；4 个并发客户端共 8 次请求全部 200，约 82–124 ms，每次设置 35 秒超时，测试服务已停止。

同机 Windows amd64 / Go 1.27.1，`-cpu=1 -count=3` 对比修改前 e630e8d 源码快照与当前版本：

- Context 仍为 152 B，TokenFrom 仍为 56 B。曾尝试把模式放到每个 TokenFrom，增加了结构大小，已撤回。
- 无声明调用：HEAD 172–173 ns/op，当前 174–176 ns/op；精确 int 调用 HEAD 178–179 ns，当前 182–197 ns；对象方法 HEAD 312–320 ns，当前 315–345 ns。均没有增加分配次数。移除了绑定器重复读写严格模式的负优化；保留实际测得的常数成本，不宣称所有热路径净加速。
- 弱 int→string 调用：HEAD 181–198 ns、48 B/1 alloc；当前 205–207 ns、64 B/2 alloc。HEAD 实际保留 int，语义错误；当前分配并保存真正的 StringValue。这两组结果不能作为同等正确性下的加速比较。

日志、基准及响应保存在 `storage/origami-debug/scalar-*`。没有修改 vendor、业务或 Blade，没有增加生产调用追踪。

### 仍需迁移

仍使用旧 Types 适配，尚未完成 TypeArena / TypeRef、ClassDescriptor Registry、ValueKind 与集中 TypeChecker；无声明/mixed/void/never 和 self/static/parent 的完整上下文语义、typed reference 后续写入约束、所有 callable/by-reference/spread 的组合矩阵、原生类完整契约仍待实施。浮点格式与溢出数字文本、PCRE 原始字节与完整匹配索引语义也未由本阶段证明。以上完成项不能替代后续架构迁移。

## 执行记录：2026-10-01，TypeArena 与 Livewire 类作用域

本阶段接入紧凑声明类型，并修复真实 Laravel / Livewire 验收暴露的核心语义和 Worker 隔离问题。没有修改 vendor、应用业务或 Blade 模板。

### 类型迁移的实际范围

- `data.TypeRef` / `SymbolID` 为 32 位值，内置声明使用保留编号；名义类型和 union/intersection 进入 TypeArena。复合类型按规范成员去重、排序、驻留；构造时加锁并发布 append-only snapshot，检查只读 snapshot。已有节点和成员视图在并发解析发布后保持稳定，当前 typeNode 为 16 B。
- Parser 的 PHP 声明入口使用 TypeRef，集中匹配和弱/严格转换位于 `data/type_checker.go`；ValueKind 由集中分类函数提供，没有在每个 Value 上新增虚方法。无声明、mixed、void、never 分开表示；void 隐式返回 null，mixed / never 跌出函数体抛 TypeError，非法 bare return / void 返回值在解析时失败。
- self / parent 根据声明类解析，static 根据被调用类解析。参数初次绑定和返回检查覆盖继承方法、静态方法、闭包及原生回调的声明上下文。普通参数进入函数体后是普通局部变量，不再把参数类型错误地施加到后续局部重赋值。
- Reflection 读取 TypeRef 的 kind / members，nullable named type 保留正确名称和 allowsNull；代码生成输出声明构造器，避免序列化只能在当前进程使用的 arena 编号。容器、JSON 默认值、HTTP 绑定、注解与 LSP 等旧入口使用冷路径适配。
- typed reference 返回初次检查解引用后的值，转换写回原槽并保留引用身份。null 数组自动创建真实 ArrayValue，支持 typed array 属性的嵌套写入。debug_backtrace 和异常 getTrace 的帧使用 PHP 数组。
- return 声明验证只遍历语法子节点，跳过已解析函数指针等运行时链接，并防止重复访问；真实 Laravel helper 的递归调用图不再令验证器栈溢出。

**尚未完成字段布局迁移**：参数、属性和返回字段仍为 `Types` interface，内含 TypeRef；不能据此宣称每个声明字段已压缩到 4 B 或旧 Types 已删除。名义类型仍经现有继承谓词检查，没有 ClassID / AncestorSet Registry。callable / resource、DNF 与全部非法声明、typed reference 后续写入约束和所有引用返回组合仍需后续实施。

### Livewire magicActions 根因与请求隔离

用户报告的 `SupportReleaseTokens` 找不到 `magicActions`，实际发生在 `SupportMagicActions::provide()` 注册的闭包读取 `self::$magicActions` 时。静态方法中定义的非 static 闭包具有声明类作用域，但没有对象接收者。旧请求重绑定把这些作用域的 nil ObjectValue 当成同一个对象身份，令不同 Feature 的闭包被重绑到另一类。

- `LambdaExpression.RequestScopeObjects` 和 `BindRequestScope` 只把真实对象身份用于重绑定；没有对象的静态方法作用域保持原声明类。补核心 Go 回归和 `closure_feature_scope_test.php`，覆盖多个 Feature、继承的 self 及随后触发的回调。
- 真实生命周期 fixture 新增官方 `Livewire\trigger('call', ..., '$refresh', ...)`，执行 vendor EventBus 和 SupportMagicActions 监听器；覆盖连续请求及 8 个并发请求。
- 对象方法帧保留调用方 RequestVM 和调用栈，PHP 对象身份不变；原生 Kernel 方法适配也从调用方 Context 建立方法帧。避免通过常驻对象的原始 Context 回到共享 VM / 启动期 CallState。
- Container callbacks 按对象身份映射到最终的请求服务副本，保留别名、闭包接收者及捕获的 ClassValue / ThisValue。先完成服务原生状态隔离，再重绑 callbacks 和 Manager 保留的闭包，避免认证服务闭包指向被第二次克隆替换的副本。保留 Worker Application 引用的服务获得请求副本并重绑 app / container。
- 配置 Repository 每请求复制，隔离登录处理中 config.set 的写入。Exception / ErrorException 构造状态归属 ObjectValue.InstanceSource，构造函数不再写共享方法元数据。函数内 static 局部变量归属 RequestVM；同一请求内重复调用共享槽，不同请求互不共享，CLI 仍保持跨调用状态。

**P0-2 仍未整体完成**：该规则覆盖容器引用、已选择服务和对象捕获，不证明任意对象图、捕获引用、循环数组或未知 native singleton 完全隔离；统一 Worker 策略注册表仍待实施。

### 本阶段验收与性能

- `go test -race ./data ./node ./runtime ./std/exception` 通过；Parser 的声明、strict_types 和递归验证回归通过 race；编译器生成代码执行回归及涉及的 PHP 标准库包测试通过。全仓 `go test ./... -run '^$'` 编译通过；前文记录的整套测试基线限制仍在，未宣称全量测试通过。
- 扩展 PHP 回归 **94/96** 通过，包括原有 85 个类型相关用例及新异常实例回归。`set_exception_handler_test.php` 与 `set_exception_handler_variadic_test.php` 在修改前 e630e8d 快照同样失败，CLI 入口绕过注册的异常处理器，仍是未修复缺口。新增/修订的 9 个用例在宿主 PHP 8.1.34 全部通过。
- 真实 Laravel `TestLoginRequestIsolation` / `TestRuntimeHTTPLifecycle` / `TestServe*` 的最终 race 验收通过（39.147 s，零 race 报告），包括 8 个并发登录请求和官方 Livewire call 事件。所有 HTTP 客户端设置 35 秒超时。
- 官方 `go run -mod=mod -tags origamidebug . serve --port=18087` 实测登录页 **200 / 31076 B**，页面实际引用的 Livewire JS **200 / 564841 B**；4 个客户端共 8 次并发登录请求均为 200，约 117–146 ms。调试构建结果不作为吞吐或修复前后加速结论，验收服务已停止。
- 在其他测试和服务器停止后独立运行 `-cpu=1 -count=3`：最终 arena 标量匹配 2.25–2.35 ns，8 成员 union 9.55–10.12 ns，均 0 B / 0 alloc；实际无声明调用 168–172 ns、精确 int 184–193 ns，均 48 B / 1 alloc；方法 334–350 ns、160 B / 3 alloc。单一内置类型初版仍经 union 掩码转换，弱 int→string 为 217–235 ns；复用原转换函数并移除这一步后为 210–211 ns、64 B / 2 alloc，前一阶段为约 205–207 ns。对象方法 Context 的两组测量分别为 53–56 ns 和 76–81 ns，均 112 B / 2 alloc；其实现未随转换调整变化，时间有明显波动。上述结果证明没有增加帧分配，不能据此宣称所有热路径净加速。

日志、HTTP 响应、race 定位和独立基准位于 `storage/origami-debug/type-arena-*`。没有在生产 Call / 方法 / json_encode 热路径增加追踪。

## 当前模型

语言层只有 VM 与 RequestVM；当前 `serve` 宿主复用启动期应用，执行路径为：

```text
Go net/http 常驻进程
  -> Origami 基础 VM 与常驻 Laravel Application
  -> 每请求 RequestVM + Kernel/Application/Router 沙箱
  -> Laravel Router / Blade / Eloquent / Livewire / Filament
  -> Go http.ResponseWriter
```

该模型接近 Octane Worker。性能收益来自只 bootstrap 一次、缓存 PHP AST、共享只读服务并按请求克隆部分可变状态；主要风险也来自常驻对象、并发请求与 PHP-FPM 请求隔离语义之间的差异。

## 建议方向：使用 VM 透明的请求级代理

Origami 控制 PHP 对象模型、方法分派、`instanceof`、Reflection 和属性存储，因此请求级代理不必采用 PHP 用户态 Wrapper，也不必暴露新的代理类型。代理可以与被代理对象共用：

- PHP 类型名
- 父类与接口列表
- 方法表和方法实现
- 属性声明及可见性规则
- Reflection 元数据

推荐的对象结构是：

```text
RequestScoped ClassValue
  Class          -> 原始 ClassStmt（类型与方法元数据）
  Context / VM   -> 当前请求 RequestVM
  InstanceSource -> 当前请求的原生实例状态
  PropertyStore  -> 请求 Overlay
                      |
                      +-- 未命中时回落到全局实例 PropertyStore
```

这样 PHP 层仍应得到：

```php
get_class($proxy) === Illuminate\Auth\AuthManager::class;
$proxy instanceof Illuminate\Auth\AuthManager;
$proxy instanceof Illuminate\Contracts\Auth\Factory;
```

方法调用继续从原始 `ClassStmt` 查找实现，但必须以请求代理作为 `$this` 执行。属性读取优先读取请求 Overlay，未覆盖的只读状态回落到常驻对象；属性赋值、数组元素写入和引用写入只能落到当前请求。

当前 `data.ClassValue.CloneRequestScoped()` 与 `data.chainStore` 已经具备这一模型的基础：它保留原 `ClassStmt`，将 Context/VM 切换到当前请求，并使用带父级回落的属性表。Application 和 Route 已经采用该方式。后续重点不是再实现一个用户态装饰器，而是把这项能力提升为统一的请求级实例机制，逐步替换 `cloneRequestServices()` 中不必要的全量 `CloneSandbox()`。

### 代理不能省略请求状态策略

透明代理能消除复制整张属性表的固定成本，但不能简单共享所有属性值。PHP 数组按值、对象按引用；父对象数组中的子对象即使随数组进入 Overlay，仍可能指向全局可变实例。因此仍需为不同服务声明请求策略：

| 服务类型 | 请求策略 |
|---|---|
| 完全只读服务 | 直接共享 |
| 大部分只读、少量属性可变 | 透明代理 + 属性 Overlay/COW |
| 带请求缓存的 Manager | 透明代理 + 清空指定缓存键 |
| 持有 Request、Session、User 的对象 | 每请求重新创建 |

典型策略包括：

- AuthManager：将 `app` 重绑到请求 Application，并清空 `guards`，使 SessionGuard 按请求重建。
- SessionManager：将 `container` 重绑到请求 Application，并清空 `drivers` / `session.store`。
- ViewFactory：隔离 render count、component stack、sections、pushes 等渲染状态，并把 `shared['__env']` 指回请求代理自身。
- BladeCompiler：共享启动期配置，隔离 raw blocks、footer、component hash stack 等编译过程状态。
- EngineResolver：共享 resolver 注册信息，失效已解析的 blade/php engine，使有状态 Engine 按请求创建。

因此应把当前的“需要 clone 哪些服务”白名单，升级为“每类服务如何进入请求作用域”的策略注册表。例如：

```go
type RequestScopePolicy struct {
    Override map[string]data.Value
    Reset    []string
    Rebind   map[string]data.Value
    Fresh    bool
}
```

### 原生实例状态必须与类元数据分离

透明代理共享 `ClassStmt` 的前提是 `ClassStmt` 只保存不可变的类型元数据和方法表。如果某个 Go 原生类把实例可变状态保存在实现 `ClassStmt` 的 Go struct 中，共享原 `ClassStmt` 仍会造成跨请求状态泄漏。

长期应将原生对象分成三层：

```text
ClassStmt       全局共享：类名、继承、接口、方法和属性声明
InstanceSource  请求级：Go 原生实例的可变状态
PropertyStore   请求级 Overlay：PHP 可见属性状态
```

Kernel 当前通过为沙箱重新创建带 `kernelState` 的 ClassStmt 达到隔离，但这会混淆“类定义”和“实例状态”。更通用的方案是在 `ClassValue` 上提供请求级 `InstanceSource`，让代理继续共享同一份类元数据。

### 必须保持的透明语义

统一代理实现需要回归以下行为：

- `get_class()`、`is_a()`、`instanceof` 与接口判断和原对象一致。
- protected/private 方法与属性仍按原声明类检查作用域。
- `$this` 始终是请求代理，方法写属性不能落到全局对象。
- Reflection 看到原类，而不是额外的 Proxy 类。
- `spl_object_id()` 在单个请求内稳定；不同请求代理拥有不同对象身份。
- WeakMap 以请求代理身份为键，不能通过全局对象身份串请求。
- 普通赋值、嵌套数组写、引用赋值和 `ArrayAccess` 路径都必须正确触发 COW。
- 原生方法不得绕过 PropertyStore，直接保存并修改全局底层指针。

### 性能判断

透明代理通常适合“大对象、启动期状态多、请求只修改少量属性”的服务，可以减少属性遍历、未使用数组复制和每请求固定分配。但它会增加 local/parent 查找、首次提升和 COW 判断；对很小或几乎所有属性都会修改的对象，直接 clone 或 fresh instance 可能更便宜。

迁移必须以 pprof 和并发回归为依据。建议先迁移 Events Dispatcher、URL Generator 等边界清晰的对象，再迁移 View/Blade，最后处理 Auth/Session。每迁移一个服务，同时比较每请求分配、延迟、QPS、对象身份和跨请求状态隔离。

### 数组也应使用透明 Overlay，而不是默认复制

请求隔离真正要求的是“写入不能污染父数组”，并不要求物理复制整张数组。当前 `chainStore` 在读取父对象的数组属性时将 `ArrayValue` 提升为请求副本，是因为代理边界停在对象属性层：属性读取返回普通可变 `*ArrayValue` 后，后续 `SetStringKey`、`UnsetKey` 或 append 会直接修改该数组，不再经过对象 PropertyStore。

长期方案应保持 `*data.ArrayValue` 这一 PHP 可见类型不变，只把内部存储抽象为：

```go
type ArrayStore interface {
    Len() int
    Get(key ArrayKey) (*ZVal, bool)
    Set(key ArrayKey, value *ZVal)
    Delete(key ArrayKey)
    Append(value *ZVal)
    Range(func(ArrayKey, *ZVal) bool)
}
```

并提供两类实现：

```text
FlatArrayStore
  -> 普通 PHP 数组的完整有序槽位

OverlayArrayStore
  -> parent：全局或上一级 ArrayStore
  -> changed：请求修改和新增的槽位
  -> deleted：请求删除键的 tombstone
  -> order：新增、删除和重插后的顺序信息
  -> children：已创建的嵌套数组代理
```

PHP 和 Go 层仍然得到 `*data.ArrayValue`，因此 `is_array()`、参数类型和现有类型断言不需要认识新的 `ArrayProxy` 类型。区别只在 `ArrayValue` 将操作委托给 Flat 或 Overlay store。

#### 局部操作只记录差异

以下操作应直接在 Overlay 上完成，不复制父数组：

```php
$array[$key];
$array[$key] = $value;
unset($array[$key]);
$array[] = $value;
isset($array[$key]);
count($array);
foreach ($array as $key => $value) { /* ... */ }
```

读取顺序为：请求 `changed`、请求 `deleted`、父 store。写入只更新 `changed`，删除只写 tombstone；遍历按 PHP 数组插入顺序合并 parent 与请求差异。

必须正确维护 PHP 数字键的下一个自动索引，以及字符串键删除后重新插入所产生的新顺序位置。

#### 全局重排操作按需实体化

没有必要为所有数组算法实现复杂的增量 Overlay。会重排大部分槽位的操作，例如：

```php
sort($array);
shuffle($array);
array_splice($array, /* ... */);
array_shift($array);
array_unshift($array, $value);
```

可以先执行 `Materialize()`：把 parent 与当前 delta 合并成请求私有 `FlatArrayStore`，然后复用现有实现。这样普通局部读写保持零全量复制，只有确实要处理整张数组时才支付一次实体化成本。

#### 嵌套数组使用路径代理

对于：

```php
$listeners[$event][] = $listener;
```

读取 `$listeners[$event]` 时，如果父值仍是数组，应创建并缓存一个子 `ArrayValue` Overlay：

```text
listeners Overlay
  -> child[$event] Overlay
       -> parent：global listeners[$event]
       -> append：$listener
```

同一路径在单个请求中必须返回同一个子代理，确保连续写入、对象身份和引用行为稳定。请求结束时整棵 Overlay 一并丢弃。

#### 引用必须提升到请求级 ZVal

引用路径不能返回父数组的 `ZVal`：

```php
$ref =& $array['foo'];
$ref['bar'] = 1;
```

`GetZValForWrite` / `GetZValForReference` 应在请求 `changed` 中创建并返回本地槽位。若槽位值为数组，则本地槽位持有子 Overlay；不得把父数组 `ZVal` 或可变子数组直接暴露给请求引用。

还需覆盖：

- `foreach ($array as &$value)`
- 多变量绑定同一个数组槽位
- 嵌套引用写入
- `unset` 后已有引用的行为
- 传值、传引用参数和返回引用
- `ArrayAccess` 间接修改

#### 数组中的对象仍需作用域策略

Array Overlay 只隔离数组结构，不能自动隔离数组元素引用的可变对象。例如请求读取全局 `guards['web']` 后调用 `setUser()`，数组没有发生结构写入，但共享 SessionGuard 已被修改。

读取对象元素时仍需根据请求策略选择：

- 不可变对象：共享。
- 请求级可变对象：返回同类型透明对象代理。
- 保存 Request、Session 或 User 的对象：每请求重建。
- 外部资源对象：按资源并发安全和事务边界决定共享或请求绑定。

### 数组 Overlay 的实施阶段

当前代码大量直接读写 `ArrayValue.List`，这些路径会绕过任何 ArrayStore。迁移应分阶段完成，避免一次重写整个解释器。

1. **封装访问 API，不改变语义**
   - 为 `ArrayValue` 增加 `Len`、`At`、`Range`、`Append`、`SetKey`、`UnsetKey`、`Snapshot`、`ReplaceAll` 等方法。
   - 优先迁移 `node/index.go`、`node/value_reference.go`、`node/foreach.go`、`runtime/context.go`、`data/value_array*.go` 以及 Container/Events 热路径。
   - 禁止核心执行路径新增直接 `.List` 访问。
2. **引入 FlatArrayStore**
   - 先让所有普通数组走 Store 抽象，保持现有语义和性能基线。
   - 在此阶段跑完整 PHP 回归，确认存储抽象本身没有改变键、顺序和引用行为。
3. **引入 OverlayArrayStore**
   - 首批支持 Get、Set、Delete、Append、Len、Range、数字键和插入顺序。
   - 全量重排操作统一回退到 `Materialize()`。
4. **补齐嵌套代理和引用语义**
   - 增加子 Overlay 缓存与请求级 ZVal 提升。
   - 覆盖 foreach 引用、返回引用、参数引用、unset 和 ArrayAccess。
5. **接入对象请求代理**
   - `CloneRequestScoped()` 遇到数组属性时返回 Overlay Array，不再在首次读取时复制数组。
   - 逐个把 Events、URL Generator、Blade、View、Auth、Session 从全量 CloneSandbox 迁移到请求策略。

### 数组 Overlay 验收矩阵

至少需要验证：

```php
// 父数组不被修改
$a['x'] = 1;

// 嵌套写
$a['x']['y'][] = 2;

// 引用写
$r =& $a['x'];
$r['z'] = 3;

// 删除并重插的顺序
unset($a['x']);
$a['x'] = 4;

// 数字键自动递增
unset($a[10]);
$a[] = 5;

// foreach 引用
foreach ($a as &$value) {
    $value = mutate($value);
}

// 数组内对象的请求隔离
$a['guard']->setUser($user);
```

还必须加入并发请求测试，证明请求 A 的局部写、嵌套代理和引用槽位不会被请求 B 或常驻父数组观察到；并用 pprof 对比当前属性提升复制与 Overlay 的 CPU、分配、延迟和 QPS。

## 类型系统设计审查

当前 `data.Types` 更接近一组分散的值谓词，还没有形成统一的 PHP 类型系统。类型声明、运行时精确匹配、弱类型转换、类关系、LSP 推断和 Reflection 信息混在一起，已经造成明确的 PHP 语义错误，也会阻碍透明请求代理。

### `nil` 混合了无声明、`mixed` 和 `void`

`data.NewBaseType()` 当前对空类型、`mixed` 和 `void` 都返回 `nil`。因此运行时无法区分：

```php
function a() {}
function b(): mixed {}
function c(): void {}
```

这使 `void` 返回具体值、`mixed` Reflection、方法签名兼容和未来 `never` 检查都无法正确实现。

`nil` 最多只能表示内部的“没有类型声明”。`mixed`、`void`、`never` 必须是实际类型节点：

```go
type TypeKind uint8

const (
    TypeUnspecified TypeKind = iota
    TypeMixed
    TypeVoid
    TypeNever
    TypeNull
    // ...
)
```

### `Types.Is()` 混合了精确匹配和弱类型转换

当前 `String.Is()` 接受 Int/Float/Bool，`Bool.Is()` 接受所有实现 `AsBool` 的值，`Float.Is()` 接受所有实现 `AsFloat` 的值。随后 `PrepareTypedValue()` 在 `ty.Is(value)` 成功时原样返回 value。

结果可能出现“声明为 string 的槽位实际保存 IntValue”等内部不变量破坏。类型系统必须拆分：

```go
// 值当前是否就是该类型，不执行转换。
ExactMatch(value Value, env TypeEnv) bool

// 当前 strict/weak 模式下能否赋值。
Accept(value Value, mode CoercionMode, env TypeEnv) bool

// 需要时生成转换后的真实 Value。
Coerce(value Value, mode CoercionMode, env TypeEnv) (Value, error)
```

弱模式下 `foo(string $x)` 接受整数时，绑定到 `$x` 的必须是 StringValue，而不是继续保存 IntValue。

### 非 nullable 类型被无条件允许接收 `null`

`PrepareTypedValue()`、Variable、Parameter 和方法返回路径都有独立的 null 放行逻辑。这不符合现代 PHP：普通 `string`、`int`、类类型等不能接收 null。

只有以下类型可以接受 null：

- `mixed`
- `null`
- `?T`
- 显式包含 `null` 的 union

`void` 的隐式返回是控制流规则，不应通过“所有返回类型接受 NullValue”实现。

### 特殊类型被错误降级

当前存在以下不合理映射或占位实现：

- `false` 被降级为 `bool`，丢失 literal type。
- `true` 和 `never` 缺少正式表示。
- `self` 被映射成 `static`，混淆声明类和 late static binding。
- `StaticType.Is()` 永远返回 true。
- Closure 类型接受 callable string/array，混淆 `Closure` 与 `callable`。
- `iterable` 通过类名特判实现，而不是 `array|Traversable` 的正式类型语义。

应分别提供 `LiteralTrue`、`LiteralFalse`、`Self`、`Parent`、`Static`、`Void`、`Never`、`Closure`、`Callable` 和 `Iterable` 类型。`self` / `parent` / `static` 的检查必须携带声明类和调用类上下文。

### `strict_types` 被解析但没有执行

`declare(strict_types=1)` 当前只消费语法，没有把严格模式写入编译单元或调用帧。应把文件级属性保存在解析结果中：

```go
type CompiledUnit struct {
    StrictTypes bool
}

type CallFrame struct {
    CallerUnit *CompiledUnit
    CalleeUnit *CompiledUnit
}
```

参数和返回值检查根据 PHP 的调用点/声明文件规则选择 strict 或 weak coercion，不能由 `Type.Is()` 自行决定。

### 参数、属性和返回值使用了不同检查路径

当前大致存在：

| 检查位置 | 当前路径 |
|---|---|
| 参数、普通变量 | `PrepareTypedValue()` |
| 属性赋值 | 直接 `property.GetType().Is(value)` |
| 函数返回 | 直接 `Ret.Is(value)` |
| 方法返回 | `Is()` + null 与 `__toString` 特判 |
| catch | `Class.Is()` |
| instanceof | 独立的 `checkClassIs()` 路径 |

应统一为单一 TypeChecker：

```go
type CheckSite uint8

const (
    CheckArgument CheckSite = iota
    CheckReturn
    CheckProperty
    CheckTypedVariable
)

type AssignRequest struct {
    Type        TypeRef
    Value       Value
    Site        CheckSite
    Strict      bool
    Declaring   *ClassDescriptor
    CalledClass *ClassDescriptor
}

func (c *TypeChecker) Assign(req AssignRequest) (Value, *TypeErrorInfo)
```

所有 typed slot 都必须通过这个入口得到已经转换且满足内部类型不变量的 Value。

### 类关系判断重复、依赖字符串且带副作用

类和接口关系目前分散在 `Class.Is`、`isClassValueInstanceOf`、`extendISClass`、`interfaceExtends`、`instanceof`、catch 及内置异常硬编码中。部分谓词还会触发 autoload 和 `ThrowControl`。

问题包括：

- 类型谓词可能产生加载和异常副作用。
- PHP 类名大小写不敏感，但当前主要只移除前导反斜线。
- 父类和接口关系被重复遍历。
- `instanceof`、catch、参数类型和 `is_a()` 可能得出不同结果。
- 内置异常需要单独维护继承 map。

应引入统一、不可变的类型注册表：

```go
type TypeID uint32

type ClassDescriptor struct {
    ID         TypeID
    Name       string
    LookupName string
    Parent     *ClassDescriptor
    Interfaces []*ClassDescriptor
    Ancestors  TypeSet
}
```

autoload 和名称解析属于 TypeResolver；解析后的 `IsA(TypeID)` 只做无副作用的 ID/祖先集合查询。显示名称保留原始大小写，查找名称统一规范化为大小写不敏感形式。

### 类定义与实例状态混合

`ClassStmt` 当前同时承担类型元数据、方法表、属性声明、实例构造，以及部分 Go 原生类的实例可变状态。透明代理若共享这类 ClassStmt，就可能共享原生状态。

目标结构应为：

```text
ClassDescriptor  全局共享且不可变：类型名、父类、接口、方法、属性声明
ObjectInstance   每实例：对象身份、Context、PropertyStore
InstanceState    每实例：Go 原生类的可变状态
```

```go
type ObjectInstance struct {
    Class *ClassDescriptor
    State InstanceState
    Props PropertyStore
}
```

请求代理共享 ClassDescriptor，但拥有请求级 InstanceState 和 PropertyStore Overlay，从而在 PHP 层保持完全相同的类型身份。

### `ObjectValue` 同时被视为 PHP array 和 object

当前 Arrays 类型接受 `*ObjectValue`，Object 类型也接受 `*ObjectValue`；ClassValue 又嵌入 ObjectValue 作为属性容器。同一个 Go 值可能同时满足 PHP array 和 object，破坏值模型。

应收敛为：

```text
ArrayValue      唯一的 PHP array（同时支持整数键和字符串键）
ObjectInstance  唯一的 PHP object
PropertyBag     仅内部属性存储，不实现 PHP Value
```

关联数组也必须迁移到 ArrayValue。ObjectValue 应逐步降级/重命名为内部 PropertyBag，不能参与 `is_array()` 或 `is_object()`。

### Closure 与 Exception 存在多套表示

Closure 目前可能表现为 FuncValue、BoundFuncValue、callable array/string 或带 `__invoke` 的对象。它们都可能 callable，但只有真实 Closure 对象应满足 `instanceof Closure`。

异常也可能是 ClassValue、无 Object 的 ThrowValue，或带 Object 的 ThrowValue，迫使 catch 维护硬编码继承树。

目标是：

- Closure 成为拥有标准 ClassDescriptor 的真实对象；`IsCallable` 与 `IsClosureObject` 分开。
- 所有 PHP 可见异常都是异常 ObjectInstance；ThrowControl 只负责非局部控制流，不再模拟异常类型。

### 运行时类型和静态分析类型混在一起

`LspTypes.Is()` 和 `Generic.Is()` 当前无条件返回 true，`MultipleReturnType` 也不是 PHP 运行时声明类型。如果这些分析占位类型进入运行时，会静默放行错误。

应拆分：

```text
runtime/phpType
  -> PHP 可执行声明类型

analysis/typeInfo
  -> PHPDoc 泛型、模板参数、推断联合、MultipleReturnType、LSP 信息
```

静态分析类型不得实现或进入运行时 Assign/Coerce 路径。

### `Value` 接口过弱，依赖具体 Go 类型断言

当前 Value 主要只有 `GetValue` 和 `AsString`，所以各处通过 `*IntValue`、`*ArrayValue`、`*ClassValue` 等具体类型断言判断 PHP 类型。这使透明代理、Closure、异常和引用值不断增加特判。

应增加稳定的 PHP 值分类：

```go
type ValueKind uint8

const (
    KindNull ValueKind = iota
    KindBool
    KindInt
    KindFloat
    KindString
    KindArray
    KindObject
    KindResource
)

type Value interface {
    GetValue
    Kind() ValueKind
    DebugString() string
}
```

对象额外暴露 ClassDescriptor 和 ObjectID。引用是 ZVal/存储语义，不应成为独立 PHP 类型；类型检查前统一 dereference。

### Type AST 应替代字符串拆分

`NewBaseType()` 当前通过字符串查找和 `strings.Split` 处理 union/intersection，无法可靠表达括号和 DNF 类型，例如 `(A&B)|C`。

Parser 应直接生成不可变 Type AST：

```go
type TypeRef struct {
    Kind     TypeKind
    Name     string
    ID       TypeID
    Members  []TypeRef
    Resolved *ClassDescriptor
}
```

TypeRef 同时服务运行时检查和 Reflection；Parser 不再先拼字符串再二次解析。

### 类型错误应统一生成 PHP `TypeError`

当前类型错误使用多个不同的普通 Go error 字符串，缺少参数位置、函数名、声明类型、实际类型和来源位置。

应统一描述：

```go
type TypeErrorInfo struct {
    Site      CheckSite
    Expected  TypeRef
    Actual    TypeRef
    Function  string
    Class     string
    Parameter string
    Position  int
    From      From
}
```

TypeChecker 负责生成 PHP 可见 `TypeError` 对象，保持参数、属性和返回错误格式及堆栈一致。

### 类型系统目标架构

```text
Parser
  -> unresolved Type AST
       |
       v
TypeResolver
  -> namespace/use/self/parent/static
  -> ClassRegistry / autoload
       |
       v
Resolved TypeRef
       |
       v
TypeChecker
  -> ExactMatch
  -> IsSubtype / IsA
  -> Assign / Coerce
  -> TypeError
       |
       v
Runtime Value
  -> ValueKind
  -> ArrayValue + ArrayStore
  -> ObjectInstance + ClassDescriptor + InstanceState
```

### 类型系统迁移顺序

1. **修复现有正确性**
   - 为 mixed、void、never、true、false、null 建立真实类型。
   - 删除所有类型无条件接收 null 的逻辑。
   - 拆分 ExactMatch 与 Coerce。
   - 参数、属性、返回值统一走 TypeChecker。
   - 实现 `strict_types`。
2. **统一 Type AST**
   - Parser 直接生成 union/intersection/nullable/DNF 类型树。
   - 正确区分 self、parent、static。
   - 将 LSP、generic、multiple-return 类型移出运行时。
   - Reflection 直接读取 TypeRef。
3. **统一类描述符**
   - 引入 ClassDescriptor、InterfaceDescriptor 与 TypeID。
   - 规范化大小写不敏感名称并缓存祖先闭包。
   - 合并 instanceof、catch、class type hint 和 `is_a()`。
   - 将 autoload 移出类型谓词。
4. **整理运行时值模型**
   - 关联数组统一为 ArrayValue。
   - ObjectValue 降级为内部 PropertyBag。
   - Closure 和 Exception 统一为真实对象。
   - 引入 ValueKind，逐步删除具体类型特判。
5. **接入透明请求代理**
   - ClassDescriptor 全局共享。
   - InstanceState、PropertyStore 和 ArrayStore 请求级化。
   - 用统一类型身份保证代理与原对象具有相同类名、继承、接口和 Reflection 结果。

迁移期间应保留旧 `Types` 的适配层，但禁止继续向其增加新特判。新增 PHP 类型能力直接落在 TypeRef/TypeChecker；每完成一个阶段，用宿主 PHP 做差分测试，并同时跑 Laravel 参数绑定、typed property、Reflection、Container autowire 和并发代理回归。

## 面向 Go 1.27+ 的类型系统定制方案

仓库根模块及 `examples/laravel13` 已声明 `go 1.27.1`，因此新类型系统不需要兼容旧 Go。设计可以直接使用 Go 1.27 泛型方法、泛型类型别名和更完整的函数类型推断，并使用现代标准库中的 `iter`、`unique`、`weak`、`slices`、`maps` 及 typed atomics。

官方依据：

- [Go 1.27 Release Notes](https://go.dev/doc/go1.27)
- [Go 1.27 Generic Methods](https://go.dev/blog/generic-methods)
- [Go Language Specification](https://go.dev/ref/spec)
- [unique package](https://pkg.go.dev/unique)
- [weak package](https://pkg.go.dev/weak)

Go 1.27 允许具体类型的方法声明自己的类型参数，但接口方法仍不能声明类型参数，泛型方法也不能用于实现普通接口方法。因此核心动态边界仍使用小而稳定的非泛型接口；泛型用于内部容器、遍历器、构建器、缓存和测试工具，避免设计无法实现的“泛型 Type 接口”。

### 使用紧凑的 TypeRef，而不是大量接口对象

PHP 类型是封闭且数量有限的代数结构。运行时不需要为 Int、String、Union、Intersection 等每种类型分配一个实现 `Types` 的 Go 对象。建议使用紧凑句柄：

```go
type TypeRef uint32

const (
    TypeInvalid TypeRef = iota
    TypeMixed
    TypeVoid
    TypeNever
    TypeNull
    TypeFalse
    TypeTrue
    TypeBool
    TypeInt
    TypeFloat
    TypeString
    TypeArray
    TypeObject
    TypeCallable
    TypeIterable
)

type typeNode struct {
    kind   TypeKind
    flags  TypeFlags
    symbol SymbolID
    first  uint32
    count  uint16
    aux    uint16
}
```

内置类型使用保留 TypeRef；类类型、union、intersection 和 unresolved contextual type 存在 TypeArena 中。TypeRef 按值传递和比较，不在热路径传递 `Types` interface，减少 interface boxing、堆分配和动态方法调用。

`typeNode` 应保持不可变并尽量控制在 80 字节以内。Go 1.27 对小于 80 字节的对象提供 size-specialized allocation；不过核心节点应优先驻留在连续 slice arena 中，分配器优化只是额外收益，不能成为依赖频繁小对象分配的理由。

### 使用 TypeArena 驻留复合类型

每个规范化后的类型只创建一次：

```go
type TypeArena struct {
    nodes    []typeNode
    members  []TypeRef
    byKey    map[TypeKey]TypeRef
}

type TypeKey struct {
    Kind TypeKind
    Hash uint64
}
```

构造 `A|B`、`A&B`、`?A` 时：

1. 递归解析成员。
2. 展平相同的 union/intersection。
3. 去重并按稳定 TypeID 排序。
4. 应用化简规则，例如 `T|never -> T`、`T|mixed -> mixed`。
5. 用规范 key 查询 arena；命中即复用 TypeRef。

这样类型相等变成整数比较；Reflection 和错误格式化再按需遍历 TypeArena。

TypeArena 在 bootstrap/解析期允许单写者扩展，服务 ready 后发布不可变 snapshot。请求热路径只能读取 snapshot，不能触发复合类型重建。

### 名称驻留采用 SymbolID，`unique.Handle` 只用于入口规范化

类名、接口名、方法名和属性名重复率极高。应建立统一符号表：

```go
type SymbolID uint32

type SymbolTable struct {
    byName map[string]SymbolID
    names  []string
}
```

输入名称先执行 PHP 规则规范化：

- 移除前导 `\`
- 类/接口/trait/enum 名建立小写 lookup key
- 保留原始大小写 display name
- 方法名按 PHP 大小写不敏感规则建立 lookup key
- 属性名保持 PHP 的大小写敏感语义

`unique.Handle[string]` 可用于解析入口的 canonical string 去重和跨缓存稳定比较；进入 Registry 后热路径统一转换为 SymbolID。不能在每次类型判断中反复调用 `strings.ToLower` 或持有大量重复 FQN。

### ClassDescriptor 使用不可变快照和数值 ID

```go
type ClassID uint32

type ClassDescriptor struct {
    ID         ClassID
    Name       SymbolID
    Parent     ClassID
    Interfaces Span[ClassID]
    Methods    MethodTableID
    Properties PropertyTableID
    Ancestors  AncestorSetID
    Flags      ClassFlags
}
```

`ClassDescriptor` 只保存类型元数据，不保存实例状态。ClassID、InterfaceID、TraitID 和 EnumID 可共享统一的 nominal TypeID 空间，通过 flags 区分类别。

Registry 采用 copy-on-publish：

```go
type Registry struct {
    current atomic.Pointer[RegistrySnapshot]
    // 写路径仅供 autoload/bootstrap 的单写者使用。
}
```

autoload 完成一批声明后构建新 snapshot，再原子发布。请求中的 `instanceof`、catch 和 typed slot 检查只读取一个稳定 snapshot，不加锁、不触发 autoload。

祖先关系在 descriptor 完成时预计算。初期采用排序的小型 `[]TypeID` + 二分查找；如果实际 profile 表明类数量和检查频率值得，可再引入分块 bitset。不要一开始为每个类分配覆盖全部类型空间的大 bitset。

### 用泛型容器消除重复的 Overlay 与 ID 表实现

Go 1.27 基线允许使用泛型类型和泛型别名统一基础设施：

```go
type Set[K comparable] = map[K]struct{}
type IDMap[K ~uint32, V any] = map[K]V

type ReadStore[K comparable, V any] interface {
    Get(K) (V, bool)
    Range(func(K, V) bool)
}

type Overlay[K comparable, V any] struct {
    parent  ReadStore[K, V]
    local   OrderedStore[K, V]
    deleted Set[K]
}
```

该基础类型可复用于：

- 对象 PropertyStore
- 请求 static property overlay
- Facade resolved instance overlay
- 符号表和请求全局量
- ArrayStore 的普通 key lookup 部分

PHP array 的插入顺序、自动数字键、ZVal 引用和重排语义仍由专用 ArrayStore 扩展；不能因为有通用 Overlay 就把 PHP 数组简化成普通 map。

### 泛型方法用于具体算法，不用于动态接口

Go 1.27 的泛型方法适合把 TypeArena/Registry 上的类型安全算法组织到接收者命名空间：

```go
func (a *TypeArena) Fold[R any](
    root TypeRef,
    init R,
    visit func(R, TypeRef) R,
) R

func (r *Registry) Collect[T any](
    ids iter.Seq[TypeID],
    mapFn func(*ClassDescriptor) (T, bool),
) []T
```

适用场景：

- Type AST 格式化与 Reflection 投影
- union/intersection 规范化
- descriptor 验证和契约测试
- 从 descriptor 生成方法/属性视图
- typed cache 和 typed pool 辅助方法

不应尝试定义：

```go
type Type interface {
    Match[T Value](T) bool // Go 1.27 仍不允许泛型接口方法
}
```

PHP 动态值检查仍应由具体 `Checker` 对 `Value`/ValueKind 做集中分派。

### 使用 `iter.Seq` 提供无中间切片遍历

类型成员、祖先、接口、方法和属性遍历统一提供 iterator：

```go
func (a *TypeArena) Members(ref TypeRef) iter.Seq[TypeRef]
func (r *RegistrySnapshot) Parents(id TypeID) iter.Seq[TypeID]
func (c *ClassDescriptor) MethodIDs() iter.Seq[MethodID]
```

调用方可以直接：

```go
for member := range arena.Members(ref) {
    // ...
}
```

这可以替换大量临时 `[]Types`、`[]string` 和队列分配。对极热、已知长度的小集合仍允许返回 `Span[T]` 视图，不能强制所有路径都使用 closure iterator；以 benchmark 决定具体 API。

### Span 用连续 arena 表示只读列表

```go
type Span[T any] struct {
    First uint32
    Count uint16
}
```

接口列表、union 成员、参数类型和属性表索引存放在共享连续 slice 中，descriptor 只持有 Span。相比每个 descriptor 保存独立 slice，可减少 slice header、底层数组分配和 GC 指针扫描。

`Span.At`、`Span.Range` 等可由泛型方法实现；构建期写入，发布后只读。

### TypeChecker 使用 ValueKind 快速分派

```go
type ValueKind uint8

type Value interface {
    Kind() ValueKind
    // ...
}
```

Checker 首先按 TypeKind 和 ValueKind 处理标量矩阵：

```go
type scalarRule func(*Checker, Value, CoercionMode) (Value, bool)

var weakScalarRules [typeKindCount][valueKindCount]scalarRule
var strictScalarRules [typeKindCount][valueKindCount]scalarRule
```

类类型再走 ClassID/AncestorSet；union/intersection 走 TypeArena 成员。不要继续在每个 Type 实现中写具体 Go 类型 switch。

Go 1.27 更完整的函数类型推断可让规则表直接引用泛型 helper，而不必在每个赋值位置显式实例化；但规则表最终仍是普通具体函数类型，避免热路径通过反射调用。

### 用泛型 Result 约束新 API 的返回形态

现有 `(GetValue, Control)` 允许出现许多无意义组合。新类型系统内部可使用：

```go
type Result[T any] struct {
    Value T
    Err   *TypeErrorInfo
}
```

或者在需要 PHP Control 的边界使用：

```go
type EvalResult[T any] struct {
    Value   T
    Control Control
}
```

TypeResolver、TypeChecker、Reflection builder 和 registry loader 各自使用明确的 T，减少 `any`、`GetValue` 和重复断言。VM 既有接口先通过适配函数接入，不要求一次重写所有 AST 节点。

### `weak.Pointer` 只用于可回收派生缓存

`weak.Pointer` 适合：

- 临时 Reflection 投影视图
- 请求代理派生对象缓存
- 大型格式化/诊断结果缓存
- 不应延长对象生命周期的 lookup side table

不能用于 ClassDescriptor、TypeArena、方法表等语义核心对象。核心 Registry 必须强引用且确定性存活；弱缓存命中与否只能影响性能，不能影响类型判断结果。

### 请求代理使用泛型策略表

```go
type ScopePolicy[S any] struct {
    Rebind func(*S, RequestScope)
    Reset  func(*S)
    Fresh  func(RequestScope) *S
}

type PolicyRegistry struct {
    // 运行时按 ClassID 查找擦除后的 policy，注册端保持强类型。
}

func (r *PolicyRegistry) Register[S any](
    class ClassID,
    policy ScopePolicy[S],
)
```

Go 1.27 的泛型方法允许 `Register` 位于具体 Registry 接收者上，注册 AuthManager、SessionManager、ViewFactory 等策略时保持编译期类型检查；运行时查找仍按 ClassID 走非泛型 fast path。

需要注意：注册后的擦除边界必须集中在 PolicyRegistry 内部，业务热路径不应通过 reflection 执行 policy。可在注册时生成闭包，将具体 `*S` 断言封装一次。

### 并发和缓存使用 typed atomics

不可变 snapshot 通过：

```go
atomic.Pointer[RegistrySnapshot]
atomic.Uint64 // generation
```

发布。Method/Property lookup cache 保存 snapshot generation，generation 变化时整体失效。这样不需要在每次类型判断中加 RWMutex，也不会让请求读取到一半更新的继承图。

autoload 重入和同一个类的并发加载使用 single-flight 风格的按 SymbolID 状态机；类型检查本身不得负责启动加载。

### 针对 Go 1.27 分配器调整热结构

Go 1.27 对小于 80 字节的分配提供更快路径。可据此检查但不能盲目拆对象：

- ZVal 当前约 48 字节，应保持紧凑，避免加入大型接口或多个指针字段。
- TypeRef、ClassID、SymbolID 使用 32 位值类型。
- typeNode、overlay delta entry 和调用期检查请求尽量保持无指针或小于 80 字节。
- 大列表放在 arena slice，通过 Span 引用，不在每个 descriptor 上分配小 slice。
- TypeError 的大字符串只在错误路径格式化。

所有布局优化以 `unsafe.Sizeof` 测试、escape analysis 和 pprof 为准；Go 1.27 分配更快不代表可以忽略分配数量。

### Go 1.27 工具链配套验证

建议增加以下工程约束：

- `go test` 默认启用的 `stdversion` vet 检查，确保各子模块的 `go` directive 一致。
- `go fix -fix=atomictypes` 等 modernizer 用于迁移旧原子操作，但只提交可审查的机械变更。
- 使用 `runtime/pprof` 的 goroutine leak profile 检查请求代理、autoload single-flight 和 iterator 是否遗留阻塞 goroutine。
- 并发 Registry/Overlay 测试使用 race detector；异步请求生命周期测试优先采用可控时钟/调度测试工具。
- benchmark 至少报告 ns/op、B/op、allocs/op，并分别覆盖标量类型检查、ClassID `IsA`、8 成员 union、typed property 写入和请求 Overlay。

### Go 1.27+ 定制实施顺序

1. **建立泛型基础设施**
   - SymbolID、TypeID、Span、Set、IDMap、Result、typed atomics snapshot。
   - 不改变现有 PHP 语义。
2. **引入 TypeArena 与 TypeRef**
   - 驻留内置、类、union、intersection 和 contextual type。
   - 为旧 `Types` 提供双向适配，仅用于迁移。
3. **实现 ClassDescriptor Registry**
   - 名称规范化、ClassID、祖先闭包、不可变 snapshot。
   - 统一 instanceof/catch/is_a/type hint。
4. **实现集中 TypeChecker**
   - ValueKind 规则矩阵、strict/weak coercion、统一 TypeError。
   - 参数、属性、返回值逐条切换。
5. **实现泛型 Overlay/ArrayStore 基础**
   - 复用 Overlay 核心，PHP ArrayStore 保留专用顺序和引用语义。
6. **分离 ClassDescriptor 与 InstanceState**
   - 迁移 Go 原生类的可变状态。
   - 接入强类型 ScopePolicy 注册。
7. **删除旧体系**
   - 删除 `Types.Is()`、字符串 union 拆分、LSP runtime type、异常继承硬编码和 ObjectValue 双重语义。

该方案的关键原则是：用 Go 泛型减少基础设施重复，用数值句柄和不可变 arena 优化热路径，用非泛型小接口维持 PHP 的动态边界。不能因为 Go 1.27 支持泛型方法，就把动态 PHP Value 包装成层层泛型抽象；PHP 类型在运行时仍是动态数据，最佳实现是紧凑 tagged representation + 集中分派。

## P0：明确的语义错误

### 1. 全局中间件没有进入请求 Pipeline

位置：

- `std/laravel/httpkernel/methods.go` 的 `kernelHandle`
- `std/laravel/httpkernel/methods.go` 的 `dispatchToRouter`

当前请求在绑定容器后直接调用 `Router::dispatch($request)`。虽然 Kernel 保存并同步了 `$middleware`、中间件组、别名和优先级，但请求执行端没有复刻官方 `sendRequestThroughRouter()` 的全局 `Pipeline`。

结果：

- `$middleware->append()` / `prepend()` 等配置表面成功，实际可能不执行。
- 维护模式、CORS、可信代理、输入规范化等全局中间件语义可能缺失。
- 当前能执行路由和 `web` / `api` 路由中间件，不能证明完整 Kernel 生命周期成立。

期望：

1. 对齐 Laravel 13 `Foundation\Http\Kernel::sendRequestThroughRouter()`。
2. 将当前请求绑定到容器后，通过全局中间件 Pipeline 再进入 Router。
3. 增加最小回归，分别证明全局、分组、路由中间件的顺序和短路行为。

### 2. 请求隔离依赖手工维护的服务白名单

位置：

- `std/laravel/httpkernel/runtime.go` 的 `Sandbox`
- `std/laravel/httpkernel/runtime.go` 的 `cloneRequestServices`

当前只针对已知可变服务克隆或重置，例如 events、view、Blade compiler、URL generator、auth、cache、session。Application 其他未命中的对象仍可能回落到全局容器实例。

这无法为任意 Laravel 应用提供 PHP-FPM 等价隔离。应用或第三方包注册的 singleton 如果保存 Request、用户、租户或其他可变状态，就可能跨请求共享；并发时还可能形成数据竞争。

当前方案的问题不是已有克隆逻辑无效，而是正确性依赖持续发现泄漏并扩充白名单，无法形成封闭证明。

期望：

- 明确选择并记录运行模型：
  - 严格 PHP-FPM 模型：每请求拥有完整请求级容器状态；或
  - 明确的 Octane 模型：实现标准 Worker reset/flush 生命周期，并公开常驻 singleton 约束。
- 为自定义 singleton、第三方 Manager、Facade 静态状态和并发请求增加隔离测试。
- 在模型确定前，不把当前白名单描述为通用请求隔离。

### 3. `StreamedResponse` 不能通过当前发送器正确输出

位置：

- `std/symfony/http-foundation/streamed_response.go`
- `std/symfony/http-foundation/runtime.go` 的 `SendResponseTo`

`StreamedResponse` 通过 `sendContent()` 执行 callback 或输出 chunks，且 `getContent()` 按 Symfony 语义返回 `false`。当前 `SendResponseTo` 只读取普通 Response content / `getContent()`，没有调用 streamed callback，却仍统一设置 `Content-Length` 后发送。

结果可能是：

- `response()->stream()` 返回空 body。
- SSE、流式下载和生成器式响应不可用。
- 无法逐块 flush，违背流式响应的核心语义。

期望：

- 按 Response 类型分派发送策略。
- StreamedResponse 在写入 headers 后执行 callback/chunks，并允许增量 flush。
- 不为未知长度流提前设置 `Content-Length`。
- 增加 callback、chunks、异常、客户端断开和多次 send 的回归测试。

### 4. HTTP 请求结束时没有运行请求级 shutdown callbacks

位置：

- `runtime/request_vm.go` 的 `AddShutdownCallback` / `RunShutdownCallbacks`
- `std/laravel/serve/serve_command.go` 的 `ServeHTTP`

CLI 的 `finish()` 会调用 `RunShutdownCallbacks()`，但 HTTP 请求链路在成功、异常、超时和客户端断开路径上都没有对请求级 RequestVM 执行该阶段。

因此 `register_shutdown_function()` 可能不在 HTTP 请求结束时运行，清理、日志和追踪逻辑会丢失。

期望：

- 为每个请求持有明确的 RequestVM 引用，并用 `defer` 保证 shutdown callbacks 恰好运行一次。
- 定义正常返回、`exit` / `die`、PHP 异常、Go panic、超时和客户端断开的顺序。
- 确保 shutdown 输出与最终 HTTP body/已发送响应之间的行为符合 PHP SAPI 语义。

## P1：生命周期和并发风险

### 5. Kernel `terminate()` 没有执行 terminable middleware

位置：`std/laravel/httpkernel/methods.go` 的 `kernelTerminate`。

当前实现只在存在时调用 `$app->terminate()`，没有对齐官方 Kernel 的 `terminateMiddleware($request, $response)`。实现了终止接口的全局或路由中间件不会收到请求结束通知。

应在补全中间件 Pipeline 后，同时记录本次请求实际执行的中间件，并按 Laravel 13 顺序完成 terminate。

### 6. 客户端取消没有传播到 PHP 执行上下文

位置：`std/laravel/serve/serve_command.go` 的 `ServeHTTP`。

当前超时上下文由 `context.Background()` 创建，而不是继承 `r.Context()`。客户端断开后，PHP 代码仍可能执行到固定 30 秒上限，并继续产生数据库或外部系统副作用。

期望至少使用：

```go
context.WithTimeout(r.Context(), serveMaxExecutionTime)
```

同时检查数据库、文件、网络扩展以及 VM 热循环是否都能感知取消，而不是只在 PHP 方法调用边界检查。

### 7. 监听端口成功早于 Laravel readiness

位置：

- `std/laravel/serve/serve_command.go` 的 `runLaravelHTTPServer`
- `std/laravel/serve/serve_command.go` 的 `ensureBase`

Server 在监听成功后立即打印 running；真正的 Kernel resolve、bootstrap 和服务预热发生在首个动态请求中。

风险：

- TCP 探活成功但 Laravel 尚不可用。
- bootstrap 错误延迟到首个用户请求。
- 首请求承担不可控的冷启动开销。

期望在对外宣告 ready 前完成必要 bootstrap。若必须懒初始化，应区分存活检查和 readiness 检查。

### 8. 在线请求与异步 vendor 预热可能并发改变类加载状态

位置：`std/laravel/serve/serve_command.go` 中异步调用 `WarmupVendorClassmap` 的逻辑。

底层 map 并发安全不等于 PHP 类加载语义并发安全。请求 autoload 与预热同时处理同一类时，类注册顺序、顶层 PHP 副作用及错误暴露时机可能随调度变化。

更稳妥的边界是：

- 监听前完成会执行 PHP 顶层代码的预热；或
- 后台阶段只做无副作用的词法/解析缓存；或
- 对类加载建立单航班机制，保证同一文件和类只由一个加载过程执行。

### 9. 同名 Go 类抢占 vendor 类，兼容性失败时不能自动回退

位置：`std/vendoraccel`、`std/laravel`、`std/symfony` 的各类 `Load`。

原生层在 Composer autoload 前注册官方 FQCN。类一旦存在，Composer 就不会加载 vendor PHP 实现；某个边缘方法、签名、可见性、反射或序列化语义不完整时，不会自动回退，而会形成局部兼容类。

期望：

- 每个启用的原生类绑定精确 Composer 版本。
- 对公开方法、继承、接口、属性可见性、默认值、异常类型、反射和序列化建立契约测试。
- 未达到契约覆盖门槛的类不应 `AddClass`。
- vendor 升级时先运行契约差异检查，再调整 `TargetVersion`。

## P2：一致性和可维护性问题

### 10. Kernel 配置 API 与实际执行能力不一致

当前 Kernel 暴露了 `setGlobalMiddleware`、`prependMiddleware`、`pushMiddleware`、middleware group、alias 和 priority 等方法，但全局 middleware 没有被请求执行链消费。

相比直接报告“不支持”，这种表面成功更容易让应用在生产请求中静默缺少安全或业务逻辑。修复前应至少通过测试明确哪些 API 已生效。

### 11. 同一异常被重复写入 stderr

位置：`std/laravel/httpkernel/methods.go` 的 `renderException`。

该函数在解析异常对象前后调用了两次 `logHandleException(thrown)`，导致同一异常重复输出。应删除重复调用，并让 Laravel Exception Handler 负责主要 report 语义，Origami 只记录 Handler 自身不可用或失败的诊断信息。

### 12. `serve` 声明了没有落实的命令行选项

位置：

- `std/laravel/serve/serve_command.go` 的 `serveGetOptionsMethod`
- `std/laravel/serve/serve_command.go` 的 `serveAddress`

命令声明 `--tries` 和 `--no-reload`，解析和运行逻辑却只消费 host/port。当前既不会按 `--tries` 尝试后续端口，也没有 `.env` 自动重载语义。

选项应完整实现，或从命令签名中删除并明确与官方 `artisan serve` 的差异，避免接受参数但静默忽略。

### 13. 静态文件和 Livewire dist 特判侵入通用 Server

当前 ServeCommand 除了模拟 `public/` 静态文件，还硬编码了 Livewire hashed dist 路径。这能绕过尚未完善的 BinaryFileResponse 路径，但会把特定生态包兼容问题固化到通用 HTTP Server。

长期应优先修复 HttpFoundation/BinaryFileResponse 和路由响应语义。Livewire 特判只能作为有删除条件的临时兼容层，并需要单独记录触发原因与移除标准。

## 建议修复顺序

1. 补全官方 Kernel 的全局 Middleware Pipeline。
2. 补全 terminable middleware 和请求级 shutdown callbacks。
3. 修复 StreamedResponse，并覆盖 HEAD、204、304、BinaryFileResponse 等发送语义。
4. 将客户端取消传播到 VM 和阻塞扩展。
5. 明确 PHP-FPM 或 Octane 隔离模型，替换无法封闭证明的服务白名单策略。
6. 将 bootstrap/readiness 前移，收紧并发预热边界。
7. 建立原生替换类与锁定 vendor 版本的契约测试矩阵。
8. 清理重复日志、无效命令选项和包特定的 HTTP 特判。

## 验收原则

每项修复都应同时包含：

- `tests/php/` 中与框架无关的最小 PHP 语义回归（适用时）。
- `examples/laravel13/tests/origami/` 中真实 Laravel 生命周期回归。
- 至少一个并发请求用例，验证状态不串请求。
- HTTP 测试必须带超时。
- 不修改 `vendor/`、示例业务代码或 Blade 模板来绕过运行时缺口。

在上述 P0 项完成前，更准确的表述应是“Origami 已能运行当前 Laravel 13 + Filament 主路径”，而不是“已完整兼容 Laravel HTTP Kernel 或 PHP-FPM 请求语义”。
