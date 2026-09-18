package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PHP：ReflectionFunction / ReflectionMethod 均继承 ReflectionFunctionAbstract。
func reflectionFunctionAbstractParent() *string {
	name := "ReflectionFunctionAbstract"
	return &name
}

// ReflectionFunctionAbstractClass 对齐 PHP 抽象基类，供 instanceof / 参数类型检查。
type ReflectionFunctionAbstractClass struct {
	node.Node
}

func (c *ReflectionFunctionAbstractClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *ReflectionFunctionAbstractClass) GetName() string { return "ReflectionFunctionAbstract" }

func (c *ReflectionFunctionAbstractClass) GetExtend() *string { return nil }

func (c *ReflectionFunctionAbstractClass) GetImplements() []string { return nil }

func (c *ReflectionFunctionAbstractClass) GetProperty(string) (data.Property, bool) {
	return nil, false
}

func (c *ReflectionFunctionAbstractClass) GetPropertyList() []data.Property { return nil }

func (c *ReflectionFunctionAbstractClass) GetMethod(string) (data.Method, bool) {
	return nil, false
}

func (c *ReflectionFunctionAbstractClass) GetMethods() []data.Method { return nil }

func (c *ReflectionFunctionAbstractClass) GetConstruct() data.Method { return nil }
