package parser

import (
	"github.com/php-any/origami/data"
	"os"
	"path/filepath"

	"github.com/php-any/origami/node"
)

type DirParser struct {
	*Parser
}

func NewDirParser(parser *Parser) StatementParser {
	return &DirParser{
		Parser: parser,
	}
}

func (p *DirParser) Parse() (data.GetValue, data.Control) {
	from := p.FromCurrentToken()

	// 获取当前文件的目录路径
	var dirPath string

	if p.source != nil {
		// 获取文件的绝对路径
		absPath, err := filepath.Abs(*p.source)
		if err == nil {
			dirPath = filepath.Dir(absPath)
		} else {
			// 如果获取绝对路径失败，使用当前工作目录
			if cwd, err := os.Getwd(); err == nil {
				dirPath = cwd
			} else {
				dirPath = "."
			}
		}
	} else {
		// 如果没有当前文件信息，使用当前工作目录
		if cwd, err := os.Getwd(); err == nil {
			dirPath = cwd
		} else {
			dirPath = "."
		}
	}

	// 移动到下一个 token
	p.next()

	// 返回目录路径的字符串字面量
	return node.NewStringLiteralByAst(from, dirPath), nil
}
