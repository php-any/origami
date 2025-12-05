package lexer

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/php-any/origami/token"
)

// TokenizeTemplate 解析包含 <?php ?> 的模板文件
func (l *Lexer) TokenizeTemplate(input string) []Token {
	var tokens []Token
	pos := 0
	line := 0
	linePos := 0
	lastWasNewline := false

	for pos < len(input) {
		// 1. 查找 <?php
		idx := strings.Index(input[pos:], "<?php")

		if idx == -1 {
			// 剩余全部是 HTML
			if pos < len(input) {
				content := input[pos:]
				tokens = append(tokens, NewWorkerToken(token.HTML_TAG, content, pos, len(input), line, linePos))
			}
			break
		}

		// 2. 处理 <?php 之前的内容 (HTML)
		if idx > 0 {
			content := input[pos : pos+idx]
			tokens = append(tokens, NewWorkerToken(token.HTML_TAG, content, pos, pos+idx, line, linePos))

			// 更新行号信息
			for _, r := range content {
				if r == '\n' {
					line++
					linePos = 0
				} else {
					linePos++
				}
			}
			pos += idx
		}

		// 3. 添加 START_TAG (<?php)
		// tokens = append(tokens, NewWorkerToken(token.START_TAG, "<?php", pos, pos+5, line, linePos))
		pos += 5
		linePos += 5

		// 4. Script 模式：复用原有的分词逻辑，直到遇到 ?>
		for pos < len(input) {
			// 检查是否是 ?> (END_TAG)
			if pos+2 <= len(input) && input[pos:pos+2] == "?>" {
				// tokens = append(tokens, NewWorkerToken(token.END_TAG, "?>", pos, pos+2, line, linePos))
				pos += 2
				linePos += 2
				break // 退出 Script 模式，回到 HTML 模式
			}

			// 跳过空白字符，但保留换行符
			if isWhitespace(input[pos]) {
				pos++
				linePos++
				continue
			}
			// 跳过全角空格
			if pos+3 <= len(input) && input[pos] == 0xe3 && input[pos+1] == 0x80 && input[pos+2] == 0x80 {
				pos += 3
				linePos += 3
				continue
			}
			if input[pos] == '\n' {
				if !lastWasNewline {
					tokens = append(tokens, NewWorkerToken(
						token.NEWLINE,
						"\n",
						pos,
						pos+1,
						line,
						linePos,
					))
					lastWasNewline = true
				}
				line++
				linePos = 0
				pos++
				continue
			}
			lastWasNewline = false

			// 处理特殊token
			if result, ok := HandleSpecialToken(input, pos, line, linePos); ok {
				tokens = append(tokens, NewWorkerToken(
					result.Token.Type,
					result.Token.Literal,
					pos,
					result.NewPos,
					line,
					linePos,
				))
				pos = result.NewPos
				line = result.NewLine
				linePos = result.NewLinePos
				continue
			}

			// 尝试匹配最长的token
			if tokDef, length, ok := l.matchLongestToken(input, pos); ok {
				tokens = append(tokens, NewWorkerToken(
					tokDef.Type,
					tokDef.Literal,
					pos,
					pos+length,
					line,
					linePos,
				))
				pos += length
				linePos += length
				continue
			}

			// 获取当前位置的 rune
			r, size := utf8.DecodeRuneInString(input[pos:])
			if r == utf8.RuneError {
				tokens = append(tokens, NewWorkerToken(
					token.UNKNOWN,
					string(input[pos]),
					pos,
					pos+1,
					line,
					linePos,
				))
				pos++
				linePos++
				continue
			}

			// 检查是否是标识符
			if unicode.IsLetter(r) || r == '_' || r >= 0x4e00 {
				start := pos
				startLinePos := linePos
				pos += size
				linePos += size

				for pos < len(input) {
					r, size := utf8.DecodeRuneInString(input[pos:])
					if r == utf8.RuneError {
						break
					}

					// 检查是否是分割符
					if IsDelimiter(r) {
						break
					}

					// 检查是否遇到 ?>
					if r == '?' && pos+1 < len(input) && input[pos+1] == '>' {
						break
					}

					if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '\\' && r < 0x4e00 {
						break
					}

					pos += size
					linePos += size
				}

				tokens = append(tokens, NewWorkerToken(
					token.IDENTIFIER,
					input[start:pos],
					start,
					pos,
					line,
					startLinePos,
				))
				continue
			}

			// 未知token
			tokens = append(tokens, NewWorkerToken(
				token.UNKNOWN,
				string(r),
				pos,
				pos+size,
				line,
				linePos,
			))
			pos += size
			linePos += size
		}
	}

	return NewPreprocessor(tokens).Process()
}
