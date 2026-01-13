package preg

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PregSplitFunction 实现 preg_split 函数
//
// preg_split(string $pattern, string $subject, int $limit = -1, int $flags = 0): array|false
//
// 根据正则表达式模式分割字符串
type PregSplitFunction struct{}

func NewPregSplitFunction() data.FuncStmt {
	return &PregSplitFunction{}
}

func (f *PregSplitFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	patternValue, _ := ctx.GetIndexValue(0)
	subjectValue, _ := ctx.GetIndexValue(1)
	limitValue, _ := ctx.GetIndexValue(2)
	flagsValue, _ := ctx.GetIndexValue(3)

	if patternValue == nil || subjectValue == nil {
		return data.NewBoolValue(false), nil
	}

	pattern := patternValue.AsString()
	subject := subjectValue.AsString()

	// 解析 limit 参数（默认为 -1，表示不限制）
	limit := -1
	if limitValue != nil {
		if intVal, ok := limitValue.(data.AsInt); ok {
			limit, _ = intVal.AsInt()
		}
	}

	// 解析 flags 参数（默认为 0）
	flags := 0
	if flagsValue != nil {
		if intVal, ok := flagsValue.(data.AsInt); ok {
			flags, _ = intVal.AsInt()
		}
	}

	// 编译正则表达式
	re, err := Compile(pattern)
	if err != nil {
		// PHP 行为: 发出 warning，返回 false；这里只返回 false
		return data.NewBoolValue(false), nil
	}

	// 检查 flags
	delimCapture := (flags & 2) != 0  // PREG_SPLIT_DELIM_CAPTURE
	offsetCapture := (flags & 4) != 0 // PREG_SPLIT_OFFSET_CAPTURE
	noEmpty := (flags & 1) != 0       // PREG_SPLIT_NO_EMPTY

	// 找到所有匹配位置
	allMatches := re.FindAllStringIndex(subject, -1)

	// 构建结果数组
	result := make([]data.Value, 0)
	currentOffset := 0
	splitCount := 0

	// 如果没有匹配，返回整个字符串（除非是空字符串且设置了 NO_EMPTY）
	if len(allMatches) == 0 {
		if !noEmpty || subject != "" {
			var value data.Value
			if offsetCapture {
				value = data.NewArrayValue([]data.Value{
					data.NewStringValue(subject),
					data.NewIntValue(0),
				})
			} else {
				value = data.NewStringValue(subject)
			}
			result = append(result, value)
		}
		return data.NewArrayValue(result), nil
	}

	// 处理每个匹配
	for _, match := range allMatches {
		// 获取当前匹配前的部分
		part := subject[currentOffset:match[0]]

		// 检查 limit（在添加元素之前检查）
		if limit > 0 && splitCount >= limit-1 {
			// 将剩余部分作为最后一个元素（包括当前部分和之后的所有内容）
			remaining := subject[currentOffset:]
			if !noEmpty || remaining != "" {
				var value data.Value
				if offsetCapture {
					value = data.NewArrayValue([]data.Value{
						data.NewStringValue(remaining),
						data.NewIntValue(currentOffset),
					})
				} else {
					value = data.NewStringValue(remaining)
				}
				result = append(result, value)
			}
			break
		}

		// 只有在不要求过滤空字符串，或者部分不为空时，才添加到结果
		if !noEmpty || part != "" {
			var value data.Value
			if offsetCapture {
				value = data.NewArrayValue([]data.Value{
					data.NewStringValue(part),
					data.NewIntValue(currentOffset),
				})
			} else {
				value = data.NewStringValue(part)
			}
			result = append(result, value)
			splitCount++
		}

		// 如果需要捕获分隔符
		if delimCapture {
			delim := subject[match[0]:match[1]]
			var delimValue data.Value
			if offsetCapture {
				delimValue = data.NewArrayValue([]data.Value{
					data.NewStringValue(delim),
					data.NewIntValue(match[0]),
				})
			} else {
				delimValue = data.NewStringValue(delim)
			}
			result = append(result, delimValue)
		}

		currentOffset = match[1]
	}

	// 处理最后一个匹配后的部分
	if limit <= 0 || splitCount < limit {
		remaining := subject[currentOffset:]
		if !noEmpty || remaining != "" {
			var value data.Value
			if offsetCapture {
				value = data.NewArrayValue([]data.Value{
					data.NewStringValue(remaining),
					data.NewIntValue(currentOffset),
				})
			} else {
				value = data.NewStringValue(remaining)
			}
			result = append(result, value)
		}
	}

	return data.NewArrayValue(result), nil
}

func (f *PregSplitFunction) GetName() string {
	return "preg_split"
}

func (f *PregSplitFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "pattern", 0, nil, nil),
		node.NewParameter(nil, "subject", 1, nil, nil),
		node.NewParameter(nil, "limit", 2, node.NewIntLiteral(nil, "-1"), nil),
		node.NewParameter(nil, "flags", 3, node.NewIntLiteral(nil, "0"), nil),
	}
}

func (f *PregSplitFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "pattern", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "subject", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "limit", 2, data.NewBaseType("int")),
		node.NewVariable(nil, "flags", 3, data.NewBaseType("int")),
	}
}
