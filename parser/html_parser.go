package parser

import (
	"errors"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
	"strings"
)

type HtmlParser struct {
	*Parser
}

func (h *HtmlParser) Parse() (data.GetValue, data.Control) {
	return h.parseHtmlContent()
}

// parseHtmlContent 解析HTML内容
func (h *HtmlParser) parseHtmlContent() (data.GetValue, data.Control) {
	tracker := h.StartTracking()

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
	if attr, ok := attributes["for"]; ok {
		// 解析for属性，创建HtmlForNode
		return h.createHtmlForNode(from, tagName, attributes, children, isSelfClosing, attr)
	}

	return node.NewHtmlNode(
		from,
		tagName,
		attributes,
		children,
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
		attrValue = node.NewStringLiteral(h.FromCurrentToken(), value)
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
	// 如果是标识符，直接作为字符串处理
	if h.current().Type == token.IDENTIFIER {
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
			return h.parseHtmlContent()
		} else if h.checkPositionIs(1, token.QUO) {
			// 结束标签，停止解析
			return nil, nil
		}
	}

	// 解析文本内容
	return h.parseHtmlText()
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
			return h.parseHtmlContent()
		} else if h.checkPositionIs(1, token.QUO) {
			// 结束标签，停止解析
			return nil, nil
		}
	}

	// 解析文本内容
	return h.parseHtmlText()
}

// parseText 解析文本内容
func (h *HtmlParser) parseText() data.GetValue {
	var text string

	// 收集所有文本内容，直到遇到 < 或 EOF
	for !h.isEOF() && h.current().Type != token.LT {
		text += h.current().Literal
		h.next()
	}

	if text == "" {
		return nil
	}

	return node.NewStringLiteral(h.FromCurrentToken(), text)
}

// parseHtmlText 解析HTML文本内容，支持插值字符串
func (h *HtmlParser) parseHtmlText() (data.GetValue, data.Control) {
	var textParts []data.GetValue
	initialPos := h.GetStart()
	lastEnd := -1
	currentText := ""

	for !h.isEOF() && h.current().Type != token.LT {
		curStart := h.GetStart()
		if lastEnd != -1 && curStart > lastEnd {
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

			// 跳过 { 和 $var
			h.nextAndCheck(token.LBRACE) // 跳过 {

			// 解析插值表达式
			// 收集表达式字符串直到遇到 }
			exprStr := ""
			for !h.isEOF() && h.current().Type != token.RBRACE {
				exprStr += h.current().Literal
				h.next()
			}

			// 使用ParseExpressionFromString解析表达式
			expr, ctl := h.Parser.ParseExpressionFromString(exprStr)
			if ctl != nil {
				return nil, ctl
			}

			// 检查是否以 } 结束
			if h.current().Type != token.RBRACE {
				return nil, data.NewErrorThrow(h.newFrom(), errors.New("插值表达式缺少结束的 }"))
			}
			h.nextAndCheck(token.RBRACE)

			// 添加表达式到结果中
			textParts = append(textParts, expr)
		} else {
			// 普通文本
			currentText += h.current().Literal
			h.next()
		}

		lastEnd = h.current().End

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

// createHtmlForNode 创建HTML for循环节点
func (h *HtmlParser) createHtmlForNode(from *node.TokenFrom, tagName string, attributes map[string]data.GetValue, children []data.GetValue, isSelfClosing bool, forAttr data.GetValue) (data.GetValue, data.Control) {
	tracker := h.StartTracking()
	// 解析for属性，格式应该是 "key, value in array" 或 "value in array"
	// 在解析阶段，forAttr应该是一个字符串字面量
	var forStr string
	if strLiteral, ok := forAttr.(*node.StringLiteral); ok {
		forStr = strLiteral.Value
	} else {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("for属性必须是字符串字面量"))
	}

	// 解析for字符串，格式：key, value in array 或 value in array
	vars, exprStr := h.parseForExpression(forStr)
	if vars == nil {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("for属性格式错误，应为：key, value in array 或 value in array"))
	}
	// 解析变量名
	keyVar := vars[0]
	val := h.scopeManager.CurrentScope().AddVariable(keyVar, nil, from)
	keyVari := node.NewVariableWithFirst(from, val)
	valueVar := vars[1]
	val = h.scopeManager.CurrentScope().AddVariable(valueVar, nil, from)
	valueVari := node.NewVariableWithFirst(from, val)

	// 使用主解释器解析表达式字符串
	arrayVari, acl := h.Parser.ParseExpressionFromString(exprStr)
	if acl != nil {
		return nil, acl
	}

	// 创建嵌套的HTML节点（不包含for属性）
	nestedAttributes := make(map[string]data.GetValue)
	for k, v := range attributes {
		if k != "for" {
			nestedAttributes[k] = v
		}
	}

	nestedHtmlNode := node.NewHtmlNode(
		from,
		tagName,
		nestedAttributes,
		children,
		isSelfClosing,
	)

	// 创建HtmlForNode
	return node.NewHtmlForNode(
		from,
		arrayVari,
		keyVari,
		valueVari,
		nestedHtmlNode,
	), nil
}

// parseForExpression 解析for表达式，返回变量信息和表达式字符串
func (h *HtmlParser) parseForExpression(forStr string) ([]string, string) {
	// 查找 " in " 分隔符
	inIndex := -1
	for i := 0; i < len(forStr)-3; i++ {
		if forStr[i:i+4] == " in " {
			inIndex = i
			break
		}
	}

	if inIndex == -1 {
		return nil, ""
	}

	// 分割变量部分和表达式部分
	varsPart := strings.TrimSpace(forStr[:inIndex])
	exprPart := strings.TrimSpace(forStr[inIndex+4:])

	// 解析变量部分
	vars := h.parseVariables(varsPart)
	if len(vars) == 1 {
		// 只有一个变量，作为value，key设为"_"
		return []string{"_", vars[0]}, exprPart
	} else if len(vars) == 2 {
		// 有两个变量，第一个是key，第二个是value
		return []string{vars[0], vars[1]}, exprPart
	} else {
		return nil, ""
	}
}

// parseVariables 解析变量列表
func (h *HtmlParser) parseVariables(varsStr string) []string {
	// 简单的逗号分割
	vars := make([]string, 0)
	current := ""

	for _, char := range varsStr {
		if char == ',' {
			if current != "" {
				vars = append(vars, strings.TrimSpace(current))
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		vars = append(vars, strings.TrimSpace(current))
	}

	return vars
}
