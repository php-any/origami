package node

import (
	"github.com/php-any/origami/data"
)

// SwitchCase 表示 switch 语句的一个分支
type SwitchCase struct {
	*Node
	CaseValue  data.GetValue   // case 值
	Statements []data.GetValue // 语句列表
}

// GetValue 执行该 case 的语句列表。
// break / return / throw 通过 Control 向上传递，供 SwitchStatement 处理 fall-through。
func (s *SwitchCase) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	var v data.GetValue
	var c data.Control
	for _, statement := range s.Statements {
		if stmt, ok := statement.(data.GetValue); ok {
			v, c = stmt.GetValue(ctx)
		} else {
			v = statement
		}
		if c != nil {
			return v, c
		}
	}
	return v, nil
}

// SwitchStatement 表示 switch 语句
type SwitchStatement struct {
	*Node
	Condition   data.GetValue   // 匹配条件
	Cases       []SwitchCase    // case 分支列表
	DefaultCase []data.GetValue // default 分支
}

// NewSwitchStatement 创建一个新的 switch 语句
func NewSwitchStatement(from data.From, condition data.GetValue, cases []SwitchCase, defaultCase []data.GetValue) *SwitchStatement {
	return &SwitchStatement{
		Node:        NewNode(from),
		Condition:   condition,
		Cases:       cases,
		DefaultCase: defaultCase,
	}
}

// GetValue 获取 switch 语句的值（支持 PHP case 穿透）
func (s *SwitchStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	conditionValue, c := s.Condition.GetValue(ctx)
	if c != nil {
		return nil, c
	}

	for i := range s.Cases {
		caseValue, c := s.Cases[i].CaseValue.GetValue(ctx)
		if c != nil {
			return nil, c
		}

		if !s.isMatch(conditionValue, caseValue) {
			continue
		}

		// 从匹配的 case 开始向下穿透，直到 break/return/throw
		var v data.GetValue
		for j := i; j < len(s.Cases); j++ {
			var ctl data.Control
			v, ctl = s.Cases[j].GetValue(ctx)
			if ctl != nil {
				switch ctl.(type) {
				case data.BreakControl:
					return v, nil
				default:
					return v, ctl
				}
			}
		}
		return v, nil
	}

	if len(s.DefaultCase) > 0 {
		var v data.GetValue
		var c data.Control
		for _, statement := range s.DefaultCase {
			if stmt, ok := statement.(data.GetValue); ok {
				v, c = stmt.GetValue(ctx)
			} else {
				v = statement
			}
			if c != nil {
				switch c.(type) {
				case data.BreakControl:
					return v, nil
				default:
					return v, c
				}
			}
		}
		return v, nil
	}

	return data.NewNullValue(), nil
}

// isMatch 检查两个值是否匹配（PHP switch 松散比较 ==）
// 优先按字符串比较，避免 StringValue 同时实现 AsInt 时把 "false"/"true" 都当成 0。
func (s *SwitchStatement) isMatch(value1, value2 data.GetValue) bool {
	_, isStr1 := value1.(*data.StringValue)
	_, isStr2 := value2.(*data.StringValue)
	if isStr1 && isStr2 {
		return value1.(*data.StringValue).Value == value2.(*data.StringValue).Value
	}

	if i1, ok := value1.(data.AsInt); ok {
		if i2, ok := value2.(data.AsInt); ok {
			n1, err1 := i1.AsInt()
			n2, err2 := i2.AsInt()
			if err1 == nil && err2 == nil {
				return n1 == n2
			}
		}
	}
	if strValue1, ok := value1.(data.AsString); ok {
		if strValue2, ok := value2.(data.AsString); ok {
			return strValue1.AsString() == strValue2.AsString()
		}
	}
	return false
}
