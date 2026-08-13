package data

func NewProxyValue(class ClassStmt, ctx Context) *ProxyValue {
	return &ProxyValue{
		ObjectValue: NewObjectValue(),
		Class:       class,
		Context:     ctx,
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
