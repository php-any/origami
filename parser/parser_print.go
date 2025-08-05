package parser

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
)

// printDetailedError 打印详细的错误信息
func (p *Parser) printDetailedError(err string, from data.From) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("🚨 解析错误")
	fmt.Println(strings.Repeat("=", 80))

	// 错误位置信息
	start, end := from.GetPosition()
	fmt.Printf("📍 位置: 第 %d 行, 第 %d 列 (位置: %d-%d)\n", p.current().Line, p.current().Pos, start, end)
	fmt.Printf("📄 文件: %s\n", from.GetSource())

	// 当前 token 信息
	currentToken := p.current()
	fmt.Printf("🔍 当前 Token: %s (类型: %d)\n", currentToken.Literal, currentToken.Type)

	// 显示错误消息
	fmt.Printf("❌ 错误: %s\n", err)

	// 显示上下文（前后几个 token）
	fmt.Println("\n📝 上下文:")
	p.printContext()

	fmt.Println(strings.Repeat("=", 80))
}

// printContext 打印当前解析位置的上下文
func (p *Parser) printContext() {
	// 保存当前位置
	originalPos := p.position

	// 显示前3个token
	fmt.Print("   前文: ")
	for i := 3; i > 0; i-- {
		if p.position-i >= 0 {
			token := p.tokens[p.position-i]
			fmt.Printf("%s ", token.Literal)
		}
	}

	// 显示当前token（高亮）
	fmt.Printf("\n   👉 当前: [%s] ", p.current().Literal)

	// 显示后3个token
	fmt.Print("\n   后文: ")
	for i := 1; i <= 3; i++ {
		if p.position+i < len(p.tokens) {
			token := p.tokens[p.position+i]
			fmt.Printf("%s ", token.Literal)
		}
	}
	fmt.Println()

	// 恢复位置
	p.position = originalPos
}
