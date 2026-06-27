package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	ciCacheKey     = "__ci_cache__"
	ciFlagsKey     = "__ci_flags__"
	ciFullCacheKey = "__ci_full_cache__"
)

// CachingIteratorClass 实现 PHP CachingIterator
type CachingIteratorClass struct {
	node.Node
	StaticProperty map[string]data.Value
}

func NewCachingIteratorClass() *CachingIteratorClass {
	return &CachingIteratorClass{
		StaticProperty: map[string]data.Value{
			"CALL_TOSTRING":        data.NewIntValue(1),
			"CATCH_GET_CHILD":      data.NewIntValue(2),
			"TOSTRING_USE_KEY":     data.NewIntValue(4),
			"TOSTRING_USE_CURRENT": data.NewIntValue(8),
			"TOSTRING_USE_INNER":   data.NewIntValue(16),
			"FULL_CACHE":           data.NewIntValue(256),
			"MASK":                 data.NewIntValue(511),
		},
	}
}

func (c *CachingIteratorClass) GetName() string { return "CachingIterator" }
func (c *CachingIteratorClass) GetExtend() *string {
	parent := "IteratorIterator"
	return &parent
}
func (c *CachingIteratorClass) GetImplements() []string {
	return []string{"ArrayAccess", "Countable"}
}
func (c *CachingIteratorClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *CachingIteratorClass) GetPropertyList() []data.Property              { return nil }
func (c *CachingIteratorClass) GetConstruct() data.Method                     { return &CIConstructMethod{} }

func (c *CachingIteratorClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := c.StaticProperty[name]; ok {
		return v, true
	}
	return nil, false
}

func (c *CachingIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(iiValidKey, data.NewBoolValue(false))
	cv.SetProperty(iiCurValKey, data.NewNullValue())
	cv.SetProperty(iiCurKeyKey, data.NewNullValue())
	cv.SetProperty(ciCacheKey, data.NewNullValue())
	cv.SetProperty(ciFlagsKey, data.NewIntValue(0))
	cv.SetProperty(ciFullCacheKey, &data.ArrayValue{List: []*data.ZVal{}})
	return cv, nil
}

func (c *CachingIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &CIConstructMethod{}, true
	case "rewind":
		return &CIRewindMethod{}, true
	case "next":
		return &CINextMethod{}, true
	case "getCache":
		return &CIGetCacheMethod{}, true
	case "hasNext":
		return &CIHasNextMethod{}, true
	case "__toString":
		return &CIToStringMethod{}, true
	case "count":
		return &CICountMethod{}, true
	case "setFlags":
		return &CISetFlagsMethod{}, true
	case "getFlags":
		return &CIGetFlagsMethod{}, true
	}
	return nil, false
}

func (c *CachingIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&CIConstructMethod{}, &CIRewindMethod{}, &CINextMethod{},
		&CIGetCacheMethod{}, &CIHasNextMethod{}, &CIToStringMethod{},
		&CICountMethod{}, &CISetFlagsMethod{}, &CIGetFlagsMethod{},
	}
}

func ciGetCache(cv *data.ClassValue) data.Value {
	v, _ := cv.ObjectValue.GetProperty(ciCacheKey)
	if v == nil {
		return data.NewNullValue()
	}
	return v
}

func ciSetCache(cv *data.ClassValue, val data.Value) {
	if val == nil {
		val = data.NewNullValue()
	}
	cv.ObjectValue.SetProperty(ciCacheKey, val)
}

func ciGetFlags(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(ciFlagsKey)
	return splAsInt(v)
}

func ciGetFullCache(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.ObjectValue.GetProperty(ciFullCacheKey)
	if arr, ok := v.(*data.ArrayValue); ok {
		return arr
	}
	arr := &data.ArrayValue{List: []*data.ZVal{}}
	cv.ObjectValue.SetProperty(ciFullCacheKey, arr)
	return arr
}

