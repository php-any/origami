package node

import "github.com/php-any/origami/data"

func (u *ContinueStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// use语句本身不返回值
	return nil, u
}

// ContinueStatement 表示continue语句
type ContinueStatement struct {
	*Node `pp:"-"`
	Level int // continue N 的层级
}

func (u *ContinueStatement) AsString() string {
	return "continue"
}

func (u *ContinueStatement) IsContinue() bool {
	return true
}

func (u *ContinueStatement) GetLabel() string {
	return ""
}

// ContinueLevel 返回 continue N 的层级，非法层级按 1 处理。
func (u *ContinueStatement) ContinueLevel() int {
	if u.Level < 1 {
		return 1
	}
	return u.Level
}

// NewContinueStatement 创建一个新的continue语句
func NewContinueStatement(token *TokenFrom) *ContinueStatement {
	return &ContinueStatement{
		Node:  NewNode(token),
		Level: 1,
	}
}

// NewContinueStatementWithLevel 创建一个带层级的continue语句
func NewContinueStatementWithLevel(token *TokenFrom, level int) *ContinueStatement {
	return &ContinueStatement{
		Node:  NewNode(token),
		Level: level,
	}
}

// continueControl 承载降级传播的 continue 层级，不作为语句节点使用。
type continueControl struct {
	level int
}

func (c *continueControl) GetValue(data.Context) (data.GetValue, data.Control) {
	return nil, c
}

func (c *continueControl) AsString() string { return "continue" }

func (c *continueControl) IsContinue() bool { return true }

func (c *continueControl) GetLabel() string { return "" }

func (c *continueControl) ContinueLevel() int { return c.level }

// consumeContinueLevel 在当前循环/switch 结构上消费 continue 的一层。
// 返回 nil 表示已在本层消化（调用方跳到下一次迭代）；
// 否则返回降级后的控制流，调用方必须原样向上返回。
func consumeContinueLevel(c data.ContinueControl) data.Control {
	level := c.ContinueLevel()
	if level <= 1 {
		return nil
	}
	next := level - 1
	if u, ok := c.(*ContinueStatement); ok {
		// 复用原节点位置；只改层级，不改共享节点本身。
		return &ContinueStatement{Node: u.Node, Level: next}
	}
	return &continueControl{level: next}
}
