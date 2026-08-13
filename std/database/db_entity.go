package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

// DbToEntityMethod 实现 DB::toEntity(Entity::class, $rows)：将查询行对象映射为实体。
type DbToEntityMethod struct{}

func (d *DbToEntityMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少实体类名参数"))
	}
	className := a0.(data.AsString).AsString()
	if className == "" {
		return nil, utils.NewThrow(errors.New("实体类名不能为空"))
	}

	dataVal, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少行数据"))
	}

	vm := ctx.GetVM()
	classStmt, acl := vm.GetOrLoadClass(className)
	if acl != nil {
		return nil, acl
	}
	if classStmt == nil {
		return nil, utils.NewThrowf("找不到类 %s", className)
	}

	scanner := NewDatabaseScanner()

	if arr, ok := dataVal.(*data.ArrayValue); ok {
		instances := make([]data.Value, 0, len(arr.List))
		for _, z := range arr.List {
			row, ok := z.Value.(*data.ObjectValue)
			if !ok {
				if cv, ok := z.Value.(*data.ClassValue); ok {
					return data.NewArrayValue([]data.Value{cv}), nil
				}
				return nil, utils.NewThrow(errors.New("行数据必须是对象"))
			}
			inst, acl := scanner.MapObjectToInstance(row, classStmt, ctx)
			if acl != nil {
				return nil, acl
			}
			instances = append(instances, inst)
		}
		return data.NewArrayValue(instances), nil
	}

	if row, ok := dataVal.(*data.ObjectValue); ok {
		return scanner.MapObjectToInstance(row, classStmt, ctx)
	}
	if cv, ok := dataVal.(*data.ClassValue); ok {
		return cv, nil
	}

	return nil, utils.NewThrow(errors.New("行数据必须是对象或对象数组"))
}

func (d *DbToEntityMethod) GetName() string            { return "toEntity" }
func (d *DbToEntityMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (d *DbToEntityMethod) GetIsStatic() bool          { return true }
func (d *DbToEntityMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("className", 0),
		data.NewParameter("rows", 1),
	}
}
func (d *DbToEntityMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("className", 0, data.NewBaseType("string")),
		data.NewVariable("rows", 1, data.NewBaseType("array")),
	}
}
func (d *DbToEntityMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}
