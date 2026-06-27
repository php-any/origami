package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	splListStorageKey = "__spl_list__"
	splListPosKey     = "__spl_pos__"
	splListModeKey    = "__spl_iter_mode__"

	SplITModeKeep   = 0
	SplITModeDelete = 1
	SplITModeFIFO   = 0
	SplITModeLIFO   = 2
)

func splListGetCV(ctx data.Context) *data.ClassValue {
	return aoGetClassValue(ctx)
}

func splListInitCV(cv *data.ClassValue) {
	cv.SetProperty(splListStorageKey, &data.ArrayValue{List: []*data.ZVal{}})
	cv.SetProperty(splListPosKey, data.NewIntValue(0))
	cv.SetProperty(splListModeKey, data.NewIntValue(SplITModeFIFO))
}

func splListGetMode(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(splListModeKey)
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	return SplITModeFIFO
}

func splListIsLIFO(cv *data.ClassValue) bool {
	return splListGetMode(cv)&SplITModeLIFO != 0
}

func splListIsDelete(cv *data.ClassValue) bool {
	return splListGetMode(cv)&SplITModeDelete != 0
}

func splListGetStorage(cv *data.ClassValue) *data.ArrayValue {
	v, _ := cv.ObjectValue.GetProperty(splListStorageKey)
	if arr, ok := v.(*data.ArrayValue); ok {
		return arr
	}
	arr := &data.ArrayValue{List: []*data.ZVal{}}
	cv.ObjectValue.SetProperty(splListStorageKey, arr)
	return arr
}

func splListGetPos(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(splListPosKey)
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	return 0
}

func splListSetPos(cv *data.ClassValue, pos int) {
	cv.ObjectValue.SetProperty(splListPosKey, data.NewIntValue(pos))
}

func splListRewind(cv *data.ClassValue) {
	arr := splListGetStorage(cv)
	if splListIsLIFO(cv) && len(arr.List) > 0 {
		splListSetPos(cv, len(arr.List)-1)
		return
	}
	splListSetPos(cv, 0)
}

func splListValid(cv *data.ClassValue) bool {
	arr := splListGetStorage(cv)
	pos := splListGetPos(cv)
	return pos >= 0 && pos < len(arr.List)
}

func splListCurrent(cv *data.ClassValue) data.Value {
	arr := splListGetStorage(cv)
	pos := splListGetPos(cv)
	if pos < 0 || pos >= len(arr.List) {
		return data.NewNullValue()
	}
	return arr.List[pos].Value
}

func splListKey(cv *data.ClassValue) data.Value {
	return data.NewIntValue(splListGetPos(cv))
}

func splListNext(cv *data.ClassValue) {
	arr := splListGetStorage(cv)
	pos := splListGetPos(cv)
	if splListIsDelete(cv) && pos >= 0 && pos < len(arr.List) {
		splListRemoveAt(arr, pos)
		if splListIsLIFO(cv) {
			if pos >= len(arr.List) {
				pos = len(arr.List) - 1
			}
			splListSetPos(cv, pos)
			return
		}
		splListSetPos(cv, pos)
		return
	}
	if splListIsLIFO(cv) {
		splListSetPos(cv, pos-1)
		return
	}
	splListSetPos(cv, pos+1)
}

func splListOffsetIndex(offset data.Value) (int, bool) {
	if iv, ok := offset.(data.AsInt); ok {
		i, err := iv.AsInt()
		if err == nil {
			return i, true
		}
	}
	return 0, false
}

func splListRemoveAt(arr *data.ArrayValue, index int) {
	if index < 0 || index >= len(arr.List) {
		return
	}
	arr.List = append(arr.List[:index], arr.List[index+1:]...)
}

// SplDoublyLinkedListClass 实现 PHP SPL �?SplDoublyLinkedList
type SplDoublyLinkedListClass struct {
	node.Node
	StaticProperty map[string]data.Value
}

