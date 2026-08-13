package node

import "github.com/php-any/origami/data"

// BreakStatement 表示break语句
type BreakStatement struct {
	*Node `pp:"-"`
	Level int // break N 的层级
}

func (u *BreakStatement) AsString() string {
	return "break"
}

func (u *BreakStatement) IsBreak() bool {
	return true
}

func (u *BreakStatement) GetLabel() string {
	return ""
}

func (u *BreakStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return nil, u
}

// NewBreakStatement 创建一个新的break语句
func NewBreakStatement(token *TokenFrom) *BreakStatement {
	return &BreakStatement{Node: NewNode(token), Level: 1}
}

// NewBreakStatementWithLevel 创建一个带层级的break语句
func NewBreakStatementWithLevel(token *TokenFrom, level int) *BreakStatement {
	return &BreakStatement{Node: NewNode(token), Level: level}
}
