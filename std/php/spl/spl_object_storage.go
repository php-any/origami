package spl

import (
	"fmt"
	"reflect"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	sosStorageKey = "__sos_storage__"
	sosPosKey     = "__sos_pos__"
)

type sosEntry struct {
	object data.Value
	info   data.Value
	hash   string
}

type sosStorageValue struct {
	entries []sosEntry
	byPtr   map[int]int
}

func (s *sosStorageValue) GetValue(ctx data.Context) (data.GetValue, data.Control) { return s, nil }
func (s *sosStorageValue) AsString() string                                        { return "sosStorage" }
func (s *sosStorageValue) Marshal(serializer data.Serializer) ([]byte, error)      { return nil, nil }
func (s *sosStorageValue) Unmarshal(b []byte, serializer data.Serializer) error    { return nil }
func (s *sosStorageValue) ToGoValue(serializer data.Serializer) (any, error)       { return nil, nil }

func sosGetCV(ctx data.Context) *data.ClassValue {
	return aoGetClassValue(ctx)
}

func sosObjectPtrID(obj data.Value) (int, bool) {
	if obj == nil {
		return 0, false
	}
	objType := data.Object{}
	if !objType.Is(obj) {
		return 0, false
	}
	rv := reflect.ValueOf(obj)
	for rv.Kind() == reflect.Interface && !rv.IsNil() {
		rv = rv.Elem()
	}
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		return int(rv.Pointer()), true
	}
	return 0, false
}

func sosObjectHash(obj data.Value) string {
	if obj == nil {
		return ""
	}
	return fmt.Sprintf("%p", obj)
}

func sosGetStorage(cv *data.ClassValue) *sosStorageValue {
	v, _ := cv.ObjectValue.GetProperty(sosStorageKey)
	if storage, ok := v.(*sosStorageValue); ok {
		return storage
	}
	storage := &sosStorageValue{entries: []sosEntry{}, byPtr: map[int]int{}}
	cv.ObjectValue.SetProperty(sosStorageKey, storage)
	return storage
}

func sosGetPos(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(sosPosKey)
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	return 0
}

func sosSetPos(cv *data.ClassValue, pos int) {
	cv.ObjectValue.SetProperty(sosPosKey, data.NewIntValue(pos))
}

func sosInitCV(cv *data.ClassValue) {
	cv.SetProperty(sosStorageKey, &sosStorageValue{entries: []sosEntry{}, byPtr: map[int]int{}})
	cv.SetProperty(sosPosKey, data.NewIntValue(0))
}

func sosFindEntry(storage *sosStorageValue, obj data.Value) (int, bool) {
	ptr, ok := sosObjectPtrID(obj)
	if !ok {
		return -1, false
	}
	idx, ok := storage.byPtr[ptr]
	return idx, ok
}

func sosRebuildIndex(storage *sosStorageValue) {
	storage.byPtr = make(map[int]int, len(storage.entries))
	for i, e := range storage.entries {
		if ptr, ok := sosObjectPtrID(e.object); ok {
			storage.byPtr[ptr] = i
		}
	}
}

// SplObjectStorageClass 实现 PHP SPL �?SplObjectStorage
type SplObjectStorageClass struct {
	node.Node
}

func NewSplObjectStorageClass() *SplObjectStorageClass {
	return &SplObjectStorageClass{}
}

func (c *SplObjectStorageClass) GetName() string    { return "SplObjectStorage" }
func (c *SplObjectStorageClass) GetExtend() *string { return nil }
func (c *SplObjectStorageClass) GetImplements() []string {
	return []string{"Iterator", "Countable", "ArrayAccess"}
}
func (c *SplObjectStorageClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SplObjectStorageClass) GetPropertyList() []data.Property              { return nil }
func (c *SplObjectStorageClass) GetConstruct() data.Method                     { return nil }
func (c *SplObjectStorageClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	sosInitCV(cv)
	return cv, nil
}

