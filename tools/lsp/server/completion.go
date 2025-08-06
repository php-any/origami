package server

import (
	"strings"
)

// handleTextDocumentCompletion 处理代码补全请求
// 当用户触发代码补全时（通常是按 Ctrl+Space），客户端会发送此请求
// 服务器需要根据当前光标位置和上下文返回合适的补全项
func (s *Server) handleTextDocumentCompletion(request map[string]interface{}) error {
	// 获取请求 ID，用于发送响应
	id := request["id"]

	// 解析请求参数
	params, ok := request["params"].(map[string]interface{})
	if !ok {
		// 参数解析失败，返回空的补全结果
		return s.sendResponse(id, map[string]interface{}{
			"isIncomplete": false,
			"items":        []map[string]interface{}{},
		})
	}

	// 获取光标位置信息
	position, ok := params["position"].(map[string]interface{})
	if !ok {
		// 位置信息解析失败，返回空的补全结果
		return s.sendResponse(id, map[string]interface{}{
			"isIncomplete": false,
			"items":        []map[string]interface{}{},
		})
	}

	// 提取行号和列号（LSP 使用 0 基索引）
	line, _ := position["line"].(float64)
	character, _ := position["character"].(float64)

	// 获取文档信息
	textDocument, ok := params["textDocument"].(map[string]interface{})
	if !ok {
		// 文档信息解析失败，返回空的补全结果
		return s.sendResponse(id, map[string]interface{}{
			"isIncomplete": false,
			"items":        []map[string]interface{}{},
		})
	}

	// 提取文档 URI
	uri, _ := textDocument["uri"].(string)

	// 获取当前输入的前缀
	// 这用于过滤补全项，只显示与前缀匹配的项目
	prefix := s.getCompletionPrefix(uri, int(line), int(character))

	// 定义所有可用的代码补全项
	allCompletionItems := s.getAllCompletionItems()

	// 根据前缀过滤补全项
	filteredItems := s.filterCompletionItems(allCompletionItems, prefix)

	// 构建响应结果
	result := map[string]interface{}{
		"isIncomplete": false,         // 表示这是完整的补全列表
		"items":        filteredItems, // 过滤后的补全项
	}

	// 发送响应给客户端
	return s.sendResponse(id, result)
}

