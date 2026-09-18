# 施工单 02：数组哈希化与真 COW

## 任务卡

- **目标**：把数组键查找从 O(n) 降到 O(1)；把变量赋值从 eager clone 换成真正的写时复制。
- **独占文件**：`data/value_array.go`、`data/cow.go`、`node/array.go`、`node/index.go`、`runtime/context.go` 的 `SetVariableValue`
- **前置依赖**：[00-baseline-and-metrics.md](00-baseline-and-metrics.md)
- **后续冲突**：施工单 03 也要改 `runtime/context.go`（`CreateContext`/`resetVariables`）。**本单必须先合并**。
- **对应基准**：`BenchmarkArrayDenseIndexRead`、`BenchmarkArrayStringKeyRead`、`BenchmarkArrayAssignLarge`、`BenchmarkArrayPassByValue`、`BenchmarkForeachLarge`
- **必须新增的回归**：`tests/php/array_key_lookup_test.php`、`tests/php/array_cow_test.php`
- **语义风险**：**高**。这是全部 8 份施工单里风险最高的一份。COW 分离时机做错的表现是"改 $b 影响了 $a"，而且往往只在特定嵌套深度下出现。Step 0 的自检机制不是可选项。

## 现状与证据

### 1. `ArrayValue` 是线性表，整数键查找先 strconv 再全表扫

```101:107:data/value_array.go
type ArrayValue struct {
	List     []*ZVal
	iterator int // 迭代器当前位置索引
	// IndirectOverloadClass 非空表示该数组来自 ArrayAccess::offsetGet 的副本，对其元素的间接修改无效
	IndirectOverloadClass string
	rc                    int // 指向该容器的 zval 数（copy-on-write）
}
```

键存在 `ZVal.Name` 里：空 `Name` 表示密集整数下标，`"6"` / `"foo"` 表示显式键。查找：

```279:293:data/value_array.go
// FindSlotByIntKey 按 PHP 整数键查找槽位（含稀疏键 Name=="6" 等）
func (a *ArrayValue) FindSlotByIntKey(i int) (*ZVal, int) {
	keyStr := IntArrayKeyName(i)
	for j, z := range a.List {
		if z != nil && z.Name == keyStr {
			return z, j
		}
	}
	if i >= 0 && i < len(a.List) {
		if z := a.List[i]; z != nil && z.Name == "" {
			return z, i
		}
	}
	return nil, -1
}
```

读一个整数键要：一次 `strconv.Itoa`（`IntArrayKeyName`）+ 全表字符串比较，只有全表 miss 后才试 `List[i]` 快路径。字符串键同理：

```265:277:data/value_array.go
func (a *ArrayValue) LookupZValByStringKey(key string) (*ZVal, bool) {
	for _, zval := range a.List {
		if zval != nil && zval.Name == key {
			return zval, true
		}
	}
	if i, ok := ParseIntArrayKeyName(key); ok {
		if z, _ := a.FindSlotByIntKey(i); z != nil {
			return z, true
		}
	}
	return nil, false
}
```

**注意：不能简单地把 `FindSlotByIntKey` 里两段的顺序对调。** 看起来"先试 `List[i]` 再扫 Name"更快，但会错：

```php
$a = [1, 2, 3];   // List: [""], [""], [""]
$a[6] = 'x';      // List 追加 Name=="6" 的槽，len=4
$a[] = 'p';       // PHP 下一个键是 7
$a[] = 'q';       // 8
$a[] = 'r';       // 9，此时 len=7
echo $a[6];       // PHP 输出 'x'
```

反转顺序后 `List[6]` 是键 9 的槽（Name==""），会返回 `'r'`。这也正是 `normalizeDenseIntKeys` 存在的原因。**正确的 O(1) 快路径必须先知道"这个数组里没有任何显式键"**，见 Step 1。

### 2. `rc` COW 已经实现，但变量赋值没用

`data/cow.go` 有完整的一套：

```30:41:data/cow.go
// CowAssign 把值写入 zval：释放旧数组引用并 AddRef 新数组（标量/对象按指针赋值）。
func CowAssign(zv *ZVal, v Value) {
	if zv == nil {
		return
	}
	zv.Defined = true
	if zv.Value == v {
		return
	}
	CowRelease(zv.Value)
	zv.Value = CowAddRef(v)
}
```

