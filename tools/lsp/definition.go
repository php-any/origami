package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"
)

// 处理定义跳转请求
func handleTextDocumentDefinition(req *jsonrpc2.Request) (interface{}, error) {
	var params DefinitionParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal definition params: %v", err)
	}

	uri := params.TextDocument.URI
	position := params.Position

	logrus.Infof("请求定义跳转：%s 位置 %d:%d; req: %v", uri, position.Line, position.Character, params)

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

	logrus.Debugf("开始查找定义，位置：(%d, %d)", position.Line, position.Character)

	// 使用 DocumentInfo.Foreach 查找光标位置的节点和对应的上下文
	var targetNode data.GetValue
	var targetCtx *LspContext
	doc.Foreach(func(ctx *LspContext, parent, child data.GetValue) bool {
		// 检查当前节点是否包含光标位置
		if isPositionInRange(child, position) {
			logrus.Debugf("找到包含位置的节点：%T，父节点：%T", child, parent)
			// 如果找到包含位置的节点，选择最小的（最精确的）
			if pickSmallerNode(targetNode, child) == child {
				targetNode = child
				targetCtx = ctx // 保存目标节点对应的上下文
				logrus.Debugf("更新目标节点：%T", targetNode)
			}
			return true // 继续遍历，寻找更精确的节点
		}

		// 如果节点已经超过了目标行号，则停止遍历
		if getFrom, ok := child.(node.GetFrom); ok {
			if from := getFrom.GetFrom(); from != nil {
				_, _, endLine, _ := from.GetRange()
				// LSP 和内部系统都使用从 0 开始的行号
				targetLine := int(position.Line)

				// 如果当前节点的结束行已经超过了目标行，停止遍历
				if endLine > targetLine {
					logrus.Debugf("节点结束行 %d 超过目标行 %d，停止遍历", endLine, targetLine)
					return false
				}
			}
		}

		return true // 继续遍历
	})

	if targetNode == nil {
		logrus.Debug("未找到包含位置的节点")
		return nil
	}

	logrus.Debugf("最终目标节点：%T", targetNode)

	// 根据节点类型查找定义，使用目标节点的上下文
	return findDefinitionFromNode(targetCtx, targetNode)
}

// findDefinitionFromNode 根据节点类型查找定义位置
func findDefinitionFromNode(ctx *LspContext, v data.GetValue) *Location {
	if v == nil {
		return nil
	}

	switch n := v.(type) {
	case *node.CallExpression:
		// 函数调用，查找函数定义
		logrus.Debugf("CallExpression FunName: %s", n.FunName)
		return findFunctionDefinition(ctx, n.FunName)
	case *node.NewExpression:
		// new 表达式，从类名查找类定义
		return findClassDefinition(ctx, n.ClassName)
	case *node.CallMethod:
		// 方法调用，查找方法定义
		// 注意：CallMethod 的 Method 字段是 data.GetValue 类型，需要进一步处理
		if method, ok := n.Method.(interface{ GetName() string }); ok {
			return findMethodDefinition(ctx, "", method.GetName())
		}
		return nil
	case *node.CallObjectMethod:
		// 对象方法调用，查找方法定义
		logrus.Debugf("CallObjectMethod: object=%T, method=%s", n.Object, n.Method)

		// 检查是否是链式调用（对象是另一个方法调用的结果）
		if chainCall, ok := n.Object.(*node.CallObjectMethod); ok {
			// 这是一个链式调用，需要递归解析
			logrus.Debug("检测到链式调用，递归解析")
			return findChainedMethodDefinition(ctx, chainCall, n.Method)
		}

		// 普通对象方法调用
		return findObjectMethodDefinition(ctx, n.Object, n.Method)
	case *node.VariableExpression:
		// 变量引用，查找变量定义
		return findVariableDefinition(ctx, n.Name)
	}

	// 对于其他类型的节点，暂时返回 nil
	return nil
}

