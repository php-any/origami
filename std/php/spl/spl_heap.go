package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	splHeapStorageKey = "__spl_heap__"
	splHeapPosKey     = "__spl_heap_pos__"
)

type heapCompareFunc func(a, b data.Value) int

func splHeapGetCV(ctx data.Context) *data.ClassValue {
	return aoGetClassValue(ctx)
}

func splHeapGetStorage(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.ObjectValue.GetProperty(splHeapStorageKey)
	if arr, ok := v.(*data.ArrayValue); ok {
		return arr
	}
	arr := &data.ArrayValue{List: []*data.ZVal{}}
	cv.ObjectValue.SetProperty(splHeapStorageKey, arr)
	return arr
}

func splHeapGetPos(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(splHeapPosKey)
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	return 0
}

func splHeapSetPos(cv *data.ClassValue, pos int) {
	cv.ObjectValue.SetProperty(splHeapPosKey, data.NewIntValue(pos))
}

func splHeapInitCV(cv *data.ClassValue) {
	cv.SetProperty(splHeapStorageKey, &data.ArrayValue{List: []*data.ZVal{}})
	cv.SetProperty(splHeapPosKey, data.NewIntValue(0))
}

func splHeapCompareForClass(cv *data.ClassValue, a, b data.Value) int {
	if cv == nil {
		return data.Compare(a, b)
	}
	method, ok := cv.GetMethod("compare")
	if !ok || method == nil {
		return data.Compare(a, b)
	}
	fnCtx := cv.CreateContext(method.GetVariables())
	if len(method.GetVariables()) >= 2 {
		fnCtx.SetVariableValue(method.GetVariables()[0], a)
		fnCtx.SetVariableValue(method.GetVariables()[1], b)
	}
	result, acl := method.Call(fnCtx)
	if acl != nil {
		return data.Compare(a, b)
	}
	if iv, ok := result.(data.AsInt); ok {
		cmp, _ := iv.AsInt()
		if cmp > 0 {
			return 1
		}
		if cmp < 0 {
			return -1
		}
		return 0
	}
	return data.Compare(a, b)
}

func splHeapBubbleUp(arr *data.ArrayValue, index int, cmp heapCompareFunc) {
	for index > 0 {
		parent := (index - 1) / 2
		if cmp(arr.List[index].Value, arr.List[parent].Value) >= 0 {
			break
		}
		arr.List[index], arr.List[parent] = arr.List[parent], arr.List[index]
		index = parent
	}
}

func splHeapBubbleDown(arr *data.ArrayValue, index int, cmp heapCompareFunc) {
	n := len(arr.List)
	for {
		smallest := index
		left := 2*index + 1
		right := 2*index + 2
		if left < n && cmp(arr.List[left].Value, arr.List[smallest].Value) < 0 {
			smallest = left
		}
		if right < n && cmp(arr.List[right].Value, arr.List[smallest].Value) < 0 {
			smallest = right
		}
		if smallest == index {
			break
		}
		arr.List[index], arr.List[smallest] = arr.List[smallest], arr.List[index]
		index = smallest
	}
}

func splHeapInsert(arr *data.ArrayValue, value data.Value, cmp heapCompareFunc) {
	arr.List = append(arr.List, data.NewZVal(value))
	splHeapBubbleUp(arr, len(arr.List)-1, cmp)
}

func splHeapExtractTop(arr *data.ArrayValue, cmp heapCompareFunc) data.Value {
	if len(arr.List) == 0 {
		return data.NewNullValue()
	}
	top := arr.List[0].Value
	last := len(arr.List) - 1
	arr.List[0] = arr.List[last]
	arr.List = arr.List[:last]
	if len(arr.List) > 0 {
		splHeapBubbleDown(arr, 0, cmp)
	}
	return top
}

func splHeapTop(arr *data.ArrayValue) data.Value {
	if len(arr.List) == 0 {
		return data.NewNullValue()
	}
	return arr.List[0].Value
}

func splHeapCompareFromCV(cv *data.ClassValue) heapCompareFunc {
	return func(a, b data.Value) int {
		return splHeapCompareForClass(cv, a, b)
	}
}

// SplHeapClass 实现 PHP SPL �?SplHeap（抽象类�?
type SplHeapClass struct {
	node.Node
}

func NewSplHeapClass() *SplHeapClass {
	return &SplHeapClass{}
}

