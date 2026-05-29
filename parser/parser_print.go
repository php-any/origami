package parser

import (
	"fmt"
	"os"
	"strings"

	"github.com/php-any/origami/data"
)

func (p *Parser) PrintDetailedError(err string, from data.From) {
	p.printDetailedError(err, from)
}

// printDetailedError 打印详细的错误信息
func (p *Parser) printDetailedError(err string, from data.From) {
	_, _ = fmt.Fprintln(os.Stderr, "\n"+strings.Repeat("=", 80))
	_, _ = fmt.Fprintln(os.Stderr, "🚨 解析错误")
	_, _ = fmt.Fprintln(os.Stderr, strings.Repeat("=", 80))

	if from == nil {
		_, _ = fmt.Fprintf(os.Stderr, "📍文件位置信息为空\n")
		// 显示错误消息
		_, _ = fmt.Fprintf(os.Stderr, "❌ 错误: %s\n", err)
		return
	}

	// 错误位置信息
	start, end := from.GetPosition()
	sl, sp := from.GetStartPosition()
	_, _ = fmt.Fprintf(os.Stderr, "📍 位置: 第 %d 行, 第 %d 列 (位置: %d-%d)\n", sl+1, sp+1, start, end)
	_, _ = fmt.Fprintf(os.Stderr, "📄 文件: %s:%d:%d\n", from.GetSource(), sl+1, sp+1)

	// 当前 token 信息
	currentToken := p.current()
	_, _ = fmt.Fprintf(os.Stderr, "🔍 当前 Token: %s (类型: %d)\n", currentToken.Literal(), currentToken.Type())

	// 显示错误消息
	_, _ = fmt.Fprintf(os.Stderr, "❌ 错误: %s\n", err)

	// 显示上下文（前后几个 token）
	_, _ = fmt.Fprintln(os.Stderr, "\n📝 上下文:")
	p.printContext()

	_, _ = fmt.Fprintln(os.Stderr, strings.Repeat("=", 80))
}

// printPHPUncaughtError 打印 PHP Uncaught Error（含 Stack trace，与 zend 一致）
func (p *Parser) printPHPUncaughtError(msg string, from data.From, frames []data.StackFrame) {
	if from == nil {
		_, _ = fmt.Fprintf(os.Stderr, "Fatal error: %s in Unknown on line 0\n", msg)
		return
	}
	sl, _ := from.GetStartPosition()
	file := from.GetSource()
	prefix := ""
	if data.HasUserOutput() {
		prefix = "\n"
	}
	_, _ = fmt.Fprintf(os.Stderr, "%sFatal error: %s in %s:%d\n", prefix, msg, file, sl+1)
	_, _ = fmt.Fprintln(os.Stderr, "Stack trace:")
	for i, frame := range frames {
		var stackSl int
		var source string
		if frame.From != nil {
			stackSl, _ = frame.From.GetStartPosition()
			source = frame.From.GetSource()
		}
		if frame.ClassName == "" {
			_, _ = fmt.Fprintf(os.Stderr, "#%d %s(%d): %s()\n", i, source, stackSl+1, frame.MethodName)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "#%d %s(%d): %s->%s()\n", i, source, stackSl+1, frame.ClassName, frame.MethodName)
		}
	}
	_, _ = fmt.Fprintf(os.Stderr, "#%d {main}\n", len(frames))
	_, _ = fmt.Fprintf(os.Stderr, "  thrown in %s on line %d\n", file, sl+1)
}

// printPHPCompileFatal 打印 PHP 编译期 Fatal error（与 zend 格式一致）
func (p *Parser) printPHPCompileFatal(msg string, from data.From) {
	prefix := ""
	if data.HasUserOutput() {
		prefix = "\n"
	}
	if from == nil {
		_, _ = fmt.Fprintf(os.Stderr, "%sFatal error: %s in Unknown on line 0\n", prefix, msg)
		return
	}
	sl, _ := from.GetStartPosition()
	_, _ = fmt.Fprintf(os.Stderr, "%sFatal error: %s in %s on line %d\n", prefix, msg, from.GetSource(), sl+1)
}

// printRuntimeError 打印运行时错误信息（例如数据库/IO等在执行阶段的异常）
func (p *Parser) printRuntimeError(err string, from data.From) {
	// 规范化错误文本：去掉前缀 "throw ", 分离 Caused by 段，避免重复展示
	normalized := strings.TrimSpace(err)
	normalized = strings.TrimPrefix(normalized, "throw ")
	mainMsg := normalized
	if idx := strings.Index(normalized, "\nCaused by: "); idx != -1 {
		mainMsg = strings.TrimSpace(normalized[:idx])
	} else if idx := strings.Index(normalized, "Caused by: "); idx != -1 {
		mainMsg = strings.TrimSpace(normalized[:idx])
	}

	if from == nil {
		_, _ = fmt.Fprintf(os.Stderr, "ZY Fatal error: %s in <unknown>:0\n", mainMsg)
		return
	}

	sl, sp := from.GetStartPosition()
	// 使用 path:line:col 形式，便于在大多数 IDE/终端中可点击跳转
	_, _ = fmt.Fprintf(os.Stderr, "ZY Fatal error: %s in %s:%d:%d\n", mainMsg, from.GetSource(), sl+1, sp+1)
}

// printContext 打印当前解析位置的上下文
func (p *Parser) printContext() {
	// 保存当前位置
	originalPos := p.position

	// 显示前3个token
	_, _ = fmt.Fprint(os.Stderr, "   前文: ")
	for i := 3; i > 0; i-- {
		if p.position-i >= 0 {
			token := p.tokens[p.position-i]
			_, _ = fmt.Fprintf(os.Stderr, "%s ", token.Literal())
		}
	}

	// 显示当前token（高亮）
	_, _ = fmt.Fprintf(os.Stderr, "\n   👉 当前: [%s] ", p.current().Literal())

	// 显示后3个token
	_, _ = fmt.Fprint(os.Stderr, "\n   后文: ")
	for i := 1; i <= 3; i++ {
		if p.position+i < len(p.tokens) {
			token := p.tokens[p.position+i]
			_, _ = fmt.Fprintf(os.Stderr, "%s ", token.Literal())
		}
	}
	_, _ = fmt.Fprintln(os.Stderr)

	// 恢复位置
	p.position = originalPos
}