// pickSmallerNode 返回最合适的节点；优先选择包含光标位置的节点
func pickSmallerNode(a, b data.GetValue) data.GetValue {
	logrus.Debugf("pickSmallerNode：a=%T，b=%T", a, b)

	if b == nil {
		logrus.Debug("pickSmallerNode：b 为空，返回 a")
		return a
	}
	if a == nil {
		logrus.Debug("pickSmallerNode：a 为空，返回 b")
		return b
	}

	af := getFromOf(a)
	bf := getFromOf(b)
	if af == nil {
		logrus.Debug("pickSmallerNode：af 为空，返回 b")
		return b
	}
	if bf == nil {
		logrus.Debug("pickSmallerNode：bf 为空，返回 a")
		return a
	}

	// 获取两个节点的范围
	slA, scA, elA, ecA := af.GetRange()
	slB, scB, elB, ecB := bf.GetRange()

	logrus.Debugf("pickSmallerNode：a 范围=(%d,%d,%d,%d)，b 范围=(%d,%d,%d,%d)",
		slA, scA, elA, ecA, slB, scB, elB, ecB)

	// 计算两个节点的范围大小（字符数）
	rangeA := (elA-slA+1)*1000 + (ecA - scA + 1)
	rangeB := (elB-slB+1)*1000 + (ecB - scB + 1)

	logrus.Debugf("pickSmallerNode：rangeA=%d，rangeB=%d", rangeA, rangeB)

	// 优先选择范围更小的节点（更精确），但前提是它们都包含光标位置
	// 如果范围差异不大，选择范围更小的；如果差异很大，选择更合适的
	if rangeB < rangeA && (rangeA-rangeB) < 100 {
		logrus.Debug("pickSmallerNode：rangeB < rangeA 且差异较小，返回 b")
		return b
	}
	if rangeA < rangeB && (rangeB-rangeA) < 100 {
		logrus.Debug("pickSmallerNode：rangeA < rangeB 且差异较小，返回 a")
		return a
	}

	// 如果范围差异很大，选择范围更小的（更精确）
	if rangeB < rangeA {
		logrus.Debug("pickSmallerNode：rangeB 明显更小，返回 b")
		return b
	}
	logrus.Debug("pickSmallerNode：rangeA <= rangeB，返回 a")
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

	// LSP 和内部系统都使用从 0 开始的行号
	lspLine := int(position.Line)

	// 添加调试信息
	logrus.Debugf("isPositionInRange：节点=%T，定位位置=(%d,%d)，正在查范围=(%d,%d,%d,%d)",
		stmt, lspLine, position.Character, startLine, startChar, endLine, endChar)

	// 检查行号是否在范围内
	if lspLine < startLine || lspLine > endLine {
		return false
	}

	// 如果在起始行，检查字符位置是否在起始字符之后
	if lspLine == startLine {
		if lspLine == endLine {
			// 单行节点：字符位置必须在起始和结束字符之间
			result := int(position.Character) >= startChar && int(position.Character) <= endChar
			logrus.Debugf("isPositionInRange：单行节点，字符在范围内：%v", result)
			return result
		} else {
			// 多行节点的起始行：字符位置必须在起始字符之后
			result := int(position.Character) >= startChar
			logrus.Debugf("isPositionInRange：多行起始，字符 >= 起始：%v", result)
			return result
		}
	}

	// 如果在结束行，检查字符位置是否在结束字符之前
	if lspLine == endLine {
		result := int(position.Character) <= endChar
		logrus.Debugf("isPositionInRange：结束行，字符 <= 结束：%v", result)
		return true // 简化：只要在结束行就认为在范围内
	}

	// 如果在中间行，肯定在范围内
	return true
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
		// 检查位置是否在表达式范围内
		if isPositionInRange(expr, position) || isPositionInLineRange(expr, position) {
			return expr
		}
		return nil
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

// 仅检查"行"是否命中，忽略列，避免 token.Pos 为末尾列导致的误差
func isPositionInLineRange(stmt node.Statement, position Position) bool {
	var from data.From
	if getFrom, ok := stmt.(node.GetFrom); ok {
		from = getFrom.GetFrom()
	}
	if from == nil {
		return false
	}
	startLine, _, endLine, _ := from.GetRange()
	lspLine := int(position.Line)

	// 简化：只要在行范围内就认为命中
	result := lspLine >= startLine && lspLine <= endLine
	logrus.Debugf("isPositionInLineRange：节点=%T，位置行=%d，范围行=[%d,%d]，结果=%v", stmt, lspLine, startLine, endLine, result)
	return result
}

// findFunctionDefinition 查找函数定义
func findFunctionDefinition(ctx *LspContext, funcName string) *Location {
	if globalLspVM == nil {
		return nil
	}

	logrus.Debugf("查找函数定义：%s", funcName)
	if function, exists := globalLspVM.GetFunc(funcName); exists {
		logrus.Debugf("找到函数：%#v，位置：%#v", function, function)
		return createLocationFromFunction(function)
	}

	logrus.Debugf("未找到函数：%s", funcName)
	return nil
}

// findClassDefinition 查找类定义
func findClassDefinition(ctx *LspContext, className string) *Location {
	if globalLspVM == nil {
		return nil
	}

	if class, exists := globalLspVM.GetClass(className); exists {
		return createLocationFromClass(class)
	}

	return nil
}

// findObjectMethodDefinition 查找对象方法定义
func findObjectMethodDefinition(ctx *LspContext, object data.GetValue, methodName string) *Location {
	if globalLspVM == nil {
		return nil
	}

	// 尝试从对象中获取类型信息
	// 这里需要根据对象的类型来查找方法定义
	// 由于对象可能是变量、表达式等，我们需要先尝试解析对象名称

	// 如果对象是变量表达式，尝试从变量类型查找对应的类
	if varExpr, ok := object.(*node.VariableExpression); ok {
		// 首先尝试从变量节点的类型信息获取
		if varExpr.Type != nil {
			// 从类型信息中获取类名
			if className := getClassNameFromType(varExpr.Type); className != "" {
				// 根据类名查找类定义
				if class, exists := globalLspVM.GetClass(className); exists {
					if methodLocation := findMethodInClass(class, methodName); methodLocation != nil {
						return methodLocation
					}
				}
			}
		}

		// 如果变量节点没有类型信息，尝试从上下文获取
		if ctx != nil {
			varType := ctx.GetVariableType(varExpr.Name)
			if varType != nil {
				logrus.Debugf("从上下文找到变量类型：%s -> %v", varExpr.Name, varType)
				if className := getClassNameFromType(varType); className != "" {
					logrus.Debugf("提取类名：%s", className)
					// 根据类名查找类定义
					if class, exists := globalLspVM.GetClass(className); exists {
						if methodLocation := findMethodInClass(class, methodName); methodLocation != nil {
							return methodLocation
						}
					}
				}
			} else {
				logrus.Debugf("在上下文中未找到变量类型：%s", varExpr.Name)
			}
		}
	}

	// 如果对象是 this 表达式，尝试查找当前类的方法
	if _, ok := object.(*node.This); ok {
		// 这里需要获取当前类的上下文，暂时简化处理
		// 可以尝试查找同名函数作为备选
		if function, exists := globalLspVM.GetFunc(methodName); exists {
			return createLocationFromFunction(function)
		}
	}

	// 如果找不到具体的类，尝试查找同名函数作为备选
	if function, exists := globalLspVM.GetFunc(methodName); exists {
		return createLocationFromFunction(function)
	}

	return nil
}

// findChainedMethodDefinition 查找链式方法调用中的方法定义
func findChainedMethodDefinition(ctx *LspContext, chainCall *node.CallObjectMethod, methodName string) *Location {
	if globalLspVM == nil {
		return nil
	}

	logrus.Debugf("解析链式调用：%s", methodName)

	// 递归解析链式调用，从最内层开始
	// 例如：$a->newB()->getC()->hello() 需要递归解析每一层

	// 获取链式调用的对象
	object := chainCall.Object
	method := chainCall.Method

	logrus.Debugf("链式调用对象：%T，方法：%s", object, method)

	// 递归解析链式调用
	return resolveChainedMethod(ctx, object, method, methodName)
}

// resolveChainedMethod 递归解析链式方法调用
func resolveChainedMethod(ctx *LspContext, object data.GetValue, currentMethod, targetMethod string) *Location {
	if globalLspVM == nil {
		return nil
	}

	logrus.Debugf("解析方法调用：%s，目标方法：%s", currentMethod, targetMethod)

	// 如果对象是变量，尝试获取其类型
	if varExpr, ok := object.(*node.VariableExpression); ok {
		logrus.Debugf("变量：%s", varExpr.Name)

		// 从上下文获取变量类型
		if ctx != nil {
			varType := ctx.GetVariableType(varExpr.Name)
			if varType != nil {
				logrus.Debugf("变量类型：%v", varType)

				// 获取类名
				if className := getClassNameFromType(varType); className != "" {
					logrus.Debugf("类名：%s", className)

					// 查找类定义
					if class, exists := globalLspVM.GetClass(className); exists {
						// 在类中查找当前方法
						if methodLocation := findMethodInClass(class, currentMethod); methodLocation != nil {
							logrus.Debugf("找到方法：%s", currentMethod)

							// 尝试推断方法的返回类型
							returnType := inferMethodReturnType(class, currentMethod)
							if returnType != nil {
								logrus.Debugf("推断返回类型：%s", returnType)

								// 仅补充 switch returnType.(type) 逻辑
								switch rt := returnType.(type) {
								case data.String:
									if function, exists := globalLspVM.GetFunc(targetMethod); exists {
										return createLocationFromFunction(function)
									}
								case data.Arrays:
									if function, exists := globalLspVM.GetFunc(targetMethod); exists {
										return createLocationFromFunction(function)
									}
								case data.Int, data.Float, data.Bool, data.Object, data.Callable:
									if function, exists := globalLspVM.GetFunc(targetMethod); exists {
										return createLocationFromFunction(function)
									}
								case data.Class:
									if class, exists := globalLspVM.GetClass(rt.Name); exists {
										if final := findMethodInClass(class, targetMethod); final != nil {
											return final
										}
									}
								case data.NullableType:
									if className := getClassNameFromType(rt.BaseType); className != "" {
										if class, exists := globalLspVM.GetClass(className); exists {
											if final := findMethodInClass(class, targetMethod); final != nil {
												return final
											}
										}
									}
								case data.MultipleReturnType:
									for _, t := range rt.Types {
										if className := getClassNameFromType(t); className != "" {
											if class, exists := globalLspVM.GetClass(className); exists {
												if final := findMethodInClass(class, targetMethod); final != nil {
													return final
												}
											}
										}
									}
								case data.Generic:
									for _, t := range rt.Types {
										if className := getClassNameFromType(t); className != "" {
											if class, exists := globalLspVM.GetClass(className); exists {
												if final := findMethodInClass(class, targetMethod); final != nil {
													return final
												}
											}
										}
									}
								}
							}

							// 如果无法推断返回类型，尝试在所有类中查找目标方法
							if globalLspVM != nil {
								// 遍历所有类，查找目标方法
								allClasses := globalLspVM.GetAllClasses()
								for className, classStmt := range allClasses {
									logrus.Debugf("在所有类中查找：%s", className)
									if methodLocation := findMethodInClass(classStmt, targetMethod); methodLocation != nil {
										logrus.Debugf("在类 %s 中找到目标方法：%s", className, targetMethod)
										return methodLocation
									}
								}
							}

							// 如果所有类中都找不到，就不需要提示了
							logrus.Debugf("在所有类中都找不到目标方法：%s", targetMethod)
						}
					}
				}
			}
		}
	}

	// 如果对象是另一个方法调用，递归解析
	if nestedCall, ok := object.(*node.CallObjectMethod); ok {
		logrus.Debug("检测到嵌套方法调用，递归解析")
		return resolveChainedMethod(ctx, nestedCall.Object, nestedCall.Method, targetMethod)
	}

	// 如果无法解析链式调用，尝试查找同名函数作为备选
	logrus.Debugf("无法解析链式调用，尝试查找同名函数：%s", targetMethod)
	if function, exists := globalLspVM.GetFunc(targetMethod); exists {
		return createLocationFromFunction(function)
	}

	return nil
}

// inferMethodReturnType 推断方法的返回类型
func inferMethodReturnType(class data.ClassStmt, methodName string) data.Types {
	// 获取方法定义
	method, exists := class.GetMethod(methodName)
	if !exists {
		return nil
	}
	// 1. 检查方法是否有返回类型注解
	if ret, ok := method.(data.GetReturnType); ok {
		ret := ret.GetReturnType()
		if ret != nil {
			return ret
		}
	}
	var inferredType data.Types
	// 方法应该实现 data.Method 接口
	if m, ok := method.(*node.ClassMethod); ok {
		docu := &DocumentInfo{}
		baseCtx := context.Background()
		lspCtx := NewLspContext(baseCtx, nil)

		for _, stmt := range m.Body {
			docu.foreachNode(lspCtx, stmt, nil, func(ctx *LspContext, parent, child data.GetValue) bool {
				switch st := child.(type) {
				case *node.ReturnStatement:
					inferredType = docu.identifyVariableTypes(ctx, st)
					return false
				}
				return true
			})
			if inferredType != nil {
				return inferredType
			}
		}
	}

	// 默认返回空字符串，表示无法推断
	return inferredType
}

// findVariableDefinition 查找变量定义
func findVariableDefinition(ctx *LspContext, varName string) *Location {
	// 变量定义查找功能暂时简化，返回 nil
	return nil
}

// findMethodDefinition 查找方法定义
func findMethodDefinition(ctx *LspContext, variableName, methodName string) *Location {
	if globalLspVM == nil {
		return nil
	}

	// 如果有变量名，查找变量对应的类
	if variableName != "" {
		if class, exists := globalLspVM.GetClass(variableName); exists {
			if methodLocation := findMethodInClass(class, methodName); methodLocation != nil {
				return methodLocation
			}
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
				URI: filePathToURI(from.GetSource()),
				Range: Range{
					Start: Position{Line: uint32(startLine) + 1, Character: uint32(startChar)},
					End:   Position{Line: uint32(endLine) + 1, Character: uint32(endChar)},
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
			URI: filePathToURI(from.GetSource()),
			Range: Range{
				Start: Position{Line: uint32(startLine) + 1, Character: uint32(startChar)},
				End:   Position{Line: uint32(endLine) + 1, Character: uint32(endChar)},
			},
		}
	}
	return nil
}

// getClassNameFromType 从类型信息中获取类名
func getClassNameFromType(typ data.Types) string {
	switch t := typ.(type) {
	case data.Class:
		return t.Name
	case data.NullableType:
		// 如果是可空类型，递归获取基础类型的类名
		return getClassNameFromType(t.BaseType)
	default:
		// 对于其他类型，尝试使用 String() 方法
		typeStr := typ.String()
		// 如果不是基础类型，可能是类名
		if !data.ISBaseType(typeStr) {
			return typeStr
		}
		return ""
	}
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
					URI: filePathToURI(from.GetSource()),
					Range: Range{
						Start: Position{Line: uint32(startLine) + 1, Character: uint32(startChar)},
						End:   Position{Line: uint32(endLine) + 1, Character: uint32(endChar)},
					},
				}
			}
		}
	}

	// 如果找不到方法或方法没有位置信息，返回类定义的位置
	if from := class.GetFrom(); from != nil {
		startLine, startChar, endLine, endChar := from.ToLSPPosition()
		return &Location{
			URI: filePathToURI(from.GetSource()),
			Range: Range{
				Start: Position{Line: uint32(startLine) + 1, Character: uint32(startChar)},
				End:   Position{Line: uint32(endLine) + 1, Character: uint32(endChar)},
			},
		}
	}
	return nil
}
