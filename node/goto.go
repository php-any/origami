package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

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
	// 标签本身不产生执行效果，只作为跳转目标（由函数体/Program 在 goto 时按节点查找）
	return nil, nil
}

// findLabelBodyIndex 在语句列表中查找标签，返回标签后下一条语句的下标。
func findLabelBodyIndex(statements []data.GetValue, name string) (int, bool) {
	for i, stmt := range statements {
		if ls, ok := stmt.(*LabelStatement); ok && ls.Name == name {
			return i + 1, true
		}
	}
	return 0, false
}

// resolveGotoBodyIndex 将 GotoControl 解析为函数体/脚本体中的跳转下标（for 循环用：返回值需再 -1）。
func resolveGotoBodyIndex(from data.From, statements []data.GetValue, gotoCtl data.GotoControl) (int, data.Control) {
	offset, ok := findLabelBodyIndex(statements, gotoCtl.GetLabel())
	if !ok {
		if g, ok := gotoCtl.(*GotoStatement); ok {
			from = g.GetFrom()
		}
		return 0, data.NewErrorThrow(from, fmt.Errorf("未定义的标签 '%s'", gotoCtl.GetLabel()))
	}
	return offset, nil
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

func (g *GotoStatement) ForWrite() data.Value {
	return g
}
