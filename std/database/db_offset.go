package database

import (
	"errors"

	"github.com/php-any/origami/data"
)

type DbOffsetMethod struct {
	source *db
}

func (d *DbOffsetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取偏移参数
	if offset, ok := ctx.GetIndexValue(0); ok {
		if offsetInt, ok := offset.(data.AsInt); ok {
			if offsetValue, _ := offsetInt.AsInt(); offsetValue >= 0 {
				d.source.setOffset(offsetValue)
				return ctx.(*data.ClassMethodContext).ClassValue, nil
			}
		}
	}
	return nil, data.NewErrorThrow(nil, errors.New("偏移必须是非负整数"))
}

func (d *DbOffsetMethod) GetName() string {
	return "offset"
}

func (d *DbOffsetMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbOffsetMethod) GetIsStatic() bool {
	return false
}

func (d *DbOffsetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("offset", 0),
	}
}

func (d *DbOffsetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("offset", 0, data.NewBaseType("int")),
	}
}

func (d *DbOffsetMethod) GetReturnType() data.Types {
	return data.NewBaseType("database\\DB")
}
