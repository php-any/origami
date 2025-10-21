package database

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/token"
)

type DbConstructMethod struct {
	source *db
}

func (d *DbConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 检查是否有连接名称参数
	if connName, ok := ctx.GetIndexValue(0); ok {
		if connNameStr, ok := connName.(data.AsString); ok {
			d.source.setConnectionName(connNameStr.AsString())
		}
	}
	return nil, nil
}

func (d *DbConstructMethod) GetName() string {
	return token.ConstructName
}

func (d *DbConstructMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbConstructMethod) GetIsStatic() bool {
	return false
}

func (d *DbConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("connectionName", 0), // 可选的连接名称参数
	}
}

func (d *DbConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("connectionName", 0, data.NewBaseType("string")),
	}
}

func (d *DbConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("")
}
