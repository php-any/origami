package parser

import (
	"fmt"
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// Variable 表示变量信息
type Variable struct {
	Name     string // 变量名
	Index    int    // 变量在作用域中的索引
	IsParam  bool   // 是否是函数参数
	IsGlobal bool   // 是否是全局变量
}

// Scope 表示作用域
type Scope struct {
	parent    *Scope                   // 父作用域
	variables map[string]data.Variable // 变量表
	nextIndex int                      // 下一个变量索引
	isLambda  bool
}

// ScopeManager 表示作用域管理器
type ScopeManager struct {
	scopes  []*Scope // 作用域栈
	current *Scope   // 当前作用域
}

// NewScopeManager 创建一个新的作用域管理器
func NewScopeManager() *ScopeManager {
	// 创建全局作用域
	globalScope := &Scope{
		variables: make(map[string]data.Variable),
	}
	return &ScopeManager{
		scopes:  []*Scope{globalScope},
		current: globalScope,
	}
}

// NewScope 创建新作用域
func (m *ScopeManager) NewScope(isLambda bool) *Scope {
	scope := &Scope{
		parent:    m.current,
		variables: make(map[string]data.Variable),
		isLambda:  isLambda,
	}
	m.scopes = append(m.scopes, scope)
	m.current = scope
	return scope
}

// PopScope 弹出当前作用域
func (m *ScopeManager) PopScope() {
	if len(m.scopes) > 1 {
		m.scopes = m.scopes[:len(m.scopes)-1]
		m.current = m.scopes[len(m.scopes)-1]
	}
}

// CurrentScope 获取当前作用域
func (m *ScopeManager) CurrentScope() *Scope {
	return m.current
}

// AddVariable 在当前作用域中添加变量
func (s *Scope) AddVariable(name string, ty data.Types, from data.From) int {
	if v, exists := s.variables[name]; exists {
		return v.GetIndex()
	}
	s.variables[name] = node.NewVariable(from, name, s.nextIndex, ty)
	s.nextIndex++
	return s.variables[name].GetIndex()
}

// LookupVariable 查找变量
func (m *ScopeManager) LookupVariable(name string) data.Variable {
	scope := m.current
	if scope != nil {
		if v, ok := scope.variables[name]; ok {
			return v
		}
	}
	return nil
}

// LookupParentVariable 查找变量, 在父级域中查找
func (m *ScopeManager) LookupParentVariable(name string) data.Variable {
	scope := m.current
	if scope != nil && scope.parent != nil {
		if v, ok := scope.parent.variables[name]; ok {
			return v
		}
	}
	return nil
}

// GetVariables 获取当前作用域中的所有变量
func (s *Scope) GetVariables() []data.Variable {
	variables := make([]data.Variable, s.nextIndex)
	for _, v := range s.variables {
		variables[v.GetIndex()] = v
	}
	return variables
}

// String 返回作用域的字符串表示（用于调试）
func (s *Scope) String() string {
	return fmt.Sprintf("Scope{variables: %v, nextIndex: %d}", s.variables, s.nextIndex)
}
