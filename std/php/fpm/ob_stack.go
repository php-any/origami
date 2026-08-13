package fpm

import "strings"

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