```43:57:data/cow.go
func CowSeparateZVal(zv *ZVal) {
	if zv == nil {
		return
	}
	switch a := zv.Value.(type) {
	case *ArrayValue:
		if a.rc > 1 {
			a.rc--
			cloned := CloneArrayValue(a)
			cloned.rc = 1
			zv.Value = cloned
		} else if a.rc == 0 {
			a.rc = 1
		}
```

但变量赋值路径完全绕开它，直接 eager clone：

```123:133:runtime/context.go
	case *data.ArrayValue:
		zv := c.variables[variable.GetIndex()]
		zv.Value = data.CloneArrayValue(v)
		zv.Defined = true
	case *data.ObjectValue:
		// PHP 中 array 是按值赋值 + copy-on-write。
		zv := c.variables[variable.GetIndex()]
		zv.Value = data.CloneObjectValue(v)
		zv.Defined = true
```

于是 `$b = $a`（1000 元素数组）立刻分配 1 个 `ArrayValue` + 1 个 1000 长的 slice + 1000 个 `ZVal`，即使 `$b` 从头到尾只读。

`CloneArrayValue` 还有个不一致：它返回的 `ArrayValue` **没设 `rc`**（默认 0），而 `CloneArrayValueForCallArgs` 设 `rc: 1`。接入 COW 时必须统一。

### 3. `.List` 是导出字段，被 748 处引用

```
rg '\.List' --type go -g '!vendor'          # 748 处引用，169 个文件
rg '\.List\s*=|\.List\[[^\]]*\]\s*=|append\([^,]*\.List' --type go -g '!vendor'
                                            # 149 处写入，49 个文件
```

`GetMethod` 甚至把 `&a.List` 的地址交出去：

```167:171:data/value_array.go
	case "push":
		return &ArrayValuePush{&a.List}, true
	case "pop":
		return &ArrayValuePop{&a.List}, true
```

**这是本单的核心工程约束**：任何"需要 149 个写入点配合维护"的索引方案都不可行。索引必须是**可自检的派生状态**——能自己发现自己脏了。

## Step 0：先建自检机制

**这一步先做，否则后面两个 Step 都无法验证正确性。**

新增 `data/array_debug.go`，用构建标记隔离，**默认不参与编译，零开销**：

```go
//go:build origamidebug

package data

const arrayDebugChecks = true

// checkIndexConsistency 在每次索引命中后做一次全表校验，不一致直接 panic。
// 只在 -tags origamidebug 下编译，用于跑全量测试套件时暴露漏掉的失效点。
func (a *ArrayValue) checkIndexConsistency(key string, gotIdx int) {
	want := -1
	for j, z := range a.List {
		if z != nil && z.Name == key {
			want = j
			break
		}
	}
	if want != gotIdx {
		panic(fmt.Sprintf("array 查找索引不一致: key=%q 索引给出 %d 实际 %d，说明有写入点未失效索引", key, gotIdx, want))
	}
}

// checkSeparated 在写入数组结构前校验该数组不是共享状态。
// 用于暴露漏掉 CowSeparateZVal 的写入点。
func (a *ArrayValue) checkSeparated(where string) {
	if a.rc > 1 {
		panic(fmt.Sprintf("在 rc=%d 的共享数组上直接写入（%s），缺少 CowSeparateZVal", a.rc, where))
	}
}
```

配一份 `data/array_debug_off.go`：

```go
//go:build !origamidebug

package data

const arrayDebugChecks = false

func (a *ArrayValue) checkIndexConsistency(string, int) {}
func (a *ArrayValue) checkSeparated(string)             {}
```

编译器会把空函数内联掉，生产构建零成本。

自检的使用方式贯穿全单：

```bash
# 每完成一个 Step 都跑一遍
go build -tags origamidebug -o origami-debug.exe .
./origami-debug.exe tests/run_tests.php
```

**如果 debug 构建下测试套件不 panic，才有资格认为这一步做对了。**

## Step 1：惰性查找索引

### 1.1 数据结构