func (c *SplHeapClass) IsBuiltinAbstractClass() bool { return true }

func (c *SplHeapClass) GetName() string    { return "SplHeap" }
func (c *SplHeapClass) GetExtend() *string { return nil }
func (c *SplHeapClass) GetImplements() []string {
	return []string{"Iterator", "Countable"}
}
func (c *SplHeapClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SplHeapClass) GetPropertyList() []data.Property              { return nil }
func (c *SplHeapClass) GetConstruct() data.Method                     { return nil }
func (c *SplHeapClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	splHeapInitCV(cv)
	return cv, nil
}

func (c *SplHeapClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "compare":
		return &SplHeapCompareMethod{}, true
	case "insert":
		return &SplHeapInsertMethod{}, true
	case "extract":
		return &SplHeapExtractMethod{}, true
	case "top":
		return &SplHeapTopMethod{}, true
	case "count":
		return &SplHeapCountMethod{}, true
	case "isEmpty":
		return &SplHeapIsEmptyMethod{}, true
	case "rewind":
		return &SplHeapRewindMethod{}, true
	case "valid":
		return &SplHeapValidMethod{}, true
	case "current":
		return &SplHeapCurrentMethod{}, true
	case "key":
		return &SplHeapKeyMethod{}, true
	case "next":
		return &SplHeapNextMethod{}, true
	}
	return nil, false
}

func (c *SplHeapClass) GetMethods() []data.Method {
	return []data.Method{
		&SplHeapCompareMethod{},
		&SplHeapInsertMethod{},
		&SplHeapExtractMethod{},
		&SplHeapTopMethod{},
		&SplHeapCountMethod{},
		&SplHeapIsEmptyMethod{},
		&SplHeapRewindMethod{},
		&SplHeapValidMethod{},
		&SplHeapCurrentMethod{},
		&SplHeapKeyMethod{},
		&SplHeapNextMethod{},
	}
}

type SplHeapCompareMethod struct{}

func (m *SplHeapCompareMethod) GetName() string            { return "compare" }
func (m *SplHeapCompareMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplHeapCompareMethod) GetIsStatic() bool          { return false }
func (m *SplHeapCompareMethod) GetReturnType() data.Types  { return data.Int{} }
func (m *SplHeapCompareMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value1", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "value2", 1, nil, data.Mixed{}),
	}
}
func (m *SplHeapCompareMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value1", 0, data.Mixed{}),
		node.NewVariable(nil, "value2", 1, data.Mixed{}),
	}
}
func (m *SplHeapCompareMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a, _ := ctx.GetIndexValue(0)
	b, _ := ctx.GetIndexValue(1)
	return data.NewIntValue(data.Compare(a, b)), nil
}

type SplHeapInsertMethod struct{}

func (m *SplHeapInsertMethod) GetName() string            { return "insert" }
func (m *SplHeapInsertMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplHeapInsertMethod) GetIsStatic() bool          { return false }
func (m *SplHeapInsertMethod) GetReturnType() data.Types  { return nil }
func (m *SplHeapInsertMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "value", 0, nil, data.Mixed{})}
}
func (m *SplHeapInsertMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "value", 0, data.Mixed{})}
}
func (m *SplHeapInsertMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splHeapGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	val, _ := ctx.GetIndexValue(0)
	cmp := splHeapCompareFromCV(cv)
	splHeapInsert(splHeapGetStorage(cv), val, cmp)
	return nil, nil
}

type SplHeapExtractMethod struct{}

func (m *SplHeapExtractMethod) GetName() string            { return "extract" }
func (m *SplHeapExtractMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplHeapExtractMethod) GetIsStatic() bool          { return false }
func (m *SplHeapExtractMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *SplHeapExtractMethod) GetParams() []data.GetValue { return nil }
func (m *SplHeapExtractMethod) GetVariables() []data.Variable {
	return nil
}
func (m *SplHeapExtractMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splHeapGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	cmp := splHeapCompareFromCV(cv)
	return splHeapExtractTop(splHeapGetStorage(cv), cmp), nil
}

type SplHeapTopMethod struct{}

