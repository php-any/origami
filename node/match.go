package node

import (
	"github.com/php-any/origami/data"
)

// MatchArm 表示match语句的一个分支
type MatchArm struct {
	*Node
	Conditions []data.GetValue // 条件表达式列表
	Expression data.GetValue   // 表达式（当Statements为空时使用）
	Statements []data.GetValue // 语句列表（当Expression为空时使用）
}

// GetValue 获取match分支的值
func (m *MatchArm) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if m.Expression != nil {
		v, c := m.Expression.GetValue(ctx)
		return v, c
	}

	// 执行语句列表
	var v data.GetValue
	var c data.Control
	for _, statement := range m.Statements {
		v, c = statement.GetValue(ctx)
		if c != nil {
			return nil, c
		}
	}
	return v, nil
}

// MatchStatement 表示match语句
type MatchStatement struct {
	*Node
	Condition data.GetValue // 匹配条件
	Arms      []MatchArm    // 匹配分支列表
	Default   []data.GetValue
}

// NewMatchStatement 创建一个新的match语句
func NewMatchStatement(from data.From, condition data.GetValue, arms []MatchArm, def []data.GetValue) *MatchStatement {
	return &MatchStatement{
		Node:      NewNode(from),
		Condition: condition,
		Arms:      arms,
		Default:   def,
	}
}

// GetValue 获取match语句的值
func (m *MatchStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 计算条件值
	conditionValue, c := m.Condition.GetValue(ctx)
	if c != nil {
		return nil, c
	}

	// 遍历所有分支，找到匹配的条件
	for _, arm := range m.Arms {
		for _, condition := range arm.Conditions {
			armConditionValue, c := condition.GetValue(ctx)
			if c != nil {
				return nil, c
			}

			// 比较条件值
			if m.isMatch(conditionValue, armConditionValue) {
				return arm.GetValue(ctx)
			}
		}
	}

	if m.Default != nil {
		var v data.GetValue
		for _, stmt := range m.Default {
			v, c = stmt.GetValue(ctx)
			if c != nil {
				return nil, c
			}
		}
		return v, nil
	}

	// 如果没有匹配的分支，返回null
	return data.NewNullValue(), nil
}

// isMatch 检查两个值是否匹配
func (m *MatchStatement) isMatch(value1, value2 data.GetValue) bool {
	// 简单的相等比较，可以根据需要扩展
	if strValue1, ok := value1.(data.AsString); ok {
		if strValue2, ok := value2.(data.AsString); ok {
			return strValue1.AsString() == strValue2.AsString()
		}
	}
	return false
}