func NewSplDoublyLinkedListClass() *SplDoublyLinkedListClass {
	return &SplDoublyLinkedListClass{
		StaticProperty: map[string]data.Value{
			"IT_MODE_FIFO":   data.NewIntValue(SplITModeFIFO),
			"IT_MODE_LIFO":   data.NewIntValue(SplITModeLIFO),
			"IT_MODE_KEEP":   data.NewIntValue(SplITModeKeep),
			"IT_MODE_DELETE": data.NewIntValue(SplITModeDelete),
		},
	}
}

func (c *SplDoublyLinkedListClass) GetName() string { return "SplDoublyLinkedList" }
func (c *SplDoublyLinkedListClass) GetExtend() *string {
	return nil
}
func (c *SplDoublyLinkedListClass) GetImplements() []string {
	return []string{"Iterator", "Countable", "ArrayAccess"}
}
func (c *SplDoublyLinkedListClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SplDoublyLinkedListClass) GetPropertyList() []data.Property              { return nil }
func (c *SplDoublyLinkedListClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := c.StaticProperty[name]; ok {
		return v, true
	}
	return nil, false
}
func (c *SplDoublyLinkedListClass) GetConstruct() data.Method { return nil }
func (c *SplDoublyLinkedListClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	splListInitCV(cv)
	return cv, nil
}

func (c *SplDoublyLinkedListClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "push":
		return &SplDLLPushMethod{}, true
	case "pop":
		return &SplDLLPopMethod{}, true
	case "shift":
		return &SplDLLShiftMethod{}, true
	case "unshift":
		return &SplDLLUnshiftMethod{}, true
	case "top":
		return &SplDLLTopMethod{}, true
	case "bottom":
		return &SplDLLBottomMethod{}, true
	case "count":
		return &SplDLLCountMethod{}, true
	case "isEmpty":
		return &SplDLLIsEmptyMethod{}, true
	case "rewind":
		return &SplDLLRewindMethod{}, true
	case "valid":
		return &SplDLLValidMethod{}, true
	case "current":
		return &SplDLLCurrentMethod{}, true
	case "key":
		return &SplDLLKeyMethod{}, true
	case "next":
		return &SplDLLNextMethod{}, true
	case "offsetExists":
		return &SplDLLOffsetExistsMethod{}, true
	case "offsetGet":
		return &SplDLLOffsetGetMethod{}, true
	case "offsetSet":
		return &SplDLLOffsetSetMethod{}, true
	case "offsetUnset":
		return &SplDLLOffsetUnsetMethod{}, true
	case "setIteratorMode":
		return &SplDLLSetIteratorModeMethod{}, true
	}
	return nil, false
}

func (c *SplDoublyLinkedListClass) GetMethods() []data.Method {
	return []data.Method{
		&SplDLLPushMethod{},
		&SplDLLPopMethod{},
		&SplDLLShiftMethod{},
		&SplDLLUnshiftMethod{},
		&SplDLLTopMethod{},
		&SplDLLBottomMethod{},
		&SplDLLCountMethod{},
		&SplDLLIsEmptyMethod{},
		&SplDLLRewindMethod{},
		&SplDLLValidMethod{},
		&SplDLLCurrentMethod{},
		&SplDLLKeyMethod{},
		&SplDLLNextMethod{},
		&SplDLLOffsetExistsMethod{},
		&SplDLLOffsetGetMethod{},
		&SplDLLOffsetSetMethod{},
		&SplDLLOffsetUnsetMethod{},
		&SplDLLSetIteratorModeMethod{},
	}
}

type SplDLLPushMethod struct{}

func (m *SplDLLPushMethod) GetName() string            { return "push" }
func (m *SplDLLPushMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplDLLPushMethod) GetIsStatic() bool          { return false }
func (m *SplDLLPushMethod) GetReturnType() data.Types  { return data.Int{} }
func (m *SplDLLPushMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "value", 0, nil, data.Mixed{})}
}
func (m *SplDLLPushMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "value", 0, data.Mixed{})}
}
func (m *SplDLLPushMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	val, _ := ctx.GetIndexValue(0)
	arr := splListGetStorage(cv)
	arr.List = append(arr.List, data.NewZVal(val))
	return data.NewIntValue(len(arr.List)), nil
}

