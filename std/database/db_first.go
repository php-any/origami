package database

import (
	"database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type DbFirstMethod struct {
	source  *db
	scanner *DatabaseScanner
}

func (d *DbFirstMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取数据库连接
	conn := d.source.getConnection()
	if conn == nil {
		return nil, utils.NewThrow(errors.New("数据库连接不可用"))
	}

	// 构建查询语句（支持注解处理）
	query := d.source.buildQueryWithContext(ctx)
	if d.source.limit == 0 {
		query += " LIMIT 1"
	}

	// 转换参数类型
	args := make([]interface{}, len(d.source.whereArgs))
	for i, arg := range d.source.whereArgs {
		args[i] = ConvertValueToGoType(arg)
	}

	// 执行查询
	rows, err := conn.Query(query, args...)
	if err != nil {
		return nil, utils.NewThrowf("first查询失败: %v", err)
	}
	defer rows.Close()

	// 检查是否有结果
	if !rows.Next() {
		// 如果没有结果，返回 null
		return data.NewNullValue(), nil
	}

	// 根据模型类型扫描结果
	if d.source.model != nil {
		return d.createModelInstance(rows, ctx)
	}

	// 如果没有模型类型，返回空对象
	return data.NewObjectValue(), nil
}

// createModelInstance 根据模型类型创建实例
func (d *DbFirstMethod) createModelInstance(rows *sql.Rows, ctx data.Context) (data.GetValue, data.Control) {
	// 获取模型类型信息
	modelType := d.source.model

	if modelType == nil {
		return data.NewNullValue(), nil
	}

	// 根据模型类型创建实例
	switch modelType := modelType.(type) {
	case data.Class:
		// 处理 Class 类型
		return d.createClassInstance(modelType, rows, ctx)
	default:
		// 其他类型，返回空对象
		return data.NewObjectValue(), nil
	}
}

// createClassInstance 创建 Class 类型的实例
func (d *DbFirstMethod) createClassInstance(classType data.Class, rows *sql.Rows, ctx data.Context) (data.GetValue, data.Control) {
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
	if d.scanner == nil {
		d.scanner = NewDatabaseScanner()
	}

	if control := d.scanner.ScanRowToInstance(instance, rows); control != nil {
		return nil, control
	}

	return instance, nil
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
