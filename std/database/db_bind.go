package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

// DbBindMethod 实现 DB::bind(Model::class) 静态工厂，等价于 DB<Model>() 泛型语法。
type DbBindMethod struct{}

func (d *DbBindMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少模型类名参数"))
	}
	className := a0.(data.AsString).AsString()

	connName := ""
	if conn, ok := ctx.GetIndexValue(1); ok {
		if connStr, ok := conn.(data.AsString); ok {
			connName = connStr.AsString()
		}
	}

	return newBuilderValue(ctx, className, connName)
}

func (d *DbBindMethod) GetName() string {
	return "bind"
}

func (d *DbBindMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbBindMethod) GetIsStatic() bool {
	return true
}

func (d *DbBindMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("className", 0),
		data.NewParameterDefault("connectionName", 1, data.NewNullValue(), nil),
	}
}

func (d *DbBindMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("className", 0, data.NewBaseType("string")),
		data.NewVariable("connectionName", 1, data.NewBaseType("string")),
	}
}

func (d *DbBindMethod) GetReturnType() data.Types {
	return data.NewBaseType("Database\\DB")
}
