package database

import (
	"errors"

	"github.com/php-any/origami/data"
)

type DbTableMethod struct {
	source *db
}

func (d *DbTableMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取表名参数
	if tableName, ok := ctx.GetIndexValue(0); ok {
		if tableNameStr, ok := tableName.(data.AsString); ok {
			d.source.setTableName(tableNameStr.AsString())
			return ctx.(*data.ClassMethodContext).ClassValue, nil
		}
	}
	return nil, data.NewErrorThrow(nil, errors.New("表名必须是字符串"))
}

func (d *DbTableMethod) GetName() string {
	return "table"
}

func (d *DbTableMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbTableMethod) GetIsStatic() bool {
	return false
}

func (d *DbTableMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("tableName", 0),
	}
}

func (d *DbTableMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("tableName", 0, data.NewBaseType("string")),
	}
}

func (d *DbTableMethod) GetReturnType() data.Types {
	return data.NewBaseType("database\\DB")
}
