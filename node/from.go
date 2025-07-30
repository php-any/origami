package node

import "strconv"

// SourceFrom 表示源代码来源
type SourceFrom struct {
	source   *string // 源代码指针
	startPos int     // 起始位置
	endPos   int     // 结束位置
	line     int
	pos      int
}

// NewSourceFrom 创建一个新的源代码来源
func NewSourceFrom(source *string, startPos, endPos int, line int, pos int) *SourceFrom {
	return &SourceFrom{
		source:   source,
		startPos: startPos,
		endPos:   endPos,
		line:     line,
		pos:      pos,
	}
}

// GetSource 返回源代码
func (f *SourceFrom) GetSource() string {
	if f.source == nil {
		return ""
	}
	return *f.source + ":" + strconv.Itoa(f.line) + ":" + strconv.Itoa(f.pos)
}

// GetPosition 返回源代码中的位置
func (f *SourceFrom) GetPosition() (start, end int) {
	return f.startPos, f.endPos
}

func (f *SourceFrom) Line() int {
	return f.line
}

// TokenFrom 表示词法单元来源
type TokenFrom struct {
	*SourceFrom
}

// NewTokenFrom 创建一个新的词法单元来源
func NewTokenFrom(source *string, startPos, endPos, line, pos int) *TokenFrom {
	return &TokenFrom{
		SourceFrom: NewSourceFrom(source, startPos, endPos, line, pos),
	}
}

// NewTokenFrom 从当前 TokenFrom 创建一个新的词法单元来源
func (f *TokenFrom) NewTokenFrom(startPos, endPos int) *TokenFrom {
	return &TokenFrom{
		SourceFrom: &SourceFrom{
			source:   f.source,
			startPos: startPos,
			endPos:   endPos,
		},
	}
}

func (f *TokenFrom) String() string {
	return ""
}

// FileFrom 表示文件来源
type FileFrom struct {
	*SourceFrom
	fileName string // 文件名
}

// NewFileFrom 创建一个新的文件来源
func NewFileFrom(fileName string, source *string) *FileFrom {
	return &FileFrom{
		SourceFrom: &SourceFrom{
			source:   source,
			startPos: 0,
			endPos:   len(*source),
		},
		fileName: fileName,
	}
}

// NewTokenFrom 从当前 FileFrom 创建一个新的词法单元来源
func (f *FileFrom) NewTokenFrom(startPos, endPos int) *TokenFrom {
	return &TokenFrom{
		SourceFrom: &SourceFrom{
			source:   f.source,
			startPos: startPos,
			endPos:   endPos,
		},
	}
}

// GetFileName 返回文件名
func (f *FileFrom) GetFileName() string {
	return f.fileName
}
