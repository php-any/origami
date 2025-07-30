package node

import (
	"errors"
	"fmt"
	"github.com/php-any/origami/data"
)

type CallStaticProperty struct {
	*Node    `pp:"-"`
	Class    string
	Property string // 属性名
}

func NewCallStaticProperty(token *TokenFrom, path string, property string) *CallStaticProperty {
	return &CallStaticProperty{
		Node:     NewNode(token),
		Class:    path,
		Property: property,
	}
}

// GetValue 获取对象属性访问表达式的值
func (pe *CallStaticProperty) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	classStmt, ok := ctx.GetVM().GetClass(pe.Class)
	if !ok {
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("类(%s)不存在。", pe.Class)))
	}

	property, ok := classStmt.GetProperty(pe.Property)
	if ok {
		if property.GetIsStatic() {
			return property.GetDefaultValue(), nil
		}
		return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("类(%s)的属性(%s)不是静态的。", pe.Class, pe.Property)))
	}

	return nil, data.NewErrorThrow(pe.GetFrom(), errors.New(fmt.Sprintf("类(%s)没有镜静态属性(%s)。", pe.Class, pe.Property)))
}
