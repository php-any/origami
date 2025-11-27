package completion

import (
	"github.com/php-any/origami/tools/lsp/defines"
)

// getDefaultCompletions 获取默认补全项
func getDefaultCompletions() []defines.CompletionItem {
	items := []defines.CompletionItem{}
	items = append(items, getKeywordCompletions()...)
	items = append(items, getSnippetCompletions()...)
	items = append(items, getGlobalFunctionCompletions()...)
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