func (m *SplHeapTopMethod) GetName() string            { return "top" }
func (m *SplHeapTopMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplHeapTopMethod) GetIsStatic() bool          { return false }
func (m *SplHeapTopMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *SplHeapTopMethod) GetParams() []data.GetValue { return nil }
func (m *SplHeapTopMethod) GetVariables() []data.Variable {
	return nil
}
func (m *SplHeapTopMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splHeapGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return splHeapTop(splHeapGetStorage(cv)), nil
}

type SplHeapCountMethod struct{}

func (m *SplHeapCountMethod) GetName() string               { return "count" }
func (m *SplHeapCountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplHeapCountMethod) GetIsStatic() bool             { return false }
func (m *SplHeapCountMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplHeapCountMethod) GetParams() []data.GetValue    { return nil }
func (m *SplHeapCountMethod) GetVariables() []data.Variable { return nil }
func (m *SplHeapCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splHeapGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(len(splHeapGetStorage(cv).List)), nil
}

type SplHeapIsEmptyMethod struct{}

func (m *SplHeapIsEmptyMethod) GetName() string            { return "isEmpty" }
func (m *SplHeapIsEmptyMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplHeapIsEmptyMethod) GetIsStatic() bool          { return false }
func (m *SplHeapIsEmptyMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *SplHeapIsEmptyMethod) GetParams() []data.GetValue { return nil }
func (m *SplHeapIsEmptyMethod) GetVariables() []data.Variable {
	return nil
}
func (m *SplHeapIsEmptyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splHeapGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(len(splHeapGetStorage(cv).List) == 0), nil
}

type SplHeapRewindMethod struct{}

func (m *SplHeapRewindMethod) GetName() string               { return "rewind" }
func (m *SplHeapRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplHeapRewindMethod) GetIsStatic() bool             { return false }
func (m *SplHeapRewindMethod) GetReturnType() data.Types     { return nil }
func (m *SplHeapRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *SplHeapRewindMethod) GetVariables() []data.Variable { return nil }
func (m *SplHeapRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv := splHeapGetCV(ctx); cv != nil {
		splHeapSetPos(cv, 0)
	}
	return nil, nil
}

type SplHeapValidMethod struct{}

func (m *SplHeapValidMethod) GetName() string               { return "valid" }
func (m *SplHeapValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplHeapValidMethod) GetIsStatic() bool             { return false }
func (m *SplHeapValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplHeapValidMethod) GetParams() []data.GetValue    { return nil }
func (m *SplHeapValidMethod) GetVariables() []data.Variable { return nil }
func (m *SplHeapValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splHeapGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	pos := splHeapGetPos(cv)
	arr := splHeapGetStorage(cv)
	return data.NewBoolValue(pos >= 0 && pos < len(arr.List)), nil
}

type SplHeapCurrentMethod struct{}

func (m *SplHeapCurrentMethod) GetName() string               { return "current" }
func (m *SplHeapCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplHeapCurrentMethod) GetIsStatic() bool             { return false }
func (m *SplHeapCurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *SplHeapCurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *SplHeapCurrentMethod) GetVariables() []data.Variable { return nil }
func (m *SplHeapCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splHeapGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	pos := splHeapGetPos(cv)
	arr := splHeapGetStorage(cv)
	if pos < 0 || pos >= len(arr.List) {
		return data.NewNullValue(), nil
	}
	return arr.List[pos].Value, nil
}

type SplHeapKeyMethod struct{}

func (m *SplHeapKeyMethod) GetName() string               { return "key" }
func (m *SplHeapKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplHeapKeyMethod) GetIsStatic() bool             { return false }
func (m *SplHeapKeyMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplHeapKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *SplHeapKeyMethod) GetVariables() []data.Variable { return nil }
func (m *SplHeapKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splHeapGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return data.NewIntValue(splHeapGetPos(cv)), nil
}

type SplHeapNextMethod struct{}

func (m *SplHeapNextMethod) GetName() string               { return "next" }
func (m *SplHeapNextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplHeapNextMethod) GetIsStatic() bool             { return false }
func (m *SplHeapNextMethod) GetReturnType() data.Types     { return nil }
func (m *SplHeapNextMethod) GetParams() []data.GetValue    { return nil }
func (m *SplHeapNextMethod) GetVariables() []data.Variable { return nil }
func (m *SplHeapNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv := splHeapGetCV(ctx); cv != nil {
		splHeapSetPos(cv, splHeapGetPos(cv)+1)
	}
	return nil, nil
}

