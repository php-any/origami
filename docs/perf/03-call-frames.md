# 施工单 03：调用帧与参数绑定

## 任务卡

- **目标**：把每次函数调用的固定开销压下来——批量分配符号表、消除无 `static` 函数的热路径锁、静态调用去掉第二帧、`return` 去分配、内置函数签名常驻。
- **独占文件**：`runtime/context.go` 的 `CreateContext`/`resetVariables`/`syncStaticLocal`、`data/static_locals.go`、`data/value_return.go`、`data/value_class.go` 的 `ClassMethodContext` 构造、`node/function.go`、`node/return.go`、`node/call.go`、`node/call_fast.go`、`node/call_static_method.go`、`node/static_var.go`、`std/php/**` 的 `GetParams`/`GetVariables`
- **前置依赖**：[00-baseline-and-metrics.md](00-baseline-and-metrics.md)、**[02-array-hashing-and-cow.md](02-array-hashing-and-cow.md)**（两者都改 `runtime/context.go`，02 先合并）
- **对应基准**：`BenchmarkFuncCallSimple`、`BenchmarkFuncCallManyLocals`、`BenchmarkMethodCall`、`BenchmarkStaticCall`、`BenchmarkBuiltinCall`、`BenchmarkIssetEmpty`、`BenchmarkReturnValue`
- **必须新增的回归**：`tests/php/static_locals_test.php`、`tests/php/call_frame_ref_test.php`、`tests/php/static_call_context_test.php`
- **语义风险**：中高。Step 1 与 Step 4 都触碰帧的生命周期，做错会表现为跨调用串值或引用失效。

## 现状与证据

### 1. 每帧为每个槽单独 new 一个 ZVal

局部变量已经是编译期槽位模型（`parser/scope_manager.go` 的 `DefaultScope` 在解析期分配递增 `Index`），运行时是 `[]*data.ZVal` 下标访问，**这部分设计是对的**。问题在建帧：

```178:194:runtime/context.go
func (c *Context) resetVariables(vars []data.Variable) {
	n := len(vars)
	// 只复用 slice 头，每个槽分配新 ZVal。禁止原地改旧 ZVal：
	// 引用返回、闭包、数组元素可能仍持有上一帧的指针。
	if cap(c.variables) < n {
		c.variables = make([]*data.ZVal, n)
	} else {
		c.variables = c.variables[:n]
	}
	for i := 0; i < n; i++ {
		name := ""
		if vars[i] != nil {
			name = vars[i].GetName()
		}
		c.variables[i] = data.NewNamedZValSlot(name)
	}
}
```

`contextPool` 只回收 `Context` 壳。一个有 20 个局部变量的函数，每次调用 20 次堆分配，**即使函数体只用到 2 个**。Laravel 的控制器/中间件全是这种小函数。

### 2. 没有 `static` 变量的函数也在每次赋值加锁

`Call` 无条件绑定 static store：

```125:128:node/function.go
func (f *FunctionStatement) Call(ctx data.Context) (data.GetValue, data.Control) {
	if b, ok := ctx.(data.StaticLocalsBinder); ok {
		b.BindStaticLocals(f.funcStaticLocals())
	}
```

`funcStaticLocals()` 用 `sync.Once` 保证只建一次，但**总会返回一个非 nil 的 store**。于是 `staticLocals != nil` 恒真，每次赋值都要过：

```147:154:runtime/context.go
func (c *Context) syncStaticLocal(index int) {
	if c.staticLocals == nil {
		return
	}
	if zv := c.GetIndexZVal(index); zv != nil {
		c.staticLocals.Update(index, zv.Value)
	}
}
```

```35:45:data/static_locals.go
// Update 在 static 变量已注册后同步最新值（用于 ++/-- 等修改）。
func (s *StaticLocals) Update(index int, val Value) {
	if s == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Vals[index]; ok {
		s.Vals[index] = val
	}
}
```

**每一次局部变量赋值** = `sync.Mutex` 加锁 + map 查找（几乎总是 miss）+ `defer` 解锁。整个仓库里带 `static $x` 的 PHP 函数是极少数，但所有函数都在为它付费。这是局部变量路径上最不该存在的开销。

`ClassMethod.Call`（`node/class.go` 约 679 行）有同样的无条件绑定。

### 3. 静态调用走双帧，还造一个空对象

