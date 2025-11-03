package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func newDB(genericMap map[string]data.Types) *db {
	return &db{
		model: genericMap["M"],
	}
}

type db struct {
	connName string  // 连接名称
	conn     *sql.DB // 数据库连接

	// 查询条件
	where     string
	whereArgs []data.Value

	// 表名
	tableName string

	// 字段选择
	selectFields []string

	// 排序
	orderBy string

	// 分组
	groupBy string

	// 限制
	limit  int
	offset int

	// 连接
	joins []string

	// 泛型对应的具体类型
	model data.Types
}

// clone 创建 db 对象的深拷贝
func (d *db) clone() *db {
	// 深拷贝 whereArgs 切片
	whereArgsCopy := make([]data.Value, len(d.whereArgs))
	copy(whereArgsCopy, d.whereArgs)

	// 深拷贝 selectFields 切片
	selectFieldsCopy := make([]string, len(d.selectFields))
	copy(selectFieldsCopy, d.selectFields)

	// 深拷贝 joins 切片
	joinsCopy := make([]string, len(d.joins))
	copy(joinsCopy, d.joins)

	return &db{
		connName:     d.connName,
		conn:         d.conn,
		where:        d.where,
		whereArgs:    whereArgsCopy,
		tableName:    d.tableName,
		selectFields: selectFieldsCopy,
		orderBy:      d.orderBy,
		groupBy:      d.groupBy,
		limit:        d.limit,
		offset:       d.offset,
		joins:        joinsCopy,
		model:        d.model,
	}
}

// getConnection 获取数据库连接
func (d *db) getConnection() *sql.DB {
	if d.conn != nil {
		return d.conn
	}

	manager := GetConnectionManager()
	if d.connName != "" {
		if conn, exists := manager.GetConnection(d.connName); exists {
			d.conn = conn
			return conn
		}
	}

	// 如果没有指定连接名称或连接不存在，使用默认连接
	if conn, exists := manager.GetDefaultConnection(); exists {
		d.conn = conn
		return conn
	}

	return nil
}

// setConnectionName 设置连接名称
func (d *db) setConnectionName(name string) {
	d.connName = name
	d.conn = nil // 重置连接，下次调用时重新获取
}

// setTableName 设置表名
func (d *db) setTableName(name string) {
	d.tableName = name
}

// addWhere 添加查询条件
func (d *db) addWhere(condition string, args ...data.Value) {
	if d.where != "" {
		d.where += " AND " + condition
	} else {
		d.where = condition
	}
	d.whereArgs = append(d.whereArgs, args...)
}

// setSelect 设置选择字段
func (d *db) setSelect(fields []string) {
	d.selectFields = fields
}

// setOrderBy 设置排序
func (d *db) setOrderBy(order string) {
	d.orderBy = order
}

// setGroupBy 设置分组
func (d *db) setGroupBy(group string) {
	d.groupBy = group
}

// setLimit 设置限制
func (d *db) setLimit(limit int) {
	d.limit = limit
}

// setOffset 设置偏移
func (d *db) setOffset(offset int) {
	d.offset = offset
}

// addJoin 添加连接
func (d *db) addJoin(join string) {
	d.joins = append(d.joins, join)
}

// buildQuery 构建查询语句
func (d *db) buildQuery() string {
	query := "SELECT "

	// 选择字段
	if len(d.selectFields) > 0 {
		for i, field := range d.selectFields {
			if i > 0 {
				query += ", "
			}
			query += field
		}
	} else {
		query += "*"
	}

	// 表名
	tableName := d.getTableName()
	if tableName == "" {
		// 如果没有表名，抛出错误
		return ""
	}
	query += " FROM " + tableName

	// 连接
	for _, join := range d.joins {
		query += " " + join
	}

	// 条件
	if d.where != "" {
		query += " WHERE " + d.where
	}

	// 分组
	if d.groupBy != "" {
		query += " GROUP BY " + d.groupBy
	}

	// 排序
	if d.orderBy != "" {
		query += " ORDER BY " + d.orderBy
	}

	// 限制
	if d.limit > 0 {
		query += " LIMIT " + fmt.Sprintf("%d", d.limit)
		if d.offset > 0 {
			query += " OFFSET " + fmt.Sprintf("%d", d.offset)
		}
	}

	return query
}

// buildQueryWithContext 构建查询语句，支持注解处理
func (d *db) buildQueryWithContext(ctx data.Context) string {
	query := "SELECT "

	// 选择字段
	if len(d.selectFields) > 0 {
		for i, field := range d.selectFields {
			if i > 0 {
				query += ", "
			}
			query += field
		}
	} else {
		query += "*"
	}

	// 表名
	tableName := d.getTableNameWithContext(ctx)
	if tableName == "" {
		// 如果没有表名，抛出错误
		return ""
	}
	query += " FROM " + tableName

	// 连接
	for _, join := range d.joins {
		query += " " + join
	}

	// 条件
	if d.where != "" {
		query += " WHERE " + d.where
	}

	// 分组
	if d.groupBy != "" {
		query += " GROUP BY " + d.groupBy
	}

	// 排序
	if d.orderBy != "" {
		query += " ORDER BY " + d.orderBy
	}

	// 限制
	if d.limit > 0 {
		query += " LIMIT " + fmt.Sprintf("%d", d.limit)
		if d.offset > 0 {
			query += " OFFSET " + fmt.Sprintf("%d", d.offset)
		}
	}

	return query
}

