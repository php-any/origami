package completion

import (
	"strings"
)

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
			return "new"
		}

		// 检查是否是代码片段关键字
		snippetKeywords := []string{"func", "class", "if", "foreach", "while", "for", "switch"}
		for _, keyword := range snippetKeywords {
			if lastCompleteWord == keyword {
				return "snippet"
			}
		}
	}

	// 如果最后一个字符是字母，可能是正在输入关键字或标识符
	lastChar := trimmedBefore[len(trimmedBefore)-1]
	if (lastChar >= 'a' && lastChar <= 'z') || (lastChar >= 'A' && lastChar <= 'Z') {
		return "keyword"
	}

	// 6. 默认情况
	return "default"
}

func isVarChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}
