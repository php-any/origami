package node

import (
	"github.com/php-any/origami/data"
)

// HtmlAttributeValue HTML属性值接口
type HtmlAttributeValue interface {
	// ProcessHtml 处理HTML节点，返回是否应该输出HTML
	ProcessHtml(ctx data.Context, htmlNode *HtmlNode) (shouldOutput bool, result string)
	// GetValue 获取原始值（用于兼容性）
	GetValue() data.GetValue
}

// AttrValueAdapter 适配器，将data.GetValue转换为HtmlAttributeValue
type AttrValueAdapter struct {
	*Node `pp:"-"`
	Value data.GetValue
}

// NewAttrValueAdapter 创建一个新的属性值适配器
func NewAttrValueAdapter(from data.From, value data.GetValue) *AttrValueAdapter {
	return &AttrValueAdapter{
		Node:  NewNode(from),
		Value: value,
	}
}

// ProcessHtml 处理普通属性值的HTML输出
func (a *AttrValueAdapter) ProcessHtml(ctx data.Context, htmlNode *HtmlNode) (bool, string) {
	// 普通属性值不控制HTML输出，总是返回true
	return true, ""
}

// GetValue 获取原始值
func (a *AttrValueAdapter) GetValue() data.GetValue {
	return a.Value
}

// AttrForValue 表示for循环属性值
type AttrForValue struct {
	*Node `pp:"-"`
	Val   data.Variable // 循环变量名（如 "item"）
	Key   data.Variable // 循环Key变量名（如 "index"，可选）
	Array data.GetValue // 要遍历的数组
}

// NewAttrForValue 创建一个新的for属性值
func NewAttrForValue(from data.From, array data.GetValue, key data.Variable, val data.Variable) *AttrForValue {
	return &AttrForValue{
		Node:  NewNode(from),
		Val:   val,
		Key:   key,
		Array: array,
	}
}

// ProcessHtml 处理for循环的HTML输出
func (a *AttrForValue) ProcessHtml(ctx data.Context, htmlNode *HtmlNode) (bool, string) {
	// 获取数组值
	arrayValue, ctl := a.Array.GetValue(ctx)
	if ctl != nil {
		return false, ""
	}

	// 检查是否是数组
	arrayData, ok := arrayValue.(*data.ArrayValue)
	if !ok {
		return false, ""
	}

	// 获取数组长度
	length := len(arrayData.Value)
	if length == 0 {
		return false, ""
	}

	// 遍历数组
	var result string
	for i := 0; i < length; i++ {
		// 获取当前元素
		item := arrayData.Value[i]

		// 设置循环变量
		if a.Val != nil {
			ctx.SetVariableValue(a.Val, item)
		}

		// 设置Key变量（如果有）
		if a.Key != nil {
			ctx.SetVariableValue(a.Key, data.NewIntValue(i))
		}

		// 生成当前迭代的HTML
		result += htmlNode.generateNormalHtml(ctx)
	}

	return true, result
}

// GetValue 获取原始值（for属性没有单一原始值，返回nil）
func (a *AttrForValue) GetValue() data.GetValue {
	return nil
}

// AttrIfValue 表示if条件链属性值（包含if、else-if、else）
type AttrIfValue struct {
	*Node     `pp:"-"`
	Condition data.GetValue // 条件表达式（if或else-if的条件，else为nil）
}

// NewAttrIfValue 创建一个新的if属性值
func NewAttrIfValue(from data.From, condition data.GetValue) *AttrIfValue {
	return &AttrIfValue{
		Node:      NewNode(from),
		Condition: condition,
	}
}

// ProcessHtml 处理if条件链的HTML输出
func (a *AttrIfValue) ProcessHtml(ctx data.Context, htmlNode *HtmlNode) (bool, string) {
	// 检查当前条件
	if a.Condition == nil {
		// else节点，总是执行
		return true, htmlNode.generateNormalHtml(ctx)
	}

	// 获取条件值
	conditionValue, ctl := a.Condition.GetValue(ctx)
	if ctl != nil {
		return false, ""
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

	// 如果条件为真，生成HTML
	if isTrue {
		return true, htmlNode.generateNormalHtml(ctx)
	}

	// 没有下一个节点，不输出
	return false, ""
}

// GetValue 获取原始值（if属性没有单一原始值，返回nil）
func (a *AttrIfValue) GetValue() data.GetValue {
	return nil
}
