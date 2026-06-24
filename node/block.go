package node

import "github.com/php-any/origami/data"

// BlockStatement 表示语句块
type BlockStatement struct {
	*Node      `pp:"-"`
	Statements []data.GetValue // 语句列表
}

// NewBlockStatement 创建一个新的语句块
func NewBlockStatement(token *TokenFrom, statements []data.GetValue) *BlockStatement {
	return &BlockStatement{
		Node:       NewNode(token),
		Statements: statements,
	}
}

// GetValue 获取语句块的值
func (bs *BlockStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	var v data.GetValue
	var ctl data.Control
	for _, stmt := range bs.Statements {
		v, ctl = stmt.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
	}
	return v, nil
}
