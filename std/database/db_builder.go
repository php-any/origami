package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

// newBuilderValue 根据模型类名创建 DB 查询构建器实例（model / DB<Model>() 的内部实现）。
func newBuilderValue(ctx data.Context, className, connName string) (*data.ClassValue, data.Control) {
	if className == "" {
		return nil, utils.NewThrow(errors.New("模型类名不能为空"))
	}

	vm := ctx.GetVM()
	classStmt, acl := vm.GetOrLoadClass(className)
	if acl != nil {
		return nil, acl
	}
	if classStmt == nil {
		return nil, utils.NewThrowf("找不到类 %s", className)
	}

	source := &db{
		model: data.Class{Name: classStmt.GetName()},
	}
	if connName != "" {
		source.setConnectionName(connName)
	}

	newDBClass := (&DBClass{}).CloneWithSource(source)
	return data.NewClassValue(newDBClass, vm.CreateContext([]data.Variable{})), nil
}

// classNameFromEntity 从实体实例或类名字符串解析模型类名。
func classNameFromEntity(val data.Value) (string, data.Control) {
	if val == nil {
		return "", utils.NewThrow(errors.New("缺少实体数据"))
	}
	if cv, ok := val.(*data.ClassValue); ok {
		if cv.Class != nil && cv.Class.GetName() != "" {
			return cv.Class.GetName(), nil
		}
		return "", utils.NewThrow(errors.New("无法从实体实例获取类名"))
	}
	if s, ok := val.(data.AsString); ok {
		name := s.AsString()
		if name != "" {
			return name, nil
		}
	}
	return "", utils.NewThrow(errors.New("插入/更新数据必须是实体对象"))
}

// callBuilderMethod 在指定模型的构建器上调用实例方法。
func callBuilderMethod(ctx data.Context, className, connName, methodName string, args []data.Value) (data.GetValue, data.Control) {
	builder, acl := newBuilderValue(ctx, className, connName)
	if acl != nil {
		return nil, acl
	}

	m, ok := builder.GetMethod(methodName)
	if !ok {
		return nil, utils.NewThrowf("DB 构建器不存在方法 %s", methodName)
	}

	vars := m.GetVariables()
	fnCtx := builder.CreateContext(vars)
	for i, arg := range args {
		if i < len(vars) {
			if ctl := fnCtx.SetVariableValue(vars[i], arg); ctl != nil {
				return nil, ctl
			}
		}
	}
	return m.Call(fnCtx)
}