索引全部是派生状态，加在 `ArrayValue` 尾部：

```go
type ArrayValue struct {
	List                  []*ZVal
	iterator              int
	IndirectOverloadClass string
	rc                    int

	// --- 以下为惰性查找索引，纯派生状态，可随时丢弃重建 ---
	// idxMap 只收录 Name != "" 的槽：Name -> List 下标
	idxMap map[string]int
	// idxAllDense 为 true 表示所有槽 Name == ""，此时整数键 i 就是 List[i]
	idxAllDense bool
	// 以下三个字段构成结构指纹，用于检测 List 被外部改动
	idxLen   int
	idxFirst *ZVal
	idxLast  *ZVal
	idxBuilt bool
}
```

### 1.2 指纹校验与构建

```go
// lookupIndex 返回可用的索引，必要时重建。
func (a *ArrayValue) lookupIndex() {
	if a.idxBuilt && a.idxLen == len(a.List) {
		n := len(a.List)
		if n == 0 {
			return
		}
		if a.idxFirst == a.List[0] && a.idxLast == a.List[n-1] {
			return // 指纹匹配，索引可用
		}
	}
	a.rebuildIndex()
}

func (a *ArrayValue) rebuildIndex() {
	n := len(a.List)
	allDense := true
	var m map[string]int
	for j, z := range a.List {
		if z == nil {
			continue
		}
		if z.Name != "" {
			allDense = false
			if m == nil {
				m = make(map[string]int, n)
			}
			// 重复 Name 时保留第一个，与原线性扫描（首次命中即返回）一致
			if _, exists := m[z.Name]; !exists {
				m[z.Name] = j
			}
		}
	}
	a.idxMap = m
	a.idxAllDense = allDense
	a.idxLen = n
	if n > 0 {
		a.idxFirst = a.List[0]
		a.idxLast = a.List[n-1]
	} else {
		a.idxFirst, a.idxLast = nil, nil
	}
	a.idxBuilt = true
}

// InvalidateLookupIndex 供内部 mutator 与外部直接操作 List 的代码调用。
func (a *ArrayValue) InvalidateLookupIndex() {
	a.idxBuilt = false
	a.idxMap = nil
}
```

**"重复 Name 时保留第一个"这条必须严格遵守**，否则会与现有线性扫描行为不一致。理论上不该有重复键，但 149 个写入点里未必都保证了这一点，索引不能比原实现更严格。

### 1.3 指纹能捕获什么、不能捕获什么

能捕获（这三类占绝大多数）：

- `append` 导致 `len` 变化
- 整体替换 `a.List = newList`（`len` 或首尾指针变化）
- 切片截断、`array_shift`/`array_pop`（`len` 变化）

**不能**捕获：

- 原地替换单槽 `a.List[i] = NewZVal(v)`，且 `i` 不是 0 或 len-1
- 原地改键名 `a.List[i].Name = "x"`

这两类必须显式失效。扫描命令：

```bash
# 原地替换中间槽
rg '\.List\[[^\]]*\]\s*=' --type go -g '!vendor'
# 原地改键名
rg '\.Name\s*=' --type go -g '!vendor'
```

`data/value_array.go` 内部已知的两处：

```312:317:data/value_array.go
	if slot := a.List[i]; slot != nil && slot.RefSlotCount > 0 {
		slot.Value = value
	} else {
		a.List[i] = NewZVal(value)
	}
```

```319:326:data/value_array.go
// normalizeDenseIntKeys 将 Name=="" 的连续槽位转为显式整数字符串键
func (a *ArrayValue) normalizeDenseIntKeys() {
	for j, z := range a.List {
		if z != nil && z.Name == "" {
			z.Name = IntArrayKeyName(j)
		}
	}
}
```

`normalizeDenseIntKeys` 会把整个数组从 dense 变成全 named，**必须在结尾调 `InvalidateLookupIndex()`**。`SetIntKey` 的 `a.List[i] = NewZVal(value)` 分支同样要失效（`Name` 都是 ""，dense 性质不变，但首尾指针可能没变而槽对象换了——保守起见直接失效）。