type SplDLLPopMethod struct{}

func (m *SplDLLPopMethod) GetName() string            { return "pop" }
func (m *SplDLLPopMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplDLLPopMethod) GetIsStatic() bool          { return false }
func (m *SplDLLPopMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *SplDLLPopMethod) GetParams() []data.GetValue { return nil }
func (m *SplDLLPopMethod) GetVariables() []data.Variable {
	return nil
}
func (m *SplDLLPopMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	arr := splListGetStorage(cv)
	if len(arr.List) == 0 {
		return data.NewNullValue(), nil
	}
	last := arr.List[len(arr.List)-1].Value
	arr.List = arr.List[:len(arr.List)-1]
	return last, nil
}

type SplDLLShiftMethod struct{}

func (m *SplDLLShiftMethod) GetName() string            { return "shift" }
func (m *SplDLLShiftMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplDLLShiftMethod) GetIsStatic() bool          { return false }
func (m *SplDLLShiftMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *SplDLLShiftMethod) GetParams() []data.GetValue { return nil }
func (m *SplDLLShiftMethod) GetVariables() []data.Variable {
	return nil
}
func (m *SplDLLShiftMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	arr := splListGetStorage(cv)
	if len(arr.List) == 0 {
		return data.NewNullValue(), nil
	}
	first := arr.List[0].Value
	arr.List = arr.List[1:]
	pos := splListGetPos(cv)
	if pos > 0 {
		splListSetPos(cv, pos-1)
	}
	return first, nil
}

type SplDLLUnshiftMethod struct{}

func (m *SplDLLUnshiftMethod) GetName() string            { return "unshift" }
func (m *SplDLLUnshiftMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplDLLUnshiftMethod) GetIsStatic() bool          { return false }
func (m *SplDLLUnshiftMethod) GetReturnType() data.Types  { return data.Int{} }
func (m *SplDLLUnshiftMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "value", 0, nil, data.Mixed{})}
}
func (m *SplDLLUnshiftMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "value", 0, data.Mixed{})}
}
func (m *SplDLLUnshiftMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	val, _ := ctx.GetIndexValue(0)
	arr := splListGetStorage(cv)
	arr.List = append([]*data.ZVal{data.NewZVal(val)}, arr.List...)
	splListSetPos(cv, splListGetPos(cv)+1)
	return data.NewIntValue(len(arr.List)), nil
}

type SplDLLTopMethod struct{}

func (m *SplDLLTopMethod) GetName() string            { return "top" }
func (m *SplDLLTopMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplDLLTopMethod) GetIsStatic() bool          { return false }
func (m *SplDLLTopMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *SplDLLTopMethod) GetParams() []data.GetValue { return nil }
func (m *SplDLLTopMethod) GetVariables() []data.Variable {
	return nil
}
func (m *SplDLLTopMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	arr := splListGetStorage(cv)
	if len(arr.List) == 0 {
		return data.NewNullValue(), nil
	}
	return arr.List[len(arr.List)-1].Value, nil
}

type SplDLLBottomMethod struct{}

func (m *SplDLLBottomMethod) GetName() string            { return "bottom" }
func (m *SplDLLBottomMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplDLLBottomMethod) GetIsStatic() bool          { return false }
func (m *SplDLLBottomMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *SplDLLBottomMethod) GetParams() []data.GetValue { return nil }
func (m *SplDLLBottomMethod) GetVariables() []data.Variable {
	return nil
}
func (m *SplDLLBottomMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	arr := splListGetStorage(cv)
	if len(arr.List) == 0 {
		return data.NewNullValue(), nil
	}
	return arr.List[0].Value, nil
}

type SplDLLCountMethod struct{}