// SplMinHeapClass 实现 PHP SPL �?SplMinHeap
type SplMinHeapClass struct {
	node.Node
}

func NewSplMinHeapClass() *SplMinHeapClass {
	return &SplMinHeapClass{}
}

func (c *SplMinHeapClass) GetName() string { return "SplMinHeap" }
func (c *SplMinHeapClass) GetExtend() *string {
	parent := "SplHeap"
	return &parent
}
func (c *SplMinHeapClass) GetImplements() []string {
	return []string{"Iterator", "Countable"}
}
func (c *SplMinHeapClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SplMinHeapClass) GetPropertyList() []data.Property              { return nil }
func (c *SplMinHeapClass) GetConstruct() data.Method                     { return nil }
func (c *SplMinHeapClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	splHeapInitCV(cv)
	return cv, nil
}

func (c *SplMinHeapClass) GetMethod(name string) (data.Method, bool) {
	return splExtendGetMethod(c, name, func(name string) (data.Method, bool) {
		if name == "compare" {
			return &SplMinHeapCompareMethod{}, true
		}
		return nil, false
	})
}

func (c *SplMinHeapClass) GetMethods() []data.Method {
	return splExtendGetMethods(c, []data.Method{&SplMinHeapCompareMethod{}})
}

type SplMinHeapCompareMethod struct{}

func (m *SplMinHeapCompareMethod) GetName() string            { return "compare" }
func (m *SplMinHeapCompareMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplMinHeapCompareMethod) GetIsStatic() bool          { return false }
func (m *SplMinHeapCompareMethod) GetReturnType() data.Types  { return data.Int{} }
func (m *SplMinHeapCompareMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value1", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "value2", 1, nil, data.Mixed{}),
	}
}
func (m *SplMinHeapCompareMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value1", 0, data.Mixed{}),
		node.NewVariable(nil, "value2", 1, data.Mixed{}),
	}
}
func (m *SplMinHeapCompareMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a, _ := ctx.GetIndexValue(0)
	b, _ := ctx.GetIndexValue(1)
	return data.NewIntValue(data.Compare(a, b)), nil
}

// SplMaxHeapClass 实现 PHP SPL �?SplMaxHeap
type SplMaxHeapClass struct {
	node.Node
}

func NewSplMaxHeapClass() *SplMaxHeapClass {
	return &SplMaxHeapClass{}
}

func (c *SplMaxHeapClass) GetName() string { return "SplMaxHeap" }
func (c *SplMaxHeapClass) GetExtend() *string {
	parent := "SplHeap"
	return &parent
}
func (c *SplMaxHeapClass) GetImplements() []string {
	return []string{"Iterator", "Countable"}
}
func (c *SplMaxHeapClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SplMaxHeapClass) GetPropertyList() []data.Property              { return nil }
func (c *SplMaxHeapClass) GetConstruct() data.Method                     { return nil }
func (c *SplMaxHeapClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	splHeapInitCV(cv)
	return cv, nil
}

func (c *SplMaxHeapClass) GetMethod(name string) (data.Method, bool) {
	return splExtendGetMethod(c, name, func(name string) (data.Method, bool) {
		if name == "compare" {
			return &SplMaxHeapCompareMethod{}, true
		}
		return nil, false
	})
}

func (c *SplMaxHeapClass) GetMethods() []data.Method {
	return splExtendGetMethods(c, []data.Method{&SplMaxHeapCompareMethod{}})
}

type SplMaxHeapCompareMethod struct{}

func (m *SplMaxHeapCompareMethod) GetName() string            { return "compare" }
func (m *SplMaxHeapCompareMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplMaxHeapCompareMethod) GetIsStatic() bool          { return false }
func (m *SplMaxHeapCompareMethod) GetReturnType() data.Types  { return data.Int{} }
func (m *SplMaxHeapCompareMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value1", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "value2", 1, nil, data.Mixed{}),
	}
}
func (m *SplMaxHeapCompareMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value1", 0, data.Mixed{}),
		node.NewVariable(nil, "value2", 1, data.Mixed{}),
	}
}
func (m *SplMaxHeapCompareMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a, _ := ctx.GetIndexValue(0)
	b, _ := ctx.GetIndexValue(1)
	cmp := data.Compare(a, b)
	return data.NewIntValue(-cmp), nil
}
