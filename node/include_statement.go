package node

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

// 缓存 include_once/require_once 已加载文件的返回值
var includeOnceCache = struct {
	mu    sync.Mutex
	files map[string]data.GetValue
}{
	files: make(map[string]data.GetValue),
}

// ClearIncludeCache 清空 include/require 返回值缓存，供开发模式热重载使用。
// 缓存为进程级全局，若不清理会导致重载时 require 命中旧值而跳过重新执行。
func ClearIncludeCache() {
	includeOnceCache.mu.Lock()
	defer includeOnceCache.mu.Unlock()
	includeOnceCache.files = make(map[string]data.GetValue)
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

// IncludeCore 统一的 include/require 逻辑，可被语句或函数重用。
//
// PHP 语义：require 在独立文件作用域执行，按名注入调用者（含 extract）已有变量。
// 注意：当前为单向注入（视图场景），并非完整共享符号表写回。
// require_once 仍缓存返回值；普通 require 每次重新执行（视图引擎需要）。
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
	filePath = utils.NormalizePhpFilePath(filePath)

	vm := ctx.GetVM()

	if once {
		includeOnceCache.mu.Lock()
		if cached, ok := includeOnceCache.files[filePath]; ok {
			includeOnceCache.mu.Unlock()
			return cached, nil
		}
		includeOnceCache.mu.Unlock()

		if vm.GetPhpFileCache(filePath) {
			return data.NewBoolValue(true), nil
		}
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

	v, acl := vm.LoadInCallerContext(ctx, filePath)
	if acl != nil {
		return nil, acl
	}

	if once {
		if !vm.GetPhpFileCache(filePath) {
			vm.SetPhpFileCache(filePath)
		}
		includeOnceCache.mu.Lock()
		if vv, ok := v.(data.Value); ok {
			includeOnceCache.files[filePath] = vv
		} else {
			includeOnceCache.files[filePath] = data.NewBoolValue(true)
		}
		includeOnceCache.mu.Unlock()
	}

	return v, nil
}
