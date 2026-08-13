package migration

import (
	"fmt"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

const (
	annotationTable          = "Database\\Annotation\\Table"
	annotationColumn         = "Database\\Annotation\\Column"
	annotationID             = "Database\\Annotation\\Id"
	annotationGeneratedValue = "Database\\Annotation\\GeneratedValue"
)

// ColumnSchema 列定义
type ColumnSchema struct {
	Name          string
	PropertyName  string
	PHPType       string
	Length        int
	Nullable      bool
	PrimaryKey    bool
	AutoIncrement bool
}

// TableSchema 表定义
type TableSchema struct {
	TableName string
	ClassName string
	Columns   []ColumnSchema
}

func parseTableSchema(cls *node.ClassStatement) (*TableSchema, error) {
	tableName := tableNameFromClass(cls)
	if tableName == "" {
		return nil, fmt.Errorf("类 %s 缺少 #[Table] 注解", cls.GetName())
	}

	schema := &TableSchema{
		TableName: tableName,
		ClassName: cls.GetName(),
	}

	for _, name := range cls.PropertiesIndex {
		prop, ok := cls.Properties[name].(*node.ClassProperty)
		if !ok || prop == nil || prop.IsStatic {
			continue
		}

		col := columnSchemaFromProperty(prop)
		schema.Columns = append(schema.Columns, col)
	}

	if len(schema.Columns) == 0 {
		return nil, fmt.Errorf("类 %s 没有可映射的字段", cls.GetName())
	}
	return schema, nil
}

func tableNameFromClass(cls *node.ClassStatement) string {
	for _, ann := range cls.Annotations {
		if ann == nil || ann.Class == nil {
			continue
		}
		if ann.Class.GetName() != annotationTable {
			continue
		}
		props := ann.GetProperties()
		if nameValue, ok := props["name"]; ok {
			if nameStr, ok := nameValue.(data.AsString); ok {
				return nameStr.AsString()
			}
		}
		for _, propValue := range props {
			if propValue == nil {
				continue
			}
			if nameStr, ok := propValue.(data.AsString); ok {
				return nameStr.AsString()
			}
		}
	}
	return ""
}

func columnSchemaFromProperty(prop *node.ClassProperty) ColumnSchema {
	colName := prop.Name
	nullable := isNullableType(prop.Type)
	length := 255

	for _, ann := range prop.Annotations {
		if ann == nil || ann.Class == nil {
			continue
		}
		switch ann.Class.GetName() {
		case annotationColumn:
			props := ann.GetProperties()
			if nameValue, ok := props["name"]; ok {
				if nameStr, ok := nameValue.(data.AsString); ok && nameStr.AsString() != "" {
					colName = nameStr.AsString()
				}
			}
			if v, ok := props["nullable"]; ok {
				if bv, ok := v.(*data.BoolValue); ok {
					nullable, _ = bv.AsBool()
				}
			}
			if v, ok := props["length"]; ok {
				if iv, ok := v.(*data.IntValue); ok {
					length, _ = iv.AsInt()
				}
			}
		case annotationID:
			// handled below
		case annotationGeneratedValue:
			// handled below
		}
	}

	col := ColumnSchema{
		Name:         colName,
		PropertyName: prop.Name,
		PHPType:      phpTypeName(prop.Type),
		Length:       length,
		Nullable:     nullable,
	}

	for _, ann := range prop.Annotations {
		if ann == nil || ann.Class == nil {
			continue
		}
		if ann.Class.GetName() == annotationID {
			col.PrimaryKey = true
		}
		if ann.Class.GetName() == annotationGeneratedValue {
			props := ann.GetProperties()
			if v, ok := props["strategy"]; ok {
				if sv, ok := v.(data.AsString); ok && strings.EqualFold(sv.AsString(), "AUTO") {
					col.AutoIncrement = true
				}
			}
		}
	}

	if col.PrimaryKey {
		col.Nullable = false
	}
	return col
}

func isNullableType(ty data.Types) bool {
	if ty == nil {
		return true
	}
	if nt, ok := ty.(data.NullableType); ok {
		_ = nt
		return true
	}
	if strings.HasPrefix(ty.String(), "?") {
		return true
	}
	return false
}

func phpTypeName(ty data.Types) string {
	typeName := ""
	if ty != nil {
		typeName = ty.String()
	}
	typeName = strings.TrimPrefix(typeName, "?")
	if typeName == "" {
		return "mixed"
	}
	return strings.ToLower(typeName)
}
