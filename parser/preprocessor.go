package parser

import (
	"regexp"
	"strings"
)

// PHP 替代语法模式的正则表达式
var (
	// 结束关键词：endif; → } 等
	reEndif      = regexp.MustCompile(`\bendif\s*;`)
	reEndforeach = regexp.MustCompile(`\bendforeach\s*;`)
	reEndwhile   = regexp.MustCompile(`\bendwhile\s*;`)
	reEndfor     = regexp.MustCompile(`\bendfor\s*;`)
	reEndswitch  = regexp.MustCompile(`\bendswitch\s*;`)
	reElseColon  = regexp.MustCompile(`\belse\s*:`)
	// Blade @end 指令（HTML 中未编译的）
	reBladeEnd = regexp.MustCompile(`@end\w+`)

	// 控制语句冒号形式（替代语法开始标记），如 if(...): / elseif(...): / foreach(...):
	// 用于检测 Laravel Blade 编译产物中可能只含开始标记而没有 endif; 结束标记的情况。
)

// hasControlColon 判断代码中是否存在替代语法的"开始标记"（如 if(...): / foreach(...):）。
// 使用平衡括号匹配，避免把三元表达式中的冒号（if ($a ? $b : $c)）误判为替代语法冒号。
func hasControlColon(code string) bool {
	keywords := []string{"if", "elseif", "foreach", "while", "for", "switch"}
	for pos := 0; pos < len(code); {
		bestStart := -1
		for _, kw := range keywords {
			searchFrom := pos
			for {
				idx := strings.Index(strings.ToLower(code[searchFrom:]), kw)
				if idx == -1 {
					break
				}
				candidateStart := searchFrom + idx
				// 确保是单词边界；不满足则继续搜索同一关键字的后续出现
				if candidateStart > 0 && isIdentChar(code[candidateStart-1]) {
					searchFrom = candidateStart + len(kw)
					continue
				}
				if candidateStart+len(kw) < len(code) && isIdentChar(code[candidateStart+len(kw)]) {
					searchFrom = candidateStart + len(kw)
					continue
				}
				if bestStart == -1 || candidateStart < bestStart {
					bestStart = candidateStart
				}
				break
			}
		}
		if bestStart == -1 {
			break
		}
		// 从 bestStart 起查找控制关键字
		found := false
		for _, kw := range keywords {
			klen := len(kw)
			if bestStart+klen > len(code) {
				continue
			}
			if !strings.EqualFold(code[bestStart:bestStart+klen], kw) {
				continue
			}
			if bestStart > 0 && isIdentChar(code[bestStart-1]) {
				continue
			}
			if bestStart+klen < len(code) && isIdentChar(code[bestStart+klen]) {
				continue
			}
			// 函数/方法名中的关键字（如 forEach）不是控制语句，跳过
			if isFunctionNameDecl(code, bestStart) {
				continue
			}
			// 查找 ( 并跳过空白
			j := bestStart + klen
			for j < len(code) && (code[j] == ' ' || code[j] == '\t' || code[j] == '\n') {
				j++
			}
			if j >= len(code) || code[j] != '(' {
				continue
			}
			// 平衡括号匹配
			depth := 0
			k := j
			for k < len(code) {
				if code[k] == '(' {
					depth++
				} else if code[k] == ')' {
					depth--
					if depth == 0 {
						break
					}
				}
				k++
			}
			if depth != 0 {
				continue
			}
			// 检查 ) 后是否有 :
			l := k + 1
			for l < len(code) && (code[l] == ' ' || code[l] == '\t' || code[l] == '\n') {
				l++
			}
			if l < len(code) && code[l] == ':' {
				found = true
			}
			break
		}
		if found {
			return true
		}
		pos = bestStart + 1
	}
	return false
}

// ConvertAltPHPSyntaxForTest 暴露 convertAltPHPSyntax 供外部诊断使用。
func ConvertAltPHPSyntaxForTest(filename, content string) string {
	return convertAltPHPSyntax(filename, content)
}

// HasControlColonForTest 暴露 hasControlColon 供外部诊断使用。
func HasControlColonForTest(code string) bool {
	return hasControlColon(code)
}

