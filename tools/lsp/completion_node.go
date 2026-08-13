package main

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/tools/lsp/completion"
	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/sirupsen/logrus"
)

// getCompletionItemsWithNodeSupport 获取补全项，支持基于节点的补全
func getCompletionItemsWithNodeSupport(doc *DocumentInfo, position defines.Position, provider *LSPSymbolProvider, vm *LspVM) []defines.CompletionItem {
	if doc == nil {
		return []defines.CompletionItem{}
	}

	// 先检查是否是 -> 或 . 操作符的情况
	// 通过 AST 查找左侧节点
	leftNode, ctx := findLeftNodeForCompletion(doc, position)
	if leftNode != nil {
		logrus.Debugf("找到左侧节点：%T，尝试基于节点获取补全", leftNode)
		items := getCompletionsFromLeftNode(leftNode, ctx, provider, doc.Content, position)
		if len(items) > 0 {
			logrus.Debugf("基于节点获取到 %d 个补全项", len(items))
			return items
		}
	}

	// 如果基于节点无法获取，回退到原有的方式（携带 docProvider 和 vmProvider，便于从 VM 中补全函数）
	return completion.GetCompletionItemsWithDoc(doc.Content, position, provider, doc, vm)
}

// findLeftNodeForCompletion 查找 -> 或 . 左侧的节点
func findLeftNodeForCompletion(doc *DocumentInfo, position defines.Position) (data.GetValue, *LspContext) {
	if doc == nil || doc.AST == nil {
		return nil, nil
	}

	logrus.Debugf("查找左侧节点，位置：(%d, %d)", position.Line, position.Character)

	var targetNode data.GetValue
	var targetCtx *LspContext

	// 使用 DocumentInfo.Foreach 查找光标位置的节点
	doc.Foreach(func(ctx *LspContext, parent, child data.GetValue) bool {
		// 检查当前节点是否包含光标位置
		if isPositionInRange(child, position) {
			logrus.Debugf("找到包含位置的节点：%T", child)
			// 检查是否是 CallObjectProperty 或 CallObjectMethod
			switch child.(type) {
			case *node.CallObjectProperty, *node.CallObjectMethod:
				// 如果找到更精确的节点，更新目标节点
				if targetNode == nil || pickSmallerNode(targetNode, child) == child {
					targetNode = child
					targetCtx = ctx
					logrus.Debugf("更新目标节点：%T", targetNode)
				}
			}
			return true // 继续遍历，寻找更精确的节点
		}

		// 如果节点已经超过了目标行号，则停止遍历
		if getFrom, ok := child.(node.GetFrom); ok {
			if from := getFrom.GetFrom(); from != nil {
				_, _, endLine, _ := from.GetRange()
				targetLine := int(position.Line)
				if endLine > targetLine {
					return false
				}
			}
		}

		return true // 继续遍历
	})

	if targetNode == nil {
		logrus.Debugf("未找到 CallObjectProperty 或 CallObjectMethod 节点")
		return nil, nil
	}

	logrus.Debugf("找到目标节点：%T", targetNode)

	// 根据节点类型获取左侧节点
	var leftNode data.GetValue
	switch n := targetNode.(type) {
	case *node.CallObjectProperty:
		leftNode = n.Object
		logrus.Debugf("CallObjectProperty，左侧节点：%T", leftNode)
	case *node.CallObjectMethod:
		leftNode = n.Object
		logrus.Debugf("CallObjectMethod，左侧节点：%T", leftNode)
	default:
		return nil, nil
	}

	return leftNode, targetCtx
}