func (c *SplObjectStorageClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "attach":
		return &SplObjectStorageAttachMethod{}, true
	case "detach":
		return &SplObjectStorageDetachMethod{}, true
	case "contains":
		return &SplObjectStorageContainsMethod{}, true
	case "count":
		return &SplObjectStorageCountMethod{}, true
	case "getHash":
		return &SplObjectStorageGetHashMethod{}, true
	case "rewind":
		return &SplObjectStorageRewindMethod{}, true
	case "valid":
		return &SplObjectStorageValidMethod{}, true
	case "current":
		return &SplObjectStorageCurrentMethod{}, true
	case "key":
		return &SplObjectStorageKeyMethod{}, true
	case "next":
		return &SplObjectStorageNextMethod{}, true
	case "offsetExists":
		return &SplObjectStorageOffsetExistsMethod{}, true
	case "offsetGet":
		return &SplObjectStorageOffsetGetMethod{}, true
	case "offsetSet":
		return &SplObjectStorageOffsetSetMethod{}, true
	case "offsetUnset":
		return &SplObjectStorageOffsetUnsetMethod{}, true
	}
	return nil, false
}

func (c *SplObjectStorageClass) GetMethods() []data.Method {
	return []data.Method{
		&SplObjectStorageAttachMethod{},
		&SplObjectStorageDetachMethod{},
		&SplObjectStorageContainsMethod{},
		&SplObjectStorageCountMethod{},
		&SplObjectStorageGetHashMethod{},
		&SplObjectStorageRewindMethod{},
		&SplObjectStorageValidMethod{},
		&SplObjectStorageCurrentMethod{},
		&SplObjectStorageKeyMethod{},
		&SplObjectStorageNextMethod{},
		&SplObjectStorageOffsetExistsMethod{},
		&SplObjectStorageOffsetGetMethod{},
		&SplObjectStorageOffsetSetMethod{},
		&SplObjectStorageOffsetUnsetMethod{},
	}
}

type SplObjectStorageAttachMethod struct{}

func (m *SplObjectStorageAttachMethod) GetName() string            { return "attach" }
func (m *SplObjectStorageAttachMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplObjectStorageAttachMethod) GetIsStatic() bool          { return false }
func (m *SplObjectStorageAttachMethod) GetReturnType() data.Types  { return nil }
func (m *SplObjectStorageAttachMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object")),
		node.NewParameter(nil, "info", 1, data.NewNullValue(), data.Mixed{}),
	}
}
func (m *SplObjectStorageAttachMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, data.NewBaseType("object")),
		node.NewVariable(nil, "info", 1, data.Mixed{}),
	}
}
func (m *SplObjectStorageAttachMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sosGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	obj, _ := ctx.GetIndexValue(0)
	info, _ := ctx.GetIndexValue(1)
	if info == nil {
		info = data.NewNullValue()
	}
	if _, ok := sosObjectPtrID(obj); !ok {
		return nil, nil
	}
	storage := sosGetStorage(cv)
	if idx, ok := sosFindEntry(storage, obj); ok {
		storage.entries[idx].info = info
		return nil, nil
	}
	hash := sosObjectHash(obj)
	storage.entries = append(storage.entries, sosEntry{object: obj, info: info, hash: hash})
	sosRebuildIndex(storage)
	return nil, nil
}

type SplObjectStorageDetachMethod struct{}

func (m *SplObjectStorageDetachMethod) GetName() string            { return "detach" }
func (m *SplObjectStorageDetachMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplObjectStorageDetachMethod) GetIsStatic() bool          { return false }
func (m *SplObjectStorageDetachMethod) GetReturnType() data.Types  { return nil }
func (m *SplObjectStorageDetachMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object"))}
}
func (m *SplObjectStorageDetachMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "object", 0, data.NewBaseType("object"))}
}
func (m *SplObjectStorageDetachMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sosGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	obj, _ := ctx.GetIndexValue(0)
	storage := sosGetStorage(cv)
	idx, ok := sosFindEntry(storage, obj)
	if !ok {
		return nil, nil
	}
	storage.entries = append(storage.entries[:idx], storage.entries[idx+1:]...)
	sosRebuildIndex(storage)
	pos := sosGetPos(cv)
	if pos > idx {
		sosSetPos(cv, pos-1)
	} else if pos >= len(storage.entries) && len(storage.entries) > 0 {
		sosSetPos(cv, len(storage.entries)-1)
	}
	return nil, nil
}

