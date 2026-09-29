package collections

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
	"github.com/php-any/origami/std/php"
)

const enumerableName = "Illuminate\\Support\\Enumerable"

type EnumerableInterface struct{ node.Node }

func NewEnumerableInterface() data.InterfaceStmt {
	return &EnumerableInterface{}
}

func (i *EnumerableInterface) GetName() string                      { return enumerableName }
func (i *EnumerableInterface) GetExtends() []string                 { return nil }
func (i *EnumerableInterface) GetMethod(string) (data.Method, bool) { return nil, false }
func (i *EnumerableInterface) GetMethods() []data.Method            { return nil }
func (i *EnumerableInterface) GetFrom() data.From                   { return nil }
func (i *EnumerableInterface) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

const collectionName = "Illuminate\\Support\\Collection"

const (
	arrayableName        = "Illuminate\\Contracts\\Support\\Arrayable"
	jsonableName         = "Illuminate\\Contracts\\Support\\Jsonable"
	jsonSerializableName = "JsonSerializable"
)

type CollectionClass struct {
	node.Node
	methods map[string]data.Method
}

func NewCollectionClass() data.ClassStmt {
	c := &CollectionClass{methods: map[string]data.Method{}}
	c.register()
	return c
}

func (c *CollectionClass) GetName() string { return collectionName }
func (c *CollectionClass) GetExtend() *string {
	return nil
}
func (c *CollectionClass) GetImplements() []string {
	return []string{
		enumerableName,
		"ArrayAccess",
		"Countable",
		"IteratorAggregate",
		"JsonSerializable",
		"Illuminate\\Contracts\\Support\\Arrayable",
		"Illuminate\\Contracts\\Support\\Jsonable",
	}
}
func (c *CollectionClass) GetProperty(name string) (data.Property, bool) {
	if name == "items" {
		return node.NewProperty(nil, "items", "protected", false, data.NewArrayValue(nil)), true
	}
	return nil, false
}
func (c *CollectionClass) GetPropertyList() []data.Property {
	return []data.Property{
		node.NewProperty(nil, "items", "protected", false, data.NewArrayValue(nil)),
	}
}
func (c *CollectionClass) GetConstruct() data.Method {
	return c.methods["__construct"]
}
func (c *CollectionClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
func (c *CollectionClass) GetMethod(name string) (data.Method, bool) {
	m, ok := c.methods[data.MethodLookupKey(name)]
	return m, ok
}
func (c *CollectionClass) GetMethods() []data.Method {
	out := make([]data.Method, 0, len(c.methods))
	for _, m := range c.methods {
		out = append(out, m)
	}
	return out
}
func (c *CollectionClass) GetStaticMethod(name string) (data.Method, bool) {
	switch strings.ToLower(name) {
	case "make", "times", "range", "wrap", "unwrap", "empty":
		return c.GetMethod(name)
	}
	return c.GetMethod(name)
}

func (c *CollectionClass) register() {
	inst := func(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[data.MethodLookupKey(name)] = kit.InstanceMethod(name, params, fn)
	}
	instOpt := func(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[data.MethodLookupKey(name)] = kit.InstanceMethodOpt(name, params, optionalFrom, fn)
	}
	stat := func(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[data.MethodLookupKey(name)] = kit.StaticMethod(name, params, -1, fn)
	}
	statOpt := func(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[data.MethodLookupKey(name)] = kit.StaticMethod(name, params, optionalFrom, fn)
	}
	va := func(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[data.MethodLookupKey(name)] = kit.InstanceMethodVariadic(name, params, fn)
	}
	inst("__construct", []string{"items"}, collectionConstruct)
	stat("make", []string{"items"}, collectionMake)
	stat("empty", nil, collectionEmpty)
	stat("wrap", []string{"value"}, collectionWrap)
	statOpt("range", []string{"from", "to", "step"}, 2, collectionRange)
	statOpt("times", []string{"number", "callback"}, 1, collectionTimes)
	stat("unwrap", []string{"value"}, collectionUnwrap)
	inst("all", nil, collectionAll)
	inst("toArray", nil, collectionToArray)
	inst("toJson", []string{"options"}, collectionToJson)
	inst("jsonSerialize", nil, collectionJsonSerialize)
	inst("map", []string{"callback"}, collectionMap)
	inst("flatMap", []string{"callback"}, collectionFlatMap)
	inst("mapWithKeys", []string{"callback"}, collectionMapWithKeys)
	inst("filter", []string{"callback"}, collectionFilter)
	inst("partition", []string{"key", "operator", "value"}, collectionPartition)
	inst("values", nil, collectionValues)
	inst("keys", nil, collectionKeys)
	inst("pluck", []string{"value", "key"}, collectionPluck)
	inst("get", []string{"key", "default"}, collectionGet)
	inst("put", []string{"key", "value"}, collectionPut)
	va("push", []string{"values"}, collectionPush)
	inst("pop", []string{"count"}, collectionPop)
	inst("first", []string{"callback", "default"}, collectionFirst)
	inst("last", []string{"callback", "default"}, collectionLast)
	inst("count", nil, collectionCount)
	inst("isEmpty", nil, collectionIsEmpty)
	inst("isNotEmpty", nil, collectionIsNotEmpty)
	inst("each", []string{"callback"}, collectionEach)
	inst("contains", []string{"key", "operator", "value"}, collectionContains)
	inst("where", []string{"key", "operator", "value"}, collectionWhere)
	inst("unique", []string{"key", "strict"}, collectionUnique)
	inst("merge", []string{"items"}, collectionMerge)
	inst("diff", []string{"items"}, collectionDiff)
	inst("diffKeys", []string{"items"}, collectionDiffKeys)
	inst("except", []string{"keys"}, collectionExcept)
	inst("only", []string{"keys"}, collectionOnly)
	inst("concat", []string{"source"}, collectionConcat)
	inst("flatten", []string{"depth"}, collectionFlatten)
	inst("sort", []string{"callback"}, collectionSort)
	c.methods["sortby"] = kit.InstanceMethodOpt("sortBy", []string{"callback", "options", "descending"}, 1, collectionSortBy)
	inst("groupBy", []string{"groupBy", "preserveKeys"}, collectionGroupBy)
	inst("keyBy", []string{"keyBy"}, collectionKeyBy)
	inst("implode", []string{"value", "glue"}, collectionImplode)
	inst("join", []string{"glue", "finalGlue"}, collectionJoin)
	inst("getIterator", nil, collectionGetIterator)
	inst("offsetExists", []string{"key"}, collectionOffsetExists)
	inst("offsetGet", []string{"key"}, collectionOffsetGet)
	inst("offsetSet", []string{"key", "value"}, collectionOffsetSet)
	inst("offsetUnset", []string{"key"}, collectionOffsetUnset)
	inst("getArrayableItems", []string{"items"}, collectionGetArrayableItems)
	inst("toBase", nil, collectionToBase)
	// EnumeratesValues::collect() 与 toBase() 函数体完全一致（都是 new Collection($this->all())）。
	inst("collect", nil, collectionToBase)
	instOpt("getCachingIterator", []string{"flags"}, 0, collectionGetCachingIterator)
	stat("proxy", []string{"method"}, collectionProxy)
	inst("__get", []string{"key"}, collectionMagicGet)
	registerCollectionMore(c)
	registerCollectionExt(c)
	inst("__call", []string{"method", "parameters"}, collectionMissing)
	kit.RegisterMacroable(c.methods, collectionName)
	kit.RegisterConditionable(c.methods, collectionWhenProxy)
}

func collectionReceiver(ctx data.Context) (*data.ClassValue, data.Control) {
	if cv := kit.Receiver(ctx); cv != nil {
		return cv, nil
	}
	return nil, data.NewErrorThrow(nil, fmt.Errorf("Collection method missing $this"))
}

func collectionItems(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.GetProperty("items")
	if av, ok := v.(*data.ArrayValue); ok && av != nil {
		return av
	}
	empty := data.NewArrayValue(nil).(*data.ArrayValue)
	_ = cv.SetProperty("items", empty)
	return empty
}

func newCollectionInstance(ctx data.Context, items data.Value) (*data.ClassValue, data.Control) {
	vm := ctx.GetVM()
	stmt, ctl := vm.GetOrLoadClass(collectionName)
	if ctl != nil {
		return nil, ctl
	}
	// ?? $this ??? Collection ???? new static?Eloquent Collection??
	// collect() ???????????ClassMethodContext ?????? Spatie Package??
	// ???? Package ?? Collection ? new?
	if classCtx, ok := ctx.(*data.ClassMethodContext); ok && classCtx.ClassValue != nil {
		if (data.Class{Name: collectionName}).Is(classCtx.ClassValue) {
			stmt = classCtx.ClassValue.Class
		}
	}
	cv := data.NewClassValue(stmt, ctx.CreateBaseContext())
	arr, ctl := getArrayableItems(ctx, items)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("items", arr)
	return cv, nil
}

func getArrayableItems(ctx data.Context, items data.Value) (*data.ArrayValue, data.Control) {
	got, ctl := arrFromValue(ctx, items, 0)
	if ctl != nil {
		return nil, ctl
	}
	if got == nil || isNull(got) {
		return data.NewArrayValue(nil).(*data.ArrayValue), nil
	}
	if av, ok := got.(*data.ArrayValue); ok {
		return data.CloneArrayValue(av), nil
	}
	if ov, ok := got.(*data.ObjectValue); ok && ov != nil {
		out := data.NewArrayValue(nil).(*data.ArrayValue)
		ov.RangeProperties(func(key string, val data.Value) bool {
			setEntry(out, key, val)
			return true
		})
		return out, nil
	}
	return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Items cannot be represented by a scalar value."), "InvalidArgumentException")
}

func collectionConstruct(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	items, _ := ctx.GetIndexValue(0)
	arr, ctl := getArrayableItems(ctx, items)
	if ctl != nil {
		return nil, ctl
	}
	_ = cv.SetProperty("items", arr)
	return data.NewNullValue(), nil
}

func collectionMake(ctx data.Context) (data.GetValue, data.Control) {
	items, _ := ctx.GetIndexValue(0)
	return newCollectionInstance(ctx, items)
}

func collectionEmpty(ctx data.Context) (data.GetValue, data.Control) {
	return newCollectionInstance(ctx, data.NewArrayValue(nil))
}

// collectionRange 对齐 Collection::range($from, $to, $step = 1)。
// 与 PHP 的 range() 一致：$from > $to 时自动降序，$step 先取绝对值。
func collectionRange(ctx data.Context) (data.GetValue, data.Control) {
	from, ctl := intFromValue(kit.Arg(ctx, 0))
	if ctl != nil {
		return nil, ctl
	}
	to, ctl := intFromValue(kit.Arg(ctx, 1))
	if ctl != nil {
		return nil, ctl
	}
	step := 1
	if v := kit.Arg(ctx, 2); !kit.IsNull(v) {
		step, ctl = intFromValue(v)
		if ctl != nil {
			return nil, ctl
		}
	}
	if step <= 0 {
		step = 1
	}
	var items []data.Value
	if from <= to {
		items = make([]data.Value, 0, (to-from)/step+1)
		for i := from; i <= to; i += step {
			items = append(items, data.NewIntValue(i))
		}
	} else {
		items = make([]data.Value, 0, (from-to)/step+1)
		for i := from; i >= to; i -= step {
			items = append(items, data.NewIntValue(i))
		}
	}
	return newCollectionInstance(ctx, data.NewArrayValue(items))
}

// collectionTimes 对齐 EnumeratesValues::times($number, $callback = null)：
// $number < 1 返回空集合；无回调时返回 range(1, $number)，否则 map 回调（键从 0 重排）。
func collectionTimes(ctx data.Context) (data.GetValue, data.Control) {
	number, ctl := intFromValue(kit.Arg(ctx, 0))
	if ctl != nil {
		return nil, ctl
	}
	cb := kit.Arg(ctx, 1)
	if number < 1 {
		return newCollectionInstance(ctx, data.NewArrayValue(nil))
	}
	items := make([]data.Value, 0, number)
	for i := 1; i <= number; i++ {
		items = append(items, data.NewIntValue(i))
	}
	arr := data.NewArrayValue(items)
	// unless($callback == null)：无回调时短路，不进入 map。
	if kit.IsNull(cb) {
		return newCollectionInstance(ctx, arr)
	}
	mapped, ctl := arrMap(withArgs(ctx, arr, cb))
	if ctl != nil {
		return nil, ctl
	}
	return newCollectionInstance(ctx, mapped.(data.Value))
}

// collectionUnwrap 对齐 EnumeratesValues::unwrap($value)：
// Enumerable 取 all()，其余原样返回（含 null / 标量 / 数组）。
func collectionUnwrap(ctx data.Context) (data.GetValue, data.Control) {
	v := unwrapValue(kit.Arg(ctx, 0))
	if cv, ok := v.(*data.ClassValue); ok && cv != nil && classIs(cv, enumerableName) {
		arr, ok, ctl := callClassNoArg(cv, "all")
		if ctl != nil {
			return nil, ctl
		}
		if !ok {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Call to undefined method %s::all()", cv.Class.GetName()), "Error")
		}
		return arr, nil
	}
	return v, nil
}

func collectionWrap(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if cv, ok := v.(*data.ClassValue); ok && cv != nil {
		name := cv.Class.GetName()
		if name == collectionName || strings.HasSuffix(name, "\\Collection") {
			return cv, nil
		}
	}
	wrapped, _ := arrWrap(withArgs(ctx, v))
	return newCollectionInstance(ctx, wrapped.(data.Value))
}

func withArgs(ctx data.Context, args ...data.Value) data.Context {
	vars := make([]data.Variable, len(args))
	for i := range args {
		vars[i] = node.NewVariable(nil, "arg"+strconv.Itoa(i), i, nil)
	}
	c := ctx.CreateContext(vars)
	for i, a := range args {
		if a == nil {
			a = data.NewNullValue()
		}
		_ = c.SetVariableValue(vars[i], a)
	}
	return c
}

func collectionAll(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return collectionItems(cv), nil
}

// collectionToArray 对齐 EnumeratesValues::toArray()：
// $this->map(fn ($v) => $v instanceof Arrayable ? $v->toArray() : $v)->all()
// 键保持原样，值里的 Arrayable 元素必须被展开——不能直接返回 items，
// 否则 Filament\Notifications\Collection::toLivewire() 之类的调用会拿到裸对象，
// Livewire 反水合时找不到对应 synth（Property type not supported in Livewire）。
func collectionToArray(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return itemsToArray(ctx, collectionItems(cv))
}

// collectionToJson 对齐 EnumeratesValues::toJson()：json_encode($this->jsonSerialize(), $options)
func collectionToJson(ctx data.Context) (data.GetValue, data.Control) {
	return collectionEncodeJson(ctx, 0)
}

// JSON_PRETTY_PRINT（见 std/php/json_constants.go）。
const collectionPrettyPrintFlag = 128

// collectionToPrettyJson 对齐 EnumeratesValues::toPrettyJson($options)。
func collectionToPrettyJson(ctx data.Context) (data.GetValue, data.Control) {
	return collectionEncodeJson(ctx, collectionPrettyPrintFlag)
}

func collectionEncodeJson(ctx data.Context, baseFlags int) (data.GetValue, data.Control) {
	arr, ctl := collectionJsonSerialize(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if v := kit.Arg(ctx, 0); v != nil && !isNull(v) {
		if iv, ok := v.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				baseFlags |= n
			}
		}
	}
	encoded, ok, jctl := php.JsonEncodeFlags(ctx, arr.(data.Value), baseFlags)
	if jctl != nil {
		return nil, jctl
	}
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(encoded), nil
}

