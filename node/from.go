package node

import (
	"strings"
)

// TokenFrom 表示基于 token 的源代码来源信息
type TokenFrom struct {
	filePath  *string // 文件路径指针（避免重复存储）
	startPos  int     // 起始字符偏移量
	endPos    int     // 结束字符偏移量
	startLine int     // 起始行号
	startChar int     // 起始列号
	endLine   int     // 结束行号
	endChar   int     // 结束列号
}

// NewTokenFrom 创建一个新的 TokenFrom
func NewTokenFrom(filePath *string, startPos, endPos, startLine, startChar int) *TokenFrom {
	return &TokenFrom{
		filePath:  filePath,
		startPos:  startPos,
		endPos:    endPos,
		startLine: startLine,
		startChar: startChar,
		endLine:   startLine, // 默认结束位置与开始位置相同
		endChar:   startChar, // 默认结束位置与开始位置相同
	}
}

// SetEndPosition 设置结束位置（在解析过程中调用）
func (tf *TokenFrom) SetEndPosition(endLine, endChar int) {
	tf.endLine = endLine
	tf.endChar = endChar
}

// GetFilePath 返回文件路径
func (tf *TokenFrom) GetFilePath() string {
	if tf.filePath == nil {
		return ""
	}
	return *tf.filePath
}

// GetStartPosition 返回起始位置
func (tf *TokenFrom) GetStartPosition() (line, char int) {
	return tf.startLine, tf.startChar
}

// GetEndPosition 返回结束位置
func (tf *TokenFrom) GetEndPosition() (line, char int) {
	return tf.endLine, tf.endChar
}

// GetRange 返回位置范围
func (tf *TokenFrom) GetRange() (startLine, startChar, endLine, endChar int) {
	return tf.startLine, tf.startChar, tf.endLine, tf.endChar
}

// GetOffsetRange 返回偏移量范围
func (tf *TokenFrom) GetOffsetRange() (start, end int) {
	return tf.startPos, tf.endPos
}

// GetSource 实现 data.From 接口，返回文件路径
func (tf *TokenFrom) GetSource() string {
	return tf.GetFilePath()
}

// GetPosition 实现 data.From 接口，返回位置范围
func (tf *TokenFrom) GetPosition() (start, end int) {
	return tf.startPos, tf.endPos
}

// CalculateEndPosition 根据内容计算结束位置（用于多行内容）
func (tf *TokenFrom) CalculateEndPosition(content string) {
	if content == "" {
		return
	}

	// 计算结束位置
	tf.endLine, tf.endChar = tf.calculateLineAndChar(content, tf.endPos)
}

// calculateLineAndChar 根据偏移量计算精确的行号和列号
func (tf *TokenFrom) calculateLineAndChar(content string, offset int) (line, char int) {
	if offset < 0 || offset > len(content) {
		return 0, 0
	}

	// 计算到指定偏移量之前的行数
	contentBefore := content[:offset]
	lines := strings.Split(contentBefore, "\n")
	line = len(lines)

	// 计算当前行的列号
	if line > 0 {
		lastLine := lines[line-1]
		char = len(lastLine)
	} else {
		char = 0
	}

	return line, char
}

// ToLSPPosition 转换为 LSP 位置信息
func (tf *TokenFrom) ToLSPPosition() (startLine, startChar, endLine, endChar int) {
	return tf.startLine, tf.startChar, tf.endLine, tf.endChar
}

// IsValid 检查位置信息是否有效
func (tf *TokenFrom) IsValid() bool {
	return tf.filePath != nil && tf.startPos >= 0 && tf.endPos >= tf.startPos
}