外部（`std/php/array/`、`std/php/spl/`、`std/symfony/` 等 48 个文件）的写入点：**逐个检查，只有改中间槽或改 Name 的才需要加失效调用**。绝大多数是整体重建 List 或 append，指纹会自动捕获。Step 0 的 debug 构建 + 全量测试套件是找出漏网之鱼的手段。

### 1.4 改写查找函数

```go
func (a *ArrayValue) FindSlotByIntKey(i int) (*ZVal, int) {
	a.lookupIndex()

	// 快路径：纯密集数组，键 i 就是 List[i]，无需 strconv、无需扫表
	if a.idxAllDense {
		if i >= 0 && i < len(a.List) {
			if z := a.List[i]; z != nil {
				return z, i
			}
		}
		return nil, -1
	}

	// 有显式键：先查显式键（保持与原实现相同的优先级）
	keyStr := IntArrayKeyName(i)
	if a.idxMap != nil {
		if j, ok := a.idxMap[keyStr]; ok {
			if arrayDebugChecks {
				a.checkIndexConsistency(keyStr, j)
			}
			return a.List[j], j
		}
	}
	// 再回退到密集位置
	if i >= 0 && i < len(a.List) {
		if z := a.List[i]; z != nil && z.Name == "" {
			return z, i
		}
	}
	return nil, -1
}
```

注意快路径里 `idxAllDense` 时**跳过了 `IntArrayKeyName`**——省掉的 `strconv.Itoa` 本身就是可观的收益。

`LookupZValByStringKey` 同构改写：先 `a.idxMap[key]`，miss 再走 `ParseIntArrayKeyName` + `FindSlotByIntKey` 的既有回退。

### 1.5 索引不该在小数组上构建

对 3 个元素的数组建 map 是净亏。加阈值：

```go
const arrayIndexThreshold = 8

func (a *ArrayValue) lookupIndex() {
	if len(a.List) < arrayIndexThreshold {
		a.idxBuilt = false // 小数组不建索引，走线性扫描
		return
	}
	// ...
}
```

但 `idxAllDense` 的快路径对小数组也有价值（省 strconv）。所以更好的做法是：**`idxAllDense` 的判定始终做**（一次 O(n) 扫描，可缓存），只有 `idxMap` 受阈值限制。阈值的具体值由 `BenchmarkArrayDenseIndexRead` 在不同规模（4/8/16/64/1000 元素）下测出的交叉点决定，不要凭感觉定。

## Step 2：接入真 COW

### 2.1 `SetVariableValue` 换成 `CowAssign`

`runtime/context.go`：

```go
	case *data.ArrayValue:
		data.CowAssign(c.variables[variable.GetIndex()], v)
	case *data.ObjectValue:
		data.CowAssign(c.variables[variable.GetIndex()], v)
```

`CowAssign` 会 `rc++` 而不复制。`$b = $a` 变成一次指针赋值 + 一次整数自增。

### 2.2 统一 `rc` 初值

`CloneArrayValue` 与 `NewArrayValue` 都要明确 `rc`：

```go
func NewArrayValue(v []Value) Value {
	list := make([]*ZVal, len(v))
	for i, val := range v {
		list[i] = NewZVal(val)
	}
	return &ArrayValue{List: list} // rc = 0：字面量/新建，尚未被任何 zval 持有
}
```

`CloneArrayValue` 保持 `rc` 为 0 还是 1？`CowSeparateZVal` 里克隆后显式设 `cloned.rc = 1`，所以 `CloneArrayValue` 本身返回 0 是对的，由调用方决定。**但 `CloneArrayValueForCallArgs` 返回 `rc: 1`**，这个不一致要么统一、要么在注释里说清为什么不同。建议统一为 0，由调用方设置。

`data/cow.go` 里 `rc == 0` 被当作"字面量，首次写入时提升到 1"处理。这个约定必须写进 `ArrayValue.rc` 的字段注释里，否则后来者会看不懂 `else if a.rc == 0 { a.rc = 1 }` 这段。

### 2.3 找齐所有写入点并插入分离

这是本单工作量最大的部分。所有"修改数组结构或元素"的位置在写之前必须 `CowSeparateZVal`（或对 `*ArrayValue` 直接用 `SeparateArrayValue`）。