// collectionJsonSerialize 对齐 EnumeratesValues::jsonSerialize() 的 match 顺序：
// JsonSerializable > Jsonable > Arrayable > 原样。
func collectionJsonSerialize(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	items := collectionItems(cv)
	if !hasClassItem(items, "") {
		return items, nil
	}
	out := make([]*data.ZVal, len(items.List))
	for i, z := range items.List {
		if z == nil {
			continue
		}
		v := z.Value
		if c, ok := unwrapValue(v).(*data.ClassValue); ok && c != nil {
			converted, ok2, ctl := callJsonSerialize(ctx, c)
			if ctl != nil {
				return nil, ctl
			}
			if ok2 {
				v = converted
			}
		}
		out[i] = data.CopyZValKeepName(z, v)
	}
	return &data.ArrayValue{List: out, IndirectOverloadClass: items.IndirectOverloadClass}, nil
}

// callJsonSerialize 按 jsonSerialize() 的 match 顺序把一个对象转成可 JSON 化的值。
// ok 为 false 表示该对象不匹配任何分支，调用方应保留原值。
func callJsonSerialize(ctx data.Context, c *data.ClassValue) (data.Value, bool, data.Control) {
	switch {
	case classIs(c, jsonSerializableName):
		v, _, ctl := callClassNoArg(c, "jsonSerialize")
		if ctl != nil {
			return nil, false, ctl
		}
		if v == nil {
			v = data.NewNullValue()
		}
		return v, true, nil
	case classIs(c, jsonableName):
		js, _, ctl := callClassNoArg(c, "toJson")
		if ctl != nil {
			return nil, false, ctl
		}
		if js == nil {
			js = data.NewNullValue()
		}
		decoded, ctl := callVMFunc(ctx, "json_decode", js, data.NewBoolValue(true))
		if ctl != nil {
			return nil, false, ctl
		}
		return decoded, true, nil
	case classIs(c, arrayableName):
		v, _, ctl := callClassNoArg(c, "toArray")
		if ctl != nil {
			return nil, false, ctl
		}
		if v == nil {
			v = data.NewNullValue()
		}
		return v, true, nil
	}
	return nil, false, nil
}