`Foo::bar($x)` 的解析结果是 `CallMethod{Method: CallStaticMethod, Args: ...}`（`parser/ident_parser.go` 与 `parser/variable_parser.go`）。`CallStaticMethod.GetValue` 每次返回一个新包装器，`CallMethod` 再调它：

```331:341:node/call_static_method.go
func (s *staticMethodFunc) Call(callCtx data.Context) (data.GetValue, data.Control) {
	// 创建类方法上下文，使用传入的 callCtx（包含已设置的参数），绑定当前类，保证 self:: 可用
	classValue := data.NewClassValue(s.class, callCtx)
	fnCtx := classValue.CreateContext(s.method.GetVariables())
	// 设置后期静态绑定类
	if cmc, ok := fnCtx.(*data.ClassMethodContext); ok {
		if s.callClass != nil {
			cmc.StaticClass = s.callClass
		}
```

`callCtx` 已经是 `CallMethod` 建好并绑定完参数的帧了，这里又建第二帧，然后把参数 ZVal 一个个拷过去：

```359:375:node/call_static_method.go
		params := s.method.GetParams()
		vars := s.method.GetVariables()
		for i := 0; i < len(vars); i++ {
			zval := callCtx.GetIndexZVal(i)
			if zval == nil {
				// ...默认值填充...
			} else {
				fnCtx.SetIndexZVal(i, zval)
			}
		}
```

而 `NewClassValue` 会造一个用不上的空对象：

```8:15:data/value_class.go
func NewClassValue(class ClassStmt, ctx Context) *ClassValue {
	return &ClassValue{
		ObjectValue: NewObjectValue(),   // 空 OrderedMap + 两个 map + 一把 RWMutex
		Class:       class,
		Context:     ctx,
	}
}
```

一次 `Foo::bar()` 的额外成本：1 个 `FuncValue` + 1 个 `staticMethodFunc` + 1 个 `ClassValue` + 1 个 `ObjectValue`（含 `OrderedMap` 的两个 map）+ 1 个池化帧（含 N 个 ZVal）+ 1 个 `ClassMethodContext` + N 次 ZVal 拷贝。Laravel 的 Facade 全是这条路径。

### 4. 每个 `return` 堆分配一个 Control

```3:8:data/value_return.go
func NewReturnControl(v Value) ReturnControl {
	return &ReturnValue{V: v}
}
```

```5:16:node/return.go
func (u *ReturnStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// ...
	return nil, data.NewReturnControl(v.(data.Value))
}
```

`FunctionStatement.Call` 收到后立刻 `rv.ReturnValue()` 拆出值，`ReturnValue` 随即变垃圾。对比 `break`/`continue`——它们让语句节点自己实现 `Control` 并 `return nil, u`，**零分配**。

### 5. 内置函数每次调用都重建签名

```37:49:std/php/strlen.go
func (f *StrlenFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
	}
}
func (f *StrlenFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
	}
}
```

`CallExpression.GetValue` 每次调用都先 `fn.GetVariables()` 和 `fn.GetParams()`：

```39:43:node/call.go
	varies := fn.GetVariables()
	params := fn.GetParams()
	fnCtx := ctx.CreateContext(varies)
```

于是每次 `strlen($s)` 分配：2 个 slice + 1 个 `Parameter` + 1 个 `Variable` + 1 个 `BaseType`。`std/php/` 下 200 多个文件都是这个模式。用户函数没有这个问题（`GetVariables` 返回已存字段）。

### 6. `empty`/`isset` 永远走慢路径

快路径判定要求形参都是 `*Parameter` 具体类型：

```11:20:node/call_fast.go
func canFastPositionalBind(params []data.GetValue, args []data.GetValue) bool {
	for _, p := range params {
		if _, ok := p.(*Parameter); !ok {
			return false
		}
	}
```

`empty`/`isset` 用 `*ParameterRawAST`（为了拿到未求值的 AST，抑制未定义变量错误），所以永远进慢路径——命名参数扫描、多次 slice 分配。而 `isset` 是 PHP 里最高频的构造之一。

## Step 1：符号表批量分配

### 1.1 思路：不复用，改批量

`resetVariables` 的注释解释了为什么不能复用旧 ZVal（引用返回、闭包、数组元素可能持有上一帧指针）。**这个约束是对的，不要试图绕过它。**

但"不能复用"不等于"必须一个个分配"。把 n 次小对象分配换成 1 次连续分配：

