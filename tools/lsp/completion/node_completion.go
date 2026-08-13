package completion

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/sirupsen/logrus"
)

// getObjectPropertyAndMethodCompletionsFromNode 基于左侧节点获取对象属性和方法补全
func getObjectPropertyAndMethodCompletionsFromNode(content string, position defines.Position, provider defines.SymbolProvider, docProvider DocumentInfoProvider, vmProvider VMProvider) []defines.CompletionItem {
	if provider == nil || docProvider == nil || vmProvider == nil {
		return nil
	}

	ast := docProvider.GetAST()
	if ast == nil {
		return nil
	}

	logrus.Debugf("基于节点获取补全，位置：(%d, %d)", position.Line, position.Character)

	// 由于无法在 completion 包中直接调用 DocumentInfo.Foreach（需要 LspContext 类型），
	// 这个函数暂时返回 nil，让代码回退到基于变量名的方式
	// 实际的基于节点的补全逻辑在 completion.go 中实现
	return nil
}

// isPositionInNode 检查位置是否在节点范围内（简化版本）
func isPositionInNode(stmt data.GetValue, position defines.Position) bool {
	if stmt == nil {
		return false
	}

	// 尝试获取节点的 From 信息
	var from data.From
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
	lspChar := int(position.Character)

	// 检查行号是否在范围内
	if lspLine < startLine || lspLine > endLine {
		return false
	}

	// 如果在起始行，检查字符位置是否在起始字符之后
	if lspLine == startLine {
		if lspLine == endLine {
			// 单行节点：字符位置必须在起始和结束字符之间
			return lspChar >= startChar && lspChar <= endChar
		} else {
			// 多行节点的起始行：字符位置必须在起始字符之后
			return lspChar >= startChar
		}
	}

	// 如果在结束行，检查字符位置是否在结束字符之前
	if lspLine == endLine {
		return lspChar <= endChar
	}

	// 如果在中间行，肯定在范围内
	return true
}

// isNodeSmaller 检查节点 a 是否比节点 b 更小（更精确）
func isNodeSmaller(a, b data.GetValue) bool {
	if b == nil {
		return true
	}
	if a == nil {
		return false
	}

	af := getFromOfNode(a)
	bf := getFromOfNode(b)
	if af == nil || bf == nil {
		return false
	}

	// 获取两个节点的范围
	slA, scA, elA, ecA := af.GetRange()
	slB, scB, elB, ecB := bf.GetRange()

	// 计算两个节点的范围大小（字符数）
	rangeA := (elA-slA+1)*1000 + (ecA - scA + 1)
	rangeB := (elB-slB+1)*1000 + (ecB - scB + 1)

	// 返回范围更小的节点
	return rangeB < rangeA
}

// getFromOfNode 提取节点的 From 信息
func getFromOfNode(v data.GetValue) data.From {
	if v == nil {
		return nil
	}
	if gf, ok := v.(node.GetFrom); ok {
		return gf.GetFrom()
	}
	return nil
}

// getTypeFromNode 从节点获取类型信息
func getTypeFromNode(n data.GetValue, ctx interface{}, provider defines.SymbolProvider, content string, position defines.Position) data.Types {
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

		// 如果变量节点没有类型信息，尝试从 provider 获取
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
		// 通过接口获取当前类名
		if lspCtx, ok := ctx.(interface {
			GetCurrentScope() string
		}); ok {
			scope := lspCtx.GetCurrentScope()
			logrus.Debugf("This 表达式的当前作用域：%s", scope)
			// 从作用域中提取类名（格式：class:ClassName）
			if strings.HasPrefix(scope, "class:") {
				className := strings.TrimPrefix(scope, "class:")
				logrus.Debugf("从作用域提取类名：%s", className)
				// 返回类的类型
				return data.NewBaseType(className)
			}
		}
		// 如果无法从作用域获取，尝试从 provider 获取
		// 这里可以进一步优化，暂时返回 nil
		logrus.Debugf("无法从上下文获取 This 表达式的类名")
	}

	// 如果节点是 new 表达式，返回类的类型
	if newExpr, ok := n.(*node.NewExpression); ok {
		return data.NewBaseType(newExpr.ClassName)
	}

	// 如果节点是方法调用，尝试推断返回类型
	if callMethod, ok := n.(*node.CallObjectMethod); ok {
		// 递归获取对象的类型
		objType := getTypeFromNode(callMethod.Object, ctx, provider, content, position)
		if objType != nil {
			// 这里可以进一步推断方法调用的返回类型
			// 暂时返回对象类型
			return objType
		}
	}

	return nil
}
