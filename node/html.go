package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

// HtmlNode 表示HTML节点
type HtmlNode struct {
	*Node         `pp:"-"`
	TagName       string                   // 标签名
	Attributes    map[string]data.GetValue // 属性
	Children      []data.GetValue          // 子节点
	IsSelfClosing bool                     // 是否是自闭合标签
}

// NewHtmlNode 创建一个新的HTML节点
func NewHtmlNode(from data.From, tagName string, attributes map[string]data.GetValue, children []data.GetValue, isSelfClosing bool) *HtmlNode {
	return &HtmlNode{
		Node:          NewNode(from),
		TagName:       tagName,
		Attributes:    attributes,
		Children:      children,
		IsSelfClosing: isSelfClosing,
	}
}

// GetValue 获取HTML节点的值
func (h *HtmlNode) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 生成HTML字符串
	html := h.generateHtml(ctx)
	return data.NewStringValue(html), nil
}

// generateHtml 生成HTML字符串
func (h *HtmlNode) generateHtml(ctx data.Context) string {
	// 开始标签
	html := "<" + h.TagName

	// 添加属性
	for name, value := range h.Attributes {
		attrValue, ctl := value.GetValue(ctx)
		if ctl != nil {
			continue
		}

		if strValue, ok := attrValue.(data.AsString); ok {
			html += fmt.Sprintf(` %s="%s"`, name, strValue.AsString())
		} else if boolValue, ok := attrValue.(data.AsBool); ok {
			if boolVal, err := boolValue.AsBool(); err == nil && boolVal {
				html += fmt.Sprintf(` %s`, name)
			}
		} else {
			html += fmt.Sprintf(` %s="%v"`, name, attrValue)
		}
	}

	if h.IsSelfClosing {
		html += " />"
	} else {
		html += ">"

		// 添加子节点
		for _, child := range h.Children {
			childValue, ctl := child.GetValue(ctx)
			if ctl != nil {
				continue
			}

			if strValue, ok := childValue.(data.AsString); ok {
				html += strValue.AsString()
			} else if htmlNode, ok := childValue.(*HtmlNode); ok {
				childHtml := htmlNode.generateHtml(ctx)
				html += childHtml
			} else {
				html += fmt.Sprintf("%v", childValue)
			}
		}

		// 结束标签
		html += "</" + h.TagName + ">"
	}

	return html
}

// GetTagName 返回标签名
func (h *HtmlNode) GetTagName() string {
	return h.TagName
}

// GetAttributes 返回属性
func (h *HtmlNode) GetAttributes() map[string]data.GetValue {
	return h.Attributes
}

// GetChildren 返回子节点
func (h *HtmlNode) GetChildren() []data.GetValue {
	return h.Children
}

// IsSelfClosing 返回是否是自闭合标签
func (h *HtmlNode) IsSelfClosingTag() bool {
	return h.IsSelfClosing
}

// HtmlForNode 表示HTML for循环节点
type HtmlForNode struct {
	*Node    `pp:"-"`
	Array    data.GetValue // 要遍历的数组
	Key      data.Variable // 键变量名（可选）
	Value    data.Variable // 值变量名
	HtmlNode *HtmlNode     // 嵌套的HTML节点
}

// NewHtmlForNode 创建一个新的HTML for循环节点
func NewHtmlForNode(from data.From, array data.GetValue, key data.Variable, value data.Variable, htmlNode *HtmlNode) *HtmlForNode {
	return &HtmlForNode{
		Node:     NewNode(from),
		Array:    array,
		Key:      key,
		Value:    value,
		HtmlNode: htmlNode,
	}
}

// GetValue 获取HTML for循环节点的值
func (h *HtmlForNode) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 获取数组值
	arrayValue, ctl := h.Array.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	var resultHtml string

	// 检查数组值是否为数组类型
	switch array := arrayValue.(type) {
	case *data.ArrayValue:
		// 遍历数组
		for i, element := range array.Value {
			// 设置值变量
			ctx.SetVariableValue(h.Value, element)

			// 如果有键变量，设置键变量
			if h.Key != nil {
				keyValue := data.NewIntValue(i)
				ctx.SetVariableValue(h.Key, keyValue)
			}

			// 执行嵌套的HTML节点
			htmlValue, ctl := h.HtmlNode.GetValue(ctx)
			if ctl != nil {
				return nil, ctl
			}

			if strValue, ok := htmlValue.(data.AsString); ok {
				resultHtml += strValue.AsString()
			} else {
				resultHtml += fmt.Sprintf("%v", htmlValue)
			}
		}
	case *data.ObjectValue:
		// 遍历对象
		for key, element := range array.GetProperties() {
			// 设置值变量
			ctx.SetVariableValue(h.Value, element)

			// 如果有键变量，设置键变量
			if h.Key != nil {
				keyValue := data.NewStringValue(key)
				ctx.SetVariableValue(h.Key, keyValue)
			}

			// 执行嵌套的HTML节点
			htmlValue, ctl := h.HtmlNode.GetValue(ctx)
			if ctl != nil {
				return nil, ctl
			}

			if strValue, ok := htmlValue.(data.AsString); ok {
				resultHtml += strValue.AsString()
			} else {
				resultHtml += fmt.Sprintf("%v", htmlValue)
			}
		}
	case *data.NullValue:
		// 空数组，返回空字符串
		return data.NewStringValue(""), nil
	default:
		return nil, data.NewErrorThrow(h.from, fmt.Errorf("for HTML 只能遍历数组或对象"))
	}

	return data.NewStringValue(resultHtml), nil
}
