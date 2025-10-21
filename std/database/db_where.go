package database

import (
	"errors"

	"github.com/php-any/origami/data"
)

type DbWhereMethod struct {
	source *db
}

func (d *DbWhereMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a1, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少 sql 参数"))
	}
	sql := a1.(data.AsString).AsString()
	d.source.where = sql

	a2, ok := ctx.GetIndexValue(1)
	if ok {
		if arr, ok := a2.(*data.ArrayValue); ok {
			d.source.whereArgs = append([]data.Value{}, arr.Value...)
		}
	}

	return ctx.(*data.ClassMethodContext).ClassValue, nil
}

func (d *DbWhereMethod) GetName() string {
	return "where"
}

func (d *DbWhereMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbWhereMethod) GetIsStatic() bool {
	return false
}

func (d *DbWhereMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("sql", 0),   // 接收字符串参数
		data.NewParameters("args", 1), // 接收剩余参数数组
	}
}

func (d *DbWhereMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("sql", 0, data.NewBaseType("string")),
		data.NewVariable("args", 1, data.NewBaseType("array")),
	}
}

// GetReturnType 返回方法返回类型
func (d *DbWhereMethod) GetReturnType() data.Types {
	return data.Generic{}
}
