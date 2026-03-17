package directory

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// FilesystemIteratorConstructMethod 实现 FilesystemIterator::__construct
type FilesystemIteratorConstructMethod struct {
	instance *FilesystemIteratorClass // 持有实例引用（skill: php-class-state-sharing-pattern）
}

func (m *FilesystemIteratorConstructMethod) GetName() string            { return "__construct" }
func (m *FilesystemIteratorConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *FilesystemIteratorConstructMethod) GetIsStatic() bool          { return false }

func (m *FilesystemIteratorConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "directory", 0, nil, data.String{}),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(FSI_DEFAULT_FLAGS), data.Int{}),
	}
}

func (m *FilesystemIteratorConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "directory", 0, data.String{}),
		node.NewVariable(nil, "flags", 1, data.Int{}),
	}
}

func (m *FilesystemIteratorConstructMethod) GetReturnType() data.Types { return nil }

func (m *FilesystemIteratorConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取 $directory 参数
	pathValue, _ := ctx.GetIndexValue(0)
	if pathValue == nil {
		return nil, utils.NewThrow(errors.New("FilesystemIterator::__construct() expects parameter 1 to be string"))
	}

	var path string
	switch p := pathValue.(type) {
	case data.AsString:
		path = p.AsString()
	default:
		path = pathValue.AsString()
	}

	if path == "" {
		return nil, utils.NewThrow(errors.New("FilesystemIterator::__construct(): Directory name must not be empty"))
	}

	// 获取 $flags 参数（默认 FSI_DEFAULT_FLAGS）
	flags := FSI_DEFAULT_FLAGS
	if flagsValue, _ := ctx.GetIndexValue(1); flagsValue != nil {
		if iv, ok := flagsValue.(interface{ AsInt() (int, error) }); ok {
			if n, err := iv.AsInt(); err == nil {
				flags = n
			}
		}
	}

	// 通过 m.instance 直接设置状态（skill 关键模式）
	m.instance.path = path
	m.instance.flags = flags
	if err := m.instance.loadEntries(); err != nil {
		return nil, utils.NewThrow(err)
	}

	return nil, nil
}