type SplObjectStorageContainsMethod struct{}

func (m *SplObjectStorageContainsMethod) GetName() string            { return "contains" }
func (m *SplObjectStorageContainsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplObjectStorageContainsMethod) GetIsStatic() bool          { return false }
func (m *SplObjectStorageContainsMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *SplObjectStorageContainsMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object"))}
}
func (m *SplObjectStorageContainsMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "object", 0, data.NewBaseType("object"))}
}
func (m *SplObjectStorageContainsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sosGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	obj, _ := ctx.GetIndexValue(0)
	_, ok := sosFindEntry(sosGetStorage(cv), obj)
	return data.NewBoolValue(ok), nil
}

type SplObjectStorageCountMethod struct{}

func (m *SplObjectStorageCountMethod) GetName() string               { return "count" }
func (m *SplObjectStorageCountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplObjectStorageCountMethod) GetIsStatic() bool             { return false }
func (m *SplObjectStorageCountMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplObjectStorageCountMethod) GetParams() []data.GetValue    { return nil }
func (m *SplObjectStorageCountMethod) GetVariables() []data.Variable { return nil }
func (m *SplObjectStorageCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sosGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(len(sosGetStorage(cv).entries)), nil
}

type SplObjectStorageGetHashMethod struct{}

func (m *SplObjectStorageGetHashMethod) GetName() string            { return "getHash" }
func (m *SplObjectStorageGetHashMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplObjectStorageGetHashMethod) GetIsStatic() bool          { return false }
func (m *SplObjectStorageGetHashMethod) GetReturnType() data.Types  { return data.String{} }
func (m *SplObjectStorageGetHashMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object"))}
}
func (m *SplObjectStorageGetHashMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "object", 0, data.NewBaseType("object"))}
}
func (m *SplObjectStorageGetHashMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	obj, _ := ctx.GetIndexValue(0)
	return data.NewStringValue(sosObjectHash(obj)), nil
}

type SplObjectStorageRewindMethod struct{}

func (m *SplObjectStorageRewindMethod) GetName() string               { return "rewind" }
func (m *SplObjectStorageRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplObjectStorageRewindMethod) GetIsStatic() bool             { return false }
func (m *SplObjectStorageRewindMethod) GetReturnType() data.Types     { return nil }
func (m *SplObjectStorageRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *SplObjectStorageRewindMethod) GetVariables() []data.Variable { return nil }
func (m *SplObjectStorageRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv := sosGetCV(ctx); cv != nil {
		sosSetPos(cv, 0)
	}
	return nil, nil
}

type SplObjectStorageValidMethod struct{}

func (m *SplObjectStorageValidMethod) GetName() string               { return "valid" }
func (m *SplObjectStorageValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplObjectStorageValidMethod) GetIsStatic() bool             { return false }
func (m *SplObjectStorageValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplObjectStorageValidMethod) GetParams() []data.GetValue    { return nil }
func (m *SplObjectStorageValidMethod) GetVariables() []data.Variable { return nil }
func (m *SplObjectStorageValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sosGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	pos := sosGetPos(cv)
	return data.NewBoolValue(pos >= 0 && pos < len(sosGetStorage(cv).entries)), nil
}

type SplObjectStorageCurrentMethod struct{}

func (m *SplObjectStorageCurrentMethod) GetName() string               { return "current" }
func (m *SplObjectStorageCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplObjectStorageCurrentMethod) GetIsStatic() bool             { return false }
func (m *SplObjectStorageCurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *SplObjectStorageCurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *SplObjectStorageCurrentMethod) GetVariables() []data.Variable { return nil }
func (m *SplObjectStorageCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sosGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	pos := sosGetPos(cv)
	entries := sosGetStorage(cv).entries
	if pos < 0 || pos >= len(entries) {
		return data.NewNullValue(), nil
	}
	return entries[pos].object, nil
}

type SplObjectStorageKeyMethod struct{}

