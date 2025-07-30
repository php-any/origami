package data

// From 表示节点的来源
type From interface {
	// GetSource 返回节点的源代码
	GetSource() string
	// GetPosition 返回节点在源代码中的位置
	GetPosition() (start, end int)
}

// BaseFrom 表示源代码来源
type BaseFrom struct {
	source string // 源代码
	start  int    // 开始位置
	end    int    // 结束位置
}

// NewBaseFrom 创建一个新的源代码来源
func NewBaseFrom(source string, start, end int) *BaseFrom {
	return &BaseFrom{
		source: source,
		start:  start,
		end:    end,
	}
}

// GetSource 返回源代码
func (s *BaseFrom) GetSource() string {
	return s.source
}

// GetPosition 返回位置信息
func (s *BaseFrom) GetPosition() (start, end int) {
	return s.start, s.end
}
