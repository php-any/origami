package node

import (
	"strings"

	"github.com/php-any/origami/data"
)

// unescapeDoubleQuoted 按 PHP 双引号规则解析转义：\n \r \t \" \' \$ \\ \e \0-\377(八进制) \xHH(十六进制)
func unescapeDoubleQuoted(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' {
			b.WriteByte(s[i])
			continue
		}
		if i+1 >= len(s) {
			b.WriteByte('\\')
			continue
		}
		next := s[i+1]
		switch next {
		case '\\':
			b.WriteByte('\\')
			i++
		case 'n':
			b.WriteByte('\n')
			i++
		case 'r':
			b.WriteByte('\r')
			i++
		case 't':
			b.WriteByte('\t')
			i++
		case '"':
			b.WriteByte('"')
			i++
		case '\'':
			b.WriteByte('\'')
			i++
		case '$':
			b.WriteByte('$')
			i++
		case 'e':
			b.WriteByte(27) // ESC，与 \033 / \x1B 相同
			i++
		case '0', '1', '2', '3', '4', '5', '6', '7':
			// 八进制 1～3 位，从 s[i+1] 起读
			oct := 0
			start := i + 1
			k := 0
			for k < 3 && start+k < len(s) {
				c := s[start+k]
				if c < '0' || c > '7' {
					break
				}
				oct = oct*8 + int(c-'0')
				k++
			}
			b.WriteByte(byte(oct & 0xFF))
			i = start + k - 1 // 主循环会 i++，故下次从 start+k 开始
		case 'x', 'X':
			// 十六进制 1～2 位（PHP：\x1B 或 \x1），从 s[i+2] 起读
			hex := 0
			start := i + 2
			k := 0
			for k < 2 && start+k < len(s) {
				c := s[start+k]
				if c >= '0' && c <= '9' {
					hex = hex*16 + int(c-'0')
				} else if c >= 'a' && c <= 'f' {
					hex = hex*16 + int(c-'a'+10)
				} else if c >= 'A' && c <= 'F' {
					hex = hex*16 + int(c-'A'+10)
				} else {
					break
				}
				k++
			}
			b.WriteByte(byte(hex & 0xFF))
			i = start + k - 1
		default:
			b.WriteByte('\\')
		}
	}
	return b.String()
}

// unescapeSingleQuoted 单引号只处理 \\ 与 \'
func unescapeSingleQuoted(s string) string {
	s = strings.ReplaceAll(s, "\\\\", "\\")
	s = strings.ReplaceAll(s, "\\'", "'")
	return s
}

// StringLiteral 表示字符串字面量
type StringLiteral struct {
	*Node `pp:"-"`
	Value string
}

// NewStringLiteral 创建一个新的字符串字面量（仅处理引号字符串；heredoc/nowdoc 见 NewHeredocLiteral / NewNowdocLiteral）
func NewStringLiteral(token *TokenFrom, value string) data.GetValue {
	if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') {
		isDouble := value[0] == '"'
		value = value[1 : len(value)-1]
		if isDouble {
			value = unescapeDoubleQuoted(value)
		} else {
			value = unescapeSingleQuoted(value)
		}
	}
	return &StringLiteral{
		Node:  NewNode(token),
		Value: value,
	}
}

// NewHeredocLiteral 创建 heredoc 字符串字面量（正文已由 HeredocParser 提取）
func NewHeredocLiteral(token *TokenFrom, body string) data.GetValue {
	return &StringLiteral{
		Node:  NewNode(token),
		Value: body,
	}
}

// NewNowdocLiteral 创建 nowdoc 字符串字面量
func NewNowdocLiteral(token *TokenFrom, body string) data.GetValue {
	return &StringLiteral{
		Node:  NewNode(token),
		Value: body,
	}
}

// GetValue 获取字符串字面量的值
func (s *StringLiteral) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewStringValue(s.Value), nil
}

// NewStringLiteralByAst 不能转义的字符串
func NewStringLiteralByAst(token *TokenFrom, value string) data.GetValue {
	return &StringLiteral{
		Node:  NewNode(token),
		Value: value,
	}
}
