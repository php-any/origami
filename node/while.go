package node

import "github.com/php-any/origami/data"

func (u *WhileStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	var v data.GetValue
	var c data.Control

	for {
		// 判断条件
		if u.Condition != nil {
			condValue, ctl := u.Condition.GetValue(ctx)
			if ctl != nil {
				return nil, ctl
			}
			shouldContinue := true
			if boolValue, ok := condValue.(data.AsBool); ok {
				b, err := boolValue.AsBool()
				if err != nil {
					return nil, data.NewErrorThrow(u.from, err)
				}
				shouldContinue = b
			} else {
				shouldContinue = condValue != nil
			}
			if !shouldContinue {
				break
			}
		}

		// 执行循环体
		for _, statement := range u.Body {
			v, c = statement.GetValue(ctx)
			if c != nil {
				// break 跳出循环
				if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
					return nil, nil
				}
				// continue 跳到条件判断
				if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
					continue
				}
				// return/throw 直接返回
				return nil, c
			}
		}
	}

	return v, nil
}

// WhileStatement 表示while语句
type WhileStatement struct {
	*Node     `pp:"-"`
	Condition data.GetValue
	Body      []data.GetValue
}

// NewWhileStatement 创建一个新的while语句
func NewWhileStatement(token *TokenFrom, condition data.GetValue, body []data.GetValue) *WhileStatement {
	return &WhileStatement{
		Node:      NewNode(token),
		Condition: condition,
		Body:      body,
	}
}
