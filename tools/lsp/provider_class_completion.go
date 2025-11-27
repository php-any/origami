package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/php-any/origami/node"
	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/sirupsen/logrus"
)

// GetClassCompletionsForContext 获取上下文相关的类补全（use导入 + 同级目录）
func (p *LSPSymbolProvider) GetClassCompletionsForContext(content string, position defines.Position) []defines.CompletionItem {
	itemsMap := make(map[string]defines.CompletionItem)

	// 1. 收集当前文件中的 use 语句
	if p.doc != nil && p.doc.AST != nil {
		for _, stmt := range p.doc.AST.Statements {
			if useStmt, ok := stmt.(*node.UseStatement); ok {
				// 提取类名（use 语句的最后一部分）
				namespace := useStmt.Namespace
				parts := strings.Split(namespace, "\\")
				className := parts[len(parts)-1]

				// 使用别名（如果有的话）
				if useStmt.Alias != "" {
					className = useStmt.Alias
				}

				itemsMap[className] = defines.CompletionItem{
					Label:  className,
					Kind:   &[]defines.CompletionItemKind{defines.CompletionItemKindClass}[0],
					Detail: &[]string{"use " + namespace}[0],
				}
			}
		}
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
					if _, exists := itemsMap[className]; !exists {
						itemsMap[className] = defines.CompletionItem{
							Label:  className,
							Kind:   &[]defines.CompletionItemKind{defines.CompletionItemKindClass}[0],
							Detail: &[]string{"同级目录"}[0],
						}
					}
				}
			}
		}
	}

	// 3. 如果从 VM 中可以获取全局类，也添加进来
	if p.vm != nil {
		allClasses := p.vm.GetAllClasses()
		for className := range allClasses {
			if _, exists := itemsMap[className]; !exists {
				itemsMap[className] = defines.CompletionItem{
					Label:  className,
					Kind:   &[]defines.CompletionItemKind{defines.CompletionItemKindClass}[0],
					Detail: &[]string{"全局类"}[0],
				}
			}
		}
	}

	// 转换为数组
	items := make([]defines.CompletionItem, 0, len(itemsMap))
	for _, item := range itemsMap {
		items = append(items, item)
	}

	logrus.Debugf("GetClassCompletionsForContext: 找到 %d 个类", len(items))
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
