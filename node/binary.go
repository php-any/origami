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
	switch operator.Type() {
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
	case token.BIT_AND:
		return NewBinaryBitAnd(from, left, right)
	case token.BIT_XOR:
		return NewBinaryBitXor(from, left, right)
	case token.BIT_OR:
		return NewBinaryBitOr(from, left, right)
	case token.SHL:
		return NewBinaryShl(from, left, right)
	case token.SHR:
		return NewBinaryShr(from, left, right)
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
	case token.CONCAT_EQ:
		// a .= b 等价于 a = a . b
		return NewBinaryAssign(from, left, NewBinaryDot(from, left, right))
	case token.SHL_EQ:
		// a <<= b 等价于 a = a << b
		return NewBinaryAssign(from, left, NewBinaryShl(from, left, right))
	case token.SHR_EQ:
		// a >>= b 等价于 a = a >> b
		return NewBinaryAssign(from, left, NewBinaryShr(from, left, right))
	case token.NULL_COALESCE_ASSIGN:
		// a ??= b 等价于 a = a ?? b（仅在 a 为 null 时赋值）
		// 需要将 from 转换为 *TokenFrom
		var tokenFrom *TokenFrom
		if tf, ok := from.(*TokenFrom); ok {
			tokenFrom = tf
		} else {
			// 如果无法转换，创建一个简单的 TokenFrom
			// 使用 NewTokenFrom 创建，传入基本参数
			tokenFrom = NewTokenFrom(nil, 0, 0, 0, 0)
		}
		return NewBinaryAssign(from, left, NewNullCoalesceExpression(tokenFrom, left, right))
	default:
		panic("unhandled default case " + operator.Literal())
	}
}
