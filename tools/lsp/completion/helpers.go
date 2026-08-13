package completion

import (
	"github.com/php-any/origami/data"
	"github.com/sirupsen/logrus"
)

// extractClassNamesFromType 从类型对象中提取所有可能的类名
func extractClassNamesFromType(typ data.Types) []string {
	if typ == nil {
		return nil
	}

	logrus.Debugf("extractClassNamesFromType: 处理类型 %T, String()=%s", typ, typ.String())

	var classNames []string

	switch t := typ.(type) {
	case *data.LspTypes:
		// LspTypes 包含多个类型，遍历所有内部类型
		logrus.Debugf("LspTypes 包含 %d 个类型", len(t.Types))
		for i, innerType := range t.Types {
			logrus.Debugf("  处理第 %d 个类型: %T", i, innerType)
			names := extractClassNamesFromType(innerType)
			classNames = append(classNames, names...)
		}
	case data.Class:
		logrus.Debugf("找到 Class 类型: %s", t.Name)
		classNames = append(classNames, t.Name)
	case data.NullableType:
		// 可空类型，递归获取基础类型
		logrus.Debugf("NullableType，递归处理基础类型")
		names := extractClassNamesFromType(t.BaseType)
		classNames = append(classNames, names...)
	default:
		// 对于其他类型，尝试使用 String() 方法
		typeStr := typ.String()
		logrus.Debugf("其他类型，String()=%s", typeStr)
		// 如果不是基础类型，可能是类名
		if !data.ISBaseType(typeStr) {
			logrus.Debugf("非基础类型，作为类名: %s", typeStr)
			classNames = append(classNames, typeStr)
		}
	}

	logrus.Debugf("extractClassNamesFromType 返回: %v", classNames)
	return classNames
}
