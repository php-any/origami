# 施工单 01：标量装箱与字符串操作

## 任务卡

- **目标**：消除标量值的重复堆分配，去掉字符串路径上的 `reflect` 与 `fmt`。
- **独占文件**：`data/value_intern.go`、`data/value_int.go`、`data/value_string.go`、`data/value_float.go`、`data/value_bool.go`、`data/scalar_fast.go`、`node/binary_dot.go`、`node/binary_link.go`、`node/string_literal.go`、`node/fused_assign.go`、`std/php/html_string_funcs.go`
- **前置依赖**：[00-baseline-and-metrics.md](00-baseline-and-metrics.md)
- **对应基准**：`BenchmarkIntArithmetic`、`BenchmarkStringConcat`、`BenchmarkStringInterpolation`、`BenchmarkStringLiteralRead`
- **必须新增的回归**：`tests/php/value_intern_test.php`、`tests/php/string_concat_semantics_test.php`
- **语义风险**：中。intern 引入**共享可变对象**的风险，Step 1 必须先做完 Step 0 的排雷。

## 现状与证据

### 只有 bool 和 null 被 intern

```1:7:data/value_intern.go
package data

var (
	internNull      = &NullValue{}
	internBoolTrue  = &BoolValue{Value: true}
	internBoolFalse = &BoolValue{Value: false}
)
```

`NewIntValue` / `NewFloatValue` / `NewStringValue` 每次 `new`：

```7:11:data/value_int.go
func NewIntValue(v int) Value {
	return &IntValue{
		Value: v,
	}
}
```

整数字面量在解析期分配一次并缓存在 AST 上（`node/number_literal.go` 的 `IntLiteral.V`），所以读字面量不分配。**但字符串字面量没有这个待遇**：

```174:177:node/string_literal.go
func (s *StringLiteral) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(s.Value), nil
}
```

循环里每次读 `"hello"` 都是一个新 `*StringValue`。

算术每次分配：`node/binary_add.go` 的 `phpNumericAdd` 结果走 `NewIntValue`/`NewFloatValue`。`strlen` 返回值也是新 `*IntValue`。

### `IntValue.AsString` 走 `fmt.Sprintf`

```26:28:data/value_int.go
func (s *IntValue) AsString() string {
	return fmt.Sprintf("%d", s.Value)
}
```

`fmt.Sprintf("%d", n)` 比 `strconv.Itoa(n)` 慢一个量级（反射 + 接口装箱 + 格式串解析）。整数转字符串在 echo、拼接、数组键、`implode` 上非常热。`BoolValue.AsString` 同样用 `fmt.Sprint`。

### `BinaryDot` 用 reflect 按类型名解包

```30:43:node/binary_dot.go
	// 使用反射检查 ZVal 类型（因为 ZVal 没有实现 GetValue 接口）
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		// 检查是否是 ZVal 类型（通过类型名称）
		if rv.Type().String() == "*data.ZVal" {
			// 获取 Value 字段
			valueField := rv.Elem().FieldByName("Value")
			if valueField.IsValid() && !valueField.IsNil() {
				if val, ok := valueField.Interface().(data.GetValue); ok {
					return unwrapValue(val)
				}
			}
		}
	}
```

每次字符串拼接都做：`reflect.ValueOf`（接口→反射值）、`rv.Type().String()`（**构造类型名字符串并比较**）、`FieldByName("Value")`（**按名字查结构体字段**）。这三步在字符串拼接热路径上，是全仓库最贵的单点之一。

而且这段逻辑本身是多余的——`*data.ZVal` 之所以"没有实现 GetValue"，只是因为没给它写方法；同文件下面几行已经有 `*data.ZValValue` 的正常类型断言分支。

同一函数里整数拼接还走 `fmt.Sprintf`：

```99:101:node/binary_dot.go
	case *data.IntValue:
		leftStr = fmt.Sprintf("%d", l.Value)
```

