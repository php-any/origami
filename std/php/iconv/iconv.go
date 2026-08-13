package iconv

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IconvFunction 实现 iconv 函数
// 当前实现做一个尽量兼容的近似：
// - 规范化编码名（忽略大小写、前后空格、去掉 //TRANSLIT 等附加标记）
// - 如果 from == to，直接返回原字符串
// - 如果目标编码是 UTF-8 / UTF8，则假定传入字符串已是 UTF-8，直接返回
// - 其他编码组合暂不做真实重编码，失败时返回 false
//
// 这样可以覆盖大量只在 UTF-8 环境下运行、主要用 iconv 做“是否支持/占位”的场景，
// 同时尽量避免引入额外的第三方依赖。
type IconvFunction struct{}

func NewIconvFunction() data.FuncStmt {
	return &IconvFunction{}
}

func (f *IconvFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	fromVal, _ := ctx.GetIndexValue(0)
	toVal, _ := ctx.GetIndexValue(1)
	strVal, _ := ctx.GetIndexValue(2)

	if fromVal == nil || toVal == nil || strVal == nil {
		return data.NewBoolValue(false), nil
	}

	from := normalizeEncodingName(fromVal.AsString())
	to := normalizeEncodingName(toVal.AsString())
	s := strVal.AsString()

	// 完全相同的编码名，直接返回
	if from == to {
		return data.NewStringValue(s), nil
	}

	// 目标是 UTF-8：在 Origami 中字符串本来就是 UTF-8，直接返回
	if isUtf8Encoding(to) {
		return data.NewStringValue(s), nil
	}

	// 其他编码暂不支持，返回 false
	return data.NewBoolValue(false), nil
}

func (f *IconvFunction) GetName() string {
	return "iconv"
}

func (f *IconvFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "from_encoding", 0, nil, nil),
		node.NewParameter(nil, "to_encoding", 1, nil, nil),
		node.NewParameter(nil, "string", 2, nil, nil),
	}
}

func (f *IconvFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "from_encoding", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "to_encoding", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "string", 2, data.NewBaseType("string")),
	}
}

// normalizeEncodingName 规范化编码名称，去除大小写差异和 //TRANSLIT/IGNORE 等后缀
func normalizeEncodingName(name string) string {
	name = strings.TrimSpace(name)
	upper := strings.ToUpper(name)

	// 去掉 //TRANSLIT、//IGNORE 等标记
	if idx := strings.Index(upper, "//"); idx != -1 {
		upper = upper[:idx]
	}

	return strings.TrimSpace(upper)
}

func isUtf8Encoding(name string) bool {
	n := normalizeEncodingName(name)
	return n == "UTF-8" || n == "UTF8"
}