// getCompletionsFromLeftNode 根据左侧节点获取补全项
func getCompletionsFromLeftNode(leftNode data.GetValue, ctx *LspContext, provider *LSPSymbolProvider, content string, position defines.Position) []defines.CompletionItem {
	if leftNode == nil || provider == nil {
		return nil
	}

	logrus.Debugf("根据左侧节点获取补全，节点类型：%T", leftNode)

	// 根据左侧节点类型获取类型信息
	varType := getTypeFromLeftNode(leftNode, ctx, provider, content, position)
	if varType == nil {
		logrus.Debugf("无法从左侧节点获取类型")
		return nil
	}

	// 提取所有可能的类名
	classNames := extractClassNamesFromTypeForCompletion(varType)
	if len(classNames) == 0 {
		logrus.Debugf("左侧节点类型不包含任何类")
		return nil
	}

	logrus.Debugf("左侧节点可能的类型: %v", classNames)

	// 收集所有可能类的成员，使用 map 去重
	itemsMap := make(map[string]defines.CompletionItem)
	for _, className := range classNames {
		members := provider.GetClassMembers(className)
		for _, member := range members {
			// 使用 Label 作为唯一键，避免重复
			if _, exists := itemsMap[member.Label]; !exists {
				itemsMap[member.Label] = member
			}
		}
	}

	// 转换为数组
	items := make([]defines.CompletionItem, 0, len(itemsMap))
	for _, item := range itemsMap {
		items = append(items, item)
	}
	logrus.Debugf("基于节点找到合并后的成员：%d 个", len(items))
	return items
}

// getTypeFromLeftNode 从左侧节点获取类型信息
func getTypeFromLeftNode(n data.GetValue, ctx *LspContext, provider *LSPSymbolProvider, content string, position defines.Position) data.Types {
	if n == nil {
		return nil
	}

	// 如果节点是变量表达式，尝试获取变量类型
	if varExpr, ok := n.(*node.VariableExpression); ok {
		// 首先尝试从变量节点的类型信息获取
		if varExpr.Type != nil {
			if typ, ok := varExpr.Type.(data.Types); ok {
				return typ
			}
		}

		// 如果变量节点没有类型信息，尝试从上下文获取
		if ctx != nil {
			varType := ctx.GetVariableType(varExpr.Name)
			if varType != nil {
				logrus.Debugf("从上下文找到变量类型：%s -> %v", varExpr.Name, varType)
				return varType
			}
		}

		// 如果上下文也没有，尝试从 provider 获取
		if provider != nil {
			varTypeObj := provider.GetVariableTypeObjectAtPosition(content, position, varExpr.Name)
			if varTypeObj != nil {
				if typ, ok := varTypeObj.(data.Types); ok {
					return typ
				}
			}
		}
	}

	// 如果节点是 this 表达式，尝试从上下文获取当前类
	if _, ok := n.(*node.This); ok {
		if ctx != nil {
			currentClassName := getCurrentClassNameFromContext(ctx)
			if currentClassName != "" {
				logrus.Debugf("从上下文找到当前类名：%s", currentClassName)
				return data.NewBaseType(currentClassName)
			}
		}
	}

	// 如果节点是 new 表达式，返回类的类型
	if newExpr, ok := n.(*node.NewExpression); ok {
		return data.NewBaseType(newExpr.ClassName)
	}

	// 如果节点是方法调用，尝试推断返回类型
	if callMethod, ok := n.(*node.CallObjectMethod); ok {
		// 递归获取对象的类型
		objType := getTypeFromLeftNode(callMethod.Object, ctx, provider, content, position)
		if objType != nil {
			// 这里可以进一步推断方法调用的返回类型
			// 暂时返回对象类型
			return objType
		}
	}

	return nil
}

// extractClassNamesFromTypeForCompletion 从类型对象中提取所有可能的类名（用于补全）
func extractClassNamesFromTypeForCompletion(typ data.Types) []string {
	if typ == nil {
		return nil
	}

	var classNames []string

	switch t := typ.(type) {
	case *data.LspTypes:
		// LspTypes 包含多个类型，遍历所有内部类型
		for _, innerType := range t.Types {
			names := extractClassNamesFromTypeForCompletion(innerType)
			classNames = append(classNames, names...)
		}
	case data.Class:
		classNames = append(classNames, t.Name)
	case data.NullableType:
		// 可空类型，递归获取基础类型
		names := extractClassNamesFromTypeForCompletion(t.BaseType)
		classNames = append(classNames, names...)
	default:
		// 对于其他类型，尝试使用 String() 方法
		typeStr := typ.String()
		// 如果不是基础类型，可能是类名
		if !data.ISBaseType(typeStr) {
			classNames = append(classNames, typeStr)
		}
	}

	return classNames
}