func ciUpdateCache(cv *data.ClassValue) {
	inner := iiGetInner(cv)
	if !iiInnerValid(inner) {
		ciSetCache(cv, data.NewNullValue())
		return
	}
	cur := iiInnerCurrent(inner)
	ciSetCache(cv, cur)
	if ciGetFlags(cv)&256 != 0 { // FULL_CACHE
		arr := ciGetFullCache(cv)
		arr.List = append(arr.List, data.NewZVal(cur))
	}
}

func ciPeekHasNext(cv *data.ClassValue) bool {
	inner := iiGetInner(cv)
	if inner == nil {
		return false
	}
	// 保存位置：复制当�?inner 状态不可行，用 next+valid 探测后回退
	// 简化：若当�?valid，next 后看是否�?valid
	if !iiIsValid(cv) {
		return false
	}
	iiCallInnerMethod(inner, "next")
	has := iiInnerValid(inner)
	// 回退：rewind 并前进到原位置（代价高但测试场景数据量小�?
	iiCallInnerMethod(inner, "rewind")
	for iiInnerValid(inner) {
		cur := iiInnerCurrent(inner)
		key := iiInnerKeyVal(inner)
		if cur.AsString() == ciGetCache(cv).AsString() && key.AsString() == iiGetCurKey(cv).AsString() {
			break
		}
		iiCallInnerMethod(inner, "next")
	}
	return has
}

type CIConstructMethod struct{}

func (m *CIConstructMethod) GetName() string            { return "__construct" }
func (m *CIConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *CIConstructMethod) GetIsStatic() bool          { return false }
func (m *CIConstructMethod) GetReturnType() data.Types  { return nil }
func (m *CIConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Iterator")),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.NewBaseType("int")),
	}
}
func (m *CIConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("Iterator")),
		node.NewVariable(nil, "flags", 1, data.NewBaseType("int")),
	}
}
func (m *CIConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	it, _ := ctx.GetIndexValue(0)
	flags, _ := ctx.GetIndexValue(1)
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	iiSetInner(cv, it)
	cv.ObjectValue.SetProperty(ciFlagsKey, data.NewIntValue(splAsInt(flags)))
	cv.ObjectValue.SetProperty(ciFullCacheKey, &data.ArrayValue{List: []*data.ZVal{}})
	iiSetValid(cv, false)
	return nil, nil
}

type CIRewindMethod struct{}

func (m *CIRewindMethod) GetName() string               { return "rewind" }
func (m *CIRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *CIRewindMethod) GetIsStatic() bool             { return false }
func (m *CIRewindMethod) GetReturnType() data.Types     { return nil }
func (m *CIRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *CIRewindMethod) GetVariables() []data.Variable { return nil }
func (m *CIRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	cv.ObjectValue.SetProperty(ciFullCacheKey, &data.ArrayValue{List: []*data.ZVal{}})
	inner := iiGetInner(cv)
	iiCallInnerMethod(inner, "rewind")
	iiSyncFromInner(cv)
	ciUpdateCache(cv)
	return nil, nil
}

type CINextMethod struct{}

func (m *CINextMethod) GetName() string               { return "next" }
func (m *CINextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *CINextMethod) GetIsStatic() bool             { return false }
func (m *CINextMethod) GetReturnType() data.Types     { return nil }
func (m *CINextMethod) GetParams() []data.GetValue    { return nil }
func (m *CINextMethod) GetVariables() []data.Variable { return nil }
func (m *CINextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	inner := iiGetInner(cv)
	iiCallInnerMethod(inner, "next")
	iiSyncFromInner(cv)
	ciUpdateCache(cv)
	return nil, nil
}

type CIGetCacheMethod struct{}

