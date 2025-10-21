package database

import (
	"errors"

	"github.com/php-any/origami/data"
)

type DbOrderByMethod struct {
	source *db
}

func (d *DbOrderByMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取排序参数
	if orderBy, ok := ctx.GetIndexValue(0); ok {
		if orderByStr, ok := orderBy.(data.AsString); ok {
			d.source.setOrderBy(orderByStr.AsString())
			return ctx.(*data.ClassMethodContext).ClassValue, nil
		}
	}
	return nil, data.NewErrorThrow(nil, errors.New("排序字段必须是字符串"))
}

func (d *DbOrderByMethod) GetName() string {
	return "orderBy"
}

func (d *DbOrderByMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbOrderByMethod) GetIsStatic() bool {
	return false
}

func (d *DbOrderByMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("orderBy", 0),
	}
}

func (d *DbOrderByMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("orderBy", 0, data.NewBaseType("string")),
	}
}

func (d *DbOrderByMethod) GetReturnType() data.Types {
	return data.NewBaseType("database\\DB")
}
