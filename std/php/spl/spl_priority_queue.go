package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	spqStorageKey      = "__spq_storage__"
	spqExtractFlagsKey = "__spq_extract_flags__"
	spqPosKey          = "__spq_pos__"

	SpqExtrData     = 0x00000001
	SpqExtrPriority = 0x00000002
	SpqExtrBoth     = 0x00000003
)

type spqEntry struct {
	value    data.Value
	priority data.Value
}

type spqEntriesValue struct {
	entries []spqEntry
}

func (s *spqEntriesValue) GetValue(ctx data.Context) (data.GetValue, data.Control) { return s, nil }
func (s *spqEntriesValue) AsString() string                                        { return "spqEntries" }
func (s *spqEntriesValue) Marshal(serializer data.Serializer) ([]byte, error)      { return nil, nil }
func (s *spqEntriesValue) Unmarshal(b []byte, serializer data.Serializer) error    { return nil }
func (s *spqEntriesValue) ToGoValue(serializer data.Serializer) (any, error)       { return nil, nil }

func spqGetCV(ctx data.Context) *data.ClassValue {
	return aoGetClassValue(ctx)
}

func spqGetEntries(cv *data.ClassValue) *spqEntriesValue {
	v, _ := cv.ObjectValue.GetProperty(spqStorageKey)
	if entries, ok := v.(*spqEntriesValue); ok {
		return entries
	}
	entries := &spqEntriesValue{entries: []spqEntry{}}
	cv.ObjectValue.SetProperty(spqStorageKey, entries)
	return entries
}

func spqGetExtractFlags(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(spqExtractFlagsKey)
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	return SpqExtrData
}

func spqGetPos(cv *data.ClassValue) int {
	v, _ := cv.ObjectValue.GetProperty(spqPosKey)
	if iv, ok := v.(*data.IntValue); ok {
		return iv.Value
	}
	return 0
}

func spqSetPos(cv *data.ClassValue, pos int) {
	cv.ObjectValue.SetProperty(spqPosKey, data.NewIntValue(pos))
}

// spqComparePriority 实现最大堆：优先级数值越大越先出队（�?PHP SplPriorityQueue 一致）
func spqComparePriority(a, b data.Value) int {
	return data.Compare(a, b)
}

func spqBubbleUp(entries []spqEntry, index int) {
	for index > 0 {
		parent := (index - 1) / 2
		if spqComparePriority(entries[index].priority, entries[parent].priority) <= 0 {
			break
		}
		entries[index], entries[parent] = entries[parent], entries[index]
		index = parent
	}
}

func spqBubbleDown(entries []spqEntry, index int) {
	n := len(entries)
	for {
		largest := index
		left := 2*index + 1
		right := 2*index + 2
		if left < n && spqComparePriority(entries[left].priority, entries[largest].priority) > 0 {
			largest = left
		}
		if right < n && spqComparePriority(entries[right].priority, entries[largest].priority) > 0 {
			largest = right
		}
		if largest == index {
			break
		}
		entries[index], entries[largest] = entries[largest], entries[index]
		index = largest
	}
}

func spqInsert(entries *spqEntriesValue, value, priority data.Value) {
	entries.entries = append(entries.entries, spqEntry{value: value, priority: priority})
	spqBubbleUp(entries.entries, len(entries.entries)-1)
}

func spqExtractTop(entries *spqEntriesValue) spqEntry {
	if len(entries.entries) == 0 {
		return spqEntry{value: data.NewNullValue(), priority: data.NewNullValue()}
	}
	top := entries.entries[0]
	last := len(entries.entries) - 1
	entries.entries[0] = entries.entries[last]
	entries.entries = entries.entries[:last]
	if len(entries.entries) > 0 {
		spqBubbleDown(entries.entries, 0)
	}
	return top
}

func spqFormatExtract(entry spqEntry, flags int) data.Value {
	switch flags {
	case SpqExtrPriority:
		return entry.priority
	case SpqExtrBoth:
		return &data.ArrayValue{
			List: []*data.ZVal{
				data.NewNamedZVal("data", entry.value),
				data.NewNamedZVal("priority", entry.priority),
			},
		}
	default:
		return entry.value
	}
}

// SplPriorityQueueClass 实现 PHP SPL �?SplPriorityQueue
type SplPriorityQueueClass struct {
	node.Node
	StaticProperty map[string]data.Value
}

