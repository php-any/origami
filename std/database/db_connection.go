package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

// newConnectionBuilder 创建绑定指定连接名的 DB 构建器实例。
func newConnectionBuilder(ctx data.Context, source *db, connName string) (*data.ClassValue, data.Control) {
	newDB := source.clone()
	newDB.setConnectionName(connName)
	newDBClass := (&DBClass{}).CloneWithSource(newDB)
	return data.NewClassValue(newDBClass, ctx.GetVM().CreateContext([]data.Variable{})), nil
}

// DbConnectionMethod 实现 ->connection($name)，链式切换命名连接。
type DbConnectionMethod struct {
	source *db
}

func (d *DbConnectionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少连接名称参数"))
	}
	connName, ok := a0.(data.AsString)
	if !ok {
		return nil, utils.NewThrow(errors.New("连接名称必须是字符串"))
	}
	return newConnectionBuilder(ctx, d.source, connName.AsString())
}

func (d *DbConnectionMethod) GetName() string            { return "connection" }
func (d *DbConnectionMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (d *DbConnectionMethod) GetIsStatic() bool          { return false }
func (d *DbConnectionMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("connectionName", 0),
	}
}
func (d *DbConnectionMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("connectionName", 0, data.NewBaseType("string")),
	}
}
func (d *DbConnectionMethod) GetReturnType() data.Types {
	return data.Generic{}
}

// DbStaticConnectionMethod 实现 DB::connection($name)，获取指定连接的构建器（用于原生 SQL 或未绑定模型的操作）。
type DbStaticConnectionMethod struct{}

func (d *DbStaticConnectionMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	a0, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少连接名称参数"))
	}
	connName, ok := a0.(data.AsString)
	if !ok {
		return nil, utils.NewThrow(errors.New("连接名称必须是字符串"))
	}
	source := &db{}
	source.setConnectionName(connName.AsString())
	return newConnectionBuilder(ctx, source, connName.AsString())
}

func (d *DbStaticConnectionMethod) GetName() string            { return "connection" }
func (d *DbStaticConnectionMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (d *DbStaticConnectionMethod) GetIsStatic() bool          { return true }
func (d *DbStaticConnectionMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("connectionName", 0),
	}
}
func (d *DbStaticConnectionMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("connectionName", 0, data.NewBaseType("string")),
	}
}
func (d *DbStaticConnectionMethod) GetReturnType() data.Types {
	return data.NewBaseType("Database\\DB")
}
