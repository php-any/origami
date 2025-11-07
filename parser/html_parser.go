package parser

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/lexer"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type HtmlParser struct {
	*Parser
	// 属性解析器映射
	attributeParsers map[string]ParseAttribute
}

// isVoidHtmlTag 判断是否为无需闭合的 HTML 空元素标签
func isVoidHtmlTag(tag string) bool {
	switch strings.ToLower(tag) {
	case "area", "base", "br", "col", "embed", "hr", "img", "input", "link", "meta", "param", "source", "track", "wbr":
		return true
	default:
		return false
	}
}

func (h *HtmlParser) Parse() (data.GetValue, data.Control) {
	return h.parseHtmlContent()
}

// parseHtmlContent 解析HTML内容
func (h *HtmlParser) parseHtmlContent() (data.GetValue, data.Control) {
	tracker := h.StartTracking()

	// 跳过开始的 < 符号
	if h.checkPositionIs(0, token.LT) {
		h.next()
	}

	// 解析标签名
	tagName := h.parseTagName()
	if tagName == "" {
		return nil, data.NewErrorThrow(h.newFrom(), errors.New("HTML标签缺少标签名"))
	}

	// 解析属性
	attributes := make(map[string]node.HtmlAttributeValue)
	for !h.isEOF() && h.current().Type != token.GT && h.current().Type != token.QUO {
		// 只解析属性名称
		attrName, isCode, acl := h.parseAttributeName()
		if acl != nil {
			return nil, acl
		}
		if attrName == "" {
			continue
		}

		// 检查是否有等号
		if h.current().Type != token.ASSIGN {
			// 没有值的属性，如 disabled
			attributes[attrName] = node.NewAttrValueAdapter(h.FromCurrentToken(), attrName, data.NewBoolValue(true))
			continue
		}

		// 检查是否有专门的属性解析器
		if parser, exists := h.attributeParsers[attrName]; exists {
			// 使用专门的属性解析器
			attrValue, acl := parser.Parser(attrName, h.Parser)
			if acl != nil {
				return nil, acl
			}
			attributes[attrName] = attrValue
		} else {
			// 对于普通属性，使用默认解析
			h.nextAndCheck(token.ASSIGN)
			var attrValue data.GetValue

			if isCode {
				attrValue = h.parseAttributeValue()
			} else {
				attrValue = node.NewStringLiteral(h.FromCurrentToken(), h.current().Literal)
				h.next()
			}

			attributes[attrName] = node.NewAttrValueAdapter(h.FromCurrentToken(), attrName, attrValue)
		}
	}

	isSelfClosing := false
	if h.checkPositionIs(0, token.QUO) {
		h.next()
		if h.checkPositionIs(0, token.GT) {
			isSelfClosing = true
			h.next()
		} else {
			return nil, data.NewErrorThrow(h.newFrom(), errors.New("自闭合标签格式错误"))
		}
	} else if h.checkPositionIs(0, token.GT) {
		// 普通开始标签闭合 '>'
		h.next()
		if isVoidHtmlTag(tagName) {
			isSelfClosing = true
		}
	} else {
		return nil, data.NewErrorThrow(h.newFrom(), errors.New("HTML标签格式错误"))
	}

	// 提前处理 <script type="text/zy">：一次性读取到 </script>
	if strings.EqualFold(tagName, "script") {
		if attr, ok := attributes["type"]; ok {
			if v := attr.GetValue(); v != nil {
				if lit, ok := v.(*node.StringLiteral); ok && strings.EqualFold(lit.Value, "text/zy") {
					// 累积原始文本直到遇到 </script>
					tokens := make([]lexer.Token, 0)
					for !h.isEOF() {
						if h.current().Type == token.LT && h.checkPositionIs(1, token.QUO) && h.checkPositionIs(2, token.IDENTIFIER) && strings.EqualFold(h.peek(2).Literal, tagName) {
							// 消费结束标签 </script>
							h.next() // <
							h.next() // /
							h.next() // script
							if h.checkPositionIs(0, token.GT) {
								h.next()
							}
							break
						}
						tokens = append(tokens, h.current())
						h.next()
					}
					// 编译脚本为 Program
					prog, acl := h.Parser.ParserTokens(tokens, *h.Parser.source)
					if acl != nil {
						return nil, acl
					}
					return node.NewScriptZyNode(tracker.EndBefore(), prog), nil
				}
			}
		}
	}

	children := make([]data.GetValue, 0)
	if !isSelfClosing {
		var acl data.Control
		children, acl = h.parseHtmlChildren()
		if acl != nil {
			return nil, acl
		}

		// 查找结束标签
		if !h.findClosingTag(tagName) {
			return nil, data.NewErrorThrow(h.newFrom(), errors.New("HTML标签缺少结束标签: "+tagName))
		}
	}

	from := tracker.EndBefore()

	// 直接创建HTML节点，所有属性都已经在ParseAttribute中处理了
	return node.NewHtmlNode(
		from,
		tagName,
		attributes,
		children,
		isSelfClosing,
	), nil
}

