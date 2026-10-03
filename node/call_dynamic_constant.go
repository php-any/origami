package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

type CallDynamicConstant struct {
	*Node
	Class, Name data.GetValue
}

func NewCallDynamicConstant(from *TokenFrom, class, name data.GetValue) *CallDynamicConstant {
	return &CallDynamicConstant{Node: NewNode(from), Class: class, Name: name}
}

func (c *CallDynamicConstant) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	value, ctl := c.Name.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}
	name, ok := value.(*data.StringValue)
	if !ok {
		return nil, data.NewErrorThrowByName(c.GetFrom(), fmt.Errorf("Class name constant must be a string"), "TypeError")
	}
	class := c.Class
	if _, ok := class.(data.GetStaticProperty); !ok {
		target, ctl := class.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if className, ok := target.(*data.StringValue); ok {
			loaded, ctl := ctx.GetVM().GetOrLoadClass(className.Value)
			if ctl != nil {
				return nil, ctl
			}
			if loaded == nil {
				return nil, data.NewErrorThrowByName(c.GetFrom(), fmt.Errorf("Class %q not found", className.Value), "Error")
			}
			class = loaded
		} else {
			class = target
		}
	}
	return NewCallStaticProperty(c.GetFrom().(*TokenFrom), class, name.Value).GetValue(ctx)
}
