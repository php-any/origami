package node

import "github.com/php-any/origami/data"

type LabelControl struct {
	Offset int
	Name   string
}

func (l LabelControl) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	//TODO implement me
	panic("implement me")
}

func (l LabelControl) AsString() string {
	return "LabelControl"
}

// LabelStatement 表示一个标签：label:
type LabelStatement struct {
	*Node `pp:"-"`
	Name  string
}

func NewLabelStatement(from *TokenFrom, name string) *LabelStatement {
	return &LabelStatement{
		Node: NewNode(from),
		Name: name,
	}
}

func (l *LabelStatement) AsString() string {
	return l.Name + ":"
}

func (l *LabelStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 标签本身不产生执行效果，只作为跳转目标
	return nil, LabelControl{Name: l.Name}
}

// GotoStatement 表示 goto 语句
type GotoStatement struct {
	*Node `pp:"-"`
	Label string
}

func NewGotoStatement(from *TokenFrom, label string) *GotoStatement {
	return &GotoStatement{
		Node:  NewNode(from),
		Label: label,
	}
}

func (g *GotoStatement) AsString() string {
	return "goto " + g.Label
}

// 实现 data.GotoControl
func (g *GotoStatement) IsGoto() bool {
	return true
}

func (g *GotoStatement) GetLabel() string {
	return g.Label
}

func (g *GotoStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 直接把自身作为控制流返回，由 Program 统一调度
	return nil, g
}
