package channel

import (
	"github.com/php-any/origami/data"
)

// ChannelClass 表示 Channel 类
type ChannelClass struct {
	channel *Channel
}

// NewChannelClass 创建一个新的 Channel 类实例
func NewChannelClass() data.ClassStmt {
	source := NewChannel()
	return &ChannelClass{
		channel: source,
	}
}

// GetValue 返回类实例
func (c *ChannelClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	source := NewChannel()

	return data.NewClassValue(&ChannelClass{
		channel: source,
	}, ctx.CreateBaseContext()), nil
}

// GetFrom 返回来源信息
func (c *ChannelClass) GetFrom() data.From {
	return nil
}

// GetName 返回类名
func (c *ChannelClass) GetName() string {
	return "Channel"
}

// GetExtend 返回父类
func (c *ChannelClass) GetExtend() *string {
	return nil
}

// GetImplements 返回实现的接口
func (c *ChannelClass) GetImplements() []string {
	return nil
}

// GetProperty 获取属性
func (c *ChannelClass) GetProperty(name string) (data.Property, bool) {
	return nil, false
}

// GetPropertyList 获取所有属性列表
func (c *ChannelClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

// GetMethod 获取方法
func (c *ChannelClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "send":
		return &ChannelSendMethod{source: c}, true
	case "receive":
		return &ChannelReceiveMethod{source: c}, true
	case "close":
		return &ChannelCloseMethod{source: c}, true
	case "isClosed":
		return &ChannelIsClosedMethod{source: c}, true
	case "len":
		return &ChannelLenMethod{source: c}, true
	case "cap":
		return &ChannelCapMethod{source: c}, true
	}
	return nil, false
}

// GetMethods 获取所有方法
func (c *ChannelClass) GetMethods() []data.Method {
	return []data.Method{
		&ChannelSendMethod{source: c},
		&ChannelReceiveMethod{source: c},
		&ChannelCloseMethod{source: c},
		&ChannelIsClosedMethod{source: c},
		&ChannelLenMethod{source: c},
		&ChannelCapMethod{source: c},
	}
}

// GetConstruct 获取构造函数
func (c *ChannelClass) GetConstruct() data.Method {
	return &ChannelConstructMethod{source: c}
}

// AddAnnotations 添加注解
func (c *ChannelClass) AddAnnotations(a *data.ClassValue) {
	// 暂时不处理注解
}
