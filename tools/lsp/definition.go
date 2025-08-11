package main

import (
	"encoding/json"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/sourcegraph/jsonrpc2"
)

// 处理定义跳转请求
func handleTextDocumentDefinition(req *jsonrpc2.Request) (interface{}, error) {
	logLSPCommunication("textDocument/definition", true, req.Params)

	var params DefinitionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal definition params: %v", err)
	}

	uri := params.TextDocument.URI
	position := params.Position

	if *logLevel > 2 {
		fmt.Printf("[INFO] Definition requested for %s at %d:%d\n", uri, position.Line, position.Character)
	}

	doc, exists := documents[uri]
	if !exists {
		return nil, nil
	}

	// 使用 AST 遍历查找定义位置
	location := findDefinitionInAST(doc, position)
	if location == nil {
		return nil, nil
	}

	result := []Location{*location}
	logLSPResponse("textDocument/definition", result, nil)
	return result, nil
}

// findDefinitionInAST 在 AST 中查找定义位置
func findDefinitionInAST(doc *DocumentInfo, position Position) *Location {
	if doc.AST == nil {
		return nil
	}

	// 获取光标位置的节点
	node := getNodeAtPositionFromAST(doc, position)
	if node == nil {
		return nil
	}

	if *logLevel > 2 {
		fmt.Printf("[INFO] Found node at position: %T\n", node)
	}

	// 根据节点类型查找定义
	return findDefinitionFromNode(node)
}

// findDefinitionFromNode 根据节点类型查找定义位置
func findDefinitionFromNode(v data.GetValue) *Location {
	if v == nil {
		return nil
	}

	switch n := v.(type) {
	case *node.CallExpression:
		// 函数调用，查找函数定义
		return findFunctionDefinition(nil, n.FunName)
	case *node.NewExpression:
		// new 表达式，从类名查找类定义
		return findClassDefinition(nil, n.ClassName)
	case *node.CallMethod:
		// 方法调用，查找方法定义
		// 注意：CallMethod 的 Method 字段是 data.GetValue 类型，需要进一步处理
		if method, ok := n.Method.(interface{ GetName() string }); ok {
			return findMethodDefinition(nil, "", method.GetName())
		}
		return nil
	case *node.VariableExpression:
		// 变量引用，查找变量定义
		return findVariableDefinition(nil, n.Name)
	}

	// 对于其他类型的节点，暂时返回 nil
	return nil
}

// getNodeAtPositionFromAST 从 AST 中获取光标位置的节点
func getNodeAtPositionFromAST(doc *DocumentInfo, position Position) data.GetValue {
	if doc.AST == nil {
		return nil
	}

	// 尝试将 AST 转换为 *node.Program
	program, ok := doc.AST.(*node.Program)
	if !ok {
		// 如果转换失败，回退到文本解析
		return getNodeAtPositionFromText(doc.Content, position)
	}

	// 遍历 AST 查找光标位置的节点
	node := findNodeAtPosition(program, position)
	if node != nil {
		return node
	}

	// 如果 AST 中没有找到，回退到简单的文本解析
	return getNodeAtPositionFromText(doc.Content, position)
}

// findNodeAtPosition 在 AST 中查找指定位置的节点
func findNodeAtPosition(program *node.Program, position Position) data.GetValue {
	if program == nil {
		return nil
	}

	// 遍历所有语句，选择“最小包含范围”的节点
	var best data.GetValue
	for _, stmt := range program.Statements {
		candidate := findNodeInStatement(stmt, position)
		best = pickSmallerNode(best, candidate)
	}

	return best
}