```go
func (c *Context) resetVariables(vars []data.Variable) {
	n := len(vars)
	if n == 0 {
		c.variables = c.variables[:0]
		return
	}
	if cap(c.variables) < n {
		c.variables = make([]*data.ZVal, n)
	} else {
		c.variables = c.variables[:n]
	}
	// 一次连续分配 n 个 ZVal，再把指针填进槽位。
	// 仍然是全新的 ZVal 对象（不复用上一帧的），引用安全性与逐个 new 完全等价；
	// 区别只是内存连续、分配次数从 n 降到 1。
	block := make([]data.ZVal, n)
	nullV := data.NewNullValue()
	for i := 0; i < n; i++ {
		z := &block[i]
		if vars[i] != nil {
			z.Name = vars[i].GetName()
		}
		z.Value = nullV
		c.variables[i] = z
	}
}
```

`data.NewNullValue()` 返回 intern 单例，提到循环外只是省接口调用。

**这一步风险极低**：每个槽拿到的仍是独立的新 `ZVal`，所有引用语义不变。唯一的行为差异是内存布局——如果某个槽逃逸（被闭包/生成器持有），整个 `block` 无法被 GC 回收。这是空间换时间，PHP 的 CV 数组也是连续布局。对于持有大量局部变量的长寿闭包，内存会略有放大，用 `BenchmarkClosureCall` 的 `B/op` 和 Laravel 的 RSS 峰值确认可接受。

### 1.2 第二阶段（可选）：条件复用

如果 1.1 之后 `BenchmarkFuncCallManyLocals` 仍显示分配是瓶颈，可以进一步做"未逃逸则复用 block"：给 `Context` 加逃逸位图，`GetIndexZVal` 之外新增一个显式的 `GetIndexZValForRef` 用于取引用的场景并标记逃逸位，归还池时只复用全未逃逸的 block。

**这一步需要审计全部 `GetIndexZVal` 调用点**，区分"读值"和"取引用"，工作量和风险都显著更高。**建议先只做 1.1，测量后再决定是否需要 1.2。** 不要在同一次提交里做两步。

## Step 2：`static` 惰性绑定

### 2.1 核心思路

不在 `Call` 开头无条件绑定，改成**执行到 `static $x` 语句时才绑定**。没有 `static` 声明的函数，`c.staticLocals` 保持 nil，`syncStaticLocal` 第一行就 return，零锁。

`node/function.go`：

```go
func (f *FunctionStatement) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 不再无条件 BindStaticLocals。
	// static 局部变量的 store 由 StaticVarStatement 在首次执行时绑定，
	// 这样没有 static 声明的函数帧 staticLocals 保持 nil，赋值路径不需要加锁。
	if f.IsGenerator { ... }
```

`node/static_var.go` 的 `StaticVarStatement.GetValue` 开头加：

```go
func (s *StaticVarStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 惰性绑定：只有真正执行到 static 声明的帧才挂 store
	if b, ok := ctx.(data.StaticLocalsBinder); ok {
		if b.StaticLocalsStore() == nil {
			b.BindStaticLocals(s.ownerStore())
		}
	}
	// ...原有逻辑...
}
```

`StaticVarStatement` 需要能拿到所属函数的 store。两种做法：

- **A（推荐）**：解析期把 `*FunctionStatement`（或它的 `funcStaticLocals` 闭包）注入 `StaticVarStatement`。`parser` 在解析函数体时已知当前函数，注入是自然的。
- **B**：让 `StaticVarStatement` 自己持有一个 `*data.StaticLocals`（每个 `static` 声明一个 store，而不是每个函数一个）。语义上等价——`static $x` 的持久化范围就是这条声明——而且更简单，不需要改 parser。**但要确认 `persistStaticLocals` 的行为**：它在函数结束时扫全部 vars 做 `Update`，如果一个函数有多条 `static` 声明分属不同 store，需要绑定的是"当前帧看到的那个 store"。若一个函数体里有多条 `static` 声明，B 方案下只有第一条的 store 会被绑定到帧上——**这会破坏语义**。

所以选 **A**。

### 2.2 闭包分支也要改

`FunctionStatement.Call` 里闭包场景会重新建帧并再次绑定：

```153:156:node/function.go
			if b, ok := execCtx.(data.StaticLocalsBinder); ok {
				b.BindStaticLocals(f.funcStaticLocals())
			}
```

这段同样删掉——`StaticVarStatement` 会在新帧上按需绑定。

