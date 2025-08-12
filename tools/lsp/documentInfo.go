package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// LspContext LSP上下文结构体，扩展了标准Context
type LspContext struct {
	context.Context
	dataCtx data.Context

	// 基础数据
	values        map[string]interface{}
	scopeStack    []string
	functionStack []string
	loopStack     []LoopInfo
	localVars     map[string]data.Value
	traces        []string

	// 作用域管理
	parent    *LspContext
	scopeName string
}

// inferTypeFromExpression 从表达式推断类型
func inferTypeFromExpression(expr data.GetValue) data.Types {
	switch e := expr.(type) {
	case *node.NewExpression:
		// new 表达式返回类的实例类型
		return data.NewBaseType(e.ClassName)
	case *node.IntLiteral:
		return data.NewBaseType("int")
	case *node.FloatLiteral:
		return data.NewBaseType("float")
	case *node.StringLiteral:
		return data.NewBaseType("string")
	case *node.BooleanLiteral:
		return data.NewBaseType("bool")
	case *node.NullLiteral:
		return data.NewBaseType("null")
	case *node.Array:
		return data.NewBaseType("array")
	case *node.CallExpression:
		// 函数调用的返回类型需要更复杂的分析，暂时返回 mixed
		return data.NewBaseType("mixed")
	case data.Variable:
		// 如果是变量，返回其已知类型
		return e.GetType()
	default:
		return nil
	}
}

// LoopInfo 循环信息
type LoopInfo struct {
	Type string // 循环类型: for, foreach, while
	Id   string // 循环标识符
}

// NewLspContext 创建新的LSP上下文
func NewLspContext(baseCtx context.Context, dataCtx data.Context) *LspContext {
	return &LspContext{
		Context:       baseCtx,
		dataCtx:       dataCtx,
		values:        make(map[string]interface{}),
		scopeStack:    []string{"global"},
		functionStack: []string{},
		loopStack:     []LoopInfo{},
		localVars:     make(map[string]data.Value),
		traces:        []string{},
		parent:        nil,
		scopeName:     "global",
	}
}

// 实现data.Context接口
func (ctx *LspContext) SetNamespace(name string) data.Context {
	if ctx.dataCtx != nil {
		return ctx.dataCtx.SetNamespace(name)
	}
	return ctx
}

func (ctx *LspContext) GetNamespace() string {
	if ctx.dataCtx != nil {
		return ctx.dataCtx.GetNamespace()
	}
	return ""
}

func (ctx *LspContext) GetVariableValue(variable data.Variable) (data.Value, data.Control) {
	if ctx.dataCtx != nil {
		return ctx.dataCtx.GetVariableValue(variable)
	}
	return nil, nil
}

func (ctx *LspContext) GetIndexValue(index int) (data.Value, bool) {
	if ctx.dataCtx != nil {
		return ctx.dataCtx.GetIndexValue(index)
	}
	return nil, false
}

func (ctx *LspContext) SetVariableValue(variable data.Variable, value data.Value) data.Control {
	if ctx.dataCtx != nil {
		return ctx.dataCtx.SetVariableValue(variable, value)
	}
	return nil
}

func (ctx *LspContext) CreateContext(vars []data.Variable) data.Context {
	if ctx.dataCtx != nil {
		return ctx.dataCtx.CreateContext(vars)
	}
	return ctx
}

func (ctx *LspContext) CreateBaseContext() data.Context {
	if ctx.dataCtx != nil {
		return ctx.dataCtx.CreateBaseContext()
	}
	return ctx
}

func (ctx *LspContext) GetVM() data.VM {
	if ctx.dataCtx != nil {
		return ctx.dataCtx.GetVM()
	}
	return nil
}

// 实现LSP特有功能
func (ctx *LspContext) SetValue(key string, value interface{}) {
	ctx.values[key] = value
}

func (ctx *LspContext) GetValue(key string) (interface{}, bool) {
	value, exists := ctx.values[key]
	return value, exists
}

func (ctx *LspContext) PushScope(name string) {
	ctx.scopeStack = append(ctx.scopeStack, name)
}

