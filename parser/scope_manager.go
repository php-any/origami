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

// Scope 表示作用域接口
type Scope interface {
	// 获取父作用域
	GetParent() Scope
	// 设置父作用域
	SetParent(parent Scope)
	// 添加变量
	AddVariable(name string, ty data.Types, from data.From) data.Variable
	// 获取变量
	GetVariable(name string) (data.Variable, bool)
	// 获取所有变量
	GetVariables() []data.Variable
	// 获取下一个变量索引
	GetNextIndex() int
	// 设置下一个变量索引
	SetNextIndex(index int)
	// 是否是 lambda 作用域
	IsLambda() bool
	// 设置是否是 lambda 作用域
	SetLambda(isLambda bool)
	// 设置变量
	SetVariable(name string, variable data.Variable)
	// 字符串表示
	String() string
}

// DefaultScope 默认作用域实现
type DefaultScope struct {
	parent    Scope                    // 父作用域
	variables map[string]data.Variable // 变量表
	nextIndex int                      // 下一个变量索引
	isLambda  bool
}

// ScopeFactory 作用域工厂函数类型
type ScopeFactory func(parent Scope, isLambda bool) Scope

// 全局作用域工厂函数，默认为 DefaultScope
var globalScopeFactory ScopeFactory = func(parent Scope, isLambda bool) Scope {
	return &DefaultScope{
		parent:    parent,
		variables: make(map[string]data.Variable),
		isLambda:  isLambda,
	}
}

// ScopeManager 表示作用域管理器
type ScopeManager struct {
	scopes  []Scope // 作用域栈
	current Scope   // 当前作用域
}

// NewScopeManager 创建一个新的作用域管理器
func NewScopeManager() *ScopeManager {
	// 使用全局工厂函数创建全局作用域
	globalScope := globalScopeFactory(nil, false)
	return &ScopeManager{
		scopes:  []Scope{globalScope},
		current: globalScope,
	}
}

// NewScope 创建新作用域
func (m *ScopeManager) NewScope(isLambda bool) Scope {
	// 使用全局工厂函数创建作用域
	scope := globalScopeFactory(m.current, isLambda)
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
func (m *ScopeManager) CurrentScope() Scope {
	return m.current
}

// AddVariable 在当前作用域中添加变量
func (s *DefaultScope) AddVariable(name string, ty data.Types, from data.From) data.Variable {
	if name[0:1] == "$" {
		name = name[1:]
	}
	if v, exists := s.variables[name]; exists {
		return v
	}
	s.variables[name] = node.NewVariable(from, name, s.nextIndex, ty)
	s.nextIndex++
	return s.variables[name]
}

// LookupVariable 查找变量（仅当前作用域）
func (m *ScopeManager) LookupVariable(name string) data.Variable {
	scope := m.current
	if scope != nil {
		if v, ok := scope.GetVariable(name); ok {
			return v
		}
	}
	return nil
}

// InLambdaScope 当前是否处于箭头/lambda 作用域（含嵌套）
func (m *ScopeManager) InLambdaScope() bool {
	for s := m.current; s != nil; s = s.GetParent() {
		if s.IsLambda() {
			return true
		}
	}
	return false
}

// CaptureFromEnclosing 在箭头/lambda 中解析自由变量：
// 从祖先作用域找到同名变量后，注入到当前作用域以及中间各层 lambda，
// 以便 parent 映射能逐层捕获（Filament 嵌套 fn 捕获方法参数 $data）。
func (m *ScopeManager) CaptureFromEnclosing(name string, from data.From) data.Variable {
	if name != "" && name[0] == '$' {
		name = name[1:]
	}
	if m.current == nil || !m.InLambdaScope() {
		return nil
	}
	if v, ok := m.current.GetVariable(name); ok {
		return v
	}

	// 收集从 current 的 parent 到找到变量的作用域链
	type frame struct {
		scope Scope
		v     data.Variable
	}
	var found data.Variable
	var between []Scope // current 的祖先，直到（不含）定义处
	for s := m.current.GetParent(); s != nil; s = s.GetParent() {
		if v, ok := s.GetVariable(name); ok {
			found = v
			break
		}
		between = append(between, s)
	}
	if found == nil {
		return nil
	}

	ty := found.GetType()
	// 先注入中间层（靠近定义处的先注入，便于外层箭头 parent 映射）
	for i := len(between) - 1; i >= 0; i-- {
		s := between[i]
		if !s.IsLambda() {
			continue
		}
		if _, ok := s.GetVariable(name); ok {
			continue
		}
		s.AddVariable(name, ty, from)
	}
	// 再注入当前 lambda
	return m.current.AddVariable(name, ty, from)
}

// LookupParentVariable 查找变量, 在父级域中查找
func (m *ScopeManager) LookupParentVariable(name string) data.Variable {
	scope := m.current
	if scope != nil && scope.GetParent() != nil {
		if v, ok := scope.GetParent().GetVariable(name); ok {
			return v
		}
	}
	return nil
}

// GetVariables 获取当前作用域中的所有变量
func (s *DefaultScope) GetVariables() []data.Variable {
	variables := make([]data.Variable, s.nextIndex)
	for _, v := range s.variables {
		variables[v.GetIndex()] = v
	}
	return variables
}

// String 返回作用域的字符串表示（用于调试）
func (s *DefaultScope) String() string {
	return fmt.Sprintf("Scope{variables: %v, nextIndex: %d}", s.variables, s.nextIndex)
}

// DefaultScope 实现 Scope 接口的所有方法

// GetParent 获取父作用域
func (s *DefaultScope) GetParent() Scope {
	return s.parent
}

// SetParent 设置父作用域
func (s *DefaultScope) SetParent(parent Scope) {
	s.parent = parent
}

// GetVariable 获取变量
func (s *DefaultScope) GetVariable(name string) (data.Variable, bool) {
	v, exists := s.variables[name]
	return v, exists
}

// GetNextIndex 获取下一个变量索引
func (s *DefaultScope) GetNextIndex() int {
	return s.nextIndex
}

// SetNextIndex 设置下一个变量索引
func (s *DefaultScope) SetNextIndex(index int) {
	s.nextIndex = index
}

// IsLambda 是否是 lambda 作用域
func (s *DefaultScope) IsLambda() bool {
	return s.isLambda
}

// SetLambda 设置是否是 lambda 作用域
func (s *DefaultScope) SetLambda(isLambda bool) {
	s.isLambda = isLambda
}

// SetVariable 设置变量
func (s *DefaultScope) SetVariable(name string, variable data.Variable) {
	s.variables[name] = variable
}

// SetGlobalScopeFactory 设置全局作用域工厂函数
func SetGlobalScopeFactory(factory ScopeFactory) {
	globalScopeFactory = factory
}

// GetGlobalScopeFactory 获取全局作用域工厂函数
func GetGlobalScopeFactory() ScopeFactory {
	return globalScopeFactory
}
