package preg

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PregReplaceFunction 实现 preg_replace
//
// 完整签名：
//
//	preg_replace(
//	  string|array $pattern,
//	  string|array $replacement,
//	  string|array $subject,
//	  int $limit = -1,
//	  int &$count = null
//	): string|array|null
//
// - 支持 pattern / replacement / subject 为字符串或一维数组
// - 支持 $limit（每个 pattern / subject 各自独立的限制，兼容 PHP 文档）
// - 使用 Compile 兼容 PHP 风格正则
type PregReplaceFunction struct{}

func NewPregReplaceFunction() data.FuncStmt {
	return &PregReplaceFunction{}
}

func (f *PregReplaceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	patternVal, _ := ctx.GetIndexValue(0)
	replVal, _ := ctx.GetIndexValue(1)
	subjectVal, _ := ctx.GetIndexValue(2)
	limitVal, _ := ctx.GetIndexValue(3)
	countVal, _ := ctx.GetIndexValue(4)

	if patternVal == nil || replVal == nil || subjectVal == nil {
		return data.NewBoolValue(false), nil
	}

	// 归一化 subject 为切片
	subjects, subjectIsArray := toStringSlice(subjectVal)

	// 归一化 pattern / replacement
	patterns, _ := toStringSlice(patternVal)
	repls, replIsArray := toStringSlice(replVal)

	// 解析 limit：默认 -1；<=0 视为无限制（与 PHP 行为保持一致）
	limit := -1
	if limitVal != nil {
		if _, isNull := limitVal.(*data.NullValue); !isNull {
			if asInt, ok := limitVal.(data.AsInt); ok {
				if v, err := asInt.AsInt(); err == nil {
					limit = v
				}
			}
		}
	}

	totalCount := 0
	results := make([]data.Value, 0, len(subjects))

	for _, subj := range subjects {
		replaced := subj
		localCount := 0

		for i, pat := range patterns {
			if pat == "" {
				// 跳过空 pattern，与 PHP 行为接近：不做任何替换
				continue
			}

			re, err := CompileAny(pat)
			if err != nil {
				// 编译错误：整体返回 false
				return data.NewBoolValue(false), nil
			}

			// 选择对应的 replacement
			var r string
			if replIsArray {
				if i < len(repls) {
					r = repls[i]
				} else {
					r = ""
				}
			} else {
				r = replVal.AsString()
			}

			// 处理 limit：<=0 表示无限制；>0 表示每个 pattern / subject 的最大替换次数
			if limit <= 0 {
				// 直接统计全部匹配次数
				matches := re.FindAllStringIndex(replaced, -1)
				if len(matches) == 0 {
					continue
				}
				localCount += len(matches)
				replaced = re.ReplaceAllString(replaced, r)
			} else {
				remaining := limit
				// 用 ReplaceAllStringFunc 控制替换次数
				replaced = re.ReplaceAllStringFunc(replaced, func(match string) string {
					if remaining > 0 {
						remaining--
						localCount++
						return r
					}
					// 超出 limit，保持原样
					return match
				})
			}
		}

		totalCount += localCount
		results = append(results, data.NewStringValue(replaced))
	}

	// 更新 $count（引用参数）
	if countVal != nil {
		if _, ok := countVal.(*data.NullValue); !ok {
			if ref, ok := countVal.(*data.ReferenceValue); ok {
				ref.Ctx.SetVariableValue(ref.Val, data.NewIntValue(totalCount))
			}
		}
	}

	if subjectIsArray {
		return data.NewArrayValue(results), nil
	}
	if len(results) == 0 {
		return data.NewStringValue(""), nil
	}
	return results[0], nil
}

func (f *PregReplaceFunction) GetName() string {
	return "preg_replace"
}

func (f *PregReplaceFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "pattern", 0, nil, nil),
		node.NewParameter(nil, "replacement", 1, nil, nil),
		node.NewParameter(nil, "subject", 2, nil, nil),
		node.NewParameter(nil, "limit", 3, node.NewIntLiteral(nil, "-1"), nil),
		node.NewParameter(nil, "count", 4, node.NewNullLiteral(nil), nil),
	}
}

func (f *PregReplaceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "pattern", 0, data.NewBaseType("mixed")),
		node.NewVariable(nil, "replacement", 1, data.NewBaseType("mixed")),
		node.NewVariable(nil, "subject", 2, data.NewBaseType("mixed")),
		node.NewVariable(nil, "limit", 3, data.NewBaseType("int")),
		node.NewVariable(nil, "count", 4, data.NewBaseType("int")),
	}
}

// toStringSlice 将 Value 归一化为字符串切片和“是否原本是数组”的标记。
func toStringSlice(v data.Value) ([]string, bool) {
	if arr, ok := v.(*data.ArrayValue); ok {
		vals := arr.ToValueList()
		out := make([]string, 0, len(vals))
		for _, item := range vals {
			out = append(out, item.AsString())
		}
		return out, true
	}
	return []string{v.AsString()}, false
}
