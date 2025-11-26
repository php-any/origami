package completion

import (
	"strings"

	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/sirupsen/logrus"
)

// GetCompletionItems 获取补全项
func GetCompletionItems(content string, position defines.Position, provider defines.SymbolProvider) []defines.CompletionItem {
	items := []defines.CompletionItem{}

	// 获取光标位置前的文本
	lines := strings.Split(content, "\n")
	if int(position.Line) >= len(lines) {
		return getDefaultCompletions()
	}

	line := lines[position.Line]
	if int(position.Character) > len(line) {
		return getDefaultCompletions()
	}

	beforeCursor := line[:position.Character]
	trimmedBefore := strings.TrimSpace(beforeCursor)

	// 如果光标前为空或只有空白，提供默认补全
	if len(trimmedBefore) == 0 {
		logrus.Infof("光标前为空，提供默认补全，位置: %d:%d", position.Line, position.Character)
		return getDefaultCompletions()
	}

	// 获取光标左边最后一个有意义的符号
	lastSymbol := getLastSymbol(beforeCursor)
	logrus.Infof("光标左边符号: %s, 位置: %d:%d", lastSymbol, position.Line, position.Character)

	// 根据光标左边的符号 switch 进入不同分支处理
	switch lastSymbol {
	case "->", ".":
		// 对象方法/属性补全：$obj->
		dynamicItems := getObjectPropertyAndMethodCompletions(content, position, provider)
		if len(dynamicItems) > 0 {
			items = append(items, dynamicItems...)
			logrus.Infof("动态获取到 %d 个对象属性/方法", len(dynamicItems))
		} else {
			// 如果没有动态获取到，添加通用方法作为备选
			items = append(items, getObjectMethodCompletions()...)
		}
		logrus.Infof("对象方法补全：%d 个项", len(items))

	case "::":
		// 静态方法/属性补全：ClassName::
		// TODO: 实现静态成员补全
		items = append(items, getObjectMethodCompletions()...)
		logrus.Infof("静态成员补全：%d 个项", len(items))

	case "$":
		// 变量补全：$
		// TODO: 实现变量补全
		items = append(items, getKeywordCompletions()...)
		items = append(items, getGlobalFunctionCompletions()...)
		logrus.Infof("变量补全：%d 个项", len(items))

	case "new":
		// 类实例化补全：new
		items = append(items, getGlobalClassCompletions()...)
		logrus.Infof("类实例化补全：%d 个项", len(items))

	case "keyword":
		// 关键字补全：正在输入关键字
		items = append(items, getKeywordCompletions()...)
		items = append(items, getGlobalFunctionCompletions()...)
		items = append(items, getGlobalClassCompletions()...)
		logrus.Infof("关键字补全：%d 个项", len(items))

	case "snippet":
		// 代码片段补全：特定关键字后
		items = append(items, getSnippetCompletions()...)
		items = append(items, getKeywordCompletions()...)
		logrus.Infof("代码片段补全：%d 个项", len(items))

	default:
		// 默认补全：提供所有类型的补全
		items = append(items, getKeywordCompletions()...)
		items = append(items, getSnippetCompletions()...)
		items = append(items, getGlobalFunctionCompletions()...)
		items = append(items, getGlobalClassCompletions()...)
		logrus.Infof("默认补全：%d 个项", len(items))
	}

	return items
}

