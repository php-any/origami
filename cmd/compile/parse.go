package compile

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
	"github.com/php-any/origami/runtime"
)

// ParsedFile 表示一个已解析的 PHP 文件
type ParsedFile struct {
	Path      string
	Program   *node.Program
	Variables []data.Variable
	Namespace string // 文件的命名空间
}

// parseFiles 批量解析 PHP 文件为 AST
func parseFiles(files []string) ([]ParsedFile, []error) {
	p := parser.NewParser()
	baseVM := runtime.NewVM(p)
	if runtimeLoader != nil {
		runtimeLoader(baseVM)
	}
	var parsed []ParsedFile
	var errs []error

	for _, file := range files {
		tempVM := runtime.NewTempVM(baseVM).(*runtime.TempVM)
		clone := p.Clone()
		clone.SetVM(tempVM)
		program, acl := clone.ParseFile(file)
		if acl != nil {
			errs = append(errs, fmt.Errorf("解析 %s 失败: %v", file, acl))
			continue
		}
		augmentProgramAST(program, tempVM, file)
		vars := clone.GetVariables()
		parsed = append(parsed, ParsedFile{
			Path:      file,
			Program:   program,
			Variables: vars,
			Namespace: clone.GetNamespace(),
		})
	}
	return parsed, errs
}

// augmentProgramAST 将解析时注册到 TempVM、但未写入 Program/Namespace 的类定义补回 AST。
func augmentProgramAST(program *node.Program, tempVM *runtime.TempVM, file string) {
	var classes []data.GetValue
	for _, c := range tempVM.AddedClasses() {
		if classSourceFile(c) != file {
			continue
		}
		if gv, ok := c.(data.GetValue); ok {
			classes = append(classes, gv)
		}
	}
	if len(classes) == 0 {
		return
	}
	for _, stmt := range program.Statements {
		ns, ok := stmt.(*node.Namespace)
		if !ok {
			continue
		}
		ns.Statements = append(filterNamespaceStmts(ns.Statements, ns.Name), classes...)
	}
}

func filterNamespaceStmts(stmts []data.GetValue, nsName string) []data.GetValue {
	out := make([]data.GetValue, 0, len(stmts))
	for _, s := range stmts {
		if lit, ok := s.(*node.StringLiteral); ok && lit.Value == nsName {
			continue
		}
		out = append(out, s)
	}
	return out
}

func classSourceFile(c data.ClassStmt) string {
	switch s := c.(type) {
	case *node.ClassStatement:
		if f := s.GetFrom(); f != nil {
			return f.GetSource()
		}
	case *node.AbstractClassStatement:
		if f := s.ClassStatement.GetFrom(); f != nil {
			return f.GetSource()
		}
	}
	return ""
}
