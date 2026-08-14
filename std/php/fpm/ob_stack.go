package fpm

import (
	"strings"

	"github.com/php-any/origami/data"
)

type obStack struct {
	buffers []*strings.Builder
}

func (s *obStack) level() int {
	return len(s.buffers)
}

func (s *obStack) push() {
	s.buffers = append(s.buffers, &strings.Builder{})
}

func (s *obStack) pop() string {
	if len(s.buffers) == 0 {
		return ""
	}
	last := s.buffers[len(s.buffers)-1].String()
	s.buffers = s.buffers[:len(s.buffers)-1]
	return last
}

func (s *obStack) contents() string {
	if len(s.buffers) == 0 {
		return ""
	}
	return s.buffers[len(s.buffers)-1].String()
}

func (s *obStack) write(str string) bool {
	if len(s.buffers) == 0 {
		return false
	}
	s.buffers[len(s.buffers)-1].WriteString(str)
	return true
}

// flush 弹出并返回栈顶缓冲内容，调用方负责把内容写到上一层或最终输出。
func (s *obStack) flush() string {
	return s.pop()
}

// cleanCurrent 清空栈顶缓冲内容但不结束缓冲，返回是否成功。
func (s *obStack) cleanCurrent() bool {
	if len(s.buffers) == 0 {
		return false
	}
	s.buffers[len(s.buffers)-1] = &strings.Builder{}
	return true
}

// length 返回栈顶缓冲的字节长度（无缓冲返回 0）。
func (s *obStack) length() int {
	if len(s.buffers) == 0 {
		return 0
	}
	return s.buffers[len(s.buffers)-1].Len()
}

// status 返回缓冲层状态；full=true 返回全部层，否则仅最顶层。
func (s *obStack) status(full bool) []data.OutputBufferStatusInfo {
	n := len(s.buffers)
	if n == 0 {
		return nil
	}
	start := 0
	if !full {
		start = n - 1
	}
	result := make([]data.OutputBufferStatusInfo, 0, n-start)
	for i := start; i < n; i++ {
		result = append(result, data.OutputBufferStatusInfo{
			Level:      i + 1,
			Type:       1,
			Flags:      0,
			ChunkSize:  0,
			BufferSize: s.buffers[i].Len(),
			Name:       "default output handler",
		})
	}
	return result
}

// handlers 返回所有激活缓冲的处理器名。
func (s *obStack) handlers() []string {
	n := len(s.buffers)
	if n == 0 {
		return nil
	}
	result := make([]string, 0, n)
	for range s.buffers {
		result = append(result, "default output handler")
	}
	return result
}
