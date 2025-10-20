package database

import (
	"github.com/php-any/origami/data"
)

type DbGetMethod struct {
	source *db
}

func (d *DbGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

func (d *DbGetMethod) GetName() string {
	return "get"
}

func (d *DbGetMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbGetMethod) GetIsStatic() bool {
	return false
}

func (d *DbGetMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (d *DbGetMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回方法返回类型
func (d *DbGetMethod) GetReturnType() data.Types {
	return data.NewBaseType("string")
}
