package node

import (
	"github.com/php-any/origami/data"
)

// LikeExpression 表示 like 表达式
type LikeExpression struct {
	*Node
	Object    data.GetValue // 对象表达式
	ClassName string        // 类名
}

// NewLikeExpression 创建一个新的 like 表达式
func NewLikeExpression(from data.From, object data.GetValue, className string) *LikeExpression {
	return &LikeExpression{
		Node:      NewNode(from),
		Object:    object,
		ClassName: className,
	}
}

// GetValue 获取 like 表达式的值
func (l *LikeExpression) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 计算对象表达式的值
	objectValue, c := l.Object.GetValue(ctx)
	if c != nil {
		return nil, c
	}

	// 检查对象值是否为类实例
	if classValue, ok := objectValue.(*data.ClassValue); ok {
		// 检查目标类或接口是否存在
		vm := ctx.GetVM()

		// 先检查是否是类
		if targetClass, ok := vm.GetClass(l.ClassName); ok {
			result := checkClassStructure(classValue.Class, targetClass)
			return data.NewBoolValue(result), nil
		}

		// 再检查是否是接口
		if targetInterface, ok := vm.GetInterface(l.ClassName); ok {
			result := checkInterfaceStructure(classValue.Class, targetInterface)
			return data.NewBoolValue(result), nil
		}
	}

	// 如果不是类实例，返回 false
	return data.NewBoolValue(false), nil
}

// checkClassStructure 检查源类是否具有目标类的所有方法（结构化检查）
func checkClassStructure(source data.ClassStmt, target data.ClassStmt) bool {
	// 获取目标类的所有方法
	targetMethods := target.GetMethods()

	// 检查源类是否实现了目标类的所有方法
	for _, targetMethod := range targetMethods {
		methodName := targetMethod.GetName()
		sourceMethod, exists := source.GetMethod(methodName)

		if !exists {
			return false
		}

		// 检查方法签名是否匹配（参数数量）
		if len(sourceMethod.GetParams()) != len(targetMethod.GetParams()) {
			return false
		}
	}

	return true
}

// checkInterfaceStructure 检查源类是否具有目标接口的所有方法（结构化检查）
func checkInterfaceStructure(source data.ClassStmt, target data.InterfaceStmt) bool {
	// 获取目标接口的所有方法
	targetMethods := target.GetMethods()

	// 检查源类是否实现了目标接口的所有方法
	for _, targetMethod := range targetMethods {
		methodName := targetMethod.GetName()
		sourceMethod, exists := source.GetMethod(methodName)

		if !exists {
			return false
		}

		// 检查方法签名是否匹配（参数数量）
		if len(sourceMethod.GetParams()) != len(targetMethod.GetParams()) {
			return false
		}
	}

	return true
}
