package node

import (
	"sync"

	"github.com/php-any/origami/data"
)

var groups sync.Map

// GlobalsNode 返回空数组占位（$GLOBALS、$_ENV）避免解析报错
type GlobalsNode struct {
	*Node `pp:"-"`
	Name  string
}

func NewGlobalsNode(from data.From, name string) *GlobalsNode {
	return &GlobalsNode{Node: NewNode(from), Name: name}
}

func (g *GlobalsNode) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	v, ok := groups.Load(g.Name)
	if !ok {
		v = data.NewObjectValue()
		groups.Store(g.Name, v)
	}

	return v.(data.GetValue), nil
}
