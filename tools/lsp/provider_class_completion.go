package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/sirupsen/logrus"
)

// GetClassCompletionsForContext 获取上下文相关的类补全（use导入 + 同级目录）
// worker: 用于过滤类名，只返回包含 worker 字母的类，如果为空则返回空列表
func (p *LSPSymbolProvider) GetClassCompletionsForContext(content string, position defines.Position, worker string) []defines.CompletionItem {
	// 如果 Worker 为空，不返回任何类
	if worker == "" {
		return []defines.CompletionItem{}
	}

	workerLower := strings.ToLower(worker)
	itemsMap := make(map[string]defines.CompletionItem)

	// 1. 收集当前文件中的 use 语句
	// 使用 p.doc.Foreach 来识别节点
	if p.doc != nil {
		p.doc.Foreach(func(ctx *LspContext, parent, child data.GetValue) bool {
			// 检查是否是 use 语句
			if useStmt, ok := child.(*node.UseStatement); ok {
				// 提取类名（use 语句的最后一部分）
				namespace := useStmt.Namespace
				parts := strings.Split(namespace, "\\")
				className := parts[len(parts)-1]

				// 使用别名（如果有的话）
				if useStmt.Alias != "" {
					className = useStmt.Alias
				}

				// 只添加包含 worker 字母的类名（不区分大小写）
				classNameLower := strings.ToLower(className)
				if strings.Contains(classNameLower, workerLower) {
					itemsMap[className] = defines.CompletionItem{
						Label:  className,
						Kind:   &[]defines.CompletionItemKind{defines.CompletionItemKindClass}[0],
						Detail: &[]string{"use " + namespace}[0],
					}
				}
			}
			return true // 继续遍历
		})
	}

	// 2. 收集当前文件同级目录下的所有类
	// 尝试从文档信息中获取文件路径
	var currentFilePath string
	if p.doc != nil && p.doc.AST != nil {
		// 尝试从 AST 的第一个语句获取 source
		for _, stmt := range p.doc.AST.Statements {
			if getFrom, ok := stmt.(node.GetFrom); ok {
				if from := getFrom.GetFrom(); from != nil {
					currentFilePath = from.GetSource()
					break
				}
			}
		}
	}

	if currentFilePath != "" {
		// 获取当前文件所在目录
		currentDir := filepath.Dir(currentFilePath)

		// 扫描目录中的所有 .zy 和 .php 文件
		files, err := os.ReadDir(currentDir)
		if err == nil {
			for _, file := range files {
				if file.IsDir() {
					continue
				}

				fileName := file.Name()
				if !strings.HasSuffix(fileName, ".zy") && !strings.HasSuffix(fileName, ".php") {
					continue
				}

				// 读取文件并解析类名
				filePath := filepath.Join(currentDir, fileName)
				classes := p.extractClassNamesFromFile(filePath)

				for _, className := range classes {
					// 只添加包含 worker 字母的类名（不区分大小写）
					classNameLower := strings.ToLower(className)
					if strings.Contains(classNameLower, workerLower) {
						if _, exists := itemsMap[className]; !exists {
							itemsMap[className] = defines.CompletionItem{
								Label:  className,
								Kind:   &[]defines.CompletionItemKind{defines.CompletionItemKindClass}[0],
								Detail: &[]string{"同包类"}[0],
							}
						}
					}
				}
			}
		}
	}

	// 3. 如果从 VM 中可以获取全局类，也添加进来（使用短名展示，Detail 中保存完整类名）
	if p.vm != nil {
		allClasses := p.vm.GetAllClasses()
		for fullName := range allClasses {
			// 提取短名：命名空间最后一段
			shortName := fullName
			if idx := strings.LastIndex(fullName, "\\"); idx >= 0 && idx+1 < len(fullName) {
				shortName = fullName[idx+1:]
			}

			// 只按短名过滤（不区分大小写）
			shortLower := strings.ToLower(shortName)
			if !strings.Contains(shortLower, workerLower) {
				continue
			}

			if _, exists := itemsMap[shortName]; !exists {
				detail := "full:" + fullName
				itemsMap[shortName] = defines.CompletionItem{
					Label:  shortName,
					Kind:   &[]defines.CompletionItemKind{defines.CompletionItemKindClass}[0],
					Detail: &detail,
				}
			}
		}
	}

	// 转换为数组
	items := make([]defines.CompletionItem, 0, len(itemsMap))
	for _, item := range itemsMap {
		items = append(items, item)
	}

	logrus.Debugf("GetClassCompletionsForContext: Worker=%s, 找到 %d 个类", worker, len(items))
	return items
}

// extractClassNamesFromFile 从文件中提取所有类名
func (p *LSPSymbolProvider) extractClassNamesFromFile(filePath string) []string {
	var classNames []string

	// 创建一个临时解析器来解析文件
	parser := NewLspParser()
	if p.vm != nil {
		parser.SetVM(p.vm)
	}

	ast, err := parser.ParseFile(filePath)
	if err != nil || ast == nil {
		return classNames
	}

	// 遍历 AST 查找类声明
	for _, stmt := range ast.Statements {
		if classStmt, ok := stmt.(*node.ClassStatement); ok {
			classNames = append(classNames, classStmt.Name)
		}
	}

	return classNames
}
