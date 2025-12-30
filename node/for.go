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
		for bodyIndex, statement := range u.Body {
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
				// yield
				if ctrl, ok := c.(data.YieldValueControl); ok {
					// 构造一个yield状态LoopControl
					return nil, &ForYieldControl{BodyIndex: bodyIndex, ForStatement: u, Value: ctrl}
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

type ForYieldControl struct {
	*ForStatement
	BodyIndex int
	Value     data.YieldValueControl
}

func (f *ForYieldControl) AsString() string {
	return "for yield"
}

func (f *ForYieldControl) GetValue(ctx data.Context) (data.Value, data.Control) {
	return nil, nil
}

func (f *ForYieldControl) Current(ctx data.Context) (data.Value, data.Control) {
	return f.Value.GetYieldValue(), nil
}

func (f *ForYieldControl) Key(ctx data.Context) (data.Value, data.Control) {
	return f.Value.GetYieldKey(), nil
}

func (f *ForYieldControl) Next(ctx data.Context) data.Control {
	f.BodyIndex++

	var c data.Control

	for bodyIndex := f.BodyIndex; bodyIndex < len(f.Body); bodyIndex++ {
		statement := f.Body[bodyIndex]
		_, c = statement.GetValue(ctx)

		if c != nil {
			// break 跳出循环
			if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
				return nil
			}
			// continue 跳到增量
			if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
				break
			}
			// yield
			if ctrl, ok := c.(data.YieldValueControl); ok {
				// 构造一个yield状态LoopControl
				return &ForYieldControl{BodyIndex: bodyIndex, ForStatement: f.ForStatement, Value: ctrl}
			}

			// return/throw 直接返回
			return checkThrowControlFrom(statement, c)
		}
	}

	// 执行增量表达式
	if f.Increment != nil {
		_, ctl := f.Increment.GetValue(ctx)
		if ctl != nil {
			return ctl
		}
	}

	return nil
}

func (f *ForYieldControl) Rewind(ctx data.Context) (data.Value, data.Control) {
	return data.NewNullValue(), nil
}

func (f *ForYieldControl) Valid(ctx data.Context) (data.Value, data.Control) {
	// 判断条件
	if f.Condition != nil {
		condValue, ctl := f.Condition.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		shouldContinue := true
		if boolValue, ok := condValue.(data.AsBool); ok {
			b, err := boolValue.AsBool()
			if err != nil {
				return nil, data.NewErrorThrow(f.from, err)
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
