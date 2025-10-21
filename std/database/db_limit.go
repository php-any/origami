package database

import (
	"errors"

	"github.com/php-any/origami/data"
)

type DbLimitMethod struct {
	source *db
}

func (d *DbLimitMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取限制参数
	if limit, ok := ctx.GetIndexValue(0); ok {
		if limitInt, ok := limit.(data.AsInt); ok {
			if limitValue, _ := limitInt.AsInt(); limitValue > 0 {
				d.source.setLimit(limitValue)
				return ctx.(*data.ClassMethodContext).ClassValue, nil
			}
		}
	}
	return nil, data.NewErrorThrow(nil, errors.New("限制必须是正整数"))
}

func (d *DbLimitMethod) GetName() string {
	return "limit"
}

func (d *DbLimitMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbLimitMethod) GetIsStatic() bool {
	return false
}

func (d *DbLimitMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("limit", 0),
	}
}

func (d *DbLimitMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("limit", 0, data.NewBaseType("int")),
	}
}

func (d *DbLimitMethod) GetReturnType() data.Types {
	return data.NewBaseType("database\\DB")
}
