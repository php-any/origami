package collections

import (
	"fmt"
	"math"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

// 本文件补齐 EnumeratesValues/Collection 中高频但此前缺失的方法。
// 全部是薄封装：核心逻辑抽成 whereCore/uniqueCore/whereInCore 等纯函数，
// 别名方法直接复用，不再经过 withArgs 造上下文（CreateContext 不保留 $this，
// 也没有必要为一次别名调用分配上下文）。

func registerCollectionExt(c *CollectionClass) {
	inst := func(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[data.MethodLookupKey(name)] = kit.InstanceMethod(name, params, fn)
	}
	opt := func(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[data.MethodLookupKey(name)] = kit.InstanceMethodOpt(name, params, optionalFrom, fn)
	}
	statOpt := func(name string, params []string, optionalFrom int, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[data.MethodLookupKey(name)] = kit.StaticMethod(name, params, optionalFrom, fn)
	}

	inst("value", []string{"key", "default"}, collectionValue)
	inst("hasMany", []string{"key", "operator", "value"}, collectionHasMany)
	inst("whereStrict", []string{"key", "value"}, collectionWhereStrict)
	inst("whereInStrict", []string{"key", "values"}, collectionWhereInStrict)
	inst("whereNotInStrict", []string{"key", "values"}, collectionWhereNotInStrict)
	opt("uniqueStrict", []string{"key"}, 0, collectionUniqueStrict)
	inst("whereBetween", []string{"key", "values"}, collectionWhereBetween)
	inst("whereNotBetween", []string{"key", "values"}, collectionWhereNotBetween)
	inst("forPage", []string{"page", "perPage"}, collectionForPage)
	opt("percentage", []string{"callback", "precision"}, 1, collectionPercentage)
	inst("reduceWithKeys", []string{"callback", "initial"}, collectionReduceWithKeys)
	inst("reduceInto", []string{"initial", "callback"}, collectionReduceInto)
	c.methods["reducespread"] = kit.InstanceMethodVariadic("reduceSpread", []string{"callback", "initial"}, collectionReduceSpread)
	inst("pipeThrough", []string{"callbacks"}, collectionPipeThrough)
	inst("pipeInto", []string{"class"}, collectionPipeInto)
	inst("mapSpread", []string{"callback"}, collectionMapSpread)
	inst("eachSpread", []string{"callback"}, collectionEachSpread)
	inst("mapInto", []string{"class"}, collectionMapInto)
	inst("mapToGroups", []string{"callback"}, collectionMapToGroups)
	inst("__toString", nil, collectionToString)
	inst("toPrettyJson", []string{"options"}, collectionToPrettyJson)
	inst("dump", nil, collectionDump)
	inst("dd", nil, collectionDd)
	opt("escapeWhenCastingToString", []string{"escape"}, 0, collectionEscapeWhenCastingToString)
	statOpt("fromJson", []string{"json", "depth", "flags"}, 1, collectionFromJson)
}

// ---- where 家族 ----

// whereCore 对齐 EnumeratesValues::operatorForWhere 的过滤部分。
// 取不到路径时按 PHP 的 data_get 给 null，而不是跳过该项。
func whereCore(items *data.ArrayValue, path, operator string, target data.Value) *data.ArrayValue {
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(items) {
		itemVal, ok := dataGetPath(e.value, path)
		if !ok || itemVal == nil {
			itemVal = data.NewNullValue()
		}
		if compareOp(itemVal, operator, target) {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return out
}

func collectionWhereStrict(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	out := whereCore(collectionItems(cv), keyToString(kit.Arg(ctx, 0)), "===", kit.Arg(ctx, 1))
	return newCollectionInstance(ctx, out)
}

func collectionWhereBetween(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	path := keyToString(kit.Arg(ctx, 0))
	values, ctl := getArrayableItems(ctx, kit.Arg(ctx, 1))
	if ctl != nil {
		return nil, ctl
	}
	low, high := bounds(values)
	// PHP: where($key,'>=',reset($values))->where($key,'<=',end($values))，一次遍历等价。
	out := whereCore(collectionItems(cv), path, ">=", low)
	out = whereCore(out, path, "<=", high)
	return newCollectionInstance(ctx, out)
}

func collectionWhereNotBetween(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	path := keyToString(kit.Arg(ctx, 0))
	values, ctl := getArrayableItems(ctx, kit.Arg(ctx, 1))
	if ctl != nil {
		return nil, ctl
	}
	low, high := bounds(values)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		itemVal, ok := dataGetPath(e.value, path)
		if !ok || itemVal == nil {
			itemVal = data.NewNullValue()
		}
		// data_get($item, $key) < reset($values) || data_get($item, $key) > end($values)
		if compareOp(itemVal, "<", low) || compareOp(itemVal, ">", high) {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return newCollectionInstance(ctx, out)
}

func bounds(values *data.ArrayValue) (data.Value, data.Value) {
	if values == nil || values.Len() == 0 {
		// PHP 的 reset()/end() 在空数组上返回 false。
		return data.NewBoolValue(false), data.NewBoolValue(false)
	}
	first, last := values.At(0), values.At(values.Len()-1)
	var low, high data.Value = data.NewBoolValue(false), data.NewBoolValue(false)
	if first != nil {
		low = first.ReadValue()
	}
	if last != nil {
		high = last.ReadValue()
	}
	return low, high
}

func collectionWhereInStrict(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	values, ctl := getArrayableItems(ctx, kit.Arg(ctx, 1))
	if ctl != nil {
		return nil, ctl
	}
	out := whereInCore(collectionItems(cv), keyToString(kit.Arg(ctx, 0)), values, true, false)
	return newCollectionInstance(ctx, out)
}

func collectionWhereNotInStrict(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	values, ctl := getArrayableItems(ctx, kit.Arg(ctx, 1))
	if ctl != nil {
		return nil, ctl
	}
	out := whereInCore(collectionItems(cv), keyToString(kit.Arg(ctx, 0)), values, true, true)
	return newCollectionInstance(ctx, out)
}

// collectionValue 对齐 EnumeratesValues::value：
// 第一个含该键的元素上取 data_get，全都没有则返回 $default。
func collectionValue(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	path := keyToString(kit.Arg(ctx, 0))
	for _, e := range toEntries(collectionItems(cv)) {
		if got, ok := dataGetPath(e.value, path); ok {
			return got, nil
		}
	}
	def := kit.Arg(ctx, 1)
	if def == nil {
		def = data.NewNullValue()
	}
	return def, nil
}

// collectionHasMany 对齐 EnumeratesValues::hasMany：
// 过滤后计数 > 1。1 参形式是回调，多参形式等价 where 的判定。
func collectionHasMany(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	argc := len(ctx.GetCallArgs())
	count := 0
	if argc > 1 {
		operator, target := whereOperands(argc, kit.Arg(ctx, 1), kit.Arg(ctx, 2))
		path := keyToString(kit.Arg(ctx, 0))
		for _, e := range toEntries(collectionItems(cv)) {
			itemVal, ok := dataGetPath(e.value, path)
			if !ok || itemVal == nil {
				itemVal = data.NewNullValue()
			}
			if compareOp(itemVal, operator, target) {
				count++
			}
		}
		return data.NewBoolValue(count > 1), nil
	}
	cb := kit.Arg(ctx, 0)
	if !isCallableValue(cb) {
		// PHP 这里要求可调用对象，非 callable 会 TypeError；按值比较兜底。
		for _, e := range toEntries(collectionItems(cv)) {
			if phpLooseEquals(e.value, cb) {
				count++
			}
		}
		return data.NewBoolValue(count > 1), nil
	}
	for _, e := range toEntries(collectionItems(cv)) {
		ok, ctl := callBool(ctx, cb, e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		if ok {
			count++
			if count > 1 {
				return data.NewBoolValue(true), nil
			}
		}
	}
	return data.NewBoolValue(false), nil
}

// ---- unique 家族 ----

// uniqueCore 对齐 EnumeratesValues::unique：$key 取到的值参与去重（取不到给 null）。
func uniqueCore(items *data.ArrayValue, path string, strict bool) *data.ArrayValue {
	eq := phpLooseEquals
	if strict {
		eq = phpStrictEquals
	}
	var seen []data.Value
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(items) {
		val := e.value
		if path != "" {
			got, ok := dataGetPath(val, path)
			if !ok || got == nil {
				got = data.NewNullValue()
			}
			val = got
		}
		dup := false
		for _, s := range seen {
			if eq(val, s) {
				dup = true
				break
			}
		}
		if dup {
			continue
		}
		seen = append(seen, val)
		setEntry(out, e.keyStr, e.value)
	}
	return out
}

func collectionUniqueStrict(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	path := ""
	if key := kit.Arg(ctx, 0); key != nil && !isNull(key) {
		path = keyToString(key)
	}
	return newCollectionInstance(ctx, uniqueCore(collectionItems(cv), path, true))
}

// ---- whereIn 核心 ----

// whereInCore 对齐 EnumeratesValues::whereIn/whereNotIn（in_array 的宽松/严格语义）。
func whereInCore(items *data.ArrayValue, path string, values *data.ArrayValue, strict, negate bool) *data.ArrayValue {
	eq := phpLooseEquals
	if strict {
		eq = phpStrictEquals
	}
	allowed := make([]data.Value, 0, values.Len())
	for _, e := range toEntries(values) {
		allowed = append(allowed, e.value)
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(items) {
		val := e.value
		if path != "" {
			got, ok := dataGetPath(e.value, path)
			if !ok || got == nil {
				got = data.NewNullValue()
			}
			val = got
		}
		in := false
		for _, a := range allowed {
			if eq(val, a) {
				in = true
				break
			}
		}
		if in != negate {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return out
}

// ---- 分页 / 统计 ----

// collectionForPage 对齐 EnumeratesValues::forPage：slice(max(0, ($page-1)*$perPage), $perPage)。
func collectionForPage(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	page, ctl := intFromValue(kit.Arg(ctx, 0))
	if ctl != nil {
		return nil, ctl
	}
	perPage, ctl := intFromValue(kit.Arg(ctx, 1))
	if ctl != nil {
		return nil, ctl
	}
	offset := (page - 1) * perPage
	if offset < 0 {
		offset = 0
	}
	return sliceCollection(ctx, cv, offset, perPage)
}

// collectionPercentage 对齐 EnumeratesValues::percentage。
func collectionPercentage(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	entries := toEntries(collectionItems(cv))
	if len(entries) == 0 {
		return data.NewNullValue(), nil
	}
	cb := kit.Arg(ctx, 0)
	matched := 0
	for _, e := range entries {
		ok, ctl := callBool(ctx, cb, e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		if ok {
			matched++
		}
	}
	precision := 2
	if p := kit.Arg(ctx, 1); p != nil && !isNull(p) {
		precision, ctl = intFromValue(p)
		if ctl != nil {
			return nil, ctl
		}
	}
	value := float64(matched) / float64(len(entries)) * 100
	return data.NewFloatValue(roundTo(value, precision)), nil
}

func roundTo(v float64, precision int) float64 {
	if precision < 0 {
		precision = 0
	}
	p := math.Pow(10, float64(precision))
	return math.Round(v*p) / p
}

// ---- reduce 家族 ----

// collectionReduceWithKeys 与 reduce 同义（对齐 EnumeratesValues::reduceWithKeys）。
func collectionReduceWithKeys(ctx data.Context) (data.GetValue, data.Control) {
	return collectionReduce(ctx)
}

// collectionReduceInto 对齐 EnumeratesValues::reduceInto：$initial 按引用累积。
func collectionReduceInto(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	carry := kit.Arg(ctx, 0)
	cb := kit.Arg(ctx, 1)
	if !isCallableValue(cb) {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Illuminate\\Support\\Collection::reduceInto(): Argument #2 ($callback) must be a valid callback"))
	}
	for _, e := range toEntries(collectionItems(cv)) {
		if _, ctl := callValue(ctx, cb, carry, e.value, e.key); ctl != nil {
			return nil, ctl
		}
	}
	return carry, nil
}

// collectionReduceSpread 对齐 EnumeratesValues::reduceSpread：
// 回调返回值必须是数组，且会展开成下一轮的参数。
func collectionReduceSpread(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb := kit.Arg(ctx, 0)
	result, ctl := getArrayableItems(ctx, kit.Arg(ctx, 1))
	if ctl != nil {
		return nil, ctl
	}
	for _, e := range toEntries(collectionItems(cv)) {
		args := make([]data.Value, 0, result.Len()+2)
		for arraySlots42, arrayPosition42 := result.View(), 0; arrayPosition42 < arraySlots42.Len(); arrayPosition42++ {
			item := arraySlots42.At(arrayPosition42)
			if item != nil {
				args = append(args, item.ReadValue())
			}
		}
		args = append(args, e.value, e.key)
		ret, ctl := callValue(ctx, cb, args...)
		if ctl != nil {
			return nil, ctl
		}
		next, ok := ret.(*data.ArrayValue)
		if !ok {
			return nil, data.NewErrorThrow(nil, fmt.Errorf(
				"Illuminate\\Support\\Collection::reduceSpread expects reducer to return an array, but got a different value."))
		}
		result = next
	}
	return result, nil
}

// ---- 映射 ----

// collectionMapSpread 对齐 EnumeratesValues::mapSpread：把元素数组展开成回调参数，末位补键。
func collectionMapSpread(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb := kit.Arg(ctx, 0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		args, ctl := spreadArgs(e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		ret, ctl := callValue(ctx, cb, args...)
		if ctl != nil {
			return nil, ctl
		}
		if ret == nil {
			continue
		}
		setEntry(out, e.keyStr, asValue(ret))
	}
	return newCollectionInstance(ctx, out)
}

// collectionEachSpread 对齐 EnumeratesValues::eachSpread。
func collectionEachSpread(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb := kit.Arg(ctx, 0)
	for _, e := range toEntries(collectionItems(cv)) {
		args, ctl := spreadArgs(e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		ret, ctl := callValue(ctx, cb, args...)
		if ctl != nil {
			return nil, ctl
		}
		if isFalse(ret) {
			break
		}
	}
	return cv, nil
}

// isFalse 对齐 PHP 的 `=== false`：each/eachSpread 只在回调**返回布尔 false** 时提前结束。
// 用宽松真值判断会把 `$arr[] = $v;`（赋值表达式返回非布尔值）这类语句当成 false 而误停。
func isFalse(v data.GetValue) bool {
	bv, ok := asValue(v).(*data.BoolValue)
	return ok && !bv.Value
}

func spreadArgs(value, key data.Value) ([]data.Value, data.Control) {
	var args []data.Value
	if arr, ok := unwrapValue(value).(*data.ArrayValue); ok && arr != nil {
		args = make([]data.Value, 0, arr.Len()+1)
		for arraySlots43, arrayPosition43 := arr.View(), 0; arrayPosition43 < arraySlots43.Len(); arrayPosition43++ {
			z := arraySlots43.At(arrayPosition43)
			if z != nil {
				args = append(args, z.ReadValue())
			}
		}
	} else {
		return nil, data.NewErrorThrow(nil, fmt.Errorf(
			"Illuminate\\Support\\Collection: spread callback requires array items."))
	}
	return append(args, key), nil
}

// collectionMapInto 对齐 EnumeratesValues::mapInto：new $class($value, $key)。
func collectionMapInto(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	class, ctl := classNameArg(kit.Arg(ctx, 0))
	if ctl != nil {
		return nil, ctl
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		obj, ctl := instantiate(ctx, class, []data.Value{e.value, e.key})
		if ctl != nil {
			return nil, ctl
		}
		setEntry(out, e.keyStr, obj)
	}
	return newCollectionInstance(ctx, out)
}

// collectionPipeInto 对齐 EnumeratesValues::pipeInto：new $class($this)。
func collectionPipeInto(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	class, ctl := classNameArg(kit.Arg(ctx, 0))
	if ctl != nil {
		return nil, ctl
	}
	return instantiate(ctx, class, []data.Value{cv})
}

// collectionPipeThrough 对齐 EnumeratesValues::pipeThrough：把 $this 依次喂给每个回调。
func collectionPipeThrough(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	callbacks, ctl := arrFromValue(ctx, kit.Arg(ctx, 0), 0)
	if ctl != nil {
		return nil, ctl
	}
	var carry data.Value = cv
	for _, e := range toEntries(callbacks) {
		ret, ctl := callValue(ctx, e.value, carry)
		if ctl != nil {
			return nil, ctl
		}
		if ret != nil {
			carry = asValue(ret)
		}
	}
	return carry, nil
}

// collectionMapToGroups 对齐 EnumeratesValues::mapToGroups：先 mapToDictionary 再把每组包成集合。
func collectionMapToGroups(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	groups, ctl := buildDictionary(ctx, collectionItems(cv), kit.Arg(ctx, 0))
	if ctl != nil {
		return nil, ctl
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(groups) {
		sub, ctl := newCollectionInstance(ctx, e.value)
		if ctl != nil {
			return nil, ctl
		}
		setEntry(out, e.keyStr, sub)
	}
	return newCollectionInstance(ctx, out)
}

// buildDictionary 对齐 EnumeratesValues::mapToDictionary：
// 回调返回的 [key => value] 追加进同键分组（PHP 是 $groups[$key][] = $value）。
func buildDictionary(ctx data.Context, items *data.ArrayValue, cb data.Value) (*data.ArrayValue, data.Control) {
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(items) {
		ret, ctl := callValue(ctx, cb, e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		keyStr, value, ok := dictionaryPair(asValue(ret))
		if !ok {
			continue
		}
		sub := dictionaryBucket(out, keyStr)
		sub.AppendValue(value)
	}
	return out, nil
}

// dictionaryPair 取回调返回值里的首个键值对（对齐 PHP 的 key()/reset()）。
// 关联数组字面量在 origami 里可能是 ObjectValue，所以要两种表示都支持。
func dictionaryPair(v data.Value) (string, data.Value, bool) {
	switch pair := unwrapValue(v).(type) {
	case *data.ArrayValue:
		for arraySlots44, i := pair.View(), 0; i < arraySlots44.Len(); i++ {
			z := arraySlots44.At(i)
			if z == nil {
				continue
			}
			switch {
			case z.EmptyStrKey:
				return "", z.ReadValue(), true
			case z.Name != "":
				return z.Name, z.ReadValue(), true
			default:
				return data.IntArrayKeyName(i), z.ReadValue(), true
			}
		}

	}
	return "", nil, false
}

// dictionaryBucket 取出/新建分组数组（setEntry 会按 int 键名还原成整数键，保持 list 语义）。
func dictionaryBucket(out *data.ArrayValue, keyStr string) *data.ArrayValue {
	if z, ok := out.LookupZValByStringKey(keyStr); ok && z != nil {
		if sub, ok := z.ReadValue().(*data.ArrayValue); ok && sub != nil {
			return sub
		}
	}
	sub := data.NewArrayValue(nil).(*data.ArrayValue)
	setEntry(out, keyStr, sub)
	return sub
}

// ---- 调试 ----

// collectionDump 对齐 EnumeratesValues::dump()：dump(...$this->all()) 后返回 $this。
func collectionDump(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	dumpFn, ok := ctx.GetVM().GetFunc("dump")
	if !ok || dumpFn == nil {
		return cv, nil
	}
	items := collectionItems(cv)
	args := make([]data.Value, 0, items.Len())
	for arraySlots45, arrayPosition45 := items.View(), 0; arrayPosition45 < arraySlots45.Len(); arrayPosition45++ {
		z := arraySlots45.At(arrayPosition45)
		if z != nil && z.ReadValue() != nil {
			args = append(args, z.ReadValue())
		}
	}
	callCtx := ctx.CreateContext(dumpFn.GetVariables())
	if ctl := data.BindDeclaredArgs(callCtx, dumpFn, args); ctl != nil {
		return nil, ctl
	}
	if _, ctl := dumpFn.Call(callCtx); ctl != nil {
		return nil, ctl
	}
	return cv, nil
}

// collectionDd 对齐 EnumeratesValues::dd()：dump 之后结束当前脚本/请求。
func collectionDd(ctx data.Context) (data.GetValue, data.Control) {
	if _, ctl := collectionDump(ctx); ctl != nil {
		return nil, ctl
	}
	return data.NewNullValue(), data.NewExitControl(1)
}

// ---- 字符串 / JSON ----

// collectionToString 对齐 EnumeratesValues::__toString。
func collectionToString(ctx data.Context) (data.GetValue, data.Control) {
	encoded, ctl := collectionJSON(ctx)
	if ctl != nil {
		return nil, ctl
	}
	s, ok := encoded.(*data.StringValue)
	if !ok {
		return encoded, nil
	}
	if collectionEscapeFlag(ctx) {
		return data.NewStringValue(escapeHTML(s.Value)), nil
	}
	return encoded, nil
}

// collectionEscapeWhenCastingToString 对齐 EnumeratesValues::escapeWhenCastingToString。
func collectionEscapeWhenCastingToString(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	escape := true
	if v := kit.Arg(ctx, 0); v != nil && !isNull(v) {
		escape = kit.Truthy(v)
	}
	_ = cv.SetProperty("escapeWhenCastingToString", data.NewBoolValue(escape))
	return cv, nil
}

func collectionEscapeFlag(ctx data.Context) bool {
	cv := kit.Receiver(ctx)
	if cv == nil {
		return false
	}
	v, _ := cv.GetProperty("escapeWhenCastingToString")
	return v != nil && !isNull(v) && kit.Truthy(v)
}

func collectionJSON(ctx data.Context) (data.GetValue, data.Control) {
	// 直接复用 collectionToJson：不能经 withArgs 造上下文，
	// 那会丢 $this（CreateContext 返回的不是 ClassMethodContext）。
	return collectionToJson(ctx)
}

// escapeHTML 对齐 Illuminate\Support\e()（htmlspecialchars, ENT_QUOTES）。
func escapeHTML(s string) string {
	var b strings.Builder
	if !strings.ContainsAny(s, "&<>\"'") {
		return s
	}
	b.Grow(len(s) + 16)
	for _, r := range s {
		switch r {
		case '&':
			b.WriteString("&amp;")
		case '<':
			b.WriteString("&lt;")
		case '>':
			b.WriteString("&gt;")
		case '"':
			b.WriteString("&quot;")
		case '\'':
			b.WriteString("&#039;")
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// collectionFromJson 对齐 EnumeratesValues::fromJson。
func collectionFromJson(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	if vm == nil {
		return nil, nil
	}
	fn, ok := vm.GetFunc("json_decode")
	if !ok || fn == nil {
		return nil, nil
	}
	args := []data.Value{kit.Arg(ctx, 0), data.NewBoolValue(true), kit.Arg(ctx, 1), kit.Arg(ctx, 2)}
	for i, a := range args {
		if a == nil {
			args[i] = data.NewNullValue()
		}
	}
	nctx := ctx.CreateContext(fn.GetVariables())
	if ctl := data.BindDeclaredArgs(nctx, fn, args); ctl != nil {
		return nil, ctl
	}
	decoded, ctl := fn.Call(nctx)
	if ctl != nil {
		return nil, ctl
	}
	return newCollectionInstance(ctx, asValue(decoded))
}

// ---- 小工具 ----

// instantiate 按类名 new 一个对象，args 走构造函数的形参绑定。
func instantiate(ctx data.Context, class string, args []data.Value) (data.Value, data.Control) {
	vm := ctx.GetVM()
	if vm == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("no VM"))
	}
	stmt, ctl := vm.GetOrLoadClass(class)
	if ctl != nil {
		return nil, ctl
	}
	cv := data.NewClassValue(stmt, ctx.CreateBaseContext())
	ctor := stmt.GetConstruct()
	if ctor == nil {
		if m, ok := cv.GetMethod("__construct"); ok {
			ctor = m
		}
	}
	if ctor == nil {
		return cv, nil
	}
	nctx := cv.CreateContext(ctor.GetVariables())
	if ctl := data.BindDeclaredArgs(nctx, ctor, args); ctl != nil {
		return nil, ctl
	}
	if _, ctl := ctor.Call(nctx); ctl != nil {
		return nil, ctl
	}
	return cv, nil
}

func classNameArg(v data.Value) (string, data.Control) {
	v = unwrapValue(v)
	switch t := v.(type) {
	case *data.StringValue:
		return t.Value, nil
	case *data.ClassValue:
		if t.Class != nil {
			return t.Class.GetName(), nil
		}
	}
	return "", data.NewErrorThrow(nil, fmt.Errorf("Collection: class name expected"))
}

func asValue(v data.GetValue) data.Value {
	if v == nil {
		return data.NewNullValue()
	}
	if val, ok := v.(data.Value); ok {
		return val
	}
	return data.NewNullValue()
}