// findNodeInStatement 在语句中查找节点
func findNodeInStatement(stmt node.Statement, position Position) data.GetValue {
	if stmt == nil {
		return nil
	}

	// 检查当前语句的位置是否包含光标位置（行内或整行命中均可）
	if !isPositionInRange(stmt, position) && !isPositionInLineRange(stmt, position) {
		return nil
	}

	// 作为兜底候选：如果没有更精确的子节点，返回当前语句本身
	var best data.GetValue = stmt

	switch s := stmt.(type) {
	case *node.FunctionStatement:
		// 合并参数和函数体成一个列表
		candidate := findSymbolInExpressions(position, append(s.GetParams(), s.GetBody()...)...)
		return pickSmallerNode(best, candidate)
	case *node.ClassStatement:
		candidate := findNodeInClass(s, position)
		return pickSmallerNode(best, candidate)
	case *node.InterfaceStatement:
		candidate := findNodeInInterface(s, position)
		return pickSmallerNode(best, candidate)
	case *node.Namespace:
		candidate := findSymbolInExpressions(position, s.Statements...)
		return pickSmallerNode(best, candidate)
	case *node.EchoStatement:
		candidate := findSymbolInExpressions(position, s.Expressions...)
		return pickSmallerNode(best, candidate)
	case *node.IfStatement:
		// 合并条件、then分支、else if分支和else分支
		allExprs := []data.GetValue{s.Condition}
		allExprs = append(allExprs, s.ThenBranch...)
		for _, elseIf := range s.ElseIf {
			allExprs = append(allExprs, elseIf.Condition)
			allExprs = append(allExprs, elseIf.ThenBranch...)
		}
		allExprs = append(allExprs, s.ElseBranch...)
		candidate := findSymbolInExpressions(position, allExprs...)
		return pickSmallerNode(best, candidate)
	case *node.ForStatement:
		// 合并初始化器、条件、增量和循环体
		allExprs := []data.GetValue{}
		if s.Initializer != nil {
			allExprs = append(allExprs, s.Initializer)
		}
		if s.Condition != nil {
			allExprs = append(allExprs, s.Condition)
		}
		if s.Increment != nil {
			allExprs = append(allExprs, s.Increment)
		}
		allExprs = append(allExprs, s.Body...)
		candidate := findSymbolInExpressions(position, allExprs...)
		return pickSmallerNode(best, candidate)
	case *node.WhileStatement:
		// 合并条件和循环体
		allExprs := []data.GetValue{}
		if s.Condition != nil {
			allExprs = append(allExprs, s.Condition)
		}
		allExprs = append(allExprs, s.Body...)
		candidate := findSymbolInExpressions(position, allExprs...)
		return pickSmallerNode(best, candidate)
	case *node.ForeachStatement:
		// 合并数组、键变量、值变量和循环体
		allExprs := []data.GetValue{s.Array}
		if s.Key != nil {
			allExprs = append(allExprs, s.Key)
		}
		if s.Value != nil {
			allExprs = append(allExprs, s.Value)
		}
		allExprs = append(allExprs, s.Body...)
		candidate := findSymbolInExpressions(position, allExprs...)
		return pickSmallerNode(best, candidate)
	case *node.SwitchStatement:
		// 合并条件、case分支和default分支
		allExprs := []data.GetValue{s.Condition}
		for _, caseStmt := range s.Cases {
			allExprs = append(allExprs, caseStmt.CaseValue)
			allExprs = append(allExprs, caseStmt.Statements...)
		}
		allExprs = append(allExprs, s.DefaultCase...)
		candidate := findSymbolInExpressions(position, allExprs...)
		return pickSmallerNode(best, candidate)
	case *node.TryStatement:
		// 合并try块、catch块和finally块
		allExprs := []data.GetValue{}
		allExprs = append(allExprs, s.TryBlock...)
		for _, catchBlock := range s.CatchBlocks {
			if catchBlock.Variable != nil {
				allExprs = append(allExprs, catchBlock.Variable)
			}
			allExprs = append(allExprs, catchBlock.Body...)
		}
		allExprs = append(allExprs, s.FinallyBlock...)
		candidate := findSymbolInExpressions(position, allExprs...)
		return pickSmallerNode(best, candidate)
	case *node.BlockStatement:
		// 块语句包含多个子语句，需要转换为 data.GetValue 类型
		allExprs := make([]data.GetValue, len(s.Statements))
		for i, stmt := range s.Statements {
			allExprs[i] = stmt
		}
		candidate := findSymbolInExpressions(position, allExprs...)
		return pickSmallerNode(best, candidate)
	case *node.ReturnStatement:
		// 返回语句包含返回值表达式
		if s.Value != nil {
			candidate := findSymbolInExpressions(position, s.Value)
			return pickSmallerNode(best, candidate)
		}
		return nil
	case *node.BinaryAssignVariable:
		// $a = expr; 赋值，优先向右侧深入（例如 new A()）
		allExprs := []data.GetValue{}
		if s.Left != nil {
			allExprs = append(allExprs, s.Left)
		}
		if s.Right != nil {
			allExprs = append(allExprs, s.Right)
		}
		candidate := findSymbolInExpressions(position, allExprs...)
		return pickSmallerNode(best, candidate)
	case *node.BinaryAssignVariableList:
		// 多变量赋值，右侧或左侧
		allExprs := []data.GetValue{}
		if s.Right != nil {
			allExprs = append(allExprs, s.Right)
		}
		if s.Left != nil {
			allExprs = append(allExprs, s.Left)
		}
		candidate := findSymbolInExpressions(position, allExprs...)
		return pickSmallerNode(best, candidate)
	case *node.BinaryAssign:
		// 赋值表达式，优先深入左右子表达式
		allExprs := []data.GetValue{}
		if s.Left != nil {
			allExprs = append(allExprs, s.Left)
		}
		if s.Right != nil {
			allExprs = append(allExprs, s.Right)
		}
		candidate := findSymbolInExpressions(position, allExprs...)
		return pickSmallerNode(best, candidate)
	case *node.VarStatement:
		// 变量声明语句包含名称和初始化器
		allExprs := []data.GetValue{}
		if s.Initializer != nil {
			allExprs = append(allExprs, s.Initializer)
		}
		candidate := findSymbolInExpressions(position, allExprs...)
		return pickSmallerNode(best, candidate)
	case *node.ConstStatement:
		// 常量声明语句包含名称和初始化器
		allExprs := []data.GetValue{}
		if s.Initializer != nil {
			allExprs = append(allExprs, s.Initializer)
		}
		candidate := findSymbolInExpressions(position, allExprs...)
		return pickSmallerNode(best, candidate)
	case *node.ThrowStatement:
		// 抛出语句包含异常表达式
		if s.Value != nil {
			candidate := findSymbolInExpressions(position, s.Value)
			return pickSmallerNode(best, candidate)
		}
		return nil
	case *node.BreakStatement:
		// break语句没有子节点
		return nil
	case *node.ContinueStatement:
		// continue语句没有子节点
		return nil
	case *node.SpawnStatement:
		// spawn语句包含要执行的方法调用
		if s.Call != nil {
			candidate := findSymbolInExpressions(position, s.Call)
			return pickSmallerNode(best, candidate)
		}
		return nil
	case *node.UseStatement:
		// use语句没有子节点，只有命名空间和别名
		return nil
	default:
		// 对于其他类型的语句，尝试作为表达式处理
		candidate := findSymbolInExpression(stmt, position)
		return pickSmallerNode(best, candidate)
	}
}