// parseAttributeName 解析属性名
func (h *HtmlParser) parseAttributeName() (string, bool, data.Control) {
	// 收集属性名，支持连字符
	var name string
	var isCode bool
	if !h.isEOF() {
		// 添加第一个标识符
		if h.checkPositionIs(0, token.SEMICOLON) && h.current().Literal == "\n" {
			// 无效的换行而已
			h.next()
			return name, isCode, nil
		}
		if h.checkPositionIs(0, token.COLON) {
			isCode = true
			h.next()
		}
		current := h.current()
		// 基于原始字符判断属性名起始：首字符需为字母或下划线
		lit := current.Literal
		if lit == "" {
			return name, isCode, data.NewErrorThrow(h.newFrom(), fmt.Errorf("html属性名不能为空"))
		}
		first, _ := utf8.DecodeRuneInString(lit)
		if !(unicode.IsLetter(first) || first == '_') {
			return name, isCode, data.NewErrorThrow(h.newFrom(), fmt.Errorf("html属性必须以字母或下划线开头: %s", lit))
		}
		name = h.current().Literal
		h.next()

		// 检查是否有连字符，如果有则继续收集
		for !h.isEOF() && h.checkPositionIs(0, token.SUB) {
			name += "-"
			h.next()

			// 添加连字符后的标识符（包括关键字）
			if !h.isEOF() && h.checkPositionIs(0, token.IDENTIFIER, token.IF, token.ELSE) {
				name += h.current().Literal
				h.next()
			} else {
				// 如果连字符后没有标识符，停止收集
				break
			}
		}
	}

	return name, isCode, nil
}

// parseAttributeValue 解析属性值
func (h *HtmlParser) parseAttributeValue() data.GetValue {
	// 如果是标识符，直接作为字符串处理
	if h.checkPositionIs(0, token.IDENTIFIER) {
		value := h.current().Literal
		h.next()
		return node.NewStringLiteral(h.FromCurrentToken(), value)
	}

	// 尝试解析为表达式
	exprParser := NewExpressionParser(h.Parser)
	var acl data.Control
	attrValue, acl := exprParser.Parse()
	if acl != nil {
		// 如果表达式解析失败，尝试作为字符串处理
		var value string
		for !h.isEOF() && h.current().Type != token.GT && h.current().Type != token.QUO {
			value += h.current().Literal
			h.next()
		}
		return node.NewStringLiteral(h.FromCurrentToken(), value)
	}

	return attrValue
}

// parseTagName 解析标签名
func (h *HtmlParser) parseTagName() string {
	// 直接使用当前token的Literal作为标签名
	if !h.isEOF() {
		name := h.current().Literal
		h.next()
		return name
	}

	return ""
}

