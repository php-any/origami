package preg

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PregQuoteFunction 实现 preg_quote(string $str, string $delimiter = null): string
// 将字符串中所有正则表达式特殊字符转义，可选额外转义 delimiter 字符。
type PregQuoteFunction struct{}

func NewPregQuoteFunction() data.FuncStmt {
	return &PregQuoteFunction{}
}

// PHP preg_quote 需要转义的特殊字符
const pregSpecialChars = `\.+*?[^]$(){}=!<>|:-#`

func (f *PregQuoteFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, ok := ctx.GetIndexValue(0)
	if !ok || strVal == nil {
		return data.NewStringValue(""), nil
	}
	str := strVal.AsString()

	var sb strings.Builder
	sb.Grow(len(str) * 2)

	// 读取可选的 delimiter
	delimiter := ""
	if delimVal, ok := ctx.GetIndexValue(1); ok && delimVal != nil {
		if _, isNull := delimVal.(*data.NullValue); !isNull {
			delimiter = delimVal.AsString()
		}
	}

	for _, ch := range str {
		s := string(ch)
		if strings.ContainsRune(pregSpecialChars, ch) {
			sb.WriteByte('\\')
		} else if delimiter != "" && strings.ContainsRune(delimiter, ch) {
			sb.WriteByte('\\')
		}
		sb.WriteString(s)
	}

	return data.NewStringValue(sb.String()), nil
}

func (f *PregQuoteFunction) GetName() string {
	return "preg_quote"
}

func (f *PregQuoteFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "str", 0, nil, nil),
		node.NewParameter(nil, "delimiter", 1, node.NewNullLiteral(nil), nil),
	}
}

func (f *PregQuoteFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "str", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "delimiter", 1, data.NewBaseType("string")),
	}
}
