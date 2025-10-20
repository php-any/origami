package node

import "github.com/php-any/origami/data"

func (u *ForStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 执行初始化语句
	if u.Initializer != nil {
		_, ctl := u.Initializer.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
	}

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
				// continue 跳到增量
				if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
					break
				}

				// return/throw 直接返回
				return nil, checkThrowControlFrom(statement, c)
			}
		}

		// 执行增量表达式
		if u.Increment != nil {
			_, ctl := u.Increment.GetValue(ctx)
			if ctl != nil {
				return nil, ctl
			}
		}
	}

	return v, nil
}

// ForStatement 表示for语句
type ForStatement struct {
	*Node       `pp:"-"`
	Initializer data.GetValue
	Condition   data.GetValue
	Increment   data.GetValue
	Body        []data.GetValue
}

// NewForStatement 创建一个新的for语句
func NewForStatement(token *TokenFrom, initializer data.GetValue, condition data.GetValue, increment data.GetValue, body []data.GetValue) *ForStatement {
	return &ForStatement{
		Node:        NewNode(token),
		Initializer: initializer,
		Condition:   condition,
		Increment:   increment,
		Body:        body,
	}
}
