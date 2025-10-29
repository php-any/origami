package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// HtmlNode 表示HTML节点
type HtmlNode struct {
	*Node         `pp:"-"`
	TagName       string                        // 标签名
	Attributes    map[string]HtmlAttributeValue // 属性（使用新的属性值接口）
	Children      []data.GetValue               // 子节点
	IsSelfClosing bool                          // 是否是自闭合标签
}

// NewHtmlNode 创建一个新的HTML节点
func NewHtmlNode(from data.From, tagName string, attributes map[string]HtmlAttributeValue, children []data.GetValue, isSelfClosing bool) *HtmlNode {
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
	// 检查是否有特殊的属性值类型
	var forValue *AttrForValue
	var ifValue *AttrIfValue

	// 处理特殊属性（按类型判断，避免硬编码属性名）
	for _, value := range h.Attributes {
		if forAttr, ok := value.(*AttrForValue); ok {
			forValue = forAttr
			continue
		}
		if ifAttr, ok := value.(*AttrIfValue); ok {
			ifValue = ifAttr
			continue
		}
	}

	// 如果有if属性，执行条件链
	if ifValue != nil {
		shouldOutput, result := ifValue.ProcessHtml(ctx, h)
		if shouldOutput {
			return result
		}
		return ""
	}

	// 如果有for属性，执行for循环
	if forValue != nil {
		shouldOutput, result := forValue.ProcessHtml(ctx, h)
		if shouldOutput {
			return result
		}
		return ""
	}

	// 普通HTML节点处理
	return h.generateNormalHtml(ctx)
}

// generateNormalHtml 生成普通HTML
func (h *HtmlNode) generateNormalHtml(ctx data.Context) string {
	// 开始标签
	html := "<" + h.TagName

	// 添加普通属性（排除特殊属性，按类型跳过）
	for name, value := range h.Attributes {
		// 跳过 for / if 系列特殊属性
		if _, ok := value.(*AttrForValue); ok {
			continue
		}
		if _, ok := value.(*AttrIfValue); ok {
			continue
		}

		attrResult, ctl := value.GetValue().GetValue(ctx)
		if ctl != nil {
			continue
		}

		if strValue, ok := attrResult.(data.AsString); ok {
			html += fmt.Sprintf(` %s="%s"`, name, strValue.AsString())
		} else if boolValue, ok := attrResult.(data.AsBool); ok {
			if boolVal, err := boolValue.AsBool(); err == nil && boolVal {
				html += fmt.Sprintf(` %s`, name)
			}
		} else {
			html += fmt.Sprintf(` %s="%v"`, name, attrResult)
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
func (h *HtmlNode) GetAttributes() map[string]HtmlAttributeValue {
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

// HtmlDocTypeNode 表示 HTML 文档类型节点，作为整份文档的根容器
type HtmlDocTypeNode struct {
	*Node    `pp:"-"`
	DocType  string
	Children []data.GetValue
}

func NewHtmlDocTypeNode(from data.From, docType string, children []data.GetValue) *HtmlDocTypeNode {
	return &HtmlDocTypeNode{
		Node:     NewNode(from),
		DocType:  docType,
		Children: children,
	}
}

func (d *HtmlDocTypeNode) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	html := "<!DOCTYPE " + d.DocType + ">"
	for _, child := range d.Children {
		v, ctl := child.GetValue(ctx)
		if ctl != nil {
			continue
		}
		if s, ok := v.(data.AsString); ok {
			html += s.AsString()
		} else if n, ok := v.(*HtmlNode); ok {
			html += n.generateHtml(ctx)
		} else {
			html += fmt.Sprintf("%v", v)
		}
	}
	return data.NewStringValue(html), nil
}

// ScriptZyNode 表示 <script type="text/zy"> 脚本节点
type ScriptZyNode struct {
	*Node   `pp:"-"`
	Program *Program
}

func NewScriptZyNode(from data.From, program *Program) *ScriptZyNode {
	return &ScriptZyNode{Node: NewNode(from), Program: program}
}

func (s *ScriptZyNode) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if s.Program != nil {
		// 执行脚本程序（内部会通过 VM ThrowControl 处理控制流）
		s.Program.GetValue(ctx)
		return nil, nil
	}
	return data.NewStringValue(""), nil
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

// HtmlIfNodeType 定义HTML if节点的类型
type HtmlIfNodeType int

const (
	HtmlIfTypeIf     HtmlIfNodeType = iota // if节点
	HtmlIfTypeElseIf                       // else-if节点
	HtmlIfTypeElse                         // else节点
)

// HtmlIfNode 表示HTML if条件节点
type HtmlIfNode struct {
	*HtmlNode                // 嵌套的HTML节点
	Type      HtmlIfNodeType // 节点类型
	Condition data.GetValue  // 条件表达式
	NextNode  *HtmlIfNode    // 下一个条件节点（else-if 或 else）
}

// NewHtmlIfNode 创建一个新的HTML if条件节点
func NewHtmlIfNode(ifType HtmlIfNodeType, condition *AttrIfValue, htmlNode *HtmlNode) *HtmlIfNode {
	ret := &HtmlIfNode{
		Type:     ifType,
		HtmlNode: htmlNode,
		NextNode: nil,
	}
	if condition != nil {
		ret.Condition = condition.Condition
	}
	return ret
}

// SetNextNode 设置下一个条件节点
func (h *HtmlIfNode) SetNextNode(next *HtmlIfNode) {
	h.NextNode = next
}

// GetValue 获取HTML if条件节点的值
func (h *HtmlIfNode) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否是else节点（没有条件）
	if h.Type == HtmlIfTypeElse {
		// else节点，直接执行
		htmlValue, ctl := h.HtmlNode.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}

		if strValue, ok := htmlValue.(data.AsString); ok {
			return data.NewStringValue(strValue.AsString()), nil
		} else {
			return data.NewStringValue(fmt.Sprintf("%v", htmlValue)), nil
		}
	}

	// 检查条件是否存在
	if h.Condition == nil {
		// 没有条件，返回空字符串
		return data.NewStringValue(""), nil
	}

	// 获取条件值
	conditionValue, ctl := h.Condition.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	// 检查条件是否为真
	var isTrue bool
	if boolValue, ok := conditionValue.(data.AsBool); ok {
		if boolVal, err := boolValue.AsBool(); err == nil {
			isTrue = boolVal
		}
	} else if _, ok := conditionValue.(*data.NullValue); ok {
		isTrue = false
	} else {
		// 非空值视为真
		isTrue = true
	}

	// 如果条件为真，执行当前节点
	if isTrue {
		htmlValue, ctl := h.HtmlNode.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}

		if strValue, ok := htmlValue.(data.AsString); ok {
			return data.NewStringValue(strValue.AsString()), nil
		} else {
			return data.NewStringValue(fmt.Sprintf("%v", htmlValue)), nil
		}
	}

	// 条件为假，检查是否有下一个节点（else-if 或 else）
	if h.NextNode != nil {
		return h.NextNode.GetValue(ctx)
	}

	// 没有下一个节点，返回空字符串
	return data.NewStringValue(""), nil
}
