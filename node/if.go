package node

import "github.com/php-any/origami/data"

// ElseIfBranch 表示一个 else if 分支
type ElseIfBranch struct {
	Condition  data.GetValue
	ThenBranch []data.GetValue
}

// IfStatement 表示if语句
type IfStatement struct {
	*Node      `pp:"-"`
	Condition  data.GetValue
	ThenBranch []data.GetValue
	ElseIf     []ElseIfBranch
	ElseBranch []data.GetValue
}

func (u *IfStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 计算条件表达式的值
	conditionValue, ctl := u.Condition.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	// 将条件值转换为布尔值
	var shouldExecuteThen bool
	if boolValue, ok := conditionValue.(data.AsBool); ok {
		b, err := boolValue.AsBool()
		if err != nil {
			return nil, data.NewErrorThrow(u.from, err)
		}
		shouldExecuteThen = b
	} else {
		// 如果无法转换为布尔值，检查是否为非空值
		shouldExecuteThen = conditionValue != nil
	}

	var v data.GetValue
	var c data.Control

	// 根据条件执行相应的分支
	if shouldExecuteThen {
		// 执行 then 分支
		for _, statement := range u.ThenBranch {
			v, c = statement.GetValue(ctx)
			if c != nil {
				return nil, c
			}
		}
	} else {
		// 检查 else if 分支
		executed := false
		for _, elseIf := range u.ElseIf {
			// 计算 else if 条件
			elseIfConditionValue, ctl := elseIf.Condition.GetValue(ctx)
			if ctl != nil {
				return nil, ctl
			}

			// 将条件值转换为布尔值
			var shouldExecuteElseIf bool
			if boolValue, ok := elseIfConditionValue.(data.AsBool); ok {
				b, err := boolValue.AsBool()
				if err != nil {
					return nil, data.NewErrorThrow(u.from, err)
				}
				shouldExecuteElseIf = b
			} else {
				// 如果无法转换为布尔值，检查是否为非空值
				shouldExecuteElseIf = elseIfConditionValue != nil
			}

			if shouldExecuteElseIf {
				// 执行 else if 分支
				for _, statement := range elseIf.ThenBranch {
					v, c = statement.GetValue(ctx)
					if c != nil {
						return nil, c
					}
				}
				executed = true
				break
			}
		}

		// 如果没有执行任何 else if 分支，且存在 else 分支，则执行 else 分支
		if !executed && len(u.ElseBranch) > 0 {
			for _, statement := range u.ElseBranch {
				v, c = statement.GetValue(ctx)
				if c != nil {
					return nil, c
				}
			}
		}
	}

	return v, nil
}

// NewIfStatement 创建一个新的if语句
func NewIfStatement(token *TokenFrom, condition data.GetValue, thenBranch []data.GetValue, elseIf []ElseIfBranch, elseBranch []data.GetValue) *IfStatement {
	return &IfStatement{
		Node:       NewNode(token),
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseIf:     elseIf,
		ElseBranch: elseBranch,
	}
}
