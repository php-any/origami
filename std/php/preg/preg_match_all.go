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
	flagsValue, _ := ctx.GetIndexValue(3)
	// offsetValue, _ := ctx.GetIndexValue(4) // 目前忽略

	if patternValue == nil || subjectValue == nil {
		return data.NewBoolValue(false), nil
	}

	pattern := patternValue.AsString()
	subject := subjectValue.AsString()

	re, err := CompileAny(pattern)
	if err != nil {
		// PHP 行为: 发出 warning，返回 false；这里只返回 false
		return data.NewBoolValue(false), nil
	}

	// 解析 flags（支持 PREG_PATTERN_ORDER / PREG_SET_ORDER、PREG_OFFSET_CAPTURE）
	flags := 0
	if flagsValue != nil {
		if asInt, ok := flagsValue.(data.AsInt); ok {
			if v, err := asInt.AsInt(); err == nil {
				flags = v
			}
		}
	}
	patternOrder := flags&2 == 0 // !(flags & PREG_SET_ORDER)
	offsetCapture := flags&256 != 0

	// 查找所有匹配以及分组（索引数组）
	allLocs := re.FindAllStringSubmatchIndex(subject, -1)
	if len(allLocs) == 0 {
		// 无匹配时返回 0，并将 $matches 设为空数组
		var empty data.Value
		if patternOrder {
			// PREG_PATTERN_ORDER: 为所有捕获组创建空数组
			// 需要确定捕获组数量：尝试对空字符串匹配获取 group count
			groupCount := 1 // 至少 group 0
			if goRe, goErr := Compile(pattern); goErr == nil {
				groupCount = goRe.NumSubexp() + 1 // NumSubexp 返回捕获组数量，+1 包含 group 0
			}
			groups := make([]data.Value, groupCount)
			for i := 0; i < groupCount; i++ {
				groups[i] = data.NewArrayValue([]data.Value{})
			}
			empty = data.NewArrayValue(groups)
		} else {
			// PREG_SET_ORDER: $matches = [] (no rows)
			empty = data.NewArrayValue([]data.Value{})
		}
		// 通过参数引用写回调用方的 $matches
		if z := ctx.GetIndexZVal(2); z != nil {
			z.Value = empty
		}
		return data.NewIntValue(0), nil
	}

	// loc: [start0,end0, start1,end1, ...]，第 0 组是完整匹配，后续是各捕获分组
	if patternOrder {
		// PREG_PATTERN_ORDER：$matches[group][matchIndex]
		groupCount := len(allLocs[0]) / 2
		groups := make([]data.Value, groupCount)

		for g := 0; g < groupCount; g++ {
			var groupMatches []data.Value
			for _, loc := range allLocs {
				start := loc[g*2]
				end := loc[g*2+1]

				// 未匹配的分组
				if start == -1 || end == -1 {
					if offsetCapture {
						// ['', 0]
						pair := []data.Value{
							data.NewStringValue(""),
							data.NewIntValue(0),
						}
						groupMatches = append(groupMatches, data.NewArrayValue(pair))
					} else {
						groupMatches = append(groupMatches, data.NewStringValue(""))
					}
					continue
				}

				text := data.NewStringValue(subject[start:end])
				if offsetCapture {
					pair := []data.Value{
						text,
						data.NewIntValue(start),
					}
					groupMatches = append(groupMatches, data.NewArrayValue(pair))
				} else {
					groupMatches = append(groupMatches, text)
				}
			}
			groups[g] = data.NewArrayValue(groupMatches)
		}

		matchesArr := data.NewArrayValue(groups)
		if z := ctx.GetIndexZVal(2); z != nil {
			z.Value = matchesArr
		}
	} else {
		// PREG_SET_ORDER：$matches[matchIndex][group]
		var rows []data.Value
		for _, loc := range allLocs {
			groupCount := len(loc) / 2
			var row []data.Value
			for g := 0; g < groupCount; g++ {
				start := loc[g*2]
				end := loc[g*2+1]
				if start == -1 || end == -1 {
					if offsetCapture {
						pair := []data.Value{
							data.NewStringValue(""),
							data.NewIntValue(0),
						}
						row = append(row, data.NewArrayValue(pair))
					} else {
						row = append(row, data.NewStringValue(""))
					}
					continue
				}
				text := data.NewStringValue(subject[start:end])
				if offsetCapture {
					pair := []data.Value{
						text,
						data.NewIntValue(start),
					}
					row = append(row, data.NewArrayValue(pair))
				} else {
					row = append(row, text)
				}
			}
			rows = append(rows, data.NewArrayValue(row))
		}
		matchesArr := data.NewArrayValue(rows)
		if z := ctx.GetIndexZVal(2); z != nil {
			z.Value = matchesArr
		}
	}

	// 返回匹配到的次数（完整匹配数量）
	return data.NewIntValue(len(allLocs)), nil
}

func (f *PregMatchAllFunction) GetName() string {
	return "preg_match_all"
}

func (f *PregMatchAllFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "pattern", 0, nil, nil),
		node.NewParameter(nil, "subject", 1, nil, nil),
		// 第三个参数为 &array $matches，按 PHP 语义需要按引用传递
		node.NewParameterReference(nil, "matches", 2, data.NewBaseType("array")),
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
