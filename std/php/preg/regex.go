package preg

import (
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
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

// parsePhpPattern 解析 PHP 风格的正则表达式，返回 (goPattern, regexp2Pattern, regexp2Flags)。
// goPattern 已转换为 Go regexp 可用的格式（无 lookahead 时），
// regexp2Pattern 是 PCRE 兼容的模式字符串（不含修饰符前缀），
// regexp2Flags 是 regexp2 的选项标志。
func parsePhpPattern(pattern string) (goPattern string, r2Pattern string, r2Flags regexp2.RegexOptions) {
	goPattern = pattern
	r2Pattern = pattern
	r2Flags = regexp2.None

	if len(pattern) < 2 {
		return
	}

	delimiter := pattern[0]
	endIndex := -1
	// 从尾部向前查找未转义的分隔符
	for i := len(pattern) - 1; i > 0; i-- {
		if pattern[i] == delimiter && pattern[i-1] != '\\' {
			endIndex = i
			break
		}
	}

	if endIndex == -1 {
		return
	}

	modifiers := pattern[endIndex+1:]
	regexBody := pattern[1:endIndex]

	// 处理占有量词
	regexBody = convertPossessiveQuantifiers(regexBody)

	// 处理修饰符
	goPrefix := ""
	extended := false
	for _, mod := range modifiers {
		switch mod {
		case 'i':
			goPrefix += "(?i)"
			r2Flags |= regexp2.IgnoreCase
		case 'm':
			goPrefix += "(?m)"
			r2Flags |= regexp2.Multiline
		case 's':
			goPrefix += "(?s)"
			r2Flags |= regexp2.Singleline
		case 'x':
			extended = true
			r2Flags |= regexp2.IgnorePatternWhitespace
		}
	}

	if extended {
		regexBody = stripExtendedWhitespace(regexBody)
	}

	goPattern = goPrefix + regexBody
	r2Pattern = regexBody
	return
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
	goPattern, _, _ := parsePhpPattern(pattern)
	return regexp.Compile(goPattern)
}

// -----------------------------------------------------------------------
// Matcher: 统一接口，兼容 Go regexp 和 regexp2（支持 lookahead/lookbehind）
// -----------------------------------------------------------------------

// Matcher 提供正则匹配所需的通用操作接口。
type Matcher interface {
	// MatchString 判断字符串是否匹配
	MatchString(s string) bool
	// FindStringSubmatchIndex 返回第一个匹配及其分组的位置 [start0,end0, start1,end1, ...]
	FindStringSubmatchIndex(s string) []int
	// FindAllStringIndex 返回所有匹配的位置列表
	FindAllStringIndex(s string, n int) [][]int
	// FindAllStringSubmatchIndex 返回所有匹配及其分组位置
	FindAllStringSubmatchIndex(s string, n int) [][]int
	// ReplaceAllString 替换所有匹配
	ReplaceAllString(src, repl string) string
	// ReplaceAllStringFunc 用函数替换所有匹配
	ReplaceAllStringFunc(src string, repl func(string) string) string
}

// ---- Go regexp 包装 ----

type goMatcher struct{ re *regexp.Regexp }

func (m *goMatcher) MatchString(s string) bool { return m.re.MatchString(s) }
func (m *goMatcher) FindStringSubmatchIndex(s string) []int {
	return m.re.FindStringSubmatchIndex(s)
}
func (m *goMatcher) FindAllStringIndex(s string, n int) [][]int {
	return m.re.FindAllStringIndex(s, n)
}
func (m *goMatcher) FindAllStringSubmatchIndex(s string, n int) [][]int {
	return m.re.FindAllStringSubmatchIndex(s, n)
}
func (m *goMatcher) ReplaceAllString(src, repl string) string {
	return m.re.ReplaceAllString(src, repl)
}
func (m *goMatcher) ReplaceAllStringFunc(src string, repl func(string) string) string {
	return m.re.ReplaceAllStringFunc(src, repl)
}

// ---- regexp2 包装 ----

type r2Matcher struct{ re *regexp2.Regexp }

func (m *r2Matcher) MatchString(s string) bool {
	ok, _ := m.re.MatchString(s)
	return ok
}

func (m *r2Matcher) FindStringSubmatchIndex(s string) []int {
	match, err := m.re.FindStringMatch(s)
	if err != nil || match == nil {
		return nil
	}
	groups := match.Groups()
	loc := make([]int, 0, len(groups)*2)
	for _, g := range groups {
		if g.Length == 0 && g.Index == 0 && len(g.Captures) == 0 {
			loc = append(loc, -1, -1)
		} else {
			loc = append(loc, g.Index, g.Index+g.Length)
		}
	}
	return loc
}

func (m *r2Matcher) FindAllStringIndex(s string, n int) [][]int {
	var result [][]int
	match, err := m.re.FindStringMatch(s)
	for err == nil && match != nil {
		if n >= 0 && len(result) >= n {
			break
		}
		g := match.Groups()[0]
		result = append(result, []int{g.Index, g.Index + g.Length})
		match, err = m.re.FindNextMatch(match)
	}
	return result
}

func (m *r2Matcher) FindAllStringSubmatchIndex(s string, n int) [][]int {
	var result [][]int
	match, err := m.re.FindStringMatch(s)
	for err == nil && match != nil {
		if n >= 0 && len(result) >= n {
			break
		}
		groups := match.Groups()
		loc := make([]int, 0, len(groups)*2)
		for _, g := range groups {
			if len(g.Captures) == 0 {
				loc = append(loc, -1, -1)
			} else {
				loc = append(loc, g.Index, g.Index+g.Length)
			}
		}
		result = append(result, loc)
		match, err = m.re.FindNextMatch(match)
	}
	return result
}

func (m *r2Matcher) ReplaceAllString(src, repl string) string {
	// regexp2 的 Replace 参数：count=-1 表示全部替换
	result, err := m.re.Replace(src, repl, -1, -1)
	if err != nil {
		return src
	}
	return result
}

func (m *r2Matcher) ReplaceAllStringFunc(src string, repl func(string) string) string {
	var sb strings.Builder
	pos := 0
	match, err := m.re.FindStringMatch(src)
	for err == nil && match != nil {
		g0 := match.Groups()[0]
		// 追加匹配之前的部分
		sb.WriteString(src[pos:g0.Index])
		// 追加回调返回的替换字符串
		sb.WriteString(repl(g0.String()))
		pos = g0.Index + g0.Length
		match, err = m.re.FindNextMatch(match)
	}
	sb.WriteString(src[pos:])
	return sb.String()
}

// CompileAny 将 PHP 风格的正则表达式编译为 Matcher。
// 优先尝试 Go 原生 regexp；若编译失败（如含 lookahead/lookbehind），
// 则使用 regexp2（完整 PCRE 支持）。
func CompileAny(pattern string) (Matcher, error) {
	goPattern, r2Pattern, r2Flags := parsePhpPattern(pattern)

	// 尝试 Go 标准 regexp
	re, err := regexp.Compile(goPattern)
	if err == nil {
		return &goMatcher{re: re}, nil
	}

	// fallback: regexp2
	r2, err2 := regexp2.Compile(r2Pattern, r2Flags)
	if err2 != nil {
		return nil, err2
	}
	return &r2Matcher{re: r2}, nil
}

// stripExtendedWhitespace 近似实现 PHP /x 修饰符的行为：
// - 在字符类 [] 外，移除未转义的空白字符（空格、制表符等）
// 这里不处理 # 注释，因为当前项目中使用 /x 的模式主要依赖空白，而非行内注释。
func stripExtendedWhitespace(pattern string) string {
	var b strings.Builder
	inCharClass := false
	escaped := false

	for i := 0; i < len(pattern); i++ {
		ch := pattern[i]

		if escaped {
			b.WriteByte(ch)
			escaped = false
			continue
		}

		if ch == '\\' {
			escaped = true
			b.WriteByte(ch)
			continue
		}

		if ch == '[' && !inCharClass {
			inCharClass = true
			b.WriteByte(ch)
			continue
		}

		if ch == ']' && inCharClass {
			inCharClass = false
			b.WriteByte(ch)
			continue
		}

		// 在字符类外，移除空白字符
		if !inCharClass && (ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' || ch == '\f') {
			continue
		}

		b.WriteByte(ch)
	}

	return b.String()
}