`ClassMethod.Call`（`node/class.go`）的无条件绑定同样删除。

### 2.3 `persistStaticLocals` 保持不变

函数结束时的 `f.persistStaticLocals(execCtx)` 内部应该先检查 store 是否为 nil。确认它在 `staticLocals == nil` 时是 no-op；如果不是，加上检查。

### 2.4 `StaticLocals` 的锁保留

不要因为"单请求单线程"就删掉 `StaticLocals.mu`。Origami 支持 `std/channel` 和 goroutine，同一函数可能被多个 goroutine 并发调用，static 变量是真正的共享状态。惰性绑定已经让绝大多数函数完全不碰这把锁，锁本身可以留着。

## Step 3：静态调用去掉第二帧

### 3.1 让 `ClassMethodContext` 能包住已有帧

`staticMethodFunc.Call` 收到的 `callCtx` 已经是绑定好参数的帧。需要的只是"在这个帧外面套一层提供 `self::`/`static::` 的包装"，而不是新建帧。

`data/value_class.go` 新增：

```go
// NewStaticMethodContext 为静态方法调用包装已有帧，不新建符号表、不造空对象。
// ObjectValue 为 nil：静态方法没有 $this，访问 $this 应当报错而不是拿到空对象。
func NewStaticMethodContext(inner Context, self ClassStmt, static ClassStmt) *ClassMethodContext {
	return &ClassMethodContext{
		ClassValue: &ClassValue{
			ObjectValue: nil,
			Class:       self,
			Context:     inner,
		},
		SelfClass:   self,
		StaticClass: static,
	}
}
```

`staticMethodFunc.Call` 简化为：

```go
func (s *staticMethodFunc) Call(callCtx data.Context) (data.GetValue, data.Control) {
	if s.method.GetName() == "__callStatic" {
		return s.callStatic(callCtx) // 参数重排路径保持原样，见 3.3
	}
	cmc := data.NewStaticMethodContext(callCtx, s.class, s.callClass)
	return s.method.Call(cmc)
}
```

省掉：1 个 `ObjectValue`（含 `OrderedMap` 的两个 map 和一把 RWMutex）、1 个池化帧、N 个 ZVal、N 次 ZVal 拷贝、以及原来那段默认值补齐循环（`callCtx` 已经由 `CallMethod` 绑定完参数并补过默认值）。

### 3.2 必须处理 `ObjectValue == nil`

`ClassMethodContext` 现在假设 `ObjectValue` 非 nil：

```352:357:data/value_class.go
func (c *ClassMethodContext) GetVariableValue(variable Variable) (Value, Control) {
	if _, ok := variable.(Property); ok {
		return c.ObjectValue.GetVariableValue(variable)
	}
	return c.Context.GetVariableValue(variable)
}
```

改为：

```go
	if _, ok := variable.(Property); ok {
		if c.ObjectValue == nil {
			// 静态上下文里访问实例属性：PHP 会 Fatal error: Using $this when not in object context
			return nil, NewErrorThrow(nil, errors.New("Using $this when not in object context"))
		}
		return c.ObjectValue.GetVariableValue(variable)
	}
```

必须把 `ClassMethodContext` 与 `ClassValue` 上**所有**解引用 `ObjectValue` 的方法都过一遍：

```bash
rg 'c\.ObjectValue\.|\.ObjectValue\.' data/value_class.go
```

每一处都要处理 nil。这是本 Step 的主要工作量，也是主要风险点。

**注意**：原来 `NewClassValue` 造空对象可能掩盖了某些"静态方法里误用 `$this`"的场景——那些代码现在会开始报错。这是**向 PHP 对齐**，但可能让 Laravel/Livewire 里原本"侥幸工作"的路径暴露出来。跑 `examples/laravel13` 时要留意新出现的 `Using $this when not in object context`；如果出现，先确认 PHP 官方在同样代码上也报错，再决定是修 Origami 的调用方判定还是接受报错。

### 3.3 `__callStatic` 路径保持原样

`__callStatic` 需要把实参重排成 `[$name, [$args]]`，签名与调用方传入的参数不匹配，**确实需要新建帧**。把这段抽成独立方法 `callStatic`，走原来的 `NewClassValue` + `CreateContext` 逻辑，不做优化。它不是热路径。

### 3.4 `instanceViaSelfFunc` 一并检查