### 插值是 `BinaryLink` 链，中间串反复分配

`"a{$x}b{$y}c"` 被解析成嵌套 `BinaryLink`，每个节点：

```35:44:node/binary_link.go
	leftStr, lCtl := ValueToDisplayString(ctx, lv)
	// ...
	rightStr, rCtl := ValueToDisplayString(ctx, rv)
	// ...
	return data.NewStringValue(leftStr + rightStr), nil
```

n 段插值产生 n-1 次中间字符串分配 + n-1 个 `*StringValue`。Blade 模板渲染是这条路径的重灾区。

### `ArrayValue.AsString` 是二次方拼接

```148:159:data/value_array.go
func (a *ArrayValue) AsString() string {
	str := "["
	for _, zval := range a.List {
		str = str + zval.Value.AsString() + ", "
	}
```

每次 `str = str + ...` 重新分配整串。1000 元素数组 = 1000 次分配、总拷贝量 O(n²)。顺带一提，PHP 里数组的字符串化结果是 `"Array"`（并发 Notice），不是 `[...]`——但**这属于语义问题，本单不改**，只把二次方拼接换成 `strings.Builder`，保持现有输出不变。

## Step 0：排雷 —— 找出所有原地修改标量的代码

**这一步必须先做完，否则 Step 1 会引入极难定位的串值 bug。**

intern 的前提是被共享的对象**不可变**。当前有两处原地修改 `IntValue.Value`：

```221:225:node/fused_assign.go
func (f *VarStmtIncr) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if zv := ctx.GetIndexZVal(f.VarIdx); zv != nil {
		if iv, ok := zv.Value.(*data.IntValue); ok {
			iv.Value++ // 原地改写，0 次分配
```

```334:337:std/php/html_string_funcs.go
			if iv, ok := percentRef.(*data.FloatValue); ok {
				iv.Value = percent
			} else if iv, ok := percentRef.(*data.IntValue); ok {
				iv.Value = int(percent)
```

如果 intern 了小整数而不改这两处，会发生：`for ($i = 0; $i < 10; $i++)` 里的 `$i++` 把 intern 表中代表 `0` 的对象改成 `1`，此后整个进程里所有 `0` 都变成 `1`。

**同时要重新扫一遍确认没有新增的原地写入点**：

```bash
rg '\.Value\+\+|\.Value--|\.Value \+= |\.Value -= ' --type go
rg 'IntValue\)\s*;?\s*\w+\.Value = |FloatValue\)\s*;?\s*\w+\.Value = ' --type go
```

处理办法：

- `VarStmtIncr`：改成 `zv.Value = data.NewIntValue(iv.Value + 1)`。**这不会变慢**——Step 1 完成后小整数走 intern 表，`NewIntValue` 对小整数是零分配，正好替代掉原地改写。循环变量通常在 intern 范围内，超出范围时退化为一次分配（和自增前的语义相同）。
- `html_string_funcs.go`：按引用写回不应该原地改 Value 对象，应该替换 ZVal 里的 Value。改为取 `*data.ZVal` 后 `zv.Value = data.NewFloatValue(percent)`。若该处拿不到 ZVal，退路是保留原地修改但**先调用 `data.SeparateScalar(v)`**（见 Step 1.3）。

## Step 1：标量 intern

### 1.1 intern 表

改写 `data/value_intern.go`：

```go
package data

const (
	internIntMin = -1024
	internIntMax = 4096
)

var (
	internNull      = &NullValue{}
	internBoolTrue  = &BoolValue{Value: true}
	internBoolFalse = &BoolValue{Value: false}

	// internInts 覆盖 [internIntMin, internIntMax]。
	// 这些对象被全局共享，必须视为不可变：禁止任何代码原地修改 .Value。
	internInts [internIntMax - internIntMin + 1]IntValue

	internEmptyString = &StringValue{Value: ""}

	// internSingleByteStrings 覆盖单字节字符串 "\x00".."\xff"。
	// 字符串下标 $s[$i]、chr()、逐字符扫描会大量命中。
	internSingleByteStrings [256]StringValue
)

func init() {
	for i := range internInts {
		internInts[i].Value = i + internIntMin
	}
	for i := range internSingleByteStrings {
		internSingleByteStrings[i].Value = string([]byte{byte(i)})
	}
}
```

