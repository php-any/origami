package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// SplTempFileObjectClass 实现 PHP SplTempFileObject（extends SplFileObject�?
type SplTempFileObjectClass struct {
	node.Node
}

func NewSplTempFileObjectClass() *SplTempFileObjectClass { return &SplTempFileObjectClass{} }

func (c *SplTempFileObjectClass) GetName() string { return "SplTempFileObject" }

func (c *SplTempFileObjectClass) GetExtend() *string {
	parent := "SplFileObject"
	return &parent
}

func (c *SplTempFileObjectClass) GetImplements() []string {
	return []string{"RecursiveIterator", "SeekableIterator"}
}

func (c *SplTempFileObjectClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SplTempFileObjectClass) GetPropertyList() []data.Property              { return nil }

func (c *SplTempFileObjectClass) GetStaticProperty(name string) (data.Value, bool) {
	return SFOConstantValue(name)
}

func (c *SplTempFileObjectClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	cv := data.NewClassValue(c, ctx.CreateBaseContext())
	cv.SetProperty(sfiPathnameKey, data.NewStringValue(""))
	cv.SetProperty(sfoStateKey, &sfoStateValue{})
	return cv, nil
}

func (c *SplTempFileObjectClass) GetMethod(name string) (data.Method, bool) {
	switch name {
	case "__construct":
		return &STFOConstructMethod{}, true
	}
	return nil, false
}

func (c *SplTempFileObjectClass) GetMethods() []data.Method {
	return []data.Method{&STFOConstructMethod{}}
}

func (c *SplTempFileObjectClass) GetConstruct() data.Method { return &STFOConstructMethod{} }

type STFOConstructMethod struct{}

func (m *STFOConstructMethod) GetName() string            { return "__construct" }
func (m *STFOConstructMethod) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *STFOConstructMethod) GetIsStatic() bool          { return false }
func (m *STFOConstructMethod) GetReturnType() data.Types  { return nil }
func (m *STFOConstructMethod) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "tempFileName", 0, data.NewStringValue("php://temp"), data.NewBaseType("string")),
		node.NewParameter(nil, "mode", 1, data.NewStringValue("w+b"), data.NewBaseType("string")),
		node.NewParameter(nil, "flags", 2, data.NewIntValue(0), data.NewBaseType("int")),
	}
}
func (m *STFOConstructMethod) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "tempFileName", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "mode", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "flags", 2, data.NewBaseType("int")),
	}
}
func (m *STFOConstructMethod) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := sfiGetCV(ctx)
	if cv == nil {
		return nil, nil
	}
	tempName := sfoCtxString(ctx, 0, "php://temp")
	mode := sfoCtxString(ctx, 1, "w+b")
	flags := sfoCtxInt(ctx, 2, 0)
	if err := sfoOpenFileForTemp(cv, tempName, mode, flags); err != nil {
		return nil, utils.NewThrow(err)
	}
	return nil, nil
}
