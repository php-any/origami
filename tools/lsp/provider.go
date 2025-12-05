package main

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/sirupsen/logrus"
)

// LSPSymbolProvider 实现了 defines.SymbolProvider 接口
type LSPSymbolProvider struct {
	doc *DocumentInfo
	vm  *LspVM
}

// GetVariableTypeAtPosition 获取指定位置变量的类型（类名）
func (p *LSPSymbolProvider) GetVariableTypeAtPosition(content string, position defines.Position, varName string) string {
	foundType := p.GetVariableTypeObjectAtPosition(content, position, varName)
	if foundType == nil {
		return ""
	}

	if typ, ok := foundType.(data.Types); ok {
		return getClassNameFromType(typ)
	}

	return ""
}

// GetVariableTypeObjectAtPosition 获取指定位置变量的类型对象（可能包含多个类型）
func (p *LSPSymbolProvider) GetVariableTypeObjectAtPosition(content string, position defines.Position, varName string) interface{} {
	if p.doc == nil || p.doc.AST == nil {
		return nil
	}

	logrus.Debugf("GetVariableTypeObjectAtPosition: varName=%s, pos=%d:%d", varName, position.Line, position.Character)

	// 确保变量名格式统一
	varNameWithDollar := varName
	if !strings.HasPrefix(varName, "$") {
		varNameWithDollar = "$" + varName
	}
	varNameWithoutDollar := strings.TrimPrefix(varNameWithDollar, "$")

	targetLine := int(position.Line)
	var collectedTypes []data.Types
	var targetContext *LspContext

	// 遍历 AST，收集光标位置之前的所有赋值类型
	p.doc.Foreach(func(ctx *LspContext, parent, child data.GetValue) bool {
		// 获取节点的位置信息
		var nodeStartLine int = -1
		if getFrom, ok := child.(node.GetFrom); ok {
			if from := getFrom.GetFrom(); from != nil {
				startLine, _, endLine, _ := from.GetRange()
				nodeStartLine = startLine

				// 更新目标上下文（找到包含或最接近目标位置的上下文）
				if startLine <= targetLine && endLine >= targetLine {
					targetContext = ctx
				}
			}
		}

		// 只处理开始行在目标行之前或同一行的节点
		if nodeStartLine > targetLine {
			return true // 跳过这个节点，继续遍历其他节点
		}

		// 检查是否是变量赋值语句
		switch n := child.(type) {
		case *node.BinaryAssignVariable:
			// 检查左侧是否是我们要查找的变量
			if leftVar, ok := n.Left.(*node.VariableExpression); ok {
				if leftVar.Name == varNameWithoutDollar || "$"+leftVar.Name == varNameWithDollar {
					// 推断右侧表达式的类型
					if inferredType := inferTypeFromExpression(n.Right); inferredType != nil {
						logrus.Debugf("找到变量 %s 的赋值，类型: %v", varName, inferredType)
						collectedTypes = append(collectedTypes, inferredType)
					}
				}
			}
		case *node.VarStatement:
			// var 声明
			if n.Name == varNameWithoutDollar || "$"+n.Name == varNameWithDollar {
				if n.Initializer != nil {
					if inferredType := inferTypeFromExpression(n.Initializer); inferredType != nil {
						logrus.Debugf("找到变量 %s 的 var 声明，类型: %v", varName, inferredType)
						collectedTypes = append(collectedTypes, inferredType)
					}
				}
			}
		case *node.Parameter:
			// 函数参数
			paramName := n.Name
			if !strings.HasPrefix(paramName, "$") {
				paramName = "$" + paramName
			}
			if paramName == varNameWithDollar {
				if n.Type != nil {
					logrus.Debugf("找到变量 %s 的参数定义，类型: %v", varName, n.Type)
					collectedTypes = append(collectedTypes, n.Type)
				}
			}
		}

		return true
	})

	// 如果没有从赋值中找到类型，尝试从上下文中查找
	if len(collectedTypes) == 0 && targetContext != nil {
		if typ := targetContext.GetVariableType(varNameWithDollar); typ != nil {
			logrus.Debugf("从上下文找到变量 %s 的类型: %v", varName, typ)
			return typ
		}
		if typ := targetContext.GetVariableType(varNameWithoutDollar); typ != nil {
			logrus.Debugf("从上下文找到变量 %s 的类型: %v", varName, typ)
			return typ
		}
	}

	// 如果只收集到一个类型，直接返回
	if len(collectedTypes) == 1 {
		logrus.Debugf("变量 %s 有单一类型: %v", varName, collectedTypes[0])
		return collectedTypes[0]
	}

	// 如果收集到多个类型，合并为 LspTypes
	if len(collectedTypes) > 1 {
		logrus.Debugf("变量 %s 有多个类型（%d个），合并为 LspTypes", varName, len(collectedTypes))
		return &data.LspTypes{
			Types: collectedTypes,
		}
	}

	logrus.Debugf("未找到变量 %s 的类型", varName)
	return nil
}

