package node

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/lexer"
	"github.com/php-any/origami/token"
)

type BinaryExpression interface {
	data.GetValue
}

func NewBinaryExpression(from data.From, left data.GetValue, operator lexer.Token, right data.GetValue) BinaryExpression {
	switch operator.Type {
	case token.ADD:
		return NewBinaryAdd(from, left, right)
	case token.SUB:
		return NewBinarySub(from, left, right)
	case token.ASSIGN:
		return NewBinaryAssign(from, left, right)
	case token.MUL:
		return NewBinaryMul(from, left, right)
	case token.QUO:
		return NewBinaryQuo(from, left, right)
	case token.REM:
		return NewBinaryRem(from, left, right)
	case token.EQ:
		return NewBinaryEq(from, left, right)
	case token.EQ_STRICT:
		return NewBinaryEqStrict(from, left, right)
	case token.NE:
		return NewBinaryNe(from, left, right)
	case token.NE_STRICT:
		return NewBinaryNeStrict(from, left, right)
	case token.LT:
		return NewBinaryLt(from, left, right)
	case token.LE:
		return NewBinaryLe(from, left, right)
	case token.GT:
		return NewBinaryGt(from, left, right)
	case token.GE:
		return NewBinaryGe(from, left, right)
	case token.LAND:
		return NewBinaryLand(from, left, right)
	case token.LOR:
		return NewBinaryLor(from, left, right)
	case token.DOT:
		return NewBinaryDot(from, left, right)
	// 复合赋值运算符
	case token.ADD_EQ:
		// a += b 等价于 a = a + b
		return NewBinaryAssign(from, left, NewBinaryAdd(from, left, right))
	case token.SUB_EQ:
		return NewBinaryAssign(from, left, NewBinarySub(from, left, right))
	case token.MUL_EQ:
		return NewBinaryAssign(from, left, NewBinaryMul(from, left, right))
	case token.QUO_EQ:
		return NewBinaryAssign(from, left, NewBinaryQuo(from, left, right))
	case token.REM_EQ:
		return NewBinaryAssign(from, left, NewBinaryRem(from, left, right))
	default:
		panic("unhandled default case " + operator.Literal)
	}
}