func (m *SplObjectStorageKeyMethod) GetName() string               { return "key" }
func (m *SplObjectStorageKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplObjectStorageKeyMethod) GetIsStatic() bool             { return false }
func (m *SplObjectStorageKeyMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *SplObjectStorageKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *SplObjectStorageKeyMethod) GetVariables() []data.Variable { return nil }
func (m *SplObjectStorageKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sosGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	pos := sosGetPos(cv)
	entries := sosGetStorage(cv).entries
	if pos < 0 || pos >= len(entries) {
		return data.NewNullValue(), nil
	}
	return data.NewStringValue(entries[pos].hash), nil
}

type SplObjectStorageNextMethod struct{}

func (m *SplObjectStorageNextMethod) GetName() string               { return "next" }
func (m *SplObjectStorageNextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplObjectStorageNextMethod) GetIsStatic() bool             { return false }
func (m *SplObjectStorageNextMethod) GetReturnType() data.Types     { return nil }
func (m *SplObjectStorageNextMethod) GetParams() []data.GetValue    { return nil }
func (m *SplObjectStorageNextMethod) GetVariables() []data.Variable { return nil }
func (m *SplObjectStorageNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv := sosGetCV(ctx); cv != nil {
		sosSetPos(cv, sosGetPos(cv)+1)
	}
	return nil, nil
}

type SplObjectStorageOffsetExistsMethod struct{}

func (m *SplObjectStorageOffsetExistsMethod) GetName() string            { return "offsetExists" }
func (m *SplObjectStorageOffsetExistsMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplObjectStorageOffsetExistsMethod) GetIsStatic() bool          { return false }
func (m *SplObjectStorageOffsetExistsMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *SplObjectStorageOffsetExistsMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object"))}
}
func (m *SplObjectStorageOffsetExistsMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "object", 0, data.NewBaseType("object"))}
}
func (m *SplObjectStorageOffsetExistsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return (&SplObjectStorageContainsMethod{}).Call(ctx)
}

type SplObjectStorageOffsetGetMethod struct{}

func (m *SplObjectStorageOffsetGetMethod) GetName() string            { return "offsetGet" }
func (m *SplObjectStorageOffsetGetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplObjectStorageOffsetGetMethod) GetIsStatic() bool          { return false }
func (m *SplObjectStorageOffsetGetMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *SplObjectStorageOffsetGetMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object"))}
}
func (m *SplObjectStorageOffsetGetMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "object", 0, data.NewBaseType("object"))}
}
func (m *SplObjectStorageOffsetGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sosGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	obj, _ := ctx.GetIndexValue(0)
	idx, ok := sosFindEntry(sosGetStorage(cv), obj)
	if !ok {
		return data.NewNullValue(), nil
	}
	return sosGetStorage(cv).entries[idx].info, nil
}

type SplObjectStorageOffsetSetMethod struct{}

func (m *SplObjectStorageOffsetSetMethod) GetName() string            { return "offsetSet" }
func (m *SplObjectStorageOffsetSetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplObjectStorageOffsetSetMethod) GetIsStatic() bool          { return false }
func (m *SplObjectStorageOffsetSetMethod) GetReturnType() data.Types  { return nil }
func (m *SplObjectStorageOffsetSetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object")),
		node.NewParameter(nil, "info", 1, data.NewNullValue(), data.Mixed{}),
	}
}
func (m *SplObjectStorageOffsetSetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "object", 0, data.NewBaseType("object")),
		node.NewVariable(nil, "info", 1, data.Mixed{}),
	}
}
func (m *SplObjectStorageOffsetSetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return (&SplObjectStorageAttachMethod{}).Call(ctx)
}

type SplObjectStorageOffsetUnsetMethod struct{}

func (m *SplObjectStorageOffsetUnsetMethod) GetName() string            { return "offsetUnset" }
func (m *SplObjectStorageOffsetUnsetMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplObjectStorageOffsetUnsetMethod) GetIsStatic() bool          { return false }
func (m *SplObjectStorageOffsetUnsetMethod) GetReturnType() data.Types  { return nil }
func (m *SplObjectStorageOffsetUnsetMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "object", 0, nil, data.NewBaseType("object"))}
}
func (m *SplObjectStorageOffsetUnsetMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "object", 0, data.NewBaseType("object"))}
}
func (m *SplObjectStorageOffsetUnsetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return (&SplObjectStorageDetachMethod{}).Call(ctx)
}
