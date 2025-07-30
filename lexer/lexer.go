package lexer

import (
	"bufio"
	"github.com/php-any/origami/token"
	"unicode"
	"unicode/utf8"
)

// Position 表示一个位置
type Position struct {
	Line   int // 行号
	Column int // 列号
	Offset int // 字节偏移量
}

// Token 表示一个词法单元
type Token struct {
	Type    token.TokenType
	Literal string
	Start   int // 起始位置
	End     int // 结束位置
	Line    int // 行号
	Pos     int // 单独一行的位置
}

// Node 表示 DAG 中的一个节点
type Node struct {
	children map[rune]*Node
	token    *token.TokenDefinition
}

// Lexer 表示词法分析器
type Lexer struct {
	input      []rune        // 输入内容
	reader     *bufio.Reader // 输入源
	pos        *Position     // 当前位置
	ch         rune          // 当前字符
	width      int           // 当前字符的宽度
	hasNext    bool
	readOffset int   // 当前读取到 input 的下标
	root       *Node // DAG 根节点
}

// NewLexer 创建一个新的词法分析器
func NewLexer() *Lexer {
	lexer := &Lexer{
		root: &Node{
			children: make(map[rune]*Node),
		},
	}

	// 构建 DAG
	for _, def := range token.TokenDefinitions {
		lexer.addTokenDefinition(def)
	}

	return lexer
}

// addTokenDefinition 添加一个 token 定义到 DAG 中
func (l *Lexer) addTokenDefinition(def token.TokenDefinition) {
	current := l.root
	for _, char := range def.Literal {
		if _, exists := current.children[char]; !exists {
			current.children[char] = &Node{
				children: make(map[rune]*Node),
			}
		}
		current = current.children[char]
	}
	current.token = &def
}

// isWhitespace 检查字符是否是空白字符（除了换行符）
func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}

// Tokenize 将输入字符串转换为 token 列表
func (l *Lexer) Tokenize(input string) []Token {
	var tokens []Token
	pos := 0
	lastWasNewline := false
	line := 1
	linePos := -1

	for pos < len(input) {
		linePos++
		// 跳过空白字符，但保留换行符
		if isWhitespace(input[pos]) {
			pos++
			continue
		}
		// 跳过全角空格
		if pos+2 <= len(input) && input[pos] == 0xe3 && input[pos+1] == 0x80 && input[pos+2] == 0x80 {
			pos += 3
			continue
		}
		if input[pos] == '\n' {
			if !lastWasNewline {
				tokens = append(tokens, Token{
					Type:    token.NEWLINE,
					Literal: "\n",
					Start:   pos,
					End:     pos + 1,
					Line:    line,
					Pos:     linePos,
				})
				lastWasNewline = true
			}
			line++
			linePos = -1
			pos++
			continue
		}
		lastWasNewline = false

		// 处理特殊token
		if special, newPos, ok := HandleSpecialToken(input, pos); ok {
			tokens = append(tokens, Token{
				Type:    special.Type,
				Literal: special.Literal,
				Start:   pos,
				End:     newPos,
				Line:    line,
				Pos:     linePos,
			})
			pos = newPos
			linePos = newPos
			continue
		}

		// 尝试匹配最长的token
		if tokDef, length, ok := l.matchLongestToken(input, pos); ok {
			tokens = append(tokens, Token{
				Type:    tokDef.Type,
				Literal: tokDef.Literal,
				Start:   pos,
				End:     pos + length,
				Line:    line,
				Pos:     linePos,
			})
			pos += length
			linePos += length
			continue
		}

		// 获取当前位置的 rune
		r, size := utf8.DecodeRuneInString(input[pos:])
		if r == utf8.RuneError {
			// 处理无效的 UTF-8 序列
			tokens = append(tokens, Token{
				Type:    token.UNKNOWN,
				Literal: string(input[pos]),
				Start:   pos,
				End:     pos + 1,
				Line:    line,
				Pos:     linePos,
			})
			pos++
			linePos++
			continue
		}

		// 检查是否是标识符
		if unicode.IsLetter(r) || r == '_' || r >= 0x4e00 {
			start := pos
			pos += size // 移动到下一个字符

			for pos < len(input) {
				r, size := utf8.DecodeRuneInString(input[pos:])
				if r == utf8.RuneError {
					break
				}

				// 检查是否是分割符
				if IsDelimiter(r) {
					break
				}

				// 检查是否是有效的标识符字符
				if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '\\' && r < 0x4e00 {
					break
				}

				pos += size
			}

			tokens = append(tokens, Token{
				Type:    token.IDENTIFIER,
				Literal: input[start:pos],
				Start:   start,
				End:     pos,
				Line:    line,
				Pos:     linePos,
			})
			continue
		}

		// 如果无法匹配任何token，将当前字符作为未知token
		tokens = append(tokens, Token{
			Type:    token.UNKNOWN,
			Literal: string(r),
			Start:   pos,
			End:     pos + size,
			Line:    line,
			Pos:     linePos,
		})
		pos += size
		linePos += size
	}

	return NewPreprocessor(tokens).Process()
}

