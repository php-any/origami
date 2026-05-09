package node

import (
	"fmt"
	"runtime/debug"

	"github.com/php-any/origami/data"
)

type CatchBlock struct {
	ExceptionType data.Types
	Variable      data.Variable
	Body          []data.GetValue
}

type TryStatement struct {
	*Node        `pp:"-"`
	TryBlock     []data.GetValue
	CatchBlocks  []CatchBlock
	FinallyBlock []data.GetValue
}

func (t *TryStatement) GetValue(ctx data.Context) (v data.GetValue, c data.Control) {
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			v, c = t.tryValue(ctx, data.NewErrorThrow(t.from, fmt.Errorf("go作用域异常退出的 panic(%v)\nstack: %s", r, stack)))
		}
	}()

	for _, statement := range t.TryBlock {
		v, c = statement.GetValue(ctx)
		if c != nil {
			if add, ok := c.(data.AddStack); ok {
				add.AddStackWithInfo(statement.(GetFrom).GetFrom(), "try: ", TryGetCallClassName(statement))
			}
			break
		}
	}

	if c != nil {
		var catchValue data.GetValue
		var catchControl data.Control
		catchValue, catchControl = t.tryValue(ctx, c)
		if catchControl == nil {
			v = catchValue
			c = nil
		} else {
			v = catchValue
			c = catchControl
		}
	}

	if len(t.FinallyBlock) > 0 {
		var nAcl data.Control
		for _, statement := range t.FinallyBlock {
			_, nAcl = statement.GetValue(ctx)
			if nAcl != nil {
				return nil, nAcl
			}
		}
	}

	return v, c
}

func (t *TryStatement) tryValue(ctx data.Context, c data.Control) (data.GetValue, data.Control) {
	if cv, ok := c.(*data.ThrowValue); ok {
		for _, catchBlock := range t.CatchBlocks {
			if catchBlock.ExceptionType != nil && catchBlock.ExceptionType.Is(cv) {
				if catchBlock.Variable != nil {
					ctx.SetVariableValue(catchBlock.Variable, c)
				}

				for _, catchStmt := range catchBlock.Body {
					_, c = catchStmt.GetValue(ctx)
					if c != nil {
						return nil, c
					}
				}

				return nil, nil
			}
		}
	} else {
		return nil, c
	}

	return nil, c
}

func NewTryStatement(token *TokenFrom, tryBlock []data.GetValue, catchBlocks []CatchBlock, finallyBlock []data.GetValue) *TryStatement {
	return &TryStatement{
		Node:         NewNode(token),
		TryBlock:     tryBlock,
		CatchBlocks:  catchBlocks,
		FinallyBlock: finallyBlock,
	}
}