// parseHtmlChildren 解析HTML子内容
func (h *HtmlParser) parseHtmlChildren() ([]data.GetValue, data.Control) {
	var children = make([]data.GetValue, 0)

	for !h.isEOF() {
		// 检查是否是结束标签
		if h.current().Type == token.LT && h.checkPositionIs(1, token.QUO) {
			break
		}

		// 解析子节点
		child, acl := h.parseHtmlChild()
		if acl != nil {
			return nil, acl
		}
		if child != nil {
			if line, ok := child.(*node.StringLiteral); ok {
				if line.Value == "\n" {
					continue
				}
			}
			children = append(children, child)
		}

		// 防止无限循环：确保token位置有变化
		if h.isEOF() {
			break
		}
	}

	// 处理if-else-if-else链式连接
	var acl data.Control
	children, acl = h.processIfElseChain(children)
	if acl != nil {
		return nil, acl
	}

	return children, nil
}

// processIfElseChain 处理if-else-if-else链式连接
func (h *HtmlParser) processIfElseChain(children []data.GetValue) ([]data.GetValue, data.Control) {
	var result []data.GetValue
	var currentIfNode *node.HtmlIfNode
	if len(children) <= 1 {
		return children, nil
	}

	for _, child := range children {
		if ifNode, ok := child.(*node.HtmlNode); ok {
			// 检查是否是有效的条件节点
			ifNode, ok := h.isValidConditionNode(ifNode)
			if ok {
				if currentIfNode == nil {
					// 第一个节点必须是if节点
					if ifNode.Type == node.HtmlIfTypeIf {
						currentIfNode = ifNode
						result = append(result, currentIfNode)
					} else {
						// else-if或else节点不能作为第一个节点，返回错误
						return nil, data.NewErrorThrow(h.newFrom(), errors.New("else-if或else节点不能作为第一个节点"))
					}
				} else {
					// 开始新链
					if ifNode.Type == node.HtmlIfTypeIf {
						currentIfNode = ifNode
						result = append(result, currentIfNode)
					} else {
						currentIfNode.SetNextNode(ifNode)
						currentIfNode = ifNode
					}
				}
			} else {
				// 无效的条件节点，添加到结果中
				currentIfNode = nil
				result = append(result, child)
			}
		} else {
			result = append(result, child)
		}
	}

	return result, nil
}

// canConnectToPrevious 检查是否可以连接到前一个节点
func (h *HtmlParser) canConnectToPrevious(prev *node.HtmlIfNode, current *node.HtmlIfNode) bool {
	// 前一个节点必须是if或else-if
	if prev.Type != node.HtmlIfTypeIf && prev.Type != node.HtmlIfTypeElseIf {
		return false
	}

	// 当前节点必须是else-if或else
	if current.Type != node.HtmlIfTypeElseIf && current.Type != node.HtmlIfTypeElse {
		return false
	}

	return true
}

// isValidConditionNode 检查是否是有效的条件节点
func (h *HtmlParser) isValidConditionNode(ifNode *node.HtmlNode) (*node.HtmlIfNode, bool) {
	// 检查节点类型
	if con, ok := ifNode.Attributes[IfAttributeName]; ok {
		return node.NewHtmlIfNode(node.HtmlIfTypeIf, con.(*node.AttrIfValue), ifNode), true
	}
	if con, ok := ifNode.Attributes[ElseIfAttributeName]; ok {
		return node.NewHtmlIfNode(node.HtmlIfTypeElseIf, con.(*node.AttrIfValue), ifNode), true
	}
	if _, ok := ifNode.Attributes[ElseAttributeName]; ok {
		return node.NewHtmlIfNode(node.HtmlIfTypeElse, nil, ifNode), true
	}
	return nil, false
}

// parseHtmlChild 解析HTML子节点
func (h *HtmlParser) parseHtmlChild() (data.GetValue, data.Control) {
	if h.current().Type == token.LT {
		// 可能是HTML标签
		// DOCTYPE 的解析移动到 ExpressionParser.parseComparison()
		// 处理 <!-- 注释 -->，直到遇到 --> 才结束，避免被注释内的 > 提前终止
		if h.checkPositionIs(1, token.NOT) && h.checkPositionIs(2, token.DECR) {
			tracker := h.StartTracking()
			// 跳过 < ! - -
			h.nextAndCheck(token.LT)
			h.nextAndCheck(token.NOT)
			h.nextAndCheck(token.DECR)
			// 扫描直到 -- >
			str := "<!--"
			for !h.isEOF() {
				if h.checkPositionIs(0, token.DECR) && h.checkPositionIs(1, token.GT) {
					h.nextAndCheck(token.DECR)
					h.nextAndCheck(token.GT)
					break
				}
				str += h.current().Literal
				h.next()
			}
			// 注释不输出内容，但返回一个空字符串字面量，保持节点类型一致
			return node.NewStringLiteral(tracker.EndBefore(), str+"-->"), nil
		}
		if h.checkPositionIs(1, token.IDENTIFIER) {
			// 解析单个HTML标签（包括开始标签、属性和结束标签）
			return h.parseSingleHtmlTag()
		} else if h.checkPositionIs(1, token.QUO) {
			// 结束标签，停止解析
			return nil, nil
		}
	}

	// 解析文本内容
	return h.parseHtmlText()
}

