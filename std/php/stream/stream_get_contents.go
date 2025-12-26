package stream

import (
	"io"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/core"
)

// StreamGetContentsFunction 实现 stream_get_contents 函数
type StreamGetContentsFunction struct{}

func NewStreamGetContentsFunction() data.FuncStmt {
	return &StreamGetContentsFunction{}
}

func (f *StreamGetContentsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取流资源
	streamValue, _ := ctx.GetIndexValue(0)
	if streamValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 从资源对象中获取 StreamInfo 或 StreamInfoFromReader
	var streamReader io.Reader
	if res, ok := streamValue.(*core.ResourceValue); ok {
		resource := res.GetResource()
		if resource == nil {
			return data.NewBoolValue(false), nil
		}
		// 检查是否是 StreamInfo 类型
		if info, ok := resource.(*StreamInfo); ok {
			streamReader = info
		} else if info, ok := resource.(*StreamInfoFromReader); ok {
			// 检查是否是 StreamInfoFromReader 类型（用于 proc_open 的管道）
			// 确保 Reader 不为 nil
			if info == nil || info.Reader == nil {
				return data.NewBoolValue(false), nil
			}
			streamReader = info
		} else {
			// 类型断言失败 - 尝试直接作为 io.Reader 使用
			if reader, ok := resource.(io.Reader); ok {
				streamReader = reader
			} else {
				return data.NewBoolValue(false), nil
			}
		}
	} else {
		return data.NewBoolValue(false), nil
	}

	// 将 streamReader 转换为 StreamInfo 接口
	var streamInfo interface {
		Read(p []byte) (int, error)
		IsClosed() bool
	}
	if info, ok := streamReader.(*StreamInfo); ok {
		streamInfo = info
	} else if info, ok := streamReader.(*StreamInfoFromReader); ok {
		streamInfo = info
	} else {
		return data.NewBoolValue(false), nil
	}

	// 检查流是否已关闭
	if streamInfo.IsClosed() {
		return data.NewBoolValue(false), nil
	}

	// 获取可选的 length 参数
	var maxLength int = -1 // -1 表示读取所有
	lengthValue, _ := ctx.GetIndexValue(1)
	if lengthValue != nil {
		// 检查是否是 NullValue（表示未提供参数）
		if _, ok := lengthValue.(*data.NullValue); !ok {
			if intVal, ok := lengthValue.(data.AsInt); ok {
				if length, err := intVal.AsInt(); err == nil && length >= 0 {
					maxLength = length
				}
			}
		}
	}

	// 获取可选的 offset 参数
	var offset int64 = -1 // -1 表示从当前位置开始
	offsetValue, _ := ctx.GetIndexValue(2)
	if offsetValue != nil {
		if intVal, ok := offsetValue.(data.AsInt); ok {
			if off, err := intVal.AsInt(); err == nil {
				offset = int64(off)
			}
		}
	}

	// 如果指定了 offset，先定位（仅对 StreamInfo 支持）
	if offset >= 0 {
		if seeker, ok := streamReader.(interface {
			Seek(offset int64, whence int) (int64, error)
		}); ok {
			_, err := seeker.Seek(offset, io.SeekStart)
			if err != nil {
				return data.NewBoolValue(false), nil
			}
		}
	}

	// 读取数据
	var content []byte
	var err error
	if maxLength < 0 {
		// 读取所有剩余数据
		// 对于管道，io.ReadAll 会阻塞直到 EOF（进程结束）
		content, err = io.ReadAll(streamReader)
		// io.ReadAll 在成功时返回 nil，在遇到 EOF 时也返回 nil
		// 只有在遇到其他错误时才返回非 nil 错误
		if err != nil {
			return data.NewBoolValue(false), nil
		}
	} else {
		// 读取指定长度的数据
		content = make([]byte, maxLength)
		n, readErr := streamReader.Read(content)
		// Read 可能返回 EOF，这是正常的（表示读取结束）
		if readErr != nil && readErr != io.EOF {
			return data.NewBoolValue(false), nil
		}
		content = content[:n]
	}

	return data.NewStringValue(string(content)), nil
}

func (f *StreamGetContentsFunction) GetName() string {
	return "stream_get_contents"
}

func (f *StreamGetContentsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "stream", 0, nil, nil),
		node.NewParameter(nil, "length", 1, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "offset", 2, node.NewIntLiteral(nil, "-1"), nil),
	}
}

func (f *StreamGetContentsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "stream", 0, data.NewBaseType("resource")),
		node.NewVariable(nil, "length", 1, data.NewBaseType("int")),
		node.NewVariable(nil, "offset", 2, data.NewBaseType("int")),
	}
}
