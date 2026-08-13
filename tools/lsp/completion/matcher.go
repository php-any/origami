package completion

import (
	"strings"
)

// 补全定位信息
type SymbolProvider struct {
	TypeString string // 进入模式类型, 比如 new
	Worker     string // 光标符号, 比如 U符号，就可以只提示有 U字母的类
}

// getLastSymbol 获取光标左边最后一个有意义的符号
func getLastSymbol(beforeCursor string) SymbolProvider {
	if len(beforeCursor) == 0 {
		return SymbolProvider{TypeString: "", Worker: ""}
	}

	trimmedBefore := strings.TrimSpace(beforeCursor)
	if len(trimmedBefore) == 0 {
		return SymbolProvider{TypeString: "", Worker: ""}
	}

	// 提取光标位置正在输入的符号作为 Worker
	worker := extractWorker(beforeCursor)

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
			return SymbolProvider{TypeString: "->", Worker: worker}
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
			return SymbolProvider{TypeString: ".", Worker: worker}
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
			return SymbolProvider{TypeString: "::", Worker: worker}
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
			return SymbolProvider{TypeString: "$", Worker: worker}
		}
	}

	// 5. 检查关键字（如 new, func, class 等）
	// 需要检查光标前的完整单词序列
	// 例如：对于 "new U"，我们需要找到 "new"

	// 移除尾部正在输入的部分单词，找到前面完整的单词
	temp := strings.TrimRight(trimmedBefore, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_")
	temp = strings.TrimSpace(temp)

	// 从后往前找最后一个完整单词
	words := strings.Fields(temp)
	if len(words) > 0 {
		lastCompleteWord := words[len(words)-1]

		// 检查是否是 new 关键字
		if lastCompleteWord == "new" {
			return SymbolProvider{TypeString: "new", Worker: worker}
		}

		// 检查是否是代码片段关键字
		snippetKeywords := []string{"func", "class", "if", "foreach", "while", "for", "switch"}
		for _, keyword := range snippetKeywords {
			if lastCompleteWord == keyword {
				return SymbolProvider{TypeString: "snippet", Worker: worker}
			}
		}
	}

	// 如果最后一个字符是字母，可能是正在输入关键字或标识符
	lastChar := trimmedBefore[len(trimmedBefore)-1]
	if (lastChar >= 'a' && lastChar <= 'z') || (lastChar >= 'A' && lastChar <= 'Z') {
		return SymbolProvider{TypeString: "keyword", Worker: worker}
	}

	// 6. 默认情况
	return SymbolProvider{TypeString: "default", Worker: worker}
}

// extractWorker 提取光标位置正在输入的符号
func extractWorker(beforeCursor string) string {
	if len(beforeCursor) == 0 {
		return ""
	}

	// 从后往前找到最后一个非标识符字符的位置
	workerStart := len(beforeCursor)
	for i := len(beforeCursor) - 1; i >= 0; i-- {
		if !isVarChar(beforeCursor[i]) {
			workerStart = i + 1
			break
		}
	}

	// 如果 workerStart 在字符串末尾，说明正在输入标识符
	if workerStart < len(beforeCursor) {
		return beforeCursor[workerStart:]
	}

	// 如果整个字符串末尾都是标识符字符，返回整个末尾部分
	// 但需要排除已经识别的操作符（->, ::, $, .）
	trimmed := strings.TrimSpace(beforeCursor)
	if len(trimmed) > 0 {
		// 检查是否以操作符结尾
		if strings.HasSuffix(trimmed, "->") || strings.HasSuffix(trimmed, "::") {
			return ""
		}
		// 提取末尾的标识符部分
		worker := ""
		for i := len(trimmed) - 1; i >= 0; i-- {
			if isVarChar(trimmed[i]) {
				worker = string(trimmed[i]) + worker
			} else {
				break
			}
		}
		return worker
	}

	return ""
}

func isVarChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}
