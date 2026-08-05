package preg

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/php-any/origami/data"
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

// convertRecursivePatterns 将 PCRE 递归子模式替换为非递归等价物。
// 处理: (?R), (?-N), (?+N), (?&name), (?P>name)
// 替换为 [^()]* 以支持常见嵌套括号场景（不完美但实用）。
func convertRecursivePatterns(pattern string) string {
	result := strings.ReplaceAll(pattern, "(?R)", "[^()]*")
	result = regexp.MustCompile(`\(\?-\d+\)`).ReplaceAllString(result, "[^()]*")
	result = regexp.MustCompile(`\(\?\+\d+\)`).ReplaceAllString(result, "[^()]*")
	result = regexp.MustCompile(`\(\?&\w+\)`).ReplaceAllString(result, "[^()]*")
	result = regexp.MustCompile(`\(\?P>\w+\)`).ReplaceAllString(result, "[^()]*")
	return result
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
	// 括号风格的分隔符映射（PHP PCRE 支持）
	closingDelimiter := delimiter
	switch delimiter {
	case '{':
		closingDelimiter = '}'
	case '(':
		closingDelimiter = ')'
	case '[':
		closingDelimiter = ']'
	case '<':
		closingDelimiter = '>'
	}
	endIndex := -1
	// 从尾部向前查找未转义的分隔符
	for i := len(pattern) - 1; i > 0; i-- {
		if pattern[i] == closingDelimiter && pattern[i-1] != '\\' {
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

	// 处理 PCRE 递归子模式 (?R) (?-1) (?+1) (?&name) 等，
	// Go 的 regexp/regexp2 不支持递归匹配，
	// 替换为匹配非括号字符 [^()]* 以支持常见嵌套括号场景。
	regexBody = convertRecursivePatterns(regexBody)

	// 移除 PCRE 动词如 (*UTF8) (*UCP) 等，regexp2 不需要这些
	pcreVerbs := []string{"(*UTF8)", "(*UCP)", "(*UTF)"}
	for _, verb := range pcreVerbs {
		regexBody = strings.Replace(regexBody, verb, "", 1)
	}

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

// matcherGroupCount returns the number of match groups (full match + captures).
func matcherGroupCount(m Matcher) int {
	switch m := m.(type) {
	case *goMatcher:
		return m.re.NumSubexp() + 1
	case *r2Matcher:
		return len(m.re.GetGroupNumbers())
	default:
		return 1
	}
}

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

// Capture 表示单个捕获分组（含命名分组）。
type Capture struct {
	Text         string
	Participated bool
	Name         string
}

// FindCaptures 在 subject[offset:] 上执行匹配，返回全部分组（含未参与分组）。
func FindCaptures(m Matcher, subject string, offset int, anchored bool) []Capture {
	if offset < 0 {
		offset = 0
	}
	if offset > len(subject) {
		return nil
	}
	search := subject[offset:]

	switch matcher := m.(type) {
	case *goMatcher:
		return findGoCaptures(matcher.re, search, anchored)
	case *r2Matcher:
		return findR2Captures(matcher.re, search, anchored)
	default:
		loc := FindSubmatchAt(m, subject, offset, anchored)
		if loc == nil {
			return nil
		}
		caps := make([]Capture, 0, len(loc)/2)
		for i := 0; i < len(loc); i += 2 {
			start, end := loc[i], loc[i+1]
			cap := Capture{Participated: start >= 0}
			if cap.Participated {
				cap.Text = subject[start:end]
			}
			caps = append(caps, cap)
		}
		return caps
	}
}

func findGoCaptures(re *regexp.Regexp, search string, anchored bool) []Capture {
	loc := re.FindStringSubmatchIndex(search)
	if loc == nil {
		return nil
	}
	if anchored && loc[0] != 0 {
		return nil
	}
	submatch := re.FindStringSubmatch(search)
	names := re.SubexpNames()
	caps := make([]Capture, 0, len(submatch))
	for i, text := range submatch {
		cap := Capture{
			Text:         text,
			Participated: i < len(loc) && loc[i*2] >= 0,
		}
		if i < len(names) {
			cap.Name = names[i]
		}
		caps = append(caps, cap)
	}
	return caps
}

func findR2Captures(re *regexp2.Regexp, search string, anchored bool) []Capture {
	match, err := re.FindStringMatch(search)
	if err != nil || match == nil {
		return nil
	}
	groups := match.Groups()
	if anchored && len(groups) > 0 && groups[0].Index != 0 {
		return nil
	}
	caps := make([]Capture, 0, len(groups))
	for _, g := range groups {
		// regexp2 未参与的可选组：Index=0、Length=0 且无 Captures
		participated := !(g.Length == 0 && g.Index == 0 && len(g.Captures) == 0)
		cap := Capture{
			Text:         g.String(),
			Participated: participated,
			Name:         g.Name,
		}
		caps = append(caps, cap)
	}
	return caps
}

// BuildMatchArray 构造 PHP preg_match 的 $matches 数组（数值键 + 命名键）。
// 默认：未参与匹配的尾部可选捕获组不写入数组（与 PHP isset($matches[n]) 语义一致，ProgressBar 依赖此点）。
// PREG_UNMATCHED_AS_NULL：保留全部捕获组，未参与的为 null（brick/math BigNumber::of 依赖此点）。
// 命名捕获与 PHP 一致：先写入命名键，再写入同值的数值键（Illuminate 路由绑定依赖 array_slice 后仍能看到命名键）。
func BuildMatchArray(captures []Capture, flags int) data.Value {
	unmatchedAsNull := flags&512 != 0 // PREG_UNMATCHED_AS_NULL

	if !unmatchedAsNull {
		last := -1
		for i, cap := range captures {
			if cap.Participated {
				last = i
			}
		}
		if last < 0 {
			return &data.ArrayValue{List: []*data.ZVal{}}
		}
		captures = captures[:last+1]
	} else if len(captures) == 0 {
		return &data.ArrayValue{List: []*data.ZVal{}}
	}

	hasNamed := false
	for _, cap := range captures {
		if cap.Name != "" {
			hasNamed = true
			break
		}
	}

	list := make([]*data.ZVal, 0, len(captures)*2)
	for i, cap := range captures {
		var val data.Value
		switch {
		case !cap.Participated && unmatchedAsNull:
			val = data.NewNullValue()
		default:
			val = data.NewStringValue(cap.Text)
		}

		if cap.Name != "" {
			list = append(list, data.NewNamedZVal(cap.Name, val))
		}
		if hasNamed {
			// 混入命名键后 List 下标不再等于 PHP 整数键，必须显式编码。
			list = append(list, data.NewNamedZVal(data.IntArrayKeyName(i), val))
		} else {
			list = append(list, data.NewZVal(val))
		}
	}
	return &data.ArrayValue{List: list}
}

// ExpandPhpReplacement 将 PHP preg_replace 替换串中的反引用展开为最终文本。
// 支持 $n / ${n} / \n；\\ 表示字面反斜杠（OutputWrapper 的 '\\1'、escape 的 '$1\\\\$2' 依赖此语义）。
func ExpandPhpReplacement(repl string, groups []string) string {
	var b strings.Builder
	for i := 0; i < len(repl); i++ {
		c := repl[i]
		switch c {
		case '\\':
			if i+1 >= len(repl) {
				b.WriteByte('\\')
				break
			}
			n := repl[i+1]
			if n >= '0' && n <= '9' {
				num, next := readReplacementIndex(repl, i+1)
				b.WriteString(groupText(groups, num))
				i = next - 1
				continue
			}
			// \\ → 一个反斜杠；其它 \x → 字面 x（与 PHP 常见行为接近）
			if n == '\\' {
				b.WriteByte('\\')
				i++
				continue
			}
			b.WriteByte(n)
			i++
		case '$':
			if i+1 >= len(repl) {
				b.WriteByte('$')
				break
			}
			n := repl[i+1]
			if n == '$' {
				b.WriteByte('$')
				i++
				continue
			}
			if n == '{' {
				end := strings.IndexByte(repl[i+2:], '}')
				if end >= 0 {
					numStr := repl[i+2 : i+2+end]
					if num, err := strconv.Atoi(numStr); err == nil {
						b.WriteString(groupText(groups, num))
						i = i + 2 + end
						continue
					}
				}
				b.WriteByte('$')
				continue
			}
			if n >= '0' && n <= '9' {
				num, next := readReplacementIndex(repl, i+1)
				b.WriteString(groupText(groups, num))
				i = next - 1
				continue
			}
			b.WriteByte('$')
		default:
			b.WriteByte(c)
		}
	}
	return b.String()
}

func readReplacementIndex(s string, start int) (num int, next int) {
	num = 0
	i := start
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		num = num*10 + int(s[i]-'0')
		i++
		// PHP 只取合理位数的分组号；此处允许多位
	}
	return num, i
}

func groupText(groups []string, n int) string {
	if n < 0 || n >= len(groups) {
		return ""
	}
	return groups[n]
}

// ReplaceAllPhp 按 PHP 语义执行替换（展开 \1 / $1，支持 limit）。
func ReplaceAllPhp(m Matcher, src, repl string, limit int) (string, int) {
	all := m.FindAllStringSubmatchIndex(src, -1)
	if len(all) == 0 {
		return src, 0
	}
	max := len(all)
	if limit > 0 && max > limit {
		max = limit
	}

	var b strings.Builder
	pos := 0
	count := 0
	for mi := 0; mi < max; mi++ {
		loc := all[mi]
		if len(loc) < 2 || loc[0] < 0 {
			continue
		}
		b.WriteString(src[pos:loc[0]])
		groups := make([]string, 0, len(loc)/2)
		for g := 0; g < len(loc); g += 2 {
			if loc[g] >= 0 && loc[g+1] >= 0 {
				groups = append(groups, src[loc[g]:loc[g+1]])
			} else {
				groups = append(groups, "")
			}
		}
		b.WriteString(ExpandPhpReplacement(repl, groups))
		pos = loc[1]
		count++
	}
	b.WriteString(src[pos:])
	return b.String(), count
}

// HasModifier 检查 PHP 正则是否带有指定修饰符（如 A、i、m）。
func HasModifier(pattern string, mod byte) bool {
	if len(pattern) < 2 {
		return false
	}
	delimiter := pattern[0]
	closingDelimiter := delimiter
	switch delimiter {
	case '{':
		closingDelimiter = '}'
	case '(':
		closingDelimiter = ')'
	case '[':
		closingDelimiter = ']'
	case '<':
		closingDelimiter = '>'
	}
	endIndex := -1
	for i := len(pattern) - 1; i > 0; i-- {
		if pattern[i] == closingDelimiter && pattern[i-1] != '\\' {
			endIndex = i
			break
		}
	}
	if endIndex == -1 {
		return false
	}
	modifiers := pattern[endIndex+1:]
	return strings.Contains(modifiers, string(mod))
}

// FindSubmatchAt 在 subject 的 offset 位置起搜索，返回相对于完整 subject 的分组下标。
// anchored 为 true 时（PHP /A 修饰符），匹配必须从 offset 处开始。
func FindSubmatchAt(m Matcher, subject string, offset int, anchored bool) []int {
	if offset < 0 {
		offset = 0
	}
	if offset > len(subject) {
		return nil
	}
	search := subject[offset:]
	loc := m.FindStringSubmatchIndex(search)
	if loc == nil {
		return nil
	}
	if anchored && loc[0] != 0 {
		return nil
	}
	for i := 0; i < len(loc); i++ {
		if loc[i] >= 0 {
			loc[i] += offset
		}
	}
	return loc
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
