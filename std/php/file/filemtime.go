package file

import (
	"os"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// FilemtimeFunction 实现 filemtime 函数
// 获取文件的修改时间（Unix 时间戳）
type FilemtimeFunction struct{}

func NewFilemtimeFunction() data.FuncStmt {
	return &FilemtimeFunction{}
}

func (f *FilemtimeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	// 获取文件信息
	fileInfo, err := os.Stat(path)
	if err != nil {
		// 如果文件不存在或其他错误，返回 false
		return data.NewBoolValue(false), nil
	}

	// 返回修改时间的 Unix 时间戳
	return data.NewIntValue(int(fileInfo.ModTime().Unix())), nil
}

func (f *FilemtimeFunction) GetName() string {
	return "filemtime"
}

func (f *FilemtimeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "filename", 0, nil, nil),
	}
}

func (f *FilemtimeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "filename", 0, data.NewBaseType("string")),
	}
}
