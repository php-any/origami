package main

import (
	"encoding/json"
	"fmt"
	"strings"

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

	// 尝试将 AST 转换为 *node.Program
	program, ok := doc.AST.(*node.Program)
	if !ok {
		// 如果转换失败，回退到文本解析
		return findDefinitionFromText(doc.Content, position)
	}

	// 遍历 AST 查找光标位置的符号
	symbolNode := findSymbolAtPosition(program, position)
	if symbolNode == nil {
		return nil
	}

	// 根据节点类型查找定义
	return findDefinitionFromNode(symbolNode, doc)
}

// findSymbolAtPosition 在 AST 中查找指定位置的符号节点
func findSymbolAtPosition(program *node.Program, position Position) node.Statement {
	if program == nil {
		return nil
	}

	// 遍历所有语句
	for _, stmt := range program.Statements {
		if symbolNode := findSymbolInStatement(stmt, position); symbolNode != nil {
			return symbolNode
		}
	}

	return nil
}

// findSymbolInStatement 在语句中查找符号
func findSymbolInStatement(stmt node.Statement, position Position) node.Statement {
	if stmt == nil {
		return nil
	}

	// 检查当前语句的位置是否包含光标位置
	if !isPositionInNode(stmt, position) {
		return nil
	}

	switch s := stmt.(type) {
	case *node.FunctionStatement:
		return findSymbolInFunction(s, position)
	case *node.ClassStatement:
		return findSymbolInClass(s, position)
	case *node.InterfaceStatement:
		return findSymbolInInterface(s, position)
	case *node.Namespace:
		return findSymbolInNamespace(s, position)
	case *node.CallExpression:
		return findSymbolInCallExpression(s, position)
	case *node.CallMethod:
		return findSymbolInCallMethod(s, position)
	case *node.NewExpression:
		return findSymbolInNewExpression(s, position)
	case *node.VariableExpression:
		return findSymbolInVariable(s, position)
	default:
		// 对于其他类型的语句，递归查找
		return findSymbolInExpression(stmt, position)
	}
}

// isPositionInNode 检查位置是否在节点范围内
func isPositionInNode(stmt node.Statement, position Position) bool {
	from := getNodeFrom(stmt)
	if from == nil {
		return false
	}

	start, end := from.GetPosition()

	// 使用更精确的位置计算
	startLine, startChar := calculatePositionFromOffset(start)
	endLine, endChar := calculatePositionFromOffset(end)

	// 检查行号是否在范围内
	if int(position.Line) < startLine || int(position.Line) > endLine {
		return false
	}

	// 如果在同一行，检查字符位置
	if int(position.Line) == startLine && int(position.Line) == endLine {
		return int(position.Character) >= startChar && int(position.Character) <= endChar
	}

	return true
}

// getNodeFrom 获取节点的 From 信息
func getNodeFrom(stmt node.Statement) data.From {
	switch n := stmt.(type) {
	case *node.FunctionStatement:
		return n.GetFrom()
	case *node.ClassStatement:
		return n.GetFrom()
	case *node.InterfaceStatement:
		return n.GetFrom()
	case *node.Namespace:
		return n.GetFrom()
	case *node.CallExpression:
		return n.GetFrom()
	case *node.CallMethod:
		return n.GetFrom()
	case *node.NewExpression:
		return n.GetFrom()
	case *node.VariableExpression:
		return n.GetFrom()
	default:
		// 尝试通过反射获取
		if nodeWithFrom, ok := stmt.(interface{ GetFrom() data.From }); ok {
			return nodeWithFrom.GetFrom()
		}
		return nil
	}
}

// calculatePositionFromOffset 根据偏移量计算行号和字符位置
func calculatePositionFromOffset(offset int) (line, char int) {
	// 这是一个简化的实现，实际应该根据源代码内容计算
	// 假设每行平均80个字符
	line = offset/80 + 1
	if line < 1 {
		line = 1
	}
	char = offset % 80
	if char < 0 {
		char = 0
	}
	return line, char
}

