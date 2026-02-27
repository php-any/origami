package preg

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PregGrepFunction 实现 preg_grep 函数
//
// 目前实现的是最常用的行为（按值匹配，忽略键类型和更复杂的 flags）：
//
//	preg_grep(string $pattern, array $input, int $flags = 0): array|false
//
// - 使用 preg.Compile 统一处理 PHP 风格的正则表达式
// - 支持 PREG_GREP_INVERT：反转匹配结果
// - 输入必须是数组，否则返回 false
// - 返回值总是普通数组（索引从 0 开始），不完全保留原始键
type PregGrepFunction struct{}

func NewPregGrepFunction() data.FuncStmt {
	return &PregGrepFunction{}
}

func (f *PregGrepFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	patternValue, _ := ctx.GetIndexValue(0)
	inputValue, _ := ctx.GetIndexValue(1)
	flagsValue, _ := ctx.GetIndexValue(2)

	if patternValue == nil || inputValue == nil {
		return data.NewBoolValue(false), nil
	}

	pattern := patternValue.AsString()

	re, err := Compile(pattern)
	if err != nil {
		// PHP 行为: 发出 warning，返回 false；这里只返回 false
		return data.NewBoolValue(false), nil
	}

	// 解析 flags（目前仅支持 PREG_GREP_INVERT）
	flags := 0
	if flagsValue != nil {
		if asInt, ok := flagsValue.(data.AsInt); ok {
			if v, err := asInt.AsInt(); err == nil {
				flags = v
			}
		}
	}
	invert := flags&1 != 0 // PREG_GREP_INVERT

	// 只支持数组输入
	arr, ok := inputValue.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	values := arr.ToValueList()
	var result []data.Value

	for _, v := range values {
		subject := v.AsString()
		matched := re.MatchString(subject)

		if invert {
			if !matched {
				result = append(result, v)
			}
		} else {
			if matched {
				result = append(result, v)
			}
		}
	}

	return data.NewArrayValue(result), nil
}

func (f *PregGrepFunction) GetName() string {
	return "preg_grep"
}

func (f *PregGrepFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "pattern", 0, nil, nil),
		node.NewParameter(nil, "input", 1, nil, data.NewBaseType("array")),
		node.NewParameter(nil, "flags", 2, node.NewIntLiteral(nil, "0"), data.NewBaseType("int")),
	}
}

func (f *PregGrepFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "pattern", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "input", 1, data.NewBaseType("array")),
		node.NewVariable(nil, "flags", 2, data.NewBaseType("int")),
	}
}
