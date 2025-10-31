package channel

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ChannelConstructMethod 构造函数
type ChannelConstructMethod struct {
	source *ChannelClass
}

func (c *ChannelConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取容量参数（可选）
	var capacity data.Value
	if value, ok := ctx.GetIndexValue(0); ok {
		capacity = value
	}

	// 调用构造函数
	c.source.channel.Construct(ctx, capacity)

	return data.NewNullValue(), nil
}

func (c *ChannelConstructMethod) GetName() string {
	return "__construct"
}

func (c *ChannelConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (c *ChannelConstructMethod) GetIsStatic() bool {
	return false
}

func (c *ChannelConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "capacity", 0, nil, data.NewBaseType("int")),
	}
}

func (c *ChannelConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "capacity", 0, nil),
	}
}

func (c *ChannelConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

// ChannelSendMethod Send 方法（对齐 Go 的用法）
type ChannelSendMethod struct {
	source *ChannelClass
}

func (c *ChannelSendMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if c.source.channel == nil {
		return nil, utils.NewThrow(errors.New("channel not initialized"))
	}

	// 获取参数
	value, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("missing parameter: value"))
	}

	// 发送数据（Go 风格的发送）
	success := c.source.channel.Send(value)
	return data.NewBoolValue(success), nil
}

func (c *ChannelSendMethod) GetName() string {
	return "send"
}

func (c *ChannelSendMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (c *ChannelSendMethod) GetIsStatic() bool {
	return false
}

func (c *ChannelSendMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (c *ChannelSendMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, nil),
	}
}

func (c *ChannelSendMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

// ChannelReceiveMethod Receive 方法（对齐 Go 的用法）
type ChannelReceiveMethod struct {
	source *ChannelClass
}

func (c *ChannelReceiveMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if c.source.channel == nil {
		return nil, utils.NewThrow(errors.New("channel not initialized"))
	}

	// 接收数据（Go 风格的接收）
	value, ok := c.source.channel.Receive()
	if !ok {
		return data.NewNullValue(), nil
	}

	return value, nil
}

func (c *ChannelReceiveMethod) GetName() string {
	return "receive"
}

func (c *ChannelReceiveMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (c *ChannelReceiveMethod) GetIsStatic() bool {
	return false
}

func (c *ChannelReceiveMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (c *ChannelReceiveMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (c *ChannelReceiveMethod) GetReturnType() data.Types {
	return nil
}

// ChannelCloseMethod Close 方法
type ChannelCloseMethod struct {
	source *ChannelClass
}

func (c *ChannelCloseMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if c.source.channel == nil {
		return nil, utils.NewThrow(errors.New("channel not initialized"))
	}

	c.source.channel.Close()
	return nil, nil
}

func (c *ChannelCloseMethod) GetName() string {
	return "close"
}

func (c *ChannelCloseMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (c *ChannelCloseMethod) GetIsStatic() bool {
	return false
}

func (c *ChannelCloseMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (c *ChannelCloseMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (c *ChannelCloseMethod) GetReturnType() data.Types {
	return data.NewBaseType("void")
}

// ChannelIsClosedMethod IsClosed 方法
type ChannelIsClosedMethod struct {
	source *ChannelClass
}

func (c *ChannelIsClosedMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if c.source.channel == nil {
		return nil, utils.NewThrow(errors.New("channel not initialized"))
	}

	closed := c.source.channel.IsClosed()
	return data.NewBoolValue(closed), nil
}

func (c *ChannelIsClosedMethod) GetName() string {
	return "isClosed"
}

func (c *ChannelIsClosedMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (c *ChannelIsClosedMethod) GetIsStatic() bool {
	return false
}

func (c *ChannelIsClosedMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (c *ChannelIsClosedMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (c *ChannelIsClosedMethod) GetReturnType() data.Types {
	return data.NewBaseType("bool")
}

// ChannelLenMethod Len 方法
type ChannelLenMethod struct {
	source *ChannelClass
}

func (c *ChannelLenMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if c.source.channel == nil {
		return nil, utils.NewThrow(errors.New("channel not initialized"))
	}

	length := c.source.channel.Len()
	return data.NewIntValue(length), nil
}

func (c *ChannelLenMethod) GetName() string {
	return "len"
}

func (c *ChannelLenMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (c *ChannelLenMethod) GetIsStatic() bool {
	return false
}

func (c *ChannelLenMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (c *ChannelLenMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (c *ChannelLenMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}

// ChannelCapMethod Cap 方法
type ChannelCapMethod struct {
	source *ChannelClass
}

func (c *ChannelCapMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	if c.source.channel == nil {
		return nil, utils.NewThrow(errors.New("channel not initialized"))
	}

	capacity := c.source.channel.Cap()
	return data.NewIntValue(capacity), nil
}

func (c *ChannelCapMethod) GetName() string {
	return "cap"
}

func (c *ChannelCapMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (c *ChannelCapMethod) GetIsStatic() bool {
	return false
}

func (c *ChannelCapMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (c *ChannelCapMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (c *ChannelCapMethod) GetReturnType() data.Types {
	return data.NewBaseType("int")
}