// findSymbolInFunction 在函数中查找符号
func findSymbolInFunction(fn *node.FunctionStatement, position Position) node.Statement {
	if fn == nil {
		return nil
	}

	// 检查是否是函数名
	if isPositionInNodeName(fn, position) {
		return fn
	}

	// 检查函数体中的符号
	for _, bodyStmt := range fn.GetBody() {
		if nodeStmt, ok := bodyStmt.(node.Statement); ok {
			if symbolNode := findSymbolInStatement(nodeStmt, position); symbolNode != nil {
				return symbolNode
			}
		}
	}

	return nil
}

// findSymbolInClass 在类中查找符号
func findSymbolInClass(cls *node.ClassStatement, position Position) node.Statement {
	if cls == nil {
		return nil
	}

	// 检查是否是类名
	if isPositionInNodeName(cls, position) {
		return cls
	}

	// 检查类的方法
	methods := cls.GetMethods()
	for _, method := range methods {
		if methodStmt, ok := method.(node.Statement); ok {
			if symbolNode := findSymbolInStatement(methodStmt, position); symbolNode != nil {
				return symbolNode
			}
		}
	}

	return nil
}

// findSymbolInInterface 在接口中查找符号
func findSymbolInInterface(iface *node.InterfaceStatement, position Position) node.Statement {
	if iface == nil {
		return nil
	}

	// 检查是否是接口名
	if isPositionInNodeName(iface, position) {
		return iface
	}

	return nil
}

// findSymbolInNamespace 在命名空间中查找符号
func findSymbolInNamespace(ns *node.Namespace, position Position) node.Statement {
	if ns == nil {
		return nil
	}

	// 遍历命名空间内的语句
	for _, stmt := range ns.Statements {
		if symbolNode := findSymbolInStatement(stmt, position); symbolNode != nil {
			return symbolNode
		}
	}

	return nil
}

// findSymbolInCallExpression 在函数调用表达式中查找符号
func findSymbolInCallExpression(call *node.CallExpression, position Position) node.Statement {
	if call == nil {
		return nil
	}

	// 检查是否是函数名
	if isPositionInNodeName(call, position) {
		return call
	}

	return nil
}

// findSymbolInCallMethod 在方法调用中查找符号
func findSymbolInCallMethod(call *node.CallMethod, position Position) node.Statement {
	if call == nil {
		return nil
	}

	// 检查是否是方法名
	if isPositionInNodeName(call, position) {
		return call
	}

	return nil
}

// findSymbolInNewExpression 在 new 表达式中查找符号
func findSymbolInNewExpression(newExpr *node.NewExpression, position Position) node.Statement {
	if newExpr == nil {
		return nil
	}

	// 检查是否是类名
	if isPositionInNodeName(newExpr, position) {
		return newExpr
	}

	return nil
}

// findSymbolInVariable 在变量中查找符号
func findSymbolInVariable(varStmt *node.VariableExpression, position Position) node.Statement {
	if varStmt == nil {
		return nil
	}

	// 检查是否是变量名
	if isPositionInNodeName(varStmt, position) {
		return varStmt
	}

	return nil
}

// findSymbolInExpression 在表达式中查找符号
func findSymbolInExpression(expr node.Statement, position Position) node.Statement {
	// 处理更多表达式类型
	switch e := expr.(type) {
	case *node.IfStatement:
		// 处理 if 语句
		if isPositionInNode(e, position) {
			return e
		}
	case *node.ForStatement:
		// 处理 for 语句
		if isPositionInNode(e, position) {
			return e
		}
	case *node.WhileStatement:
		// 处理 while 语句
		if isPositionInNode(e, position) {
			return e
		}
	case *node.MatchStatement:
		// 处理 match 语句
		if isPositionInNode(e, position) {
			return e
		}
	case *node.TryStatement:
		// 处理 try 语句
		if isPositionInNode(e, position) {
			return e
		}
	}

	return nil
}