func (m *CIGetCacheMethod) GetName() string               { return "getCache" }
func (m *CIGetCacheMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *CIGetCacheMethod) GetIsStatic() bool             { return false }
func (m *CIGetCacheMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *CIGetCacheMethod) GetParams() []data.GetValue    { return nil }
func (m *CIGetCacheMethod) GetVariables() []data.Variable { return nil }
func (m *CIGetCacheMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return ciGetCache(cv), nil
}

type CIHasNextMethod struct{}

func (m *CIHasNextMethod) GetName() string               { return "hasNext" }
func (m *CIHasNextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *CIHasNextMethod) GetIsStatic() bool             { return false }
func (m *CIHasNextMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *CIHasNextMethod) GetParams() []data.GetValue    { return nil }
func (m *CIHasNextMethod) GetVariables() []data.Variable { return nil }
func (m *CIHasNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(ciPeekHasNext(cv)), nil
}

type CIToStringMethod struct{}

func (m *CIToStringMethod) GetName() string               { return "__toString" }
func (m *CIToStringMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *CIToStringMethod) GetIsStatic() bool             { return false }
func (m *CIToStringMethod) GetReturnType() data.Types     { return data.String{} }
func (m *CIToStringMethod) GetParams() []data.GetValue    { return nil }
func (m *CIToStringMethod) GetVariables() []data.Variable { return nil }
func (m *CIToStringMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	flags := ciGetFlags(cv)
	if flags&4 != 0 {
		return iiGetCurKey(cv), nil
	}
	if flags&8 != 0 {
		return ciGetCache(cv), nil
	}
	return ciGetCache(cv), nil
}

type CICountMethod struct{}

func (m *CICountMethod) GetName() string               { return "count" }
func (m *CICountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *CICountMethod) GetIsStatic() bool             { return false }
func (m *CICountMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *CICountMethod) GetParams() []data.GetValue    { return nil }
func (m *CICountMethod) GetVariables() []data.Variable { return nil }
func (m *CICountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	if ciGetFlags(cv)&256 != 0 {
		return data.NewIntValue(len(ciGetFullCache(cv).List)), nil
	}
	return data.NewIntValue(0), nil
}

type CISetFlagsMethod struct{}

func (m *CISetFlagsMethod) GetName() string            { return "setFlags" }
func (m *CISetFlagsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *CISetFlagsMethod) GetIsStatic() bool          { return false }
func (m *CISetFlagsMethod) GetReturnType() data.Types  { return nil }
func (m *CISetFlagsMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "flags", 0, nil, data.NewBaseType("int")),
	}
}
func (m *CISetFlagsMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "flags", 0, data.NewBaseType("int")),
	}
}
func (m *CISetFlagsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	flags, _ := ctx.GetIndexValue(0)
	cv := splGetClassValue(ctx)
	if cv != nil {
		cv.ObjectValue.SetProperty(ciFlagsKey, data.NewIntValue(splAsInt(flags)))
	}
	return nil, nil
}

type CIGetFlagsMethod struct{}

func (m *CIGetFlagsMethod) GetName() string               { return "getFlags" }
func (m *CIGetFlagsMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *CIGetFlagsMethod) GetIsStatic() bool             { return false }
func (m *CIGetFlagsMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *CIGetFlagsMethod) GetParams() []data.GetValue    { return nil }
func (m *CIGetFlagsMethod) GetVariables() []data.Variable { return nil }
func (m *CIGetFlagsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(ciGetFlags(cv)), nil
}

// RecursiveCachingIteratorClass 实现 PHP RecursiveCachingIterator
type RecursiveCachingIteratorClass struct {
	node.Node
}

func NewRecursiveCachingIteratorClass() *RecursiveCachingIteratorClass {
	return &RecursiveCachingIteratorClass{}
}

func (c *RecursiveCachingIteratorClass) GetName() string { return "RecursiveCachingIterator" }
func (c *RecursiveCachingIteratorClass) GetExtend() *string {
	parent := "CachingIterator"
	return &parent
}
func (c *RecursiveCachingIteratorClass) GetImplements() []string {
	return []string{"RecursiveIterator"}
}
func (c *RecursiveCachingIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}
func (c *RecursiveCachingIteratorClass) GetPropertyList() []data.Property { return nil }
func (c *RecursiveCachingIteratorClass) GetConstruct() data.Method        { return &RCIConstructMethod{} }

