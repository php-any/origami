package core

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RecursiveIteratorIterator 内部状态属性名（存储在 ClassValue.ObjectValue.property）
// 类似 FilterIterator 的模式，保证 PHP 子类继承时状态隔离
const (
	riiInnerKey  = "__rii_inner__"  // 根迭代器（*data.ClassValue 序列化为 Value）
	riiValidKey  = "__rii_valid__"  // bool
	riiCurValKey = "__rii_curval__" // current value
	riiCurKeyKey = "__rii_curkey__" // current key
	riiModeKey   = "__rii_mode__"   // int
	riiStackKey  = "__rii_stack__"  // *riiStackValue（伪装成 data.Value）
)

// riiStackEntry 栈帧
type riiStackEntry struct {
	iter *data.ClassValue
}

// riiStackValue 包装迭代器栈，实现 data.Value 接口
type riiStackValue struct {
	frames []riiStackEntry
}

func (s *riiStackValue) GetValue(ctx data.Context) (data.GetValue, data.Control) { return s, nil }
func (s *riiStackValue) AsString() string                                        { return "riiStack" }
func (s *riiStackValue) Marshal(serializer data.Serializer) ([]byte, error)      { return nil, nil }
func (s *riiStackValue) Unmarshal(b []byte, serializer data.Serializer) error    { return nil }
func (s *riiStackValue) ToGoValue(serializer data.Serializer) (any, error)       { return nil, nil }

// RecursiveIteratorIteratorClass 实现 PHP 的 RecursiveIteratorIterator 类
// 状态存储在 ClassValue.ObjectValue.property，支持 PHP 子类继承（skill: php-class-state-sharing-pattern）
type RecursiveIteratorIteratorClass struct {
	node.Node
	StaticProperty map[string]data.Value // 类常量（静态，不随实例变化）
}

func NewRecursiveIteratorIteratorClass() *RecursiveIteratorIteratorClass {
	return &RecursiveIteratorIteratorClass{
		StaticProperty: map[string]data.Value{
			"SELF_FIRST":       data.NewIntValue(1),
			"CHILD_FIRST":      data.NewIntValue(2),
			"LEAVES_ONLY":      data.NewIntValue(0),
			"SELF_FIRST_SELF":  data.NewIntValue(4),
			"CHILD_FIRST_SELF": data.NewIntValue(8),
		},
	}
}

func (r *RecursiveIteratorIteratorClass) GetName() string { return "RecursiveIteratorIterator" }

func (r *RecursiveIteratorIteratorClass) GetExtend() *string { return nil }

func (r *RecursiveIteratorIteratorClass) GetImplements() []string {
	return []string{"OuterIterator", "Iterator"}
}

func (r *RecursiveIteratorIteratorClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (r *RecursiveIteratorIteratorClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := r.StaticProperty[name]; ok {
		return v, true
	}
	return nil, false
}

func (r *RecursiveIteratorIteratorClass) GetPropertyList() []data.Property { return nil }

func (r *RecursiveIteratorIteratorClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(r, ctx.CreateBaseContext())
	// 初始化实例属性
	cv.SetProperty(riiValidKey, data.NewBoolValue(false))
	cv.SetProperty(riiCurValKey, data.NewNullValue())
	cv.SetProperty(riiCurKeyKey, data.NewNullValue())
	cv.SetProperty(riiModeKey, data.NewIntValue(0))
	cv.SetProperty(riiStackKey, &riiStackValue{})
	return cv, nil
}

func (r *RecursiveIteratorIteratorClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &RIIConstruct{}, true
	case "rewind":
		return &RIIRewind{}, true
	case "current":
		return &RIICurrent{}, true
	case "key":
		return &RIIKey{}, true
	case "next":
		return &RIINext{}, true
	case "valid":
		return &RIIValid{}, true
	case "getInnerIterator":
		return &RIIGetInnerIterator{}, true
	case "getDepth":
		return &RIIGetDepth{}, true
	case "getSubIterator":
		return &RIIGetSubIterator{}, true
	}
	return nil, false
}

func (r *RecursiveIteratorIteratorClass) GetMethods() []data.Method {
	return []data.Method{
		&RIIConstruct{},
		&RIIRewind{},
		&RIICurrent{},
		&RIIKey{},
		&RIINext{},
		&RIIValid{},
		&RIIGetInnerIterator{},
		&RIIGetDepth{},
		&RIIGetSubIterator{},
	}
}

func (r *RecursiveIteratorIteratorClass) GetConstruct() data.Method {
	return &RIIConstruct{}
}

// ---- 辅助：从 ctx 获取 ClassValue ----

func riiGetCV(ctx data.Context) *data.ClassValue {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		return cmc.ClassValue
	}
	if cv, ok := ctx.(*data.ClassValue); ok {
		return cv
	}
	return nil
}

