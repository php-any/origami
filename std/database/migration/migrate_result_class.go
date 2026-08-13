package migration

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// NewMigrateResultClass 创建迁移结果类（用于类型声明与 IDE 提示）
func NewMigrateResultClass() data.ClassStmt {
	return &MigrateResultClass{
		propCreated:      node.NewProperty(nil, "created", "public", false, nil, data.NewBaseType("array")),
		propAltered:      node.NewProperty(nil, "altered", "public", false, nil, data.NewBaseType("array")),
		propTables:       node.NewProperty(nil, "tables", "public", false, nil, data.NewBaseType("array")),
		propCreatedCount: node.NewProperty(nil, "createdCount", "public", false, nil, data.NewBaseType("int")),
		propAlteredCount: node.NewProperty(nil, "alteredCount", "public", false, nil, data.NewBaseType("int")),
	}
}

// MigrateResultClass Database\Migration\MigrateResult
type MigrateResultClass struct {
	node.Node
	propCreated      data.Property
	propAltered      data.Property
	propTables       data.Property
	propCreatedCount data.Property
	propAlteredCount data.Property
}

func (s *MigrateResultClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *s
	return &clone, nil
}

func (s *MigrateResultClass) GetName() string         { return "Database\\Migration\\MigrateResult" }
func (s *MigrateResultClass) GetExtend() *string      { return nil }
func (s *MigrateResultClass) GetImplements() []string { return nil }
func (s *MigrateResultClass) AsString() string        { return "MigrateResult{}" }

func (s *MigrateResultClass) GetProperty(name string) (data.Property, bool) {
	switch name {
	case "created":
		return s.propCreated, true
	case "altered":
		return s.propAltered, true
	case "tables":
		return s.propTables, true
	case "createdCount":
		return s.propCreatedCount, true
	case "alteredCount":
		return s.propAlteredCount, true
	}
	return nil, false
}

func (s *MigrateResultClass) GetPropertyList() []data.Property {
	return []data.Property{
		s.propCreated,
		s.propAltered,
		s.propTables,
		s.propCreatedCount,
		s.propAlteredCount,
	}
}

func (s *MigrateResultClass) GetMethod(_ string) (data.Method, bool) { return nil, false }
func (s *MigrateResultClass) GetMethods() []data.Method              { return nil }
func (s *MigrateResultClass) GetConstruct() data.Method              { return nil }

func newMigrateResultValue(ctx data.Context, created, altered, tables []data.Value) data.GetValue {
	result := data.NewClassValue(NewMigrateResultClass(), ctx)
	result.SetProperty("created", data.NewArrayValue(created))
	result.SetProperty("altered", data.NewArrayValue(altered))
	result.SetProperty("tables", data.NewArrayValue(tables))
	result.SetProperty("createdCount", data.NewIntValue(len(created)))
	result.SetProperty("alteredCount", data.NewIntValue(len(altered)))
	return result
}