`node/call_static_method.go` 里还有 `instanceViaSelfFunc`（支持实例方法里 `self::nonStaticMethod()`）。检查它是否有同样的双帧问题，若有则同法处理；它有真实的 `$this`，不涉及 nil ObjectValue。

## Step 4：`return` 去分配

### 4.1 每帧一个可复用的 ReturnValue

`ReturnValue` 的生命周期极短：`ReturnStatement` 创建 → 沿 AST 向上冒泡 → `Call` 拆开取值 → 丢弃。它从不跨帧存活。

给 `Context` 加一个内嵌的返回槽：

```go
// runtime/context.go
type Context struct {
	// ...
	// retSlot 是本帧复用的 return 载体。ReturnControl 只在本帧内向上冒泡，
	// 被 Call 拆值后即失效，因此每帧一个实例就够，不需要每次 return 都分配。
	retSlot data.ReturnValue
}

func (c *Context) ReturnSlot(v data.Value) data.ReturnControl {
	c.retSlot.V = v
	return &c.retSlot
}
```

`Context` 接口加 `ReturnSlot(Value) ReturnControl`，`node/return.go` 改为：

```go
	return nil, ctx.ReturnSlot(v.(data.Value))
```

`ClassMethodContext` / `BoundContext` 等包装器透传给内层。

### 4.2 风险与验证

危险场景是"`ReturnValue` 被保存到超出本帧生命周期的地方"。要审计：

```bash
rg 'ReturnControl|ReturnValue\b' --type go -g '!vendor' -g '!*_test.go'
```

重点看：

- `try`/`finally`：`finally` 里若有自己的 `return`，会覆盖同一帧的 `retSlot`。**这正好符合 PHP 语义**（`finally` 的 `return` 覆盖 `try` 的），但要确认实现不是先保存 try 的 Control 再执行 finally 再比较——如果是，两者会指向同一个对象。**这是本 Step 最可能出错的地方**，必须看 `node/try.go` 的实现后再决定。
- 生成器：`YieldControl` 是另一套，不受影响。但生成器帧被 `MarkEscaped` 后不还池，`retSlot` 也随帧存活，安全。
- 异常栈：`AddStack` 处理的是 `ThrowValue`，不是 `ReturnValue`。

**如果 `try`/`finally` 的实现确实会同时持有两个 Control，放弃本 Step，改用 `sync.Pool` 池化 `ReturnValue`**（由 `Call` 在拆值后归还）。收益略低但没有别名风险。

在 `-tags origamidebug` 下可以加一个校验：`ReturnSlot` 被同一帧连续调用两次且上一次尚未被消费时 panic。

## Step 5：内置函数签名常驻化

### 5.1 模式

以 `strlen` 为例：

```go
type StrlenFunction struct {
	params []data.GetValue
	vars   []data.Variable
}

func NewStrlenFunction() data.FuncStmt {
	return &StrlenFunction{
		params: []data.GetValue{node.NewParameter(nil, "string", 0, nil, nil)},
		vars:   []data.Variable{node.NewVariable(nil, "string", 0, data.NewBaseType("string"))},
	}
}

func (f *StrlenFunction) GetParams() []data.GetValue    { return f.params }
func (f *StrlenFunction) GetVariables() []data.Variable { return f.vars }
```

内置函数在 `std/php/load.go` 里每个只 `New` 一次并注册进 VM，所以签名对象是进程级共享的。

### 5.2 前置确认：`Parameter` 与 `Variable` 必须无可变状态

共享的前提是这些对象不可变。`Parameter.SetValue` 只写 `ctx`，`Parameter.GetValue` 读 `DefaultValue`，看起来无状态。**但必须逐项确认**：

```bash
rg 'func \(p \*Parameter\)' node/function.go
rg 'func \(v \*Variable\)' node/variable.go
```

任何写 `p.xxx = ` 或 `v.xxx = ` 的方法都是阻塞点。`Variable` 若持有 `value` 字段并在 `SetValue` 里写自己，共享就会串值——**这一步必须先查清，不能假设**。

用户函数的 `GetVariables()` 已经返回共享的 `f.vars` 字段，说明共享模式在现有设计下是成立的。这是一个有力的旁证，但仍要确认内置用的 `NewParameter`/`NewVariable` 构造出的对象与用户函数用的是同一套。

### 5.3 工作量与做法

`std/php/` 下 200+ 文件。**不要手工改 200 个文件**。做法：

