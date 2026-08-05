package pdo

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// PdoMysqlClass 实现 PHP 8.4 的 Pdo\Mysql（extends PDO）。
// Laravel 13 config/database.php 通过 Pdo\Mysql::ATTR_* 引用驱动属性；
// Symfony polyfill 仅在 PHP_VERSION_ID < 80400 时提供同名类，故在此用 Go 注册原生类。
type PdoMysqlClass struct {
	node.Node
}

func (c *PdoMysqlClass) GetName() string {
	return "Pdo\\Mysql"
}

func (c *PdoMysqlClass) GetExtend() *string {
	ext := "PDO"
	return &ext
}

func (c *PdoMysqlClass) GetImplements() []string                    { return nil }
func (c *PdoMysqlClass) GetProperty(_ string) (data.Property, bool) { return nil, false }
func (c *PdoMysqlClass) GetPropertyList() []data.Property           { return nil }

func (c *PdoMysqlClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}

func (c *PdoMysqlClass) GetConstruct() data.Method {
	return &pdoConstructMethod{}
}

func (c *PdoMysqlClass) GetMethod(name string) (data.Method, bool) {
	// 方法走 PDO 继承链（call_object_method 会沿 GetExtend 查找）
	if name == "__construct" {
		return &pdoConstructMethod{}, true
	}
	return nil, false
}

func (c *PdoMysqlClass) GetMethods() []data.Method { return nil }

// GetStaticProperty 暴露 Pdo\Mysql::ATTR_*（对齐 PHP 8.4 / polyfill 常量名）
func (c *PdoMysqlClass) GetStaticProperty(name string) (data.Value, bool) {
	switch name {
	case "ATTR_COMPRESS":
		return data.NewIntValue(PDO_MYSQL_ATTR_COMPRESS), true
	case "ATTR_DIRECT_QUERY":
		return data.NewIntValue(PDO_MYSQL_ATTR_DIRECT_QUERY), true
	case "ATTR_FOUND_ROWS":
		return data.NewIntValue(PDO_MYSQL_ATTR_FOUND_ROWS), true
	case "ATTR_IGNORE_SPACE":
		return data.NewIntValue(PDO_MYSQL_ATTR_IGNORE_SPACE), true
	case "ATTR_INIT_COMMAND":
		return data.NewIntValue(PDO_MYSQL_ATTR_INIT_COMMAND), true
	case "ATTR_LOCAL_INFILE":
		return data.NewIntValue(PDO_MYSQL_ATTR_LOCAL_INFILE), true
	case "ATTR_LOCAL_INFILE_DIRECTORY":
		return data.NewIntValue(PDO_MYSQL_ATTR_LOCAL_INFILE_DIRECTORY), true
	case "ATTR_MAX_BUFFER_SIZE":
		return data.NewIntValue(PDO_MYSQL_ATTR_MAX_BUFFER_SIZE), true
	case "ATTR_MULTI_STATEMENTS":
		return data.NewIntValue(PDO_MYSQL_ATTR_MULTI_STATEMENTS), true
	case "ATTR_READ_DEFAULT_FILE":
		return data.NewIntValue(PDO_MYSQL_ATTR_READ_DEFAULT_FILE), true
	case "ATTR_READ_DEFAULT_GROUP":
		return data.NewIntValue(PDO_MYSQL_ATTR_READ_DEFAULT_GROUP), true
	case "ATTR_SERVER_PUBLIC_KEY":
		return data.NewIntValue(PDO_MYSQL_ATTR_SERVER_PUBLIC_KEY), true
	case "ATTR_SSL_CA":
		return data.NewIntValue(PDO_MYSQL_ATTR_SSL_CA), true
	case "ATTR_SSL_CAPATH":
		return data.NewIntValue(PDO_MYSQL_ATTR_SSL_CAPATH), true
	case "ATTR_SSL_CERT":
		return data.NewIntValue(PDO_MYSQL_ATTR_SSL_CERT), true
	case "ATTR_SSL_CIPHER":
		return data.NewIntValue(PDO_MYSQL_ATTR_SSL_CIPHER), true
	case "ATTR_SSL_KEY":
		return data.NewIntValue(PDO_MYSQL_ATTR_SSL_KEY), true
	case "ATTR_SSL_VERIFY_SERVER_CERT":
		return data.NewIntValue(PDO_MYSQL_ATTR_SSL_VERIFY_SERVER_CERT), true
	case "ATTR_USE_BUFFERED_QUERY":
		return data.NewIntValue(PDO_MYSQL_ATTR_USE_BUFFERED_QUERY), true
	}
	return nil, false
}