// parseSingleHtmlTag 解析单个HTML标签（包括开始标签、属性和结束标签）
func (h *HtmlParser) parseSingleHtmlTag() (data.GetValue, data.Control) {
	tracker := h.StartTracking()

	// 跳过开始的 < 符号
	if h.checkPositionIs(0, token.LT) {
		h.next()
	}

	// 解析标签名
	tagName := h.parseTagName()
	if tagName == "" {
		return nil, data.NewErrorThrow(h.newFrom(), errors.New("HTML标签缺少标签名"))
	}

	// 解析属性
	attributes := make(map[string]node.HtmlAttributeValue)
	for !h.isEOF() && h.current().Type != token.GT && h.current().Type != token.QUO {
		// 只解析属性名称
		attrName, isCode, acl := h.parseAttributeName()
		if acl != nil {
			return nil, acl
		}
		if attrName == "" {
			continue
		}

		// 检查是否有等号
		if h.current().Type != token.ASSIGN {
			// 没有值的属性，如 disabled
			attributes[attrName] = node.NewAttrValueAdapter(h.FromCurrentToken(), attrName, data.NewBoolValue(true))
			continue
		}

		// 检查是否有专门的属性解析器
		if parser, exists := h.attributeParsers[attrName]; exists {
			// 使用专门的属性解析器
			attrValue, acl := parser.Parser(attrName, h.Parser)
			if acl != nil {
				return nil, acl
			}
			attributes[attrName] = attrValue
		} else {
			// 对于普通属性，使用默认解析
			h.next()
			var attrValue data.GetValue
			if isCode {
				attrValue = h.parseAttributeValue()
			} else {
				attrValue = node.NewStringLiteral(h.FromCurrentToken(), h.current().Literal)
				h.next()
			}
			attributes[attrName] = node.NewAttrValueAdapter(h.FromCurrentToken(), attrName, attrValue)
		}
	}

	isSelfClosing := false
	// void 标签无需显式自闭合，直接消费 '>' 并返回
	if isVoidHtmlTag(tagName) {
		if h.checkPositionIs(0, token.QUO) {
			h.next()
		}
		if h.checkPositionIs(0, token.GT) {
			h.next()
		}
		from := tracker.EndBefore()
		return node.NewHtmlNode(
			from,
			tagName,
			attributes,
			nil,
			true,
		), nil
	}
	if h.checkPositionIs(0, token.QUO) {
		h.next()
		if h.checkPositionIs(0, token.GT) {
			isSelfClosing = true
			h.next()
		} else {
			return nil, data.NewErrorThrow(h.newFrom(), errors.New("自闭合标签格式错误"))
		}
	} else if h.checkPositionIs(0, token.GT) {
		// 普通开始标签闭合 '>'
		h.next()
	} else {
		return nil, data.NewErrorThrow(h.newFrom(), errors.New("HTML标签格式错误"))
	}

	// 在进入子节点解析之前，优先处理 <script type="text/zy">：一次性读取源码到 </script>
	if !isSelfClosing && strings.EqualFold(tagName, "script") {
		if attr, ok := attributes["type"]; ok {
			if v := attr.GetValue(); v != nil {
				if lit, ok := v.(*node.StringLiteral); ok && strings.EqualFold(lit.Value, "text/zy") {
					// 累积原始文本直到遇到 </script>
					tokens := make([]lexer.Token, 0)
					for !h.isEOF() {
						if h.current().Type == token.LT && h.checkPositionIs(1, token.QUO) && h.checkPositionIs(2, token.IDENTIFIER) && strings.EqualFold(h.peek(2).Literal, tagName) {
							// 消费结束标签 </script>
							h.next() // <
							h.next() // /
							h.next() // script
							if h.checkPositionIs(0, token.GT) {
								h.next()
							}
							break
						}
						tokens = append(tokens, h.current())
						h.next()
					}

					// 编译脚本为 Program 并返回 ScriptZyNode
					prog, acl := h.Parser.ParserTokens(tokens, *h.Parser.source)
					if acl != nil {
						return nil, acl
					}
					return node.NewScriptZyNode(tracker.EndBefore(), prog), nil
				}
			}
		}
	}

	children := make([]data.GetValue, 0)

	if !isSelfClosing {
		var acl data.Control
		children, acl = h.parseHtmlChildren()
		if acl != nil {
			return nil, acl
		}

		// 查找结束标签
		if !h.findClosingTag(tagName) {
			return nil, data.NewErrorThrow(h.newFrom(), errors.New("HTML标签缺少结束标签: "+tagName))
		}
	}

	from := tracker.EndBefore()

	// 直接创建HTML节点，所有属性都已经在ParseAttribute中处理了
	return node.NewHtmlNode(
		from,
		tagName,
		attributes,
		children,
		isSelfClosing,
	), nil
}

