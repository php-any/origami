package completion

import (
	"strings"

	"github.com/php-any/origami/tools/lsp/defines"
)

// getDefaultCompletions 获取默认补全项
// 默认情况下：关键字 + 代码片段 + 内置函数 + 项目函数 + 内置类
func getDefaultCompletions(vmProvider VMProvider) []defines.CompletionItem {
	items := []defines.CompletionItem{}
	items = append(items, getKeywordCompletions()...)
	items = append(items, getSnippetCompletions()...)
	items = append(items, getGlobalFunctionCompletions()...)
	items = append(items, getGlobalFunctionCompletionsWithVM("", vmProvider)...)
	items = append(items, getGlobalClassCompletions()...)
	return items
}

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

// getGlobalFunctionCompletionsWithVM 返回来自 VM 的项目函数（包括带路径的函数，如 Net\Http\app）。
// - worker：当前正在输入的标识符，用于按函数短名（最后一段）前缀过滤
// - Label：短名（例如 "app"）
// - Detail：完整函数名（包含路径），例如 "Net\\Http\\app"
func getGlobalFunctionCompletionsWithVM(worker string, vmProvider VMProvider) []defines.CompletionItem {
	if vmProvider == nil {
		return nil
	}

	funcs := vmProvider.GetAllFunctions()
	if len(funcs) == 0 {
		return nil
	}

	items := make([]defines.CompletionItem, 0, len(funcs))

	kind := defines.CompletionItemKindFunction
	snippetFormat := defines.InsertTextFormatSnippet

	for fullName := range funcs {
		// 取短名：最后一个 '\' 或 '/' 之后的部分
		shortName := fullName
		if idx := strings.LastIndexAny(fullName, "\\/"); idx >= 0 && idx+1 < len(fullName) {
			shortName = fullName[idx+1:]
		}

		// 只按短名与 worker 的前缀进行匹配，例如输入 "app" 匹配 "Net\\Http\\app"
		if !strings.HasPrefix(shortName, worker) {
			continue
		}

		insertText := shortName + "($1)"

		// Label 用短名，Detail 保存完整名，方便后续生成 use 语句
		detail := fullName
		item := defines.CompletionItem{
			Label:            shortName,
			Kind:             &kind,
			Detail:           &detail,
			InsertText:       &insertText,
			InsertTextFormat: &snippetFormat,
		}

		// FilterText/SortText 使用短名即可

		items = append(items, item)
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

// getGlobalClassCompletionsFromVM 返回来自 VM 的项目类（包括带命名空间的类，如 Net\Http\Server）。
// - worker：当前正在输入的标识符，用于按类短名（最后一段）做前缀过滤
// - Label：短名（例如 "Server"）
// - Detail："full:Net\\Http\\Server" 这样的完整类名描述，方便后续自动添加 use 语句
func getGlobalClassCompletionsFromVM(worker string, vmProvider VMProvider) []defines.CompletionItem {
	if vmProvider == nil {
		return nil
	}

	allClasses := vmProvider.GetAllClasses()
	if len(allClasses) == 0 {
		return nil
	}

	items := make([]defines.CompletionItem, 0, len(allClasses))
	kind := defines.CompletionItemKindClass

	for fullName := range allClasses {
		// 提取短名：命名空间最后一段
		shortName := fullName
		if idx := strings.LastIndexAny(fullName, "\\/"); idx >= 0 && idx+1 < len(fullName) {
			shortName = fullName[idx+1:]
		}

		// 只按短名前缀过滤，例如输入 "Ser" 匹配 "Net\\Http\\Server"
		if !strings.HasPrefix(shortName, worker) {
			continue
		}

		detail := "full:" + fullName
		item := defines.CompletionItem{
			Label:  shortName,
			Kind:   &kind,
			Detail: &detail,
		}

		items = append(items, item)
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