// GetClassMembers 获取类的所有成员（属性和方法）作为补全项
// 支持继承链查找
func (p *LSPSymbolProvider) GetClassMembers(className string) []defines.CompletionItem {
	if p.vm == nil {
		return nil
	}

	// 使用 map 去重，键为成员名
	itemsMap := make(map[string]defines.CompletionItem)

	currentClassName := className
	visited := make(map[string]bool)

	logrus.Debugf("开始查找类成员: %s", className)

	// 循环向上查找父类
	for currentClassName != "" {
		// 防止循环继承死循环
		if visited[currentClassName] {
			break
		}
		visited[currentClassName] = true

		class, exists := p.vm.GetClass(currentClassName)
		if !exists {
			logrus.Debugf("未找到类定义: %s", currentClassName)
			break
		}

		logrus.Debugf("正在分析类: %s", currentClassName)

		if classStmt, ok := class.(*node.ClassStatement); ok {
			// 添加方法
			for _, method := range classStmt.Methods {
				label := method.GetName()

				// 如果子类已经定义了该方法，跳过父类的
				if _, exists := itemsMap[label]; exists {
					continue
				}

				// 检查可见性
				visibility := "public"
				modifier := method.GetModifier()
				if modifier == data.ModifierPrivate {
					visibility = "private"
					// 私有方法只能在当前类访问，如果是父类则跳过
					if currentClassName != className {
						continue
					}
				} else if modifier == data.ModifierProtected {
					visibility = "protected"
				}

				// 格式化参数
				paramDisplay, paramSnippet := formatMethodParams(method)

				// 构建详细信息：public method method(params): ReturnType in ClassName
				returnType := "mixed"
				if ret := method.GetReturnType(); ret != nil {
					returnType = ret.String()
				}
				detail := fmt.Sprintf("%s method %s%s: %s in %s", visibility, label, paramDisplay, returnType, currentClassName)

				insertText := label + paramSnippet
				insertTextFormat := defines.InsertTextFormatSnippet

				item := defines.CompletionItem{
					Label:            label,
					Kind:             &[]defines.CompletionItemKind{defines.CompletionItemKindMethod}[0],
					Detail:           &detail,
					InsertText:       &insertText,
					InsertTextFormat: &insertTextFormat,
				}
				itemsMap[label] = item
			}

			// 添加属性
			for _, prop := range classStmt.Properties {
				label := prop.GetName()

				if _, exists := itemsMap[label]; exists {
					continue
				}

				visibility := "public"
				modifier := prop.GetModifier()
				if modifier == data.ModifierPrivate {
					visibility = "private"
					if currentClassName != className {
						continue
					}
				} else if modifier == data.ModifierProtected {
					visibility = "protected"
				}

				detail := fmt.Sprintf("%s property in %s", visibility, currentClassName)

				item := defines.CompletionItem{
					Label:  label,
					Kind:   &[]defines.CompletionItemKind{defines.CompletionItemKindProperty}[0],
					Detail: &detail,
				}
				itemsMap[label] = item
			}

			// 查找父类
			if extends := classStmt.GetExtend(); extends != nil {
				currentClassName = *extends
			} else {
				currentClassName = ""
			}
		} else {
			break
		}
	}

	// 将 map 转换为 slice
	items := make([]defines.CompletionItem, 0, len(itemsMap))
	for _, item := range itemsMap {
		items = append(items, item)
	}

	return items
}

// GetStaticClassMembers 获取类的静态成员
func (p *LSPSymbolProvider) GetStaticClassMembers(className string) []defines.CompletionItem {
	if p.vm == nil {
		return nil
	}

	itemsMap := make(map[string]defines.CompletionItem)
	currentClassName := className
	visited := make(map[string]bool)

	logrus.Debugf("开始查找类静态成员: %s", className)

	for currentClassName != "" {
		if visited[currentClassName] {
			break
		}
		visited[currentClassName] = true

		class, exists := p.vm.GetClass(currentClassName)
		if !exists {
			break
		}

		if classStmt, ok := class.(*node.ClassStatement); ok {
			// 添加静态方法
			for _, method := range classStmt.Methods {
				// 检查是否是静态方法
				// 暂时无法直接判断是否为静态方法，先全部包含
				// 实际项目中应该检查 method.GetModifier() 是否包含 static 标志

				label := method.GetName()
				if _, exists := itemsMap[label]; exists {
					continue
				}

				visibility := "public"
				modifier := method.GetModifier()
				if modifier == data.ModifierPrivate {
					visibility = "private"
					if currentClassName != className {
						continue
					}
				} else if modifier == data.ModifierProtected {
					visibility = "protected"
				}

				// 格式化参数
				paramDisplay, paramSnippet := formatMethodParams(method)

				// 构建详细信息
				returnType := "mixed"
				if ret := method.GetReturnType(); ret != nil {
					returnType = ret.String()
				}
				detail := fmt.Sprintf("%s static method %s%s: %s in %s", visibility, label, paramDisplay, returnType, currentClassName)

				insertText := label + paramSnippet
				insertTextFormat := defines.InsertTextFormatSnippet

				item := defines.CompletionItem{
					Label:            label,
					Kind:             &[]defines.CompletionItemKind{defines.CompletionItemKindMethod}[0],
					Detail:           &detail,
					InsertText:       &insertText,
					InsertTextFormat: &insertTextFormat,
				}
				itemsMap[label] = item
			}

			// 添加静态属性
			for _, prop := range classStmt.Properties {
				// 同样需要检查是否静态
				label := prop.GetName()
				if _, exists := itemsMap[label]; exists {
					continue
				}

				// 属性必须是静态的才能用 :: 访问
				// 假设 IsStatic 存在
				// if !prop.IsStatic() { continue }

				visibility := "public"
				modifier := prop.GetModifier()
				if modifier == data.ModifierPrivate {
					visibility = "private"
					if currentClassName != className {
						continue
					}
				} else if modifier == data.ModifierProtected {
					visibility = "protected"
				}

				detail := fmt.Sprintf("%s static property in %s", visibility, currentClassName)
				item := defines.CompletionItem{
					Label:  label,
					Kind:   &[]defines.CompletionItemKind{defines.CompletionItemKindProperty}[0],
					Detail: &detail,
				}
				itemsMap[label] = item
			}

			if extends := classStmt.GetExtend(); extends != nil {
				currentClassName = *extends
			} else {
				currentClassName = ""
			}
		} else {
			break
		}
	}

	items := make([]defines.CompletionItem, 0, len(itemsMap))
	for _, item := range itemsMap {
		items = append(items, item)
	}
	return items
}

