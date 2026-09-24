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

func (i *EnumerableInterface) GetName() string                     { return enumerableName }
func (i *EnumerableInterface) GetExtends() []string                { return nil }
func (i *EnumerableInterface) GetMethod(string) (data.Method, bool) { return nil, false }
func (i *EnumerableInterface) GetMethods() []data.Method            { return nil }
func (i *EnumerableInterface) GetFrom() data.From                   { return nil }
func (i *EnumerableInterface) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

const collectionName = "Illuminate\\Support\\Collection"

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
	stat := func(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[data.MethodLookupKey(name)] = kit.StaticMethod(name, params, -1, fn)
	}
	inst("__construct", []string{"items"}, collectionConstruct)
	stat("make", []string{"items"}, collectionMake)
	stat("empty", nil, collectionEmpty)
	stat("wrap", []string{"value"}, collectionWrap)
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
	inst("push", []string{"values"}, collectionPush)
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
	inst("__get", []string{"key"}, collectionMagicGet)
	registerCollectionMore(c)
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

func collectionToArray(ctx data.Context) (data.GetValue, data.Control) {
	return collectionAll(ctx)
}

func collectionToJson(ctx data.Context) (data.GetValue, data.Control) {
	arr, ctl := collectionAll(ctx)
	if ctl != nil {
		return nil, ctl
	}
	encoded, ok, jctl := php.JsonEncode(ctx, arr.(data.Value))
	if jctl != nil {
		return nil, jctl
	}
	if !ok {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(encoded), nil
}

func collectionJsonSerialize(ctx data.Context) (data.GetValue, data.Control) {
	return collectionAll(ctx)
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
	return arrGet(withArgs(ctx, collectionItems(cv), key, def))
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
		if ret != nil {
			if b, ok := ret.(data.AsBool); ok {
				if okv, _ := b.AsBool(); !okv {
					break
				}
			}
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
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		itemVal, ok := dataGetPath(e.value, keyToString(key))
		if !ok {
			continue
		}
		match := false
		if val == nil || isNull(val) {
			// where($key, $value)
			match = itemVal.AsString() == keyToString(op)
		} else {
			match = compareOp(itemVal, keyToString(op), val)
		}
		if match {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return newCollectionInstance(ctx, out)
}

func compareOp(left data.Value, op string, right data.Value) bool {
	ls, rs := left.AsString(), right.AsString()
	switch op {
	case "=", "==":
		return ls == rs
	case "!=", "<>":
		return ls != rs
	case ">":
		return ls > rs
	case "<":
		return ls < rs
	case ">=":
		return ls >= rs
	case "<=":
		return ls <= rs
	default:
		return ls == rs
	}
}

func collectionUnique(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	seen := map[string]bool{}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		k := e.value.AsString()
		if seen[k] {
			continue
		}
		seen[k] = true
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
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

func compareSortKeys(a, b data.Value) int {
	a = kit.Unwrap(a)
	b = kit.Unwrap(b)
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}
	if ai, ok := a.(*data.IntValue); ok {
		if bi, ok := b.(*data.IntValue); ok {
			switch {
			case ai.Value < bi.Value:
				return -1
			case ai.Value > bi.Value:
				return 1
			default:
				return 0
			}
		}
	}
	if af, ok := a.(*data.FloatValue); ok {
		if bf, ok := b.(*data.FloatValue); ok {
			switch {
			case af.Value < bf.Value:
				return -1
			case af.Value > bf.Value:
				return 1
			default:
				return 0
			}
		}
	}
	as, bs := a.AsString(), b.AsString()
	switch {
	case as < bs:
		return -1
	case as > bs:
		return 1
	default:
		return 0
	}
}

func collectionGroupBy(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	groupBy, _ := ctx.GetIndexValue(0)
	groups := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		gk := ""
		if v, ok := dataGetPath(e.value, keyToString(groupBy)); ok {
			gk = keyToString(v)
		}
		g, ok := dataGetPath(groups, gk)
		var garr *data.ArrayValue
		if ok {
			garr, _ = g.(*data.ArrayValue)
		}
		if garr == nil {
			garr = data.NewArrayValue(nil).(*data.ArrayValue)
			setEntry(groups, gk, garr)
		}
		garr.List = append(garr.List, data.NewZVal(e.value))
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
	return collectionAll(ctx)
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
	if _, ok := collectionProxies[data.MethodLookupKey(key)]; ok {
		return newHigherOrderProxy(ctx, cv, key)
	}
	return data.NewNullValue(), nil
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
