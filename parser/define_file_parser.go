package parser

import (
	"path/filepath"

	"github.com/php-any/origami/data"

	"github.com/php-any/origami/node"
)

// FileParser 解析 __FILE__ 魔术常量
type FileParser struct {
	*Parser
}

// NewFileParser 创建一个新的 FileParser
func NewFileParser(p *Parser) StatementParser {
	return &FileParser{
		Parser: p,
	}
}

// Parse 返回当前文件的绝对路径字符串
func (p *FileParser) Parse() (data.GetValue, data.Control) {
	from := p.FromCurrentToken()

	var filePath string

	if p.source != nil {
		// 将当前源文件转换为绝对路径
		if absPath, err := filepath.Abs(*p.source); err == nil {
			filePath = absPath
		} else {
			// 获取绝对路径失败时，退化为原始 source 路径
			filePath = *p.source
		}
	} else {
		// 没有 source 信息时，退化为空字符串（与 PHP 在 CLI/EVAL 中的行为类似）
		filePath = ""
	}

	// 移动到下一个 token
	p.next()

	// 返回文件路径的字符串字面量
	return node.NewStringLiteralByAst(from, filePath), nil
}