// isPositionInNodeName 检查位置是否在节点名称范围内
func isPositionInNodeName(stmt node.Statement, position Position) bool {
	from := getNodeFrom(stmt)
	if from == nil {
		return false
	}

	start, _ := from.GetPosition()
	startLine, startChar := calculatePositionFromOffset(start)

	// 检查是否在名称范围内（假设名称在开始位置附近）
	if int(position.Line) == startLine {
		endChar := startChar + 50 // 假设名称最大50个字符
		return int(position.Character) >= startChar && int(position.Character) <= endChar
	}

	return false
}

// findDefinitionFromNode 根据节点类型查找定义
func findDefinitionFromNode(symbolNode node.Statement, doc *DocumentInfo) *Location {
	switch node := symbolNode.(type) {
	case *node.FunctionStatement:
		return findFunctionDefinition(doc, node.GetName())
	case *node.ClassStatement:
		return findClassDefinition(doc, node.GetName())
	case *node.InterfaceStatement:
		return findInterfaceDefinition(doc, node.GetName())
	case *node.VariableExpression:
		return findVariableDefinition(doc, node.GetName())
	case *node.CallExpression:
		return findFunctionDefinition(doc, node.FunName)
	case *node.CallMethod:
		// 对于方法调用，需要从方法名中提取类名和方法名
		if method, ok := node.Method.(data.Method); ok {
			return findMethodDefinition(doc, "", method.GetName())
		}
		return nil
	case *node.NewExpression:
		return findClassDefinition(doc, node.ClassName)
	default:
		return nil
	}
}

