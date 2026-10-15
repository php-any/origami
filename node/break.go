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

// BreakLevel 返回 break N 的层级，非法层级按 1 处理。
func (u *BreakStatement) BreakLevel() int {
	if u.Level < 1 {
		return 1
	}
	return u.Level
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

// breakControl 承载降级传播的 break 层级，不作为语句节点使用。
type breakControl struct {
	level int
}

func (b *breakControl) GetValue(data.Context) (data.GetValue, data.Control) {
	return nil, b
}

func (b *breakControl) AsString() string { return "break" }

func (b *breakControl) IsBreak() bool { return true }

func (b *breakControl) GetLabel() string { return "" }

func (b *breakControl) BreakLevel() int { return b.level }

// consumeBreakLevel 在当前循环/switch 结构上消费 break 的一层。
// 返回 nil 表示已在本层消化（调用方退出本层）；
// 否则返回降级后的控制流，调用方必须原样向上返回。
func consumeBreakLevel(c data.BreakControl) data.Control {
	level := c.BreakLevel()
	if level <= 1 {
		return nil
	}
	next := level - 1
	if b, ok := c.(*BreakStatement); ok {
		// 复用原节点位置，错误信息仍指回源码行；只改层级，不改共享节点本身。
		return &BreakStatement{Node: b.Node, Level: next}
	}
	return &breakControl{level: next}
}
