package main

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ASTWalker 用于遍历 AST 节点
type ASTWalker struct {
	functions []FunctionInfo
	classes   []ClassInfo
	variables []VariableInfo
}

// FunctionInfo 存储函数信息
type FunctionInfo struct {
	Name     string
	Position int
	Line     int
	Node     *node.FunctionStatement
}

// ClassInfo 存储类信息
type ClassInfo struct {
	Name     string
	Position int
	Line     int
	Node     *node.ClassStatement
	Methods  []MethodInfo
}

// MethodInfo 存储方法信息
type MethodInfo struct {
	Name      string
	Position  int
	Line      int
	ClassName string
	Method    data.Method
}

// VariableInfo 存储变量信息
type VariableInfo struct {
	Name     string
	Position int
	Line     int
}

// NewASTWalker 创建新的 AST 遍历器
func NewASTWalker() *ASTWalker {
	return &ASTWalker{
		functions: make([]FunctionInfo, 0),
		classes:   make([]ClassInfo, 0),
		variables: make([]VariableInfo, 0),
	}
}

// Walk 遍历整个 AST
func (w *ASTWalker) Walk(program *node.Program) {
	if program == nil {
		return
	}

	// 遍历所有顶层语句
	for _, stmt := range program.Statements {
		w.walkStatement(stmt)
	}
}

// walkStatement 遍历语句节点
func (w *ASTWalker) walkStatement(stmt node.Statement) {
	if stmt == nil {
		return
	}

	switch s := stmt.(type) {
	case *node.FunctionStatement:
		w.processFunctionStatement(s)
	case *node.ClassStatement:
		w.processClassStatement(s)
	case *node.Namespace:
		// 遍历命名空间内的语句
		for _, nsStmt := range s.Statements {
			w.walkStatement(nsStmt)
		}
	case *node.IfStatement:
		// 遍历 if 语句的分支
		w.walkStatements(s.ThenBranch)
		for _, elseIf := range s.ElseIf {
			w.walkStatements(elseIf.ThenBranch)
		}
		w.walkStatements(s.ElseBranch)
	}
}

// walkStatements 遍历语句列表
func (w *ASTWalker) walkStatements(statements []data.GetValue) {
	for _, stmt := range statements {
		if nodeStmt, ok := stmt.(node.Statement); ok {
			w.walkStatement(nodeStmt)
		}
	}
}

// processFunctionStatement 处理函数声明
func (w *ASTWalker) processFunctionStatement(fn *node.FunctionStatement) {
	if fn == nil || fn.Name == "" {
		return
	}

	position := 0
	line := 1
	if fn.GetFrom() != nil {
		start, end := fn.GetFrom().GetPosition()
		position = start
		// 简化的行号计算，实际可能需要更复杂的逻辑
		line = position / 100 // 假设每行平均100个字符
		if line < 1 {
			line = 1
		}
		_ = end // 忽略 end 值
	}

	funcInfo := FunctionInfo{
		Name:     fn.Name,
		Position: position,
		Line:     line,
		Node:     fn,
	}

	w.functions = append(w.functions, funcInfo)

	// 遍历函数体
	w.walkStatements(fn.Body)
}

// processClassStatement 处理类声明
func (w *ASTWalker) processClassStatement(cls *node.ClassStatement) {
	if cls == nil || cls.Name == "" {
		return
	}

	position := 0
	line := 1
	if cls.GetFrom() != nil {
		start, end := cls.GetFrom().GetPosition()
		position = start
		line = position / 100
		if line < 1 {
			line = 1
		}
		_ = end // 忽略 end 值
	}

	classInfo := ClassInfo{
		Name:     cls.Name,
		Position: position,
		Line:     line,
		Node:     cls,
		Methods:  make([]MethodInfo, 0),
	}

	// 遍历类的方法
	for methodName, method := range cls.Methods {
		methodPosition := 0
		methodLine := 1
		// data.Method 接口没有 GetFrom 方法，使用默认位置

		methodInfo := MethodInfo{
			Name:      methodName,
			Position:  methodPosition,
			Line:      methodLine,
			ClassName: cls.Name,
			Method:    method,
		}

		classInfo.Methods = append(classInfo.Methods, methodInfo)

		// data.Method 接口没有 GetBody 方法，跳过方法体遍历
		// 如果需要遍历方法体，需要通过其他方式获取
	}

	w.classes = append(w.classes, classInfo)
}

// GetFunctions 获取所有函数信息
func (w *ASTWalker) GetFunctions() []FunctionInfo {
	return w.functions
}

// GetClasses 获取所有类信息
func (w *ASTWalker) GetClasses() []ClassInfo {
	return w.classes
}

// GetVariables 获取所有变量信息
func (w *ASTWalker) GetVariables() []VariableInfo {
	return w.variables
}

// FindFunctionByName 根据名称查找函数
func (w *ASTWalker) FindFunctionByName(name string) *FunctionInfo {
	for _, fn := range w.functions {
		if fn.Name == name {
			return &fn
		}
	}
	return nil
}

// FindClassByName 根据名称查找类
func (w *ASTWalker) FindClassByName(name string) *ClassInfo {
	for _, cls := range w.classes {
		if cls.Name == name {
			return &cls
		}
	}
	return nil
}

// FindMethodByName 根据类名和方法名查找方法
func (w *ASTWalker) FindMethodByName(className, methodName string) *MethodInfo {
	for _, cls := range w.classes {
		if cls.Name == className {
			for _, method := range cls.Methods {
				if method.Name == methodName {
					return &method
				}
			}
		}
	}
	return nil
}

// Reset 重置遍历器状态
func (w *ASTWalker) Reset() {
	w.functions = make([]FunctionInfo, 0)
	w.classes = make([]ClassInfo, 0)
	w.variables = make([]VariableInfo, 0)
}