func findSymbolInExpressions(position Position, exprs ...data.GetValue) data.GetValue {
	var best data.GetValue
	fmt.Printf("[DEBUG] findSymbolInExpressions: position=(%d,%d), exprs count=%d\n",
		position.Line, position.Character, len(exprs))

	// 循环处理每个表达式参数
	for i, expr := range exprs {
		if expr == nil {
			continue
		}

		fmt.Printf("[DEBUG] Processing expr[%d]: %T\n", i, expr)

		// 使用精确的位置检查，确保表达式真正包含光标位置
		if stmt, ok := expr.(node.Statement); ok {
			if !isPositionInRange(stmt, position) {
				fmt.Printf("[DEBUG] Expr[%d] position check failed, skipping\n", i)
				continue
			}
			fmt.Printf("[DEBUG] Expr[%d] position check passed\n", i)
		}

		var candidate data.GetValue
		switch expr.(type) {
		// 原子字面量和变量：直接作为候选
		case *node.VariableExpression, *node.StringLiteral, *node.BooleanLiteral, *node.NullLiteral, *node.IntLiteral, *node.FloatLiteral:
			candidate = expr
			fmt.Printf("[DEBUG] Expr[%d] is atomic, candidate=%T\n", i, candidate)
		default:
			if stmt, ok := expr.(node.Statement); ok {
				fmt.Printf("[DEBUG] Expr[%d] calling findSymbolInExpression\n", i)
				candidate = findSymbolInExpression(stmt, position)
				fmt.Printf("[DEBUG] Expr[%d] findSymbolInExpression returned: %T\n", i, candidate)
			} else {
				candidate = expr
				fmt.Printf("[DEBUG] Expr[%d] is not statement, candidate=%T\n", i, candidate)
			}
		}

		if candidate != nil {
			fmt.Printf("[DEBUG] Expr[%d] has candidate, calling pickSmallerNode\n", i)
			best = pickSmallerNode(best, candidate)
			fmt.Printf("[DEBUG] After pickSmallerNode, best=%T\n", best)
		}
	}

	fmt.Printf("[DEBUG] findSymbolInExpressions final result: %T\n", best)
	return best
}

