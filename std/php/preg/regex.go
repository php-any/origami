package preg

import (
	"regexp"
)

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