1. 先改 10 个最热的（`strlen`、`count`、`isset`、`empty`、`is_array`、`is_string`、`in_array`、`implode`、`explode`、`array_key_exists`），用 `BenchmarkBuiltinCall` 确认收益。
2. 收益确认后写一个改写脚本（放 `scripts/` 下，参考已有的 `scripts/migrate_illuminate_http.py` 的字符串改写风格），批量把 `GetParams`/`GetVariables` 的返回体提到构造函数。
3. 脚本改完必须 `go build ./...` + 全量测试，并人工 review diff——`std/php/` 里有些函数的签名是动态的（参数个数依赖调用），这些不能机械改写。

**收益排序清楚的话，第 1 步就可能拿到大部分收益**，第 2 步可以作为独立提交甚至后续任务。

## Step 6：让 `isset`/`empty` 进快路径

`canFastPositionalBind` 拒绝非 `*Parameter` 的形参。`*ParameterRawAST` 的绑定逻辑（把未求值 AST 包成 `ASTValue` 塞进槽）本身很简单，可以在快路径里支持：

```go
func canFastPositionalBind(params []data.GetValue, args []data.GetValue) bool {
	for _, p := range params {
		switch p.(type) {
		case *Parameter, *ParameterRawAST:
			// 这两种都能按位置直接绑定
		default:
			return false
		}
	}
	for _, a := range args {
		switch a.(type) {
		case *NamedArgument, *SpreadArgument, *CallerContextParameter:
			return false
		}
	}
	return true
}
```

`bindPositionalParameters` 里对应加分支：`*ParameterRawAST` 时**不求值实参**，直接包 AST 写入槽。注意它现在是先 `arg.GetValue(ctx)` 再绑定——`ParameterRawAST` 的整个存在意义就是不求值，所以快路径里必须先判类型再决定是否求值：

```go
	for i, arg := range args {
		if i < len(params) {
			if raw, ok := params[i].(*ParameterRawAST); ok {
				// 不求值，保留 AST（isset/empty 依赖它抑制未定义变量错误）
				if acl := raw.SetRawAST(fnCtx, arg); acl != nil {
					return acl
				}
				continue
			}
		}
		v, acl := arg.GetValue(ctx)
		// ...原有逻辑...
	}
```

`flat` / `SetFlatCallArgs` 对 raw 参数要写什么，取决于 `func_get_args` 在 `isset` 上的语义——`isset` 不是真函数，不支持 `func_get_args`，写 nil 或跳过都可以，但要和慢路径保持一致。

### 顺带：`SetCallArgs`/`SetFlatCallArgs` 的无条件分配

```go
	flat := make([]data.Value, 0, len(args))
	// ...
	fnCtx.SetCallArgs(args)
	fnCtx.SetFlatCallArgs(flat)
```

`flat` 只被 `func_get_args()` / `func_num_args()` 消费，绝大多数调用不需要。可以改成惰性：`Context` 保存 `args []GetValue` 和调用方 `ctx`，`func_get_args` 被调用时才求值成 `flat`。

**但要小心**：`func_get_args` 在函数体中间调用时，PHP 返回的是**调用时传入的原始值**，不是当前的参数变量值（参数被修改后 `func_get_args` 仍返回原值）。惰性求值会返回"调用 `func_get_args` 那一刻"的值——**语义不一致**。所以惰性化只有在"实参表达式求值结果被缓存"的前提下才安全，而那又需要分配。

结论：`flat` 的分配可以优化成"只在函数体内出现 `func_get_args`/`func_num_args` 时才做"，由解析期在 `FunctionStatement` 上打标记。**这需要 AST 扫描，作为可选项，不在本单的必做范围。**

## 语义风险与 PHP 对照

| 改动 | 风险 | 缓解 |
|---|---|---|
| ZVal 批量分配 | 槽逃逸时整个 block 不被 GC | 测 RSS；每槽仍是独立新对象，引用语义不变 |
| static 惰性绑定 | `static` 变量跨调用丢值 | `tests/php/static_locals_test.php` 覆盖多条声明、条件分支、递归、闭包 |
| 静态调用去第二帧 | `ObjectValue == nil` 引发 panic | 逐个方法加 nil 检查；`$this` 误用改为按 PHP 报错 |
| 静态调用去第二帧 | 参数默认值不再被补齐 | 确认 `CallMethod` 侧已补齐；回归覆盖带默认值的静态方法 |
| `ReturnSlot` 复用 | `try`/`finally` 双 Control 别名 | 先读 `node/try.go`；有风险则改用 `sync.Pool` |
| 内置签名共享 | `Parameter`/`Variable` 若可变会串值 | 5.2 的前置确认是硬性门槛 |
| `isset` 进快路径 | 实参被提前求值，破坏未定义变量抑制 | 快路径里先判类型再决定求值 |