已有的基础设施可以复用：

```69:78:data/cow.go
// CowSeparateIndex 对调用帧第 index 个参数做写前分离（array_push 等引用参数）。
func CowSeparateIndex(ctx Context, index int) Value {
	zv := ctx.GetIndexZVal(index)
	if zv == nil {
		v, _ := ctx.GetIndexValue(index)
		return v
	}
	CowSeparateZVal(zv)
	return zv.Value
}
```

`array_push` 等按引用接收数组的内置函数已经在用它。要处理的是：

1. **`node/index.go` 的写入路径**：`$a[k] = v`、`$a[] = v`、`$a[k][j] = v`、`unset($a[k])`、`$a[k]++`、`$a[k] .= x`。这些都要在拿到目标数组后、写入前分离。嵌套写入（`$a['x']['y'] = 1`）要**沿路径逐层分离**，这是最容易漏的地方。
2. **`node/array.go`**：数组字面量构造是新建对象，`rc` 为 0，不需要分离。
3. **`std/php/array/` 下按引用修改数组的函数**：`array_push`/`pop`/`shift`/`unshift`/`splice`/`sort` 系列/`shuffle` 等。检查每一个是否走了 `CowSeparateIndex`。
4. **`ArrayValue.GetMethod` 交出的 `&a.List`**：`push`/`pop`/`shift`/`unshift`/`splice`/`sort` 这些 Origami 方言方法会直接改 List。它们拿不到 ZVal，无法分离。**处理办法**：这些方法只在 `.zy` 方言里用，不在 PHP 语义路径上；先在方法实现处加 `a.checkSeparated("ArrayValue.push")`，让 debug 构建告诉你是否真的会在共享数组上被调用。若确实会，需要把这些方法改为接收 `*ZVal` 而不是 `*[]*ZVal`。

### 2.4 移除现在的"半吊子 COW"

`node/unset.go` 有一段不看 `rc`、只看父表达式类型就无条件克隆的逻辑：

```141:155:node/unset.go
func cowSeparateNestedArray(arrayExpr data.GetValue, container data.GetValue) data.GetValue {
	if _, ok := arrayExpr.(*IndexExpression); !ok {
		return container
	}
	switch v := container.(type) {
	case *data.ArrayValue:
		return data.CloneArrayValue(v)
	case *data.ObjectValue:
		return data.CloneObjectValue(v)
```

真 COW 接上后，这段应该改为 `data.SeparateArrayValue(v)` / `SeparateObjectValue(v)`——只在 `rc > 1` 时才复制。**这是本单的净收益之一**：嵌套写入不再无条件复制。

但要小心：这段代码存在说明当时遇到过实际 bug。改动后必须确认相关测试仍通过（`rg 'cowSeparateNestedArray' -A5` 找到它的调用点，看注释里提到的场景）。

### 2.5 `ObjectValue` 也走同一套

关联数组在 Origami 里可能由 `ObjectValue` 表示（`$_GET`、从 null vivify 出来的 `$a['k']=1`、json 对象）。`cow.go` 已经对 `*ObjectValue` 做了对称处理，写入点同样要分离。`ObjectValue` 的属性存储是 `OrderedMap`，其克隆成本比 `ArrayValue` 更高（双 map 重建），COW 收益更大。

**注意与施工单 07 的边界**：07 改 `OrderedMap` 的内部结构与锁，本单只改 `ObjectValue` 的 COW 时机。两者理论上不冲突，但若 07 先合并，本单需要 rebase。

## Step 3：字面量稀疏键不再填洞

```118:127:node/array.go
func setArrayLiteralIntKey(av *data.ArrayValue, i int, val data.Value) {
	if i < 0 {
		setArrayLiteralStringKey(av, data.IntArrayKeyName(i), val)
		return
	}
	for len(av.List) <= i {
		av.List = append(av.List, data.NewZVal(data.NewNullValue()))
	}
	av.List[i] = data.NewZVal(val)
}
```

`[10000 => 1]` 会分配 10001 个槽。PHP 的稀疏数组是哈希表，只占 1 个。改为与运行时 `SetIntKey` 一致的稀疏语义（追加 `Name` 为 `"10000"` 的槽）。

