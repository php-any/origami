package node

import "github.com/php-any/origami/data"

// BreakStatement 表示break语句
type BreakStatement struct {
	*Node `pp:"-"`
}

func (u *BreakStatement) AsString() string {
	return "break"
}

func (u *BreakStatement) IsBreak() bool {
	return true
}

func (u *BreakStatement) GetLabel() string {
	//TODO implement me
	panic("implement me")
}

func (u *BreakStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// use语句本身不返回值
	return nil, u
}

// NewBreakStatement 创建一个新的break语句
func NewBreakStatement(token *TokenFrom) *BreakStatement {
	return &BreakStatement{
		Node: NewNode(token),
	}
}
