package parser

import (
	"github.com/php-any/origami/node"
)

// FromBuilder 用于构建From信息的接口
type FromBuilder interface {
	SetStart(pos int) FromBuilder
	SetEnd(pos int) FromBuilder
	SetLine(line int) FromBuilder
	SetPos(pos int) FromBuilder
	Build() *node.TokenFrom
}

// fromBuilder 实现FromBuilder接口
type fromBuilder struct {
	parser   *Parser
	startPos int
	endPos   int
	line     int
	pos      int
	hasStart bool
	hasEnd   bool
	hasLine  bool
	hasPos   bool
}

// NewFromBuilder 创建一个新的FromBuilder
func (p *Parser) NewFromBuilder() FromBuilder {
	return &fromBuilder{
		parser: p,
	}
}

// SetStart 设置起始位置
func (fb *fromBuilder) SetStart(pos int) FromBuilder {
	fb.startPos = pos
	fb.hasStart = true
	return fb
}

// SetEnd 设置结束位置
func (fb *fromBuilder) SetEnd(pos int) FromBuilder {
	fb.endPos = pos
	fb.hasEnd = true
	return fb
}

// SetLine 设置行号
func (fb *fromBuilder) SetLine(line int) FromBuilder {
	fb.line = line
	fb.hasLine = true
	return fb
}

// SetPos 设置行内位置
func (fb *fromBuilder) SetPos(pos int) FromBuilder {
	fb.pos = pos
	fb.hasPos = true
	return fb
}

// Build 构建TokenFrom对象
func (fb *fromBuilder) Build() *node.TokenFrom {
	// 如果没有设置起始位置，使用当前token的起始位置
	if !fb.hasStart {
		fb.startPos = fb.parser.current().Start()
	}

	// 如果没有设置结束位置，使用当前token的结束位置
	if !fb.hasEnd {
		fb.endPos = fb.parser.current().End()
	}

	// 如果没有设置行号，使用起始位置对应的行号
	if !fb.hasLine {
		fb.line = fb.parser.getLineByPosition(fb.startPos)
	}

	// 如果没有设置行内位置，使用起始位置对应的行内位置
	if !fb.hasPos {
		fb.pos = fb.parser.getPosInLineByPosition(fb.startPos)
	}

	// 创建 TokenFrom 并设置结束位置
	tf := node.NewTokenFrom(fb.parser.source, fb.startPos, fb.endPos, fb.line, fb.pos)

	// 如果结束位置与开始位置不同，需要计算结束位置的行号和列号
	if fb.hasEnd && fb.endPos != fb.startPos {
		endLine := fb.parser.getLineByPosition(fb.endPos)
		endPos := fb.parser.getPosInLineByPosition(fb.endPos)
		tf.SetEndPosition(endLine, endPos)
	}

	return tf
}

// FromCurrentToken 从当前token创建From信息
// deprecated: p.StartTracking()
func (p *Parser) FromCurrentToken() *node.TokenFrom {
	current := p.current()
	return node.NewTokenFrom(p.source, current.Start(), current.End(), current.Line(), current.Pos())
}

// FromRange 从指定范围创建From信息
func (p *Parser) FromRange(start, end int) *node.TokenFrom {
	line := p.getLineByPosition(start)
	pos := p.getPosInLineByPosition(start)
	return node.NewTokenFrom(p.source, start, end, line, pos)
}

// FromTokenRange 从token范围创建From信息
func (p *Parser) FromTokenRange(startToken, endToken int) *node.TokenFrom {
	if startToken >= len(p.tokens) || endToken >= len(p.tokens) {
		return p.FromCurrentToken()
	}

	start := p.tokens[startToken].Start()
	end := p.tokens[endToken].End()
	line := p.tokens[startToken].Line()
	pos := p.tokens[startToken].Pos()

	return node.NewTokenFrom(p.source, start, end, line, pos)
}

// getLineByPosition 根据位置获取行号
func (p *Parser) getLineByPosition(pos int) int {
	if p.source == nil || pos < 0 || pos >= len(*p.source) {
		return 0
	}

	line := 0
	for i := 0; i < pos && i < len(*p.source); i++ {
		if (*p.source)[i] == '\n' {
			line++
		}
	}
	return line
}

// getPosInLineByPosition 根据位置获取行内位置
func (p *Parser) getPosInLineByPosition(pos int) int {
	if p.source == nil || pos < 0 || pos >= len(*p.source) {
		return 0
	}

	linePos := 0
	for i := 0; i < pos && i < len(*p.source); i++ {
		if (*p.source)[i] == '\n' {
			linePos = 0
		} else {
			linePos++
		}
	}
	return linePos
}
