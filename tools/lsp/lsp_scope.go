package main

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/parser"
)

// LspScope 是专门为 LSP 服务器设计的作用域实现
// 它与 DefaultScope 保持完全一致的结构和行为
type LspScope struct {
	parent    parser.Scope             // 父作用域
	variables map[string]data.Variable // 变量表
	nextIndex int                      // 下一个变量索引
	isLambda  bool
}

// NewLspScope 创建一个新的 LSP 作用域 - 与 DefaultScope 保持一致
func NewLspScope(parent parser.Scope, scopeName, scopeType, filePath string) *LspScope {
	return &LspScope{
		parent:    parent,
		variables: make(map[string]data.Variable),
		nextIndex: 0,
		isLambda:  false,
	}
}

// 实现 parser.Scope 接口的所有方法 - 与 DefaultScope 完全一致

// GetParent 获取父作用域
func (s *LspScope) GetParent() parser.Scope {
	return s.parent
}

// SetParent 设置父作用域
func (s *LspScope) SetParent(parent parser.Scope) {
	s.parent = parent
}

// AddVariable 添加变量 - 与 DefaultScope 保持一致
func (s *LspScope) AddVariable(name string, ty data.Types, from data.From) data.Variable {
	if name[0:1] == "$" {
		name = name[1:]
	}
	if v, exists := s.variables[name]; exists {
		return v
	}

	// LSP 特殊处理：如果类型为 nil，创建 LspTypes
	if ty == nil {
		ty = &data.LspTypes{
			Types: []data.Types{},
		}
	}

	s.variables[name] = node.NewVariable(from, name, s.nextIndex, ty)
	s.nextIndex++
	return s.variables[name]
}

// GetVariable 获取变量 - 与 DefaultScope 保持一致
func (s *LspScope) GetVariable(name string) (data.Variable, bool) {
	v, exists := s.variables[name]
	return v, exists
}

// GetVariables 获取所有变量 - 与 DefaultScope 保持一致
func (s *LspScope) GetVariables() []data.Variable {
	variables := make([]data.Variable, s.nextIndex)
	for _, v := range s.variables {
		variables[v.GetIndex()] = v
	}
	return variables
}

// GetNextIndex 获取下一个变量索引 - 与 DefaultScope 保持一致
func (s *LspScope) GetNextIndex() int {
	return s.nextIndex
}

// SetNextIndex 设置下一个变量索引 - 与 DefaultScope 保持一致
func (s *LspScope) SetNextIndex(index int) {
	s.nextIndex = index
}

// IsLambda 是否是 lambda 作用域 - 与 DefaultScope 保持一致
func (s *LspScope) IsLambda() bool {
	return s.isLambda
}

// SetLambda 设置是否是 lambda 作用域 - 与 DefaultScope 保持一致
func (s *LspScope) SetLambda(isLambda bool) {
	s.isLambda = isLambda
}

// SetVariable 设置变量 - 与 DefaultScope 保持一致
func (s *LspScope) SetVariable(name string, variable data.Variable) {
	s.variables[name] = variable
}

// String 返回作用域的字符串表示 - 与 DefaultScope 保持一致
func (s *LspScope) String() string {
	return fmt.Sprintf("Scope{variables: %v, nextIndex: %d}", s.variables, s.nextIndex)
}
