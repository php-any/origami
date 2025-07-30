package parser

import (
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

type HtmlParser struct {
	*Parser
}

func (h *HtmlParser) Parse() (data.GetValue, data.Control) {
	// 解析HTML标签
	htmlNode, acl := h.parseHtmlContent()
	if acl != nil {
		return nil, acl
	}

	return htmlNode, nil
}

// parseHtmlContent 解析HTML内容（支持自闭合和非自闭合标签）
func (h *HtmlParser) parseHtmlContent() (data.GetValue, data.Control) {
	start := h.GetStart()

	// 跳过开始的 < 符号
	if h.current().Type == token.LT {
		h.next()
	}

	// 解析标签名
	tagName := h.parseTagName()
	if tagName == "" {
		return nil, data.NewErrorThrow(h.newFrom(), errors.New("HTML标签缺少标签名"))
	}

	// 解析属性
	attributes := make(map[string]data.GetValue)
	for !h.isEOF() && h.current().Type != token.GT && h.current().Type != token.QUO {
		attrName, attrValue, acl := h.parseAttribute()
		if acl != nil {
			return nil, acl
		}
		if attrName != "" {
			attributes[attrName] = attrValue
		}
	}

	isSelfClosing := false
	if h.current().Type == token.QUO {
		h.next()
		if h.current().Type == token.GT {
			isSelfClosing = true
			h.next()
		} else {
			return nil, data.NewErrorThrow(h.newFrom(), errors.New("自闭合标签格式错误"))
		}
	} else if h.current().Type == token.GT {
		h.next()
	} else {
		return nil, data.NewErrorThrow(h.newFrom(), errors.New("HTML标签格式错误"))
	}

	if !isSelfClosing {
		children, acl := h.parseHtmlChildren()
		if acl != nil {
			return nil, acl
		}

		// 查找结束标签
		if !h.findClosingTag(tagName) {
			return nil, data.NewErrorThrow(h.newFrom(), errors.New("HTML标签缺少结束标签: "+tagName))
		}

		return node.NewHtmlNode(
			h.NewTokenFrom(start),
			tagName,
			attributes,
			children,
			isSelfClosing,
		), nil
	}

	return node.NewHtmlNode(
		h.NewTokenFrom(start),
		tagName,
		attributes,
		nil,
		isSelfClosing,
	), nil
}

// parseAttribute 解析HTML属性
func (h *HtmlParser) parseAttribute() (string, data.GetValue, data.Control) {
	// 解析属性名
	attrName := h.parseAttributeName()
	if attrName == "" {
		return "", nil, nil
	}

	// 检查是否有等号
	if h.current().Type != token.ASSIGN {
		// 没有值的属性，如 disabled
		return attrName, data.NewBoolValue(true), nil
	}
	h.next()

	// 解析属性值
	var attrValue data.GetValue
	if h.current().Type == token.STRING {
		// 字符串值
		value := h.current().Literal
		h.next()
		attrValue = node.NewStringLiteral(h.NewTokenFrom(h.GetStart()), value)
	} else {
		// 其他类型的值，尝试解析为表达式或直接作为字符串
		attrValue = h.parseAttributeValue()
	}

	return attrName, attrValue, nil
}

// parseAttributeName 解析属性名
func (h *HtmlParser) parseAttributeName() string {
	// 直接使用当前token的Literal作为属性名
	if !h.isEOF() {
		name := h.current().Literal
		h.next()
		return name
	}

	return ""
}

// parseAttributeValue 解析属性值
func (h *HtmlParser) parseAttributeValue() data.GetValue {
	start := h.GetStart()

	// 如果是标识符，直接作为字符串处理
	if h.current().Type == token.IDENTIFIER {
		value := h.current().Literal
		h.next()
		return node.NewStringLiteral(h.NewTokenFrom(start), value)
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
		return node.NewStringLiteral(h.NewTokenFrom(start), value)
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
	var children []data.GetValue

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
			children = append(children, child)
		}

		// 防止无限循环：确保token位置有变化
		if h.isEOF() {
			break
		}
	}

	return children, nil
}

// parseHtmlChild 解析HTML子节点
func (h *HtmlParser) parseHtmlChild() (data.GetValue, data.Control) {
	if h.current().Type == token.LT {
		// 可能是HTML标签
		if h.checkPositionIs(1, token.IDENTIFIER) {
			// 直接在这里解析子标签，避免递归调用
			return h.parseSubHtmlTag()
		} else if h.checkPositionIs(1, token.QUO) {
			// 结束标签，停止解析
			return nil, nil
		}
	}

	// 解析文本内容
	return h.parseHtmlText()
}

