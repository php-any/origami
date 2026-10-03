package stream

import (
	"io"
	"os"
	"sync"
)

// StreamInfo 存储流信息
type StreamInfo struct {
	File          *os.File
	Mode          string // 打开模式，如 "r", "w", "a" 等
	Closed        bool
	mutex         sync.RWMutex
	stopCancel    func() bool
	temporaryPath string
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
	if s.Closed {
		s.mutex.Unlock()
		return nil
	}
	s.Closed = true
	file, stop := s.File, s.stopCancel
	s.stopCancel = nil
	s.mutex.Unlock()
	if stop != nil {
		stop()
	}
	if file != nil {
		// 对于标准流（stdin/stdout/stderr），不要真正关闭它们
		// 只标记为已关闭
		if file == os.Stdin || file == os.Stdout || file == os.Stderr {
			return nil
		}
		err := file.Close()
		if s.temporaryPath != "" {
			_ = os.Remove(s.temporaryPath)
		}
		return err
	}
	return nil
}

func (s *StreamInfo) BindRequestCancel(stop func() bool) {
	s.mutex.Lock()
	closed := s.Closed
	if !closed {
		s.stopCancel = stop
	}
	s.mutex.Unlock()
	if closed {
		stop()
	}
}

func (s *StreamInfo) openFile() *os.File {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if s.Closed {
		return nil
	}
	return s.File
}

// IsClosed 检查流是否已关闭
func (s *StreamInfo) IsClosed() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Closed
}

// Read 读取数据
func (s *StreamInfo) Read(p []byte) (int, error) {
	file := s.openFile()
	if file == nil {
		return 0, io.EOF
	}
	return file.Read(p)
}

// Write 写入数据
func (s *StreamInfo) Write(p []byte) (int, error) {
	file := s.openFile()
	if file == nil {
		return 0, io.ErrClosedPipe
	}
	return file.Write(p)
}

// Flush 将缓冲数据刷新到底层（如 os.File.Sync），便于立即看到 stdout 等输出
func (s *StreamInfo) Flush() error {
	file := s.openFile()
	if file == nil {
		return io.ErrClosedPipe
	}
	return file.Sync()
}

// Seek 设置文件偏移量
func (s *StreamInfo) Seek(offset int64, whence int) (int64, error) {
	file := s.openFile()
	if file == nil {
		return 0, io.ErrClosedPipe
	}
	return file.Seek(offset, whence)
}

// ReadAt 从指定位置读取
func (s *StreamInfo) ReadAt(p []byte, off int64) (int, error) {
	file := s.openFile()
	if file == nil {
		return 0, io.EOF
	}
	return file.ReadAt(p, off)
}
