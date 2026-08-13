package reflection

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ReflectionObjectClass 实现 PHP ReflectionObject（继承 ReflectionClass，仅接受对象）。
type ReflectionObjectClass struct {
	node.Node
}

func (c *ReflectionObjectClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

func (c *ReflectionObjectClass) GetName() string { return "ReflectionObject" }

func (c *ReflectionObjectClass) GetExtend() *string {
	parent := "ReflectionClass"
	return &parent
}

func (c *ReflectionObjectClass) GetImplements() []string { return nil }

func (c *ReflectionObjectClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

func (c *ReflectionObjectClass) GetPropertyList() []data.Property { return nil }

func (c *ReflectionObjectClass) GetMethod(name string) (data.Method, bool) {
	// 复用 ReflectionClass 方法表；构造函数仍走 ReflectionClassConstructMethod（支持 object）
	base := &ReflectionClassClass{}
	return base.GetMethod(name)
}

func (c *ReflectionObjectClass) GetMethods() []data.Method {
	return (&ReflectionClassClass{}).GetMethods()
}

func (c *ReflectionObjectClass) GetConstruct() data.Method {
	return &ReflectionClassConstructMethod{}
}
