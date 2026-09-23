package http

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	httpfoundation "github.com/php-any/origami/std/symfony/http-foundation"
)

const (
	fqnIlluminateFile         = "Illuminate\\Http\\File"
	fqnIlluminateUploadedFile = "Illuminate\\Http\\UploadedFile"
)

// IlluminateFileClass 继承 Symfony File（Illuminate\Http\File）。
type IlluminateFileClass struct {
	node.Node
}

func NewIlluminateFileClass() data.ClassStmt {
	return &IlluminateFileClass{}
}

func (c *IlluminateFileClass) GetName() string { return fqnIlluminateFile }
func (c *IlluminateFileClass) GetExtend() *string {
	parent := httpfoundation.FqnFile
	return &parent
}
func (c *IlluminateFileClass) GetImplements() []string          { return nil }
func (c *IlluminateFileClass) GetProperty(string) (data.Property, bool) { return nil, false }
func (c *IlluminateFileClass) GetPropertyList() []data.Property { return nil }
func (c *IlluminateFileClass) GetConstruct() data.Method {
	parent := httpfoundation.NewFileClass()
	return parent.GetConstruct()
}
func (c *IlluminateFileClass) GetMethod(name string) (data.Method, bool) {
	return httpfoundation.NewFileClass().GetMethod(name)
}
func (c *IlluminateFileClass) GetMethods() []data.Method {
	return httpfoundation.NewFileClass().GetMethods()
}
func (c *IlluminateFileClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}

// IlluminateUploadedFileClass 继承 Symfony UploadedFile。
type IlluminateUploadedFileClass struct {
	node.Node
}

func NewIlluminateUploadedFileClass() data.ClassStmt {
	return &IlluminateUploadedFileClass{}
}

func (c *IlluminateUploadedFileClass) GetName() string { return fqnIlluminateUploadedFile }
func (c *IlluminateUploadedFileClass) GetExtend() *string {
	parent := httpfoundation.FqnUploadedFile
	return &parent
}
func (c *IlluminateUploadedFileClass) GetImplements() []string          { return nil }
func (c *IlluminateUploadedFileClass) GetProperty(string) (data.Property, bool) {
	return nil, false
}
func (c *IlluminateUploadedFileClass) GetPropertyList() []data.Property { return nil }
func (c *IlluminateUploadedFileClass) GetConstruct() data.Method {
	return httpfoundation.NewUploadedFileClass().GetConstruct()
}
func (c *IlluminateUploadedFileClass) GetMethod(name string) (data.Method, bool) {
	return httpfoundation.NewUploadedFileClass().GetMethod(name)
}
func (c *IlluminateUploadedFileClass) GetMethods() []data.Method {
	return httpfoundation.NewUploadedFileClass().GetMethods()
}
func (c *IlluminateUploadedFileClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx.CreateBaseContext()), nil
}
