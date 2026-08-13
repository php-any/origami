package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// SplQueueClass 实现 PHP SPL �?SplQueue（FIFO），继承 SplDoublyLinkedList
type SplQueueClass struct {
	node.Node
}

func NewSplQueueClass() *SplQueueClass {
	return &SplQueueClass{}
}

func (c *SplQueueClass) GetName() string { return "SplQueue" }
func (c *SplQueueClass) GetExtend() *string {
	parent := "SplDoublyLinkedList"
	return &parent
}
func (c *SplQueueClass) GetImplements() []string {
	return []string{"Iterator", "Countable"}
}
func (c *SplQueueClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SplQueueClass) GetPropertyList() []data.Property              { return nil }
func (c *SplQueueClass) GetConstruct() data.Method                     { return nil }
func (c *SplQueueClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	splListInitCV(cv)
	return cv, nil
}

func (c *SplQueueClass) GetMethod(name string) (data.Method, bool) {
	return splExtendGetMethod(c, name, func(name string) (data.Method, bool) {
		switch name {
		case "enqueue":
			return &SplQueueEnqueueMethod{}, true
		case "dequeue":
			return &SplQueueDequeueMethod{}, true
		}
		return nil, false
	})
}

func (c *SplQueueClass) GetMethods() []data.Method {
	return splExtendGetMethods(c, []data.Method{
		&SplQueueEnqueueMethod{},
		&SplQueueDequeueMethod{},
	})
}

type SplQueueEnqueueMethod struct{}

func (m *SplQueueEnqueueMethod) GetName() string            { return "enqueue" }
func (m *SplQueueEnqueueMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplQueueEnqueueMethod) GetIsStatic() bool          { return false }
func (m *SplQueueEnqueueMethod) GetReturnType() data.Types  { return nil }
func (m *SplQueueEnqueueMethod) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameter(nil, "value", 0, nil, data.Mixed{})}
}
func (m *SplQueueEnqueueMethod) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "value", 0, data.Mixed{})}
}
func (m *SplQueueEnqueueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := splListGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	val, _ := ctx.GetIndexValue(0)
	arr := splListGetStorage(cv)
	arr.List = append(arr.List, data.NewZVal(val))
	return nil, nil
}

type SplQueueDequeueMethod struct{}

func (m *SplQueueDequeueMethod) GetName() string            { return "dequeue" }
func (m *SplQueueDequeueMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *SplQueueDequeueMethod) GetIsStatic() bool          { return false }
func (m *SplQueueDequeueMethod) GetReturnType() data.Types  { return data.Mixed{} }
func (m *SplQueueDequeueMethod) GetParams() []data.GetValue { return nil }
func (m *SplQueueDequeueMethod) GetVariables() []data.Variable {
	return nil
}
func (m *SplQueueDequeueMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
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
