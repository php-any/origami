package preg

import (
	"regexp"
	"strings"
)

// convertPossessiveQuantifiers 将 PHP 的占有量词转换为 Go regexp 兼容的语法
// ++ -> + (占有量词转换为贪婪量词)
// *+ -> * (占有量词转换为贪婪量词)
// ?+ -> ? (占有量词转换为贪婪量词)
// {n}+ -> {n} (占有量词转换为贪婪量词)
// {n,}+ -> {n,} (占有量词转换为贪婪量词)
// {n,m}+ -> {n,m} (占有量词转换为贪婪量词)
func convertPossessiveQuantifiers(pattern string) string {
	// 简单的字符串替换，将 ++ 替换为 +
	// 注意：需要避免在字符类 [] 内替换
	result := strings.Builder{}
	inCharClass := false
	escaped := false

	for i := 0; i < len(pattern); i++ {
		char := pattern[i]

		if escaped {
			result.WriteByte(char)
			escaped = false
			continue
		}

		if char == '\\' {
			escaped = true
			result.WriteByte(char)
			continue
		}

		if char == '[' && !escaped {
			inCharClass = true
			result.WriteByte(char)
			continue
		}

		if char == ']' && inCharClass {
			inCharClass = false
			result.WriteByte(char)
			continue
		}

		if !inCharClass && char == '+' && i+1 < len(pattern) && pattern[i+1] == '+' {
			// 发现 ++，替换为 +
			result.WriteByte('+')
			i++ // 跳过下一个 +
			continue
		}

		if !inCharClass && char == '*' && i+1 < len(pattern) && pattern[i+1] == '+' {
			// 发现 *+，替换为 *
			result.WriteByte('*')
			i++ // 跳过 +
			continue
		}

		if !inCharClass && char == '?' && i+1 < len(pattern) && pattern[i+1] == '+' {
			// 发现 ?+，替换为 ?
			result.WriteByte('?')
			i++ // 跳过 +
			continue
		}

		// 处理 {n}+, {n,}+, {n,m}+ 形式的占有量词
		if !inCharClass && char == '}' && i+1 < len(pattern) && pattern[i+1] == '+' {
			// 找到 }+，替换为 }
			result.WriteByte('}')
			i++ // 跳过 +
			continue
		}

		result.WriteByte(char)
	}

	return result.String()
}

// Compile 将 PHP 风格的正则表达式转换为 Go 的 regexp 并编译。
//
// 支持的特性:
//   - 分隔符: /pattern/ 或 #pattern# 等，使用第一个字符作为分隔符
//   - 转义分隔符: 通过反斜杠进行转义，例如: \/ 不视为结束分隔符
//   - 修饰符:
//   - i: 大小写不敏感
//   - m: 多行模式
//   - s: 点号匹配换行
//
// 例如:
//
//	/abc/i  -> (?i)abc
//	/abc/ms -> (?m)(?s)abc
func Compile(pattern string) (*regexp.Regexp, error) {
	if len(pattern) >= 2 {
		delimiter := pattern[0]
		endIndex := -1
		// 从尾部向前查找未转义的分隔符
		for i := len(pattern) - 1; i > 0; i-- {
			if pattern[i] == delimiter && pattern[i-1] != '\\' {
				endIndex = i
				break
			}
		}

		if endIndex != -1 {
			modifiers := pattern[endIndex+1:]
			regexBody := pattern[1:endIndex]

			// 处理 PHP 特有的语法转换为 Go regexp 兼容的语法
			// 1. 占有量词 ++ 转换为 +（Go 不支持占有量词，但贪婪量词在大多数情况下行为相同）
			regexBody = convertPossessiveQuantifiers(regexBody)

			// 处理修饰符
			prefix := ""
			if len(modifiers) > 0 {
				for _, mod := range modifiers {
					switch mod {
					case 'i':
						prefix += "(?i)"
					case 'm':
						prefix += "(?m)"
					case 's':
						prefix += "(?s)"
					}
				}
			}

			pattern = prefix + regexBody
		}
	}

	return regexp.Compile(pattern)
}