// convertAltPHPSyntax 将 PHP 替代语法（if: endif; 等）转换为标准花括号语法。
// 仅在包含替代语法的文件中进行转换，安全跳过字符串和注释。
func convertAltPHPSyntax(filename, content string) string {
	_ = filename // 保留参数，可能用于日志

	// 快速检查：不包含替代语法的文件直接跳过。
	// 除结束标记（endif; 等）和 @end 指令外，还需检测替代语法的"开始标记"
	// （如 if(...): 带冒号形式），因为 Blade 编译产物可能只含开始标记而结束用花括号。
	if !strings.Contains(content, "endif;") &&
		!strings.Contains(content, "endforeach;") &&
		!strings.Contains(content, "endwhile;") &&
		!strings.Contains(content, "endfor;") &&
		!strings.Contains(content, "endswitch;") &&
		!strings.Contains(content, "else:") &&
		!strings.Contains(content, " @end") &&
		!strings.Contains(content, "\t@end") &&
		!hasControlColon(content) {
		return content
	}

	var result strings.Builder
	pos := 0

	for pos < len(content) {
		phpStart := strings.Index(content[pos:], "<?php")
		if phpStart == -1 {
			// 剩余 HTML：处理 @end 指令
			rest := content[pos:]
			if strings.Contains(rest, "@end") {
				rest = reBladeEnd.ReplaceAllString(rest, "<?php } ?>")
			}
			result.WriteString(rest)
			break
		}

		// HTML 内容（在 <?php 之前）
		html := content[pos : pos+phpStart]
		if strings.Contains(html, "@end") {
			html = reBladeEnd.ReplaceAllString(html, "<?php } ?>")
		}
		result.WriteString(html)
		pos += phpStart

		// 输出 <?php 标签
		result.WriteString("<?php")
		pos += 5

		// 找到关闭的 ?>
		phpEnd := strings.Index(content[pos:], "?>")
		if phpEnd == -1 {
			result.WriteString(convertPHPBlock(content[pos:]))
			break
		}

		result.WriteString(convertPHPBlock(content[pos : pos+phpEnd]))
		result.WriteString("?>")
		pos += phpEnd + 2
	}

	return result.String()
}

// convertPHPBlock 转换单个 PHP 代码块中的替代语法。
// 两阶段处理：
// 1. 控制关键词 (if/foreach/while/for/switch...:) 应用于整个块
// 2. 短模式 (endif;/else: 等) 只应用于代码区域（跳过字符串/注释）
func convertPHPBlock(block string) string {
	// 第一阶段：控制关键词转换应用到整个块
	// 这些模式需要完整的条件来匹配，不受字符串内容影响
	block = convertControlKeywords(block)

	// 第二阶段：短模式替换需要跳过字符串/注释
	block = convertShortPatterns(block)

	return block
}

// region 表示代码中的一段区间（用于标识字符串/注释等非代码区域）。
type region struct{ start, end int }

// convertShortPatterns 替换 endif;/else: 等短模式，跳过字符串/注释。
func convertShortPatterns(block string) string {
	var regions []region

	i := 0
	codeStart := 0

	for i < len(block) {
		switch {
		case block[i] == '\'' || block[i] == '"':
			if i > codeStart {
				regions = append(regions, region{codeStart, i})
			}
			quote := block[i]
			i++
			for i < len(block) {
				if block[i] == '\\' && i+1 < len(block) {
					i += 2
					continue
				}
				if block[i] == quote {
					i++
					break
				}
				i++
			}
			codeStart = i

		case i+1 < len(block) && block[i] == '/' && block[i+1] == '/':
			if i > codeStart {
				regions = append(regions, region{codeStart, i})
			}
			i += 2
			for i < len(block) && block[i] != '\n' {
				i++
			}
			codeStart = i

		case i+1 < len(block) && block[i] == '/' && block[i+1] == '*':
			if i > codeStart {
				regions = append(regions, region{codeStart, i})
			}
			i += 2
			for i+1 < len(block) && !(block[i] == '*' && block[i+1] == '/') {
				i++
			}
			if i+1 < len(block) {
				i += 2
			}
			codeStart = i

		case block[i] == '#' && (i == 0 || isLineStart(block[i-1])):
			if i > codeStart {
				regions = append(regions, region{codeStart, i})
			}
			i++
			for i < len(block) && block[i] != '\n' {
				i++
			}
			codeStart = i

		default:
			i++
		}
	}

	if i > codeStart {
		regions = append(regions, region{codeStart, i})
	}

	if len(regions) == 0 {
		return block
	}

	var result strings.Builder
	lastEnd := 0
	for _, r := range regions {
		result.WriteString(block[lastEnd:r.start])
		code := block[r.start:r.end]
		code = reEndif.ReplaceAllString(code, "}")
		code = reEndforeach.ReplaceAllString(code, "}")
		code = reEndwhile.ReplaceAllString(code, "}")
		code = reEndfor.ReplaceAllString(code, "}")
		code = reEndswitch.ReplaceAllString(code, "}")
		code = reElseColon.ReplaceAllString(code, "} else {")
		result.WriteString(code)
		lastEnd = r.end
	}
	result.WriteString(block[lastEnd:])

	return result.String()
}

