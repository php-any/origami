package log

import (
	"fmt"
	"os"
	"time"

	"github.com/php-any/origami/data"
)

// 颜色常量
const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
)

// Log 日志结构体
type Log struct {
	output *os.File
}

// NewLog 创建新的日志实例
func NewLog() *Log {
	return &Log{
		output: os.Stdout,
	}
}

// SetOutput 设置输出流
func (l *Log) SetOutput(output *os.File) {
	l.output = output
}

// SetOutputFile 设置输出到文件
func (l *Log) SetOutputFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	l.output = file
	return nil
}

// formatMessage 格式化日志消息
func (l *Log) formatMessage(level, levelColor, msg string, args data.ArrayValue) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 检查args是否为空数组
	argsStr := ""
	if args.Value != nil && len(args.Value) > 0 {
		argsStr = fmt.Sprintf(" %v", args.AsString())
	}

	return fmt.Sprintf("%s[%s]%s %s %s%s%s\n",
		levelColor, level, ColorReset,
		timestamp, msg, argsStr, ColorReset)
}

// Fatal 致命错误级别日志 - 会终止程序
func (l *Log) Fatal(msg string, args data.ArrayValue) {
	formatted := l.formatMessage("FATAL", ColorBold+ColorRed, msg, args)
	fmt.Fprint(l.output, formatted)
	os.Exit(1)
}

// Error 错误级别日志
func (l *Log) Error(msg string, args data.ArrayValue) {
	formatted := l.formatMessage("ERROR", ColorRed, msg, args)
	fmt.Fprint(l.output, formatted)
}

// Warn 警告级别日志
func (l *Log) Warn(msg string, args data.ArrayValue) {
	formatted := l.formatMessage("WARN", ColorYellow, msg, args)
	fmt.Fprint(l.output, formatted)
}

// Notice 通知级别日志
func (l *Log) Notice(msg string, args data.ArrayValue) {
	formatted := l.formatMessage("NOTICE", ColorPurple, msg, args)
	fmt.Fprint(l.output, formatted)
}

// Info 信息级别日志
func (l *Log) Info(msg string, args data.ArrayValue) {
	formatted := l.formatMessage("INFO", ColorGreen, msg, args)
	fmt.Fprint(l.output, formatted)
}

// Debug 调试级别日志
func (l *Log) Debug(msg string, args data.ArrayValue) {
	formatted := l.formatMessage("DEBUG", ColorBlue, msg, args)
	fmt.Fprint(l.output, formatted)
}

// Trace 跟踪级别日志 - 最低优先级
func (l *Log) Trace(msg string, args data.ArrayValue) {
	formatted := l.formatMessage("TRACE", ColorCyan, msg, args)
	fmt.Fprint(l.output, formatted)
}