**这会改变可观察行为**：`count()` 结果、`foreach` 顺序、`var_dump` 输出都会变——而且是**变得和 PHP 一致**。所以这一步既是性能修复也是语义修复。必须单独提交，配 `tests/php/array_key_lookup_test.php` 里的稀疏键用例。

如果发现有测试依赖当前的填洞行为，那些测试本身就是错的（与 PHP 不符），应该修测试而不是回退这一步。

## 语义风险与 PHP 对照

| 改动 | 风险 | 缓解 |
|---|---|---|
| 惰性索引 | 漏失效点 → 读到错误的槽 | Step 0 的 `origamidebug` 全表校验 + 全量测试套件 |
| 索引重复键处理 | 与线性扫描优先级不一致 | 严格"保留第一个" |
| `CowAssign` 替代 eager clone | 漏分离点 → 改 `$b` 影响 `$a` | Step 0 的 `checkSeparated` + 专项回归 |
| 嵌套写入分离 | 只分离最外层 → 内层仍共享 | 沿索引链逐层分离，回归覆盖三层嵌套 |
| `&$arr[i]` 引用槽 | 分离时必须保持引用槽共享 | `CloneArrayValue` 已按 `RefSlotCount > 0` 共享，不要动这段逻辑 |
| `ArrayValue.GetMethod` 的 `&a.List` | 绕过分离 | `checkSeparated` 探测；必要时改签名 |
| 稀疏键不填洞 | `count()`/`foreach` 行为变化 | 这是向 PHP 对齐；修测试而非回退 |

PHP 对照要点：数组是值语义 + 延迟 COW，对象是引用语义。本单让 Origami 的数组从"eager 值拷贝"变成"延迟 COW"，方向上是向 PHP 收敛。`&$arr[i]` 建立的引用槽在 PHP 里也会阻止该元素被分离，现有的 `RefSlotCount` 机制已经对应上了。

## 回归测试

### `tests/php/array_key_lookup_test.php`

覆盖索引正确性，尤其是密集/稀疏混合的边界：

```php
<?php
// 密集数组
$a = [10, 20, 30];
assertEq($a[0], 10); assertEq($a[2], 30);

// 稀疏键与后续追加：这是 FindSlotByIntKey 顺序敏感的核心用例
$b = [1, 2, 3];
$b[6] = 'x';
$b[] = 'p';   // 键 7
$b[] = 'q';   // 8
$b[] = 'r';   // 9
assertEq($b[6], 'x');   // 不能因为快路径返回 'r'
assertEq($b[7], 'p');
assertEq($b[9], 'r');
assertEq(count($b), 7);

// 数字字符串键与整数键等价
$c = ["5" => 'five'];
assertEq($c[5], 'five');
$c[6] = 'six';
assertEq($c["6"], 'six');

// 负键
$d = [-3 => 'neg'];
assertEq($d[-3], 'neg');

// unset 后查找
$e = [0=>'a', 1=>'b', 2=>'c'];
unset($e[1]);
assertEq($e[2], 'c');
assertEq(isset($e[1]), false);
assertEq(count($e), 2);

// 大数组跨过索引阈值
$big = [];
for ($i = 0; $i < 200; $i++) { $big[$i] = $i * 2; }
assertEq($big[150], 300);
$big["k"] = 'v';        // 引入显式键，索引须重建
assertEq($big[150], 300);
assertEq($big["k"], 'v');
unset($big[100]);        // 结构变化，索引须失效
assertEq($big[150], 300);

// 稀疏字面量（Step 3）
$sparse = [10000 => 1];
assertEq(count($sparse), 1);
assertEq($sparse[10000], 1);

// foreach 顺序
$order = [];
foreach ($b as $k => $v) { $order[] = $k; }
assertEq($order, [0,1,2,6,7,8,9]);

// 混合键 foreach
$mix = ['a'=>1, 0=>2, 'b'=>3];
$keys = [];
foreach ($mix as $k => $v) { $keys[] = $k; }
assertEq($keys, ['a', 0, 'b']);
```

### `tests/php/array_cow_test.php`

覆盖 COW 分离时机：

