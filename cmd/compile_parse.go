package cmd

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
)

// ParsedFile 表示一个已解析的 PHP 文件
type ParsedFile struct {
	Path      string
	Program   *node.Program
	Variables []data.Variable
}

// parseFiles 批量解析 PHP 文件为 AST
func parseFiles(files []string) ([]ParsedFile, []error) {
	p := parser.NewParser()
	var parsed []ParsedFile
	var errs []error

	for _, file := range files {
		clone := p.Clone()
		program, acl := clone.ParseFile(file)
		if acl != nil {
			errs = append(errs, fmt.Errorf("解析 %s 失败: %v", file, acl))
			continue
		}
		vars := clone.GetVariables()
		parsed = append(parsed, ParsedFile{
			Path:      file,
			Program:   program,
			Variables: vars,
		})
	}
	return parsed, errs
}
