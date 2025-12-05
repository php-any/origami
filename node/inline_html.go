package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// InlineHTMLNode 表示标签外的 HTML 内容
type InlineHTMLNode struct {
	*Node   `pp:"-"`
	Content string
}

// NewInlineHTMLNode 创建一个新的 InlineHTMLNode
func NewInlineHTMLNode(token *TokenFrom, content string) *InlineHTMLNode {
	return &InlineHTMLNode{
		Node:    NewNode(token),
		Content: content,
	}
}

// GetValue 直接输出内容
func (n *InlineHTMLNode) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	fmt.Print(n.Content)
	return nil, nil
}