// getLastSymbol 获取光标左边最后一个有意义的符号
func getLastSymbol(beforeCursor string) string {
	if len(beforeCursor) == 0 {
		return ""
	}

	trimmedBefore := strings.TrimSpace(beforeCursor)
	if len(trimmedBefore) == 0 {
		return ""
	}

	// 从后往前查找最后一个有意义的符号
	// 优先级：-> > :: > $ > 关键字 > 其他

	// 1. 检查 -> (对象方法调用)
	lastArrowIdx := strings.LastIndex(beforeCursor, "->")
	if lastArrowIdx != -1 {
		// 检查 -> 后面是否只包含合法的标识符字符（或者是空的）
		afterArrow := beforeCursor[lastArrowIdx+2:]
		isValid := true
		for _, c := range afterArrow {
			if !isVarChar(byte(c)) {
				isValid = false
				break
			}
		}
		if isValid {
			return "->"
		}
	}

	// 2. 检查 . (类静态属性访问，兼容类似类名.常量/函数调用)
	lastDotIdx := strings.LastIndex(beforeCursor, ".")
	if lastDotIdx != -1 {
		afterDot := beforeCursor[lastDotIdx+1:]
		isValidAfter := true
		for _, c := range afterDot {
			if !isVarChar(byte(c)) {
				isValidAfter = false
				break
			}
		}
		if isValidAfter {
			return "."
		}
	}

	// 3. 检查 :: (静态成员访问)
	lastColonIdx := strings.LastIndex(beforeCursor, "::")
	if lastColonIdx != -1 {
		// 检查 :: 后面是否只包含合法的标识符字符（或者是空的）
		afterColon := beforeCursor[lastColonIdx+2:]
		isValid := true
		for _, c := range afterColon {
			if !isVarChar(byte(c)) {
				isValid = false
				break
			}
		}
		if isValid {
			return "::"
		}
	}

	// 4. 检查 $ (变量)
	lastDollarIdx := strings.LastIndex(beforeCursor, "$")
	if lastDollarIdx != -1 {
		// 检查 $ 后面是否只包含合法的标识符字符（或者是空的）
		afterDollar := beforeCursor[lastDollarIdx+1:]
		isValid := true
		for _, c := range afterDollar {
			if !isVarChar(byte(c)) {
				isValid = false
				break
			}
		}
		if isValid {
			return "$"
		}
	}

	// 5. 检查关键字（如 new, func, class 等）
	// 从后往前查找最后一个完整的单词
	// 先去除末尾可能的空白和标识符字符，找到单词边界
	wordEnd := len(trimmedBefore)
	wordStart := wordEnd
	// 从后往前找到单词的起始位置
	for wordStart > 0 {
		c := trimmedBefore[wordStart-1]
		if isVarChar(c) {
			wordStart--
		} else {
			break
		}
	}

	if wordStart < wordEnd {
		lastWord := trimmedBefore[wordStart:wordEnd]
		// 检查是否是代码片段关键字
		snippetKeywords := []string{"func", "class", "if", "foreach", "while", "for", "switch"}
		for _, keyword := range snippetKeywords {
			if lastWord == keyword {
				return "snippet"
			}
		}
		// 检查是否是 new 关键字
		if lastWord == "new" {
			return "new"
		}
	}

	// 如果最后一个字符是字母，可能是正在输入关键字
	lastChar := trimmedBefore[len(trimmedBefore)-1]
	if (lastChar >= 'a' && lastChar <= 'z') || (lastChar >= 'A' && lastChar <= 'Z') {
		return "keyword"
	}

	// 6. 默认情况
	return "default"
}

// getDefaultCompletions 获取默认补全项
func getDefaultCompletions() []defines.CompletionItem {
	items := []defines.CompletionItem{}
	items = append(items, getKeywordCompletions()...)
	items = append(items, getSnippetCompletions()...)
	items = append(items, getGlobalFunctionCompletions()...)
	items = append(items, getGlobalClassCompletions()...)
	return items
}

// getObjectPropertyAndMethodCompletions 获取对象属性和方法补全
func getObjectPropertyAndMethodCompletions(content string, position defines.Position, provider defines.SymbolProvider) []defines.CompletionItem {
	if provider == nil {
		logrus.Warn("SymbolProvider 为空，无法获取动态补全")
		return nil
	}

	// 1. 获取光标前的变量名
	// 简单的字符串分析：查找 -> 前面的单词
	lines := strings.Split(content, "\n")
	if int(position.Line) >= len(lines) {
		return nil
	}
	line := lines[position.Line]
	beforeCursor := line[:position.Character]

	// 找到最后一次出现的 -> 或 .
	idxArrow := strings.LastIndex(beforeCursor, "->")
	idxDot := strings.LastIndex(beforeCursor, ".")

	token := ""
	idx := -1
	if idxArrow == -1 && idxDot == -1 {
		return nil
	}

	if idxArrow > idxDot {
		token = "->"
		idx = idxArrow
	} else {
		token = "."
		idx = idxDot
	}

	// 提取变量名部分，例如 $user-> 中的 $user
	varPart := strings.TrimSpace(beforeCursor[:idx])
	// 取出最后一个单词，假设是变量
	// 从后往前找，直到遇到非变量字符
	varEnd := len(varPart)
	varStart := varEnd
	for varStart > 0 {
		c := varPart[varStart-1]
		if isVarChar(c) || c == '$' {
			varStart--
		} else {
			break
		}
	}
	varName := varPart[varStart:varEnd]

	if varName == "" {
		return nil
	}

	logrus.Infof("尝试获取变量 %s 的类型，触发符号: %s", varName, token)

	// 2. 获取变量类型
	className := provider.GetVariableTypeAtPosition(content, position, varName)
	if className == "" {
		logrus.Infof("未找到变量 %s 的类型", varName)
		return nil
	}

	logrus.Infof("变量 %s 的类型为 %s", varName, className)

	// 3. 获取类成员
	items := provider.GetClassMembers(className)
	logrus.Infof("找到类 %s 的成员：%d 个", className, len(items))
	return items
}

func isVarChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// ... (以下是静态补全数据函数，保持不变) ...

func getSnippetCompletions() []defines.CompletionItem {
	insertTextFormatSnippet := defines.InsertTextFormatSnippet

	// Helper variables for pointers
	funcInsert := "func ${1:name}(${2:params}) {\n\t${0}\n}"
	classInsert := "class ${1:Name} {\n\tfunc ${1:Name}() {\n\t\t${0}\n\t}\n}"
	ifInsert := "if (${1:condition}) {\n\t${0}\n}"
	foreachInsert := "foreach (${1:array} as ${2:key} => ${3:value}) {\n\t${0}\n}"

	return []defines.CompletionItem{
		{
			Label:            "func",
			Kind:             &[]defines.CompletionItemKind{defines.CompletionItemKindSnippet}[0],
			Detail:           &[]string{"Function definition"}[0],
			Documentation:    &defines.MarkupContent{Kind: defines.MarkupKindMarkdown, Value: "Define a new function"},
			InsertText:       &funcInsert,
			InsertTextFormat: &insertTextFormatSnippet,
		},
		{
			Label:            "class",
			Kind:             &[]defines.CompletionItemKind{defines.CompletionItemKindSnippet}[0],
			Detail:           &[]string{"Class definition"}[0],
			Documentation:    &defines.MarkupContent{Kind: defines.MarkupKindMarkdown, Value: "Define a new class"},
			InsertText:       &classInsert,
			InsertTextFormat: &insertTextFormatSnippet,
		},
		{
			Label:            "if",
			Kind:             &[]defines.CompletionItemKind{defines.CompletionItemKindSnippet}[0],
			Detail:           &[]string{"If statement"}[0],
			InsertText:       &ifInsert,
			InsertTextFormat: &insertTextFormatSnippet,
		},
		{
			Label:            "foreach",
			Kind:             &[]defines.CompletionItemKind{defines.CompletionItemKindSnippet}[0],
			Detail:           &[]string{"Foreach loop"}[0],
			InsertText:       &foreachInsert,
			InsertTextFormat: &insertTextFormatSnippet,
		},
	}
}

func getKeywordCompletions() []defines.CompletionItem {
	keywords := []string{
		"func", "class", "if", "else", "elseif", "while", "for", "foreach",
		"return", "break", "continue", "switch", "case", "default",
		"try", "catch", "finally", "throw", "new", "var", "const",
		"public", "private", "protected", "static", "extends", "implements",
		"interface", "trait", "use", "namespace", "echo", "print",
		"true", "false", "null", "this", "parent",
	}

	items := make([]defines.CompletionItem, len(keywords))
	for i, keyword := range keywords {
		items[i] = defines.CompletionItem{
			Label: keyword,
			Kind:  &[]defines.CompletionItemKind{defines.CompletionItemKindKeyword}[0],
		}
	}
	return items
}

func getGlobalFunctionCompletions() []defines.CompletionItem {
	functions := []string{
		"print_r", "var_dump", "count", "strlen", "strpos", "substr",
		"array_merge", "array_push", "json_encode", "json_decode",
		"file_get_contents", "file_put_contents", "date", "time",
	}

	items := make([]defines.CompletionItem, len(functions))
	for i, fn := range functions {
		insertText := fn + "($1)"
		items[i] = defines.CompletionItem{
			Label:      fn,
			Kind:       &[]defines.CompletionItemKind{defines.CompletionItemKindFunction}[0],
			InsertText: &insertText,
			InsertTextFormat: &[]defines.InsertTextFormat{
				defines.InsertTextFormatSnippet,
			}[0],
		}
	}
	return items
}

func getGlobalClassCompletions() []defines.CompletionItem {
	classes := []string{
		"Exception", "Error", "DateTime", "PDO", "stdClass",
	}

	items := make([]defines.CompletionItem, len(classes))
	for i, cls := range classes {
		items[i] = defines.CompletionItem{
			Label: cls,
			Kind:  &[]defines.CompletionItemKind{defines.CompletionItemKindClass}[0],
		}
	}
	return items
}

func getObjectMethodCompletions() []defines.CompletionItem {
	methods := []string{
		"toString", "toJSON",
	}

	items := make([]defines.CompletionItem, len(methods))
	for i, method := range methods {
		insertText := method + "($1)"
		items[i] = defines.CompletionItem{
			Label:      method,
			Kind:       &[]defines.CompletionItemKind{defines.CompletionItemKindMethod}[0],
			InsertText: &insertText,
			InsertTextFormat: &[]defines.InsertTextFormat{
				defines.InsertTextFormatSnippet,
			}[0],
		}
	}
	return items
}
