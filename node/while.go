package node

import "github.com/php-any/origami/data"

func (u *WhileStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	var v data.GetValue
	var c data.Control

	for {
		checkTimeLimit(u.GetFrom())
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
		for bodyIndex, statement := range u.Body {
			v, c = statement.GetValue(ctx)
			if c != nil {
				// break 跳出循环
				if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
					return nil, nil
				}
				// continue 跳到条件判断
				if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
					break
				}
				// yield：保存循环恢复状态（vlucas/phpdotenv Lexer::lex 依赖 while+yield）
				if ctrl, ok := c.(data.YieldValueControl); ok {
					return nil, NewWhileYieldControl(u, bodyIndex+1, ctrl)
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

func NewWhileYieldControl(stmt *WhileStatement, index int, v data.YieldValueControl) data.YieldControl {
	return &WhileYieldControl{BodyIndex: index, WhileStatement: stmt, Value: v}
}

// WhileYieldControl 表示 while 循环体内遇到 yield 时的暂停状态。
// 对齐 ForYieldControl：Resume 时继续执行同一次迭代剩余语句，再重新判断条件。
type WhileYieldControl struct {
	*WhileStatement
	BodyIndex int
	Value     data.YieldValueControl
}

func (w *WhileYieldControl) GetYieldKey() data.Value {
	return w.Value.GetYieldKey()
}

func (w *WhileYieldControl) GetYieldValue() data.Value {
	return w.Value.GetYieldValue()
}

func (w *WhileYieldControl) AsString() string {
	return "while yield"
}

func (w *WhileYieldControl) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 同一次迭代中途 resume：先跑完 yield 之后的语句，再重新判断条件
	// （否则会在条件已变 false 时错误地跳过剩余 body）
	if w.BodyIndex > 0 {
		if acl := w.Next(ctx); acl != nil {
			return nil, acl
		}
	}

	valid, acl := w.Valid(ctx)
	if acl != nil {
		return nil, acl
	}
	for valid.(*data.BoolValue).Value == true {
		acl = w.Next(ctx)
		if acl != nil {
			return nil, acl
		}
		valid, acl = w.Valid(ctx)
		if acl != nil {
			return nil, acl
		}
	}

	return w.Value, nil
}

func (w *WhileYieldControl) Current(ctx data.Context) (data.Value, data.Control) {
	return w.Value.GetYieldValue(), nil
}

func (w *WhileYieldControl) Key(ctx data.Context) (data.Value, data.Control) {
	return w.Value.GetYieldKey(), nil
}

func (w *WhileYieldControl) Next(ctx data.Context) data.Control {
	index := w.BodyIndex
	w.BodyIndex = 0
	var c data.Control

	for bodyIndex := index; bodyIndex < len(w.Body); bodyIndex++ {
		statement := w.Body[bodyIndex]
		_, c = statement.GetValue(ctx)

		if c != nil {
			if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
				return nil
			}
			if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
				break
			}
			if ctrl, ok := c.(data.YieldValueControl); ok {
				w.Value = ctrl
				w.BodyIndex = bodyIndex + 1
				return w
			}
			return c
		}
	}

	return nil
}

func (w *WhileYieldControl) Rewind(ctx data.Context) (data.Value, data.Control) {
	return data.NewNullValue(), nil
}

func (w *WhileYieldControl) Valid(ctx data.Context) (data.Value, data.Control) {
	if w.Condition != nil {
		condValue, ctl := w.Condition.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		shouldContinue := true
		if boolValue, ok := condValue.(data.AsBool); ok {
			b, err := boolValue.AsBool()
			if err != nil {
				return nil, data.NewErrorThrow(w.from, err)
			}
			shouldContinue = b
		} else {
			shouldContinue = condValue != nil
		}
		if !shouldContinue {
			return data.NewBoolValue(false), nil
		}
	}

	return data.NewBoolValue(true), nil
}

func (w *WhileYieldControl) CreateStackState(ctx data.Context, fn data.FuncStmt, originalBody []data.GetValue, bodyIndex int) data.Generator {
	newBody := originalBody[:bodyIndex]
	newBody = append(newBody, w)
	newBody = append(newBody, originalBody[bodyIndex+1:]...)
	currentKey := w.Value.GetYieldKey()
	currentValue := w.Value.GetYieldValue()
	return NewFuncYieldStackState(ctx, fn, newBody, bodyIndex, currentKey, currentValue)
}
