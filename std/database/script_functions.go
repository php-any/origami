package database

import (
	"database/sql"
	"errors"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	sqlpkg "github.com/php-any/origami/std/database/sql"
	"github.com/php-any/origami/utils"
)

// NewRegisterConnectionFunction 创建注册连接函数
func NewRegisterConnectionFunction() data.FuncStmt {
	return &RegisterConnectionFunction{}
}

type RegisterConnectionFunction struct{}

func (f *RegisterConnectionFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取参数
	nameValue, _ := ctx.GetIndexValue(0)
	dbValue, _ := ctx.GetIndexValue(1)

	// 获取连接名称
	var name string
	if nameStr, ok := nameValue.(data.AsString); ok {
		name = nameStr.AsString()
	} else {
		return nil, utils.NewThrow(errors.New("连接名称必须是字符串"))
	}

	// 获取数据库连接对象
	manager := GetConnectionManager()
	if dbClass, ok := dbValue.(*data.ClassValue); ok {
		// 获取 ClassValue 的 Class，然后获取 source
		if dbClassStmt, ok := dbClass.Class.(*sqlpkg.DBClass); ok {
			if sqlDB, ok := dbClassStmt.GetSource().(*sql.DB); ok {
				manager.AddConnection(name, sqlDB)
				return data.NewBoolValue(true), nil
			}
		}
	}

	return nil, utils.NewThrow(errors.New("无效的数据库连接对象"))
}

func (f *RegisterConnectionFunction) GetName() string {
	return "Database\\registerConnection"
}

func (f *RegisterConnectionFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.NewBaseType("string")),
		node.NewParameter(nil, "db", 1, nil, data.NewBaseType("object")),
	}
}

func (f *RegisterConnectionFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "db", 1, data.NewBaseType("object")),
	}
}

// NewRegisterDefaultConnectionFunction 创建注册默认连接函数
func NewRegisterDefaultConnectionFunction() data.FuncStmt {
	return &RegisterDefaultConnectionFunction{}
}

type RegisterDefaultConnectionFunction struct{}

func (f *RegisterDefaultConnectionFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取数据库连接对象
	dbValue, _ := ctx.GetIndexValue(0)

	manager := GetConnectionManager()
	if dbClass, ok := dbValue.(*data.ClassValue); ok {
		// 获取 ClassValue 的 Class，然后获取 source
		if dbClassStmt, ok := dbClass.Class.(*sqlpkg.DBClass); ok {
			if sqlDB, ok := dbClassStmt.GetSource().(*sql.DB); ok {
				manager.AddConnection("default", sqlDB)
				return data.NewBoolValue(true), nil
			}
		}
	}

	return nil, utils.NewThrow(errors.New("无效的数据库连接对象"))
}

func (f *RegisterDefaultConnectionFunction) GetName() string {
	return "Database\\registerDefaultConnection"
}

func (f *RegisterDefaultConnectionFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "db", 0, nil, data.NewBaseType("object")),
	}
}

func (f *RegisterDefaultConnectionFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "db", 0, data.NewBaseType("object")),
	}
}

// NewGetConnectionFunction 创建获取连接函数
func NewGetConnectionFunction() data.FuncStmt {
	return &GetConnectionFunction{}
}

type GetConnectionFunction struct{}

func (f *GetConnectionFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取连接名称
	nameValue, _ := ctx.GetIndexValue(0)

	var name string
	if nameStr, ok := nameValue.(data.AsString); ok {
		name = nameStr.AsString()
	} else {
		return data.NewNullValue(), nil
	}

	manager := GetConnectionManager()
	conn, exists := manager.GetConnection(name)
	if !exists {
		return data.NewNullValue(), nil
	}

	// 将 *sql.DB 包装为 AnyValue 返回给脚本
	return data.NewAnyValue(conn), nil
}

func (f *GetConnectionFunction) GetName() string {
	return "Database\\getConnection"
}

func (f *GetConnectionFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.NewBaseType("string")),
	}
}

func (f *GetConnectionFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.NewBaseType("string")),
	}
}

// NewGetDefaultConnectionFunction 创建获取默认连接函数
func NewGetDefaultConnectionFunction() data.FuncStmt {
	return &GetDefaultConnectionFunction{}
}

type GetDefaultConnectionFunction struct{}

func (f *GetDefaultConnectionFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	manager := GetConnectionManager()
	conn, exists := manager.GetConnection("default")
	if !exists {
		return data.NewNullValue(), nil
	}

	// 将 *sql.DB 包装为 AnyValue 返回给脚本
	return data.NewAnyValue(conn), nil
}

func (f *GetDefaultConnectionFunction) GetName() string {
	return "Database\\getDefaultConnection"
}

func (f *GetDefaultConnectionFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (f *GetDefaultConnectionFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}

// NewRemoveConnectionFunction 创建移除连接函数
func NewRemoveConnectionFunction() data.FuncStmt {
	return &RemoveConnectionFunction{}
}

type RemoveConnectionFunction struct{}

func (f *RemoveConnectionFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取连接名称
	nameValue, _ := ctx.GetIndexValue(0)

	var name string
	if nameStr, ok := nameValue.(data.AsString); ok {
		name = nameStr.AsString()
	} else {
		return nil, utils.NewThrow(errors.New("连接名称必须是字符串"))
	}

	manager := GetConnectionManager()
	manager.RemoveConnection(name)
	return data.NewBoolValue(true), nil
}

func (f *RemoveConnectionFunction) GetName() string {
	return "Database\\removeConnection"
}

func (f *RemoveConnectionFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "name", 0, nil, data.NewBaseType("string")),
	}
}

func (f *RemoveConnectionFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "name", 0, data.NewBaseType("string")),
	}
}

// NewListConnectionsFunction 创建列出连接函数
func NewListConnectionsFunction() data.FuncStmt {
	return &ListConnectionsFunction{}
}

type ListConnectionsFunction struct{}

func (f *ListConnectionsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	manager := GetConnectionManager()
	names := manager.ListConnections()

	// 将字符串数组转换为 ArrayValue
	values := make([]data.Value, len(names))
	for i, name := range names {
		values[i] = data.NewStringValue(name)
	}

	return data.NewArrayValue(values), nil
}

func (f *ListConnectionsFunction) GetName() string {
	return "Database\\listConnections"
}

func (f *ListConnectionsFunction) GetParams() []data.GetValue {
	return []data.GetValue{}
}

func (f *ListConnectionsFunction) GetVariables() []data.Variable {
	return []data.Variable{}
}
