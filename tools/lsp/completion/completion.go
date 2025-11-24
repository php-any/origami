package completion

import (
	"strings"

	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/sirupsen/logrus"
)

// GetCompletionItems 获取补全项
func GetCompletionItems(content string, position defines.Position, provider defines.SymbolProvider) []defines.CompletionItem {
	items := []defines.CompletionItem{}

	// 获取当前位置的上下文
	context, triggerChar := getCompletionContext(content, position)
	logrus.Infof("补全上下文: %s, 触发字符: %s, 位置: %d:%d", context, triggerChar, position.Line, position.Character)

	// 根据上下文提供相应的补全项
	switch context {
	case "snippet":
		// 代码片段补全
		items = append(items, getSnippetCompletions()...)
		// 同时提供关键字补全
		items = append(items, getKeywordCompletions()...)
		logrus.Infof("代码片段补全：%d 个项", len(items))
	case "keyword":
		// 关键字补全
		items = append(items, getKeywordCompletions()...)
		// 也可以包含一些常用的全局函数或类
		items = append(items, getGlobalFunctionCompletions()...)
		items = append(items, getGlobalClassCompletions()...)
		logrus.Infof("关键字补全：%d 个项", len(items))
	case "object_method":
		// 对象方法补全
		// 动态获取对象属性和方法
		dynamicItems := getObjectPropertyAndMethodCompletions(content, position, provider)
		if len(dynamicItems) > 0 {
			items = append(items, dynamicItems...)
			logrus.Infof("动态获取到 %d 个对象属性/方法", len(dynamicItems))
		} else {
			// 如果没有动态获取到，添加通用方法作为备选
			items = append(items, getObjectMethodCompletions()...)
		}
		logrus.Infof("对象方法补全：%d 个项", len(items))
	default:
		// 默认提供所有类型的补全，但要注意排序和优先级
		items = append(items, getKeywordCompletions()...)
		items = append(items, getSnippetCompletions()...)
		items = append(items, getGlobalFunctionCompletions()...)
		items = append(items, getGlobalClassCompletions()...)
		logrus.Infof("默认补全：%d 个项", len(items))
	}

	return items
}

// getCompletionContext 分析当前位置的上下文
func getCompletionContext(content string, position defines.Position) (string, string) {
	lines := strings.Split(content, "\n")
	if int(position.Line) >= len(lines) {
		return "default", ""
	}

	line := lines[position.Line]
	if int(position.Character) > len(line) {
		return "default", ""
	}

	beforeCursor := line[:position.Character]
	trimmedBefore := strings.TrimSpace(beforeCursor)

	// 检查是否在对象方法调用中 (如 $str-> 或 $str->le)
	// 只要行内包含 -> 且光标在 -> 之后，我们就认为是 object_method
	// 更精确的判断：找到最后一个 ->，如果光标在它后面
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
			return "object_method", "->"
		}
	}

	// 简单的关键字触发逻辑 (实际应用中需要更复杂的词法分析)
	if len(trimmedBefore) > 0 {
		lastChar := trimmedBefore[len(trimmedBefore)-1]
		if (lastChar >= 'a' && lastChar <= 'z') || (lastChar >= 'A' && lastChar <= 'Z') {
			// 可能是正在输入关键字
			// 这里简单假设如果是单词字符结尾，可能是关键字补全
			// 但如果前面有 -> 则已经被上面的逻辑捕获
			return "keyword", ""
		}
	}

	// 默认情况
	return "default", ""
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

	// 找到最后一个 ->
	idx := strings.LastIndex(beforeCursor, "->")
	if idx == -1 {
		return nil
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

	logrus.Infof("尝试获取变量 %s 的类型", varName)

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
