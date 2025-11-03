package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

// DatabaseScanner 数据库扫描器，提供通用的数据库行扫描功能
type DatabaseScanner struct{}

// NewDatabaseScanner 创建数据库扫描器
func NewDatabaseScanner() *DatabaseScanner {
	return &DatabaseScanner{}
}

// ScanRowToInstance 扫描数据库行并设置实例属性
func (ds *DatabaseScanner) ScanRowToInstance(instance *data.ClassValue, rows *sql.Rows) data.Control {
	// 获取类的属性定义
	classStmt := instance.Class
	if classStmt == nil {
		return utils.NewThrow(errors.New("实例没有类定义"))
	}

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		return utils.NewThrowf("获取列信息失败: %v", err)
	}

	// 创建列名到索引的映射
	columnMap := make(map[string]int)
	for i, col := range columns {
		columnMap[col] = i
	}

	// 动态创建扫描目标
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	// 扫描数据
	err = rows.Scan(valuePtrs...)
	if err != nil {
		return utils.NewThrowf("扫描数据库行失败: %v", err)
	}

	// 根据列名映射到类属性
	for _, property := range classStmt.GetPropertyList() {
		propertyName := property.GetName()
		// 尝试不同的列名匹配策略
		var value data.Value

		// 1. 优先直接匹配属性名（性能优化：匹配成功就不需要检查注解）
		if columnIndex, exists := columnMap[propertyName]; exists {
			value = ds.convertToValue(values[columnIndex])
		} else {
			// 2. 直接匹配失败，检查注解 @Column("name")
			annotationColumnName := getColumnName(classStmt, propertyName)
			if annotationColumnName != propertyName {
				// 注解中的列名和属性名不同，使用注解中的列名
				if columnIndex, exists := columnMap[annotationColumnName]; exists {
					value = ds.convertToValue(values[columnIndex])
				} else {
					value = data.NewNullValue()
				}
			} else {
				// 3. 尝试下划线命名转换 (user_name -> userName)
				camelCaseName := ds.toCamelCase(propertyName)
				if columnIndex, exists := columnMap[camelCaseName]; exists {
					value = ds.convertToValue(values[columnIndex])
				} else {
					// 4. 尝试蛇形命名转换 (userName -> user_name)
					snakeCaseName := ds.toSnakeCase(propertyName)
					if columnIndex, exists := columnMap[snakeCaseName]; exists {
						value = ds.convertToValue(values[columnIndex])
					} else {
						// 5. 如果都不匹配，设置为 null
						value = data.NewNullValue()
					}
				}
			}
		}

		instance.SetProperty(propertyName, value)
	}

	return nil
}

// ScanRowsToInstances 扫描多行数据到实例数组
func (ds *DatabaseScanner) ScanRowsToInstances(rows *sql.Rows, classStmt data.ClassStmt, ctx data.Context) ([]*data.ClassValue, data.Control) {
	var instances []*data.ClassValue

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		return nil, utils.NewThrowf("获取列信息失败: %v", err)
	}

	// 创建列名到索引的映射
	columnMap := make(map[string]int)
	for i, col := range columns {
		columnMap[col] = i
	}

	// 创建扫描目标
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		// 创建新的实例
		instance := data.NewClassValue(classStmt, ctx)

		// 扫描当前行
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, utils.NewThrowf("扫描行失败: %v", err)
		}

		// 将扫描结果映射到实例属性
		for _, property := range classStmt.GetPropertyList() {
			propertyName := property.GetName()
			var value data.Value

			// 1. 直接匹配
			if columnIndex, exists := columnMap[propertyName]; exists {
				value = ds.convertToValue(values[columnIndex])
			} else {
				// 2. 尝试下划线命名转换
				camelCaseName := ds.toCamelCase(propertyName)
				if columnIndex, exists := columnMap[camelCaseName]; exists {
					value = ds.convertToValue(values[columnIndex])
				} else {
					// 3. 尝试蛇形命名转换
					snakeCaseName := ds.toSnakeCase(propertyName)
					if columnIndex, exists := columnMap[snakeCaseName]; exists {
						value = ds.convertToValue(values[columnIndex])
					} else {
						// 4. 如果都不匹配，设置为 null
						value = data.NewNullValue()
					}
				}
			}

			instance.SetProperty(propertyName, value)
		}

		instances = append(instances, instance)
	}

	// 检查是否有错误
	if err := rows.Err(); err != nil {
		return nil, utils.NewThrowf("遍历行时出错: %v", err)
	}

	return instances, nil
}

// convertToValue 将数据库值转换为脚本值
func (ds *DatabaseScanner) convertToValue(val interface{}) data.Value {
	if val == nil {
		return data.NewNullValue()
	}

	switch v := val.(type) {
	case int:
		return data.NewIntValue(v)
	case int8:
		return data.NewIntValue(int(v))
	case int16:
		return data.NewIntValue(int(v))
	case int32:
		return data.NewIntValue(int(v))
	case int64:
		// 检查是否超出 int 范围
		if v > int64(^uint(0)>>1) || v < int64(-1<<63) {
			// 如果超出 int 范围，转换为字符串
			return data.NewStringValue(fmt.Sprintf("%d", v))
		}
		return data.NewIntValue(int(v))
	case uint:
		if v > uint(^uint(0)>>1) {
			return data.NewStringValue(fmt.Sprintf("%d", v))
		}
		return data.NewIntValue(int(v))
	case uint8:
		return data.NewIntValue(int(v))
	case uint16:
		return data.NewIntValue(int(v))
	case uint32:
		return data.NewIntValue(int(v))
	case uint64:
		if v > uint64(^uint(0)>>1) {
			return data.NewStringValue(fmt.Sprintf("%d", v))
		}
		return data.NewIntValue(int(v))
	case float32:
		return data.NewFloatValue(float64(v))
	case float64:
		return data.NewFloatValue(v)
	case string:
		return data.NewStringValue(v)
	case []byte:
		return data.NewStringValue(string(v))
	case bool:
		return data.NewBoolValue(v)
	default:
		// 对于其他类型，转换为字符串
		return data.NewStringValue(fmt.Sprintf("%v", v))
	}
}

// toCamelCase 将下划线命名转换为驼峰命名
func (ds *DatabaseScanner) toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 1 {
		return s
	}

	result := parts[0]
	for _, part := range parts[1:] {
		if len(part) > 0 {
			result += strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return result
}

// toSnakeCase 将驼峰命名转换为下划线命名
func (ds *DatabaseScanner) toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
