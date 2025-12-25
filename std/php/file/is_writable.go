package file

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// IsWritableFunction 实现 is_writable 函数
// 检查文件或目录是否可写
type IsWritableFunction struct{}

func NewIsWritableFunction() data.FuncStmt {
	return &IsWritableFunction{}
}

func (f *IsWritableFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取路径参数
	pathValue, _ := ctx.GetIndexValue(0)

	var path string
	switch p := pathValue.(type) {
	case data.AsString:
		path = p.AsString()
	default:
		return data.NewBoolValue(false), nil
	}

	// 检查路径是否为空
	if path == "" {
		return data.NewBoolValue(false), nil
	}

	// 检查文件或目录是否存在
	_, err := os.Stat(path)
	if err != nil {
		// 如果文件不存在，检查父目录是否可写（可以创建新文件）
		dir := path
		// 尝试获取父目录
		for i := len(dir) - 1; i >= 0; i-- {
			if dir[i] == os.PathSeparator || dir[i] == '/' {
				dir = dir[:i]
				break
			}
		}
		if dir == "" {
			dir = "."
		}
		// 检查父目录是否可写
		return data.NewBoolValue(isWritablePath(dir)), nil
	}

	// 检查是否可写
	return data.NewBoolValue(isWritablePath(path)), nil
}

// isWritablePath 检查路径是否可写
func isWritablePath(path string) bool {
	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

func (f *IsWritableFunction) GetName() string {
	return "is_writable"
}

func (f *IsWritableFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
	}
}

func (f *IsWritableFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
	}
}
