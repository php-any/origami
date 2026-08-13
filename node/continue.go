package node

import "github.com/php-any/origami/data"

func (u *ContinueStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// use语句本身不返回值
	return nil, u
}

// ContinueStatement 表示continue语句
type ContinueStatement struct {
	*Node `pp:"-"`
}

func (u *ContinueStatement) AsString() string {
	return "continue"
}

func (u *ContinueStatement) IsContinue() bool {
	return true
}

func (u *ContinueStatement) GetLabel() string {
	//TODO implement me
	panic("implement me")
}

// NewContinueStatement 创建一个新的continue语句
func NewContinueStatement(token *TokenFrom) *ContinueStatement {
	return &ContinueStatement{
		Node: NewNode(token),
	}
}
