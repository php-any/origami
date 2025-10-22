package database

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type DbInsertMethod struct {
	source *db
}

func (d *DbInsertMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取数据库连接
	conn := d.source.getConnection()
	if conn == nil {
		return nil, data.NewErrorThrow(nil, errors.New("数据库连接不可用"))
	}

	// 获取插入数据
	dataValue, ok := ctx.GetIndexValue(0)
	if !ok {
		return nil, data.NewErrorThrow(nil, errors.New("缺少插入数据"))
	}

	// 构建插入语句
	tableName := d.source.getTableNameWithContext(ctx)

	// 处理插入数据
	var insertData map[string]interface{}
	if classValue, ok := dataValue.(*data.ClassValue); ok {
		// 从类实例获取属性，并处理注解映射
		insertData = make(map[string]interface{})
		properties := classValue.GetProperties()
		classStmt := classValue.Class

		for name, value := range properties {
			if value != nil {
				// 检查值是否为空
				if !d.isNullValue(value) {
					// 获取数据库列名
					columnName := d.getColumnName(classStmt, name)
					insertData[columnName] = ConvertValueToGoType(value)
				}
			}
		}
	} else if objectValue, ok := dataValue.(*data.ObjectValue); ok {
		// 从对象获取属性
		insertData = make(map[string]interface{})
		properties := objectValue.GetProperties()
		for name, value := range properties {
			if value != nil {
				insertData[name] = ConvertValueToGoType(value)
			}
		}
	} else {
		return nil, data.NewErrorThrow(nil, errors.New("插入数据必须是对象或类实例"))
	}

	// 构建 SQL 语句
	columns := make([]string, 0, len(insertData))
	placeholders := make([]string, 0, len(insertData))
	values := make([]interface{}, 0, len(insertData))

	for column, value := range insertData {
		columns = append(columns, column)
		placeholders = append(placeholders, "?")
		values = append(values, value)
	}

	// 构建列名字符串
	columnStr := ""
	placeholderStr := ""
	for i, col := range columns {
		if i > 0 {
			columnStr += ", "
			placeholderStr += ", "
		}
		columnStr += col
		placeholderStr += "?"
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, columnStr, placeholderStr)

	// 执行插入
	result, err := conn.Exec(query, values...)
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("插入失败: %w", err))
	}

	// 获取插入的 ID
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("获取插入ID失败: %w", err))
	}

	// 获取影响的行数
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("获取影响行数失败: %w", err))
	}

	// 返回插入结果
	resultObj := data.NewObjectValue()
	resultObj.SetProperty("insertId", data.NewIntValue(int(lastInsertId)))
	resultObj.SetProperty("rowsAffected", data.NewIntValue(int(rowsAffected)))
	resultObj.SetProperty("success", data.NewBoolValue(true))

	return resultObj, nil
}

func (d *DbInsertMethod) GetName() string {
	return "insert"
}

func (d *DbInsertMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbInsertMethod) GetIsStatic() bool {
	return false
}

func (d *DbInsertMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("data", 0),
	}
}

func (d *DbInsertMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("data", 0, data.NewBaseType("object")),
	}
}

func (d *DbInsertMethod) GetReturnType() data.Types {
	return data.NewBaseType("object")
}

// getColumnName 获取数据库列名，支持注解映射
func (d *DbInsertMethod) getColumnName(classStmt data.ClassStmt, propertyName string) string {
	// 获取属性定义
	properties := classStmt.GetPropertyList()
	var property data.Property
	var exists bool
	for _, prop := range properties {
		if prop.GetName() == propertyName {
			property = prop
			exists = true
			break
		}
	}
	if !exists {
		return propertyName
	}

	// 检查是否有注解
	if classProperty, ok := property.(*node.ClassProperty); ok {
		// 遍历属性的注解列表
		for _, annotation := range classProperty.Annotations {
			if annotation == nil {
				continue
			}

			// 检查是否是 Column 注解
			if annotation.Class != nil {
				className := annotation.Class.GetName()
				if className == "database\\annotation\\Column" {
					// 获取注解实例的属性
					annotationProps := annotation.GetProperties()

					// 查找 name 属性（Column 注解的第一个参数）
					if nameValue, exists := annotationProps["name"]; exists {
						if nameStr, ok := nameValue.(data.AsString); ok {
							return nameStr.AsString()
						}
					}

					// 如果 name 属性不存在，尝试其他可能的属性名
					for propName, propValue := range annotationProps {
						if propName != "name" && propValue != nil {
							if nameStr, ok := propValue.(data.AsString); ok {
								// 如果找到字符串值，可能是列名
								return nameStr.AsString()
							}
						}
					}
				}
			}
		}
	}

	// 如果没有注解，使用属性名
	return propertyName
}

// isNullValue 检查值是否为空
func (d *DbInsertMethod) isNullValue(val data.Value) bool {
	if val == nil {
		return true
	}

	switch v := val.(type) {
	case *data.NullValue:
		return true
	case *data.StringValue:
		return v.Value == ""
	case *data.IntValue:
		return v.Value == 0
	case *data.FloatValue:
		return v.Value == 0.0
	case *data.BoolValue:
		return !v.Value
	default:
		return false
	}
}
