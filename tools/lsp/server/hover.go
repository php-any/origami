package server

import (
	"strings"
)

// handleTextDocumentHover 处理悬停信息请求
// 当用户将鼠标悬停在代码上时，客户端会发送此请求
// 服务器需要返回相关的文档信息、类型信息或其他有用的提示
func (s *Server) handleTextDocumentHover(request map[string]interface{}) error {
	// 获取请求 ID，用于发送响应
	id := request["id"]

	// 解析请求参数
	params, ok := request["params"].(map[string]interface{})
	if !ok {
		// 参数解析失败，返回空的悬停结果
		return s.sendResponse(id, nil)
	}

	// 获取光标位置信息
	position, ok := params["position"].(map[string]interface{})
	if !ok {
		// 位置信息解析失败，返回空的悬停结果
		return s.sendResponse(id, nil)
	}

	// 提取行号和列号（LSP 使用 0 基索引）
	line, _ := position["line"].(float64)
	character, _ := position["character"].(float64)

	// 获取文档信息
	textDocument, ok := params["textDocument"].(map[string]interface{})
	if !ok {
		// 文档信息解析失败，返回空的悬停结果
		return s.sendResponse(id, nil)
	}

	// 提取文档 URI
	uri, _ := textDocument["uri"].(string)

	// 获取悬停位置的单词或符号
	word := s.getWordAtPosition(uri, int(line), int(character))
	if word == "" {
		// 没有找到单词，返回空的悬停结果
		return s.sendResponse(id, nil)
	}

	// 获取该单词的悬停信息
	hoverInfo := s.getHoverInfo(word)
	if hoverInfo == nil {
		// 没有找到相关信息，返回空的悬停结果
		return s.sendResponse(id, nil)
	}

	// 发送悬停信息响应
	return s.sendResponse(id, hoverInfo)
}

// getWordAtPosition 获取指定位置的单词
// 从光标位置提取完整的标识符或关键字
func (s *Server) getWordAtPosition(uri string, line, character int) string {
	// 从文档缓存中获取内容
	content, exists := s.documents[uri]
	if !exists {
		return "" // 文档不存在
	}

	// 按行分割文档内容
	lines := strings.Split(content, "\n")
	if line >= len(lines) {
		return "" // 行号超出范围
	}

	// 获取当前行的内容
	currentLine := lines[line]
	if character >= len(currentLine) {
		return "" // 字符位置超出行长度
	}

	// 向前查找单词开始位置
	start := character
	for start > 0 && isIdentifierChar(currentLine[start-1]) {
		start--
	}

	// 向后查找单词结束位置
	end := character
	for end < len(currentLine) && isIdentifierChar(currentLine[end]) {
		end++
	}

	// 提取单词
	if start < end {
		return currentLine[start:end]
	}

	return ""
}

// getHoverInfo 获取指定单词的悬停信息
// 返回包含文档说明、类型信息等的悬停内容
func (s *Server) getHoverInfo(word string) map[string]interface{} {
	// 定义 Origami 语言关键字和内置函数的文档信息
	hoverData := map[string]map[string]interface{}{
		// 控制流关键字
		"if": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**if** - 条件语句\n\n用于根据条件执行不同的代码分支。\n\n```origami\nif (condition) {\n    // 当条件为真时执行\n}\n```",
			},
		},
		"else": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**else** - 条件语句的 else 分支\n\n与 if 语句配合使用，当 if 条件为假时执行。\n\n```origami\nif (condition) {\n    // 条件为真\n} else {\n    // 条件为假\n}\n```",
			},
		},
		"for": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**for** - for 循环语句\n\n用于重复执行代码块。\n\n```origami\nfor (i = 0; i < 10; i++) {\n    // 循环体\n}\n```",
			},
		},
		"foreach": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**foreach** - foreach 循环语句\n\n用于遍历数组或集合中的每个元素。\n\n```origami\nforeach (array as value) {\n    // 处理每个元素\n}\n```",
			},
		},
		"while": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**while** - while 循环语句\n\n当条件为真时重复执行代码块。\n\n```origami\nwhile (condition) {\n    // 循环体\n}\n```",
			},
		},

		// 函数和类定义
		"function": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**function** - 函数定义关键字\n\n用于定义可重用的代码块。\n\n```origami\nfunction functionName(param1, param2) {\n    // 函数体\n    return value;\n}\n```",
			},
		},
		"class": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**class** - 类定义关键字\n\n用于定义对象的模板。\n\n```origami\nclass ClassName {\n    // 属性和方法\n}\n```",
			},
		},

		// 异常处理
		"try": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**try** - 异常处理语句\n\n用于捕获和处理可能发生的异常。\n\n```origami\ntry {\n    // 可能抛出异常的代码\n} catch (Exception e) {\n    // 异常处理\n}\n```",
			},
		},
		"catch": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**catch** - 异常捕获语句\n\n与 try 语句配合使用，用于捕获和处理异常。\n\n```origami\ntry {\n    // 代码\n} catch (Exception e) {\n    // 处理异常\n}\n```",
			},
		},

		// 分支语句
		"switch": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**switch** - 分支语句\n\n根据变量的值执行不同的代码分支。\n\n```origami\nswitch (variable) {\n    case value1:\n        // 代码\n        break;\n    case value2:\n        // 代码\n        break;\n    default:\n        // 默认代码\n}\n```",
			},
		},
		"case": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**case** - switch 语句的分支标签\n\n定义 switch 语句中的一个分支。\n\n```origami\ncase value:\n    // 当 switch 变量等于 value 时执行\n    break;\n```",
			},
		},
		"default": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**default** - switch 语句的默认分支\n\n当所有 case 都不匹配时执行的默认分支。\n\n```origami\ndefault:\n    // 默认执行的代码\n```",
			},
		},
		"break": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**break** - 跳出语句\n\n用于跳出循环或 switch 语句。\n\n```origami\nfor (i = 0; i < 10; i++) {\n    if (condition) {\n        break; // 跳出循环\n    }\n}\n```",
			},
		},
		"continue": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**continue** - 继续语句\n\n跳过当前循环迭代，继续下一次迭代。\n\n```origami\nfor (i = 0; i < 10; i++) {\n    if (condition) {\n        continue; // 跳过本次迭代\n    }\n    // 其他代码\n}\n```",
			},
		},

		// 内置函数
		"echo": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**echo** - 输出函数\n\n将内容输出到标准输出。\n\n```origami\necho \"Hello, World!\";\necho variable;\n```",
			},
		},

		// 返回语句
		"return": {
			"contents": map[string]interface{}{
				"kind":  "markdown",
				"value": "**return** - 返回语句\n\n从函数中返回值并结束函数执行。\n\n```origami\nfunction getName() {\n    return \"Origami\";\n}\n```",
			},
		},
	}

	// 查找单词对应的悬停信息
	if info, exists := hoverData[word]; exists {
		return info
	}

	// 如果没有找到预定义的信息，返回基本信息
	return map[string]interface{}{
		"contents": map[string]interface{}{
			"kind":  "markdown",
			"value": "**" + word + "**\n\nOrigami 语言标识符",
		},
	}
}
