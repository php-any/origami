package node

import "github.com/php-any/origami/data"

func (u *ForStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 执行初始化语句（支持多个）
	for _, initializer := range u.Initializers {
		if initializer != nil {
			_, ctl := initializer.GetValue(ctx)
			if ctl != nil {
				return nil, ctl
			}
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
					return nil, NewForYieldControl(u, bodyIndex+1, ctrl)
				}

				// return/throw 直接返回
				return nil, checkThrowControlFrom(statement, c)
			}
		}

		// 执行增量表达式（支持多个）
		for _, increment := range u.Increments {
			if increment != nil {
				_, ctl := increment.GetValue(ctx)
				if ctl != nil {
					return nil, ctl
				}
			}
		}
	}

	return v, nil
}

// ForStatement 表示for语句
type ForStatement struct {
	*Node        `pp:"-"`
	Initializers []data.GetValue
	Condition    data.GetValue
	Increments   []data.GetValue
	Body         []data.GetValue
}

// NewForStatement 创建一个新的for语句
func NewForStatement(token *TokenFrom, initializers []data.GetValue, condition data.GetValue, increments []data.GetValue, body []data.GetValue) *ForStatement {
	return &ForStatement{
		Node:         NewNode(token),
		Initializers: initializers,
		Condition:    condition,
		Increments:   increments,
		Body:         body,
	}
}

func NewForYieldControl(stmt *ForStatement, index int, v data.YieldValueControl) data.YieldControl {
	return &ForYieldControl{BodyIndex: index, ForStatement: stmt, Value: v}
}

type ForYieldControl struct {
	*ForStatement
	BodyIndex int
	Value     data.YieldValueControl
}

func (f *ForYieldControl) GetYieldKey() data.Value {
	return f.Value.GetYieldKey()
}

func (f *ForYieldControl) GetYieldValue() data.Value {
	return f.Value.GetYieldValue()
}

func (f *ForYieldControl) AsString() string {
	return "for yield"
}

func (f *ForYieldControl) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	valid, acl := f.Valid(ctx)
	if acl != nil {
		return nil, acl
	}
	for valid.(*data.BoolValue).Value == true {
		acl = f.Next(ctx)
		if acl != nil {
			return nil, acl
		}
		valid, acl = f.Valid(ctx)
		if acl != nil {
			return nil, acl
		}
	}

	return f.Value, nil
}

func (f *ForYieldControl) Current(ctx data.Context) (data.Value, data.Control) {
	return f.Value.GetYieldValue(), nil
}

func (f *ForYieldControl) Key(ctx data.Context) (data.Value, data.Control) {
	return f.Value.GetYieldKey(), nil
}

func (f *ForYieldControl) Next(ctx data.Context) data.Control {
	index := f.BodyIndex
	f.BodyIndex = 0
	var c data.Control

	for bodyIndex := index; bodyIndex < len(f.Body); bodyIndex++ {
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
				// 更新当前值并保存位置
				f.Value = ctrl
				f.BodyIndex = bodyIndex + 1
				return f
			}

			// return/throw 直接返回
			return checkThrowControlFrom(statement, c)
		}
	}

	// 执行增量表达式（支持多个）
	for _, increment := range f.Increments {
		if increment != nil {
			_, ctl := increment.GetValue(ctx)
			if ctl != nil {
				return ctl
			}
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

func (f *ForYieldControl) CreateStackState(ctx data.Context, fn data.FuncStmt, originalBody []data.GetValue, bodyIndex int) data.Generator {
	// 构建 newBody：将 bodyIndex 位置的元素（产生 YieldControl 的语句，如 for 循环）替换成 f 自身
	// 这样在 Next() 执行时，遇到 bodyIndex 位置的就是 f（YieldControl），可以直接处理，不会死循环
	newBody := originalBody[:bodyIndex]
	newBody = append(newBody, f) // 替换 bodyIndex 位置的元素为 f 自身（将 for 内部转成 yield 状态语句）
	newBody = append(newBody, originalBody[bodyIndex+1:]...)
	// 获取当前值（f 已经是 YieldControl，可以直接获取）
	currentKey := f.Value.GetYieldKey()
	currentValue := f.Value.GetYieldValue()
	// 创建 FuncYieldStackState
	return NewFuncYieldStackState(ctx, fn, newBody, bodyIndex, currentKey, currentValue)
}
