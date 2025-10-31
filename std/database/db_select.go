package database

import (
	"errors"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type DbSelectMethod struct {
	source *db
}

func (d *DbSelectMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取字段参数
	if fields, ok := ctx.GetIndexValue(0); ok {
		if fieldsStr, ok := fields.(data.AsString); ok {
			// 分割字段字符串
			fieldList := strings.Split(fieldsStr.AsString(), ",")
			// 去除空格
			for i, field := range fieldList {
				fieldList[i] = strings.TrimSpace(field)
			}

			// 创建新的 db 对象实例
			newDB := d.source.clone()
			newDB.setSelect(fieldList)

			// 创建新的 DBClass 实例
			newDBClass := (&DBClass{}).CloneWithSource(newDB)

			// 返回新的 DB 类实例
			return data.NewClassValue(newDBClass, ctx.GetVM().CreateContext([]data.Variable{})), nil
		}
	}
	return nil, utils.NewThrow(errors.New("字段必须是字符串"))
}

func (d *DbSelectMethod) GetName() string {
	return "select"
}

func (d *DbSelectMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbSelectMethod) GetIsStatic() bool {
	return false
}

func (d *DbSelectMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("fields", 0),
	}
}

func (d *DbSelectMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("fields", 0, data.NewBaseType("string")),
	}
}

func (d *DbSelectMethod) GetReturnType() data.Types {
	return data.NewBaseType("Database\\DB")
}
