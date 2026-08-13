package core

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// HeadersSentFunction 实现 headers_sent 函数
// headers_sent(string &$file = null, int &$line = null): bool
// 检查 HTTP 响应头是否已经发送
// 如果头已发送，返回 true，否则返回 false
// 如果提供了引用参数，还会返回文件名和行号
type HeadersSentFunction struct{}

func NewHeadersSentFunction() data.FuncStmt {
	return &HeadersSentFunction{}
}

func (f *HeadersSentFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// headers_sent 函数在 CLI 模式下通常返回 false
	// 因为 CLI 模式下没有 HTTP 响应头
	// 这里简化实现，总是返回 false（表示头未发送）
	return data.NewBoolValue(false), nil
}

func (f *HeadersSentFunction) GetName() string {
	return "headers_sent"
}

func (f *HeadersSentFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "file", 0, data.NewNullValue(), data.Mixed{}),
		node.NewParameter(nil, "line", 1, data.NewIntValue(0), data.Mixed{}),
	}
}

func (f *HeadersSentFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "file", 0, data.Mixed{}),
		node.NewVariable(nil, "line", 1, data.Mixed{}),
	}
}
