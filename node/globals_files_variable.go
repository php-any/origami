package node

import "github.com/php-any/origami/data"

// $_FILES

type FilesVariable struct {
	*Node `pp:"-"`
}

var filesValue *data.ArrayValue

func NewFilesVariable(from data.From) data.Variable {
	return &FilesVariable{Node: NewNode(from)}
}

func (v *FilesVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return superglobalArray(ctx, "_FILES", func() *data.ArrayValue {
		_, files := multipartInput(ctx)
		return files

	}), nil
}

func (v *FilesVariable) GetIndex() int       { return -1 }
func (v *FilesVariable) GetName() string     { return "$_FILES" }
func (v *FilesVariable) GetType() data.Types { return nil }
func (v *FilesVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return setSuperglobalArray(ctx, "_FILES", value)
}

func (v *FilesVariable) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	if _, ctl := v.GetValue(ctx); ctl != nil {
		return nil, ctl
	}
	return ctx.GetVM().EnsureGlobalZVal("_FILES"), nil
}
func (v *FilesVariable) SuperglobalName() string { return "_FILES" }
