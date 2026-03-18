package node

import "github.com/php-any/origami/data"

func (u *VarStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// use语句本身不返回值
	return nil, nil
}

// VarStatement 表示变量声明语句
type VarStatement struct {
	*Node       `pp:"-"`
	Name        string
	Initializer data.GetValue
}

// NewVarStatement 创建一个新的变量声明语句
func NewVarStatement(token *TokenFrom, name string, initializer data.GetValue) *VarStatement {
	return &VarStatement{
		Node:        NewNode(token),
		Name:        name,
		Initializer: initializer,
	}
}

// StaticVarStatement 表示静态局部变量声明语句
type StaticVarStatement struct {
	*Node       `pp:"-"`
	Name        string
	Initializer data.GetValue
}

// NewStaticVarStatement 创建一个新的静态局部变量声明语句
func NewStaticVarStatement(token *TokenFrom, name string, initializer data.GetValue) *StaticVarStatement {
	return &StaticVarStatement{
		Node:        NewNode(token),
		Name:        name,
		Initializer: initializer,
	}
}

func (s *StaticVarStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 静态局部变量语句本身不返回值
	return s.Initializer, nil
}
