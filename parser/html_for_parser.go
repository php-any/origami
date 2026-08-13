package parser

import (
	"errors"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// HtmlForAttributeParser 专门处理for属性的解析器
type HtmlForAttributeParser struct {
}

// Parser 实现 ParseAttribute 接口 - 处理for循环属性
func (h *HtmlForAttributeParser) Parser(name string, parser *Parser) (node.HtmlAttributeValue, data.Control) {
	tracker := parser.StartTracking()
	// 解析for属性值
	if !parser.checkPositionIs(0, token.ASSIGN) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("for属性缺少等号"))
	}
	parser.next()

	// 解析属性值
	if !parser.checkPositionIs(0, token.STRING) {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("for表达式需要使用双引号括起来"))
	}

	// 字符串值
	codes := parser.current().Literal()
	if len(codes) > 3 {
		if codes[0] == '"' && codes[len(codes)-1] == '"' {
			codes = strings.Trim(codes, "\"")
		} else if codes[0] == '\'' && codes[len(codes)-1] == '\'' {
			codes = strings.Trim(codes, "'")
		}
	}
	parser.next()
	vars, exprStr := h.parseForExpression(codes)

	if vars == nil {
		return nil, data.NewErrorThrow(tracker.EndBefore(), errors.New("for属性格式错误，应为：key, value in array 或 value in array"))
	}

	// 解析变量名
	keyVar := vars[0]
	keyVari := parser.scopeManager.LookupVariable(keyVar)
	if keyVari == nil {
		keyVari = parser.scopeManager.CurrentScope().AddVariable(keyVar, nil, tracker.EndBefore())
	}

	valueVar := vars[1]
	valueVari := parser.scopeManager.LookupVariable(valueVar)
	if valueVari == nil {
		valueVari = parser.scopeManager.CurrentScope().AddVariable(valueVar, nil, tracker.EndBefore())
	}

	// 使用主解释器解析表达式字符串
	arrayVari, acl := parser.ParseExpressionFromString(exprStr)
	if acl != nil {
		return nil, acl
	}

	// 创建AttrForValue
	return node.NewAttrForValue(tracker.EndBefore(), arrayVari, keyVari, valueVari), nil

}

// parseForExpression 解析for表达式，返回变量信息和表达式字符串
func (h *HtmlForAttributeParser) parseForExpression(forStr string) ([]string, string) {
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
func (h *HtmlForAttributeParser) parseVariables(varsStr string) []string {
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