用**数组而不是 map**：数组是一次连续分配，索引访问无哈希开销；`&internInts[idx]` 取到的指针在数组生命周期内稳定。

范围选 `[-1024, 4096]` 的依据：覆盖绝大多数循环计数器、数组下标、HTTP 状态码、`strlen` 结果、小尺寸计算。范围应在 00 的 `BenchmarkIntArithmetic` 与 Laravel profile 上验证后微调——**不要凭感觉扩到十万级**，那会让 init 分配 MB 级内存且命中率提升有限。

### 1.2 工厂函数

```go
func NewIntValue(v int) Value {
	if v >= internIntMin && v <= internIntMax {
		return &internInts[v-internIntMin]
	}
	return &IntValue{Value: v}
}

func NewStringValue(s string) Value {
	switch len(s) {
	case 0:
		return internEmptyString
	case 1:
		return &internSingleByteStrings[s[0]]
	}
	return &StringValue{Value: s}
}
```

`NewFloatValue` **不 intern**。理由：浮点相等语义有 NaN/±0 的坑，而 `node/binary_eq.go` 有指针相等短路（见下），intern NaN 会让 `NAN == NAN` 错误地返回 true。

### 1.3 给"必须可变"的场景留后门

少数地方确实需要独占的标量对象（按引用参数写回、未来可能的原地优化）。提供显式分离函数：

```go
// SeparateScalar 返回一个保证不被共享的标量副本，用于需要原地修改的场景。
// 正常求值路径不应调用它——共享 intern 对象是本文件的设计前提。
func SeparateScalar(v Value) Value {
	switch s := v.(type) {
	case *IntValue:
		return &IntValue{Value: s.Value}
	case *FloatValue:
		return &FloatValue{Value: s.Value}
	case *StringValue:
		return &StringValue{Value: s.Value}
	}
	return v
}
```

### 1.4 确认与已有指针相等短路兼容

`node/binary_eq.go` 和 `binary_ne.go` 开头有：

```33:35:node/binary_eq.go
	if lv == rv {
		return data.NewBoolValue(true), nil
	}
```

intern 让这条短路的命中率大幅上升，**且结果正确**：指针相等 ⇒ 同一对象 ⇒ 值相等 ⇒ PHP 的 `==` 为 true。这是 intern 的额外收益，不是风险。

`node/binary_eq_strict.go` 的 `isStrictEqual` 按值比较：

```45:50:node/binary_eq_strict.go
	switch v1 := value1.(type) {
	case *data.IntValue:
		if v2, ok2 := value2.(*data.IntValue); ok2 {
			return v1.Value == v2.Value
		}
```

不受 intern 影响。

`std/php/array/array_search.go` 的 `valuesIdentical` 对 `*FuncValue`/`*ClassValue` 用指针比较，对标量按值比较——也不受影响。

**但要检查一处**：`$a == $a` 当 `$a` 是 NaN 时，指针短路已经会错误返回 true。这是**既存缺陷**，本单不负责修，但要在 PR 描述里记录，避免被误认为 intern 引入。

## Step 2：字符串字面量缓存到 AST

`StringLiteral` 与 `IntLiteral` 保持一致——解析期分配一次，求值直接返回。

```go
type StringLiteral struct {
	*Node `pp:"-"`
	Value string
	v     data.Value // 解析期构造，求值直接复用
}

func (s *StringLiteral) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if s.v == nil { // 兼容通过结构体字面量直接构造的节点（cmd/compile 生成的代码）
		s.v = data.NewStringValue(s.Value)
	}
	return s.v, nil
}
```