// pickSmallerNode 返回最合适的节点；优先选择包含光标位置的节点
func pickSmallerNode(a, b data.GetValue) data.GetValue {
	fmt.Printf("[DEBUG] pickSmallerNode: a=%T, b=%T\n", a, b)

	if b == nil {
		fmt.Printf("[DEBUG] pickSmallerNode: b is nil, returning a\n")
		return a
	}
	if a == nil {
		fmt.Printf("[DEBUG] pickSmallerNode: a is nil, returning b\n")
		return b
	}

	af := getFromOf(a)
	bf := getFromOf(b)
	if af == nil {
		fmt.Printf("[DEBUG] pickSmallerNode: af is nil, returning b\n")
		return b
	}
	if bf == nil {
		fmt.Printf("[DEBUG] pickSmallerNode: bf is nil, returning a\n")
		return a
	}

	// 获取两个节点的范围
	slA, scA, elA, ecA := af.GetRange()
	slB, scB, elB, ecB := bf.GetRange()

	fmt.Printf("[DEBUG] pickSmallerNode: a range=(%d,%d,%d,%d), b range=(%d,%d,%d,%d)\n",
		slA, scA, elA, ecA, slB, scB, elB, ecB)

	// 计算两个节点的范围大小（字符数）
	rangeA := (elA-slA+1)*1000 + (ecA - scA + 1)
	rangeB := (elB-slB+1)*1000 + (ecB - scB + 1)

	fmt.Printf("[DEBUG] pickSmallerNode: rangeA=%d, rangeB=%d\n", rangeA, rangeB)

	// 优先选择范围更小的节点（更精确），但前提是它们都包含光标位置
	// 如果范围差异不大，选择范围更小的；如果差异很大，选择更合适的
	if rangeB < rangeA && (rangeA-rangeB) < 100 {
		fmt.Printf("[DEBUG] pickSmallerNode: rangeB < rangeA and difference small, returning b\n")
		return b
	}
	if rangeA < rangeB && (rangeB-rangeA) < 100 {
		fmt.Printf("[DEBUG] pickSmallerNode: rangeA < rangeB and difference small, returning a\n")
		return a
	}

	// 如果范围差异很大，选择范围更小的（更精确）
	if rangeB < rangeA {
		fmt.Printf("[DEBUG] pickSmallerNode: rangeB significantly smaller, returning b\n")
		return b
	}
	fmt.Printf("[DEBUG] pickSmallerNode: rangeA <= rangeB, returning a\n")
	return a
}

// rangeWeight 估算范围大小，用于比较"更小的包含区间"
func rangeWeight(from data.From) int64 {
	sl, sc, el, ec := from.GetRange()
	lineSpan := el - sl
	charSpan := ec - sc
	if lineSpan < 0 {
		lineSpan = 0
	}
	if charSpan < 0 {
		charSpan = 0
	}
	// 以行跨度为主、列跨度为辅进行排序权重估算
	return int64(lineSpan)*100000 + int64(charSpan)
}

// getFromOf 提取节点的 From 信息
func getFromOf(v data.GetValue) data.From {
	if v == nil {
		return nil
	}
	if gf, ok := v.(node.GetFrom); ok {
		return gf.GetFrom()
	}
	if st, ok := v.(node.Statement); ok {
		if gf, ok2 := st.(node.GetFrom); ok2 {
			return gf.GetFrom()
		}
	}
	return nil
}

