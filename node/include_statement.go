package node

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

// 缓存 include_once/require_once 已加载文件
var includeOnceCache = struct {
	mu    sync.Mutex
	files map[string]struct{}
}{
	files: make(map[string]struct{}),
}

// IncludeStatement 表示 include/require/include_once/require_once 语句
type IncludeStatement struct {
	*Node
	Expr     data.GetValue
	Once     bool
	Required bool
}

func NewIncludeStatement(from data.From, expr data.GetValue, once bool, required bool) *IncludeStatement {
	return &IncludeStatement{
		Node:     NewNode(from),
		Expr:     expr,
		Once:     once,
		Required: required,
	}
}

func (s *IncludeStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	val, acl := s.Expr.GetValue(ctx)
	if acl != nil {
		return nil, acl
	}
	v, ok := val.(data.Value)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	return IncludeCore(ctx, v, s.Once, s.Required, s.from)
}

// IncludeCore 统一的 include/require 逻辑，可被语句或函数重用
func IncludeCore(ctx data.Context, pathVal data.Value, once bool, required bool, from data.From) (data.GetValue, data.Control) {
	var filePath string
	switch p := pathVal.(type) {
	case data.AsString:
		filePath = p.AsString()
	default:
		if required {
			return nil, utils.NewThrowf("require 文件失败: 非字符串路径")
		}
		return data.NewBoolValue(false), nil
	}

	if filePath == "" {
		if required {
			return nil, utils.NewThrowf("require 文件失败: 空路径")
		}
		return data.NewBoolValue(false), nil
	}

	if !filepath.IsAbs(filePath) {
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, utils.NewThrowf("include 文件失败: %s, 错误: %v", filePath, err)
		}
		filePath = filepath.Join(currentDir, filePath)
	}
	filePath = filepath.Clean(filePath)

	if once {
		includeOnceCache.mu.Lock()
		if _, ok := includeOnceCache.files[filePath]; ok {
			includeOnceCache.mu.Unlock()
			return data.NewBoolValue(true), nil
		}
		includeOnceCache.mu.Unlock()
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if required {
			return nil, utils.NewThrowf("require 文件失败: %s, 错误: %v", filePath, err)
		}
		return data.NewBoolValue(false), nil
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if required {
			return nil, utils.NewThrowf("require 文件失败: %s, 错误: %v", filePath, err)
		}
		return data.NewBoolValue(false), nil
	}
	if fileInfo.IsDir() {
		if required {
			return nil, utils.NewThrowf("require 文件失败: %s, 错误: 无法引入目录", filePath)
		}
		return data.NewBoolValue(false), nil
	}

	vm := ctx.GetVM()
	v, acl := vm.LoadAndRun(filePath)
	if acl != nil {
		if required {
			return nil, acl
		}
		return data.NewBoolValue(false), nil
	}

	if once {
		includeOnceCache.mu.Lock()
		includeOnceCache.files[filePath] = struct{}{}
		includeOnceCache.mu.Unlock()
	}

	return v, nil
}
