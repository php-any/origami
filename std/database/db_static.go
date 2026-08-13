package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// DbStaticInsertMethod 实现 DB::insert($entity)，从实体自动推断模型并插入。
type DbStaticInsertMethod struct{}

func (d *DbStaticInsertMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	entity, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少插入数据"))
	}
	val, ok := entity.(data.Value)
	if !ok {
		return nil, utils.NewThrow(errors.New("插入数据必须是对象或类实例"))
	}

	className, acl := classNameFromEntity(val)
	if acl != nil {
		return nil, acl
	}

	connName := ""
	if conn, ok := ctx.GetIndexValue(1); ok {
		if connStr, ok := conn.(data.AsString); ok {
			connName = connStr.AsString()
		}
	}

	return callBuilderMethod(ctx, className, connName, "insert", []data.Value{val})
}

func (d *DbStaticInsertMethod) GetName() string            { return "insert" }
func (d *DbStaticInsertMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (d *DbStaticInsertMethod) GetIsStatic() bool          { return true }
func (d *DbStaticInsertMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("entity", 0),
		data.NewParameterDefault("connectionName", 1, data.NewNullValue(), nil),
	}
}
func (d *DbStaticInsertMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("entity", 0, data.NewBaseType("object")),
		data.NewVariable("connectionName", 1, data.NewBaseType("string")),
	}
}
func (d *DbStaticInsertMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}

// DbStaticQueryMethod 实现 DB::query($sql, ...$args)，执行原生 SELECT 并返回行对象数组。
type DbStaticQueryMethod struct{}

func (d *DbStaticQueryMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	sqlStr, goArgs, ctl := parseSQLCallArgs(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return (&db{}).runQuery(ctx, sqlStr, goArgs)
}

func (d *DbStaticQueryMethod) GetName() string            { return "query" }
func (d *DbStaticQueryMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (d *DbStaticQueryMethod) GetIsStatic() bool          { return true }
func (d *DbStaticQueryMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("sql", 0),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}
func (d *DbStaticQueryMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("sql", 0, data.NewBaseType("string")),
		data.NewVariable("args", 1, data.NewBaseType("array")),
	}
}
func (d *DbStaticQueryMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}

// DbStaticExecuteMethod 实现 DB::execute($sql, ...$args)，执行原生 INSERT/UPDATE/DELETE/DDL。
type DbStaticExecuteMethod struct{}

func (d *DbStaticExecuteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	sqlStr, goArgs, ctl := parseSQLCallArgs(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return (&db{}).runExecute(ctx, sqlStr, goArgs)
}

func (d *DbStaticExecuteMethod) GetName() string            { return "execute" }
func (d *DbStaticExecuteMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (d *DbStaticExecuteMethod) GetIsStatic() bool          { return true }
func (d *DbStaticExecuteMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("sql", 0),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}
func (d *DbStaticExecuteMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("sql", 0, data.NewBaseType("string")),
		data.NewVariable("args", 1, data.NewBaseType("array")),
	}
}
func (d *DbStaticExecuteMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}
