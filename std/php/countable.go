package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// CountableInterface 表示 PHP 的 Countable 接口
type CountableInterface struct {
	*node.InterfaceStatement
}

// NewCountableInterface 创建一个新的 Countable 接口
func NewCountableInterface() *CountableInterface {
	return &CountableInterface{
		InterfaceStatement: node.NewInterfaceStatement(
			nil,
			"Countable",
			nil,
			[]data.Method{
				&CountableCountMethod{},
			},
		),
	}
}

// CountableCountMethod 表示 Countable::count() 方法
type CountableCountMethod struct{}

func (m *CountableCountMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 这个方法应该由实现 Countable 接口的类来具体实现
	// 这里返回一个占位实现
	return data.NewIntValue(0), nil
}

func (m *CountableCountMethod) GetName() string            { return "count" }
func (m *CountableCountMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *CountableCountMethod) GetIsStatic() bool          { return false }

func (m *CountableCountMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (m *CountableCountMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (m *CountableCountMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}
