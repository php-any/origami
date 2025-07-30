package node

import (
	"fmt"
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

// GetValue 获取 echo 语句的值
func (e *EchoStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	for _, expr := range e.Expressions {
		v, c := expr.GetValue(ctx)
		if c != nil {
			return nil, c
		}

		if s, ok := v.(data.Value); ok {
			fmt.Printf("%s", s.AsString())
		} else if v != nil {
			fmt.Printf("%s", v)
		}
	}

	return nil, nil
}