// itemsToArray 是 toArray() 的核心：保持键，把 Arrayable 元素换成 $v->toArray() 的结果。
// 没有任何对象元素时直接返回原数组，标量集合不付额外分配。
func itemsToArray(ctx data.Context, items *data.ArrayValue) (data.Value, data.Control) {
	if items == nil || !hasClassItem(items, arrayableName) {
		return items, nil
	}
	src := items.List
	out := make([]*data.ZVal, len(src))
	for i, z := range src {
		if z == nil {
			continue
		}
		v := z.Value
		if c, ok := unwrapValue(v).(*data.ClassValue); ok && c != nil && classIs(c, arrayableName) {
			conv, _, ctl := callClassNoArg(c, "toArray")
			if ctl != nil {
				return nil, ctl
			}
			if conv == nil {
				conv = data.NewNullValue()
			}
			v = conv
		}
		out[i] = data.CopyZValKeepName(z, v)
	}
	return &data.ArrayValue{List: out, IndirectOverloadClass: items.IndirectOverloadClass}, nil
}

// hasClassItem 快速判断 items 里是否存在对象元素（filter 非空时还要求它实现该接口）。
// 标量集合在这里只付一次类型断言的代价，不会分配新数组。
func hasClassItem(items *data.ArrayValue, filter string) bool {
	if items == nil {
		return false
	}
	for _, z := range items.List {
		if z == nil {
			continue
		}
		c, ok := unwrapValue(z.Value).(*data.ClassValue)
		if !ok || c == nil {
			continue
		}
		if filter == "" || classIs(c, filter) {
			return true
		}
	}
	return false
}