func NewSplPriorityQueueClass() *SplPriorityQueueClass {
	return &SplPriorityQueueClass{
		StaticProperty: map[string]data.Value{
			"EXTR_DATA":     data.NewIntValue(SpqExtrData),
			"EXTR_PRIORITY": data.NewIntValue(SpqExtrPriority),
			"EXTR_BOTH":     data.NewIntValue(SpqExtrBoth),
		},
	}
}

func (c *SplPriorityQueueClass) GetName() string    { return "SplPriorityQueue" }
func (c *SplPriorityQueueClass) GetExtend() *string { return nil }
func (c *SplPriorityQueueClass) GetImplements() []string {
	return []string{"Iterator", "Countable"}
}
func (c *SplPriorityQueueClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SplPriorityQueueClass) GetPropertyList() []data.Property              { return nil }
func (c *SplPriorityQueueClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := c.StaticProperty[name]; ok {
		return v, true
	}
	return nil, false
}
func (c *SplPriorityQueueClass) GetConstruct() data.Method { return nil }
func (c *SplPriorityQueueClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(spqStorageKey, &spqEntriesValue{entries: []spqEntry{}})
	cv.SetProperty(spqExtractFlagsKey, data.NewIntValue(SpqExtrData))
	cv.SetProperty(spqPosKey, data.NewIntValue(0))
	return cv, nil
}

func (c *SplPriorityQueueClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "insert":
		return &SplPriorityQueueInsertMethod{}, true
	case "extract":
		return &SplPriorityQueueExtractMethod{}, true
	case "top":
		return &SplPriorityQueueTopMethod{}, true
	case "count":
		return &SplPriorityQueueCountMethod{}, true
	case "isEmpty":
		return &SplPriorityQueueIsEmptyMethod{}, true
	case "setExtractFlags":
		return &SplPriorityQueueSetExtractFlagsMethod{}, true
	case "getExtractFlags":
		return &SplPriorityQueueGetExtractFlagsMethod{}, true
	case "rewind":
		return &SplPriorityQueueRewindMethod{}, true
	case "valid":
		return &SplPriorityQueueValidMethod{}, true
	case "current":
		return &SplPriorityQueueCurrentMethod{}, true
	case "key":
		return &SplPriorityQueueKeyMethod{}, true
	case "next":
		return &SplPriorityQueueNextMethod{}, true
	}
	return nil, false
}

func (c *SplPriorityQueueClass) GetMethods() []data.Method {
	return []data.Method{
		&SplPriorityQueueInsertMethod{},
		&SplPriorityQueueExtractMethod{},
		&SplPriorityQueueTopMethod{},
		&SplPriorityQueueCountMethod{},
		&SplPriorityQueueIsEmptyMethod{},
		&SplPriorityQueueSetExtractFlagsMethod{},
		&SplPriorityQueueGetExtractFlagsMethod{},
		&SplPriorityQueueRewindMethod{},
		&SplPriorityQueueValidMethod{},
		&SplPriorityQueueCurrentMethod{},
		&SplPriorityQueueKeyMethod{},
		&SplPriorityQueueNextMethod{},
	}
}

type SplPriorityQueueInsertMethod struct{}

func (m *SplPriorityQueueInsertMethod) GetName() string            { return "insert" }
func (m *SplPriorityQueueInsertMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplPriorityQueueInsertMethod) GetIsStatic() bool          { return false }
func (m *SplPriorityQueueInsertMethod) GetReturnType() data.Types  { return nil }
func (m *SplPriorityQueueInsertMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "priority", 1, data.NewIntValue(0), data.Mixed{}),
	}
}
func (m *SplPriorityQueueInsertMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.Mixed{}),
		node.NewVariable(nil, "priority", 1, data.Mixed{}),
	}
}
func (m *SplPriorityQueueInsertMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := spqGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	val, _ := ctx.GetIndexValue(0)
	priority, _ := ctx.GetIndexValue(1)
	if priority == nil {
		priority = data.NewIntValue(0)
	}
	spqInsert(spqGetEntries(cv), val, priority)
	return nil, nil
}

type SplPriorityQueueExtractMethod struct{}

