package stream

import "os"

// OpenFile 以 PHP fopen 模式打开文件并返回 StreamInfo。
func OpenFile(filename, mode string) (*StreamInfo, error) {
	if mode == "" {
		mode = "r"
	}
	file, err := os.OpenFile(filename, parseMode(mode), 0644)
	if err != nil {
		return nil, err
	}
	return NewStreamInfo(file, mode), nil
}
