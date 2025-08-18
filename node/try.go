package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

// CatchBlock 表示一个 catch 块
type CatchBlock struct {
	ExceptionType string
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

func (t *TryStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	defer func() {
		if r := recover(); r != nil {
			t.tryValue(ctx, data.NewErrorThrow(t.from, fmt.Errorf("go作用域异常退出的 panic(%v)", r)))
		}
	}()

	var v data.GetValue
	var c data.Control

	// 执行 try 块
	for _, statement := range t.TryBlock {
		v, c = statement.GetValue(ctx)
		if c != nil {
			break
		}
	}

	if c != nil {
		t.tryValue(ctx, c)
	}

	// 执行 finally 块（如果存在）
	if len(t.FinallyBlock) > 0 {
		for _, statement := range t.FinallyBlock {
			_, c = statement.GetValue(ctx)
			if c != nil {
				// finally 块中的异常会覆盖之前的异常
				return nil, c
			}
		}
	}

	return v, nil
}

func (t *TryStatement) tryValue(ctx data.Context, c data.Control) (data.GetValue, data.Control) {
	// 检查是否是异常控制
	if cv, ok := c.(*data.ClassValue); ok {
		// 查找匹配的 catch 块
		for _, catchBlock := range t.CatchBlocks {
			// 这里简化处理，假设所有异常都能被捕获
			// 在实际实现中，需要检查异常类型是否匹配

			// 将异常对象设置到 catch 变量中
			if catchBlock.Variable != nil {
				// 这里需要将异常对象设置到变量中
				if checkClassIs(ctx, cv.Class, catchBlock.Variable.GetType().String()) {
					ctx.SetVariableValue(catchBlock.Variable, c)
				} else {
					continue
				}
			}

			// 执行 catch 块
			for _, catchStmt := range catchBlock.Body {
				_, c = catchStmt.GetValue(ctx)
				if c != nil {
					// 执行 finally 块（如果存在）
					if len(t.FinallyBlock) > 0 {
						for _, statement := range t.FinallyBlock {
							_, c = statement.GetValue(ctx)
							if c != nil {
								// finally 块中的异常会覆盖之前的异常
								return nil, c
							}
						}
					}
					return nil, c
				}
			}

			// 异常已被处理，继续执行 finally 块
			break
		}
	} else if cv, ok := c.(*data.ThrowValue); ok {
		// 这里是 go 作用域返回的异常处理
		// 查找匹配的 catch 块
		for _, catchBlock := range t.CatchBlocks {
			// 这里简化处理，假设所有异常都能被捕获
			// 在实际实现中，需要检查异常类型是否匹配

			// 将异常对象设置到 catch 变量中
			if catchBlock.Variable != nil {
				// 这里需要将异常对象设置到变量中
				if checkClassIs(ctx, cv, catchBlock.Variable.GetType().String()) {
					ctx.SetVariableValue(catchBlock.Variable, c)
				} else {
					continue
				}
			}

			// 执行 catch 块
			for _, catchStmt := range catchBlock.Body {
				_, c = catchStmt.GetValue(ctx)
				if c != nil {
					// 执行 finally 块（如果存在）
					if len(t.FinallyBlock) > 0 {
						for _, statement := range t.FinallyBlock {
							_, c = statement.GetValue(ctx)
							if c != nil {
								// finally 块中的异常会覆盖之前的异常
								return nil, c
							}
						}
					}
					return nil, c
				}
			}

			// 异常已被处理，继续执行 finally 块
			break
		}
	} else {
		// 其他类型的控制流（如 return、break 等）
		return nil, c
	}

	return nil, nil
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
