package node

import (
	"github.com/php-any/origami/data"
)

// Node 表示语法树中的基本节点
type Node struct {
	from data.From

	value data.GetValue
}

// NewNode 创建一个新的节点
func NewNode(from data.From) *Node {
	return &Node{
		from: from,
	}
}

// GetFrom 返回节点的来源
func (n *Node) GetFrom() data.From {
	return n.from
}

// Statement 表示语句节点
type Statement interface {
	data.GetValue
}

type GetFrom interface {
	GetFrom() data.From
}

// Program 表示程序节点
type Program struct {
	*Node      `pp:"-"`
	Statements []Statement
}

// NewProgram 创建一个新的程序节点
func NewProgram(from data.From, statements []Statement) *Program {
	return &Program{
		Node:       NewNode(from),
		Statements: statements,
	}
}

// GetValue 获取程序节点的值
func (p *Program) GetValue(ctx data.Context) data.GetValue {
	var v data.GetValue
	var c data.Control
	for _, statement := range p.Statements {
		v, c = statement.GetValue(ctx)
		if c != nil {
			ctx.GetVM().ThrowControl(c)
			break
		}
	}

	return v
}