// isLineStart 检查字符是否为行首或语句开始的分隔符
func isLineStart(c byte) bool {
	return c == '\n' || c == ';' || c == '{' || c == '}' || c == ' ' || c == '\t'
}

// isIdentChar 检查字符是否为 PHP 标识符字符
func isIdentChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

// isFunctionNameDecl 判断 code[start] 处的关键字是否是一个函数/方法名定义
// （前面紧跟 function 关键字），例如 forEach(callable $cb): void。
// 这种场景中的关键字不是控制语句，不应被当作替代语法处理。
func isFunctionNameDecl(code string, start int) bool {
	// 向前跳过空白，检查是否以 function 关键字结尾
	i := start
	for i > 0 && (code[i-1] == ' ' || code[i-1] == '\t' || code[i-1] == '\n') {
		i--
	}
	if i < 8 {
		return false
	}
	// 向前跳过空白后检查前面是否为 "function"
	j := i
	for j > 0 && (code[j-1] == ' ' || code[j-1] == '\t' || code[j-1] == '\n') {
		j--
	}
	if j >= 8 && strings.EqualFold(code[j-8:j], "function") {
		// 确认 function 前面是词边界
		if j-8 == 0 || !isIdentChar(code[j-9]) {
			return true
		}
	}
	return false
}

// convertControlKeywords 将 if(...):、foreach(...): 等转换为花括号语法。
// 使用平衡括号匹配正确处理嵌套括号。为避免误伤字符串/注释中的
// 伪控制语句（如 "<?php if(...): ?>"），先扫描并屏蔽字符串/注释内容，
// 转换完成后再还原。
func convertControlKeywords(code string) string {
	// 扫描字符串/注释区域，记录其 [start, end) 及原始内容。
	// isString 标记该区域是字符串字面量（保留首尾引号）还是注释（整体屏蔽）。
	type maskedRegion struct {
		start, end int
		content    string
		isString   bool
	}
	var masked []maskedRegion
	i := 0

	for i < len(code) {
		switch {
		case code[i] == '\'' || code[i] == '"':
			quote := code[i]
			strStart := i
			i++
			for i < len(code) {
				if code[i] == '\\' && i+1 < len(code) {
					i += 2
					continue
				}
				if code[i] == quote {
					i++
					break
				}
				i++
			}
			// 记录整个字符串字面量（含引号）为一个屏蔽区域。
			// 注意：空字符串（如 '' ）屏蔽后不含哨兵字节，若仍记录会
			// 导致还原阶段 mi 计数与 masked 错位，故必须跳过。
			if i-strStart > 2 {
				masked = append(masked, maskedRegion{strStart, i, code[strStart:i], true})
			}

		case i+1 < len(code) && code[i] == '/' && code[i+1] == '/':
			comStart := i
			i += 2
			for i < len(code) && code[i] != '\n' {
				i++
			}
			masked = append(masked, maskedRegion{comStart, i, code[comStart:i], false})

		case i+1 < len(code) && code[i] == '/' && code[i+1] == '*':
			comStart := i
			i += 2
			for i+1 < len(code) && !(code[i] == '*' && code[i+1] == '/') {
				i++
			}
			if i+1 < len(code) {
				i += 2
			}
			masked = append(masked, maskedRegion{comStart, i, code[comStart:i], false})

		case code[i] == '#' && (i == 0 || isLineStart(code[i-1])):
			comStart := i
			i++
			for i < len(code) && code[i] != '\n' {
				i++
			}
			masked = append(masked, maskedRegion{comStart, i, code[comStart:i], false})

		default:
			i++
		}
	}

	// 用哨兵字符（NUL）屏蔽字符串/注释内容，使关键字匹配只在真实代码上进行。
	// 字符串保留首尾引号（保证括号结构完整），只屏蔽中间内容；注释整体屏蔽。
	// 转换后再按哨兵位置依次还原，从而避免转换导致的位置偏移。
	// 注意：扫描索引基于字节，故统一用 []byte 操作以保证位置一致。
	const sentinel = '\x00'
	buf := []byte(code)
	for _, r := range masked {
		if r.content == "" {
			continue
		}
		if r.isString && r.end-r.start >= 2 {
			// 保留首尾引号，屏蔽中间内容
			for j := r.start + 1; j < r.end-1; j++ {
				buf[j] = sentinel
			}
		} else {
			for j := r.start; j < r.end; j++ {
				buf[j] = sentinel
			}
		}
	}
	maskedCode := string(buf)

	// 在屏蔽后的代码上做 if(...): 转换
	converted := convertControlKeywordsInCode(maskedCode)

	// 还原：按顺序把连续的哨兵序列替换回字符串/注释内容。
	cb := []byte(converted)
	mi := 0
	j := 0
	for j < len(cb) {
		if cb[j] != sentinel {
			j++
			continue
		}
		// 找到连续哨兵序列 [j, k)
		k := j
		for k < len(cb) && cb[k] == sentinel {
			k++
		}
		if mi < len(masked) {
			repl := []byte(masked[mi].content)
			// 字符串区域保留了引号，因此替换内容要去掉首尾引号，保留引号本身。
			if masked[mi].isString && len(repl) >= 2 {
				repl = repl[1 : len(repl)-1]
			}
			// 哨兵序列长度应恰好等于要替换的内容长度
			n := k - j
			if len(repl) == n {
				copy(cb[j:k], repl)
			}
			mi++
		}
		j = k
	}

	return string(cb)
}