// ---- 状态读写 ----

func riiGetInner(cv *data.ClassValue) *data.ClassValue {
	v, _ := cv.ObjectValue.GetProperty(riiInnerKey)
	if v == nil {
		return nil
	}
	if _, ok := v.(*data.NullValue); ok {
		return nil
	}
	if inner, ok := v.(*data.ClassValue); ok {
		return inner
	}
	return nil
}

func riiSetInner(cv *data.ClassValue, inner *data.ClassValue) {
	if inner == nil {
		cv.ObjectValue.SetProperty(riiInnerKey, data.NewNullValue())
	} else {
		cv.ObjectValue.SetProperty(riiInnerKey, inner)
	}
}

func riiIsValid(cv *data.ClassValue) bool {
	v, _ := cv.ObjectValue.GetProperty(riiValidKey)
	if bv, ok := v.(*data.BoolValue); ok {
		return bv.Value
	}
	return false
}

func riiSetValid(cv *data.ClassValue, valid bool) {
	cv.ObjectValue.SetProperty(riiValidKey, data.NewBoolValue(valid))
}

func riiGetCurVal(cv *data.ClassValue) data.Value {
	v, _ := cv.ObjectValue.GetProperty(riiCurValKey)
	if v == nil {
		return data.NewNullValue()
	}
	return v
}

func riiSetCurVal(cv *data.ClassValue, val data.Value) {
	if val == nil {
		val = data.NewNullValue()
	}
	cv.ObjectValue.SetProperty(riiCurValKey, val)
}

func riiGetCurKey(cv *data.ClassValue) data.Value {
	v, _ := cv.ObjectValue.GetProperty(riiCurKeyKey)
	if v == nil {
		return data.NewNullValue()
	}
	return v
}

func riiSetCurKey(cv *data.ClassValue, key data.Value) {
	if key == nil {
		key = data.NewNullValue()
	}
	cv.ObjectValue.SetProperty(riiCurKeyKey, key)
}

func riiGetMode(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(riiModeKey)
	if iv, ok := v.(*data.IntValue); ok {
		return int(iv.Value)
	}
	return 0
}

func riiSetMode(cv *data.ClassValue, mode int) {
	cv.ObjectValue.SetProperty(riiModeKey, data.NewIntValue(mode))
}

func riiGetStack(cv *data.ClassValue) *riiStackValue {
	v, _ := cv.ObjectValue.GetProperty(riiStackKey)
	if sv, ok := v.(*riiStackValue); ok {
		return sv
	}
	sv := &riiStackValue{}
	cv.ObjectValue.SetProperty(riiStackKey, sv)
	return sv
}

// ---- 迭代器方法调用辅助（复用内部迭代器的 ClassValue 上下文） ----

func riiCallMethod(iter *data.ClassValue, name string) (data.GetValue, data.Control) {
	if iter == nil {
		return nil, nil
	}
	if m, ok := iter.GetMethod(name); ok {
		fnCtx := iter.CreateContext(m.GetVariables())
		return m.Call(fnCtx)
	}
	return nil, nil
}

func riiCallBool(iter *data.ClassValue, name string) (bool, data.Control) {
	v, ctl := riiCallMethod(iter, name)
	if ctl != nil {
		return false, ctl
	}
	if bv, ok := v.(*data.BoolValue); ok {
		return bv.Value, nil
	}
	if av, ok := v.(data.AsBool); ok {
		result, err := av.AsBool()
		if err != nil {
			return false, nil
		}
		return result, nil
	}
	return v != nil, nil
}

func riiCallValue(iter *data.ClassValue, name string) (data.Value, data.Control) {
	v, ctl := riiCallMethod(iter, name)
	if ctl != nil {
		return nil, ctl
	}
	if val, ok := v.(data.Value); ok {
		return val, nil
	}
	return data.NewNullValue(), nil
}

// ---- currentIter：返回当前栈顶迭代器 ----

func riiCurrentIter(cv *data.ClassValue) *data.ClassValue {
	stack := riiGetStack(cv)
	if len(stack.frames) == 0 {
		return riiGetInner(cv)
	}
	return stack.frames[len(stack.frames)-1].iter
}

// ---- advance：从当前迭代器获取当前值更新缓存 ----

func riiAdvance(cv *data.ClassValue) data.Control {
	iter := riiCurrentIter(cv)
	if iter == nil {
		riiSetValid(cv, false)
		return nil
	}
	valid, ctl := riiCallBool(iter, "valid")
	if ctl != nil {
		return ctl
	}
	if !valid {
		riiSetValid(cv, false)
		return nil
	}
	val, ctl := riiCallValue(iter, "current")
	if ctl != nil {
		return ctl
	}
	key, ctl := riiCallValue(iter, "key")
	if ctl != nil {
		return ctl
	}
	riiSetCurVal(cv, val)
	riiSetCurKey(cv, key)
	riiSetValid(cv, true)
	return nil
}