// parseHtmlText 解析HTML文本内容，支持插值字符串
func (h *HtmlParser) parseHtmlText() (data.GetValue, data.Control) {
	var textParts []data.GetValue
	initialPos := h.GetStart()
	currentText := ""

	for !h.isEOF() && h.current().Type != token.LT {
		if h.position > 1 && h.current().Start != h.peek(-1).End {
			// token 之间有间隔，补空格
			currentText += " "
		}

		// 检查是否是插值开始 {$
		if h.current().Type == token.LBRACE && h.checkPositionIs(1, token.VARIABLE) {
			// 如果有累积的文本，先添加到结果中
			if currentText != "" {
				textParts = append(textParts, node.NewStringLiteral(h.FromCurrentToken(), currentText))
				currentText = ""
			}

			// 跳过 {
			h.nextAndCheck(token.LBRACE)

			// 收集并解析表达式 tokens 直到匹配的 }
			exprTokens, acl := h.collectExprTokensInsideBraces()
			if acl != nil {
				return nil, acl
			}

			expr, acl := h.parseExprFromTokens(exprTokens)
			if acl != nil {
				return nil, acl
			}

			// 结束的 }
			h.nextAndCheck(token.RBRACE)

			// 添加表达式到结果中
			textParts = append(textParts, expr)
		} else if h.current().Type == token.AT && h.checkPositionIs(1, token.LBRACE) && !h.checkPositionIs(-1, token.LBRACE) {
			// 处理函数/表达式插值 @{ ... } 但不能是 \@{
			// 输出累积文本
			if currentText != "" {
				textParts = append(textParts, node.NewStringLiteral(h.FromCurrentToken(), currentText))
				currentText = ""
			}

			// 跳过 '@' '{'
			h.nextAndCheck(token.AT)
			h.nextAndCheck(token.LBRACE)

			// 收集直到匹配的 '}'（支持嵌套），采集原始 tokens
			exprTokens, acl := h.collectExprTokensInsideBraces()
			if acl != nil {
				return nil, acl
			}
			// 结束的 }
			h.nextAndCheck(token.RBRACE)

			// 解析表达式（基于采集的 tokens）
			expr, acl := h.parseExprFromTokens(exprTokens)
			if acl != nil {
				return nil, acl
			}
			textParts = append(textParts, expr)
		} else {
			// 普通文本
			currentText += h.current().Literal
			h.next()
		}

		// 防止无限循环：检查位置是否变化
		if h.GetStart() == initialPos {
			break
		}
		initialPos = h.GetStart()
	}

	// 添加剩余的文本
	if currentText != "" {
		textParts = append(textParts, node.NewStringLiteral(h.FromCurrentToken(), currentText))
	}

	if len(textParts) == 0 {
		return nil, nil
	}

	// 如果只有一个部分，直接返回
	if len(textParts) == 1 {
		return textParts[0], nil
	}

	// 多个部分，需要创建字符串连接节点
	// 使用BinaryAdd节点连接多个部分
	result := textParts[0]
	for i := 1; i < len(textParts); i++ {
		result = node.NewBinaryAdd(h.FromCurrentToken(), result, textParts[i])
	}

	return result, nil
}