```php
<?php
// 基本值语义
$a = [1,2,3];
$b = $a;
$b[0] = 99;
assertEq($a[0], 1);      // 未被污染
assertEq($b[0], 99);

// 反向
$c = [1,2,3];
$d = $c;
$c[1] = 88;
assertEq($d[1], 2);

// 嵌套两层
$n = ['x' => ['y' => 1]];
$m = $n;
$m['x']['y'] = 2;
assertEq($n['x']['y'], 1);   // 内层也不能被污染

// 嵌套三层
$p = ['a' => ['b' => ['c' => 1]]];
$q = $p;
$q['a']['b']['c'] = 2;
assertEq($p['a']['b']['c'], 1);

// 传参按值
function mutate($arr) { $arr[0] = 'changed'; return $arr; }
$orig = ['keep'];
$ret = mutate($orig);
assertEq($orig[0], 'keep');
assertEq($ret[0], 'changed');

// 传参按引用
function mutateRef(&$arr) { $arr[0] = 'changed'; }
$orig2 = ['keep'];
mutateRef($orig2);
assertEq($orig2[0], 'changed');

// 引用槽在分离时保持共享
$base = [1, 2, 3];
$ref = &$base[1];
$copy = $base;
$ref = 99;
assertEq($base[1], 99);
assertEq($copy[1], 2);   // 副本不受引用写入影响

// 内置函数按引用修改
$push = [1,2];
$pushCopy = $push;
array_push($push, 3);
assertEq(count($pushCopy), 2);
assertEq(count($push), 3);

$sort = [3,1,2];
$sortCopy = $sort;
sort($sort);
assertEq($sortCopy, [3,1,2]);

// unset 不影响副本
$u = ['a'=>1,'b'=>2];
$uc = $u;
unset($u['a']);
assertEq(isset($uc['a']), true);

// 对象属性里的数组
class Holder { public array $items = [1,2,3]; }
$h = new Holder();
$local = $h->items;
$local[0] = 'x';
assertEq($h->items[0], 1);

// 关联数组走 ObjectValue 路径
$assoc = null;
$assoc['k'] = 'v';
$assocCopy = $assoc;
$assocCopy['k'] = 'w';
assertEq($assoc['k'], 'v');

// foreach 值拷贝
$fe = [[1],[2]];
foreach ($fe as $row) { $row[0] = 'mut'; }
assertEq($fe[0][0], 1);

// foreach 引用
foreach ($fe as &$row) { $row[0] = 'mut'; }
unset($row);
assertEq($fe[0][0], 'mut');
```

`assertEq` 按 `tests/php/` 现有惯例实现（失败走 `Log::fatal`）。

### 运行

```bash
go run ./zy.go tests/php/array_key_lookup_test.php
go run ./zy.go tests/php/array_cow_test.php

# 关键：debug 构建下跑全量套件
go build -tags origamidebug -o origami-debug.exe .
./origami-debug.exe tests/run_tests.php
```

## 完成定义

- [ ] Step 0 的 `origamidebug` 自检机制就位，生产构建零开销
- [ ] `origamidebug` 构建下 `tests/run_tests.php` 全程无 panic
- [ ] 惰性索引实现，指纹校验 + 显式失效点齐备
- [ ] `idxAllDense` 快路径生效，密集数组读取不再 `strconv`
- [ ] 索引阈值由 benchmark 交叉点确定，写进代码注释
- [ ] `SetVariableValue` 改为 `CowAssign`，不再 eager clone
- [ ] `rc` 初值约定统一并写进字段注释
- [ ] `node/index.go` 的写入路径逐层分离，含三层嵌套
- [ ] `std/php/array/` 下按引用修改的函数全部走 `CowSeparateIndex`
- [ ] `cowSeparateNestedArray` 改为看 `rc` 的 `SeparateArrayValue`
- [ ] 字面量稀疏键不再填洞（单独提交）
- [ ] 两份 `tests/php/` 回归通过
- [ ] `tests/run_tests.php` 全绿，`examples/laravel13` 能起
- [ ] `benchstat` 显示 `BenchmarkArrayDenseIndexRead` / `BenchmarkArrayAssignLarge` 显著改善并归档
