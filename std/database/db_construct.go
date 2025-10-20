package database

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/token"
)

type DbConstructMethod struct {
	source *db
}

func (d *DbConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
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
	return []data.GetValue{}
}

func (d *DbConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{}
}

func (d *DbConstructMethod) GetReturnType() data.Types {
	return data.NewBaseType("")
}
