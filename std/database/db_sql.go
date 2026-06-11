package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// DbSqlMethod 实现 DB::model(Entity::class)->sql($sql, $params)。
// 在已绑定模型的构建器上执行原生 SQL，返回行对象；可用 DB::toEntity() 再映射为实体。
type DbSqlMethod struct {
	source *db
}

func (d *DbSqlMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	sqlStr, goArgs, ctl := parseSQLCallArgs(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return d.source.runSQL(ctx, sqlStr, goArgs)
}

func (d *DbSqlMethod) GetName() string            { return "sql" }
func (d *DbSqlMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (d *DbSqlMethod) GetIsStatic() bool          { return false }
func (d *DbSqlMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("sql", 0),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}
func (d *DbSqlMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("sql", 0, data.NewBaseType("string")),
		data.NewVariable("args", 1, data.NewBaseType("array")),
	}
}
func (d *DbSqlMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}

// DbStaticSqlMethod 实现 DB::sql($sql, ...$args)，执行原生 SQL 并返回行对象数组。
type DbStaticSqlMethod struct{}

func (d *DbStaticSqlMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	sqlStr, goArgs, ctl := parseSQLCallArgs(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return (&db{}).runSQL(ctx, sqlStr, goArgs)
}

func (d *DbStaticSqlMethod) GetName() string            { return "sql" }
func (d *DbStaticSqlMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (d *DbStaticSqlMethod) GetIsStatic() bool          { return true }
func (d *DbStaticSqlMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("sql", 0),
		node.NewParameters(nil, "args", 1, nil, nil),
	}
}
func (d *DbStaticSqlMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("sql", 0, data.NewBaseType("string")),
		data.NewVariable("args", 1, data.NewBaseType("array")),
	}
}
func (d *DbStaticSqlMethod) GetReturnType() data.Types {
	return data.NewBaseType("array")
}

func parseSQLCallArgs(ctx data.Context) (string, []interface{}, data.Control) {
	sqlValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return "", nil, utils.NewThrow(errors.New("缺少 SQL 语句"))
	}
	sqlStrValue, ok := sqlValue.(data.AsString)
	if !ok {
		return "", nil, utils.NewThrow(errors.New("SQL 语句必须是字符串"))
	}
	return sqlStrValue.AsString(), collectBindArgs(ctx, 1), nil
}

// runSQL 执行原生 SQL，返回行对象数组。
func (d *db) runSQL(ctx data.Context, sqlStr string, goArgs []interface{}) (data.GetValue, data.Control) {
	conn := d.getConnection()
	if conn == nil {
		return nil, utils.NewThrow(errors.New("数据库连接不可用"))
	}

	rows, err := conn.Query(sqlStr, goArgs...)
	if err != nil {
		return nil, utils.NewThrowf("执行 SQL 失败: %v", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, utils.NewThrowf("获取列信息失败: %v", err)
	}

	scanner := NewDatabaseScanner()
	var results []data.Value
	scanValues := make([]interface{}, len(columns))
	scanPtrs := make([]interface{}, len(columns))
	for i := range scanValues {
		scanPtrs[i] = &scanValues[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanPtrs...); err != nil {
			return nil, utils.NewThrowf("扫描行失败: %v", err)
		}
		rowObj := data.NewObjectValue()
		for i, col := range columns {
			rowObj.SetProperty(col, scanner.convertToValue(scanValues[i]))
		}
		results = append(results, rowObj)
	}
	if err := rows.Err(); err != nil {
		return nil, utils.NewThrowf("遍历行时出错: %v", err)
	}
	return data.NewArrayValue(results), nil
}
