package directory

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// DirectoryIteratorConstructMethod 实现 DirectoryIterator::__construct
type DirectoryIteratorConstructMethod struct{}

func (m *DirectoryIteratorConstructMethod) GetName() string { return "__construct" }

func (m *DirectoryIteratorConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }

func (m *DirectoryIteratorConstructMethod) GetIsStatic() bool { return false }

func (m *DirectoryIteratorConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "directory", 0, nil, data.String{}),
	}
}

func (m *DirectoryIteratorConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "directory", 0, data.String{}),
	}
}

func (m *DirectoryIteratorConstructMethod) GetReturnType() data.Types { return nil }

func (m *DirectoryIteratorConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取目录路径参数
	pathValue, _ := ctx.GetIndexValue(0)
	if pathValue == nil {
		return nil, utils.NewThrow(errors.New("DirectoryIterator::__construct() expects parameter 1 to be string"))
	}

	var path string
	switch p := pathValue.(type) {
	case data.AsString:
		path = p.AsString()
	default:
		path = pathValue.AsString()
	}

	// 检查路径是否为空
	if path == "" {
		return nil, utils.NewThrow(errors.New("DirectoryIterator::__construct(): Directory name must not be empty"))
	}

	// 创建 DirectoryIteratorData
	iterData, err := NewDirectoryIteratorData(path)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	// 将 DirectoryIteratorData 存储到实例属性中
	if objCtx, ok := ctx.(*data.ClassMethodContext); ok {
		// 存储迭代器数据到 ObjectValue 的实例属性中
		objCtx.ObjectValue.SetProperty("_iterator", iterData)
	}

	return nil, nil
}
