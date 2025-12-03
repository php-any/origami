package main

import (
	"context"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/tools/lsp/defines"
	"github.com/sirupsen/logrus"
)

// handleCallObjectMethodDefinition 处理 *node.CallObjectMethod 的定义跳转
func handleCallObjectMethodDefinition(ctx *LspContext, n *node.CallObjectMethod) []*defines.Location {
	if n == nil {
		return nil
	}

	logrus.Debugf("CallObjectMethod: object=%T, method=%s", n.Object, n.Method)

	// 检查是否是链式调用（对象是另一个方法调用的结果）
	if chainCall, ok := n.Object.(*node.CallObjectMethod); ok {
		logrus.Debug("检测到链式调用，递归解析")
		if loc := findChainedMethodDefinition(ctx, chainCall, n.Method); loc != nil {
			return []*defines.Location{loc}
		}
		return nil
	}

	// 普通对象方法调用
	return findObjectMethodDefinition(ctx, n.Object, n.Method)
}

// findObjectMethodDefinition 查找对象方法定义
func findObjectMethodDefinition(ctx *LspContext, object data.GetValue, methodName string) []*defines.Location {
	if globalLspVM == nil {
		return nil
	}

	// 1. 对象是变量：从类型信息或上下文中推断所属类
	if varExpr, ok := object.(*node.VariableExpression); ok {
		// 1.1 优先使用变量节点上已有的类型信息
		if varExpr.Type != nil {
			switch t := varExpr.Type.(type) {
			case *data.LspTypes:
				var ret []*defines.Location
				for _, tt := range t.Types {
					if className := getClassNameFromType(tt); className != "" {
						if class, exists := globalLspVM.GetClass(className); exists {
							if methodLocation := findMethodInClass(class, methodName); methodLocation != nil {
								ret = append(ret, methodLocation)
							}
						}
					}
				}
				if len(ret) > 0 {
					return ret
				}
			default:
				if className := getClassNameFromType(varExpr.Type); className != "" {
					if class, exists := globalLspVM.GetClass(className); exists {
						if methodLocation := findMethodInClass(class, methodName); methodLocation != nil {
							return []*defines.Location{methodLocation}
						}
					}
				}
			}
		}

		// 1.2 再从上下文里根据变量名拿类型
		if ctx != nil {
			varType := ctx.GetVariableType(varExpr.Name)
			if varType != nil {
				logrus.Debugf("从上下文找到变量类型：%s -> %v", varExpr.Name, varType)
				if className := getClassNameFromType(varType); className != "" {
					logrus.Debugf("提取类名：%s", className)
					if class, exists := globalLspVM.GetClass(className); exists {
						if methodLocation := findMethodInClass(class, methodName); methodLocation != nil {
							return []*defines.Location{methodLocation}
						}
					}
				}
			} else {
				logrus.Debugf("在上下文中未找到变量类型：%s", varExpr.Name)
			}
		}
	}

	// 2. 对象是 this：查找当前类的方法
	if _, ok := object.(*node.This); ok {
		if ctx != nil {
			currentClassName := getCurrentClassNameFromContext(ctx)
			if currentClassName != "" {
				if class, exists := globalLspVM.GetClass(currentClassName); exists {
					if methodLocation := findMethodInClass(class, methodName); methodLocation != nil {
						return []*defines.Location{methodLocation}
					}
				}
			}
		}
	}

	// 3. 兜底：在所有类里按照方法名扫描一遍，只收集“真正的方法定义”，不退回类定义
	var candidates []*defines.Location
	if globalLspVM != nil {
		allClasses := globalLspVM.GetAllClasses()
		for className, classStmt := range allClasses {
			if methodLocation := findMethodInClassOnlyMethod(classStmt, methodName); methodLocation != nil {
				logrus.Debugf("在类 %s 中找到方法 %s（兜底扫描候选）", className, methodName)
				candidates = append(candidates, methodLocation)
			}
		}
	}
	// 若全局只有一个候选，就跳到它；多个候选则全部返回，让客户端决定如何展示/选择
	if len(candidates) == 1 {
		return candidates
	}
	if len(candidates) > 1 {
		return candidates
	}

	// 4. 最终兜底：尝试同名全局函数（保持与旧逻辑兼容）
	if function, exists := globalLspVM.GetFunc(methodName); exists {
		return []*defines.Location{createLocationFromFunction(function)}
	}

	return nil
}

