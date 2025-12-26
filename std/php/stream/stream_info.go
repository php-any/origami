package stream

import (
	"io"
	"os"
	"sync"
)

// StreamInfo 存储流信息
type StreamInfo struct {
	File   *os.File
	Mode   string // 打开模式，如 "r", "w", "a" 等
	Closed bool
	mutex  sync.RWMutex
}

// NewStreamInfo 创建流信息
func NewStreamInfo(file *os.File, mode string) *StreamInfo {
	return &StreamInfo{
		File:   file,
		Mode:   mode,
		Closed: false,
	}
}

// Close 关闭流
func (s *StreamInfo) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.Closed {
		return nil
	}
	s.Closed = true
	if s.File != nil {
		return s.File.Close()
	}
	return nil
}

// IsClosed 检查流是否已关闭
func (s *StreamInfo) IsClosed() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Closed
}

// Read 读取数据
func (s *StreamInfo) Read(p []byte) (int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.Closed || s.File == nil {
		return 0, io.EOF
	}
	return s.File.Read(p)
}

// Write 写入数据
func (s *StreamInfo) Write(p []byte) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.Closed || s.File == nil {
		return 0, io.ErrClosedPipe
	}
	return s.File.Write(p)
}

// Seek 设置文件偏移量
func (s *StreamInfo) Seek(offset int64, whence int) (int64, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.Closed || s.File == nil {
		return 0, io.ErrClosedPipe
	}
	return s.File.Seek(offset, whence)
}

// ReadAt 从指定位置读取
func (s *StreamInfo) ReadAt(p []byte, off int64) (int, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.Closed || s.File == nil {
		return 0, io.EOF
	}
	return s.File.ReadAt(p, off)
}
