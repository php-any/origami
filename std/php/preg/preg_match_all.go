package preg

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PregMatchAllFunction 实现 preg_match_all 函数
//
// 目前实现的是最常用的行为（PREG_PATTERN_ORDER，忽略 flags 和 offset）:
//
//	preg_match_all(string $pattern, string $subject, array &$matches = null, int $flags = 0, int $offset = 0): int|false
//
// $matches 的结构为:
//   - $matches[0]: 所有完整匹配
//   - $matches[1]: 第 1 个捕获分组的所有匹配
//   - ...
type PregMatchAllFunction struct{}

func NewPregMatchAllFunction() data.FuncStmt {
	return &PregMatchAllFunction{}
}

func (f *PregMatchAllFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	patternValue, _ := ctx.GetIndexValue(0)
	subjectValue, _ := ctx.GetIndexValue(1)
	matchesValue, _ := ctx.GetIndexValue(2)
	// flagsValue, _ := ctx.GetIndexValue(3)  // 目前忽略
	// offsetValue, _ := ctx.GetIndexValue(4) // 目前忽略

	if patternValue == nil || subjectValue == nil {
		return data.NewBoolValue(false), nil
	}

	pattern := patternValue.AsString()
	subject := subjectValue.AsString()

	re, err := Compile(pattern)
	if err != nil {
		// PHP 行为: 发出 warning，返回 false；这里只返回 false
		return data.NewBoolValue(false), nil
	}

	// 查找所有匹配以及分组
	allLocs := re.FindAllStringSubmatchIndex(subject, -1)
	if len(allLocs) == 0 {
		// 无匹配时返回 0，并将 $matches 设为空数组
		if matchesValue != nil {
			empty := data.NewArrayValue([]data.Value{})
			if r, ok := matchesValue.(*data.ReferenceValue); ok {
				r.Ctx.SetVariableValue(r.Val, empty)
			} else if arr, ok := matchesValue.(*data.ArrayValue); ok {
				arr.Value = []data.Value{}
			}
		}
		return data.NewIntValue(0), nil
	}

	// 每个匹配 loc: [start0,end0, start1,end1, ...]
	groupCount := len(allLocs[0]) / 2
	groups := make([]data.Value, groupCount)

	for g := 0; g < groupCount; g++ {
		var groupMatches []data.Value
		for _, loc := range allLocs {
			start := loc[g*2]
			end := loc[g*2+1]
			if start == -1 || end == -1 {
				groupMatches = append(groupMatches, data.NewStringValue(""))
			} else {
				groupMatches = append(groupMatches, data.NewStringValue(subject[start:end]))
			}
		}
		groups[g] = data.NewArrayValue(groupMatches)
	}

	if matchesValue != nil {
		newMatches := data.NewArrayValue(groups)
		if r, ok := matchesValue.(*data.ReferenceValue); ok {
			r.Ctx.SetVariableValue(r.Val, newMatches)
		} else if arr, ok := matchesValue.(*data.ArrayValue); ok {
			arr.Value = groups
		}
	}

	// 返回匹配到的次数
	return data.NewIntValue(len(allLocs)), nil
}

func (f *PregMatchAllFunction) GetName() string {
	return "preg_match_all"
}

func (f *PregMatchAllFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "pattern", 0, nil, nil),
		node.NewParameter(nil, "subject", 1, nil, nil),
		node.NewParameter(nil, "matches", 2, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "flags", 3, node.NewIntLiteral(nil, "0"), nil),
		node.NewParameter(nil, "offset", 4, node.NewIntLiteral(nil, "0"), nil),
	}
}

func (f *PregMatchAllFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "pattern", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "subject", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "matches", 2, data.NewBaseType("array")),
		node.NewVariable(nil, "flags", 3, data.NewBaseType("int")),
		node.NewVariable(nil, "offset", 4, data.NewBaseType("int")),
	}
}
