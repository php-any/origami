package parser

import (
	"github.com/php-any/origami/node"
)

// PositionTracker 位置跟踪器，用于自动跟踪解析范围
type PositionTracker struct {
	parser    *Parser
	startPos  int
	endPos    int
	isStarted bool
	isEnded   bool
}

// StartTracking 开始位置跟踪
func (p *Parser) StartTracking() *PositionTracker {
	return &PositionTracker{
		parser:    p,
		startPos:  p.position,
		isStarted: true,
	}
}

// StartTrackingAt 从指定位置开始跟踪
func (p *Parser) StartTrackingAt(pos int) *PositionTracker {
	return &PositionTracker{
		parser:    p,
		startPos:  pos,
		isStarted: true,
	}
}

// End 结束位置跟踪并返回From信息
func (pt *PositionTracker) End() *node.TokenFrom {
	if !pt.isStarted {
		return pt.parser.FromCurrentToken()
	}

	pt.endPos = pt.parser.position
	pt.isEnded = true

	return pt.parser.FromPositionRange(pt.startPos, pt.endPos)
}

// EndAt 在指定位置结束跟踪
func (pt *PositionTracker) EndAt(pos int) *node.TokenFrom {
	if !pt.isStarted {
		return pt.parser.FromCurrentToken()
	}

	pt.endPos = pos
	pt.isEnded = true

	return pt.parser.FromPositionRange(pt.startPos, pt.endPos)
}

// EndBefore 在当前位置之前结束跟踪（不包含当前token）
func (pt *PositionTracker) EndBefore() *node.TokenFrom {
	if !pt.isStarted {
		return pt.parser.FromCurrentToken()
	}

	endPos := pt.parser.position - 1
	if endPos < pt.startPos {
		endPos = pt.startPos
	}

	pt.endPos = endPos
	pt.isEnded = true

	return pt.parser.FromPositionRange(pt.startPos, pt.endPos)
}

// GetStartPos 获取开始位置
func (pt *PositionTracker) GetStartPos() int {
	return pt.startPos
}

// GetEndPos 获取结束位置
func (pt *PositionTracker) GetEndPos() int {
	if pt.isEnded {
		return pt.endPos
	}
	return pt.parser.position
}

// IsValid 检查跟踪器是否有效
func (pt *PositionTracker) IsValid() bool {
	return pt.isStarted && pt.startPos < len(pt.parser.tokens)
}

// UpdateStart 更新开始位置
func (pt *PositionTracker) UpdateStart(pos int) *PositionTracker {
	pt.startPos = pos
	return pt
}

// Clone 克隆位置跟踪器
func (pt *PositionTracker) Clone() *PositionTracker {
	return &PositionTracker{
		parser:    pt.parser,
		startPos:  pt.startPos,
		endPos:    pt.endPos,
		isStarted: pt.isStarted,
		isEnded:   pt.isEnded,
	}
}
