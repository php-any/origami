package channel

import (
	"github.com/php-any/origami/data"
)

// Channel 表示一个 channel 实例，封装 Go 的 chan
type Channel struct {
	closed  bool
	channel chan data.Value
}

// NewChannel 创建一个新的 channel
func NewChannel() *Channel {
	return &Channel{}
}

// Construct 构造函数，支持传入容量参数
func (c *Channel) Construct(ctx data.Context, capacity data.Value) {
	// 如果已经初始化，先关闭旧的 channel
	if c.channel != nil && !c.closed {
		close(c.channel)
	}

	// 获取容量参数
	capValue := 0 // 默认无缓冲
	if capacity != nil {
		if intVal, ok := capacity.(*data.IntValue); ok {
			capValue, _ = intVal.AsInt()
		}
	}

	// 确保容量为非负数
	if capValue < 0 {
		capValue = 0
	}

	// 创建新的 channel
	c.channel = make(chan data.Value, capValue)
	c.closed = false
}

// Send 发送数据到 channel（对齐 Go 的用法）
func (c *Channel) Send(value data.Value) bool {
	if c.closed || c.channel == nil {
		return false
	}

	// 发送数据（Go 风格的发送）
	c.channel <- value
	return true
}

// Receive 从 channel 接收数据（对齐 Go 的用法）
func (c *Channel) Receive() (data.Value, bool) {
	if c.channel == nil {
		return nil, false
	}

	// 接收数据（Go 风格的接收）
	value, ok := <-c.channel
	return value, ok
}

// Close 关闭 channel
func (c *Channel) Close() {
	if !c.closed && c.channel != nil {
		c.closed = true
		close(c.channel)
	}
}

// IsClosed 检查 channel 是否已关闭
func (c *Channel) IsClosed() bool {
	return c.closed
}

// Len 返回 channel 缓冲区长度
func (c *Channel) Len() int {
	if c.channel == nil {
		return 0
	}
	return len(c.channel)
}

// Cap 返回 channel 容量
func (c *Channel) Cap() int {
	if c.channel == nil {
		return 0
	}
	return cap(c.channel)
}
