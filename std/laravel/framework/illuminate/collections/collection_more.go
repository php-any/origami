package collections

import (
	"fmt"
	"math/rand"
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/laravel/framework/illuminate/conditionable"
	"github.com/php-any/origami/std/laravel/framework/internal/kit"
)

func registerCollectionMore(c *CollectionClass) {
	add := func(name string, params []string, fn func(data.Context) (data.GetValue, data.Control)) {
		c.methods[data.MethodLookupKey(name)] = kit.InstanceMethod(name, params, fn)
	}
	add("add", []string{"key", "value"}, collectionAdd)
	add("after", []string{"value", "key"}, collectionAfter)
	add("before", []string{"value", "key"}, collectionBefore)
	add("chunk", []string{"size"}, collectionChunk)
	add("chunkWhile", []string{"callback"}, collectionChunkWhile)
	add("collapse", nil, collectionCollapse)
	add("collapseWithKeys", nil, collectionCollapseWithKeys)
	add("combine", []string{"values"}, collectionCombine)
	add("containsManyItems", []string{"items"}, collectionContainsManyItems)
	add("containsOneItem", []string{"items"}, collectionContainsOneItem)
	add("containsStrict", []string{"key", "operator", "value"}, collectionContainsStrict)
	add("countBy", []string{"callback"}, collectionCountBy)
	add("crossJoin", []string{"lists", "lists2", "lists3", "lists4"}, collectionCrossJoin)
	add("diffAssoc", []string{"items"}, collectionDiffAssoc)
	add("diffAssocUsing", []string{"items", "callback"}, collectionDiffAssocUsing)
	add("diffKeysUsing", []string{"items", "callback"}, collectionDiffKeysUsing)
	add("diffUsing", []string{"items", "callback"}, collectionDiffUsing)
	add("doesntContain", []string{"key", "operator", "value"}, collectionDoesntContain)
	add("doesntContainStrict", []string{"key", "operator", "value"}, collectionDoesntContainStrict)
	add("dot", nil, collectionDot)
	add("duplicates", []string{"callback", "strict"}, collectionDuplicates)
	add("duplicatesStrict", []string{"callback"}, collectionDuplicatesStrict)
	add("firstOrFail", []string{"callback"}, collectionFirstOrFail)
	add("flip", nil, collectionFlip)
	add("forget", []string{"keys"}, collectionForget)
	add("getOrPut", []string{"key", "value"}, collectionGetOrPut)
	add("has", []string{"key"}, collectionHas)
	add("hasAny", []string{"keys"}, collectionHasAny)
	add("hasSole", []string{"key", "operator", "value"}, collectionHasSole)
	add("intersect", []string{"items"}, collectionIntersect)
	add("intersectAssoc", []string{"items"}, collectionIntersectAssoc)
	add("intersectAssocUsing", []string{"items", "callback"}, collectionIntersectAssocUsing)
	add("intersectByKeys", []string{"items"}, collectionIntersectByKeys)
	add("intersectUsing", []string{"items", "callback"}, collectionIntersectUsing)
	add("lazy", nil, collectionLazy)
	add("mapToDictionary", []string{"callback"}, collectionMapToDictionary)
	add("median", []string{"key"}, collectionMedian)
	add("mergeRecursive", []string{"items"}, collectionMergeRecursive)
	add("mode", []string{"key"}, collectionMode)
	add("multiply", []string{"multiplier"}, collectionMultiply)
	add("nth", []string{"step", "offset"}, collectionNth)
	add("pad", []string{"size", "value"}, collectionPad)
	add("prepend", []string{"value", "key"}, collectionPrepend)
	add("pull", []string{"key", "default"}, collectionPull)
	add("random", []string{"number"}, collectionRandom)
	add("reject", []string{"callback"}, collectionReject)
	add("whereInstanceOf", []string{"type"}, collectionWhereInstanceOf)
	add("whereNull", []string{"key"}, collectionWhereNull)
	add("whereNotNull", []string{"key"}, collectionWhereNotNull)
	add("whereIn", []string{"key", "values", "strict"}, collectionWhereIn)
	add("whereNotIn", []string{"key", "values", "strict"}, collectionWhereNotIn)
	add("reduce", []string{"callback", "initial"}, collectionReduce)
	add("every", []string{"callback"}, collectionEvery)
	add("some", []string{"key", "operator", "value"}, collectionSome)
	add("avg", []string{"callback"}, collectionAvg)
	add("average", []string{"callback"}, collectionAvg)
	add("sum", []string{"callback"}, collectionSum)
	add("max", []string{"callback"}, collectionMax)
	add("min", []string{"callback"}, collectionMin)
	add("tap", []string{"callback"}, collectionTap)
	add("pipe", []string{"callback"}, collectionPipe)
	add("whenEmpty", []string{"callback", "default"}, collectionWhenEmpty)
	add("whenNotEmpty", []string{"callback", "default"}, collectionWhenNotEmpty)
	add("unlessEmpty", []string{"callback", "default"}, collectionUnlessEmpty)
	add("unlessNotEmpty", []string{"callback", "default"}, collectionUnlessNotEmpty)
	add("replace", []string{"items"}, collectionReplace)
	add("replaceRecursive", []string{"items"}, collectionReplaceRecursive)
	add("reverse", nil, collectionReverse)
	add("search", []string{"value", "strict"}, collectionSearch)
	add("select", []string{"keys"}, collectionSelect)
	add("shift", nil, collectionShift)
	add("shuffle", []string{"seed"}, collectionShuffle)
	add("skip", []string{"count"}, collectionSkip)
	add("skipUntil", []string{"callback"}, collectionSkipUntil)
	add("skipWhile", []string{"callback"}, collectionSkipWhile)
	add("slice", []string{"offset", "length"}, collectionSlice)
	add("sliding", []string{"size", "step"}, collectionSliding)
	add("sole", []string{"key", "operator", "value"}, collectionSole)
	c.methods["sortbydesc"] = kit.InstanceMethodOpt("sortByDesc", []string{"callback", "options"}, 1, collectionSortByDesc)
	add("sortDesc", []string{"options"}, collectionSortDesc)
	add("sortKeys", nil, collectionSortKeys)
	add("sortKeysDesc", nil, collectionSortKeysDesc)
	add("sortKeysUsing", []string{"callback"}, collectionSortKeysUsing)
	add("splice", []string{"offset", "length", "replacement"}, collectionSplice)
	add("split", []string{"numberOfGroups"}, collectionSplit)
	add("splitIn", []string{"numberOfGroups"}, collectionSplitIn)
	add("take", []string{"limit"}, collectionTake)
	add("takeUntil", []string{"callback"}, collectionTakeUntil)
	add("takeWhile", []string{"callback"}, collectionTakeWhile)
	add("transform", []string{"callback"}, collectionTransform)
	add("undot", nil, collectionUndot)
	add("union", []string{"items"}, collectionUnion)
	add("unshift", []string{"values"}, collectionUnshift)
	add("zip", []string{"items", "items2", "items3", "items4", "items5"}, collectionZip)
}