func (m *SplPriorityQueueExtractMethod) GetName() string            { return "extract" }
func (m *SplPriorityQueueExtractMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplPriorityQueueExtractMethod) GetIsStatic() bool          { return false }
func (m *SplPriorityQueueExtractMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *SplPriorityQueueExtractMethod) GetParams() []data.GetValue { return nil }
func (m *SplPriorityQueueExtractMethod) GetVariables() []data.Variable {
	return nil
}
func (m *SplPriorityQueueExtractMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := spqGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	entry := spqExtractTop(spqGetEntries(cv))
	if entry.value == nil {
		return data.NewNullValue(), nil
	}
	return spqFormatExtract(entry, spqGetExtractFlags(cv)), nil
}

type SplPriorityQueueTopMethod struct{}

func (m *SplPriorityQueueTopMethod) GetName() string            { return "top" }
func (m *SplPriorityQueueTopMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplPriorityQueueTopMethod) GetIsStatic() bool          { return false }
func (m *SplPriorityQueueTopMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *SplPriorityQueueTopMethod) GetParams() []data.GetValue { return nil }
func (m *SplPriorityQueueTopMethod) GetVariables() []data.Variable {
	return nil
}
func (m *SplPriorityQueueTopMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := spqGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	entries := spqGetEntries(cv)
	if len(entries.entries) == 0 {
		return data.NewNullValue(), nil
	}
	return spqFormatExtract(entries.entries[0], spqGetExtractFlags(cv)), nil
}

type SplPriorityQueueCountMethod struct{}

func (m *SplPriorityQueueCountMethod) GetName() string               { return "count" }
func (m *SplPriorityQueueCountMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplPriorityQueueCountMethod) GetIsStatic() bool             { return false }
func (m *SplPriorityQueueCountMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplPriorityQueueCountMethod) GetParams() []data.GetValue    { return nil }
func (m *SplPriorityQueueCountMethod) GetVariables() []data.Variable { return nil }
func (m *SplPriorityQueueCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := spqGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(0), nil
	}
	return data.NewIntValue(len(spqGetEntries(cv).entries)), nil
}

type SplPriorityQueueIsEmptyMethod struct{}

func (m *SplPriorityQueueIsEmptyMethod) GetName() string            { return "isEmpty" }
func (m *SplPriorityQueueIsEmptyMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplPriorityQueueIsEmptyMethod) GetIsStatic() bool          { return false }
func (m *SplPriorityQueueIsEmptyMethod) GetReturnType() data.Types  { return data.Bool{} }
func (m *SplPriorityQueueIsEmptyMethod) GetParams() []data.GetValue { return nil }
func (m *SplPriorityQueueIsEmptyMethod) GetVariables() []data.Variable {
	return nil
}
func (m *SplPriorityQueueIsEmptyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := spqGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(true), nil
	}
	return data.NewBoolValue(len(spqGetEntries(cv).entries) == 0), nil
}

type SplPriorityQueueSetExtractFlagsMethod struct{}

func (m *SplPriorityQueueSetExtractFlagsMethod) GetName() string { return "setExtractFlags" }
func (m *SplPriorityQueueSetExtractFlagsMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *SplPriorityQueueSetExtractFlagsMethod) GetIsStatic() bool         { return false }
func (m *SplPriorityQueueSetExtractFlagsMethod) GetReturnType() data.Types { return nil }
func (m *SplPriorityQueueSetExtractFlagsMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "flags", 0, data.NewIntValue(SpqExtrData), data.Int{})}
}
func (m *SplPriorityQueueSetExtractFlagsMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "flags", 0, data.Int{})}
}
func (m *SplPriorityQueueSetExtractFlagsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := spqGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	flags, _ := ctx.GetIndexValue(0)
	if iv, ok := flags.(data.AsInt); ok {
		v, _ := iv.AsInt()
		cv.SetProperty(spqExtractFlagsKey, data.NewIntValue(v))
	}
	return nil, nil
}

type SplPriorityQueueGetExtractFlagsMethod struct{}

func (m *SplPriorityQueueGetExtractFlagsMethod) GetName() string { return "getExtractFlags" }
func (m *SplPriorityQueueGetExtractFlagsMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}
func (m *SplPriorityQueueGetExtractFlagsMethod) GetIsStatic() bool          { return false }
func (m *SplPriorityQueueGetExtractFlagsMethod) GetReturnType() data.Types  { return data.Int{} }
func (m *SplPriorityQueueGetExtractFlagsMethod) GetParams() []data.GetValue { return nil }
func (m *SplPriorityQueueGetExtractFlagsMethod) GetVariables() []data.Variable {
	return nil
}
func (m *SplPriorityQueueGetExtractFlagsMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := spqGetCV(ctx)
	if cv == nil {
		return data.NewIntValue(SpqExtrData), nil
	}
	return data.NewIntValue(spqGetExtractFlags(cv)), nil
}