func (m *SplDLLCountMethod) GetName() string               { return "count" }
func (m *SplDLLCountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplDLLCountMethod) GetIsStatic() bool             { return false }
func (m *SplDLLCountMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplDLLCountMethod) GetParams() []data.GetValue    { return nil }
func (m *SplDLLCountMethod) GetVariables() []data.Variable { return nil }
func (m *SplDLLCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(len(splListGetStorage(cv).List)), nil
}

type SplDLLIsEmptyMethod struct{}

func (m *SplDLLIsEmptyMethod) GetName() string            { return "isEmpty" }
func (m *SplDLLIsEmptyMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplDLLIsEmptyMethod) GetIsStatic() bool          { return false }
func (m *SplDLLIsEmptyMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *SplDLLIsEmptyMethod) GetParams() []data.GetValue { return nil }
func (m *SplDLLIsEmptyMethod) GetVariables() []data.Variable {
	return nil
}
func (m *SplDLLIsEmptyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(len(splListGetStorage(cv).List) == 0), nil
}

type SplDLLRewindMethod struct{}

func (m *SplDLLRewindMethod) GetName() string               { return "rewind" }
func (m *SplDLLRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplDLLRewindMethod) GetIsStatic() bool             { return false }
func (m *SplDLLRewindMethod) GetReturnType() data.Types     { return nil }
func (m *SplDLLRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *SplDLLRewindMethod) GetVariables() []data.Variable { return nil }
func (m *SplDLLRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv := splListGetCV(ctx); cv != nil {
		splListRewind(cv)
	}
	return nil, nil
}

type SplDLLValidMethod struct{}

func (m *SplDLLValidMethod) GetName() string               { return "valid" }
func (m *SplDLLValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplDLLValidMethod) GetIsStatic() bool             { return false }
func (m *SplDLLValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplDLLValidMethod) GetParams() []data.GetValue    { return nil }
func (m *SplDLLValidMethod) GetVariables() []data.Variable { return nil }
func (m *SplDLLValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(splListValid(cv)), nil
}

type SplDLLCurrentMethod struct{}

func (m *SplDLLCurrentMethod) GetName() string               { return "current" }
func (m *SplDLLCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplDLLCurrentMethod) GetIsStatic() bool             { return false }
func (m *SplDLLCurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *SplDLLCurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *SplDLLCurrentMethod) GetVariables() []data.Variable { return nil }
func (m *SplDLLCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return splListCurrent(cv), nil
}

type SplDLLKeyMethod struct{}

func (m *SplDLLKeyMethod) GetName() string               { return "key" }
func (m *SplDLLKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplDLLKeyMethod) GetIsStatic() bool             { return false }
func (m *SplDLLKeyMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplDLLKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *SplDLLKeyMethod) GetVariables() []data.Variable { return nil }
func (m *SplDLLKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return splListKey(cv), nil
}

type SplDLLNextMethod struct{}

func (m *SplDLLNextMethod) GetName() string               { return "next" }
func (m *SplDLLNextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplDLLNextMethod) GetIsStatic() bool             { return false }
func (m *SplDLLNextMethod) GetReturnType() data.Types     { return nil }
func (m *SplDLLNextMethod) GetParams() []data.GetValue    { return nil }
func (m *SplDLLNextMethod) GetVariables() []data.Variable { return nil }
func (m *SplDLLNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv := splListGetCV(ctx); cv != nil {
		splListNext(cv)
	}
	return nil, nil
}

type SplDLLOffsetExistsMethod struct{}

func (m *SplDLLOffsetExistsMethod) GetName() string            { return "offsetExists" }
func (m *SplDLLOffsetExistsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplDLLOffsetExistsMethod) GetIsStatic() bool          { return false }
func (m *SplDLLOffsetExistsMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *SplDLLOffsetExistsMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "index", 0, nil, data.Mixed{})}
}
func (m *SplDLLOffsetExistsMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "index", 0, data.Mixed{})}
}
func (m *SplDLLOffsetExistsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	offset, _ := ctx.GetIndexValue(0)
	i, ok := splListOffsetIndex(offset)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	arr := splListGetStorage(cv)
	return data.NewBoolValue(i >= 0 && i < len(arr.List)), nil
}

