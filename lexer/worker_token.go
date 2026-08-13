package lexer

import "github.com/php-any/origami/token"

// WorkerToken 表示一个单词 token（普通词法单元）
type WorkerToken struct {
	type_   token.TokenType
	literal string // 原始值, 换行替换为;符号也不能替换Literal
	start   int    // 起始位置
	end     int    // 结束位置
	line    int    // 行号
	pos     int    // 单独一行的位置
}

// NewWorkerToken 创建一个新的 WorkerToken
func NewWorkerToken(type_ token.TokenType, literal string, start, end, line, pos int) *WorkerToken {
	return &WorkerToken{
		type_:   type_,
		literal: literal,
		start:   start,
		end:     end,
		line:    line,
		pos:     pos,
	}
}

// Type 返回 token 的类型
func (w *WorkerToken) Type() token.TokenType {
	return w.type_
}

// Literal 返回 token 的字面值
func (w *WorkerToken) Literal() string {
	return w.literal
}

// Start 返回 token 的起始位置
func (w *WorkerToken) Start() int {
	return w.start
}

// End 返回 token 的结束位置
func (w *WorkerToken) End() int {
	return w.end
}

// Line 返回 token 所在的行号
func (w *WorkerToken) Line() int {
	return w.line
}

// Pos 返回 token 在单独一行中的位置
func (w *WorkerToken) Pos() int {
	return w.pos
}