// findDefinitionFromText 从文本中查找定义（回退方法）
func findDefinitionFromText(content string, position Position) *Location {
	lines := strings.Split(content, "\n")
	if int(position.Line) >= len(lines) {
		return nil
	}

	line := lines[position.Line]
	if int(position.Character) >= len(line) {
		return nil
	}

	// 获取光标位置的单词
	word := getWordAtPosition(content, position)
	if word == "" {
		return nil
	}

	// 检查是否是方法调用 (如 $a->hello())
	if strings.Contains(line, "->") && strings.Contains(word, "->") {
		parts := strings.Split(word, "->")
		if len(parts) == 2 {
			// 这里需要从 LspVM 中查找方法定义
			return findMethodDefinition(nil, parts[0], parts[1])
		}
	}

	// 检查是否是变量 (如 $str)
	if strings.HasPrefix(word, "$") {
		varName := strings.TrimPrefix(word, "$")
		return findVariableDefinition(nil, varName)
	}

	// 检查是否是普通符号（可能是函数或类）
	// 先尝试作为函数查找
	if location := findFunctionDefinition(nil, word); location != nil {
		return location
	}

	// 再尝试作为类查找
	if location := findClassDefinition(nil, word); location != nil {
		return location
	}

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

// findInterfaceDefinition 查找接口定义
func findInterfaceDefinition(doc *DocumentInfo, interfaceName string) *Location {
	if globalLspVM == nil {
		return nil
	}

	if interfaceInfo, exists := globalLspVM.GetInterface(interfaceName); exists {
		return createLocationFromInterface(interfaceInfo)
	}

	return nil
}

// findVariableDefinition 查找变量定义
func findVariableDefinition(doc *DocumentInfo, varName string) *Location {
	if globalLspVM == nil {
		return nil
	}

	if variable, exists := globalLspVM.GetVariable(varName); exists {
		return createLocationFromVariable(variable)
	}

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

// createLocationFromVariable 从变量信息创建位置信息
func createLocationFromVariable(variable *LspVariableInfo) *Location {
	return &Location{
		URI: fmt.Sprintf("file://%s", variable.FilePath),
		Range: Range{
			Start: Position{Line: uint32(variable.Line - 1), Character: 0},
			End:   Position{Line: uint32(variable.Line - 1), Character: 100},
		},
	}
}

// createLocationFromFunction 从函数定义创建位置信息
func createLocationFromFunction(function interface{}) *Location {
	// 尝试从 SimpleFunction 获取位置信息
	if simpleFunction, ok := function.(*SimpleFunction); ok {
		return &Location{
			URI: fmt.Sprintf("file://%s", simpleFunction.GetFilePath()),
			Range: Range{
				Start: Position{Line: uint32(simpleFunction.GetLine() - 1), Character: 0},
				End:   Position{Line: uint32(simpleFunction.GetLine() - 1), Character: uint32(len(simpleFunction.GetContent()))},
			},
		}
	}

	// 尝试从其他类型获取位置信息
	if funcWithPath, ok := function.(interface{ GetFilePath() string }); ok {
		if funcWithLine, ok := function.(interface{ GetLine() int }); ok {
			return &Location{
				URI: fmt.Sprintf("file://%s", funcWithPath.GetFilePath()),
				Range: Range{
					Start: Position{Line: uint32(funcWithLine.GetLine() - 1), Character: 0},
					End:   Position{Line: uint32(funcWithLine.GetLine() - 1), Character: 100},
				},
			}
		}
	}

	return nil
}

// createLocationFromClass 从类定义创建位置信息
func createLocationFromClass(class interface{}) *Location {
	// 尝试从 SimpleClass 获取位置信息
	if simpleClass, ok := class.(*SimpleClass); ok {
		return &Location{
			URI: fmt.Sprintf("file://%s", simpleClass.GetFilePath()),
			Range: Range{
				Start: Position{Line: uint32(simpleClass.GetLine() - 1), Character: 0},
				End:   Position{Line: uint32(simpleClass.GetLine() - 1), Character: uint32(len(simpleClass.GetContent()))},
			},
		}
	}

	// 尝试从其他类型获取位置信息
	if classWithPath, ok := class.(interface{ GetFilePath() string }); ok {
		if classWithLine, ok := class.(interface{ GetLine() int }); ok {
			return &Location{
				URI: fmt.Sprintf("file://%s", classWithPath.GetFilePath()),
				Range: Range{
					Start: Position{Line: uint32(classWithLine.GetLine() - 1), Character: 0},
					End:   Position{Line: uint32(classWithLine.GetLine() - 1), Character: 100},
				},
			}
		}
	}

	return nil
}

// createLocationFromInterface 从接口定义创建位置信息
func createLocationFromInterface(interfaceInfo interface{}) *Location {
	// 尝试从 SimpleInterface 获取位置信息
	if simpleInterface, ok := interfaceInfo.(*SimpleInterface); ok {
		return &Location{
			URI: fmt.Sprintf("file://%s", simpleInterface.GetFilePath()),
			Range: Range{
				Start: Position{Line: uint32(simpleInterface.GetLine() - 1), Character: 0},
				End:   Position{Line: uint32(simpleInterface.GetLine() - 1), Character: uint32(len(simpleInterface.GetContent()))},
			},
		}
	}

	// 尝试从其他类型获取位置信息
	if interfaceWithPath, ok := interfaceInfo.(interface{ GetFilePath() string }); ok {
		if interfaceWithLine, ok := interfaceInfo.(interface{ GetLine() int }); ok {
			return &Location{
				URI: fmt.Sprintf("file://%s", interfaceWithPath.GetFilePath()),
				Range: Range{
					Start: Position{Line: uint32(interfaceWithLine.GetLine() - 1), Character: 0},
					End:   Position{Line: uint32(interfaceWithLine.GetLine() - 1), Character: 100},
				},
			}
		}
	}

	return nil
}

// findMethodInClass 在类中查找方法
func findMethodInClass(class interface{}, methodName string) *Location {
	// 这里需要根据实际的类结构来实现
	// 目前简化处理，返回类定义的位置
	if classWithPath, ok := class.(interface{ GetFilePath() string }); ok {
		if classWithLine, ok := class.(interface{ GetLine() int }); ok {
			return &Location{
				URI: fmt.Sprintf("file://%s", classWithPath.GetFilePath()),
				Range: Range{
					Start: Position{Line: uint32(classWithLine.GetLine() - 1), Character: 0},
					End:   Position{Line: uint32(classWithLine.GetLine() - 1), Character: 100},
				},
			}
		}
	}
	return nil
}