type SplDLLOffsetGetMethod struct{}

func (m *SplDLLOffsetGetMethod) GetName() string            { return "offsetGet" }
func (m *SplDLLOffsetGetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplDLLOffsetGetMethod) GetIsStatic() bool          { return false }
func (m *SplDLLOffsetGetMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *SplDLLOffsetGetMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "index", 0, nil, data.Mixed{})}
}
func (m *SplDLLOffsetGetMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "index", 0, data.Mixed{})}
}
func (m *SplDLLOffsetGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	offset, _ := ctx.GetIndexValue(0)
	i, ok := splListOffsetIndex(offset)
	if !ok {
		return data.NewNullValue(), nil
	}
	arr := splListGetStorage(cv)
	if i < 0 || i >= len(arr.List) {
		return data.NewNullValue(), nil
	}
	return arr.List[i].Value, nil
}

type SplDLLOffsetSetMethod struct{}

func (m *SplDLLOffsetSetMethod) GetName() string            { return "offsetSet" }
func (m *SplDLLOffsetSetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplDLLOffsetSetMethod) GetIsStatic() bool          { return false }
func (m *SplDLLOffsetSetMethod) GetReturnType() data.Types  { return nil }
func (m *SplDLLOffsetSetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "index", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "newval", 1, nil, data.Mixed{}),
	}
}
func (m *SplDLLOffsetSetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "index", 0, data.Mixed{}),
		node.NewVariable(nil, "newval", 1, data.Mixed{}),
	}
}
func (m *SplDLLOffsetSetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	offset, _ := ctx.GetIndexValue(0)
	val, _ := ctx.GetIndexValue(1)
	arr := splListGetStorage(cv)
	if offset == nil {
		arr.List = append(arr.List, data.NewZVal(val))
		return nil, nil
	}
	i, ok := splListOffsetIndex(offset)
	if !ok || i < 0 || i >= len(arr.List) {
		return nil, nil
	}
	arr.List[i].Value = val
	return nil, nil
}

type SplDLLOffsetUnsetMethod struct{}

func (m *SplDLLOffsetUnsetMethod) GetName() string            { return "offsetUnset" }
func (m *SplDLLOffsetUnsetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplDLLOffsetUnsetMethod) GetIsStatic() bool          { return false }
func (m *SplDLLOffsetUnsetMethod) GetReturnType() data.Types  { return nil }
func (m *SplDLLOffsetUnsetMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "index", 0, nil, data.Mixed{})}
}
func (m *SplDLLOffsetUnsetMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "index", 0, data.Mixed{})}
}
func (m *SplDLLOffsetUnsetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	offset, _ := ctx.GetIndexValue(0)
	i, ok := splListOffsetIndex(offset)
	if !ok {
		return nil, nil
	}
	arr := splListGetStorage(cv)
	splListRemoveAt(arr, i)
	pos := splListGetPos(cv)
	if pos > i {
		splListSetPos(cv, pos-1)
	}
	return nil, nil
}

type SplDLLSetIteratorModeMethod struct{}

func (m *SplDLLSetIteratorModeMethod) GetName() string            { return "setIteratorMode" }
func (m *SplDLLSetIteratorModeMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplDLLSetIteratorModeMethod) GetIsStatic() bool          { return false }
func (m *SplDLLSetIteratorModeMethod) GetReturnType() data.Types  { return nil }
func (m *SplDLLSetIteratorModeMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "mode", 0, data.NewIntValue(0), data.Int{})}
}
func (m *SplDLLSetIteratorModeMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "mode", 0, data.Int{})}
}
func (m *SplDLLSetIteratorModeMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	mode, _ := ctx.GetIndexValue(0)
	if mode != nil {
		cv.ObjectValue.SetProperty(splListModeKey, mode)
	}
	return nil, nil
}
