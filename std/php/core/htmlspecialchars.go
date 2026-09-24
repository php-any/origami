package core

import (
	"fmt"
	"html"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// HtmlspecialcharsFunction 实现 htmlspecialchars 函数
// PHP 8：参数须为 string（或可转 string 的标量）；对象仅在有 __toString 时可用。
type HtmlspecialcharsFunction struct{}

func NewHtmlspecialcharsFunction() data.FuncStmt {
	return &HtmlspecialcharsFunction{}
}

func (f *HtmlspecialcharsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	if strVal == nil {
		return data.NewStringValue(""), nil
	}

	if _, isNull := strVal.(*data.NullValue); isNull {
		return data.NewStringValue(""), nil
	}

	if cv, ok := strVal.(*data.ClassValue); ok {
		str, acl := node.ValueToDisplayString(ctx, cv)
		if acl != nil {
			return nil, acl
		}
		// 无 __toString 时 ValueToDisplayString 仍会返回 Object(...)；对齐 PHP 抛 TypeError
		if _, has := cv.GetMethod("__toString"); !has {
			name := "object"
			if cv.Class != nil {
				name = cv.Class.GetName()
			}
			return nil, data.NewErrorThrow(nil, fmt.Errorf("htmlspecialchars(): Argument #1 ($string) must be of type string, %s given", name))
		}
		return data.NewStringValue(html.EscapeString(str)), nil
	}

	str := strVal.AsString()
	result := html.EscapeString(str)

	return data.NewStringValue(result), nil
}

func (f *HtmlspecialcharsFunction) GetName() string {
	return "htmlspecialchars"
}

var htmlspecialcharsFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "string", 0, nil, nil),
	node.NewParameter(nil, "flags", 1, nil, nil),
	node.NewParameter(nil, "encoding", 2, nil, nil),
	node.NewParameter(nil, "double_encode", 3, nil, nil),
}

func (f *HtmlspecialcharsFunction) GetParams() []data.GetValue {
	return htmlspecialcharsFunctionGetParams
}

var htmlspecialcharsFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "string", 0, nil),
	node.NewVariable(nil, "flags", 1, nil),
	node.NewVariable(nil, "encoding", 2, nil),
	node.NewVariable(nil, "double_encode", 3, nil),
}

func (f *HtmlspecialcharsFunction) GetVariables() []data.Variable {
	return htmlspecialcharsFunctionGetVariables
}