**为什么这是安全的**：`StringValue` 在 Step 0 之后是不可变的，字符串赋值走 `SetVariableValue` 的 `default` 分支只做指针赋值，不会改到共享对象。这与 `IntLiteral` 已经在做的事完全同构。

`nil` 兜底那一行是必需的：`cmd/compile` 生成的 Go 代码是用结构体字面量直接构造节点（`&node.StringLiteral{...}`），不走 `NewStringLiteral` 构造函数，`v` 会是 nil。同时也要在 `cmd/compile` 的 special handler 里确认新加的私有字段不会被 emit（`gen.go` 只处理可导出字段，但要跑一次 `zy compile` 验证）。

## Step 3：`BinaryDot` 去 reflect

### 3.1 给 `ZVal` 加 `GetValue` 方法，删掉 reflect 分支

reflect 那段的注释说明了动机："因为 ZVal 没有实现 GetValue 接口"。直接给它实现：

```go
// data/zval.go
func (z *ZVal) GetValue(ctx Context) (GetValue, Control) {
	if z == nil || z.Value == nil {
		return NewNullValue(), nil
	}
	return z.Value, nil
}
```

**先确认这不会引起歧义**：`ZVal` 变成 `GetValue` 后，所有 `switch v.(type)` 里有 `data.GetValue` 分支的地方可能改变匹配结果。必须搜一遍：

```bash
rg 'case data\.GetValue|\.\(data\.GetValue\)' --type go
```

如果风险不可控，退路是**不**给 `ZVal` 加方法，而是把 reflect 换成直接类型断言——`unwrapValue` 的参数类型是 `data.GetValue`，`*data.ZVal` 能被传进来只可能是通过 `any`/接口转换，直接断言即可：

```go
func unwrapValue(v data.GetValue) data.GetValue {
	for {
		switch t := v.(type) {
		case *data.ZValValue:
			if t == nil || t.ZVal == nil || t.ZVal.Value == nil {
				return v
			}
			v = t.ZVal.Value
		case *data.ReferenceValue:
			actual, _ := t.Val.GetValue(t.Ctx)
			v = actual
		case *data.IndexReferenceValue:
			actual, _ := t.Expr.GetValue(t.Ctx)
			v = actual
		default:
			return v
		}
	}
}
```

注意这里同时把递归改成了循环——`unwrapValue` 原本是递归的，深层引用链会有栈开销。

**退路方案更安全，建议优先选它**：给 `ZVal` 加接口方法会影响全仓库的类型断言，收益却只是省掉一个分支。

### 3.2 整数拼接改 `strconv`

`binary_dot.go` 里两处 `fmt.Sprintf("%d", l.Value)` → `strconv.Itoa(l.Value)`。

### 3.3 用 `strings.Builder` 一次成型

左右两侧都转成 string 后，`leftStr + rightStr` 会分配一次结果串。这一次分配是必需的，但可以用预估容量避免 Go 内部的二次扩容：

```go
var sb strings.Builder
sb.Grow(len(leftStr) + len(rightStr))
sb.WriteString(leftStr)
sb.WriteString(rightStr)
return data.NewStringValue(sb.String()), nil
```

对单次拼接收益有限。**真正的收益在 Step 4 的链式折叠**，这里保持简单的 `+` 也可接受。

## Step 4：插值链折叠

`"a{$x}b{$y}c"` 目前是 `BinaryLink(BinaryLink(BinaryLink("a", $x), "b"), ...)`，产生 n-1 次中间串。

改法：在解析期把连续的 `BinaryLink` 折叠成一个多元节点。

