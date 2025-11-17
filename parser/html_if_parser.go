package parser

import (
	"errors"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// HtmlIfAttributeParser 专门处理if属性的解析器
type HtmlIfAttributeParser struct {
	parser *Parser
}

// Parser 实现 ParseAttribute 接口 - 处理if条件属性
func (h *HtmlIfAttributeParser) Parser(name string, parser *Parser) (node.HtmlAttributeValue, data.Control) {
	tracker := h.parser.StartTracking()
	// 解析if属性值
	if !parser.checkPositionIs(0, token.ASSIGN) {
		return nil, data.NewErrorThrow(parser.newFrom(), errors.New("if属性缺少等号"))
	}
	parser.next()

	// 解析属性值
	codes := parser.current().Literal()
	if len(codes) > 3 {
		if codes[0] == '"' && codes[len(codes)-1] == '"' {
			codes = strings.Trim(codes, "\"")
		} else if codes[0] == '\'' && codes[len(codes)-1] == '\'' {
			codes = strings.Trim(codes, "'")
		}
	}
	parser.next()

	condition, acl := parser.ParseExpressionFromString(codes)
	if acl != nil {
		return nil, acl
	}

	return node.NewAttrIfValue(tracker.EndBefore(), condition), nil
}

// HtmlElseIfAttributeParser 专门处理else-if属性的解析器
type HtmlElseIfAttributeParser struct {
	parser *Parser
}

// Parser 实现 ParseAttribute 接口 - 处理else-if条件属性
func (h *HtmlElseIfAttributeParser) Parser(name string, parser *Parser) (node.HtmlAttributeValue, data.Control) {
	tracker := h.parser.StartTracking()
	// 解析else-if属性值
	if !parser.checkPositionIs(0, token.ASSIGN) {
		return nil, data.NewErrorThrow(parser.newFrom(), errors.New("else-if属性缺少等号"))
	}
	parser.next()

	// 解析属性值
	var attrValue data.GetValue
	if parser.checkPositionIs(0, token.STRING) {
		// 字符串值
		value := parser.current().Literal()
		parser.next()
		attrValue = node.NewStringLiteral(tracker.EndBefore(), value)
	} else {
		// 其他类型的值，尝试解析为表达式
		exprParser := NewExpressionParser(parser)
		var acl data.Control
		attrValue, acl = exprParser.Parse()
		if acl != nil {
			return nil, acl
		}
	}

	// 检查else-if属性是否是字符串字面量
	if strLiteral, ok := attrValue.(*node.StringLiteral); ok {
		// 如果是字符串，尝试解析为表达式
		condition, acl := parser.ParseExpressionFromString(strLiteral.Value)
		if acl != nil {
			return nil, acl
		}
		return node.NewAttrIfValue(tracker.EndBefore(), condition), nil
	} else {
		// 如果不是字符串，直接使用
		return node.NewAttrIfValue(tracker.EndBefore(), attrValue), nil
	}
}

// HtmlElseAttributeParser 专门处理else属性的解析器
type HtmlElseAttributeParser struct {
	parser *Parser
}

// Parser 实现 ParseAttribute 接口 - 处理else条件属性
func (h *HtmlElseAttributeParser) Parser(name string, parser *Parser) (node.HtmlAttributeValue, data.Control) {
	tracker := h.parser.StartTracking()
	// else属性不需要值，直接返回
	return node.NewAttrIfValue(tracker.EndBefore(), nil), nil
}
