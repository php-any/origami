package stream

import (
	"io"
	"sync"
)

// StreamInfoFromReader 从 io.ReadCloser 创建 StreamInfo
// 用于包装 proc_open 返回的管道
type StreamInfoFromReader struct {
	Reader io.ReadCloser
	Mode   string
	Closed bool
	mutex  sync.RWMutex
}

// NewStreamInfoFromReader 从 io.ReadCloser 创建流信息
func NewStreamInfoFromReader(reader io.ReadCloser, mode string) *StreamInfoFromReader {
	return &StreamInfoFromReader{
		Reader: reader,
		Mode:   mode,
		Closed: false,
	}
}

// Close 关闭流
func (s *StreamInfoFromReader) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.Closed {
		return nil
	}
	s.Closed = true
	if s.Reader != nil {
		return s.Reader.Close()
	}
	return nil
}

// IsClosed 检查流是否已关闭
func (s *StreamInfoFromReader) IsClosed() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Closed
}

// Read 读取数据
func (s *StreamInfoFromReader) Read(p []byte) (int, error) {
	s.mutex.RLock()
	closed := s.Closed
	reader := s.Reader
	s.mutex.RUnlock()

	if closed || reader == nil {
		return 0, io.EOF
	}
	// 在锁外执行实际的读取操作，避免阻塞
	return reader.Read(p)
}

// Write 写入数据（不支持）
func (s *StreamInfoFromReader) Write(p []byte) (int, error) {
	return 0, io.ErrClosedPipe
}

// Seek 设置文件偏移量（不支持）
func (s *StreamInfoFromReader) Seek(offset int64, whence int) (int64, error) {
	return 0, io.ErrClosedPipe
}

// ReadAt 从指定位置读取（不支持）
func (s *StreamInfoFromReader) ReadAt(p []byte, off int64) (int, error) {
	return 0, io.EOF
}
