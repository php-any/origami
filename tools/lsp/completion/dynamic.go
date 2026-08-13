package completion

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/sirupsen/logrus"
)

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

	// 2. 获取变量类型（可能包含多个类型）
	varTypeObj := provider.GetVariableTypeObjectAtPosition(content, position, varName)
	if varTypeObj == nil {
		logrus.Infof("未找到变量 %s 的类型", varName)
		return nil
	}

	// 类型断言
	varType, ok := varTypeObj.(data.Types)
	if !ok {
		logrus.Infof("变量 %s 的类型不是 data.Types", varName)
		return nil
	}

	// 3. 提取所有可能的类名
	classNames := extractClassNamesFromType(varType)
	if len(classNames) == 0 {
		logrus.Infof("变量 %s 的类型不包含任何类", varName)
		return nil
	}

	logrus.Infof("变量 %s 可能的类型: %v", varName, classNames)

	// 4. 收集所有可能类的成员，使用 map 去重
	itemsMap := make(map[string]defines.CompletionItem)
	for _, className := range classNames {
		members := provider.GetClassMembers(className)
		for _, member := range members {
			// 使用 Label 作为唯一键，避免重复
			if _, exists := itemsMap[member.Label]; !exists {
				itemsMap[member.Label] = member
			}
		}
	}

	// 转换为数组
	items := make([]defines.CompletionItem, 0, len(itemsMap))
	for _, item := range itemsMap {
		items = append(items, item)
	}
	logrus.Infof("找到合并后的成员：%d 个", len(items))
	return items
}

// getStaticMemberCompletions 获取静态成员补全
func getStaticMemberCompletions(content string, position defines.Position, provider defines.SymbolProvider) []defines.CompletionItem {
	if provider == nil {
		return nil
	}

	// 1. 获取光标前的类名
	lines := strings.Split(content, "\n")
	if int(position.Line) >= len(lines) {
		return nil
	}
	line := lines[position.Line]
	beforeCursor := line[:position.Character]

	// 找到最后一次出现的 ::
	idxColon := strings.LastIndex(beforeCursor, "::")
	if idxColon == -1 {
		return nil
	}

	// 提取类名部分
	classPart := strings.TrimSpace(beforeCursor[:idxColon])
	// 取出最后一个单词，假设是类名
	classEnd := len(classPart)
	classStart := classEnd
	for classStart > 0 {
		c := classPart[classStart-1]
		if isVarChar(c) || c == '\\' { // 类名可能包含命名空间分隔符
			classStart--
		} else {
			break
		}
	}
	className := classPart[classStart:classEnd]

	if className == "" {
		return nil
	}

	logrus.Infof("尝试获取类 %s 的静态成员", className)

	// 2. 获取静态成员
	items := provider.GetStaticClassMembers(className)
	logrus.Infof("找到类 %s 的静态成员：%d 个", className, len(items))
	return items
}

// getVariableCompletions 获取变量补全
func getVariableCompletions(content string, position defines.Position, provider defines.SymbolProvider) []defines.CompletionItem {
	if provider == nil {
		return nil
	}

	logrus.Infof("尝试获取当前位置的变量补全")
	items := provider.GetVariablesAtPosition(content, position)
	logrus.Infof("找到变量补全：%d 个", len(items))
	return items
}

// getClassCompletionsForNew 获取 new 关键字后的类名补全
func getClassCompletionsForNew(content string, position defines.Position, provider defines.SymbolProvider, lastSymbol SymbolProvider) []defines.CompletionItem {
	if provider == nil {
		return nil
	}

	logrus.Infof("尝试获取 new 关键字的类补全，Worker: %s", lastSymbol.Worker)

	// 调用 provider 的方法获取上下文相关的类（已经在 GetClassCompletionsForContext 中根据 Worker 过滤）
	items := provider.GetClassCompletionsForContext(content, position, lastSymbol.Worker)

	// 按照优先级排序：use 导入的类 > 同级目录的类 > 全局类
	sortedItems := sortClassesByPriority(items)

	logrus.Infof("找到类补全：%d 个", len(sortedItems))
	return sortedItems
}

// sortClassesByPriority 按照优先级排序类补全项
// 优先级：use 导入的类 > 同级目录的类 > 全局类
func sortClassesByPriority(items []defines.CompletionItem) []defines.CompletionItem {
	useItems := make([]defines.CompletionItem, 0)
	sameDirItems := make([]defines.CompletionItem, 0)
	globalItems := make([]defines.CompletionItem, 0)

	for _, item := range items {
		if item.Detail != nil {
			detail := *item.Detail
			if strings.HasPrefix(detail, "use ") {
				useItems = append(useItems, item)
			} else if detail == "同级目录" {
				sameDirItems = append(sameDirItems, item)
			} else if detail == "全局类" {
				globalItems = append(globalItems, item)
			} else {
				// 未知类型，放到最后
				globalItems = append(globalItems, item)
			}
		} else {
			// 没有 Detail 信息，放到最后
			globalItems = append(globalItems, item)
		}
	}

	// 合并结果：use 导入的类 > 同级目录的类 > 全局类
	result := make([]defines.CompletionItem, 0, len(items))
	result = append(result, useItems...)
	result = append(result, sameDirItems...)
	result = append(result, globalItems...)

	return result
}
