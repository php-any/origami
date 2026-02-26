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

// NewStringLiteral 创建一个新的字符串字面量
func NewStringLiteral(token *TokenFrom, value string) data.GetValue {
	// 检查是否是 heredoc/nowdoc
	if len(value) >= 3 && value[:3] == "<<<" {
		// heredoc/nowdoc 格式: <<<'IDENTIFIER'\n内容\nIDENTIFIER 或 <<<IDENTIFIER\n内容\nIDENTIFIER
		// 找到第一个换行符（内容开始）
		firstNewline := strings.IndexByte(value, '\n')
		if firstNewline == -1 {
			firstNewline = strings.IndexByte(value, '\r')
		}
		if firstNewline != -1 {
			// 找到最后一个换行符（结束标记前）
			lastNewline := strings.LastIndexByte(value, '\n')
			if lastNewline == -1 {
				lastNewline = strings.LastIndexByte(value, '\r')
			}
			if lastNewline > firstNewline {
				// 提取内容部分（第一个换行符后到最后一个换行符前）
				value = value[firstNewline+1 : lastNewline]
			}
		}
	} else {
		// 普通引号字符串：区分双引号与单引号，再去掉首尾引号并解析转义
		if len(value) >= 2 && (value[0] == '"' || value[0] == '\'') {
			isDouble := value[0] == '"'
			value = value[1 : len(value)-1]
			if isDouble {
				value = unescapeDoubleQuoted(value) // 支持 \033 \x1B \e \n \t 等
			} else {
				value = unescapeSingleQuoted(value) // 仅 \\ 与 \'
			}
		}
	}

	return &StringLiteral{
		Node:  NewNode(token),
		Value: value,
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