// GetVariablesAtPosition 获取指定位置的所有可用变量
func (p *LSPSymbolProvider) GetVariablesAtPosition(content string, position defines.Position) []defines.CompletionItem {
	if p.doc == nil || p.doc.AST == nil {
		return nil
	}

	itemsMap := make(map[string]defines.CompletionItem)
	var targetContext *LspContext

	targetLine := int(position.Line)

	// 遍历 AST 找到包含光标位置的最深层上下文
	p.doc.Foreach(func(ctx *LspContext, parent, child data.GetValue) bool {
		// 使用行号比较，而不是严格的位置检查
		// 找到包含目标行的最内层节点
		if getFrom, ok := child.(node.GetFrom); ok {
			if from := getFrom.GetFrom(); from != nil {
				startLine, _, endLine, _ := from.GetRange()

				// 如果当前节点包含目标行
				if targetLine >= startLine && targetLine <= endLine {
					// 更新目标上下文为更深层的上下文
					targetContext = ctx
				}
			}
		}
		return true
	})

	// 如果没有找到特定上下文，使用根上下文
	if targetContext == nil {
		// 创建一个默认上下文并遍历一次以填充
		p.doc.Foreach(func(ctx *LspContext, parent, child data.GetValue) bool {
			if targetContext == nil {
				targetContext = ctx
			}
			return false // 只需要第一个
		})
	}

	// 从目标上下文开始，向上遍历所有父上下文，收集所有变量
	curr := targetContext
	for curr != nil {
		// 1. 从 values 中获取类型推断的变量
		for key, val := range curr.values {
			if strings.HasPrefix(key, "var_type:") {
				varName := strings.TrimPrefix(key, "var_type:")
				if _, exists := itemsMap[varName]; !exists {
					detail := "Variable"
					if t, ok := val.(data.Types); ok {
						detail = getClassNameFromType(t)
					}
					itemsMap[varName] = defines.CompletionItem{
						Label:  varName,
						Kind:   &[]defines.CompletionItemKind{defines.CompletionItemKindVariable}[0],
						Detail: &detail,
					}
				}
			}
		}

		// 2. 从 localVars 获取
		for name := range curr.localVars {
			if _, exists := itemsMap[name]; !exists {
				itemsMap[name] = defines.CompletionItem{
					Label: name,
					Kind:  &[]defines.CompletionItemKind{defines.CompletionItemKindVariable}[0],
				}
			}
		}

		curr = curr.parent
	}

	items := make([]defines.CompletionItem, 0, len(itemsMap))
	for _, item := range itemsMap {
		items = append(items, item)
	}
	return items
}

// formatMethodParams 格式化方法参数，返回 (参数显示字符串, snippet插入字符串)
func formatMethodParams(method data.Method) (string, string) {
	var paramDisplays []string
	var paramSnippets []string

	params := method.GetParams()
	for i, p := range params {
		// 尝试转换为 data.Parameter 接口
		if param, ok := p.(data.Parameter); ok {
			name := param.GetName()
			// 确保变量名带 $
			if !strings.HasPrefix(name, "$") {
				name = "$" + name
			}

			// 获取类型
			typeName := ""
			if t := param.GetType(); t != nil {
				typeName = t.String() + " "
			}

			// 显示格式: Type $name
			paramDisplays = append(paramDisplays, fmt.Sprintf("%s%s", typeName, name))

			// Snippet 格式: ${1:$name}
			paramSnippets = append(paramSnippets, fmt.Sprintf("${%d:%s}", i+1, name))
		}
	}

	display := "(" + strings.Join(paramDisplays, ", ") + ")"
	snippet := "(" + strings.Join(paramSnippets, ", ") + ")"

	return display, snippet
}