PHP 对照要点：`static` 变量的持久化范围是"声明所在的函数"，递归调用共享同一份；`self::`/`static::` 的后期静态绑定必须在无 `$this` 的上下文里也正确；静态方法里访问 `$this` 是 Fatal error。

## 回归测试

### `tests/php/static_locals_test.php`

```php
<?php
function counter() { static $n = 0; $n++; return $n; }
assertEq(counter(), 1); assertEq(counter(), 2); assertEq(counter(), 3);

// 多条 static 声明
function multi() {
    static $a = 1;
    static $b = 10;
    $a++; $b += 2;
    return [$a, $b];
}
assertEq(multi(), [2, 12]);
assertEq(multi(), [3, 14]);

// 条件分支里的 static（首次调用不执行到声明）
function cond($run) {
    if ($run) { static $s = 0; $s++; return $s; }
    return -1;
}
assertEq(cond(false), -1);
assertEq(cond(true), 1);
assertEq(cond(true), 2);

// 递归共享同一份
function rec($d) { static $calls = 0; $calls++; if ($d > 0) rec($d - 1); return $calls; }
assertEq(rec(3), 4);

// 无 static 的函数不受影响（本单主要优化目标）
function plain($x) { $y = $x * 2; $z = $y + 1; return $z; }
assertEq(plain(5), 11);

// 方法里的 static
class C {
    public function m() { static $k = 0; $k++; return $k; }
    public static function sm() { static $j = 100; $j++; return $j; }
}
$c1 = new C(); $c2 = new C();
assertEq($c1->m(), 1);
assertEq($c2->m(), 2);   // PHP：非静态方法的 static 变量在同类所有实例间共享
assertEq(C::sm(), 101);
assertEq(C::sm(), 102);

// 闭包里的 static
$f = function() { static $q = 0; $q++; return $q; };
assertEq($f(), 1); assertEq($f(), 2);

// static 数组
function acc($v) { static $arr = []; $arr[] = $v; return count($arr); }
assertEq(acc('a'), 1); assertEq(acc('b'), 2);
```

### `tests/php/call_frame_ref_test.php`

覆盖帧生命周期与引用逃逸：

```php
<?php
// 引用参数
function addOne(&$x) { $x++; }
$n = 5; addOne($n); assertEq($n, 6);

// 引用返回
class Box { public $v = 1; public function &ref() { return $this->v; } }
$b = new Box(); $r = &$b->ref(); $r = 99; assertEq($b->v, 99);

// 闭包按值捕获
function makeVal() { $x = 1; return function() use ($x) { return $x; }; }
$g = makeVal(); assertEq($g(), 1);

// 闭包按引用捕获，且父帧已返回
function makeRef() { $x = 1; $inc = function() use (&$x) { $x++; return $x; }; return $inc; }
$h = makeRef();
assertEq($h(), 2); assertEq($h(), 3);

// 多个闭包共享同一父帧变量
function makePair() {
    $n = 0;
    return [function() use (&$n) { $n++; }, function() use (&$n) { return $n; }];
}
list($inc, $get) = makePair();
$inc(); $inc();
assertEq($get(), 2);

// 生成器持有帧
function gen() { $local = 'kept'; yield 1; yield $local; }
$it = gen();
$vals = [];
foreach ($it as $v) { $vals[] = $v; }
assertEq($vals, [1, 'kept']);

// 递归中各帧独立
function fib($n) { if ($n < 2) return $n; return fib($n-1) + fib($n-2); }
assertEq(fib(10), 55);

// 大量局部变量（对应 BenchmarkFuncCallManyLocals）
function manyLocals() {
    $a=1;$b=2;$c=3;$d=4;$e=5;$f=6;$g=7;$h=8;$i=9;$j=10;
    return $a+$j;
}
assertEq(manyLocals(), 11);

// return 在 try/finally 中（对应 Step 4 的风险点）
function tf() { try { return 'try'; } finally { } }
assertEq(tf(), 'try');
function tf2() { try { return 'try'; } finally { return 'finally'; } }
assertEq(tf2(), 'finally');
function tf3() { try { throw new Exception('e'); } catch (Exception $ex) { return 'catch'; } finally { } }
assertEq(tf3(), 'catch');

// 嵌套 return
function outer() { $r = inner(); return $r . '-outer'; }
function inner() { return 'inner'; }
assertEq(outer(), 'inner-outer');
```