// getAllCompletionItems 获取所有可用的代码补全项
// 返回 Origami 语言支持的所有语法结构和关键字的补全项
func (s *Server) getAllCompletionItems() []map[string]interface{} {
	return []map[string]interface{}{
		// echo 语句 - 用于输出内容
		{
			"label":            "echo",               // 显示标签
			"kind":             3,                    // 补全项类型：3 = Function
			"detail":           "echo statement",     // 详细描述
			"documentation":    "输出内容到标准输出",          // 文档说明
			"insertText":       "echo ${1:content};", // 插入的代码模板
			"insertTextFormat": 2,                    // 插入格式：2 = Snippet（支持占位符）
		},
		// if 条件语句
		{
			"label":            "if",
			"kind":             15, // 15 = Snippet
			"detail":           "if statement",
			"documentation":    "条件语句",
			"insertText":       "if (${1:condition}) {\n\t${2:// code}\n}",
			"insertTextFormat": 2,
		},
		// if-else 条件语句
		{
			"label":            "if-else",
			"kind":             15,
			"detail":           "if-else statement",
			"documentation":    "条件语句（带 else 分支）",
			"insertText":       "if (${1:condition}) {\n\t${2:// code}\n} else {\n\t${3:// code}\n}",
			"insertTextFormat": 2,
		},
		// for 循环语句
		{
			"label":            "for",
			"kind":             15,
			"detail":           "for loop",
			"documentation":    "for 循环语句",
			"insertText":       "for (${1:i} = ${2:0}; ${1:i} < ${3:length}; ${1:i}++) {\n\t${4:// code}\n}",
			"insertTextFormat": 2,
		},
		// foreach 循环语句
		{
			"label":            "foreach",
			"kind":             15,
			"detail":           "foreach loop",
			"documentation":    "foreach 循环语句，用于遍历数组或集合",
			"insertText":       "foreach (${1:array} as ${2:value}) {\n\t${3:// code}\n}",
			"insertTextFormat": 2,
		},
		// while 循环语句
		{
			"label":            "while",
			"kind":             15,
			"detail":           "while loop",
			"documentation":    "while 循环语句",
			"insertText":       "while (${1:condition}) {\n\t${2:// code}\n}",
			"insertTextFormat": 2,
		},
		// function 函数定义
		{
			"label":            "function",
			"kind":             15,
			"detail":           "function declaration",
			"documentation":    "定义一个函数",
			"insertText":       "function ${1:name}(${2:params}) {\n\t${3:// code}\n\treturn ${4:value};\n}",
			"insertTextFormat": 2,
		},
		// class 类定义
		{
			"label":            "class",
			"kind":             15,
			"detail":           "class declaration",
			"documentation":    "定义一个类",
			"insertText":       "class ${1:ClassName} {\n\t${2:// properties and methods}\n}",
			"insertTextFormat": 2,
		},
		// try-catch 异常处理
		{
			"label":            "try",
			"kind":             15,
			"detail":           "try-catch statement",
			"documentation":    "异常处理语句",
			"insertText":       "try {\n\t${1:// code}\n} catch (${2:Exception} ${3:e}) {\n\t${4:// handle exception}\n}",
			"insertTextFormat": 2,
		},
		// switch 分支语句
		{
			"label":            "switch",
			"kind":             15,
			"detail":           "switch statement",
			"documentation":    "switch 分支语句",
			"insertText":       "switch (${1:variable}) {\n\tcase ${2:value1}:\n\t\t${3:// code}\n\t\tbreak;\n\tcase ${4:value2}:\n\t\t${5:// code}\n\t\tbreak;\n\tdefault:\n\t\t${6:// code}\n}",
			"insertTextFormat": 2,
		},
	}
}

// filterCompletionItems 根据前缀过滤补全项
// 只返回标签以指定前缀开头的补全项（不区分大小写）
func (s *Server) filterCompletionItems(items []map[string]interface{}, prefix string) []map[string]interface{} {
	// 如果没有前缀，返回所有补全项
	if prefix == "" {
		return items
	}

	// 创建过滤后的结果数组
	filteredItems := make([]map[string]interface{}, 0)

	// 遍历所有补全项
	for _, item := range items {
		// 获取补全项的标签
		label, ok := item["label"].(string)
		if !ok {
			continue // 跳过无效的标签
		}

		// 检查标签是否以前缀开头（不区分大小写）
		if strings.HasPrefix(strings.ToLower(label), strings.ToLower(prefix)) {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

// getCompletionPrefix 获取当前输入的前缀
// 从光标位置向前查找，提取当前正在输入的单词
func (s *Server) getCompletionPrefix(uri string, line, character int) string {
	// 从文档缓存中获取内容
	content, exists := s.documents[uri]
	if !exists {
		return "" // 文档不存在，返回空前缀
	}

	// 按行分割文档内容
	lines := strings.Split(content, "\n")
	if line >= len(lines) {
		return "" // 行号超出范围，返回空前缀
	}

	// 获取当前行的内容
	currentLine := lines[line]
	if character > len(currentLine) {
		character = len(currentLine) // 确保字符位置不超出行长度
	}

	// 从光标位置向前查找单词边界
	start := character
	for start > 0 {
		char := currentLine[start-1]
		// 检查字符是否为标识符的一部分（字母、数字或下划线）
		if !isIdentifierChar(char) {
			break // 遇到非标识符字符，停止查找
		}
		start--
	}

	// 提取前缀
	if start < character {
		return currentLine[start:character]
	}

	return "" // 没有找到前缀
}

// isIdentifierChar 检查字符是否为标识符的一部分
// 标识符可以包含字母、数字和下划线
func isIdentifierChar(char byte) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		(char >= '0' && char <= '9') ||
		char == '_'
}
