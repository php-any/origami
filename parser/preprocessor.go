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
)

// convertAltPHPSyntax 将 PHP 替代语法（if: endif; 等）转换为标准花括号语法。
// 仅在包含替代语法的文件中进行转换，安全跳过字符串和注释。
func convertAltPHPSyntax(filename, content string) string {
	_ = filename // 保留参数，可能用于日志

	// 快速检查：不包含替代语法的文件直接跳过
	if !strings.Contains(content, "endif;") &&
		!strings.Contains(content, "endforeach;") &&
		!strings.Contains(content, "endwhile;") &&
		!strings.Contains(content, "endfor;") &&
		!strings.Contains(content, "endswitch;") &&
		!strings.Contains(content, "else:") &&
		!strings.Contains(content, " @end") &&
		!strings.Contains(content, "\t@end") {
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

// convertShortPatterns 替换 endif;/else: 等短模式，跳过字符串/注释。
func convertShortPatterns(block string) string {
	type region struct{ start, end int }
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

// convertControlKeywords 将 if(...):、foreach(...): 等转换为花括号语法。
// 使用平衡括号匹配正确处理嵌套括号。
func convertControlKeywords(code string) string {
	keywords := []string{"if", "elseif", "foreach", "while", "for", "switch"}

	result := code
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
