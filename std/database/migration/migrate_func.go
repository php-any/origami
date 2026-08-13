package migration

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewMigrateFunction 创建 migrate 函数
func NewMigrateFunction() data.FuncStmt {
	return &MigrateFunction{}
}

// MigrateFunction Database\migrate($connection, $modelDir)
type MigrateFunction struct{}

func (f *MigrateFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	connection, _ := ctx.GetIndexValue(0)
	modelDirValue, ok := ctx.GetIndexValue(1)
	if !ok {
		return nil, nil
	}
	modelDir := modelDirValue.AsString()
	return Migrate(ctx, connection, modelDir)
}

func (f *MigrateFunction) GetName() string {
	return "Database\\migrate"
}

func (f *MigrateFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "connection", 0, data.NewNullValue(), data.NewBaseType("object")),
		node.NewParameter(nil, "modelDir", 1, nil, data.NewBaseType("string")),
	}
}

func (f *MigrateFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "connection", 0, data.NewBaseType("object")),
		node.NewVariable(nil, "modelDir", 1, data.NewBaseType("string")),
	}
}

func (f *MigrateFunction) GetReturnType() data.Types {
	return data.Class{Name: "Database\\Migration\\MigrateResult"}
}