// parseSubHtmlTag 解析子HTML标签（支持自闭合和非自闭合标签，避免递归）
func (h *HtmlParser) parseSubHtmlTag() (data.GetValue, data.Control) {
	start := h.GetStart()

	// 跳过开始的 < 符号
	if h.current().Type == token.LT {
		h.next()
	}

	// 解析标签名
	tagName := h.parseTagName()
	if tagName == "" {
		return nil, data.NewErrorThrow(h.newFrom(), errors.New("HTML标签缺少标签名"))
	}

	// 解析属性
	attributes := make(map[string]data.GetValue)
	for !h.isEOF() && h.current().Type != token.GT && h.current().Type != token.QUO {
		attrName, attrValue, acl := h.parseAttribute()
		if acl != nil {
			return nil, acl
		}
		if attrName != "" {
			attributes[attrName] = attrValue
		}
	}

	isSelfClosing := false
	if h.current().Type == token.QUO {
		h.next()
		if h.current().Type == token.GT {
			isSelfClosing = true
			h.next()
		} else {
			return nil, data.NewErrorThrow(h.newFrom(), errors.New("自闭合标签格式错误"))
		}
	} else if h.current().Type == token.GT {
		h.next()
	} else {
		return nil, data.NewErrorThrow(h.newFrom(), errors.New("HTML标签格式错误"))
	}

	if !isSelfClosing {
		children, acl := h.parseSubHtmlChildren()
		if acl != nil {
			return nil, acl
		}

		// 查找结束标签
		if !h.findClosingTag(tagName) {
			return nil, data.NewErrorThrow(h.newFrom(), errors.New("HTML标签缺少结束标签: "+tagName))
		}

		return node.NewHtmlNode(
			h.NewTokenFrom(start),
			tagName,
			attributes,
			children,
			isSelfClosing,
		), nil
	}

	return node.NewHtmlNode(
		h.NewTokenFrom(start),
		tagName,
		attributes,
		nil,
		isSelfClosing,
	), nil
}

// parseSubHtmlChildren 解析子HTML内容（避免递归）
func (h *HtmlParser) parseSubHtmlChildren() ([]data.GetValue, data.Control) {
	var children []data.GetValue

	for !h.isEOF() {
		// 检查是否是结束标签
		if h.current().Type == token.LT && h.checkPositionIs(1, token.QUO) {
			break
		}

		// 解析子节点
		child, acl := h.parseSubHtmlChild()
		if acl != nil {
			return nil, acl
		}
		if child != nil {
			children = append(children, child)
		}

		// 防止无限循环：确保token位置有变化
		if h.isEOF() {
			break
		}
	}

	return children, nil
}

// parseSubHtmlChild 解析子HTML节点（避免递归）
func (h *HtmlParser) parseSubHtmlChild() (data.GetValue, data.Control) {
	if h.current().Type == token.LT {
		// 可能是HTML标签
		if h.checkPositionIs(1, token.IDENTIFIER) {
			// 递归调用，但限制深度
			return h.parseSubHtmlTag()
		} else if h.checkPositionIs(1, token.QUO) {
			// 结束标签，停止解析
			return nil, nil
		}
	}

	// 解析文本内容
	return h.parseHtmlText()
}

// parseHtmlText 解析HTML文本内容
func (h *HtmlParser) parseHtmlText() (data.GetValue, data.Control) {
	start := h.GetStart()

	var text string
	initialPos := h.GetStart()
	lastEnd := -1

	for !h.isEOF() && h.current().Type != token.LT {
		curStart := h.GetStart()
		if lastEnd != -1 && curStart > lastEnd {
			// token 之间有间隔，补空格
			text += " "
		}
		text += h.current().Literal
		lastEnd = h.current().End
		h.next()

		// 防止无限循环：检查位置是否变化
		if h.GetStart() == initialPos {
			break
		}
		initialPos = h.GetStart()
	}

	if text == "" {
		return nil, nil
	}

	// 去除首尾空白
	text = trimSpace(text)
	if text == "" {
		return nil, nil
	}

	return node.NewStringLiteral(h.NewTokenFrom(start), text), nil
}

// findClosingTag 查找结束标签
func (h *HtmlParser) findClosingTag(tagName string) bool {
	// 检查是否是结束标签
	if h.current().Type == token.LT && h.checkPositionIs(1, token.QUO) {
		h.next() // 跳过 <
		h.next() // 跳过 /

		// 检查标签名是否匹配
		if h.current().Type == token.IDENTIFIER && h.current().Literal == tagName {
			h.next()

			// 检查结束的 >
			if h.current().Type == token.GT {
				h.next()
				return true
			}
		}
	}

	return false
}

// trimSpace 去除字符串首尾空白
func trimSpace(s string) string {
	// 简单的空白字符去除
	start := 0
	end := len(s)

	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}

func NewHtmlParser(parser *Parser) StatementParser {
	return &HtmlParser{
		Parser: parser,
	}
}
