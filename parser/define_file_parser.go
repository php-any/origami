package parser

import (
	"github.com/php-any/origami/data"
	"path/filepath"

	"github.com/php-any/origami/node"
)

type FileParser struct {
	*Parser
}

func NewFileParser(parser *Parser) StatementParser {
	return &FileParser{
		Parser: parser,
	}
}

func (p *FileParser) Parse() (data.GetValue, data.Control) {
	from := p.NewTokenFrom(p.current().Start)

	// 获取当前文件的目录路径
	var dirPath string

	if p.source != nil {
		// 获取文件的绝对路径
		absPath, err := filepath.Abs(*p.source)
		if err == nil {
			dirPath = filepath.Dir(absPath)
		}
	}

	// 移动到下一个 token
	p.next()

	// 返回目录路径的字符串字面量
	return node.NewStringLiteralByAst(from, dirPath), nil
}
