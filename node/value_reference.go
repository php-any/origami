package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

// ValueReference 表示引用取值表达式 &$var
type ValueReference struct {
	*Node `pp:"-"`
	Value data.GetValue
}

// NewValueReference 创建一个新的引用取值表达式
func NewValueReference(token *TokenFrom, value data.GetValue) *ValueReference {
	return &ValueReference{
		Node:  NewNode(token),
		Value: value,
	}
}

// GetValue 获取引用的值
func (v *ValueReference) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 检查 Value 是否是变量
	if variable, ok := v.Value.(data.Variable); ok {
		return data.NewReferenceValue(variable, ctx), nil
	}
	return nil, data.NewErrorThrow(v.from, errors.New("只能引用变量"))
}
