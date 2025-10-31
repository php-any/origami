package database

import (
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
)

type DbJoinMethod struct {
	source *db
}

func (d *DbJoinMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取连接参数
	if join, ok := ctx.GetIndexValue(0); ok {
		if joinStr, ok := join.(data.AsString); ok {
			// 创建新的 db 对象实例
			newDB := d.source.clone()
			newDB.addJoin(joinStr.AsString())

			// 创建新的 DBClass 实例
			newDBClass := (&DBClass{}).CloneWithSource(newDB)

			// 返回新的 DB 类实例
			return data.NewClassValue(newDBClass, ctx.GetVM().CreateContext([]data.Variable{})), nil
		}
	}
	return nil, utils.NewThrow(errors.New("连接语句必须是字符串"))
}

func (d *DbJoinMethod) GetName() string {
	return "join"
}

func (d *DbJoinMethod) GetModifier() data.Modifier {
	return data.ModifierPublic
}

func (d *DbJoinMethod) GetIsStatic() bool {
	return false
}

func (d *DbJoinMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		data.NewParameter("join", 0),
	}
}

func (d *DbJoinMethod) GetVariables() []data.Variable {
	return []data.Variable{
		data.NewVariable("join", 0, data.NewBaseType("string")),
	}
}

func (d *DbJoinMethod) GetReturnType() data.Types {
	return data.NewBaseType("Database\\DB")
}
