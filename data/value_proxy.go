package data

func NewProxyValue(class ClassStmt, ctx Context) *ProxyValue {
	var vm VM
	if ctx != nil {
		vm = ctx.GetVM()
	}
	return &ProxyValue{
		ObjectValue: NewObjectValue(),
		Class:       class,
		Context:     ctx,
		vm:          vm,
	}
}

// ProxyValue 代理类的值
type ProxyValue = ClassValue

func (c *ProxyValue) GetSource() any {
	if p, ok := c.Class.(GetSource); ok {
		return p.GetSource()
	}
	return nil
}