// matchLongestToken 尝试匹配最长的 token
func (l *Lexer) matchLongestToken(input string, pos int) (*token.TokenDefinition, int, bool) {
	// 获取当前位置的rune
	r, _ := utf8.DecodeRuneInString(input[pos:])
	if r == utf8.RuneError {
		return nil, 0, false
	}

	// 如果不是字母、下划线或中文字符开头，尝试匹配其他token
	if (!unicode.IsLetter(r) && r != '_' && r < 0x4e00) || IsDelimiter(r) {
		// 使用 DAG 进行高效匹配
		return l.matchTokenWithDAG(input, pos)
	}

	// 如果是标识符开头，使用 DAG 匹配关键字
	return l.matchKeywordWithDAG(input, pos)
}

// matchTokenWithDAG 使用 DAG 匹配 token
func (l *Lexer) matchTokenWithDAG(input string, pos int) (*token.TokenDefinition, int, bool) {
	current := l.root
	var longestMatch *token.TokenDefinition
	longestLength := 0
	currentPos := pos

	// 遍历输入字符串，在 DAG 中查找匹配
	for currentPos < len(input) {
		r, size := utf8.DecodeRuneInString(input[currentPos:])
		if r == utf8.RuneError {
			break
		}

		// 检查当前节点是否有子节点
		if child, exists := current.children[r]; exists {
			current = child
			currentPos += size

			// 如果当前节点有 token 定义，记录为可能的匹配
			if current.token != nil {
				longestMatch = current.token
				longestLength = currentPos - pos
			}
		} else {
			// 没有更多匹配，退出循环
			break
		}
	}

	if longestMatch != nil {
		return longestMatch, longestLength, true
	}
	return nil, 0, false
}

// matchKeywordWithDAG 使用 DAG 匹配关键字
func (l *Lexer) matchKeywordWithDAG(input string, pos int) (*token.TokenDefinition, int, bool) {
	current := l.root
	var longestMatch *token.TokenDefinition
	longestLength := 0
	currentPos := pos

	// 遍历输入字符串，在 DAG 中查找匹配
	for currentPos < len(input) {
		r, size := utf8.DecodeRuneInString(input[currentPos:])
		if r == utf8.RuneError {
			break
		}

		// 检查当前节点是否有子节点
		if child, exists := current.children[r]; exists {
			current = child
			currentPos += size

			// 如果当前节点有 token 定义，且是关键字类型，记录为可能的匹配
			if current.token != nil &&
				((current.token.Type >= token.KEYWORD_START && current.token.Type <= token.KEYWORD_END) ||
					(current.token.Type >= token.VALUE_START && current.token.Type <= token.VALUE_END)) {
				longestMatch = current.token
				longestLength = currentPos - pos
			}
		} else {
			// 没有更多匹配，退出循环
			break
		}
	}

	// 检查匹配的关键字后面是否还有更多标识符字符
	if longestMatch != nil {
		// 检查关键字后面是否还有更多字符
		if pos+longestLength < len(input) {
			nextRune, _ := utf8.DecodeRuneInString(input[pos+longestLength:])
			// 如果后面还有标识符字符，不匹配关键字
			if unicode.IsLetter(nextRune) || unicode.IsDigit(nextRune) || nextRune == '_' || nextRune >= 0x4e00 {
				return nil, 0, false
			}
		}
		return longestMatch, longestLength, true
	}
	return nil, 0, false
}