// getTableName 获取表名，优先使用显式设置的表名，否则从注解中获取
func (d *db) getTableName() string {
	// 如果显式设置了表名，使用它
	if d.tableName != "" {
		return d.tableName
	}

	// 尝试从模型类型的注解中获取表名
	if d.model != nil {
		if classType, ok := d.model.(data.Class); ok {
			// 从类名推断表名
			className := classType.Name
			if className != "" {
				// 如果没有注解，使用类名转换
				return d.convertClassNameToTableName(className)
			}
		}
	}

	// 如果没有模型类型，返回空字符串
	return ""
}

// getTableNameWithContext 获取表名，支持注解处理（需要 Context）
func (d *db) getTableNameWithContext(ctx data.Context) string {
	// 如果显式设置了表名，使用它
	if d.tableName != "" {
		return d.tableName
	}

	// 尝试从模型类型的注解中获取表名
	if d.model != nil {
		if classType, ok := d.model.(data.Class); ok {
			// 从类名推断表名
			className := classType.Name
			if className != "" {
				// 首先尝试从注解获取表名
				tableName := d.getTableNameFromAnnotation(classType, ctx)
				if tableName != "" {
					return tableName
				}

				// 如果没有注解，使用类名转换
				return d.convertClassNameToTableName(className)
			}
		}
	}

	// 如果没有模型类型，返回空字符串
	return ""
}

// convertClassNameToTableName 将类名转换为表名
func (d *db) convertClassNameToTableName(className string) string {
	if className == "" {
		return ""
	}

	// 处理命名空间：App\User -> User
	// 只取最后一个反斜杠后的部分
	if lastSlash := strings.LastIndex(className, "\\"); lastSlash != -1 {
		className = className[lastSlash+1:]
	}

	// 直接返回类名作为表名
	// User -> User, UserInfo -> UserInfo
	return className
}

// getTableNameFromAnnotation 从类定义中获取 @Table 注解的表名
func (d *db) getTableNameFromAnnotation(classType data.Class, ctx data.Context) string {
	// 获取 VM 实例来查找类定义
	vm := ctx.GetVM()
	if vm == nil {
		return ""
	}

	// 根据类名获取类定义
	classStmt, acl := vm.GetOrLoadClass(classType.Name)
	if acl != nil || classStmt == nil {
		return ""
	}

	// 检查类是否有 @Table 注解
	// 需要将 data.ClassStmt 转换为 node.ClassStatement 来访问注解
	if classStatement, ok := classStmt.(*node.ClassStatement); ok {
		// 遍历类的注解列表
		for _, annotation := range classStatement.Annotations {
			if annotation == nil {
				continue
			}

			// 检查是否是 Table 注解
			if annotation.Class != nil {
				className := annotation.Class.GetName()
				if className == "Database\\Annotation\\Table" {
					// 获取注解实例的属性
					annotationProps := annotation.GetProperties()

					// 查找 name 属性（Table 注解的第一个参数）
					if nameValue, exists := annotationProps["name"]; exists {
						if nameStr, ok := nameValue.(data.AsString); ok {
							return nameStr.AsString()
						}
					}

					// 如果 name 属性不存在，尝试其他可能的属性名
					for propName, propValue := range annotationProps {
						if propName != "name" && propValue != nil {
							if nameStr, ok := propValue.(data.AsString); ok {
								// 如果找到字符串值，可能是表名
								return nameStr.AsString()
							}
						}
					}
				}
			}
		}
	}

	// 如果没有找到 @Table 注解，返回空字符串
	return ""
}

// getColumnName 获取数据库列名，支持注解映射（通用方法）
// 优先使用属性名，只有当注解中的列名和属性名不匹配时，才使用注解中的列名
func getColumnName(classStmt data.ClassStmt, propertyName string) string {
	// 直接通过属性名获取属性定义
	property, exists := classStmt.GetProperty(propertyName)
	if !exists {
		return propertyName
	}

	// 检查是否有注解，只有在需要时才检查注解
	classProperty, ok := property.(*node.ClassProperty)
	if !ok || len(classProperty.Annotations) == 0 {
		// 没有注解，直接使用属性名
		return propertyName
	}

	// 有注解，查找 Column 注解
	for _, annotation := range classProperty.Annotations {
		if annotation == nil {
			continue
		}

		// 检查是否是 Column 注解
		if annotation.Class != nil {
			className := annotation.Class.GetName()
			if className == "Database\\Annotation\\Column" {
				// 获取注解实例的属性
				annotationProps := annotation.GetProperties()

				// 查找 name 属性（Column 注解的第一个参数）
				if nameValue, exists := annotationProps["name"]; exists && nameValue != nil {
					if nameStr, ok := nameValue.(data.AsString); ok {
						annotationColumnName := nameStr.AsString()
						// 如果注解中的列名和属性名相同，直接返回属性名（避免不必要处理）
						if annotationColumnName == propertyName {
							return propertyName
						}
						// 只有注解中的列名和属性名不匹配时，才使用注解中的列名
						return annotationColumnName
					}
				}
			}
		}
	}

	// 没有找到有效的 Column 注解，直接使用属性名
	return propertyName
}
