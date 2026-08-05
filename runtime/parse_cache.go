package runtime

import (
	"github.com/php-any/origami/data"
)

type parsedPHPFile struct {
	program data.GetValue
	vars    []data.Variable
}

// ParseFileCached 解析 PHP 文件并缓存 AST（进程级，按规范化路径去重）。
// 多次执行同一入口/被 include 的文件时不应重复读盘与词法/语法分析。
func (vm *VM) ParseFileCached(file string) (data.GetValue, []data.Variable, data.Control) {
	file = normalizePhpFilePath(file)
	if file == "" {
		return nil, nil, nil
	}
	if cached, ok := syncMapLoad[*parsedPHPFile](&vm.parsedFiles, file); ok {
		return cached.program, cached.vars, nil
	}

	p := vm.parser.Clone()
	program, acl := p.ParseFile(file)
	if acl != nil {
		return nil, nil, acl
	}
	entry := &parsedPHPFile{
		program: program,
		vars:    p.GetVariables(),
	}
	actual, _ := vm.parsedFiles.LoadOrStore(file, entry)
	cached := actual.(*parsedPHPFile)
	return cached.program, cached.vars, nil
}

// ClearParsedFileCache 清空已缓存的 AST，供热重载使用。
func (vm *VM) ClearParsedFileCache() {
	syncMapClear(&vm.parsedFiles)
}
