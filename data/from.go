package data

// From 表示节点的来源
type From interface {
	// GetSource 返回节点的源代码（文件路径）
	GetSource() string
	// GetPosition 返回节点在源代码中的位置（字符偏移量）
	GetPosition() (start, end int)
	// GetStartPosition 返回起始位置（行号和列号）
	GetStartPosition() (line, char int)
	// GetEndPosition 返回结束位置（行号和列号）
	GetEndPosition() (line, char int)
	// GetRange 返回完整的位置范围
	GetRange() (startLine, startChar, endLine, endChar int)
	// ToLSPPosition 转换为 LSP 位置信息
	ToLSPPosition() (startLine, startChar, endLine, endChar int)
}

// BaseFrom 表示源代码来源
type BaseFrom struct {
	source string // 源代码
	start  int    // 开始位置
	end    int    // 结束位置
	// 行号和列号信息（可选，用于 LSP 支持）
	startLine int // 起始行号
	startChar int // 起始列号
	endLine   int // 结束行号
	endChar   int // 结束列号
}

// NewBaseFrom 创建一个新的源代码来源
func NewBaseFrom(source string, start, end int) *BaseFrom {
	return &BaseFrom{
		source: source,
		start:  start,
		end:    end,
		// 默认行号和列号信息
		startLine: 0,
		startChar: 0,
		endLine:   0,
		endChar:   0,
	}
}

// NewBaseFromWithPosition 创建一个带有行号和列号信息的源代码来源
func NewBaseFromWithPosition(source string, start, end, startLine, startChar, endLine, endChar int) *BaseFrom {
	return &BaseFrom{
		source:    source,
		start:     start,
		end:       end,
		startLine: startLine,
		startChar: startChar,
		endLine:   endLine,
		endChar:   endChar,
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

// GetStartPosition 返回起始位置（行号和列号）
func (s *BaseFrom) GetStartPosition() (line, char int) {
	return s.startLine, s.startChar
}

// GetEndPosition 返回结束位置（行号和列号）
func (s *BaseFrom) GetEndPosition() (line, char int) {
	return s.endLine, s.endChar
}

// GetRange 返回完整的位置范围
func (s *BaseFrom) GetRange() (startLine, startChar, endLine, endChar int) {
	return s.startLine, s.startChar, s.endLine, s.endChar
}

// ToLSPPosition 转换为 LSP 位置信息
func (s *BaseFrom) ToLSPPosition() (startLine, startChar, endLine, endChar int) {
	return s.startLine, s.startChar, s.endLine, s.endChar
}

// SetPosition 设置行号和列号信息
func (s *BaseFrom) SetPosition(startLine, startChar, endLine, endChar int) {
	s.startLine = startLine
	s.startChar = startChar
	s.endLine = endLine
	s.endChar = endChar
}
