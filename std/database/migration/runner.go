package migration

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// ConnectionResolver 按名称解析数据库连接
type ConnectionResolver func(connectionName string) (*sql.DB, bool)

var connectionResolver ConnectionResolver

// SetConnectionResolver 注入连接解析器（由 database 包在 Load 时设置）
func SetConnectionResolver(resolver ConnectionResolver) {
	connectionResolver = resolver
}

type classLister interface {
	AllClasses() []data.ClassStmt
}

// Migrate 根据模型目录中的 Entity 注解同步数据库 Schema
func Migrate(ctx data.Context, connection data.GetValue, modelDir string) (data.GetValue, data.Control) {
	conn, ctl := resolveConnectionArg(connection)
	if ctl != nil {
		return nil, ctl
	}
	if conn == nil {
		return nil, utils.NewThrow(fmt.Errorf("无效的数据库连接"))
	}

	modelDir = strings.TrimSpace(modelDir)
	if modelDir == "" {
		return nil, utils.NewThrow(fmt.Errorf("模型目录不能为空"))
	}
	absDir, err := filepath.Abs(modelDir)
	if err != nil {
		return nil, utils.NewThrow(fmt.Errorf("模型目录无效: %w", err))
	}
	if st, err := os.Stat(absDir); err != nil || !st.IsDir() {
		return nil, utils.NewThrow(fmt.Errorf("模型目录不存在: %s", modelDir))
	}

	files, err := collectPhpFiles(absDir)
	if err != nil {
		return nil, utils.NewThrow(fmt.Errorf("扫描模型目录失败: %w", err))
	}

	vm := ctx.GetVM()
	for _, file := range files {
		if _, acl := vm.LoadAndRun(file); acl != nil {
			return nil, acl
		}
	}

	schemas, ctl := collectTableSchemas(vm, files)
	if ctl != nil {
		return nil, ctl
	}
	if len(schemas) == 0 {
		return nil, utils.NewThrow(fmt.Errorf("模型目录中未找到带 #[Table] 注解的实体类: %s", modelDir))
	}

	dialect, err := detectDialect(conn)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	tx, err := conn.Begin()
	if err != nil {
		return nil, utils.NewThrow(fmt.Errorf("开启迁移事务失败: %w", err))
	}
	defer tx.Rollback()

	created, altered, err := applyMigrations(tx, dialect, schemas)
	if err != nil {
		return nil, utils.NewThrow(err)
	}

	if err := tx.Commit(); err != nil {
		return nil, utils.NewThrow(fmt.Errorf("提交迁移事务失败: %w", err))
	}

	tables := make([]data.Value, len(schemas))
	for i, schema := range schemas {
		obj := data.NewObjectValue()
		obj.SetProperty("table", data.NewStringValue(schema.TableName))
		obj.SetProperty("class", data.NewStringValue(schema.ClassName))
		tables[i] = obj
	}

	createdValues := migrationItemsToValues(created, false)
	alteredValues := migrationItemsToValues(altered, true)

	return newMigrateResultValue(ctx, createdValues, alteredValues, tables), nil
}

type migrationItem struct {
	table  string
	column string
	class  string
}

func applyMigrations(exec execQuerier, dialect Dialect, schemas []*TableSchema) (created, altered []migrationItem, err error) {
	created = make([]migrationItem, 0)
	altered = make([]migrationItem, 0)

	for _, schema := range schemas {
		exists, err := dialect.TableExists(exec, schema.TableName)
		if err != nil {
			return nil, nil, fmt.Errorf("检查表 %s 失败: %w", schema.TableName, err)
		}

		if !exists {
			sqlStr := dialect.BuildCreateTableSQL(schema)
			if _, err := exec.Exec(sqlStr); err != nil {
				return nil, nil, fmt.Errorf("创建表 %s 失败: %w", schema.TableName, err)
			}
			created = append(created, migrationItem{table: schema.TableName, class: schema.ClassName})
			continue
		}

		existingCols, err := dialect.TableColumns(exec, schema.TableName)
		if err != nil {
			return nil, nil, fmt.Errorf("读取表 %s 结构失败: %w", schema.TableName, err)
		}

		for _, col := range schema.Columns {
			if existingCols[col.Name] {
				continue
			}
			sqlStr := dialect.BuildAddColumnSQL(schema.TableName, col)
			if _, err := exec.Exec(sqlStr); err != nil {
				return nil, nil, fmt.Errorf("表 %s 添加列 %s 失败: %w", schema.TableName, col.Name, err)
			}
			altered = append(altered, migrationItem{
				table:  schema.TableName,
				column: col.Name,
				class:  schema.ClassName,
			})
		}
	}
	return created, altered, nil
}

