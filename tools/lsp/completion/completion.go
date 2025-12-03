package completion

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/sirupsen/logrus"
)

// DocumentInfoProvider 提供访问 DocumentInfo 的接口
type DocumentInfoProvider interface {
	GetAST() *node.Program
}

// VMProvider 提供访问 VM 的接口
type VMProvider interface {
	GetClass(className string) (data.ClassStmt, bool)
	GetAllClasses() map[string]data.ClassStmt
	// GetAllFunctions 返回所有已知函数（包含路径名，例如 http\app）
	GetAllFunctions() map[string]data.FuncStmt
}

// GetCompletionItems 获取补全项（保持向后兼容）
func GetCompletionItems(content string, position defines.Position, provider defines.SymbolProvider) []defines.CompletionItem {
	return GetCompletionItemsWithDoc(content, position, provider, nil, nil)
}

// GetCompletionItemsWithDoc 获取补全项，支持基于节点的补全
func GetCompletionItemsWithDoc(content string, position defines.Position, provider defines.SymbolProvider, docProvider DocumentInfoProvider, vmProvider VMProvider) []defines.CompletionItem {
	items := []defines.CompletionItem{}

	// 获取光标位置前的文本
	lines := strings.Split(content, "\n")
	if int(position.Line) >= len(lines) {
		return getDefaultCompletions(vmProvider)
	}

	line := lines[position.Line]
	if int(position.Character) > len(line) {
		return getDefaultCompletions(vmProvider)
	}

	beforeCursor := line[:position.Character]
	trimmedBefore := strings.TrimSpace(beforeCursor)

	// 如果光标前为空或只有空白，提供默认补全
	if len(trimmedBefore) == 0 {
		logrus.Infof("光标前为空，提供默认补全，位置: %d:%d", position.Line, position.Character)
		return getDefaultCompletions(vmProvider)
	}

	// 获取光标左边最后一个有意义的符号
	lastSymbol := getLastSymbol(beforeCursor)
	logrus.Infof("光标左边符号: %s, Worker: %s, 位置: %d:%d", lastSymbol.TypeString, lastSymbol.Worker, position.Line, position.Character)

	// 根据光标左边的符号 switch 进入不同分支处理
	switch lastSymbol.TypeString {
	case "->", ".":
		// 对象方法/属性补全：$obj->
		// 根据左边节点来提示
		dynamicItems := getObjectPropertyAndMethodCompletionsFromNode(content, position, provider, docProvider, vmProvider)
		if len(dynamicItems) > 0 {
			items = append(items, dynamicItems...)
			logrus.Infof("基于节点获取到 %d 个对象属性/方法", len(dynamicItems))
		} else {
			// 如果基于节点无法获取，回退到基于变量名的方式
			dynamicItems = getObjectPropertyAndMethodCompletions(content, position, provider)
			if len(dynamicItems) > 0 {
				items = append(items, dynamicItems...)
				logrus.Infof("基于变量名获取到 %d 个对象属性/方法", len(dynamicItems))
			} else {
				// 如果都没有获取到，添加通用方法作为备选
				items = append(items, getObjectMethodCompletions()...)
			}
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
		// 根据 lastSymbol.Worker 过滤，优先显示 use 导入的类
		dynamicItems := getClassCompletionsForNew(content, position, provider, lastSymbol)
		// 对于通过 VM 提供的全局类（Detail 以 "full:" 开头），自动添加 use 语句
		dynamicItems = addUseEditsForClassItems(content, dynamicItems)
		items = append(items, dynamicItems...)
		logrus.Infof("类实例化补全：%d 个项", len(items))

	case "keyword":
		// 关键字补全：正在输入关键字
		worker := lastSymbol.Worker

		// 关键字 / 内置函数 / 内置类使用通用前缀过滤
		kwItems := filterCompletionsByPrefix(getKeywordCompletions(), worker)
		builtinFuncItems := filterCompletionsByPrefix(getGlobalFunctionCompletions(), worker)
		builtinClassItems := filterCompletionsByPrefix(getGlobalClassCompletions(), worker)
		// 项目函数只按短名（不含路径）与 worker 匹配，并自动添加 use 语句的附加编辑
		vmFuncItems := getGlobalFunctionCompletionsWithVM(worker, vmProvider)
		vmFuncItems = addUseEditsForFunctionItems(content, vmFuncItems)
		// 项目类（用于 ClassName::staticMethod 或 new 之前的类名输入）：
		// Label 为短名，Detail 中保存完整类名，便于自动添加 use
		vmClassItems := getGlobalClassCompletionsFromVM(worker, vmProvider)
		vmClassItems = addUseEditsForClassItems(content, vmClassItems)

		items = append(items, kwItems...)
		items = append(items, builtinFuncItems...)
		items = append(items, vmFuncItems...)
		items = append(items, builtinClassItems...)
		items = append(items, vmClassItems...)
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

// filterCompletionsByPrefix 根据当前正在输入的 worker（前缀）过滤补全项：
// - 优先使用 CompletionItem.FilterText 前缀匹配
// - 若无 FilterText，则使用 Label 前缀匹配
// worker 为空时直接返回原列表
func filterCompletionsByPrefix(items []defines.CompletionItem, worker string) []defines.CompletionItem {
	filtered := make([]defines.CompletionItem, 0, len(items))
	for _, it := range items {
		match := false

		if it.FilterText != nil && *it.FilterText != "" {
			if strings.HasPrefix(*it.FilterText, worker) {
				match = true
			}
		} else {
			if strings.HasPrefix(it.Label, worker) {
				match = true
			}
		}

		if match {
			filtered = append(filtered, it)
		}
	}
	return filtered
}

// addUseEditsForFunctionItems 为带命名空间的函数补全项自动添加 use 语句的 AdditionalTextEdits。
// 仅针对项目函数（来自 VM）的补全：
// - Label 为短名（例如 "app"）
// - Detail 为完整函数名（例如 "Net\\Http\\app"）
// - 如果当前文件尚未存在包含该全名的 use 行，则在文件头合适位置插入一行 "use Net\\Http\\app;"
func addUseEditsForFunctionItems(content string, items []defines.CompletionItem) []defines.CompletionItem {
	if len(items) == 0 || content == "" {
		return items
	}

	lines := strings.Split(content, "\n")

	// 收集已有 use 行，做一个简单的包含判断，避免重复插入
	existingUseLines := make([]string, 0)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "use ") {
			existingUseLines = append(existingUseLines, trimmed)
		}
	}

	hasUseFor := func(fullName string) bool {
		for _, u := range existingUseLines {
			if strings.Contains(u, fullName) {
				return true
			}
		}
		return false
	}

	// 计算插入 use 语句的行号：
	// - 如果已有 use 行，则插在最后一个 use 行之后
	// - 否则尽量在文件头、开头标签或第一行代码之后
	// 额外约束：至少从第 3 行开始插入（0-based 行号 >= 2）
	insertLine := 0
	lastUseLine := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "use ") {
			lastUseLine = i
			continue
		}
		// 跳过类似 <?origami / <?php 这样的文件头标签
		if strings.HasPrefix(trimmed, "<?") {
			insertLine = i + 1
			continue
		}
		// 遇到第一行非 use 的代码，就可以停止
		if lastUseLine >= 0 {
			insertLine = lastUseLine + 1
		} else if insertLine == 0 {
			insertLine = i
		}
		break
	}
	if lastUseLine >= 0 {
		insertLine = lastUseLine + 1
	}

	// 保证 use 至少从第三行开始插入（例如预留前两行给头部注释或标签）
	if insertLine < 2 {
		insertLine = 2
	}

	for i := range items {
		// 优先从 Detail 中读取完整函数名；若没有则回退到 Label
		fullName := items[i].Label
		if items[i].Detail != nil && *items[i].Detail != "" {
			fullName = *items[i].Detail
		}

		// 只处理带命名空间分隔符的函数，例如 http\app 或 http/app
		if !strings.Contains(fullName, `\`) && !strings.Contains(fullName, `/`) {
			continue
		}

		// 已经存在对应 use 行，则不再重复插入
		if hasUseFor(fullName) {
			continue
		}

		newText := "use " + fullName + ";\n"
		edit := defines.TextEdit{
			Range: defines.Range{
				Start: defines.Position{Line: uint32(insertLine), Character: 0},
				End:   defines.Position{Line: uint32(insertLine), Character: 0},
			},
			NewText: newText,
		}

		items[i].AdditionalTextEdits = append(items[i].AdditionalTextEdits, edit)
	}

	return items
}

// addUseEditsForClassItems 为 new 关键字后的类补全自动添加 use 语句：
// - 仅处理 Detail 以 "full:" 开头的类（来自 VM 的全局类）
// - Label 为短名（例如 "Server"），Detail 为 "full:Net\\Http\\Server"
// - 如果当前文件尚未存在对应的 use 行，则在文件头合适位置插入 "use Net\\Http\\Server;"
func addUseEditsForClassItems(content string, items []defines.CompletionItem) []defines.CompletionItem {
	if len(items) == 0 || content == "" {
		return items
	}

	lines := strings.Split(content, "\n")

	// 收集已有 use 行
	existingUseLines := make([]string, 0)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "use ") {
			existingUseLines = append(existingUseLines, trimmed)
		}
	}

	hasUseFor := func(fullName string) bool {
		for _, u := range existingUseLines {
			if strings.Contains(u, fullName) {
				return true
			}
		}
		return false
	}

	// 计算插入 use 语句的行号，与函数 use 逻辑保持一致
	insertLine := 0
	lastUseLine := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "use ") {
			lastUseLine = i
			continue
		}
		if strings.HasPrefix(trimmed, "<?") {
			insertLine = i + 1
			continue
		}
		if lastUseLine >= 0 {
			insertLine = lastUseLine + 1
		} else if insertLine == 0 {
			insertLine = i
		}
		break
	}
	if lastUseLine >= 0 {
		insertLine = lastUseLine + 1
	}
	if insertLine < 2 {
		insertLine = 2
	}

	for i := range items {
		if items[i].Detail == nil {
			continue
		}
		detail := *items[i].Detail
		if !strings.HasPrefix(detail, "full:") {
			continue
		}
		fullName := strings.TrimPrefix(detail, "full:")
		if fullName == "" {
			continue
		}

		// 已存在对应 use 则跳过
		if hasUseFor(fullName) {
			continue
		}

		newText := "use " + fullName + ";\n"
		edit := defines.TextEdit{
			Range: defines.Range{
				Start: defines.Position{Line: uint32(insertLine), Character: 0},
				End:   defines.Position{Line: uint32(insertLine), Character: 0},
			},
			NewText: newText,
		}

		items[i].AdditionalTextEdits = append(items[i].AdditionalTextEdits, edit)
	}

	return items
}
