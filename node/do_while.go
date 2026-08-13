package node

import "github.com/php-any/origami/data"

// DoWhileStatement 表示 do-while 语句
type DoWhileStatement struct {
	*Node     `pp:"-"`
	Condition data.GetValue
	Body      []data.GetValue
}

// NewDoWhileStatement 创建一个新的 do-while 语句
func NewDoWhileStatement(token *TokenFrom, condition data.GetValue, body []data.GetValue) *DoWhileStatement {
	return &DoWhileStatement{
		Node:      NewNode(token),
		Condition: condition,
		Body:      body,
	}
}

// GetValue 执行 do-while 循环
func (d *DoWhileStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	var v data.GetValue
	var c data.Control

	for {
		// 先执行循环体（至少执行一次）
		for _, statement := range d.Body {
			v, c = statement.GetValue(ctx)
			if c != nil {
				// break 跳出循环
				if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
					return v, nil
				}
				// continue 跳到条件判断
				if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
					break
				}
				// return/throw 直接返回
				return nil, c
			}
		}

		// 判断条件
		if d.Condition != nil {
			condValue, ctl := d.Condition.GetValue(ctx)
			if ctl != nil {
				return nil, ctl
			}
			shouldContinue := true
			if boolValue, ok := condValue.(data.AsBool); ok {
				b, err := boolValue.AsBool()
				if err != nil {
					return nil, data.NewErrorThrow(d.from, err)
				}
				shouldContinue = b
			} else {
				shouldContinue = condValue != nil
			}
			if !shouldContinue {
				break
			}
		} else {
			// 如果没有条件，只执行一次
			break
		}
	}

	return v, nil
}