func (c *RecursiveCachingIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(iiValidKey, data.NewBoolValue(false))
	cv.SetProperty(iiCurValKey, data.NewNullValue())
	cv.SetProperty(iiCurKeyKey, data.NewNullValue())
	cv.SetProperty(ciCacheKey, data.NewNullValue())
	cv.SetProperty(ciFlagsKey, data.NewIntValue(0))
	cv.SetProperty(ciFullCacheKey, &data.ArrayValue{List: []*data.ZVal{}})
	return cv, nil
}

func (c *RecursiveCachingIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &RCIConstructMethod{}, true
	case "hasChildren":
		return &RCIHasChildrenMethod{}, true
	case "getChildren":
		return &RCIGetChildrenMethod{}, true
	}
	return nil, false
}

func (c *RecursiveCachingIteratorClass) GetMethods() []data.Method {
	return []data.Method{&RCIConstructMethod{}, &RCIHasChildrenMethod{}, &RCIGetChildrenMethod{}}
}

type RCIConstructMethod struct{}

func (m *RCIConstructMethod) GetName() string            { return "__construct" }
func (m *RCIConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *RCIConstructMethod) GetIsStatic() bool          { return false }
func (m *RCIConstructMethod) GetReturnType() data.Types  { return nil }
func (m *RCIConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("RecursiveIterator")),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.NewBaseType("int")),
	}
}
func (m *RCIConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "iterator", 0, data.NewBaseType("RecursiveIterator")),
		node.NewVariable(nil, "flags", 1, data.NewBaseType("int")),
	}
}
func (m *RCIConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	it, _ := ctx.GetIndexValue(0)
	flags, _ := ctx.GetIndexValue(1)
	cv := splGetClassValue(ctx)
	if cv == nil {
		return nil, nil
	}
	iiSetInner(cv, it)
	cv.ObjectValue.SetProperty(ciFlagsKey, data.NewIntValue(splAsInt(flags)))
	cv.ObjectValue.SetProperty(ciFullCacheKey, &data.ArrayValue{List: []*data.ZVal{}})
	iiSetValid(cv, false)
	return nil, nil
}

type RCIHasChildrenMethod struct{}

func (m *RCIHasChildrenMethod) GetName() string               { return "hasChildren" }
func (m *RCIHasChildrenMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RCIHasChildrenMethod) GetIsStatic() bool             { return false }
func (m *RCIHasChildrenMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *RCIHasChildrenMethod) GetParams() []data.GetValue    { return nil }
func (m *RCIHasChildrenMethod) GetVariables() []data.Variable { return nil }
func (m *RCIHasChildrenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	inner := iiGetInner(cv)
	result, _ := filterCallInnerMethod(inner, "hasChildren")
	if bv, ok := result.(*data.BoolValue); ok {
		return bv, nil
	}
	return data.NewBoolValue(false), nil
}

type RCIGetChildrenMethod struct{}

func (m *RCIGetChildrenMethod) GetName() string               { return "getChildren" }
func (m *RCIGetChildrenMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RCIGetChildrenMethod) GetIsStatic() bool             { return false }
func (m *RCIGetChildrenMethod) GetReturnType() data.Types     { return nil }
func (m *RCIGetChildrenMethod) GetParams() []data.GetValue    { return nil }
func (m *RCIGetChildrenMethod) GetVariables() []data.Variable { return nil }
func (m *RCIGetChildrenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splGetClassValue(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	inner := iiGetInner(cv)
	result, _ := filterCallInnerMethod(inner, "getChildren")
	if result != nil {
		return result, nil
	}
	return data.NewNullValue(), nil
}
