package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

// DbModelMethod 实现 DB::model(Model::class)，用于运行时类名或无法使用泛型时的构建器工厂。
// 静态类型场景优先使用 DB<Model>()。
type DbModelMethod struct{}

func (d *DbModelMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
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

func (d *DbModelMethod) GetName() string            { return "model" }
func (d *DbModelMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (d *DbModelMethod) GetIsStatic() bool          { return true }
func (d *DbModelMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("className", 0),
		data.NewParameterDefault("connectionName", 1, data.NewNullValue(), nil),
	}
}
func (d *DbModelMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("className", 0, data.NewBaseType("string")),
		data.NewVariable("connectionName", 1, data.NewBaseType("string")),
	}
}
func (d *DbModelMethod) GetReturnType() data.Types {
	return data.NewBaseType("Database\\DB")
}