// isPositionInRange 检查位置是否在节点范围内
func isPositionInRange(stmt node.Statement, position Position) bool {
	// 尝试获取节点的 From 信息
	var from data.From

	// 检查不同类型的节点
	if getFrom, ok := stmt.(node.GetFrom); ok {
		from = getFrom.GetFrom()
	}

	if from == nil {
		return false
	}

	// 直接使用 GetRange 获取位置范围
	startLine, startChar, endLine, endChar := from.GetRange()

	// 注意：LSP 使用从 0 开始的行号，我们的 from 系统使用从 1 开始的行号
	lspLine := int(position.Line) + 1

	// 添加调试信息
	fmt.Printf("[DEBUG] isPositionInRange: node=%T, position=(%d,%d), range=(%d,%d,%d,%d)\n",
		stmt, lspLine, position.Character, startLine, startChar, endLine, endChar)

	// 检查行号是否在范围内
	if lspLine < startLine || lspLine > endLine {
		fmt.Printf("[DEBUG] isPositionInRange: line out of range\n")
		return false
	}

	// 如果在起始行，检查字符位置是否在起始字符之后
	if lspLine == startLine {
		result := int(position.Character) >= startChar
		if *logLevel > 2 {
			fmt.Printf("[DEBUG] isPositionInRange: start line, char check: %d >= %d = %v\n",
				position.Character, startChar, result)
		}
		// 还需要检查是否在结束字符之前
		if result {
			result = int(position.Character) <= endChar
			if *logLevel > 2 {
				fmt.Printf("[DEBUG] isPositionInRange: start line, end char check: %d <= %d = %v\n",
					position.Character, endChar, result)
			}
		}
		return result
	}

	// 如果在结束行，检查字符位置是否在结束字符之前
	if lspLine == endLine {
		result := int(position.Character) <= endChar
		if *logLevel > 2 {
			fmt.Printf("[DEBUG] isPositionInRange: end line, char check: %d <= %d = %v\n",
				position.Character, endChar, result)
		}
		return result
	}

	// 如果在中间行，肯定在范围内
	if *logLevel > 2 {
		fmt.Printf("[DEBUG] isPositionInRange: middle line, in range\n")
	}
	return true
}

// findNodeInFunction 在函数中查找节点
func findNodeInFunction(fn *node.FunctionStatement, position Position) data.GetValue {
	if fn == nil {
		return nil
	}

	// 检查是否是函数名
	if isPositionInNodeName(fn, position) {
		return fn
	}

	// 检查函数体中的节点
	for _, bodyStmt := range fn.GetBody() {
		if nodeStmt, ok := bodyStmt.(node.Statement); ok {
			if node := findNodeInStatement(nodeStmt, position); node != nil {
				return node
			}
		}
	}

	return nil
}

// findNodeInClass 在类中查找节点
func findNodeInClass(cls *node.ClassStatement, position Position) data.GetValue {
	if cls == nil {
		return nil
	}

	// 检查是否是类名
	if isPositionInNodeName(cls, position) {
		return cls
	}

	// 检查类的属性
	properties := cls.GetProperties()
	for _, property := range properties {
		// 检查属性名
		if isPositionInNodeName(property, position) {
			return property
		}
		// 检查属性的默认值
		if defaultValue := property.GetDefaultValue(); defaultValue != nil {
			if node := findSymbolInExpressions(position, defaultValue); node != nil {
				return node
			}
		}
	}

	// 检查类的方法
	methods := cls.GetMethods()
	for _, method := range methods {
		// 检查方法体
		if methodStmt, ok := method.(node.Statement); ok {
			if node := findNodeInStatement(methodStmt, position); node != nil {
				return node
			}
		}
	}

	// 检查构造函数
	if construct := cls.GetConstruct(); construct != nil {
		// 检查构造函数体
		if constructStmt, ok := construct.(node.Statement); ok {
			if node := findNodeInStatement(constructStmt, position); node != nil {
				return node
			}
		}
	}

	// 注解检查暂时跳过，因为 Annotations 是私有字段
	// 如果需要检查注解，需要添加公共访问方法

	return nil
}