### `tests/php/static_call_context_test.php`

覆盖 Step 3 的静态调用上下文：

```php
<?php
class Base {
    public static $shared = 'base';
    const K = 'base-k';
    public static function who() { return 'Base'; }
    public static function callWho() { return static::who(); }   // 后期静态绑定
    public static function selfWho() { return self::who(); }     // 词法绑定
    public static function withDefault($a, $b = 'def') { return "$a|$b"; }
    public static function readConst() { return self::K; }
    public static function readStatic() { return static::$shared; }
}
class Child extends Base {
    public static $shared = 'child';
    const K = 'child-k';
    public static function who() { return 'Child'; }
}

assertEq(Base::who(), 'Base');
assertEq(Child::who(), 'Child');
assertEq(Child::callWho(), 'Child');    // static:: 解析到 Child
assertEq(Child::selfWho(), 'Base');     // self:: 解析到 Base
assertEq(Base::withDefault('x'), 'x|def');
assertEq(Base::withDefault('x', 'y'), 'x|y');
assertEq(Child::readConst(), 'base-k'); // self::K 是词法的
assertEq(Child::readStatic(), 'child'); // static::$shared 是动态的

// parent::
class P { public static function f() { return 'P'; } }
class Q extends P { public static function f() { return 'Q+' . parent::f(); } }
assertEq(Q::f(), 'Q+P');

// 静态方法里用 $this 应当报错（Step 3.2 的行为变化）
class R {
    public $prop = 'v';
    public static function bad() { return isset($this); }
}
assertEq(R::bad(), false);   // PHP: isset($this) 在静态上下文为 false，不是 Fatal

// __callStatic 仍然工作
class M {
    public static function __callStatic($name, $args) { return $name . ':' . implode(',', $args); }
}
assertEq(M::anything('a', 'b'), 'anything:a,b');

// 静态方法作为回调
assertEq(call_user_func(['Base', 'who']), 'Base');
assertEq(call_user_func('Base::who'), 'Base');

// 实例上调用静态方法
$c = new Child();
assertEq($c::who(), 'Child');
```

`isset($this)` 那条用例值得注意：PHP 里静态上下文的 `isset($this)` 是 `false` 而不是 Fatal error，只有**实际使用** `$this->x` 才 Fatal。Step 3.2 的 nil 处理要能区分这两种情况——如果做不到，先让 `isset($this)` 走通比让它报错更重要。

### 运行

```bash
go run ./zy.go tests/php/static_locals_test.php
go run ./zy.go tests/php/call_frame_ref_test.php
go run ./zy.go tests/php/static_call_context_test.php
go run ./zy.go tests/run_tests.php
```

若 02 的 `origamidebug` 自检已就位，也用 debug 构建跑一遍。

## 完成定义

- [ ] `resetVariables` 改批量分配，`BenchmarkFuncCallManyLocals` 的 `allocs/op` 显著下降
- [ ] `static` 改惰性绑定，`FunctionStatement.Call` / `ClassMethod.Call` / 闭包分支的无条件 `BindStaticLocals` 全部移除
- [ ] memprofile / CPU profile 确认 `StaticLocals.Update` 从热路径消失
- [ ] `staticMethodFunc.Call` 不再新建帧、不再造空 `ObjectValue`
- [ ] `ClassMethodContext` / `ClassValue` 所有解引用 `ObjectValue` 的方法都处理了 nil
- [ ] `__callStatic` 路径独立且行为不变
- [ ] `return` 去分配（`ReturnSlot` 或 `sync.Pool`），`try`/`finally` 语义已验证
- [ ] 至少 10 个最热内置函数签名常驻化，`Parameter`/`Variable` 不可变性已确认
- [ ] `isset`/`empty` 进入快路径，且未定义变量抑制行为不变
- [ ] 三份 `tests/php/` 回归通过
- [ ] `tests/run_tests.php` 全绿，`examples/laravel13` 能起并响应首页
- [ ] `benchstat` 前后对比归档，`BenchmarkStaticCall` 改善明显
