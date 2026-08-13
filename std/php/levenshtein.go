package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// LevenshteinFunction 实现 levenshtein 函数
//
// 对齐 PHP 语义的常用部分：
//
//	levenshtein(
//	  string $string1,
//	  string $string2,
//	  ?int $insertion_cost = 1,
//	  ?int $replacement_cost = 1,
//	  ?int $deletion_cost = 1
//	): int
//
// - 当任一字符串为 null 时，按空字符串处理。
// - 当任一 cost 参数为 null 或未提供时，使用默认值 1。
// - 当前实现不校验过大的字符串长度，也不实现返回 false 的错误分支，始终返回距离（int）。
type LevenshteinFunction struct{}

func NewLevenshteinFunction() data.FuncStmt {
	return &LevenshteinFunction{}
}

func (f *LevenshteinFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	s1Val, _ := ctx.GetIndexValue(0)
	s2Val, _ := ctx.GetIndexValue(1)
	insCostVal, _ := ctx.GetIndexValue(2)
	repCostVal, _ := ctx.GetIndexValue(3)
	delCostVal, _ := ctx.GetIndexValue(4)

	var s1, s2 string
	if s1Val != nil {
		if str, ok := s1Val.(data.AsString); ok {
			s1 = str.AsString()
		} else {
			s1 = s1Val.AsString()
		}
	}
	if s2Val != nil {
		if str, ok := s2Val.(data.AsString); ok {
			s2 = str.AsString()
		} else {
			s2 = s2Val.AsString()
		}
	}

	insertionCost := 1
	if insCostVal != nil {
		if iv, ok := insCostVal.(data.AsInt); ok {
			if v, err := iv.AsInt(); err == nil {
				insertionCost = v
			}
		}
	}
	replacementCost := 1
	if repCostVal != nil {
		if iv, ok := repCostVal.(data.AsInt); ok {
			if v, err := iv.AsInt(); err == nil {
				replacementCost = v
			}
		}
	}
	deletionCost := 1
	if delCostVal != nil {
		if iv, ok := delCostVal.(data.AsInt); ok {
			if v, err := iv.AsInt(); err == nil {
				deletionCost = v
			}
		}
	}

	d := levenshteinDistance(s1, s2, insertionCost, replacementCost, deletionCost)
	return data.NewIntValue(d), nil
}

func (f *LevenshteinFunction) GetName() string {
	return "levenshtein"
}

func (f *LevenshteinFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string1", 0, nil, nil),
		node.NewParameter(nil, "string2", 1, nil, nil),
		node.NewParameter(nil, "insertion_cost", 2, node.NewIntLiteral(nil, "1"), data.NewBaseType("int")),
		node.NewParameter(nil, "replacement_cost", 3, node.NewIntLiteral(nil, "1"), data.NewBaseType("int")),
		node.NewParameter(nil, "deletion_cost", 4, node.NewIntLiteral(nil, "1"), data.NewBaseType("int")),
	}
}

func (f *LevenshteinFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string1", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "string2", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "insertion_cost", 2, data.NewBaseType("int")),
		node.NewVariable(nil, "replacement_cost", 3, data.NewBaseType("int")),
		node.NewVariable(nil, "deletion_cost", 4, data.NewBaseType("int")),
	}
}

// levenshteinDistance 计算带权 Levenshtein 距离
func levenshteinDistance(s1, s2 string, insCost, repCost, delCost int) int {
	len1 := len(s1)
	len2 := len(s2)

	if len1 == 0 {
		return len2 * insCost
	}
	if len2 == 0 {
		return len1 * delCost
	}

	// 使用两行 DP，节省内存
	prev := make([]int, len2+1)
	curr := make([]int, len2+1)

	for j := 0; j <= len2; j++ {
		prev[j] = j * insCost
	}

	for i := 1; i <= len1; i++ {
		curr[0] = i * delCost
		c1 := s1[i-1]
		for j := 1; j <= len2; j++ {
			c2 := s2[j-1]
			cost := 0
			if c1 == c2 {
				cost = 0
			} else {
				cost = repCost
			}

			del := prev[j] + delCost
			ins := curr[j-1] + insCost
			rep := prev[j-1] + cost

			min := del
			if ins < min {
				min = ins
			}
			if rep < min {
				min = rep
			}
			curr[j] = min
		}
		prev, curr = curr, prev
	}

	return prev[len2]
}
