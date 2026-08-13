package protowire

import (
	"github.com/php-any/origami/data"
)

// NewProtowireClass creates the Protowire standard library class.
// All methods are static, mirroring PHP's typical utility class pattern.
func NewProtowireClass() data.ClassStmt {
	return &ProtowireClass{
		parse:         NewParseMethod(),
		serialize:     NewSerializeMethod(),
		encodeTag:     NewEncodeTagMethod(),
		encodeVarint:  NewEncodeVarintMethod(),
		encodeBytes:   NewEncodeBytesMethod(),
		encodeFixed32: NewEncodeFixed32Method(),
		encodeFixed64: NewEncodeFixed64Method(),
	}
}

// ProtowireClass provides protobuf binary wire-format read and write
// operations. All methods are static.
//
// PHP namespace: Protowire
type ProtowireClass struct {
	parse         data.Method
	serialize     data.Method
	encodeTag     data.Method
	encodeVarint  data.Method
	encodeBytes   data.Method
	encodeFixed32 data.Method
	encodeFixed64 data.Method
}

func (c *ProtowireClass) GetValue(_ data.Context) (data.GetValue, data.Control) {
	clone := *c
	return &clone, nil
}

func (c *ProtowireClass) GetName() string {
	return "Protowire"
}

func (c *ProtowireClass) GetFrom() data.From {
	return nil
}

func (c *ProtowireClass) GetExtend() *string {
	return nil
}

func (c *ProtowireClass) GetImplements() []string {
	return nil
}

func (c *ProtowireClass) GetProperty(_ string) (data.Property, bool) {
	return nil, false
}

func (c *ProtowireClass) GetPropertyList() []data.Property {
	return []data.Property{}
}

func (c *ProtowireClass) GetMethod(name string) (data.Method, bool) {
	return c.GetStaticMethod(name)
}

func (c *ProtowireClass) GetStaticMethod(name string) (data.Method, bool) {
	switch name {
	case "parse":
		return c.parse, true
	case "serialize":
		return c.serialize, true
	case "encodeTag":
		return c.encodeTag, true
	case "encodeVarint":
		return c.encodeVarint, true
	case "encodeBytes":
		return c.encodeBytes, true
	case "encodeFixed32":
		return c.encodeFixed32, true
	case "encodeFixed64":
		return c.encodeFixed64, true
	}
	return nil, false
}

func (c *ProtowireClass) GetMethods() []data.Method {
	return []data.Method{
		c.parse,
		c.serialize,
		c.encodeTag,
		c.encodeVarint,
		c.encodeBytes,
		c.encodeFixed32,
		c.encodeFixed64,
	}
}

func (c *ProtowireClass) GetConstruct() data.Method {
	return nil
}
