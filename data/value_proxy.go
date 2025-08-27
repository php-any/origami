package data

func NewProxyValue(class ClassStmt, ctx Context) *ProxyValue {
	return &ProxyValue{
		ObjectValue: NewObjectValue(),
		Class:       class,
		Context:     ctx,
	}
}

// ProxyValue 代理类的值
type ProxyValue ClassValue

func (o *ProxyValue) SetProperty(name string, value Value) {
	if set, ok := o.Class.(SetProperty); ok {
		set.SetProperty(name, value)
	} else {
		o.property.Store(name, value)
	}
}

func (o *ProxyValue) GetSource() any {
	if p, ok := o.Class.(interface{ GetSource() any }); ok {
		return p.GetSource()
	}
	return nil
}
