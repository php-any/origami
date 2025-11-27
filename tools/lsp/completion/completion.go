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
		// 只提供对象成员，不提供其他类型的补全
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
		// 只提供静态成员，不提供其他类型的补全
		dynamicItems := getStaticMemberCompletions(content, position, provider)
		if len(dynamicItems) > 0 {
			items = append(items, dynamicItems...)
		} else {
			// 如果没有动态获取到，添加通用方法作为备选
			items = append(items, getObjectMethodCompletions()...)
		}
		logrus.Infof("静态成员补全：%d 个项", len(items))

	case "$":
		// 变量补全：$
		// 只提供变量名，不提供关键字或函数
		dynamicItems := getVariableCompletions(content, position, provider)
		if len(dynamicItems) > 0 {
			items = append(items, dynamicItems...)
		}
		logrus.Infof("变量补全：%d 个项", len(items))

	case "new":
		// 类实例化补全：new
		// 只提供上下文相关的类名（use导入的类 + 同级目录的类）
		dynamicItems := getClassCompletionsForNew(content, position, provider)
		items = append(items, dynamicItems...)
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