// findNodeInInterface 在接口中查找节点
func findNodeInInterface(iface *node.InterfaceStatement, position Position) data.GetValue {
	if iface == nil {
		return nil
	}

	// 检查是否是接口名
	if isPositionInNodeName(iface, position) {
		return iface
	}

	return nil
}

// findSymbolInExpression 在表达式中查找符号，直接返回节点
func findSymbolInExpression(expr node.Statement, position Position) data.GetValue {
	if expr == nil {
		return nil
	}

	// 检查位置是否在当前表达式范围内
	if !isPositionInRange(expr, position) {
		return nil
	}

	// 根据表达式类型进行具体处理
	switch e := expr.(type) {
	case *node.CallExpression:
		return findSymbolInCallExpression(e, position)
	case *node.CallMethod:
		return findSymbolInCallMethod(e, position)
	case *node.CallObjectMethod:
		return findSymbolInObjectMethod(e, position)
	case *node.CallObjectProperty:
		return findSymbolInObjectProperty(e, position)
	case *node.CallStaticMethod:
		return findSymbolInStaticMethod(e, position)
	case *node.NewExpression:
		return findSymbolInNewExpression(e, position)
	case *node.VariableExpression:
		return findSymbolInVariableExpression(e, position)
	default:
		// 回退到通用语句处理，继续向下递归
		return findNodeInStatement(expr, position)
	}
}

// findSymbolInCallExpression 在函数调用表达式中查找符号
func findSymbolInCallExpression(call *node.CallExpression, position Position) data.GetValue {
	if call == nil {
		return nil
	}

	if !isPositionInLineRange(call, position) {
		return nil
	}

	// 先在参数中寻找更精确的命中
	var best data.GetValue
	for _, arg := range call.Args {
		if arg == nil {
			continue
		}
		if st, ok := arg.(node.Statement); ok {
			if !isPositionInRange(st, position) {
				continue
			}
			candidate := findSymbolInExpression(st, position)
			best = pickSmallerNode(best, candidate)
		} else {
			best = pickSmallerNode(best, arg)
		}
	}

	// 若没有更小的子节点，则返回调用本身
	if best == nil {
		return call
	}
	return best
}

// findSymbolInCallMethod 在方法调用表达式中查找符号
func findSymbolInCallMethod(call *node.CallMethod, position Position) data.GetValue {
	if call == nil {
		return nil
	}

	if !isPositionInLineRange(call, position) {
		return nil
	}

	var best data.GetValue

	// 先检查方法表达式本身
	if call.Method != nil {
		if st, ok := call.Method.(node.Statement); ok {
			if isPositionInRange(st, position) {
				cand := findSymbolInExpression(st, position)
				best = pickSmallerNode(best, cand)
			}
		} else {
			best = pickSmallerNode(best, call.Method)
		}
	}

	// 再检查参数
	for _, arg := range call.Args {
		if arg == nil {
			continue
		}
		if st, ok := arg.(node.Statement); ok {
			if !isPositionInRange(st, position) {
				continue
			}
			cand := findSymbolInExpression(st, position)
			best = pickSmallerNode(best, cand)
		} else {
			best = pickSmallerNode(best, arg)
		}
	}

	if best == nil {
		return call
	}
	return best
}

// findSymbolInNewExpression 在 new 表达式中查找符号
func findSymbolInNewExpression(newExpr *node.NewExpression, position Position) data.GetValue {
	if newExpr == nil {
		return nil
	}

	if !isPositionInLineRange(newExpr, position) {
		return nil
	}

	// 如果光标在 new 表达式的范围内，直接返回 newExpr
	// 因为类名是 new 表达式的核心部分
	if isPositionInRange(newExpr, position) {
		return newExpr
	}

	var best data.GetValue
	for _, arg := range newExpr.Arguments {
		if arg == nil {
			continue
		}
		if st, ok := arg.(node.Statement); ok {
			if !isPositionInRange(st, position) {
				continue
			}
			cand := findSymbolInExpression(st, position)
			best = pickSmallerNode(best, cand)
		} else {
			best = pickSmallerNode(best, arg)
		}
	}

	if best == nil {
		return newExpr
	}
	return best
}

