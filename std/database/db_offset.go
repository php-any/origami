package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type DbOffsetMethod struct {
	source *db
}

func (d *DbOffsetMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取偏移参数
	if offset, ok := ctx.GetIndexValue(0); ok {
		if offsetInt, ok := offset.(data.AsInt); ok {
			if offsetValue, _ := offsetInt.AsInt(); offsetValue >= 0 {
				// 创建新的 db 对象实例
				newDB := d.source.clone()
				newDB.setOffset(offsetValue)

				// 创建新的 DBClass 实例
				newDBClass := (&DBClass{}).CloneWithSource(newDB)

				// 返回新的 DB 类实例
				return data.NewClassValue(newDBClass, ctx.GetVM().CreateContext([]data.Variable{})), nil
			}
		}
	}
	return nil, utils.NewThrow(errors.New("偏移必须是非负整数"))
}

func (d *DbOffsetMethod) GetName() string {
	return "offset"
}

func (d *DbOffsetMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbOffsetMethod) GetIsStatic() bool {
	return false
}

func (d *DbOffsetMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("offset", 0),
	}
}

func (d *DbOffsetMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("offset", 0, data.NewBaseType("int")),
	}
}

func (d *DbOffsetMethod) GetReturnType() data.Types {
	return data.NewBaseType("Database\\DB")
}