// findChainedMethodDefinition 查找链式方法调用中的方法定义
func findChainedMethodDefinition(ctx *LspContext, chainCall *node.CallObjectMethod, methodName string) *defines.Location {
	if globalLspVM == nil {
		return nil
	}

	logrus.Debugf("解析链式调用：%s", methodName)

	object := chainCall.Object
	method := chainCall.Method

	logrus.Debugf("链式调用对象：%T，方法：%s", object, method)

	return resolveChainedMethod(ctx, object, method, methodName)
}

// resolveChainedMethod 递归解析链式方法调用
func resolveChainedMethod(ctx *LspContext, object data.GetValue, currentMethod, targetMethod string) *defines.Location {
	if globalLspVM == nil {
		return nil
	}

	logrus.Debugf("解析方法调用：%s，目标方法：%s", currentMethod, targetMethod)

	// 对象是变量：先根据变量类型推断所属类
	if varExpr, ok := object.(*node.VariableExpression); ok {
		logrus.Debugf("变量：%s", varExpr.Name)

		if ctx != nil {
			varType := ctx.GetVariableType(varExpr.Name)
			if varType != nil {
				logrus.Debugf("变量类型：%v", varType)

				if className := getClassNameFromType(varType); className != "" {
					logrus.Debugf("类名：%s", className)

					if class, exists := globalLspVM.GetClass(className); exists {
						if methodLocation := findMethodInClass(class, currentMethod); methodLocation != nil {
							logrus.Debugf("找到方法：%s", currentMethod)

							// 推断返回类型，并基于返回类型继续解析 targetMethod
							returnType := inferMethodReturnType(class, currentMethod)
							if returnType != nil {
								logrus.Debugf("推断返回类型：%s", returnType)

								switch rt := returnType.(type) {
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
								default:
									// 其他返回类型（string/int 等）统一退化为同名函数
									if function, exists := globalLspVM.GetFunc(targetMethod); exists {
										return createLocationFromFunction(function)
									}
								}
							}

							// 如果无法从返回类型中解析，兜底在所有类中扫描一遍目标方法（只返回真正的方法定义）
							allClasses := globalLspVM.GetAllClasses()
							for className, classStmt := range allClasses {
								logrus.Debugf("在所有类中查找：%s", className)
								if methodLocation := findMethodInClassOnlyMethod(classStmt, targetMethod); methodLocation != nil {
									logrus.Debugf("在类 %s 中找到目标方法：%s", className, targetMethod)
									return methodLocation
								}
							}
							logrus.Debugf("在所有类中都找不到目标方法：%s", targetMethod)
						}
					}
				}
			}
		}
	}

	// 对象本身又是一个 CallObjectMethod：继续向里递归
	if nestedCall, ok := object.(*node.CallObjectMethod); ok {
		logrus.Debug("检测到嵌套方法调用，递归解析")
		return resolveChainedMethod(ctx, nestedCall.Object, nestedCall.Method, targetMethod)
	}

	// 兜底：同名函数
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

	// 1. 方法上有返回类型注解，直接用
	if ret, ok := method.(data.GetReturnType); ok {
		ret := ret.GetReturnType()
		if ret != nil {
			return ret
		}
	}

	var inferredType data.Types

	// 2. 遍历方法体中的 return 语句，尝试从表达式推断
	if m, ok := method.(*node.ClassMethod); ok {
		docu := &DocumentInfo{}
		baseCtx := context.Background()
		lspCtx := NewLspContext(baseCtx, nil)

		for _, stmt := range m.Body {
			docu.foreachNode(lspCtx, stmt, nil, func(ctx *LspContext, parent, child data.GetValue) bool {
				if st, ok := child.(*node.ReturnStatement); ok {
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

	return inferredType
}
