package http

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/runtime"
)

type AppFunction struct{}

func NewAppFunction() data.FuncStmt { return &AppFunction{} }

func (h *AppFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取 request 和 response 参数
	requestValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数: request"))
	}

	responseValue, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少参数: response"))
	}

	// 获取可选的 filePath 参数，默认为 "./main.zy"
	var filePath string
	if filePathValue, ok := ctx.GetIndexValue(2); ok {
		if pathStr, ok := filePathValue.(data.AsString); ok {
			filePath = pathStr.AsString()
		}
	}
	if filePath == "" {
		filePath = "./main.zy"
	}

	// 获取当前工作目录
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("获取当前目录失败: %v", err))
	}

	// 如果路径是相对路径，则基于当前目录解析
	if !filepath.IsAbs(filePath) {
		filePath = filepath.Join(currentDir, filePath)
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("文件不存在: %s", filePath))
	}

	// 检查文件是否为目录
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("无法访问文件: %s", filePath))
	}
	if fileInfo.IsDir() {
		return nil, data.NewErrorThrow(nil, errors.New("无法引入目录"))
	}

	// 获取 VM 实例，使用 LoadAndRun 加载并执行文件（参考 std/include.go）
	vm := runtime.NewTempVM(ctx.GetVM())
	_, acl := vm.LoadAndRun(filePath)
	if acl != nil {
		return nil, acl
	}

	// 查找 app 函数
	fullName := "App\\main"
	fn, ok := vm.GetFunc(fullName)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("未找到 app($request, $response) 函数"))
	}

	// 调用 app($request, $response) - 参考 NextHandler 的简洁方式
	vars := fn.GetVariables()
	if len(vars) < 2 {
		return nil, data.NewErrorThrow(nil, errors.New("app 函数需要至少 2 个参数: $request, $response"))
	}

	fnCtx := vm.CreateContext(vars)
	fnCtx.SetVariableValue(vars[0], requestValue)
	fnCtx.SetVariableValue(vars[1], responseValue)

	return fn.Call(fnCtx)
}

func (h *AppFunction) GetName() string            { return "Net\\Http\\app" }
func (h *AppFunction) GetModifier() data.Modifier { return data.ModifierPublic }
func (h *AppFunction) GetIsStatic() bool          { return true }
func (h *AppFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "request", 0, nil, nil),
		node.NewParameter(nil, "response", 1, nil, nil),
		node.NewParameter(nil, "filePath", 2, data.NewStringValue("./main.zy"), data.String{}),
	}
}
func (h *AppFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "request", 0, nil),
		node.NewVariable(nil, "response", 1, nil),
		node.NewVariable(nil, "filePath", 2, nil),
	}
}
func (h *AppFunction) GetReturnType() data.Types { return data.NewBaseType("mixed") }
