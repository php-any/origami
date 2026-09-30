# Laravel 13 运行时设计审查

本文记录 `examples/laravel13` 当前常驻 HTTP 运行模型中已确认的语义缺口和高风险设计，作为后续对齐官方 Laravel 13 / PHP 请求生命周期的修复清单。

本文只描述问题，不代表对应能力已经修复。判断标准是：官方 Laravel 应用不应为了适配 Origami 而修改 `vendor/`、业务代码或隐藏错误；通用差异应回推到 Origami 核心、标准库或 Laravel/Symfony 原生兼容层。

## 当前模型

当前 `serve` 不是传统的 `php -S` 或 PHP-FPM 模型，而是：

```text
Go net/http 常驻进程
  -> Origami 基础 VM 与常驻 Laravel Application
  -> 每请求 TempVM + Kernel/Application/Router 沙箱
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
  Context / VM   -> 当前请求 TempVM
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

- `runtime/vm_temp.go` 的 `AddShutdownCallback` / `RunShutdownCallbacks`
- `std/laravel/serve/serve_command.go` 的 `ServeHTTP`

CLI 的 `finish()` 会调用 `RunShutdownCallbacks()`，但 HTTP 请求链路在成功、异常、超时和客户端断开路径上都没有对请求级 TempVM 执行该阶段。

因此 `register_shutdown_function()` 可能不在 HTTP 请求结束时运行，清理、日志和追踪逻辑会丢失。

期望：

- 为每个请求持有明确的 TempVM 引用，并用 `defer` 保证 shutdown callbacks 恰好运行一次。
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
