package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IteratorInterface 实现 PHP 的 Iterator 接口
type IteratorInterface struct {
	node.Node
}

func NewIteratorInterface() *IteratorInterface {
	return &IteratorInterface{}
}

func (i *IteratorInterface) GetName() string {
	return "Iterator"
}

func (i *IteratorInterface) GetExtend() *string {
	// Iterator 继承自 Traversable
	traversable := "Traversable"
	return &traversable
}

func (i *IteratorInterface) GetImplements() []string {
	return nil
}

func (i *IteratorInterface) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (i *IteratorInterface) GetPropertyList() []data.Property {
	return nil
}

func (i *IteratorInterface) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(i, ctx.CreateBaseContext()), nil
}

func (i *IteratorInterface) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "current":
		return &IteratorCurrent{}, true
	case "key":
		return &IteratorKey{}, true
	case "next":
		return &IteratorNext{}, true
	case "rewind":
		return &IteratorRewind{}, true
	case "valid":
		return &IteratorValid{}, true
	}
	return nil, false
}

func (i *IteratorInterface) GetMethods() []data.Method {
	return []data.Method{
		&IteratorCurrent{},
		&IteratorKey{},
		&IteratorNext{},
		&IteratorRewind{},
		&IteratorValid{},
	}
}

func (i *IteratorInterface) GetConstruct() data.Method {
	return nil
}

// IteratorCurrent current 方法
type IteratorCurrent struct{}

func (m *IteratorCurrent) GetName() string {
	return "current"
}

func (m *IteratorCurrent) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *IteratorCurrent) GetIsStatic() bool {
	return false
}

func (m *IteratorCurrent) GetVariables() []data.Variable {
	return nil
}

func (m *IteratorCurrent) GetReturnType() data.Types {
	return nil
}

func (m *IteratorCurrent) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *IteratorCurrent) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

// IteratorKey key 方法
type IteratorKey struct{}

func (m *IteratorKey) GetName() string {
	return "key"
}

func (m *IteratorKey) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *IteratorKey) GetIsStatic() bool {
	return false
}

func (m *IteratorKey) GetVariables() []data.Variable {
	return nil
}

func (m *IteratorKey) GetReturnType() data.Types {
	return nil
}

func (m *IteratorKey) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *IteratorKey) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

// IteratorNext next 方法
type IteratorNext struct{}

func (m *IteratorNext) GetName() string {
	return "next"
}

func (m *IteratorNext) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *IteratorNext) GetIsStatic() bool {
	return false
}

func (m *IteratorNext) GetVariables() []data.Variable {
	return nil
}

func (m *IteratorNext) GetReturnType() data.Types {
	return nil
}

func (m *IteratorNext) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *IteratorNext) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

// IteratorRewind rewind 方法
type IteratorRewind struct{}

func (m *IteratorRewind) GetName() string {
	return "rewind"
}

func (m *IteratorRewind) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *IteratorRewind) GetIsStatic() bool {
	return false
}

func (m *IteratorRewind) GetVariables() []data.Variable {
	return nil
}

func (m *IteratorRewind) GetReturnType() data.Types {
	return nil
}

func (m *IteratorRewind) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *IteratorRewind) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

// IteratorValid valid 方法
type IteratorValid struct{}

func (m *IteratorValid) GetName() string {
	return "valid"
}

func (m *IteratorValid) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (m *IteratorValid) GetIsStatic() bool {
	return false
}

func (m *IteratorValid) GetVariables() []data.Variable {
	return nil
}

func (m *IteratorValid) GetReturnType() data.Types {
	return nil
}

func (m *IteratorValid) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *IteratorValid) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}
