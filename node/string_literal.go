package node

import (
	"strings"

	"github.com/php-any/origami/data"
)

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
		// 普通引号字符串：去掉字符串前后的引号
		if len(value) >= 2 {
			if value[0] == '"' || value[0] == '\'' {
				value = value[1 : len(value)-1]
			}
		}

		// 处理转义字符
		value = strings.ReplaceAll(value, "\\n", "\n")
		value = strings.ReplaceAll(value, "\\r", "\r")
		value = strings.ReplaceAll(value, "\\t", "\t")
		value = strings.ReplaceAll(value, "\\\"", "\"")
		value = strings.ReplaceAll(value, "\\'", "'")
		value = strings.ReplaceAll(value, "\\\\", "\\")
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