```go
// node/binary_link.go
// InterpolationExpr 表示一段字符串插值，避免 BinaryLink 链的中间字符串分配。
type InterpolationExpr struct {
	*Node `pp:"-"`
	Parts []data.GetValue
}

func (e *InterpolationExpr) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 先求值全部片段，再一次性拼接
	strs := make([]string, len(e.Parts))
	total := 0
	for i, p := range e.Parts {
		v, ctl := p.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		s, ctl := ValueToDisplayString(ctx, v)
		if ctl != nil {
			return nil, ctl
		}
		strs[i] = s
		total += len(s)
	}
	var sb strings.Builder
	sb.Grow(total)
	for _, s := range strs {
		sb.WriteString(s)
	}
	return data.NewStringValue(sb.String()), nil
}
```

**求值顺序必须严格从左到右**，因为片段里可能有 `__toString()` 或函数调用带副作用。上面的实现满足这一点：先按顺序全部求值，再拼接。

构造点在 `parser/parser.go` 的 `parseLingToken`（它现在产出 `BinaryLink` 链）。改成收集全部片段后一次性构造 `InterpolationExpr`。**保留 `BinaryLink` 类型不删**，因为 `cmd/compile` 生成的代码和可能的其他构造点还在引用它。

`strs` 切片本身是一次分配。若要更进一步可以复用 `sync.Pool`，但先测量——`BenchmarkStringInterpolation` 会告诉你值不值得。

## Step 5：`AsString` 全面去 `fmt`

- `data/value_int.go`：`fmt.Sprintf("%d", s.Value)` → `strconv.Itoa(s.Value)`
- `data/value_bool.go`：`fmt.Sprint(...)` → 直接返回常量 `"1"` / `""`（**先确认现有返回值**：PHP 里 `(string)true === "1"`、`(string)false === ""`。如果当前实现返回的是 `"true"`/`"false"`，那是个既存语义缺陷，**本单不改**，只把 `fmt` 换成常量并保持相同输出，在 PR 里记录该缺陷）
- `data/value_array.go` 的 `AsString`：换 `strings.Builder`，输出格式保持不变

`data/value_float.go` 的 `strconv.FormatFloat(..., 'g', 14, 64)` 已经在用 strconv，不改。

## Step 6：验证分配确实消失

改完跑：

```bash
go test ./bench/ -bench='BenchmarkIntArithmetic|BenchmarkStringConcat|BenchmarkStringInterpolation|BenchmarkStringLiteralRead' -benchmem -count=10 > after.txt
benchstat before.txt after.txt
```

再用 memprofile 确认：

```bash
$env:ORIGAMI_MEMPROFILE="mem-after.out"; ./origami.exe bench/testdata/string_heavy.php
go tool pprof -sample_index=alloc_objects -top mem-after.out
```

`data.NewIntValue` 和 `data.NewStringValue` 应该从 top 里消失或大幅下降。如果没有，说明命中的整数超出了 intern 范围，回到 1.1 调整边界。

## 语义风险与 PHP 对照

| 改动 | 风险 | 缓解 |
|---|---|---|
| 小整数 intern | 任何原地修改 `IntValue.Value` 的代码会污染全局 | Step 0 排雷 + `SeparateScalar` 后门 + 代码注释标注不可变 |
| 单字节字符串 intern | 同上 | 同上 |
| 字符串字面量缓存 | 若某处原地改 `StringValue.Value`，同一字面量在所有位置被改 | Step 0 的扫描要覆盖 `StringValue` |
| `ZVal` 加 `GetValue` | 影响全仓库类型断言匹配 | 优先选 3.1 的退路方案（不加方法） |
| 插值折叠 | 求值顺序、副作用时机 | 严格左到右，`tests/php/string_concat_semantics_test.php` 覆盖 |
| `AsString` 换实现 | 输出格式若有差异会引起大面积 diff | 只换实现不换输出，先跑全量回归 |

不引入 PHP 语义变化——PHP 本身也有 interned string 和小整数优化，本单是在向 PHP 靠拢。

## 回归测试

新建 `tests/php/value_intern_test.php`，覆盖 intern 的共享安全性：