// findSymbolInVariableExpression 在变量表达式中查找符号
func findSymbolInVariableExpression(varExpr *node.VariableExpression, position Position) data.GetValue {
	if varExpr == nil {
		return nil
	}

	return varExpr
}

// 额外补充：对象方法、对象属性、静态方法
func findSymbolInObjectMethod(call *node.CallObjectMethod, position Position) data.GetValue {
	if call == nil || !isPositionInLineRange(call, position) {
		return nil
	}

	var best data.GetValue

	if call.Object != nil {
		if st, ok := call.Object.(node.Statement); ok {
			if isPositionInRange(st, position) {
				cand := findSymbolInExpression(st, position)
				best = pickSmallerNode(best, cand)
			}
		} else {
			best = pickSmallerNode(best, call.Object)
		}
	}
	for _, arg := range call.Args {
		if arg == nil {
			continue
		}
		if st, ok := arg.(node.Statement); ok {
			if !isPositionInRange(st, position) {
				continue
			}
			cand := findSymbolInExpression(st, position)
			best = pickSmallerNode(best, cand)
		} else {
			best = pickSmallerNode(best, arg)
		}
	}
	if best == nil {
		return call
	}
	return best
}

func findSymbolInObjectProperty(call *node.CallObjectProperty, position Position) data.GetValue {
	if call == nil || !isPositionInLineRange(call, position) {
		return nil
	}
	var best data.GetValue
	if call.Object != nil {
		if st, ok := call.Object.(node.Statement); ok {
			if isPositionInRange(st, position) {
				cand := findSymbolInExpression(st, position)
				best = pickSmallerNode(best, cand)
			}
		} else {
			best = pickSmallerNode(best, call.Object)
		}
	}
	if best == nil {
		return call
	}
	return best
}

func findSymbolInStaticMethod(call *node.CallStaticMethod, position Position) data.GetValue {
	if call == nil || !isPositionInLineRange(call, position) {
		return nil
	}
	return call
}

// 仅检查“行”是否命中，忽略列，避免 token.Pos 为末尾列导致的误差
func isPositionInLineRange(stmt node.Statement, position Position) bool {
	var from data.From
	if getFrom, ok := stmt.(node.GetFrom); ok {
		from = getFrom.GetFrom()
	}
	if from == nil {
		return false
	}
	startLine, _, endLine, _ := from.GetRange()
	lspLine := int(position.Line) + 1
	return lspLine >= startLine && lspLine <= endLine
}

// isPositionInNodeName 检查位置是否在节点名称范围内
func isPositionInNodeName(stmt node.Statement, position Position) bool {
	// 尝试获取节点的 From 信息
	var from data.From

	// 检查不同类型的节点
	switch n := stmt.(type) {
	case *node.FunctionStatement:
		from = n.GetFrom()
	case *node.ClassStatement:
		from = n.GetFrom()
	case *node.InterfaceStatement:
		from = n.GetFrom()
	case *node.Namespace:
		from = n.GetFrom()
	default:
		return false
	}

	if from == nil {
		return false
	}

	// 直接使用 GetRange 获取位置范围
	startLine, startChar, endLine, endChar := from.GetRange()

	// 注意：LSP 使用从 0 开始的行号，我们的 from 系统使用从 1 开始的行号
	lspLine := int(position.Line)

	// 检查是否在名称范围内（假设名称在开始位置附近）
	if lspLine+1 == startLine {
		// 如果在同一行，检查字符位置
		return int(position.Character) >= startChar && int(position.Character) <= endChar
	}

	// 如果跨行，检查是否在范围内
	if lspLine+1 >= startLine && lspLine+1 <= endLine {
		if lspLine+1 == startLine {
			return int(position.Character) >= startChar
		}
		if lspLine+1 == endLine {
			return int(position.Character) <= endChar
		}
		return true
	}

	return false
}