// ---- __construct ----

type RIIConstruct struct{}

func (m *RIIConstruct) GetName() string            { return "__construct" }
func (m *RIIConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *RIIConstruct) GetIsStatic() bool          { return false }
func (m *RIIConstruct) GetReturnType() data.Types  { return nil }
func (m *RIIConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("iterator", 0, data.NewBaseType("Traversable")),
		data.NewVariable("mode", 1, data.NewBaseType("int")),
		data.NewVariable("flags", 2, data.NewBaseType("int")),
	}
}
func (m *RIIConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "iterator", 0, nil, data.NewBaseType("Traversable")),
		node.NewParameter(nil, "mode", 1, nil, data.NewBaseType("int")),
		node.NewParameter(nil, "flags", 2, nil, data.NewBaseType("int")),
	}
}
func (m *RIIConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return nil, data.NewErrorThrow(nil, errors.New("RecursiveIteratorIterator: 无法获取实例"))
	}

	iterVal, exists := ctx.GetIndexValue(0)
	if !exists {
		return nil, data.NewErrorThrow(nil, errors.New("缺少必需的迭代器参数"))
	}

	var classVal *data.ClassValue
	switch v := iterVal.(type) {
	case *data.ClassValue:
		classVal = v
	case *data.ThisValue:
		classVal = v.ClassValue
	default:
		return nil, data.NewErrorThrow(nil, errors.New("RecursiveIteratorIterator: 参数必须是对象"))
	}

	modeVal, hasMode := ctx.GetIndexValue(1)
	mode := 0
	if hasMode {
		if mi, ok := modeVal.(interface{ AsInt() int }); ok {
			mode = mi.AsInt()
		} else if mi64, ok := modeVal.(interface{ AsInt64() int64 }); ok {
			mode = int(mi64.AsInt64())
		}
	}

	riiSetInner(cv, classVal)
	riiSetMode(cv, mode)
	riiSetValid(cv, false)
	riiSetCurVal(cv, data.NewNullValue())
	riiSetCurKey(cv, data.NewNullValue())
	// 清空栈
	stack := riiGetStack(cv)
	stack.frames = stack.frames[:0]

	return nil, nil
}

// ---- rewind ----

type RIIRewind struct{}

func (m *RIIRewind) GetName() string               { return "rewind" }
func (m *RIIRewind) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RIIRewind) GetIsStatic() bool             { return false }
func (m *RIIRewind) GetVariables() []data.Variable { return nil }
func (m *RIIRewind) GetReturnType() data.Types     { return nil }
func (m *RIIRewind) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *RIIRewind) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}

	inner := riiGetInner(cv)
	if inner == nil {
		riiSetValid(cv, false)
		return nil, nil
	}

	// 清空栈
	stack := riiGetStack(cv)
	stack.frames = stack.frames[:0]

	// rewind 根迭代器
	if _, ctl := riiCallMethod(inner, "rewind"); ctl != nil {
		return nil, ctl
	}

	return nil, riiAdvance(cv)
}

// ---- next（深度优先 SELF_FIRST） ----

type RIINext struct{}

func (m *RIINext) GetName() string               { return "next" }
func (m *RIINext) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RIINext) GetIsStatic() bool             { return false }
func (m *RIINext) GetVariables() []data.Variable { return nil }
func (m *RIINext) GetReturnType() data.Types     { return nil }
func (m *RIINext) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *RIINext) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}

	inner := riiGetInner(cv)
	if inner == nil {
		riiSetValid(cv, false)
		return nil, nil
	}

	mode := riiGetMode(cv)
	stack := riiGetStack(cv)
	iter := riiCurrentIter(cv)

	// SELF_FIRST (1) 或 LEAVES_ONLY (0)：尝试进入子节点
	if mode == 1 || mode == 0 {
		hasChildren, ctl := riiCallBool(iter, "hasChildren")
		if ctl == nil && hasChildren {
			childVal, ctl2 := riiCallValue(iter, "getChildren")
			if ctl2 == nil {
				if childClass, ok := childVal.(*data.ClassValue); ok {
					if _, ctl3 := riiCallMethod(childClass, "rewind"); ctl3 == nil {
						childValid, ctl4 := riiCallBool(childClass, "valid")
						if ctl4 == nil && childValid {
							stack.frames = append(stack.frames, riiStackEntry{iter: childClass})
							return nil, riiAdvance(cv)
						}
					}
				}
			}
		}
	}

	// 无子节点，当前迭代器 next
	if _, ctl := riiCallMethod(iter, "next"); ctl != nil {
		return nil, ctl
	}

	valid, ctl := riiCallBool(iter, "valid")
	if ctl != nil {
		return nil, ctl
	}
	if valid {
		return nil, riiAdvance(cv)
	}

	// 当前迭代器耗尽，弹栈
	for len(stack.frames) > 0 {
		stack.frames = stack.frames[:len(stack.frames)-1]
		parentIter := riiCurrentIter(cv)

		if _, ctl := riiCallMethod(parentIter, "next"); ctl != nil {
			return nil, ctl
		}

		parentValid, ctl := riiCallBool(parentIter, "valid")
		if ctl != nil {
			return nil, ctl
		}
		if parentValid {
			return nil, riiAdvance(cv)
		}
	}

	riiSetValid(cv, false)
	riiSetCurVal(cv, data.NewNullValue())
	return nil, nil
}

