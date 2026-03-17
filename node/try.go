package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// CatchBlock 表示一个 catch 块（单类型为 BaseType，多类型为 NewUnionType）
type CatchBlock struct {
	ExceptionType data.Types
	Variable      data.Variable
	Body          []data.GetValue
}

// TryStatement 表示try语句
type TryStatement struct {
	*Node        `pp:"-"`
	TryBlock     []data.GetValue
	CatchBlocks  []CatchBlock
	FinallyBlock []data.GetValue
}

func (t *TryStatement) GetValue(ctx data.Context) (v data.GetValue, c data.Control) {
	defer func() {
		if r := recover(); r != nil {
			v, c = t.tryValue(ctx, data.NewErrorThrow(t.from, fmt.Errorf("go作用域异常退出的 panic(%v)", r)))
		}
	}()

	// 执行 try 块
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

		// 如果 catch 块处理了异常，使用 catch 块的值和控制流
		// 如果 catch 块没有处理异常（没有匹配的 catch），catchControl 会是原来的异常
		if catchControl == nil {
			// 异常已被 catch 处理，使用 catch 块的值
			v = catchValue
			c = nil
		} else {
			// 异常未被处理，继续传播
			// 但 finally 块仍然需要执行
			v = catchValue
			c = catchControl
		}
	}

	// 执行 finally 块（如果存在）
	if len(t.FinallyBlock) > 0 {
		var nAcl data.Control
		for _, statement := range t.FinallyBlock {
			_, nAcl = statement.GetValue(ctx)
			if nAcl != nil {
				// finally 块中的异常会覆盖之前的异常
				return nil, nAcl
			}
		}
	}

	return v, c
}

func (t *TryStatement) tryValue(ctx data.Context, c data.Control) (data.GetValue, data.Control) {
	// 检查是否是异常控制
	if cv, ok := c.(*data.ThrowValue); ok {
		// 这里是 go 作用域返回的异常处理
		// 查找匹配的 catch 块（直接使用 ExceptionType.Is(异常值) 判断，UnionType 内任一匹配即可）
		for _, catchBlock := range t.CatchBlocks {
			if catchBlock.ExceptionType != nil && catchBlock.ExceptionType.Is(cv) {
				if catchBlock.Variable != nil {
					ctx.SetVariableValue(catchBlock.Variable, c)
				}

				// 执行 catch 块
				for _, catchStmt := range catchBlock.Body {
					_, c = catchStmt.GetValue(ctx)
					if c != nil {
						// catch 块中有新的异常或 return，直接返回
						// finally 块会在 GetValue 方法中执行
						return nil, c
					}
				}

				// 异常已被处理，返回 nil 表示已处理
				// finally 块会在 GetValue 方法中执行
				return nil, nil
			}
		}
	} else {
		// 其他类型的控制流（如 return、break 等）
		return nil, c
	}

	return nil, c
}

// NewTryStatement 创建一个新的try语句
func NewTryStatement(token *TokenFrom, tryBlock []data.GetValue, catchBlocks []CatchBlock, finallyBlock []data.GetValue) *TryStatement {
	return &TryStatement{
		Node:         NewNode(token),
		TryBlock:     tryBlock,
		CatchBlocks:  catchBlocks,
		FinallyBlock: finallyBlock,
	}
}
