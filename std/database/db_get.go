package database

import (
	"database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type DbGetMethod struct {
	source *db
}

func (d *DbGetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取数据库连接
	conn := d.source.getConnection()
	if conn == nil {
		return nil, utils.NewThrow(errors.New("数据库连接不可用"))
	}

	// 构建查询语句（支持注解处理）
	query := d.source.buildQueryWithContext(ctx)

	// 转换参数类型
	args := make([]interface{}, len(d.source.whereArgs))
	for i, arg := range d.source.whereArgs {
		args[i] = ConvertValueToGoType(arg)
	}

	// 执行查询
	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, utils.NewThrowf("get查询失败: err(%v); sql(%s)", err, query)
	}
	defer rows.Close()

	// 创建结果数组
	var results []data.GetValue

	// 遍历所有行
	for rows.Next() {
		// 创建模型实例
		instance, ctl := d.createClassInstance(d.source.model.(data.Class), rows, ctx)
		if ctl != nil {
			return nil, ctl
		}
		results = append(results, instance)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.NewThrowf("遍历行时出错: %v", err)
	}

	// 将 []data.GetValue 转换为 []data.Value
	values := make([]data.Value, len(results))
	for i, result := range results {
		if val, ok := result.(data.Value); ok {
			values[i] = val
		} else {
			// 如果转换失败，创建一个包装器
			values[i] = data.NewAnyValue(result)
		}
	}
	return data.NewArrayValue(values), nil
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
	return data.NewBaseType("array")
}

// createClassInstance 创建 Class 类型的实例
func (d *DbGetMethod) createClassInstance(classType data.Class, rows *sql.Rows, ctx data.Context) (data.GetValue, data.Control) {
	// 获取 VM 实例来查找类定义
	vm := ctx.GetVM()
	if vm == nil {
		return data.NewNullValue(), nil
	}

	// 根据类名获取类定义
	classStmt, exists := vm.GetClass(classType.Name)
	if !exists {
		return data.NewNullValue(), nil
	}

	// 创建类实例
	instance := data.NewClassValue(classStmt, vm.CreateContext([]data.Variable{}))

	// 使用扫描器扫描数据库行并设置属性
	scanner := NewDatabaseScanner()
	ctl := scanner.ScanRowToInstance(instance, rows)
	if ctl != nil {
		return nil, ctl
	}

	return instance, nil
}
