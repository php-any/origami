package database

import (
	"github.com/php-any/origami/data"
)

type DbFirstMethod struct {
	source *db
}

func (d *DbFirstMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	return nil, nil
}

func (d *DbFirstMethod) GetName() string {
	return "first"
}

func (d *DbFirstMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbFirstMethod) GetIsStatic() bool {
	return false
}

func (d *DbFirstMethod) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (d *DbFirstMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

// GetReturnType 返回方法返回类型
func (d *DbFirstMethod) GetReturnType() data.Types {
	return data.NewBaseType("<T>")
}