func collectionMap(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	mapped, ctl := arrMap(withArgs(ctx, collectionItems(cv), cb))
	if ctl != nil {
		return nil, ctl
	}
	return newCollectionInstance(ctx, mapped.(data.Value))
}

func collectionFlatMap(ctx data.Context) (data.GetValue, data.Control) {
	mapped, ctl := collectionMap(ctx)
	if ctl != nil {
		return nil, ctl
	}
	mcv, ok := mapped.(*data.ClassValue)
	if !ok {
		return mapped, nil
	}
	items := collectionItems(mcv)
	collapsed, ctl := arrCollapse(withArgs(ctx, items))
	if ctl != nil {
		return nil, ctl
	}
	return newCollectionInstance(ctx, collapsed.(data.Value))
}

func collectionMapWithKeys(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	mapped, ctl := arrMapWithKeys(withArgs(ctx, collectionItems(cv), cb))
	if ctl != nil {
		return nil, ctl
	}
	return newCollectionInstance(ctx, mapped.(data.Value))
}

func collectionFilter(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	var filtered data.GetValue
	if cb == nil || isNull(cb) {
		out := data.NewArrayValue(nil).(*data.ArrayValue)
		for _, e := range toEntries(collectionItems(cv)) {
			if truthy(e.value) {
				setEntry(out, e.keyStr, e.value)
			}
		}
		filtered = out
	} else {
		var err data.Control
		filtered, err = arrWhere(withArgs(ctx, collectionItems(cv), cb))
		if err != nil {
			return nil, err
		}
	}
	return newCollectionInstance(ctx, filtered.(data.Value))
}

