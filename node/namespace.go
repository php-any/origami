package node

import "github.com/php-any/origami/data"

// Namespace 表示命名空间节点
type Namespace struct {
	*Node      `pp:"-"`
	Name       string      // 命名空间名称
	Statements []Statement // 命名空间内的语句
}

// NewNamespace 创建一个新的命名空间节点
func NewNamespace(from data.From, name string, statements []Statement) *Namespace {
	return &Namespace{
		Node:       NewNode(from),
		Name:       name,
		Statements: statements,
	}
}

// GetName 返回命名空间名称
func (n *Namespace) GetName() string {
	return n.Name
}

// GetStatements 返回命名空间内的语句
func (n *Namespace) GetStatements() []Statement {
	return n.Statements
}

// GetValue 获取命名空间节点的值
func (n *Namespace) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	ctx.SetNamespace(n.Name)

	// 执行命名空间内的语句
	var value data.GetValue
	var ctl data.Control
	for _, stmt := range n.Statements {
		value, ctl = stmt.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
	}

	return value, nil
}

// Scope 表示作用域
type Scope struct {
	parent    *Scope                   // 父作用域
	namespace string                   // 当前命名空间
	variables map[string]data.GetValue // 变量表
}

// NewScope 创建一个新的作用域
func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:    parent,
		variables: make(map[string]data.GetValue),
	}
}

// SetNamespace 设置当前命名空间
func (s *Scope) SetNamespace(namespace string) {
	s.namespace = namespace
}

// GetNamespace 获取当前命名空间
func (s *Scope) GetNamespace() string {
	return s.namespace
}

// SetVariable 设置变量
func (s *Scope) SetVariable(name string, value data.GetValue) {
	s.variables[name] = value
}

// GetVariable 获取变量
func (s *Scope) GetVariable(name string) (data.GetValue, bool) {
	value, ok := s.variables[name]
	if !ok && s.parent != nil {
		// 在父作用域中查找
		return s.parent.GetVariable(name)
	}
	return value, ok
}