func collectionContainsItems(ctx data.Context, cv *data.ClassValue, strict bool) (data.GetValue, data.Control) {
	key, _ := ctx.GetIndexValue(0)
	op, _ := ctx.GetIndexValue(1)
	val, _ := ctx.GetIndexValue(2)
	items := collectionItems(cv)
	eq := valueLooseEqual
	if strict {
		eq = valueStrictEqual
	}
	if val == nil || isNull(val) {
		if op == nil || isNull(op) {
			if isCallable(key) {
				v, ctl := arrFirst(withArgs(ctx, collectionItems(cv), key, data.NewNullValue()))
				if ctl != nil {
					return nil, ctl
				}
				if v == nil {
					return data.NewBoolValue(false), nil
				}
				if val, ok := v.(data.Value); ok {
					return data.NewBoolValue(!isNull(val)), nil
				}
				return data.NewBoolValue(true), nil
			}
			for _, e := range toEntries(items) {
				if eq(e.value, key) {
					return data.NewBoolValue(true), nil
				}
			}
			return data.NewBoolValue(false), nil
		}
		val = op
		op = data.NewStringValue("==")
	}
	for _, e := range toEntries(items) {
		itemVal, ok := dataGetPath(e.value, keyToString(key))
		if !ok {
			continue
		}
		if strict {
			if valueStrictEqual(itemVal, val) {
				return data.NewBoolValue(true), nil
			}
		} else if compareOp(itemVal, keyToString(op), val) {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func isCallable(v data.Value) bool {
	switch v.(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		return true
	default:
		return false
	}
}

func valueLooseEqual(a, b data.Value) bool {
	if a == nil || b == nil {
		return a == b
	}
	a, b = unwrapValue(a), unwrapValue(b)
	if isNull(a) && isNull(b) {
		return true
	}
	return a.AsString() == b.AsString()
}

func valueStrictEqual(a, b data.Value) bool {
	if a == nil || b == nil {
		return a == b
	}
	a, b = unwrapValue(a), unwrapValue(b)
	if fmt.Sprintf("%T", a) != fmt.Sprintf("%T", b) {
		return false
	}
	switch ta := a.(type) {
	case *data.IntValue:
		tb, ok := b.(*data.IntValue)
		return ok && ta.Value == tb.Value
	case *data.FloatValue:
		tb, ok := b.(*data.FloatValue)
		return ok && ta.Value == tb.Value
	case *data.StringValue:
		tb, ok := b.(*data.StringValue)
		return ok && ta.Value == tb.Value
	case *data.BoolValue:
		tb, ok := b.(*data.BoolValue)
		return ok && ta.Value == tb.Value
	case *data.NullValue:
		_, ok := b.(*data.NullValue)
		return ok
	default:
		return a == b
	}
}

func collectionAdd(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, _ := ctx.GetIndexValue(0)
	val, _ := ctx.GetIndexValue(1)
	items := collectionItems(cv)
	if _, ok := dataGetPath(items, keyToString(key)); !ok {
		setEntry(items, keyToString(key), val)
		_ = cv.SetProperty("items", items)
	}
	return cv, nil
}

func collectionAfter(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	val, _ := ctx.GetIndexValue(0)
	searchKey, _ := ctx.GetIndexValue(1)
	return insertRelative(ctx, cv, val, searchKey, false)
}

func collectionBefore(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	val, _ := ctx.GetIndexValue(0)
	searchKey, _ := ctx.GetIndexValue(1)
	return insertRelative(ctx, cv, val, searchKey, true)
}

func insertRelative(ctx data.Context, cv *data.ClassValue, val, searchKey data.Value, before bool) (data.GetValue, data.Control) {
	entries := toEntries(collectionItems(cv))
	sk := keyToString(searchKey)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	inserted := false
	for _, e := range entries {
		if !inserted && e.keyStr == sk {
			if before {
				out.List = append(out.List, data.NewZVal(val))
				out.List = append(out.List, data.NewZVal(e.value))
			} else {
				out.List = append(out.List, data.NewZVal(e.value))
				out.List = append(out.List, data.NewZVal(val))
			}
			inserted = true
			continue
		}
		setEntry(out, e.keyStr, e.value)
	}
	if !inserted {
		out.List = append(out.List, data.NewZVal(val))
	}
	return newCollectionInstance(ctx, out)
}

func collectionChunk(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	sizeArg, _ := ctx.GetIndexValue(0)
	size, ctl := intFromValue(sizeArg)
	if ctl != nil {
		return nil, ctl
	}
	if size < 1 {
		size = 1
	}
	entries := toEntries(collectionItems(cv))
	chunks := data.NewArrayValue(nil).(*data.ArrayValue)
	cur := data.NewArrayValue(nil).(*data.ArrayValue)
	n := 0
	for _, e := range entries {
		setEntry(cur, e.keyStr, e.value)
		n++
		if n >= size {
			inst, err := newCollectionInstance(ctx, cur)
			if err != nil {
				return nil, err
			}
			chunks.List = append(chunks.List, data.NewZVal(inst))
			cur = data.NewArrayValue(nil).(*data.ArrayValue)
			n = 0
		}
	}
	if n > 0 {
		inst, err := newCollectionInstance(ctx, cur)
		if err != nil {
			return nil, err
		}
		chunks.List = append(chunks.List, data.NewZVal(inst))
	}
	return newCollectionInstance(ctx, chunks)
}

func collectionChunkWhile(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	entries := toEntries(collectionItems(cv))
	chunks := data.NewArrayValue(nil).(*data.ArrayValue)
	var cur *data.ArrayValue
	for _, e := range entries {
		if cur == nil {
			cur = data.NewArrayValue(nil).(*data.ArrayValue)
		}
		setEntry(cur, e.keyStr, e.value)
		ok, ctl := callBool(ctx, cb, e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		if !ok {
			inst, err := newCollectionInstance(ctx, cur)
			if err != nil {
				return nil, err
			}
			chunks.List = append(chunks.List, data.NewZVal(inst))
			cur = nil
		}
	}
	if cur != nil && len(cur.List) > 0 {
		inst, err := newCollectionInstance(ctx, cur)
		if err != nil {
			return nil, err
		}
		chunks.List = append(chunks.List, data.NewZVal(inst))
	}
	return newCollectionInstance(ctx, chunks)
}

func collectionCollapse(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	collapsed, ctl := arrCollapse(withArgs(ctx, collectionItems(cv)))
	if ctl != nil {
		return nil, ctl
	}
	return newCollectionInstance(ctx, collapsed.(data.Value))
}

func collectionCollapseWithKeys(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		layer, ok, ctl := collapseLayer(ctx, e.value)
		if ctl != nil {
			return nil, ctl
		}
		if !ok {
			setEntry(out, e.keyStr, e.value)
			continue
		}
		for _, ne := range layer {
			setEntry(out, ne.keyStr, ne.value)
		}
	}
	return newCollectionInstance(ctx, out)
}

func collectionCombine(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	values, _ := ctx.GetIndexValue(0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	keys := toEntries(collectionItems(cv))
	vals := toEntries(values)
	for i, k := range keys {
		if i >= len(vals) {
			break
		}
		setEntry(out, keyToString(k.value), vals[i].value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionContainsManyItems(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	items, _ := ctx.GetIndexValue(0)
	for _, need := range toEntries(items) {
		found := false
		for _, e := range toEntries(collectionItems(cv)) {
			if valueLooseEqual(e.value, need.value) {
				found = true
				break
			}
		}
		if !found {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func collectionContainsOneItem(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	items, _ := ctx.GetIndexValue(0)
	for _, need := range toEntries(items) {
		for _, e := range toEntries(collectionItems(cv)) {
			if valueLooseEqual(e.value, need.value) {
				return data.NewBoolValue(true), nil
			}
		}
	}
	return data.NewBoolValue(false), nil
}

func collectionContainsStrict(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return collectionContainsItems(ctx, cv, true)
}

func collectionDoesntContain(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := collectionContains(ctx)
	if ctl != nil {
		return nil, ctl
	}
	b, _ := v.(data.AsBool).AsBool()
	return data.NewBoolValue(!b), nil
}

func collectionDoesntContainStrict(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := collectionContainsStrict(ctx)
	if ctl != nil {
		return nil, ctl
	}
	b, _ := v.(data.AsBool).AsBool()
	return data.NewBoolValue(!b), nil
}

func collectionCountBy(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	counts := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		k := e.value.AsString()
		if cb != nil && !isNull(cb) {
			ret, ctl := callValue(ctx, cb, e.value, e.key)
			if ctl != nil {
				return nil, ctl
			}
			if ret != nil {
				if v, ok := ret.(data.Value); ok {
					k = keyToString(v)
				}
			}
		}
		if z, ok := counts.LookupZValByStringKey(k); ok && z != nil {
			if iv, ok := z.Value.(*data.IntValue); ok {
				iv.Value++
			} else {
				setEntry(counts, k, data.NewIntValue(1))
			}
		} else {
			setEntry(counts, k, data.NewIntValue(1))
		}
	}
	return newCollectionInstance(ctx, counts)
}

func collectionCrossJoin(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	lists := [][]data.Value{entriesToValues(toEntries(collectionItems(cv)))}
	for i := 1; ; i++ {
		v, ok := ctx.GetIndexValue(i)
		if !ok || v == nil || isNull(v) {
			break
		}
		other, ctl := getArrayableItems(ctx, v)
		if ctl != nil {
			return nil, ctl
		}
		lists = append(lists, entriesToValues(toEntries(other)))
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	crossJoinInto(out, lists, 0, nil)
	return newCollectionInstance(ctx, out)
}

func entriesToValues(entries []kv) []data.Value {
	out := make([]data.Value, len(entries))
	for i, e := range entries {
		out[i] = e.value
	}
	return out
}

func crossJoinInto(out *data.ArrayValue, lists [][]data.Value, depth int, prefix []data.Value) {
	if depth >= len(lists) {
		row := data.NewArrayValue(nil).(*data.ArrayValue)
		for _, v := range prefix {
			row.List = append(row.List, data.NewZVal(v))
		}
		out.List = append(out.List, data.NewZVal(row))
		return
	}
	for _, v := range lists[depth] {
		crossJoinInto(out, lists, depth+1, append(prefix, v))
	}
}

func collectionDiffAssoc(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	other, _ := ctx.GetIndexValue(0)
	otherItems, ctl := getArrayableItems(ctx, other)
	if ctl != nil {
		return nil, ctl
	}
	otherMap := map[string]data.Value{}
	for _, e := range toEntries(otherItems) {
		otherMap[e.keyStr] = e.value
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		ov, ok := otherMap[e.keyStr]
		if ok && valueLooseEqual(ov, e.value) {
			continue
		}
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionDiffUsing(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	other, _ := ctx.GetIndexValue(0)
	cb, _ := ctx.GetIndexValue(1)
	otherItems, ctl := getArrayableItems(ctx, other)
	if ctl != nil {
		return nil, ctl
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
outer:
	for _, e := range toEntries(collectionItems(cv)) {
		for _, o := range toEntries(otherItems) {
			cmp, ctl := callValue(ctx, cb, e.value, o.value)
			if ctl != nil {
				return nil, ctl
			}
			if iv, ok := cmp.(*data.IntValue); ok && iv.Value == 0 {
				continue outer
			}
		}
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionDiffAssocUsing(ctx data.Context) (data.GetValue, data.Control) {
	return collectionDiffUsing(ctx)
}

func collectionDiffKeysUsing(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	other, _ := ctx.GetIndexValue(0)
	cb, _ := ctx.GetIndexValue(1)
	otherItems, ctl := getArrayableItems(ctx, other)
	if ctl != nil {
		return nil, ctl
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		skip := false
		for _, o := range toEntries(otherItems) {
			cmp, ctl := callValue(ctx, cb, data.NewStringValue(e.keyStr), o.key)
			if ctl != nil {
				return nil, ctl
			}
			if iv, ok := cmp.(*data.IntValue); ok && iv.Value == 0 {
				skip = true
				break
			}
		}
		if !skip {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return newCollectionInstance(ctx, out)
}

func collectionDot(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	dotted, ctl := arrDot(withArgs(ctx, collectionItems(cv)))
	if ctl != nil {
		return nil, ctl
	}
	return newCollectionInstance(ctx, dotted.(data.Value))
}

func collectionUndot(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	undotted, ctl := arrUndot(withArgs(ctx, collectionItems(cv)))
	if ctl != nil {
		return nil, ctl
	}
	return newCollectionInstance(ctx, undotted.(data.Value))
}

func collectionDuplicates(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	strict := false
	if s, ok := ctx.GetIndexValue(1); ok && s != nil && truthy(s) {
		strict = true
	}
	return duplicatesCollection(ctx, cv, cb, strict)
}

func collectionDuplicatesStrict(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	return duplicatesCollection(ctx, cv, cb, true)
}

func duplicatesCollection(ctx data.Context, cv *data.ClassValue, cb data.Value, strict bool) (data.GetValue, data.Control) {
	_ = strict
	seen := map[string]int{}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		k := e.value
		if cb != nil && !isNull(cb) {
			ret, ctl := callValue(ctx, cb, e.value, e.key)
			if ctl != nil {
				return nil, ctl
			}
			if ret != nil {
				if v, ok := ret.(data.Value); ok {
					k = v
				}
			}
		}
		key := keyToString(k)
		seen[key]++
		if seen[key] > 1 {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return newCollectionInstance(ctx, out)
}

func collectionFirstOrFail(ctx data.Context) (data.GetValue, data.Control) {
	v, ctl := collectionFirst(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if v == nil {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Item not found."), "Illuminate\\Support\\ItemNotFoundException")
	}
	if val, ok := v.(data.Value); ok && isNull(val) {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Item not found."), "Illuminate\\Support\\ItemNotFoundException")
	}
	return v, nil
}

func collectionFlip(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		setEntry(out, keyToString(e.value), e.key)
	}
	return newCollectionInstance(ctx, out)
}

func collectionForget(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	keys, _ := ctx.GetIndexValue(0)
	items := collectionItems(cv)
	for _, k := range keysToStrings(keys) {
		dataForgetPath(items, k)
	}
	_ = cv.SetProperty("items", items)
	return cv, nil
}

func collectionGetOrPut(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, _ := ctx.GetIndexValue(0)
	def, _ := ctx.GetIndexValue(1)
	items := collectionItems(cv)
	if v, ok := dataGetPath(items, keyToString(key)); ok {
		return v, nil
	}
	val, ctl := laravelValue(ctx, def)
	if ctl != nil {
		return nil, ctl
	}
	if v, ok := val.(data.Value); ok {
		setEntry(items, keyToString(key), v)
		_ = cv.SetProperty("items", items)
		return v, nil
	}
	return data.NewNullValue(), nil
}

func collectionHas(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, _ := ctx.GetIndexValue(0)
	_, ok := dataGetPath(collectionItems(cv), keyToString(key))
	return data.NewBoolValue(ok), nil
}

func collectionHasAny(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	keys, _ := ctx.GetIndexValue(0)
	items := collectionItems(cv)
	for _, k := range keysToStrings(keys) {
		if _, ok := dataGetPath(items, k); ok {
			return data.NewBoolValue(true), nil
		}
	}
	return data.NewBoolValue(false), nil
}

func collectionHasSole(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	n := countWhereMatches(ctx, cv)
	return data.NewBoolValue(n == 1), nil
}

func collectionIntersect(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	other, _ := ctx.GetIndexValue(0)
	seen := map[string]struct{}{}
	for _, e := range toEntries(collectionItems(cv)) {
		seen[e.value.AsString()] = struct{}{}
	}
	otherItems, ctl := getArrayableItems(ctx, other)
	if ctl != nil {
		return nil, ctl
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(otherItems) {
		if _, ok := seen[e.value.AsString()]; ok {
			out.List = append(out.List, data.NewZVal(e.value))
		}
	}
	return newCollectionInstance(ctx, out)
}

func collectionIntersectAssoc(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	other, _ := ctx.GetIndexValue(0)
	self := collectionItems(cv)
	otherItems, ctl := getArrayableItems(ctx, other)
	if ctl != nil {
		return nil, ctl
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(otherItems) {
		if v, ok := dataGetPath(self, e.keyStr); ok && valueLooseEqual(v, e.value) {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return newCollectionInstance(ctx, out)
}

func collectionIntersectUsing(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	other, _ := ctx.GetIndexValue(0)
	cb, _ := ctx.GetIndexValue(1)
	otherItems, ctl := getArrayableItems(ctx, other)
	if ctl != nil {
		return nil, ctl
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		for _, o := range toEntries(otherItems) {
			cmp, ctl := callValue(ctx, cb, e.value, o.value)
			if ctl != nil {
				return nil, ctl
			}
			if iv, ok := cmp.(*data.IntValue); ok && iv.Value == 0 {
				out.List = append(out.List, data.NewZVal(e.value))
				break
			}
		}
	}
	return newCollectionInstance(ctx, out)
}

func collectionIntersectAssocUsing(ctx data.Context) (data.GetValue, data.Control) {
	return collectionIntersectUsing(ctx)
}

func collectionIntersectByKeys(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	other, _ := ctx.GetIndexValue(0)
	keys := map[string]struct{}{}
	otherItems, ctl := getArrayableItems(ctx, other)
	if ctl != nil {
		return nil, ctl
	}
	for _, e := range toEntries(otherItems) {
		keys[e.keyStr] = struct{}{}
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		if _, ok := keys[e.keyStr]; ok {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return newCollectionInstance(ctx, out)
}

// collectionLazy 暂无 LazyCollection 原生实现，返回自身即可满足 Enumerable 契约。
func collectionLazy(ctx data.Context) (data.GetValue, data.Control) {
	return collectionReceiver(ctx)
}

func collectionMapToDictionary(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		ret, ctl := callValue(ctx, cb, e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		if ret == nil {
			continue
		}
		if pair, ok := ret.(*data.ArrayValue); ok && len(pair.List) >= 2 {
			k := pair.List[0].Value
			v := pair.List[1].Value
			setEntry(out, keyToString(k), v)
		}
	}
	return newCollectionInstance(ctx, out)
}

func collectionMedian(ctx data.Context) (data.GetValue, data.Control) {
	return collectionNumericAggregate(ctx, "median")
}

func collectionMode(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, _ := ctx.GetIndexValue(0) // optional column key
	counts := map[string]int{}
	for _, e := range toEntries(collectionItems(cv)) {
		v := e.value
		if key != nil && !isNull(key) {
			if got, ok := dataGetPath(e.value, keyToString(key)); ok {
				v = got
			}
		}
		counts[v.AsString()]++
	}
	max := 0
	for _, c := range counts {
		if c > max {
			max = c
		}
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for s, c := range counts {
		if c == max {
			out.List = append(out.List, data.NewZVal(data.NewStringValue(s)))
		}
	}
	return out, nil
}

func collectionNumericAggregate(ctx data.Context, kind string) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, _ := ctx.GetIndexValue(0)
	nums := make([]float64, 0)
	for _, e := range toEntries(collectionItems(cv)) {
		v := e.value
		if key != nil && !isNull(key) {
			if got, ok := dataGetPath(e.value, keyToString(key)); ok {
				v = got
			}
		}
		if iv, ok := v.(*data.IntValue); ok {
			nums = append(nums, float64(iv.Value))
		} else if fv, ok := v.(*data.FloatValue); ok {
			nums = append(nums, fv.Value)
		}
	}
	if len(nums) == 0 {
		return data.NewNullValue(), nil
	}
	sort.Float64s(nums)
	if kind == "median" {
		mid := len(nums) / 2
		if len(nums)%2 == 0 {
			return data.NewFloatValue((nums[mid-1] + nums[mid]) / 2), nil
		}
		return data.NewFloatValue(nums[mid]), nil
	}
	return data.NewNullValue(), nil
}

func collectionMergeRecursive(ctx data.Context) (data.GetValue, data.Control) {
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
	mergeRecursiveInto(out, otherItems)
	return newCollectionInstance(ctx, out)
}

func mergeRecursiveInto(dst *data.ArrayValue, src data.Value) {
	for _, e := range toEntries(src) {
		if existing, ok := dataGetPath(dst, e.keyStr); ok {
			if ea, ok := existing.(*data.ArrayValue); ok {
				if sa, ok := e.value.(*data.ArrayValue); ok {
					mergeRecursiveInto(ea, sa)
					continue
				}
			}
		}
		setEntry(dst, e.keyStr, e.value)
	}
}

func collectionMultiply(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	mul, _ := ctx.GetIndexValue(0)
	factor := 1.0
	if iv, ok := mul.(*data.IntValue); ok {
		factor = float64(iv.Value)
	} else if fv, ok := mul.(*data.FloatValue); ok {
		factor = fv.Value
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		if iv, ok := e.value.(*data.IntValue); ok {
			setEntry(out, e.keyStr, data.NewIntValue(int(float64(iv.Value)*factor)))
		} else if fv, ok := e.value.(*data.FloatValue); ok {
			setEntry(out, e.keyStr, data.NewFloatValue(fv.Value*factor))
		} else {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return newCollectionInstance(ctx, out)
}

func collectionNth(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	stepArg, _ := ctx.GetIndexValue(0)
	step, ctl := intFromValue(stepArg)
	if ctl != nil {
		return nil, ctl
	}
	offset := 0
	if o, ok := ctx.GetIndexValue(1); ok && o != nil && !isNull(o) {
		offset, ctl = intFromValue(o)
		if ctl != nil {
			return nil, ctl
		}
	}
	if step < 1 {
		step = 1
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for i, e := range toEntries(collectionItems(cv)) {
		if i < offset {
			continue
		}
		if (i-offset)%step == 0 {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return newCollectionInstance(ctx, out)
}

func collectionPad(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	sizeArg, _ := ctx.GetIndexValue(0)
	size, ctl := intFromValue(sizeArg)
	if ctl != nil {
		return nil, ctl
	}
	val, _ := ctx.GetIndexValue(1)
	entries := toEntries(collectionItems(cv))
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range entries {
		setEntry(out, e.keyStr, e.value)
	}
	for len(out.List) < size {
		out.List = append(out.List, data.NewZVal(val))
	}
	return newCollectionInstance(ctx, out)
}

func collectionPrepend(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	val, _ := ctx.GetIndexValue(0)
	key, _ := ctx.GetIndexValue(1)
	prepended, ctl := arrPrepend(withArgs(ctx, collectionItems(cv), val, key))
	if ctl != nil {
		return nil, ctl
	}
	return newCollectionInstance(ctx, prepended.(data.Value))
}

func collectionPull(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key, _ := ctx.GetIndexValue(0)
	def, _ := ctx.GetIndexValue(1)
	items := collectionItems(cv)
	ks := keyToString(key)
	v, ok := dataGetPath(items, ks)
	dataForgetPath(items, ks)
	_ = cv.SetProperty("items", items)
	if ok {
		return v, nil
	}
	return laravelValue(ctx, def)
}

func collectionRandom(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	entries := toEntries(collectionItems(cv))
	if len(entries) == 0 {
		return data.NewNullValue(), nil
	}
	n := 1
	if num, ok := ctx.GetIndexValue(0); ok && num != nil && !isNull(num) {
		var err data.Control
		n, err = intFromValue(num)
		if err != nil {
			return nil, err
		}
	}
	if n <= 1 {
		return entries[rand.Intn(len(entries))].value, nil
	}
	// 简化：返回前 n 个随机项的新集合
	idx := rand.Perm(len(entries))
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	if n > len(entries) {
		n = len(entries)
	}
	for i := 0; i < n; i++ {
		e := entries[idx[i]]
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionReplace(ctx data.Context) (data.GetValue, data.Control) {
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

func collectionReplaceRecursive(ctx data.Context) (data.GetValue, data.Control) {
	return collectionMergeRecursive(ctx)
}

func collectionReverse(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	entries := toEntries(collectionItems(cv))
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionSearch(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	needle, _ := ctx.GetIndexValue(0)
	strict := false
	if s, ok := ctx.GetIndexValue(1); ok && s != nil && truthy(s) {
		strict = true
	}
	eq := valueLooseEqual
	if strict {
		eq = valueStrictEqual
	}
	for _, e := range toEntries(collectionItems(cv)) {
		if eq(e.value, needle) {
			return e.key, nil
		}
	}
	return data.NewBoolValue(false), nil
}

func collectionSelect(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	keys, _ := ctx.GetIndexValue(0)
	want := keysToStrings(keys)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		row := data.NewArrayValue(nil).(*data.ArrayValue)
		for _, k := range want {
			if v, ok := dataGetPath(e.value, k); ok {
				setEntry(row, k, v)
			}
		}
		out.List = append(out.List, data.NewZVal(row))
	}
	return newCollectionInstance(ctx, out)
}

func collectionShift(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	items := collectionItems(cv)
	if len(items.List) == 0 {
		return data.NewNullValue(), nil
	}
	first := items.List[0]
	items.List = items.List[1:]
	_ = cv.SetProperty("items", items)
	if first == nil {
		return data.NewNullValue(), nil
	}
	return first.Value, nil
}

func collectionShuffle(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if seed, ok := ctx.GetIndexValue(0); ok && seed != nil && !isNull(seed) {
		if iv, ok := seed.(*data.IntValue); ok {
			rand.Seed(int64(iv.Value))
		}
	}
	entries := toEntries(collectionItems(cv))
	rand.Shuffle(len(entries), func(i, j int) {
		entries[i], entries[j] = entries[j], entries[i]
	})
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range entries {
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionSkip(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	countArg, _ := ctx.GetIndexValue(0)
	count, ctl := intFromValue(countArg)
	if ctl != nil {
		return nil, ctl
	}
	return sliceCollection(ctx, cv, count, -1)
}

func collectionTake(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	limitArg, _ := ctx.GetIndexValue(0)
	limit, ctl := intFromValue(limitArg)
	if ctl != nil {
		return nil, ctl
	}
	return sliceCollection(ctx, cv, 0, limit)
}

func collectionSlice(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	offsetArg, _ := ctx.GetIndexValue(0)
	offset, ctl := intFromValue(offsetArg)
	if ctl != nil {
		return nil, ctl
	}
	length := -1
	if l, ok := ctx.GetIndexValue(1); ok && l != nil && !isNull(l) {
		length, ctl = intFromValue(l)
		if ctl != nil {
			return nil, ctl
		}
	}
	return sliceCollection(ctx, cv, offset, length)
}

func sliceCollection(ctx data.Context, cv *data.ClassValue, offset, length int) (data.GetValue, data.Control) {
	entries := toEntries(collectionItems(cv))
	if offset < 0 {
		offset = 0
	}
	if offset > len(entries) {
		entries = nil
	} else {
		entries = entries[offset:]
	}
	if length >= 0 && length < len(entries) {
		entries = entries[:length]
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range entries {
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionSkipWhile(ctx data.Context) (data.GetValue, data.Control) {
	return takeSkipWhile(ctx, false, true)
}

func collectionSkipUntil(ctx data.Context) (data.GetValue, data.Control) {
	return takeSkipWhile(ctx, false, false)
}

func collectionTakeWhile(ctx data.Context) (data.GetValue, data.Control) {
	return takeSkipWhile(ctx, true, true)
}

func collectionTakeUntil(ctx data.Context) (data.GetValue, data.Control) {
	return takeSkipWhile(ctx, true, false)
}

func takeSkipWhile(ctx data.Context, take bool, while bool) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	active := take
	for _, e := range toEntries(collectionItems(cv)) {
		ok, ctl := callBool(ctx, cb, e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		match := ok
		if !while {
			match = !ok
		}
		if take {
			if match {
				setEntry(out, e.keyStr, e.value)
			} else {
				break
			}
		} else {
			if match {
				break
			}
			setEntry(out, e.keyStr, e.value)
		}
	}
	if !take && !active {
		_ = active
	}
	return newCollectionInstance(ctx, out)
}

func collectionSliding(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	sizeArg, _ := ctx.GetIndexValue(0)
	size, ctl := intFromValue(sizeArg)
	if ctl != nil {
		return nil, ctl
	}
	step := 1
	if s, ok := ctx.GetIndexValue(1); ok && s != nil && !isNull(s) {
		step, ctl = intFromValue(s)
		if ctl != nil {
			return nil, ctl
		}
	}
	if size < 1 {
		size = 1
	}
	if step < 1 {
		step = 1
	}
	entries := toEntries(collectionItems(cv))
	chunks := data.NewArrayValue(nil).(*data.ArrayValue)
	for i := 0; i+size <= len(entries); i += step {
		cur := data.NewArrayValue(nil).(*data.ArrayValue)
		for j := 0; j < size; j++ {
			e := entries[i+j]
			setEntry(cur, e.keyStr, e.value)
		}
		inst, err := newCollectionInstance(ctx, cur)
		if err != nil {
			return nil, err
		}
		chunks.List = append(chunks.List, data.NewZVal(inst))
	}
	return newCollectionInstance(ctx, chunks)
}

func collectionSole(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	matches := whereMatches(ctx, cv)
	if len(matches) == 0 {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Item not found."), "Illuminate\\Support\\ItemNotFoundException")
	}
	if len(matches) > 1 {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Multiple items found."), "Illuminate\\Support\\MultipleItemsFoundException")
	}
	return matches[0].value, nil
}

func countWhereMatches(ctx data.Context, cv *data.ClassValue) int {
	return len(whereMatches(ctx, cv))
}

func whereMatches(ctx data.Context, cv *data.ClassValue) []kv {
	key, _ := ctx.GetIndexValue(0)
	op, _ := ctx.GetIndexValue(1)
	val, _ := ctx.GetIndexValue(2)
	var out []kv
	for _, e := range toEntries(collectionItems(cv)) {
		if val == nil || isNull(val) {
			if op == nil || isNull(op) {
				continue
			}
			itemVal, ok := dataGetPath(e.value, keyToString(key))
			if !ok {
				continue
			}
			if itemVal.AsString() == keyToString(op) {
				out = append(out, e)
			}
			continue
		}
		itemVal, ok := dataGetPath(e.value, keyToString(key))
		if !ok {
			continue
		}
		if compareOp(itemVal, keyToString(op), val) {
			out = append(out, e)
		}
	}
	return out
}

func collectionSortDesc(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	sorted, ctl := arrSortDesc(withArgs(ctx, collectionItems(cv)))
	if ctl != nil {
		return nil, ctl
	}
	return newCollectionInstance(ctx, sorted.(data.Value))
}

func collectionSortByDesc(ctx data.Context) (data.GetValue, data.Control) {
	return collectionSortByDir(ctx, true)
}

func collectionSortKeys(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	entries := toEntries(collectionItems(cv))
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].keyStr < entries[j].keyStr
	})
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range entries {
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionSortKeysDesc(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	entries := toEntries(collectionItems(cv))
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].keyStr > entries[j].keyStr
	})
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range entries {
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionSortKeysUsing(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	entries := toEntries(collectionItems(cv))
	sort.SliceStable(entries, func(i, j int) bool {
		cmp, ctl := callValue(ctx, cb, data.NewStringValue(entries[i].keyStr), data.NewStringValue(entries[j].keyStr))
		if ctl != nil {
			return false
		}
		if iv, ok := cmp.(*data.IntValue); ok {
			return iv.Value < 0
		}
		return entries[i].keyStr < entries[j].keyStr
	})
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range entries {
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionSplice(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	spliceOffsetArg, _ := ctx.GetIndexValue(0)
	offset, ctl := intFromValue(spliceOffsetArg)
	if ctl != nil {
		return nil, ctl
	}
	length := -1
	if l, ok := ctx.GetIndexValue(1); ok && l != nil && !isNull(l) {
		length, ctl = intFromValue(l)
		if ctl != nil {
			return nil, ctl
		}
	}
	replacement, _ := ctx.GetIndexValue(2)
	items := collectionItems(cv)
	entries := toEntries(items)
	if offset < 0 {
		offset = len(entries) + offset
	}
	if offset < 0 {
		offset = 0
	}
	end := len(entries)
	if length >= 0 {
		end = offset + length
		if end > len(entries) {
			end = len(entries)
		}
	}
	removed := data.NewArrayValue(nil).(*data.ArrayValue)
	for i := offset; i < end && i < len(entries); i++ {
		removed.List = append(removed.List, data.NewZVal(entries[i].value))
	}
	newEntries := append([]kv{}, entries[:offset]...)
	if replacement != nil && !isNull(replacement) {
		replItems, replCtl := getArrayableItems(ctx, replacement)
		if replCtl != nil {
			return nil, replCtl
		}
		for _, e := range toEntries(replItems) {
			newEntries = append(newEntries, e)
		}
	}
	newEntries = append(newEntries, entries[end:]...)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range newEntries {
		out.List = append(out.List, data.NewZVal(e.value))
	}
	_ = cv.SetProperty("items", out)
	removedInst, ctl := newCollectionInstance(ctx, removed)
	if ctl != nil {
		return nil, ctl
	}
	return removedInst, nil
}

func collectionSplit(ctx data.Context) (data.GetValue, data.Control) {
	return splitCollection(ctx, false)
}

func collectionSplitIn(ctx data.Context) (data.GetValue, data.Control) {
	return splitCollection(ctx, true)
}

func splitCollection(ctx data.Context, inGroups bool) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	splitArg, _ := ctx.GetIndexValue(0)
	n, ctl := intFromValue(splitArg)
	if ctl != nil {
		return nil, ctl
	}
	if n < 1 {
		n = 1
	}
	entries := toEntries(collectionItems(cv))
	chunks := data.NewArrayValue(nil).(*data.ArrayValue)
	if inGroups {
		size := (len(entries) + n - 1) / n
		if size < 1 {
			size = 1
		}
		for i := 0; i < len(entries); i += size {
			end := i + size
			if end > len(entries) {
				end = len(entries)
			}
			cur := data.NewArrayValue(nil).(*data.ArrayValue)
			for _, e := range entries[i:end] {
				setEntry(cur, e.keyStr, e.value)
			}
			inst, err := newCollectionInstance(ctx, cur)
			if err != nil {
				return nil, err
			}
			chunks.List = append(chunks.List, data.NewZVal(inst))
		}
	} else {
		for i := 0; i < n; i++ {
			cur := data.NewArrayValue(nil).(*data.ArrayValue)
			for j := i; j < len(entries); j += n {
				e := entries[j]
				setEntry(cur, e.keyStr, e.value)
			}
			inst, err := newCollectionInstance(ctx, cur)
			if err != nil {
				return nil, err
			}
			chunks.List = append(chunks.List, data.NewZVal(inst))
		}
	}
	return newCollectionInstance(ctx, chunks)
}

func collectionTransform(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb, _ := ctx.GetIndexValue(0)
	items := collectionItems(cv)
	for _, e := range toEntries(items) {
		ret, ctl := callValue(ctx, cb, e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		if ret != nil {
			if v, ok := ret.(data.Value); ok {
				setEntry(items, e.keyStr, v)
			}
		}
	}
	_ = cv.SetProperty("items", items)
	return cv, nil
}

func collectionUnion(ctx data.Context) (data.GetValue, data.Control) {
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
		if _, ok := dataGetPath(out, e.keyStr); ok {
			continue
		}
		setEntry(out, e.keyStr, e.value)
	}
	return newCollectionInstance(ctx, out)
}

func collectionUnshift(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	items := collectionItems(cv)
	for i := 0; ; i++ {
		v, ok := ctx.GetIndexValue(i)
		if !ok || v == nil || isNull(v) {
			break
		}
		items.List = append([]*data.ZVal{data.NewZVal(v)}, items.List...)
	}
	_ = cv.SetProperty("items", items)
	return cv, nil
}

func collectionZip(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	lists := []data.Value{collectionItems(cv)}
	for i := 0; ; i++ {
		v, ok := ctx.GetIndexValue(i)
		if !ok || v == nil || isNull(v) {
			break
		}
		av, err := getArrayableItems(ctx, v)
		if err != nil {
			return nil, err
		}
		lists = append(lists, av)
	}
	maxLen := 0
	for _, l := range lists {
		if n := len(toEntries(l)); n > maxLen {
			maxLen = n
		}
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for i := 0; i < maxLen; i++ {
		row := data.NewArrayValue(nil).(*data.ArrayValue)
		for _, l := range lists {
			entries := toEntries(l)
			if i < len(entries) {
				row.List = append(row.List, data.NewZVal(entries[i].value))
			} else {
				row.List = append(row.List, data.NewZVal(data.NewNullValue()))
			}
		}
		out.List = append(out.List, data.NewZVal(row))
	}
	return newCollectionInstance(ctx, out)
}

func intFromValue(v data.Value) (int, data.Control) {
	if v == nil || isNull(v) {
		return 0, nil
	}
	if iv, ok := v.(data.AsInt); ok {
		n, err := iv.AsInt()
		if err != nil {
			return 0, data.NewErrorThrow(nil, err)
		}
		return n, nil
	}
	return 0, nil
}

func collectionWhereInstanceOf(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	typeArg := kit.Arg(ctx, 0)
	var types []string
	if av, ok := kit.Unwrap(typeArg).(*data.ArrayValue); ok {
		for _, e := range toEntries(av) {
			types = append(types, e.value.AsString())
		}
	} else if typeArg != nil {
		types = []string{typeArg.AsString()}
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		ok := false
		for _, t := range types {
			if valueInstanceOf(e.value, t) {
				ok = true
				break
			}
		}
		if ok {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return newCollectionInstance(ctx, out)
}

func collectionWhereNull(ctx data.Context) (data.GetValue, data.Control) {
	return collectionWhereNullness(ctx, true)
}

func collectionWhereNotNull(ctx data.Context) (data.GetValue, data.Control) {
	return collectionWhereNullness(ctx, false)
}

func collectionWhereNullness(ctx data.Context, wantNull bool) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := ""
	if v := kit.Arg(ctx, 0); v != nil && !isNull(v) {
		key = v.AsString()
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		val := e.value
		if key != "" {
			if got, ok := dataGetPath(e.value, key); ok {
				val = got
			} else {
				val = data.NewNullValue()
			}
		}
		isN := val == nil || isNull(val)
		if isN == wantNull {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return newCollectionInstance(ctx, out)
}

func collectionWhereIn(ctx data.Context) (data.GetValue, data.Control) {
	return collectionWhereInOut(ctx, false)
}

func collectionWhereNotIn(ctx data.Context) (data.GetValue, data.Control) {
	return collectionWhereInOut(ctx, true)
}

func collectionWhereInOut(ctx data.Context, negate bool) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	key := ""
	if v := kit.Arg(ctx, 0); v != nil {
		key = v.AsString()
	}
	values, ctl := getArrayableItems(ctx, kit.Arg(ctx, 1))
	if ctl != nil {
		return nil, ctl
	}
	allowed := map[string]struct{}{}
	for _, e := range toEntries(values) {
		allowed[e.value.AsString()] = struct{}{}
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(collectionItems(cv)) {
		val := e.value
		if key != "" {
			if got, ok := dataGetPath(e.value, key); ok {
				val = got
			} else {
				continue
			}
		}
		_, in := allowed[val.AsString()]
		if in != negate {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return newCollectionInstance(ctx, out)
}

func valueInstanceOf(v data.Value, typeName string) bool {
	if typeName == "" {
		return false
	}
	cv, ok := kit.Unwrap(v).(*data.ClassValue)
	if !ok || cv == nil || cv.Class == nil {
		return false
	}
	return (data.Class{Name: typeName}).Is(cv)
}

// collectionReject 对齐 EnumeratesValues::reject。
func collectionReject(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb := kit.Arg(ctx, 0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	useCallable := cb != nil && !isNull(cb) && isCallableValue(cb)
	for _, e := range toEntries(collectionItems(cv)) {
		keep := true
		if useCallable {
			ok, err := callBool(ctx, cb, e.value, e.key)
			if err != nil {
				return nil, err
			}
			keep = !ok
		} else {
			// reject($value) — 去掉等于该值的元素
			target := cb
			if target == nil {
				target = data.NewBoolValue(true)
			}
			keep = !valuesEqual(e.value, target, false)
		}
		if keep {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return newCollectionInstance(ctx, out)
}

func collectionReduce(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb := kit.Arg(ctx, 0)
	carry := kit.Arg(ctx, 1)
	if carry == nil {
		carry = data.NewNullValue()
	}
	for _, e := range toEntries(collectionItems(cv)) {
		ret, err := callValue(ctx, cb, carry, e.value, e.key)
		if err != nil {
			return nil, err
		}
		if v, ok := ret.(data.Value); ok {
			carry = v
		} else {
			carry = data.NewNullValue()
		}
	}
	return carry, nil
}

func collectionEvery(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb := kit.Arg(ctx, 0)
	for _, e := range toEntries(collectionItems(cv)) {
		ok, err := callBool(ctx, cb, e.value, e.key)
		if err != nil {
			return nil, err
		}
		if !ok {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func collectionSome(ctx data.Context) (data.GetValue, data.Control) {
	// some(...) ≡ contains(...)
	return collectionContains(ctx)
}

func collectionAvg(ctx data.Context) (data.GetValue, data.Control) {
	return collectionNumericAggregate(ctx, "avg")
}

func collectionSum(ctx data.Context) (data.GetValue, data.Control) {
	return collectionNumericAggregate(ctx, "sum")
}

func collectionMax(ctx data.Context) (data.GetValue, data.Control) {
	return collectionNumericAggregate(ctx, "max")
}

func collectionMin(ctx data.Context) (data.GetValue, data.Control) {
	return collectionNumericAggregate(ctx, "min")
}

func collectionTap(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb := kit.Arg(ctx, 0)
	if cb != nil && !isNull(cb) {
		if _, err := callValue(ctx, cb, cv); err != nil {
			return nil, err
		}
	}
	return cv, nil
}

func collectionPipe(ctx data.Context) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	cb := kit.Arg(ctx, 0)
	if cb == nil || isNull(cb) {
		return cv, nil
	}
	return callValue(ctx, cb, cv)
}

func collectionWhenEmpty(ctx data.Context) (data.GetValue, data.Control) {
	return collectionWhenByEmpty(ctx, true)
}

func collectionWhenNotEmpty(ctx data.Context) (data.GetValue, data.Control) {
	return collectionWhenByEmpty(ctx, false)
}

func collectionUnlessEmpty(ctx data.Context) (data.GetValue, data.Control) {
	return collectionWhenByEmpty(ctx, false)
}

func collectionUnlessNotEmpty(ctx data.Context) (data.GetValue, data.Control) {
	return collectionWhenByEmpty(ctx, true)
}

func collectionWhenByEmpty(ctx data.Context, wantEmpty bool) (data.GetValue, data.Control) {
	cv, ctl := collectionReceiver(ctx)
	if ctl != nil {
		return nil, ctl
	}
	empty := len(toEntries(collectionItems(cv))) == 0
	run := empty == wantEmpty
	cb := kit.Arg(ctx, 0)
	def := kit.Arg(ctx, 1)
	if run && cb != nil && !isNull(cb) {
		ret, err := callValue(ctx, cb, cv)
		if err != nil {
			return nil, err
		}
		if v, ok := ret.(data.Value); ok && v != nil && !isNull(v) {
			return v, nil
		}
		return cv, nil
	}
	if !run && def != nil && !isNull(def) {
		ret, err := callValue(ctx, def, cv)
		if err != nil {
			return nil, err
		}
		if v, ok := ret.(data.Value); ok && v != nil && !isNull(v) {
			return v, nil
		}
	}
	return cv, nil
}

func collectionWhenProxy(ctx data.Context, target data.Value) (data.Value, data.Control) {
	return conditionable.NewWhenProxy(ctx, target)
}

func isCallableValue(v data.Value) bool {
	v = kit.Unwrap(v)
	switch v.(type) {
	case *data.FuncValue, *data.BoundFuncValue:
		return true
	}
	return false
}

func valuesEqual(a, b data.Value, strict bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if strict {
		return a.AsString() == b.AsString() && fmt.Sprintf("%T", a) == fmt.Sprintf("%T", b)
	}
	return a.AsString() == b.AsString()
}
