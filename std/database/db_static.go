package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// DbModelMethod 实现 DB::model(Model::class)，语义更清晰的构建器工厂（等价于 bind / DB<Model>()）。
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

// DbStaticExecuteMethod 实现 DB::execute($sql, $params)，执行原生 INSERT/UPDATE/DELETE/DDL。
type DbStaticExecuteMethod struct{}

func (d *DbStaticExecuteMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	sqlVal, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, utils.NewThrow(errors.New("缺少 SQL 语句"))
	}
	sql, ok := sqlVal.(data.Value)
	if !ok {
		return nil, utils.NewThrow(errors.New("SQL 语句必须是字符串"))
	}

	args := []data.Value{sql}
	if paramVal, ok := ctx.GetIndexValue(1); ok {
		if _, isNull := paramVal.(*data.NullValue); !isNull {
			args = append(args, paramVal.(data.Value))
		}
	}

	return callDefaultBuilderMethod(ctx, "", "exec", args)
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
