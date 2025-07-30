package node

import "github.com/php-any/origami/data"

// BlockStatement 表示语句块
type BlockStatement struct {
	*Node      `pp:"-"`
	Statements []Statement // 语句列表
}

// NewBlockStatement 创建一个新的语句块
func NewBlockStatement(token *TokenFrom, statements []Statement) *BlockStatement {
	return &BlockStatement{
		Node:       NewNode(token),
		Statements: statements,
	}
}

// GetValue 获取语句块的值
func (bs *BlockStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}
