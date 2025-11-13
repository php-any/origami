package lexer

import (
	"fmt"
	"os"
	"strings"

	"github.com/php-any/origami/token"
)

// PrintTokens 以人类可读的格式打印 token 列表
func PrintTokens(tokens []Token, title string) {
	if len(tokens) == 0 {
		fmt.Fprintf(os.Stderr, "\n=== %s ===\n", title)
		fmt.Fprintf(os.Stderr, "  (无 tokens)\n")
		fmt.Fprintf(os.Stderr, "==================\n\n")
		return
	}

	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "%s\n", strings.Repeat("=", 80))
	fmt.Fprintf(os.Stderr, "  %s (共 %d 个 tokens)\n", title, len(tokens))
	fmt.Fprintf(os.Stderr, "%s\n", strings.Repeat("=", 80))

	// 打印表头
	fmt.Fprintf(os.Stderr, "%-6s | %-25s | %-60s\n",
		"索引", "类型", "字面值")
	fmt.Fprintf(os.Stderr, "%s\n", strings.Repeat("-", 100))

	// 打印每个 token
	for i, tok := range tokens {
		typeName := getTokenTypeName(tok.Type())
		literal := formatLiteral(tok.Literal())

		fmt.Fprintf(os.Stderr, "%-6d | %-25s | %-60s\n",
			i, typeName, literal)
	}

	fmt.Fprintf(os.Stderr, "%s\n", strings.Repeat("=", 80))
	fmt.Fprintf(os.Stderr, "\n")
}

// formatLiteral 格式化字面值，使其更易读
func formatLiteral(literal string) string {
	// 限制长度，避免过长
	if len(literal) > 80 {
		return literal[:77] + "..."
	}

	// 转义特殊字符
	literal = strings.ReplaceAll(literal, "\n", "\\n")
	literal = strings.ReplaceAll(literal, "\t", "\\t")
	literal = strings.ReplaceAll(literal, "\r", "\\r")

	// 如果包含不可见字符，用引号包裹
	if strings.ContainsAny(literal, " \t\n\r") {
		return fmt.Sprintf("%q", literal)
	}

	return literal
}

// getTokenTypeName 获取 token 类型的可读名称
func getTokenTypeName(t token.TokenType) string {
	switch t {
	// 关键字
	case token.IF:
		return "IF"
	case token.ELSE:
		return "ELSE"
	case token.ELSE_IF:
		return "ELSE_IF"
	case token.WHILE:
		return "WHILE"
	case token.FOR:
		return "FOR"
	case token.FOREACH:
		return "FOREACH"
	case token.FUNC:
		return "FUNC"
	case token.CLASS:
		return "CLASS"
	case token.PUBLIC:
		return "PUBLIC"
	case token.PRIVATE:
		return "PRIVATE"
	case token.PROTECTED:
		return "PROTECTED"
	case token.STATIC:
		return "STATIC"
	case token.NEW:
		return "NEW"
	case token.RETURN:
		return "RETURN"
	case token.ECHO:
		return "ECHO"
	case token.TRY:
		return "TRY"
	case token.CATCH:
		return "CATCH"
	case token.THROW:
		return "THROW"

	// 运算符
	case token.ADD:
		return "ADD (+)"
	case token.SUB:
		return "SUB (-)"
	case token.MUL:
		return "MUL (*)"
	case token.QUO:
		return "QUO (/)"
	case token.REM:
		return "REM (%)"
	case token.ASSIGN:
		return "ASSIGN (=)"
	case token.EQ:
		return "EQ (==)"
	case token.NE:
		return "NE (!=)"
	case token.LT:
		return "LT (<)"
	case token.GT:
		return "GT (>)"
	case token.LE:
		return "LE (<=)"
	case token.GE:
		return "GE (>=)"
	case token.OBJECT_OPERATOR:
		return "OBJECT_OPERATOR (->)"
	case token.SCOPE_RESOLUTION:
		return "SCOPE_RESOLUTION (::)"
	case token.DOT:
		return "DOT (.)"
	case token.COLON:
		return "COLON (:)"

	// 分隔符
	case token.COMMA:
		return "COMMA (,)"
	case token.SEMICOLON:
		return "SEMICOLON (;)"
	case token.LPAREN:
		return "LPAREN (()"
	case token.RPAREN:
		return "RPAREN ())"
	case token.LBRACE:
		return "LBRACE ({)"
	case token.RBRACE:
		return "RBRACE (})"
	case token.LBRACKET:
		return "LBRACKET ([)"
	case token.RBRACKET:
		return "RBRACKET (])"

	// 字面量
	case token.STRING:
		return "STRING"
	case token.INT:
		return "INT"
	case token.FLOAT:
		return "FLOAT"
	case token.TRUE:
		return "TRUE"
	case token.FALSE:
		return "FALSE"
	case token.NULL:
		return "NULL"

	// 标识符和变量
	case token.IDENTIFIER:
		return "IDENTIFIER"
	case token.VARIABLE:
		return "VARIABLE"

	// 特殊
	case token.NEWLINE:
		return "NEWLINE (\\n)"
	case token.EOF:
		return "EOF"
	case token.WHITESPACE:
		return "WHITESPACE"
	case token.COMMENT:
		return "COMMENT"
	case token.MULTILINE_COMMENT:
		return "MULTILINE_COMMENT"

	case token.INTERPOLATION_TOKEN:
		return "INTERPOLATION_TOKEN"
	case token.INTERPOLATION_VALUE:
		return "INTERPOLATION_VALUE"

	default:
		return fmt.Sprintf("Type%d", t)
	}
}