// collectExprTokensInsideBraces 收集当前位于 { 之后的表达式 tokens，直到匹配到对应的 }
func (h *HtmlParser) collectExprTokensInsideBraces() ([]lexer.Token, data.Control) {
	braceDepth := 1
	exprTokens := make([]lexer.Token, 0)
	for !h.isEOF() && braceDepth > 0 {
		if h.current().Type == token.LBRACE {
			braceDepth++
			exprTokens = append(exprTokens, h.current())
			h.next()
			continue
		}
		if h.current().Type == token.RBRACE {
			braceDepth--
			if braceDepth == 0 {
				break
			}
			exprTokens = append(exprTokens, h.current())
			h.next()
			continue
		}
		exprTokens = append(exprTokens, h.current())
		h.next()
	}
	if braceDepth != 0 {
		return nil, data.NewErrorThrow(h.newFrom(), errors.New("插值表达式缺少匹配的 }"))
	}
	return exprTokens, nil
}

// parseExprFromTokens 使用表达式解析器解析一段 tokens 为表达式值
func (h *HtmlParser) parseExprFromTokens(exprTokens []lexer.Token) (data.GetValue, data.Control) {
	// 为了保持与当前解析作用域一致（变量索引一致），临时切换 tokens 并用当前 Parser 的 ExpressionParser 解析
	originalTokens := h.Parser.tokens
	originalPosition := h.Parser.position
	originalSource := h.Parser.source

	h.Parser.tokens = exprTokens
	h.Parser.position = 0

	exprParser := NewExpressionParser(h.Parser)
	expr, ctl := exprParser.Parse()

	h.Parser.tokens = originalTokens
	h.Parser.position = originalPosition
	h.Parser.source = originalSource

	if ctl != nil {
		return nil, ctl
	}
	return expr, nil
}

// findClosingTag 查找结束标签
func (h *HtmlParser) findClosingTag(tagName string) bool {
	// 检查是否是结束标签
	if h.current().Type == token.LT && h.checkPositionIs(1, token.QUO) {
		h.next() // 跳过 <
		h.next() // 跳过 /

		// 检查标签名是否匹配
		if h.checkPositionIs(0, token.IDENTIFIER) && h.current().Literal == tagName {
			h.next()

			// 检查结束的 >
			if h.checkPositionIs(0, token.GT) {
				h.next()
				return true
			}
		}
	}

	return false
}

// trimSpace 去除字符串首尾空白
//func trimSpace(s string) string {
//	// 简单的空白字符去除
//	start := 0
//	end := len(s)
//
//	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
//		start++
//	}
//
//	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
//		end--
//	}
//
//	return s[start:end]
//}

func NewHtmlParser(parser *Parser) StatementParser {
	htmlParser := &HtmlParser{
		Parser:           parser,
		attributeParsers: make(map[string]ParseAttribute),
	}

	// 初始化属性解析器映射
	htmlParser.attributeParsers[IfAttributeName] = &HtmlIfAttributeParser{parser}
	htmlParser.attributeParsers[ElseIfAttributeName] = &HtmlElseIfAttributeParser{parser}
	htmlParser.attributeParsers[ElseAttributeName] = &HtmlElseAttributeParser{parser}
	htmlParser.attributeParsers[ForAttributeName] = &HtmlForAttributeParser{}

	return htmlParser
}

// ParseAttribute 解释属性时, 根据不同名称调用不同的
type ParseAttribute interface {
	Parser(name string, parser *Parser) (node.HtmlAttributeValue, data.Control)
}