// ---- current ----

type RIICurrent struct{}

func (m *RIICurrent) GetName() string               { return "current" }
func (m *RIICurrent) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RIICurrent) GetIsStatic() bool             { return false }
func (m *RIICurrent) GetVariables() []data.Variable { return nil }
func (m *RIICurrent) GetReturnType() data.Types     { return nil }
func (m *RIICurrent) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *RIICurrent) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return riiGetCurVal(cv), nil
}

// ---- key ----

type RIIKey struct{}

func (m *RIIKey) GetName() string               { return "key" }
func (m *RIIKey) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RIIKey) GetIsStatic() bool             { return false }
func (m *RIIKey) GetVariables() []data.Variable { return nil }
func (m *RIIKey) GetReturnType() data.Types     { return nil }
func (m *RIIKey) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *RIIKey) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewStringValue(""), nil
	}
	return riiGetCurKey(cv), nil
}

// ---- valid ----

type RIIValid struct{}

func (m *RIIValid) GetName() string               { return "valid" }
func (m *RIIValid) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RIIValid) GetIsStatic() bool             { return false }
func (m *RIIValid) GetVariables() []data.Variable { return nil }
func (m *RIIValid) GetReturnType() data.Types     { return nil }
func (m *RIIValid) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *RIIValid) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(riiIsValid(cv)), nil
}

// ---- getInnerIterator ----

type RIIGetInnerIterator struct{}

func (m *RIIGetInnerIterator) GetName() string               { return "getInnerIterator" }
func (m *RIIGetInnerIterator) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RIIGetInnerIterator) GetIsStatic() bool             { return false }
func (m *RIIGetInnerIterator) GetVariables() []data.Variable { return nil }
func (m *RIIGetInnerIterator) GetReturnType() data.Types     { return nil }
func (m *RIIGetInnerIterator) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *RIIGetInnerIterator) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	inner := riiGetInner(cv)
	if inner == nil {
		return data.NewNullValue(), nil
	}
	return inner, nil
}

// ---- getDepth ----

type RIIGetDepth struct{}

func (m *RIIGetDepth) GetName() string               { return "getDepth" }
func (m *RIIGetDepth) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RIIGetDepth) GetIsStatic() bool             { return false }
func (m *RIIGetDepth) GetVariables() []data.Variable { return nil }
func (m *RIIGetDepth) GetReturnType() data.Types     { return data.Int{} }
func (m *RIIGetDepth) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *RIIGetDepth) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	stack := riiGetStack(cv)
	return data.NewIntValue(len(stack.frames)), nil
}

// ---- getSubIterator ----

type RIIGetSubIterator struct{}

func (m *RIIGetSubIterator) GetName() string               { return "getSubIterator" }
func (m *RIIGetSubIterator) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *RIIGetSubIterator) GetIsStatic() bool             { return false }
func (m *RIIGetSubIterator) GetVariables() []data.Variable { return nil }
func (m *RIIGetSubIterator) GetReturnType() data.Types     { return nil }
func (m *RIIGetSubIterator) GetParams() []data.GetValue    { return []data.GetValue{} }
func (m *RIIGetSubIterator) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := riiGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	depthVal, hasDepth := ctx.GetIndexValue(0)
	stack := riiGetStack(cv)
	depth := len(stack.frames) // 默认当前深度
	if hasDepth {
		if di, ok := depthVal.(interface{ AsInt() int }); ok {
			depth = di.AsInt()
		} else if di64, ok := depthVal.(interface{ AsInt64() int64 }); ok {
			depth = int(di64.AsInt64())
		}
	}
	if depth == 0 {
		return riiGetInner(cv), nil
	}
	if depth > 0 && depth <= len(stack.frames) {
		return stack.frames[depth-1].iter, nil
	}
	return data.NewNullValue(), nil
}