type SplPriorityQueueRewindMethod struct{}

func (m *SplPriorityQueueRewindMethod) GetName() string               { return "rewind" }
func (m *SplPriorityQueueRewindMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplPriorityQueueRewindMethod) GetIsStatic() bool             { return false }
func (m *SplPriorityQueueRewindMethod) GetReturnType() data.Types     { return nil }
func (m *SplPriorityQueueRewindMethod) GetParams() []data.GetValue    { return nil }
func (m *SplPriorityQueueRewindMethod) GetVariables() []data.Variable { return nil }
func (m *SplPriorityQueueRewindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv := spqGetCV(ctx); cv != nil {
		spqSetPos(cv, 0)
	}
	return nil, nil
}

type SplPriorityQueueValidMethod struct{}

func (m *SplPriorityQueueValidMethod) GetName() string               { return "valid" }
func (m *SplPriorityQueueValidMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplPriorityQueueValidMethod) GetIsStatic() bool             { return false }
func (m *SplPriorityQueueValidMethod) GetReturnType() data.Types     { return data.Bool{} }
func (m *SplPriorityQueueValidMethod) GetParams() []data.GetValue    { return nil }
func (m *SplPriorityQueueValidMethod) GetVariables() []data.Variable { return nil }
func (m *SplPriorityQueueValidMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := spqGetCV(ctx)
	if cv == nil {
		return data.NewBoolValue(false), nil
	}
	pos := spqGetPos(cv)
	return data.NewBoolValue(pos >= 0 && pos < len(spqGetEntries(cv).entries)), nil
}

type SplPriorityQueueCurrentMethod struct{}

func (m *SplPriorityQueueCurrentMethod) GetName() string               { return "current" }
func (m *SplPriorityQueueCurrentMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplPriorityQueueCurrentMethod) GetIsStatic() bool             { return false }
func (m *SplPriorityQueueCurrentMethod) GetReturnType() data.Types     { return data.Mixed{} }
func (m *SplPriorityQueueCurrentMethod) GetParams() []data.GetValue    { return nil }
func (m *SplPriorityQueueCurrentMethod) GetVariables() []data.Variable { return nil }
func (m *SplPriorityQueueCurrentMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := spqGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	pos := spqGetPos(cv)
	entries := spqGetEntries(cv)
	if pos < 0 || pos >= len(entries.entries) {
		return data.NewNullValue(), nil
	}
	return spqFormatExtract(entries.entries[pos], spqGetExtractFlags(cv)), nil
}

type SplPriorityQueueKeyMethod struct{}

func (m *SplPriorityQueueKeyMethod) GetName() string               { return "key" }
func (m *SplPriorityQueueKeyMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplPriorityQueueKeyMethod) GetIsStatic() bool             { return false }
func (m *SplPriorityQueueKeyMethod) GetReturnType() data.Types     { return data.Int{} }
func (m *SplPriorityQueueKeyMethod) GetParams() []data.GetValue    { return nil }
func (m *SplPriorityQueueKeyMethod) GetVariables() []data.Variable { return nil }
func (m *SplPriorityQueueKeyMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := spqGetCV(ctx)
	if cv == nil {
		return data.NewNullValue(), nil
	}
	return data.NewIntValue(spqGetPos(cv)), nil
}

type SplPriorityQueueNextMethod struct{}

func (m *SplPriorityQueueNextMethod) GetName() string               { return "next" }
func (m *SplPriorityQueueNextMethod) GetModifier() data.Modifier    { return data.ModifierPublic }
func (m *SplPriorityQueueNextMethod) GetIsStatic() bool             { return false }
func (m *SplPriorityQueueNextMethod) GetReturnType() data.Types     { return nil }
func (m *SplPriorityQueueNextMethod) GetParams() []data.GetValue    { return nil }
func (m *SplPriorityQueueNextMethod) GetVariables() []data.Variable { return nil }
func (m *SplPriorityQueueNextMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if cv := spqGetCV(ctx); cv != nil {
		spqSetPos(cv, spqGetPos(cv)+1)
	}
	return nil, nil
}