func (ctx *LspContext) PopScope() string {
	if len(ctx.scopeStack) > 1 {
		popped := ctx.scopeStack[len(ctx.scopeStack)-1]
		ctx.scopeStack = ctx.scopeStack[:len(ctx.scopeStack)-1]
		return popped
	}
	return ""
}

func (ctx *LspContext) GetCurrentScope() string {
	if len(ctx.scopeStack) > 0 {
		return ctx.scopeStack[len(ctx.scopeStack)-1]
	}
	return "global"
}

func (ctx *LspContext) PushFunction(name string) {
	ctx.functionStack = append(ctx.functionStack, name)
}

func (ctx *LspContext) PopFunction() string {
	if len(ctx.functionStack) > 0 {
		popped := ctx.functionStack[len(ctx.functionStack)-1]
		ctx.functionStack = ctx.functionStack[:len(ctx.functionStack)-1]
		return popped
	}
	return ""
}

func (ctx *LspContext) GetCurrentFunction() string {
	if len(ctx.functionStack) > 0 {
		return ctx.functionStack[len(ctx.functionStack)-1]
	}
	return ""
}

func (ctx *LspContext) PushLoop(info LoopInfo) {
	ctx.loopStack = append(ctx.loopStack, info)
}

func (ctx *LspContext) PopLoop() *LoopInfo {
	if len(ctx.loopStack) > 0 {
		popped := ctx.loopStack[len(ctx.loopStack)-1]
		ctx.loopStack = ctx.loopStack[:len(ctx.loopStack)-1]
		return &popped
	}
	return nil
}

func (ctx *LspContext) GetCurrentLoop() *LoopInfo {
	if len(ctx.loopStack) > 0 {
		return &ctx.loopStack[len(ctx.loopStack)-1]
	}
	return nil
}

func (ctx *LspContext) SetLocalVar(name string, value data.Value) {
	ctx.localVars[name] = value
}

func (ctx *LspContext) GetLocalVar(name string) (data.Value, bool) {
	value, exists := ctx.localVars[name]
	return value, exists
}

// 变量类型管理
func (ctx *LspContext) SetVariableType(varName string, typ data.Types) {
	key := "var_type:" + varName
	ctx.values[key] = typ
}

func (ctx *LspContext) GetVariableType(varName string) data.Types {
	key := "var_type:" + varName

	// 在当前作用域查找
	if typ, exists := ctx.values[key]; exists {
		if dataType, ok := typ.(data.Types); ok {
			return dataType
		}
	}

	// 在父作用域查找
	if ctx.parent != nil {
		return ctx.parent.GetVariableType(varName)
	}

	return nil
}

// 作用域管理
func (ctx *LspContext) CreateChildScope(scopeName string) *LspContext {
	childCtx := &LspContext{
		Context:       ctx.Context,
		dataCtx:       ctx.dataCtx,
		values:        make(map[string]interface{}),
		scopeStack:    append([]string{}, ctx.scopeStack...),
		functionStack: append([]string{}, ctx.functionStack...),
		loopStack:     append([]LoopInfo{}, ctx.loopStack...),
		localVars:     make(map[string]data.Value),
		traces:        append([]string{}, ctx.traces...),
		parent:        ctx,
		scopeName:     scopeName,
	}

	// 继承父作用域的类型信息
	for key, value := range ctx.values {
		if strings.HasPrefix(key, "var_type:") {
			childCtx.values[key] = value
		}
	}

	return childCtx
}

func (ctx *LspContext) GetParentScope() *LspContext {
	return ctx.parent
}

// identifyVariableTypes 识别变量类型
func (d *DocumentInfo) identifyVariableTypes(ctx *LspContext, stmt data.GetValue) {
	switch n := stmt.(type) {
	case *node.BinaryAssignVariable:
		// 变量赋值：$a = new A()
		if leftVar, ok := n.Left.(*node.VariableExpression); ok {
			if inferredType := inferTypeFromExpression(n.Right); inferredType != nil {
				ctx.SetVariableType(leftVar.Name, inferredType)
			}
		}

	case *node.VarStatement:
		// var 声明：var $a = "string"
		if n.Initializer != nil {
			if inferredType := inferTypeFromExpression(n.Initializer); inferredType != nil {
				ctx.SetVariableType(n.Name, inferredType)
			}
		}

	case *node.ConstStatement:
		// const 声明：const $a = 123
		if n.Initializer != nil {
			if inferredType := inferTypeFromExpression(n.Initializer); inferredType != nil {
				ctx.SetVariableType(n.Name, inferredType)
			}
		}
	}
}

