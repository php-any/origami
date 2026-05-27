package node

import (
	"fmt"

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

type GetFrom interface {
	GetFrom() data.From
}

// Program 表示程序节点
type Program struct {
	*Node      `pp:"-"`
	Statements []data.GetValue
}

// NewProgram 创建一个新的程序节点
func NewProgram(from data.From, statements []data.GetValue) *Program {
	return &Program{
		Node:       NewNode(from),
		Statements: statements,
	}
}

// GetValue 获取程序节点的值
func (p *Program) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	var v data.GetValue
	var c data.Control
	for offset, statement := range p.Statements {
		v, c = statement.GetValue(ctx)
		if c != nil {
			switch acl := c.(type) {
			case data.ReturnControl:
				return c.GetValue(ctx)
			case LabelControl:
				acl.Offset = offset + 1
				return p.runLabel(ctx, acl)
			case data.GotoControl:
				return p.runGoto(ctx, acl)
			default:
				if acl, ok := acl.(data.AddStack); ok {
					if statement, ok := statement.(GetFrom); ok {
						acl.AddStackWithInfo(statement.GetFrom(), "program", "")
					}
				}
				ctx.GetVM().ThrowControl(c)
				return v, nil
			}
		}
	}

	return v, nil
}

func (p *Program) findLabelOffset(name string) (int, bool) {
	for i, stmt := range p.Statements {
		if ls, ok := stmt.(*LabelStatement); ok && ls.Name == name {
			return i + 1, true
		}
	}
	return 0, false
}

func (p *Program) runGoto(ctx data.Context, gotoCtl data.GotoControl) (data.GetValue, data.Control) {
	offset, ok := p.findLabelOffset(gotoCtl.GetLabel())
	if !ok {
		var from data.From
		if g, ok := gotoCtl.(*GotoStatement); ok {
			from = g.GetFrom()
		}
		return nil, data.NewErrorThrow(from, fmt.Errorf("未定义的标签 '%s'", gotoCtl.GetLabel()))
	}
	return p.runLabel(ctx, LabelControl{Name: gotoCtl.GetLabel(), Offset: offset})
}

func (p *Program) runLabel(ctx data.Context, label LabelControl) (data.GetValue, data.Control) {
	var v data.GetValue
	var c data.Control
	for offset := label.Offset; offset < len(p.Statements); offset++ {
		statement := p.Statements[offset]
		v, c = statement.GetValue(ctx)
		if c != nil {
			switch acl := c.(type) {
			case data.ReturnControl:
				return c.GetValue(ctx)
			case LabelControl:
				acl.Offset = offset + 1
				return p.runLabel(ctx, acl)
			case data.GotoControl:
				if acl.GetLabel() == label.Name {
					return p.runLabel(ctx, label)
				}
				return p.runGoto(ctx, acl)
			default:
				ctx.GetVM().ThrowControl(c)
				return v, nil
			}
		}
	}

	return v, nil
}