func migrationItemsToValues(items []migrationItem, withColumn bool) []data.Value {
	out := make([]data.Value, 0, len(items))
	for _, item := range items {
		obj := data.NewObjectValue()
		obj.SetProperty("table", data.NewStringValue(item.table))
		obj.SetProperty("class", data.NewStringValue(item.class))
		if withColumn {
			obj.SetProperty("column", data.NewStringValue(item.column))
		}
		out = append(out, obj)
	}
	return out
}

func resolveConnectionArg(connection data.GetValue) (*sql.DB, data.Control) {
	if connection == nil {
		if connectionResolver == nil {
			return nil, utils.NewThrow(fmt.Errorf("迁移模块未初始化连接解析器"))
		}
		conn, ok := connectionResolver("default")
		if !ok || conn == nil {
			return nil, utils.NewThrow(fmt.Errorf("数据库连接不可用: default"))
		}
		return conn, nil
	}

	if _, isNull := connection.(*data.NullValue); isNull {
		if connectionResolver == nil {
			return nil, utils.NewThrow(fmt.Errorf("迁移模块未初始化连接解析器"))
		}
		conn, ok := connectionResolver("default")
		if !ok || conn == nil {
			return nil, utils.NewThrow(fmt.Errorf("数据库连接不可用: default"))
		}
		return conn, nil
	}

	if conn, ok := sqlDBFromValue(connection); ok {
		return conn, nil
	}

	if connectionResolver != nil {
		if name, ok := connection.(data.AsString); ok {
			connName := strings.TrimSpace(name.AsString())
			if connName == "" {
				connName = "default"
			}
			conn, ok := connectionResolver(connName)
			if ok && conn != nil {
				return conn, nil
			}
			return nil, utils.NewThrow(fmt.Errorf("数据库连接不可用: %s", connName))
		}
	}

	return nil, utils.NewThrow(fmt.Errorf("无效的数据库连接参数"))
}

func collectTableSchemas(vm data.VM, modelFiles []string) ([]*TableSchema, data.Control) {
	lister, ok := vm.(classLister)
	if !ok {
		return nil, utils.NewThrow(fmt.Errorf("运行时无法枚举实体类"))
	}

	fileSet := make(map[string]struct{}, len(modelFiles))
	for _, f := range modelFiles {
		fileSet[normalizePath(f)] = struct{}{}
	}

	var schemas []*TableSchema
	seen := make(map[string]struct{})
	for _, cls := range lister.AllClasses() {
		cs, ok := cls.(*node.ClassStatement)
		if !ok || cs == nil {
			continue
		}
		source := classSourcePath(cs)
		if source == "" || !pathInSet(source, fileSet) {
			continue
		}
		schema, err := parseTableSchema(cs)
		if err != nil {
			continue
		}
		if _, dup := seen[schema.TableName]; dup {
			return nil, utils.NewThrow(fmt.Errorf("重复的表名 %s", schema.TableName))
		}
		seen[schema.TableName] = struct{}{}
		schemas = append(schemas, schema)
	}
	return schemas, nil
}

func classSourcePath(cls *node.ClassStatement) string {
	if cls == nil {
		return ""
	}
	from := cls.GetFrom()
	if from == nil {
		return ""
	}
	return normalizePath(from.GetSource())
}

func pathInSet(path string, fileSet map[string]struct{}) bool {
	if _, ok := fileSet[path]; ok {
		return true
	}
	for f := range fileSet {
		if strings.EqualFold(f, path) {
			return true
		}
	}
	return false
}
