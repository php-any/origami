package lexer

import "github.com/php-any/origami/token"

// LingToken 表示一个插值块 token（包含多个子 token）
type LingToken struct {
	type_    token.TokenType
	literal  string  // 原始值
	start    int     // 起始位置
	end      int     // 结束位置
	line     int     // 行号
	pos      int     // 单独一行的位置
	children []Token // 子 token 列表（插值块内的 token）
}

// NewLingToken 创建一个新的 LingToken
func NewLingToken(type_ token.TokenType, literal string, start, end, line, pos int, children []Token) *LingToken {
	return &LingToken{
		type_:    type_,
		literal:  literal,
		start:    start,
		end:      end,
		line:     line,
		pos:      pos,
		children: children,
	}
}

// Type 返回 token 的类型
func (l *LingToken) Type() token.TokenType {
	return l.type_
}

// Literal 返回 token 的字面值
func (l *LingToken) Literal() string {
	return l.literal
}

// Start 返回 token 的起始位置
func (l *LingToken) Start() int {
	return l.start
}

// End 返回 token 的结束位置
func (l *LingToken) End() int {
	return l.end
}

// Line 返回 token 所在的行号
func (l *LingToken) Line() int {
	return l.line
}

// Pos 返回 token 在单独一行中的位置
func (l *LingToken) Pos() int {
	return l.pos
}

// Children 返回插值块内的子 token 列表
func (l *LingToken) Children() []Token {
	return l.children
}