func collectionPartition(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	passed := data.NewArrayValue(nil).(*data.ArrayValue)
	failed := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		ok := false
		if cb != nil && !isNull(cb) {
			ok, ctl = callBool(ctx, cb, e.value, e.key)
			if ctl != nil {
				return nil, ctl
			}
		}
		if ok {
			setEntry(passed, e.keyStr, e.value)
		} else {
			setEntry(failed, e.keyStr, e.value)
		}
	}
	left, ctl := newCollectionInstance(ctx, passed)
	if ctl != nil {
		return nil, ctl
	}
	right, ctl := newCollectionInstance(ctx, failed)
	if ctl != nil {
		return nil, ctl
	}
	both := data.NewArrayValue([]data.Value{left, right}).(*data.ArrayValue)
	return newCollectionInstance(ctx, both)
}

func collectionValues(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		out.List = append(out.List, data.NewZVal(e.value))
	}
	return newCollectionInstance(ctx, out)
}

func collectionKeys(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		out.List = append(out.List, data.NewZVal(e.key))
	}
	return newCollectionInstance(ctx, out)
}

func collectionPluck(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	valueKey, _ := ctx.GetIndexValue(0)
	keyKey, _ := ctx.GetIndexValue(1)
	plucked, err := arrPluck(withArgs(ctx, collectionItems(cv), valueKey, keyKey))
	if err != nil {
		return nil, err
	}
	return newCollectionInstance(ctx, plucked.(data.Value))
}

func collectionGet(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, _ := ctx.GetIndexValue(0)
	def, _ := ctx.GetIndexValue(1)
	items := collectionItems(cv)
	ks := keyToString(key)
	// Collection::get 是纯键访问（对齐 offsetExists），不能直接走 dataGetPath：
	// 空路径在 data_get 里是「原样返回整体」，会让 get('', collect()) 把 items
	// 本身（ArrayValue）当命中值返回，破坏声明的 Collection 返回类型。
	if z, ok := items.LookupZValByStringKey(ks); ok && z != nil {
		return z.Value, nil
	}
	if ks != "" {
		if v, ok := dataGetPath(items, ks); ok {
			return v, nil
		}
	}
	return laravelValue(ctx, def)
}

func collectionPut(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, _ := ctx.GetIndexValue(0)
	val, _ := ctx.GetIndexValue(1)
	items := collectionItems(cv)
	setEntry(items, keyToString(key), val)
	_ = cv.SetProperty("items", items)
	return cv, nil
}

func collectionPush(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	items := collectionItems(cv)
	// variadic: first arg may be array of values
	if v, ok := ctx.GetIndexValue(0); ok && v != nil {
		if av, ok := v.(*data.ArrayValue); ok {
			for _, e := range toEntries(av) {
				items.List = append(items.List, data.NewZVal(e.value))
			}
		} else {
			items.List = append(items.List, data.NewZVal(v))
		}
	}
	_ = cv.SetProperty("items", items)
	return cv, nil
}

func collectionPop(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	items := collectionItems(cv)
	if len(items.List) == 0 {
		return data.NewNullValue(), nil
	}
	last := items.List[len(items.List)-1]
	items.List = items.List[:len(items.List)-1]
	_ = cv.SetProperty("items", items)
	if last == nil {
		return data.NewNullValue(), nil
	}
	return last.Value, nil
}

func collectionFirst(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	def, _ := ctx.GetIndexValue(1)
	return arrFirst(withArgs(ctx, collectionItems(cv), cb, def))
}

func collectionLast(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	def, _ := ctx.GetIndexValue(1)
	return arrLast(withArgs(ctx, collectionItems(cv), cb, def))
}

func collectionCount(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewIntValue(len(toEntries(collectionItems(cv)))), nil
}

func collectionIsEmpty(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return data.NewBoolValue(len(toEntries(collectionItems(cv))) == 0), nil
}

func collectionIsNotEmpty(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := collectionIsEmpty(ctx)
	if ctl != nil {
		return nil, ctl
	}
	b, _ := v.(data.AsBool).AsBool()
	return data.NewBoolValue(!b), nil
}

func collectionEach(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	for _, e := range toEntries(collectionItems(cv)) {
		ret, err := callValue(ctx, cb, e.value, e.key)
		if err != nil {
			return nil, err
		}
		if isFalse(ret) {
			break
		}
	}
	return cv, nil
}

func collectionContains(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return collectionContainsItems(ctx, cv, false)
}

func collectionWhere(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, _ := ctx.GetIndexValue(0)
	op, _ := ctx.GetIndexValue(1)
	val, _ := ctx.GetIndexValue(2)
	operator, target := whereOperands(len(ctx.GetCallArgs()), op, val)
	out := whereCore(collectionItems(cv), keyToString(key), operator, target)
	return newCollectionInstance(ctx, out)
}

// whereOperands 对齐 EnumeratesValues::operatorForWhere 的形式判定：
// where($k) → $k == true；where($k, $v) → $k == $v；where($k, $op, $v) → $k $op $v。
// 未传的形参在帧里是 null，所以多出来的第三参必须靠实参个数（func_num_args 语义）区分。
func whereOperands(argc int, op, val data.Value) (string, data.Value) {
	switch {
	case argc == 3:
		return keyToString(op), val
	case isNull(op) && isNull(val):
		return "=", data.NewBoolValue(true)
	case isNull(val):
		return "=", op
	default:
		return keyToString(op), val
	}
}