```php
<?php
// 自增不能污染其他变量
$a = 5;
$b = 5;
$a++;
if ($a !== 6 || $b !== 5) { /* 报错 */ }

// 循环变量自增不能污染字面量
for ($i = 0; $i < 3; $i++) {
    $lit = 0;
    if ($lit !== 0) { /* 报错：intern 表被污染 */ }
}

// 数组元素独立性
$arr = [1, 1, 1];
$arr[0]++;
if ($arr[1] !== 1 || $arr[2] !== 1) { /* 报错 */ }

// 引用自增
$x = 10;
$r = &$x;
$r++;
if ($x !== 11) { /* 报错 */ }

// intern 边界外
$big = 999999;
$big2 = 999999;
$big++;
if ($big2 !== 999999) { /* 报错 */ }

// 单字节字符串
$s = "abc";
$c1 = $s[0];
$c2 = "a";
if ($c1 !== "a" || $c2 !== "a") { /* 报错 */ }

// 严格比较与松散比较
if (!(5 === 5) || !(0 == "0") || (0 === "0")) { /* 报错 */ }
```

新建 `tests/php/string_concat_semantics_test.php`，覆盖拼接与插值的语义：

```php
<?php
// 各类型拼接
$parts = [1 . "", 1.5 . "", true . "", false . "", null . ""];
// 期望: "1", "1.5", "1", "", ""

// 插值求值顺序与副作用
$log = [];
function mark($x) { global $log; $log[] = $x; return $x; }
$s = "" . mark("a") . mark("b") . mark("c");
// 期望 $log === ["a","b","c"]

$n = 0;
function bump() { global $n; $n++; return $n; }
$s2 = "x{$GLOBALS['n']}y";

// __toString
class S { public function __toString(): string { return "S!"; } }
$o = new S();
if ("v={$o}" !== "v=S!") { /* 报错 */ }
if ("v=" . $o !== "v=S!") { /* 报错 */ }

// 引用与 ZVal 解包
$a = "A"; $r = &$a;
if (("p" . $r) !== "pA") { /* 报错 */ }

// 数组元素拼接
$arr = ["k" => "V"];
if (("q" . $arr["k"]) !== "qV") { /* 报错 */ }

// .= 复合赋值
$acc = "";
for ($i = 0; $i < 3; $i++) { $acc .= $i; }
if ($acc !== "012") { /* 报错 */ }
```

按仓库惯例，失败用 `Log::fatal` 输出（参考 `tests/php/` 下现有文件的写法），这样 `tests/run_tests.php` 能捕获。

运行：

```bash
go run ./zy.go tests/php/value_intern_test.php
go run ./zy.go tests/php/string_concat_semantics_test.php
```

## 完成定义

- [ ] Step 0 的扫描已执行，`VarStmtIncr` 与 `html_string_funcs.go` 的原地修改已消除
- [ ] intern 表实现，边界经 benchmark 验证
- [ ] `NewIntValue` / `NewStringValue` 走 intern，`NewFloatValue` 明确不 intern 且注释说明原因
- [ ] `SeparateScalar` 提供并注释"正常路径不应调用"
- [ ] `StringLiteral` 缓存值，含 nil 兜底，`zy compile` 验证通过
- [ ] `binary_dot.go` 无 `reflect` 引用，`unwrapValue` 改为循环
- [ ] `fmt.Sprintf("%d")` / `fmt.Sprint` 在本单独占文件里全部消除
- [ ] 插值折叠为 `InterpolationExpr`，求值顺序左到右
- [ ] `ArrayValue.AsString` 用 `strings.Builder`，输出不变
- [ ] 两份 `tests/php/` 回归通过
- [ ] `tests/run_tests.php` 全绿，`examples/laravel13` 能起
- [ ] `benchstat` 前后对比已归档到 `bench/results/`
- [ ] memprofile 证明 `NewIntValue`/`NewStringValue` 从 top 分配点消失
