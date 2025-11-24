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
	if p.doc == nil || p.doc.AST == nil {
		return ""
	}

	logrus.Debugf("GetVariableTypeAtPosition: varName=%s, pos=%d:%d", varName, position.Line, position.Character)

	// 确保变量名格式统一，优先尝试带 $ 的
	varNameWithDollar := varName
	if !strings.HasPrefix(varName, "$") {
		varNameWithDollar = "$" + varName
	}
	varNameWithoutDollar := strings.TrimPrefix(varNameWithDollar, "$")

	// 使用 Foreach 遍历 AST，找到包含光标位置的节点，并检查其中的变量类型
	var foundType data.Types

	p.doc.Foreach(func(ctx *LspContext, parent, child data.GetValue) bool {
		// 检查当前节点是否包含光标位置
		// 注意：如果当前是在输入新代码，光标可能不在任何现有语句的范围内
		// 但通常我们会在一个 BlockStatement 或 FunctionStatement 内部
		if isPositionInRange(child, position) {
			// 在当前上下文中查找变量类型
			// 优先查找带 $ 的
			if typ := ctx.GetVariableType(varNameWithDollar); typ != nil {
				foundType = typ
				logrus.Debugf("在上下文 %s 中找到变量 %s 的类型: %v", ctx.scopeName, varNameWithDollar, typ)
			} else if typ := ctx.GetVariableType(varNameWithoutDollar); typ != nil {
				foundType = typ
				logrus.Debugf("在上下文 %s 中找到变量 %s 的类型: %v", ctx.scopeName, varNameWithoutDollar, typ)
			}
			return true
		}

		// 如果节点已经超过了目标行号，则停止遍历该分支的后续部分（优化）
		// 但因为 AST 结构不是完全线性的，简单的行号判断可能不够
		// 这里主要依赖 Foreach 的逻辑
		if getFrom, ok := child.(node.GetFrom); ok {
			if from := getFrom.GetFrom(); from != nil {
				_, _, endLine, _ := from.GetRange()
				targetLine := int(position.Line) + 1 // convert 0-based to 1-based
				// 如果当前节点结束行在目标行之前，它不包含目标，继续找下一个
				// 如果当前节点开始行在目标行之后，可以跳过（仅当父节点是有序语句块时）
				// 这里简单处理：不做过早退出的优化，保证正确性
				_ = endLine
				_ = targetLine
			}
		}
		return true
	})

	if foundType != nil {
		// 使用 definition.go 中的 getClassNameFromType
		return getClassNameFromType(foundType)
	}

	return ""
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

				detail := fmt.Sprintf("%s method in %s", visibility, currentClassName)

				insertText := label + "(${1})"
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
