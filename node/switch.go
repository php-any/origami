package node

import (
	"github.com/php-any/origami/data"
)

// SwitchCase 表示switch语句的一个分支
type SwitchCase struct {
	*Node
	CaseValue  data.GetValue   // case值
	Statements []data.GetValue // 语句列表
}

// GetValue 获取switch分支的值
func (s *SwitchCase) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 执行语句列表
	var v data.GetValue
	var c data.Control
	for _, statement := range s.Statements {
		if stmt, ok := statement.(data.GetValue); ok {
			v, c = stmt.GetValue(ctx)
		} else {
			// 如果是表达式，直接获取值
			v = statement
		}
		if c != nil {
			if _, ok := c.(data.BreakControl); ok {
				return v, nil
			}
		}
	}
	return v, nil
}

// SwitchStatement 表示switch语句
type SwitchStatement struct {
	*Node
	Condition   data.GetValue   // 匹配条件
	Cases       []SwitchCase    // case分支列表
	DefaultCase []data.GetValue // default分支
}

// NewSwitchStatement 创建一个新的switch语句
func NewSwitchStatement(from data.From, condition data.GetValue, cases []SwitchCase, defaultCase []data.GetValue) *SwitchStatement {
	return &SwitchStatement{
		Node:        NewNode(from),
		Condition:   condition,
		Cases:       cases,
		DefaultCase: defaultCase,
	}
}

// GetValue 获取switch语句的值
func (s *SwitchStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 计算条件值
	conditionValue, c := s.Condition.GetValue(ctx)
	if c != nil {
		return nil, c
	}

	// 遍历所有case分支，找到匹配的条件
	for _, caseStmt := range s.Cases {
		caseValue, c := caseStmt.CaseValue.GetValue(ctx)
		if c != nil {
			return nil, c
		}

		// 比较条件值
		if s.isMatch(conditionValue, caseValue) {
			return caseStmt.GetValue(ctx)
		}
	}

	// 如果没有匹配的case，执行default分支
	if len(s.DefaultCase) > 0 {
		var v data.GetValue
		var c data.Control
		for _, statement := range s.DefaultCase {
			if stmt, ok := statement.(data.GetValue); ok {
				v, c = stmt.GetValue(ctx)
			} else {
				// 如果是表达式，直接获取值
				v = statement
			}
			if c != nil {
				if _, ok := c.(data.BreakControl); ok {
					return v, nil
				}
			}
		}
		return v, nil
	}

	// 如果没有匹配的分支，返回null
	return data.NewNullValue(), nil
}

// isMatch 检查两个值是否匹配
func (s *SwitchStatement) isMatch(value1, value2 data.GetValue) bool {
	// 简单的相等比较，可以根据需要扩展
	if strValue1, ok := value1.(data.AsString); ok {
		if strValue2, ok := value2.(data.AsString); ok {
			return strValue1.AsString() == strValue2.AsString()
		}
	}
	return false
}