// convertControlKeywordsInCode 对纯代码区域（不含字符串/注释）进行 if(...): 等转换。
func convertControlKeywordsInCode(result string) string {
	keywords := []string{"if", "elseif", "foreach", "while", "for", "switch"}

	pos := 0

	for pos < len(result) {
		bestStart := -1
		bestKW := ""
		bestParenEnd := 0
		bestColonPos := 0

		for _, kw := range keywords {
			searchFrom := pos
			for {
				idx := strings.Index(strings.ToLower(result[searchFrom:]), kw)
				if idx == -1 {
					break
				}
				candidateStart := searchFrom + idx

				// 确保是单词边界
				if candidateStart > 0 && isIdentChar(result[candidateStart-1]) {
					searchFrom = candidateStart + len(kw)
					continue
				}
				if candidateStart+len(kw) < len(result) && isIdentChar(result[candidateStart+len(kw)]) {
					searchFrom = candidateStart + len(kw)
					continue
				}
				// 函数/方法名中的关键字（如 forEach）不是控制语句，跳过
				if isFunctionNameDecl(result, candidateStart) {
					searchFrom = candidateStart + len(kw)
					continue
				}

				// 查找 ( 并跳过空白
				j := candidateStart + len(kw)
				for j < len(result) && (result[j] == ' ' || result[j] == '\t' || result[j] == '\n') {
					j++
				}
				if j >= len(result) || result[j] != '(' {
					searchFrom = candidateStart + len(kw)
					continue
				}

				// 平衡括号匹配
				depth := 0
				k := j
				for k < len(result) {
					if result[k] == '(' {
						depth++
					} else if result[k] == ')' {
						depth--
						if depth == 0 {
							break
						}
					}
					k++
				}
				if depth != 0 {
					searchFrom = candidateStart + len(kw)
					continue
				}

				// 检查 ) 后是否有 :
				l := k + 1
				for l < len(result) && (result[l] == ' ' || result[l] == '\t' || result[l] == '\n') {
					l++
				}
				if l >= len(result) || result[l] != ':' {
					searchFrom = candidateStart + len(kw)
					continue
				}

				if bestStart == -1 || candidateStart < bestStart {
					bestStart = candidateStart
					bestKW = kw
					bestParenEnd = k
					bestColonPos = l
				}
				break
			}
		}

		if bestStart == -1 {
			break
		}

		var converted strings.Builder
		converted.WriteString(result[:bestStart])

		if bestKW == "elseif" {
			converted.WriteString("} elseif (")
		} else {
			converted.WriteString(bestKW + " (")
		}
		// 找到 ( 的位置
		j := bestStart + len(bestKW)
		for j < len(result) && (result[j] == ' ' || result[j] == '\t' || result[j] == '\n') {
			j++
		}
		converted.WriteString(result[j+1 : bestParenEnd])
		converted.WriteString(") {")

		converted.WriteString(result[bestColonPos+1:])
		result = converted.String()
		pos = bestStart + 1
	}

	return result
}
