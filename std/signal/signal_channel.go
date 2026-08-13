package signal

import "os"

type SignalChannel struct {
	closed  bool
	channel chan os.Signal
}

func NewSignalChannel() *SignalChannel {
	return &SignalChannel{}
}

func (c *SignalChannel) Construct(capacity int) {
	if c.channel != nil && !c.closed {
		close(c.channel)
	}
	if capacity < 1 {
		capacity = 1
	}
	c.channel = make(chan os.Signal, capacity)
	c.closed = false
}

func (c *SignalChannel) Receive() (os.Signal, bool) {
	if c.channel == nil {
		return nil, false
	}
	sig, ok := <-c.channel
	return sig, ok
}

func (c *SignalChannel) Close() {
	if !c.closed && c.channel != nil {
		c.closed = true
		close(c.channel)
	}
}

func (c *SignalChannel) Chan() chan<- os.Signal {
	return c.channel
}