// compareOp 对齐 PHP 的比较运算符（含 === / !== 的严格语义）。
func compareOp(left data.Value, op string, right data.Value) bool {
	switch op {
	case "===":
		return phpStrictEquals(left, right)
	case "!==":
		return !phpStrictEquals(left, right)
	case "=", "==":
		return phpLooseEquals(left, right)
	case "!=", "<>":
		return !phpLooseEquals(left, right)
	case "<=>":
		cmp, ok := phpCompareValues(left, right)
		return ok && cmp != 0
	}
	cmp, ok := phpCompareValues(left, right)
	if !ok {
		return false
	}
	switch op {
	case ">":
		return cmp > 0
	case "<":
		return cmp < 0
	case ">=":
		return cmp >= 0
	case "<=":
		return cmp <= 0
	default:
		return phpLooseEquals(left, right)
	}
}

func collectionUnique(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, _ := ctx.GetIndexValue(0)
	strict := false
	if v, _ := ctx.GetIndexValue(1); v != nil && !isNull(v) {
		strict = phpToBool(v)
	}
	path := ""
	if key != nil && !isNull(key) {
		path = keyToString(key)
	}
	return newCollectionInstance(ctx, uniqueCore(collectionItems(cv), path, strict))
}

func collectionDiff(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	other, _ := ctx.GetIndexValue(0)
	seen := map[string]struct{}{}
	otherItems, ctl := getArrayableItems(ctx, other)
	if ctl != nil {
		return nil, ctl
	}
	for _, e := range toEntries(otherItems) {
		seen[e.value.AsString()] = struct{}{}
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		if _, ok := seen[e.value.AsString()]; ok {
			continue
		}
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionDiffKeys(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	other, _ := ctx.GetIndexValue(0)
	exclude := map[string]struct{}{}
	otherItems, ctl := getArrayableItems(ctx, other)
	if ctl != nil {
		return nil, ctl
	}
	for _, e := range toEntries(otherItems) {
		exclude[e.keyStr] = struct{}{}
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		if _, ok := exclude[e.keyStr]; ok {
			continue
		}
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionExcept(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	keys, _ := ctx.GetIndexValue(0)
	exclude := map[string]struct{}{}
	for _, k := range keysToStrings(keys) {
		exclude[k] = struct{}{}
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		if _, ok := exclude[e.keyStr]; ok {
			continue
		}
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionOnly(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	keys, _ := ctx.GetIndexValue(0)
	keep := map[string]struct{}{}
	for _, k := range keysToStrings(keys) {
		keep[k] = struct{}{}
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		if _, ok := keep[e.keyStr]; !ok {
			continue
		}
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionMerge(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	other, _ := ctx.GetIndexValue(0)
	out := data.CloneArrayValue(collectionItems(cv))
	otherItems, ctl := getArrayableItems(ctx, other)
	if ctl != nil {
		return nil, ctl
	}
	for _, e := range toEntries(otherItems) {
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionConcat(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	other, _ := ctx.GetIndexValue(0)
	out := data.CloneArrayValue(collectionItems(cv))
	otherItems, ctl := getArrayableItems(ctx, other)
	if ctl != nil {
		return nil, ctl
	}
	for _, e := range toEntries(otherItems) {
		out.List = append(out.List, data.NewZVal(e.value))
	}
	return newCollectionInstance(ctx, out)
}

func collectionFlatten(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	depth := data.NewIntValue(1)
	if d, ok := ctx.GetIndexValue(0); ok && d != nil {
		depth = d
	}
	flat, err := arrFlatten(withArgs(ctx, collectionItems(cv), depth))
	if err != nil {
		return nil, err
	}
	return newCollectionInstance(ctx, flat.(data.Value))
}

func collectionSort(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	sorted, err := arrSort(withArgs(ctx, collectionItems(cv)))
	if err != nil {
		return nil, err
	}
	return newCollectionInstance(ctx, sorted.(data.Value))
}

func collectionSortBy(ctx data.Context) (data.GetValue, data.Control) {
	return collectionSortByDir(ctx, false)
}

func collectionSortByDir(ctx data.Context, forceDesc bool) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	callback := kit.Arg(ctx, 0)
	descending := forceDesc
	if !forceDesc {
		if d := kit.Arg(ctx, 2); d != nil && !kit.IsNull(d) {
			if b, ok := d.(data.AsBool); ok {
				descending, _ = b.AsBool()
			} else {
				descending = kit.Truthy(d)
			}
		}
	}
	type scored struct {
		key   data.Value
		entry kv
		idx   int
	}
	entries := toEntries(collectionItems(cv))
	scoredList := make([]scored, 0, len(entries))
	for i, e := range entries {
		var sortKey data.Value = e.value
		if callback != nil && !kit.IsNull(callback) && isCallableValue(callback) {
			ret, ctl := kit.Call(ctx, callback, e.value, data.NewStringValue(e.keyStr))
			if ctl != nil {
				return nil, ctl
			}
			if ret != nil {
				if v, ok := ret.(data.Value); ok {
					sortKey = v
				}
			}
		} else if callback != nil && !kit.IsNull(callback) {
			if got, ok := dataGetPath(e.value, callback.AsString()); ok {
				sortKey = got
			}
		}
		scoredList = append(scoredList, scored{key: sortKey, entry: e, idx: i})
	}
	sort.SliceStable(scoredList, func(i, j int) bool {
		cmp := compareSortKeys(scoredList[i].key, scoredList[j].key)
		if descending {
			return cmp > 0
		}
		return cmp < 0
	})
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, s := range scoredList {
		setEntry(out, s.entry.keyStr, s.entry.value)
	}
	return newCollectionInstance(ctx, out)
}

// compareSortKeys 对齐 PHP 的默认排序比较：数字/数字字符串按数值，其余按字符串；
// null 排在前面（与 PHP 的 NULL 处理一致）。
func compareSortKeys(a, b data.Value) int {
	a = kit.Unwrap(a)
	b = kit.Unwrap(b)
	aNil, bNil := a == nil || isNull(a), b == nil || isNull(b)
	if aNil || bNil {
		return cmpFloat(boolToFloat(!aNil), boolToFloat(!bNil))
	}
	if cmp, ok := phpCompareValues(a, b); ok {
		return cmp
	}
	return strings.Compare(a.AsString(), b.AsString())
}

func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func collectionGroupBy(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	groupBy := kit.Arg(ctx, 0)
	preserveKeys := false
	if preserve := kit.Arg(ctx, 1); preserve != nil && !kit.IsNull(preserve) {
		preserveKeys = kit.Truthy(preserve)
	}
	groups := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		var groupKeys []data.Value
		if groupBy != nil && !kit.IsNull(groupBy) && isCallableValue(groupBy) {
			ret, callCtl := kit.Call(ctx, groupBy, e.value, e.key)
			if callCtl != nil {
				return nil, callCtl
			}
			if value, ok := ret.(data.Value); ok && value != nil {
				value = kit.Unwrap(value)
				if keys, ok := value.(*data.ArrayValue); ok {
					for _, key := range toEntries(keys) {
						groupKeys = append(groupKeys, key.value)
					}
				} else {
					groupKeys = append(groupKeys, value)
				}
			}
		} else if groupBy != nil && !kit.IsNull(groupBy) {
			if value, ok := dataGetPath(e.value, keyToString(groupBy)); ok {
				groupKeys = append(groupKeys, value)
			}
		}
		if len(groupKeys) == 0 {
			groupKeys = append(groupKeys, data.NewNullValue())
		}
		for _, groupKey := range groupKeys {
			gk := collectionGroupKey(groupKey)
			var garr *data.ArrayValue
			// 分组结果的键是 PHP 数组字面键，不是 data_get() 路径。
			// 尤其空字符串键表示“无父导航项”；把它当空路径会错误地
			// 返回 groups 自身，随后把元素追加到分组数组的外层。
			if slot, ok := groups.LookupZValByStringKey(gk); ok && slot != nil {
				garr, _ = slot.Value.(*data.ArrayValue)
			}
			if garr == nil {
				garr = data.NewArrayValue(nil).(*data.ArrayValue)
				setEntry(groups, gk, garr)
			}
			if preserveKeys {
				setEntry(garr, e.keyStr, e.value)
			} else {
				garr.List = append(garr.List, data.NewZVal(e.value))
			}
		}
	}
	// wrap each group as Collection
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(groups) {
		inst, err := newCollectionInstance(ctx, e.value)
		if err != nil {
			return nil, err
		}
		setEntry(out, e.keyStr, inst)
	}
	return newCollectionInstance(ctx, out)
}

func collectionGroupKey(value data.Value) string {
	value = kit.Unwrap(value)
	switch key := value.(type) {
	case nil, *data.NullValue:
		return ""
	case *data.BoolValue:
		if key.Value {
			return "1"
		}
		return "0"
	default:
		return keyToString(value)
	}
}

func collectionKeyBy(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	keyBy, _ := ctx.GetIndexValue(0)
	keyed, err := arrKeyBy(withArgs(ctx, collectionItems(cv), keyBy))
	if err != nil {
		return nil, err
	}
	return newCollectionInstance(ctx, keyed.(data.Value))
}

func collectionImplode(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	value, _ := ctx.GetIndexValue(0)
	glue, _ := ctx.GetIndexValue(1)
	items := collectionItems(cv)
	if glue == nil || isNull(glue) {
		return arrJoin(withArgs(ctx, items, value, nil))
	}
	plucked, err := arrPluck(withArgs(ctx, items, value, nil))
	if err != nil {
		return nil, err
	}
	return arrJoin(withArgs(ctx, plucked.(data.Value), glue, nil))
}

func collectionJoin(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	glue, _ := ctx.GetIndexValue(0)
	finalGlue, _ := ctx.GetIndexValue(1)
	return arrJoin(withArgs(ctx, collectionItems(cv), glue, finalGlue))
}

func collectionGetIterator(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return collectionIterator(ctx, cv)
}

// collectionIterator 抽成纯函数：getIterator() / getCachingIterator() 都要用它。
func collectionIterator(ctx data.Context, cv *data.ClassValue) (*data.ClassValue, data.Control) {
	return newNativeInstance(ctx, ctx.GetVM(), "ArrayIterator", []data.Value{collectionItems(cv)})
}

// newNativeInstance 按类名到 VM 取（或自动加载）一个原生类，构造出实例。
// 模式沿用本文件既有写法：CreateContext + BindDeclaredArgs + Call，$this 才成立。
func newNativeInstance(ctx data.Context, vm data.VM, className string, args []data.Value) (*data.ClassValue, data.Control) {
	cls, ok := vm.GetClass(className)
	if !ok {
		var ctl data.Control
		cls, ctl = vm.GetOrLoadClass(className)
		if ctl != nil {
			return nil, ctl
		}
	}
	if cls == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Class %s not found", className))
	}
	cv := data.NewClassValue(cls, ctx.CreateBaseContext())
	if ctor := cls.GetConstruct(); ctor != nil {
		nctx := cv.CreateContext(ctor.GetVariables())
		data.BindDeclaredArgs(nctx, ctor, args)
		if _, ctl := ctor.Call(nctx); ctl != nil {
			return nil, ctl
		}
	}
	return cv, nil
}

// collectionGetCachingIterator 对齐 EnumeratesValues::getCachingIterator($flags = CachingIterator::CALL_TOSTRING)。
func collectionGetCachingIterator(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	// CachingIterator::CALL_TOSTRING
	const callToString = 1
	flags := data.Value(data.NewIntValue(callToString))
	if v := kit.Arg(ctx, 0); v != nil && !kit.IsNull(v) {
		flags = v
	}
	inner, ctl := collectionIterator(ctx, cv)
	if ctl != nil {
		return nil, ctl
	}
	return newNativeInstance(ctx, ctx.GetVM(), "CachingIterator", []data.Value{inner, flags})
}

// collectionProxy 对齐 EnumeratesValues::proxy($method)：static::$proxies[] = $method。
func collectionProxy(ctx data.Context) (data.GetValue, data.Control) {
	if v := kit.Arg(ctx, 0); v != nil {
		addCollectionProxy(v.AsString())
	}
	return data.NewNullValue(), nil
}

func collectionOffsetExists(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, _ := ctx.GetIndexValue(0)
	_, ok := arrayGet(collectionItems(cv), keyToString(key))
	return data.NewBoolValue(ok), nil
}

func collectionOffsetGet(ctx data.Context) (data.GetValue, data.Control) {
	return collectionGet(ctx)
}

func collectionOffsetSet(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, _ := ctx.GetIndexValue(0)
	val, _ := ctx.GetIndexValue(1)
	items := collectionItems(cv)
	if key == nil || isNull(key) {
		items.List = append(items.List, data.NewZVal(val))
	} else {
		setEntry(items, keyToString(key), val)
	}
	_ = cv.SetProperty("items", items)
	return data.NewNullValue(), nil
}

func collectionOffsetUnset(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, _ := ctx.GetIndexValue(0)
	items := collectionItems(cv)
	items.UnsetKey(key)
	_ = cv.SetProperty("items", items)
	return data.NewNullValue(), nil
}

func collectionGetArrayableItems(ctx data.Context) (data.GetValue, data.Control) {
	items, _ := ctx.GetIndexValue(0)
	arr, ctl := getArrayableItems(ctx, items)
	if ctl != nil {
		return nil, ctl
	}
	return arr, nil
}

func collectionToBase(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	// ???? Support\Collection????????
	vm := ctx.GetVM()
	stmt, err := vm.GetOrLoadClass(collectionName)
	if err != nil {
		return nil, err
	}
	base := data.NewClassValue(stmt, ctx.CreateBaseContext())
	_ = base.SetProperty("items", data.CloneArrayValue(collectionItems(cv)))
	return base, nil
}

func collectionMagicGet(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := ""
	if v := kit.Arg(ctx, 0); v != nil {
		key = v.AsString()
	}
	if _, ok := collectionProxyNames()[data.MethodLookupKey(key)]; ok {
		return newHigherOrderProxy(ctx, cv, key)
	}
	// 对齐 EnumeratesValues::__get：非代理属性直接抛，不能静默返回 null。
	return nil, data.NewErrorThrow(nil, fmt.Errorf(
		"Property [%s] does not exist on this collection instance.", key))
}

func jsonEncodeSimple(v data.Value) (string, bool) {
	switch t := v.(type) {
	case *data.NullValue:
		return "null", true
	case *data.BoolValue:
		if t.Value {
			return "true", true
		}
		return "false", true
	case *data.IntValue:
		return strconv.Itoa(t.Value), true
	case *data.StringValue:
		b, err := jsonMarshalString(t.Value)
		return string(b), err == nil
	case *data.ArrayValue:
		if isListArray(t) {
			parts := make([]string, 0, len(t.List))
			for _, z := range t.List {
				if z == nil {
					parts = append(parts, "null")
					continue
				}
				s, ok := jsonEncodeSimple(z.Value)
				if !ok {
					return "", false
				}
				parts = append(parts, s)
			}
			return "[" + strings.Join(parts, ",") + "]", true
		}
		parts := make([]string, 0)
		for _, e := range toEntries(t) {
			ks, err := jsonMarshalString(e.keyStr)
			if err != nil {
				return "", false
			}
			vs, ok := jsonEncodeSimple(e.value)
			if !ok {
				return "", false
			}
			parts = append(parts, string(ks)+":"+vs)
		}
		return "{" + strings.Join(parts, ",") + "}", true
	default:
		b, err := jsonMarshalString(v.AsString())
		return string(b), err == nil
	}
}

func jsonMarshalString(s string) ([]byte, error) {
	return []byte(`"` + strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `"`, `\"`) + `"`), nil
}