// getNodeAtPositionFromText 从文本中获取节点信息（回退方法）
func getNodeAtPositionFromText(content string, position Position) data.GetValue {
	// 这是一个简化的实现，实际应该返回一个表示文本位置的节点
	// 由于我们主要使用 AST，这个函数暂时返回 nil
	return nil
}

// findFunctionDefinition 查找函数定义
func findFunctionDefinition(doc *DocumentInfo, funcName string) *Location {
	if globalLspVM == nil {
		return nil
	}

	if function, exists := globalLspVM.GetFunc(funcName); exists {
		return createLocationFromFunction(function)
	}

	return nil
}

// findClassDefinition 查找类定义
func findClassDefinition(doc *DocumentInfo, className string) *Location {
	if globalLspVM == nil {
		return nil
	}

	if class, exists := globalLspVM.GetClass(className); exists {
		return createLocationFromClass(class)
	}

	return nil
}

// findVariableDefinition 查找变量定义
func findVariableDefinition(doc *DocumentInfo, varName string) *Location {
	// 变量定义查找功能暂时简化，返回 nil
	return nil
}

// findMethodDefinition 查找方法定义
func findMethodDefinition(doc *DocumentInfo, variableName, methodName string) *Location {
	if globalLspVM == nil {
		return nil
	}

	// 查找变量对应的类
	if class, exists := globalLspVM.GetClass(variableName); exists {
		if methodLocation := findMethodInClass(class, methodName); methodLocation != nil {
			return methodLocation
		}
	}

	// 如果找不到类，尝试查找同名函数
	if function, exists := globalLspVM.GetFunc(methodName); exists {
		return createLocationFromFunction(function)
	}

	return nil
}

// createLocationFromFunction 从函数定义创建位置信息
func createLocationFromFunction(function data.FuncStmt) *Location {
	// 使用 node.GetFrom 接口获取位置信息
	if fnStmt, ok := function.(node.GetFrom); ok {
		if from := fnStmt.GetFrom(); from != nil {
			startLine, startChar, endLine, endChar := from.ToLSPPosition()
			return &Location{
				URI: fmt.Sprintf("file://%s", from.GetSource()),
				Range: Range{
					Start: Position{Line: uint32(startLine), Character: uint32(startChar)},
					End:   Position{Line: uint32(endLine), Character: uint32(endChar)},
				},
			}
		}
	}
	return nil
}

// createLocationFromClass 从类定义创建位置信息
func createLocationFromClass(class data.ClassStmt) *Location {
	// 从类节点获取位置信息
	if from := class.GetFrom(); from != nil {
		startLine, startChar, endLine, endChar := from.ToLSPPosition()
		return &Location{
			URI: fmt.Sprintf("file://%s", from.GetSource()),
			Range: Range{
				Start: Position{Line: uint32(startLine), Character: uint32(startChar)},
				End:   Position{Line: uint32(endLine), Character: uint32(endChar)},
			},
		}
	}
	return nil
}

// findMethodInClass 在类中查找方法
func findMethodInClass(class data.ClassStmt, methodName string) *Location {
	// 尝试在类中查找指定名称的方法
	if method, exists := class.GetMethod(methodName); exists {
		// 如果方法有位置信息，返回方法的位置
		if methodWithFrom, ok := method.(node.GetFrom); ok {
			if from := methodWithFrom.GetFrom(); from != nil {
				startLine, startChar, endLine, endChar := from.ToLSPPosition()
				return &Location{
					URI: fmt.Sprintf("file://%s", from.GetSource()),
					Range: Range{
						Start: Position{Line: uint32(startLine), Character: uint32(startChar)},
						End:   Position{Line: uint32(endLine), Character: uint32(endChar)},
					},
				}
			}
		}
	}

	// 如果找不到方法或方法没有位置信息，返回类定义的位置
	if from := class.GetFrom(); from != nil {
		startLine, startChar, endLine, endChar := from.ToLSPPosition()
		return &Location{
			URI: fmt.Sprintf("file://%s", from.GetSource()),
			Range: Range{
				Start: Position{Line: uint32(startLine), Character: uint32(startChar)},
				End:   Position{Line: uint32(endLine), Character: uint32(endChar)},
			},
		}
	}
	return nil
}
