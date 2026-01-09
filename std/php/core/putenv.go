package core

import (
	"os"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PutenvFunction 实现 putenv 函数
// 设置环境变量的值
type PutenvFunction struct{}

func NewPutenvFunction() data.FuncStmt {
	return &PutenvFunction{}
}

func (f *PutenvFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数：环境变量赋值字符串，格式为 "KEY=VALUE"
	assignmentValue, _ := ctx.GetIndexValue(0)
	if assignmentValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 将参数转换为字符串
	var assignment string
	if str, ok := assignmentValue.(data.AsString); ok {
		assignment = str.AsString()
	} else {
		assignment = assignmentValue.AsString()
	}

	// 检查字符串是否为空
	if assignment == "" {
		return data.NewBoolValue(false), nil
	}

	// 解析 "KEY=VALUE" 格式
	parts := strings.SplitN(assignment, "=", 2)
	if len(parts) != 2 {
		// 如果格式不正确，返回 false
		return data.NewBoolValue(false), nil
	}

	key := parts[0]
	value := parts[1]

	// 设置环境变量
	err := os.Setenv(key, value)
	if err != nil {
		return data.NewBoolValue(false), nil
	}

	// putenv 函数返回 true 表示成功
	return data.NewBoolValue(true), nil
}

func (f *PutenvFunction) GetName() string {
	return "putenv"
}

func (f *PutenvFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "assignment", 0, nil, data.String{}),
	}
}

func (f *PutenvFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "assignment", 0, data.String{}),
	}
}
