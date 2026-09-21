package node

import (
	"github.com/php-any/origami/data"
)

// EchoStatement 表示 echo 语句
type EchoStatement struct {
	*Node       `pp:"-"`
	Expressions []data.GetValue
}

// NewEchoStatement 创建一个新的 echo 语句
func NewEchoStatement(token *TokenFrom, expr []data.GetValue) *EchoStatement {
	return &EchoStatement{
		Node:        NewNode(token),
		Expressions: expr,
	}
}

// emitEchoExpr 把表达式写到输出：字面量/拼接不分配中间 StringValue。
func emitEchoExpr(ctx data.Context, expr data.GetValue) data.Control {
	switch t := expr.(type) {
	case *StringLiteral:
		return data.EmitOutput(ctx, t.Value)
	case *BinaryDot:
		if t == nil || t.Left == nil {
			return nil
		}
		if c := emitEchoExpr(ctx, t.Left); c != nil {
			return c
		}
		return emitEchoExpr(ctx, t.Right)
	case *VariableExpression:
		zv := ctx.GetIndexZVal(t.Index)
		if zv == nil || zv.Value == nil {
			return data.EmitOutput(ctx, "")
		}
		s, c := concatOperandString(ctx, unwrapValue(zv.Value))
		if c != nil {
			return c
		}
		return data.EmitOutput(ctx, s)
	default:
		v, c := expr.GetValue(ctx)
		if c != nil {
			return c
		}
		s, c := ValueToDisplayString(ctx, v)
		if c != nil {
			return c
		}
		return data.EmitOutput(ctx, s)
	}
}

// GetValue 获取 echo 语句的值
func (e *EchoStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if MarkHeaderOutputStarted != nil {
		MarkHeaderOutputStarted()
	}
	for _, expr := range e.Expressions {
		if c := emitEchoExpr(ctx, expr); c != nil {
			return nil, c
		}
	}

	return nil, nil
}
