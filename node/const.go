package node

import "github.com/php-any/origami/data"

func (u *ConstStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	v, acl := u.Initializer.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	// 跳过类型检查，直接赋值
	ctx.SetVariableValue(u.Val, v.(data.Value))

	return v, nil
}

// ConstStatement 表示常量声明语句
type ConstStatement struct {
	*Node       `pp:"-"`
	Val         data.Variable
	Initializer data.GetValue
}

// NewConstStatement 创建一个新的常量声明语句
func NewConstStatement(token *TokenFrom, val data.Variable, initializer data.GetValue) *ConstStatement {
	return &ConstStatement{
		Node:        NewNode(token),
		Val:         val,
		Initializer: initializer,
	}
}