type DocumentInfo struct {
	Content string
	Version int32
	AST     *node.Program
	Parser  *LspParser
}

// CheckNode 遍历节点的回调函数
type CheckNode func(ctx *LspContext, parent, child data.GetValue) bool

// Foreach 穷举 AST 结构，遍历所有子节点传入check
func (d *DocumentInfo) Foreach(check CheckNode) {
	if d.AST == nil {
		return
	}

	// 创建LSP上下文
	baseCtx := context.Background()
	lspCtx := NewLspContext(baseCtx, nil)

	// 从根节点开始遍历
	for _, stmt := range d.AST.Statements {
		d.foreachNode(lspCtx, stmt, nil, check)
	}
}

// foreachNode 递归遍历单个节点及其所有子节点
func (d *DocumentInfo) foreachNode(ctx *LspContext, stmt data.GetValue, parent data.GetValue, check CheckNode) {
	// 先识别变量类型
	d.identifyVariableTypes(ctx, stmt)

	// 根据节点类型递归遍历子节点
	next := check(ctx, parent, stmt)
	if !next {
		return
	}
	parent = stmt
	switch n := stmt.(type) {
	case *node.FunctionStatement:
		// 函数定义：创建新的子作用域
		funcName := n.Name
		if funcName == "" {
			funcName = "anonymous"
		}

		// 创建函数作用域
		funcCtx := ctx.CreateChildScope("function:" + funcName)

		ctx.PushFunction(funcName)
		ctx.PushScope("function:" + funcName)

		// 在函数作用域中遍历参数和函数体
		for _, param := range n.Params {
			d.foreachNode(funcCtx, param, parent, check)
		}
		for _, bodyStmt := range n.Body {
			d.foreachNode(funcCtx, bodyStmt, parent, check)
		}

		ctx.PopScope()
		ctx.PopFunction()

	case *node.ClassStatement:
		// 类定义：创建新的子作用域
		className := n.Name
		classCtx := ctx.CreateChildScope("class:" + className)

		// 在类作用域中遍历属性和方法
		for _, property := range n.Properties {
			d.foreachNode(classCtx, property, parent, check)
			if defaultValue := property.GetDefaultValue(); defaultValue != nil {
				d.foreachNode(classCtx, defaultValue, property, check)
			}
		}
		// 注意：方法是data.Method接口，不实现data.GetValue，所以不能直接遍历
		// 如果需要遍历方法内容，需要通过其他方式获取方法的AST节点

	case *node.InterfaceStatement:
		// 接口定义：遍历方法
		// 注意：方法是data.Method接口，不实现data.GetValue，所以不能直接遍历

	case *node.IfStatement:
		// if语句：遍历条件、then分支、elseif分支、else分支
		d.foreachNode(ctx, n.Condition, parent, check)
		for _, stmt := range n.ThenBranch {
			d.foreachNode(ctx, stmt, parent, check)
		}
		for _, elseIf := range n.ElseIf {
			d.foreachNode(ctx, elseIf.Condition, parent, check)
			for _, stmt := range elseIf.ThenBranch {
				d.foreachNode(ctx, stmt, parent, check)
			}
		}
		for _, stmt := range n.ElseBranch {
			d.foreachNode(ctx, stmt, parent, check)
		}

	case *node.ForStatement:
		// for语句：遍历初始化、条件、增量、循环体
		if n.Initializer != nil {
			d.foreachNode(ctx, n.Initializer, parent, check)
		}
		if n.Condition != nil {
			d.foreachNode(ctx, n.Condition, parent, check)
		}
		if n.Increment != nil {
			d.foreachNode(ctx, n.Increment, parent, check)
		}

		for _, stmt := range n.Body {
			d.foreachNode(ctx, stmt, parent, check)
		}

	case *node.ForeachStatement:
		// foreach语句：遍历数组、键、值、循环体
		d.foreachNode(ctx, n.Array, parent, check)
		if n.Key != nil {
			d.foreachNode(ctx, n.Key, parent, check)
		}
		d.foreachNode(ctx, n.Value, parent, check)
		for _, stmt := range n.Body {
			d.foreachNode(ctx, stmt, parent, check)
		}

	case *node.WhileStatement:
		// while语句：遍历条件和循环体
		d.foreachNode(ctx, n.Condition, parent, check)

		for _, stmt := range n.Body {
			d.foreachNode(ctx, stmt, parent, check)
		}

	case *node.TryStatement:
		// try语句：遍历try块、catch块、finally块
		for _, stmt := range n.TryBlock {
			d.foreachNode(ctx, stmt, parent, check)
		}
		for _, catchBlock := range n.CatchBlocks {
			d.foreachNode(ctx, catchBlock.Variable, parent, check)
			for _, stmt := range catchBlock.Body {
				d.foreachNode(ctx, stmt, parent, check)
			}
		}
		if n.FinallyBlock != nil {
			for _, stmt := range n.FinallyBlock {
				d.foreachNode(ctx, stmt, parent, check)
			}
		}

	case *node.SwitchStatement:
		// switch语句：遍历条件和所有case
		d.foreachNode(ctx, n.Condition, parent, check)
		for _, caseStmt := range n.Cases {
			if caseStmt.CaseValue != nil {
				d.foreachNode(ctx, caseStmt.CaseValue, parent, check)
			}
			for _, stmt := range caseStmt.Statements {
				d.foreachNode(ctx, stmt, parent, check)
			}
		}
		// 遍历default分支
		for _, stmt := range n.DefaultCase {
			d.foreachNode(ctx, stmt, parent, check)
		}

	case *node.BlockStatement:
		// 块语句：遍历所有语句
		for _, stmt := range n.Statements {
			d.foreachNode(ctx, stmt, parent, check)
		}

	case *node.CallExpression:
		// 函数调用：遍历所有参数
		for _, arg := range n.Args {
			d.foreachNode(ctx, arg, parent, check)
		}

	case *node.NewExpression:
		// new表达式：遍历所有参数
		for _, arg := range n.Arguments {
			d.foreachNode(ctx, arg, parent, check)
		}

	case *node.Array:
		// 数组：遍历所有元素
		for _, element := range n.V {
			d.foreachNode(ctx, element, parent, check)
		}

	case *node.IndexExpression:
		// 索引表达式：遍历数组和索引
		d.foreachNode(ctx, n.Array, parent, check)
		if n.Index != nil {
			d.foreachNode(ctx, n.Index, parent, check)
		}

	case *node.TernaryExpression:
		// 三元表达式：遍历条件、真值、假值
		d.foreachNode(ctx, n.Condition, parent, check)
		d.foreachNode(ctx, n.TrueValue, parent, check)
		d.foreachNode(ctx, n.FalseValue, parent, check)

	case *node.LambdaExpression:
		// Lambda表达式：遍历参数和函数体

		for _, param := range n.Params {
			d.foreachNode(ctx, param, parent, check)
		}
		for _, stmt := range n.Body {
			d.foreachNode(ctx, stmt, parent, check)
		}

	case *node.ReturnStatement:
		// return语句：遍历返回值
		if n.Value != nil {
			d.foreachNode(ctx, n.Value, parent, check)
		}

	case *node.ThrowStatement:
		// throw语句：遍历抛出的值
		d.foreachNode(ctx, n.Value, parent, check)

	case *node.EchoStatement:
		// echo语句：遍历所有表达式
		for _, expr := range n.Expressions {
			d.foreachNode(ctx, expr, parent, check)
		}

	case *node.VarStatement:
		// var语句：遍历初始化值
		if n.Initializer != nil {
			d.foreachNode(ctx, n.Initializer, parent, check)

			// 简单的类型推断：从初始化值推断变量类型
			if inferredType := inferTypeFromExpression(n.Initializer); inferredType != nil {
				// 这里暂时无法直接修改 VarStatement 的类型，因为它没有 Type 字段
				// 类型推断主要在变量使用时生效
				// 但是可以将类型信息记录到上下文中
				ctx.SetVariableType(n.Name, inferredType)
			}
		}

	case *node.ConstStatement:
		// const语句：遍历初始化值
		if n.Initializer != nil {
			d.foreachNode(ctx, n.Initializer, parent, check)
		}

	// 二元表达式
	case *node.BinaryAdd:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinarySub:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryMul:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryQuo:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryRem:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryEq:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryNe:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryLt:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryLe:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryGt:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryGe:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryLand:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryLor:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryDot:
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryAssignVariable:
		// 变量赋值：遍历左右操作数
		// 类型推断：从右侧表达式推断左侧变量类型
		if leftVar, ok := n.Left.(*node.VariableExpression); ok {
			if inferredType := inferTypeFromExpression(n.Right); inferredType != nil {
				// 如果变量还没有类型，直接设置类型
				if leftVar.Type == nil {
					leftVar.Type = inferredType
				}
			}
		}

		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.BinaryAssignVariableList:
		// 变量列表赋值：遍历左右操作数
		// 类型推断：尝试从右侧表达式推断类型
		if inferredType := inferTypeFromExpression(n.Right); inferredType != nil {
			// 对于列表赋值，为每个变量分配相同类型
			leftVarList := n.Left // n.Left 已经是 *node.VariableList 类型
			for _, variable := range leftVarList.Vars {
				// variable 已经是 *node.VariableExpression 类型
				if variable.Type == nil {
					// 直接设置变量的类型
					variable.Type = inferredType
				}
			}
		}
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)
	case *node.InstanceOfExpression:
		// instanceof表达式：遍历对象表达式
		d.foreachNode(ctx, n.Object, parent, check)

	case *node.LikeExpression:
		// like表达式：遍历对象表达式
		d.foreachNode(ctx, n.Object, parent, check)

	case *node.UnaryExpression:
		// 一元表达式：遍历右操作数
		d.foreachNode(ctx, n.Right, parent, check)

	case *node.Kv:
		// 键值对：遍历所有键值
		for key, value := range n.V {
			d.foreachNode(ctx, key, parent, check)
			d.foreachNode(ctx, value, parent, check)
		}

	case *node.VariableList:
		// 变量列表：遍历所有变量
		for _, variable := range n.Vars {
			d.foreachNode(ctx, variable, parent, check)
		}

	case *node.Parameter:
		// 参数：遍历默认值
		if n.DefaultValue != nil {
			d.foreachNode(ctx, n.DefaultValue, parent, check)
		}

	case *node.Parameters:
		// 多值参数：遍历默认值
		if n.DefaultValue != nil {
			d.foreachNode(ctx, n.DefaultValue, parent, check)
		}
	case *node.SpawnStatement:
		// spawn语句：遍历调用表达式
		d.foreachNode(ctx, n.Call, parent, check)

	case *node.CallMethod:
		// 方法调用：遍历对象和参数
		if n.Method != nil {
			d.foreachNode(ctx, n.Method, parent, check)
		}
		for _, arg := range n.Args {
			d.foreachNode(ctx, arg, parent, check)
		}

	case *node.CallObjectMethod:
		// 对象方法调用：遍历对象和参数
		if n.Object != nil {
			d.foreachNode(ctx, n.Object, parent, check)
		}
		for _, arg := range n.Args {
			d.foreachNode(ctx, arg, parent, check)
		}

	case *node.CallObjectProperty:
		// 对象属性访问：遍历对象
		if n.Object != nil {
			d.foreachNode(ctx, n.Object, parent, check)
		}

	case *node.CallStaticMethod:
		// 静态方法调用：叶子节点（无参数）

	case *node.CallStaticProperty:
		// 静态属性访问：叶子节点

	case *node.CallParentMethod:
		// 父类方法调用：叶子节点（无参数）

	case *node.CallParentProperty:
		// 父类属性访问：叶子节点

	case *node.MatchStatement:
		// match语句：遍历条件和所有分支
		d.foreachNode(ctx, n.Condition, parent, check)
		for _, arm := range n.Arms {
			for _, condition := range arm.Conditions {
				d.foreachNode(ctx, condition, parent, check)
			}
			if arm.Expression != nil {
				d.foreachNode(ctx, arm.Expression, parent, check)
			}
			for _, stmt := range arm.Statements {
				d.foreachNode(ctx, stmt, parent, check)
			}
		}
		// 遍历默认分支
		for _, stmt := range n.Default {
			d.foreachNode(ctx, stmt, parent, check)
		}

	case *node.Range:
		// range表达式：遍历开始和结束值
		d.foreachNode(ctx, n.Start, parent, check)
		d.foreachNode(ctx, n.Stop, parent, check)

	case *node.NullCoalesceExpression:
		// null合并表达式：遍历左右操作数
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)

	case *node.PostfixIncr:
		// 后缀递增：遍历操作数
		d.foreachNode(ctx, n.Left, parent, check)

	case *node.PostfixDecr:
		// 后缀递减：遍历操作数
		d.foreachNode(ctx, n.Left, parent, check)

	case *node.UnaryIncr:
		// 前缀递增：遍历操作数
		d.foreachNode(ctx, n.Right, parent, check)

	case *node.UnaryDecr:
		// 前缀递减：遍历操作数
		d.foreachNode(ctx, n.Right, parent, check)

	case *node.BinaryEqStrict:
		// 严格相等：遍历左右操作数
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)

	case *node.BinaryNeStrict:
		// 严格不等：遍历左右操作数
		d.foreachNode(ctx, n.Left, parent, check)
		d.foreachNode(ctx, n.Right, parent, check)

	case *node.BooleanLiteral:
		// 布尔字面量：叶子节点

	case *node.Annotation:
		// 注解：遍历参数
		for _, arg := range n.Arguments {
			d.foreachNode(ctx, arg, parent, check)
		}

	case *node.ClassConstant:
	case *node.ClassProperty:
		// 类属性：遍历默认值
		if n.DefaultValue != nil {
			d.foreachNode(ctx, n.DefaultValue, parent, check)
		}

	case *node.ClassMethod:
		// 类方法：遍历参数和函数体
		for _, param := range n.Params {
			d.foreachNode(ctx, param, parent, check)
		}
		for _, stmt := range n.Body {
			d.foreachNode(ctx, stmt, parent, check)
		}
	case *node.NamedArgument:
		// 命名参数：遍历值
		d.foreachNode(ctx, n.Value, parent, check)

	case *node.HtmlNode:
		// HTML节点：遍历子节点
		for _, child := range n.Children {
			d.foreachNode(ctx, child, parent, check)
		}

	case *node.HtmlForNode:
		// HTML for节点：遍历数组和循环体
		d.foreachNode(ctx, n.Array, parent, check)
		d.foreachNode(ctx, n.Key, parent, check)
		d.foreachNode(ctx, n.Value, parent, check)
	case *node.Namespace:
		// 命名空间：遍历所有语句
		for _, stmt := range n.Statements {
			d.foreachNode(ctx, stmt, parent, check)
		}
	case *node.Todo:
		// Todo：叶子节点

	// 叶子节点，无需进一步遍历
	case *node.VariableExpression:
		// 变量表达式：叶子节点
	case *node.IntLiteral:
		// 整数字面量：叶子节点
	case *node.FloatLiteral:
		// 浮点数字面量：叶子节点
	case *node.StringLiteral:
		// 字符串字面量：叶子节点
	case *node.NullLiteral:
		// null字面量：叶子节点
	case *node.This:
		// this：叶子节点
	case *node.Parent:
		// parent：叶子节点
	case *node.BreakStatement:
		// break语句：叶子节点
	case *node.ContinueStatement:
		// continue语句：叶子节点
	case *node.UseStatement:
		// use语句：叶子节点

	// 默认情况：遇到未处理的节点类型时报错
	default:
		// 报错：发现未处理的节点类型
		panic(fmt.Sprintf("foreachNode: 未处理的节点类型 %T", stmt))
	}
}
