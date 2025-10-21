package database

import (
	"errors"

	"github.com/php-any/origami/data"
)

type DbGroupByMethod struct {
	source *db
}

func (d *DbGroupByMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取分组参数
	if groupBy, ok := ctx.GetIndexValue(0); ok {
		if groupByStr, ok := groupBy.(data.AsString); ok {
			d.source.setGroupBy(groupByStr.AsString())
			return ctx.(*data.ClassMethodContext).ClassValue, nil
		}
	}
	return nil, data.NewErrorThrow(nil, errors.New("分组字段必须是字符串"))
}

func (d *DbGroupByMethod) GetName() string {
	return "groupBy"
}

func (d *DbGroupByMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbGroupByMethod) GetIsStatic() bool {
	return false
}

func (d *DbGroupByMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("groupBy", 0),
	}
}

func (d *DbGroupByMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("groupBy", 0, data.NewBaseType("string")),
	}
}

func (d *DbGroupByMethod) GetReturnType() data.Types {
	return data.NewBaseType("database\\DB")
}
