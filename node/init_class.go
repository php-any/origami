package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

// InitClass 表示 new 表达式
type InitClass struct {
	*Node     `pp:"-"`
	ClassName string
	KV        map[string]data.GetValue
}

func NewInitClass(from *TokenFrom, className string, KV map[string]data.GetValue) *InitClass {
	return &InitClass{
		Node:      NewNode(from),
		ClassName: className,
		KV:        KV,
	}
}

// GetValue 实现 Value 接口
func (n *InitClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	vm := ctx.GetVM()
	stmt, ok := vm.GetClass(n.ClassName)
	if !ok {
		return nil, data.NewErrorThrow(n.from, errors.New(fmt.Sprintf("类 %s 不存在", n.ClassName)))
	}

	object, acl := stmt.GetValue(ctx.CreateBaseContext())
	if acl != nil {
		return nil, acl
	}

	if object, ok := object.(*data.ClassValue); ok {
		for k, v := range n.KV {
			value, acl := v.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			object.SetProperty(k, value.(data.Value))
		}
	}

	return object, acl
}
