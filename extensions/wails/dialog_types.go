package wails

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/token"
)

// ============================================================================
// Wails\Dialog\FileFilter — 对话框文件过滤器
//
//	$filter = new Wails\Dialog\FileFilter("Image Files", "*.jpg;*.png");
// ============================================================================

type FileFilterClass struct{}

func NewFileFilterClass() data.ClassStmt { return &FileFilterClass{} }

func (c *FileFilterClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *FileFilterClass) GetFrom() data.From                            { return nil }
func (c *FileFilterClass) GetName() string                               { return "Wails\\Dialog\\FileFilter" }
func (c *FileFilterClass) GetExtend() *string                            { return nil }
func (c *FileFilterClass) GetImplements() []string                       { return nil }
func (c *FileFilterClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *FileFilterClass) GetPropertyList() []data.Property              { return nil }
func (c *FileFilterClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *FileFilterClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *FileFilterClass) GetMethods() []data.Method { return nil }
func (c *FileFilterClass) GetConstruct() data.Method { return &fileFilterConstruct{} }

type fileFilterConstruct struct{}

func (m *fileFilterConstruct) GetName() string            { return token.ConstructName }
func (m *fileFilterConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *fileFilterConstruct) GetIsStatic() bool          { return false }
func (m *fileFilterConstruct) GetReturnType() data.Types  { return nil }

func (m *fileFilterConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "displayName", 0, data.NewStringValue(""), data.NewBaseType("string")),
		node.NewParameter(nil, "pattern", 1, data.NewStringValue("*"), data.NewBaseType("string")),
	}
}
func (m *fileFilterConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "displayName", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "pattern", 1, data.NewBaseType("string")),
	}
}

func (m *fileFilterConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	if v, ok := ctx.GetIndexValue(0); ok {
		cv.SetProperty("DisplayName", data.NewStringValue(toString(v)))
	}
	if v, ok := ctx.GetIndexValue(1); ok {
		cv.SetProperty("Pattern", data.NewStringValue(toString(v)))
	}
	return nil, nil
}

// ============================================================================
// Wails\Dialog\OpenDialogOptions — 打开文件对话框选项
// ============================================================================

type OpenDialogOptionsClass struct{}

func NewOpenDialogOptionsClass() data.ClassStmt { return &OpenDialogOptionsClass{} }

func (c *OpenDialogOptionsClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *OpenDialogOptionsClass) GetFrom() data.From                            { return nil }
func (c *OpenDialogOptionsClass) GetName() string                               { return "Wails\\Dialog\\OpenDialogOptions" }
func (c *OpenDialogOptionsClass) GetExtend() *string                            { return nil }
func (c *OpenDialogOptionsClass) GetImplements() []string                       { return nil }
func (c *OpenDialogOptionsClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *OpenDialogOptionsClass) GetPropertyList() []data.Property              { return nil }
func (c *OpenDialogOptionsClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *OpenDialogOptionsClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *OpenDialogOptionsClass) GetMethods() []data.Method { return nil }
func (c *OpenDialogOptionsClass) GetConstruct() data.Method { return &openDialogOptionsConstruct{} }

type openDialogOptionsConstruct struct{}

func (m *openDialogOptionsConstruct) GetName() string            { return token.ConstructName }
func (m *openDialogOptionsConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *openDialogOptionsConstruct) GetIsStatic() bool          { return false }
func (m *openDialogOptionsConstruct) GetReturnType() data.Types  { return nil }

func (m *openDialogOptionsConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("array")),
	}
}
func (m *openDialogOptionsConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, data.NewBaseType("array")),
	}
}

func (m *openDialogOptionsConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	setDefaultStringProperty(cv, "DefaultDirectory", "")
	setDefaultStringProperty(cv, "DefaultFilename", "")
	setDefaultStringProperty(cv, "Title", "")
	setDefaultBoolProperty(cv, "ShowHiddenFiles", false)
	setDefaultBoolProperty(cv, "CanCreateDirectories", false)
	setDefaultBoolProperty(cv, "ResolvesAliases", false)
	setDefaultBoolProperty(cv, "TreatPackagesAsDirectories", false)
	cv.SetProperty("Filters", data.NewArrayValue(nil))

	if v, ok := ctx.GetIndexValue(0); ok {
		if av, ok := v.(*data.ArrayValue); ok {
			applyArrayToClassValue(cv, av, []string{
				"DefaultDirectory", "DefaultFilename", "Title",
				"ShowHiddenFiles", "CanCreateDirectories",
				"ResolvesAliases", "TreatPackagesAsDirectories",
			})
			if v, ok := arrayGet(av, "Filters"); ok {
				cv.SetProperty("Filters", v)
			}
		}
	}
	return nil, nil
}

// ============================================================================
// Wails\Dialog\SaveDialogOptions — 保存文件对话框选项
// ============================================================================

type SaveDialogOptionsClass struct{}

func NewSaveDialogOptionsClass() data.ClassStmt { return &SaveDialogOptionsClass{} }

func (c *SaveDialogOptionsClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *SaveDialogOptionsClass) GetFrom() data.From                            { return nil }
func (c *SaveDialogOptionsClass) GetName() string                               { return "Wails\\Dialog\\SaveDialogOptions" }
func (c *SaveDialogOptionsClass) GetExtend() *string                            { return nil }
func (c *SaveDialogOptionsClass) GetImplements() []string                       { return nil }
func (c *SaveDialogOptionsClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *SaveDialogOptionsClass) GetPropertyList() []data.Property              { return nil }
func (c *SaveDialogOptionsClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *SaveDialogOptionsClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *SaveDialogOptionsClass) GetMethods() []data.Method { return nil }
func (c *SaveDialogOptionsClass) GetConstruct() data.Method { return &saveDialogOptionsConstruct{} }

type saveDialogOptionsConstruct struct{}

func (m *saveDialogOptionsConstruct) GetName() string            { return token.ConstructName }
func (m *saveDialogOptionsConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *saveDialogOptionsConstruct) GetIsStatic() bool          { return false }
func (m *saveDialogOptionsConstruct) GetReturnType() data.Types  { return nil }

func (m *saveDialogOptionsConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("array")),
	}
}
func (m *saveDialogOptionsConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, data.NewBaseType("array")),
	}
}

func (m *saveDialogOptionsConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	setDefaultStringProperty(cv, "DefaultDirectory", "")
	setDefaultStringProperty(cv, "DefaultFilename", "")
	setDefaultStringProperty(cv, "Title", "")
	setDefaultBoolProperty(cv, "ShowHiddenFiles", false)
	setDefaultBoolProperty(cv, "CanCreateDirectories", false)
	setDefaultBoolProperty(cv, "TreatPackagesAsDirectories", false)
	cv.SetProperty("Filters", data.NewArrayValue(nil))

	if v, ok := ctx.GetIndexValue(0); ok {
		if av, ok := v.(*data.ArrayValue); ok {
			applyArrayToClassValue(cv, av, []string{
				"DefaultDirectory", "DefaultFilename", "Title",
				"ShowHiddenFiles", "CanCreateDirectories",
				"TreatPackagesAsDirectories",
			})
			if v, ok := arrayGet(av, "Filters"); ok {
				cv.SetProperty("Filters", v)
			}
		}
	}
	return nil, nil
}

// ============================================================================
// Wails\Dialog\MessageDialogOptions — 消息对话框选项
// ============================================================================

type MessageDialogOptionsClass struct{}

func NewMessageDialogOptionsClass() data.ClassStmt { return &MessageDialogOptionsClass{} }

func (c *MessageDialogOptionsClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *MessageDialogOptionsClass) GetFrom() data.From                            { return nil }
func (c *MessageDialogOptionsClass) GetName() string                               { return "Wails\\Dialog\\MessageDialogOptions" }
func (c *MessageDialogOptionsClass) GetExtend() *string                            { return nil }
func (c *MessageDialogOptionsClass) GetImplements() []string                       { return nil }
func (c *MessageDialogOptionsClass) GetProperty(name string) (data.Property, bool) { return nil, false }
func (c *MessageDialogOptionsClass) GetPropertyList() []data.Property              { return nil }
func (c *MessageDialogOptionsClass) GetMethod(name string) (data.Method, bool)     { return nil, false }
func (c *MessageDialogOptionsClass) GetStaticMethod(name string) (data.Method, bool) {
	return nil, false
}
func (c *MessageDialogOptionsClass) GetMethods() []data.Method { return nil }
func (c *MessageDialogOptionsClass) GetConstruct() data.Method { return &messageDialogOptionsConstruct{} }

type messageDialogOptionsConstruct struct{}

func (m *messageDialogOptionsConstruct) GetName() string            { return token.ConstructName }
func (m *messageDialogOptionsConstruct) GetModifier() data.Modifier { return data.ModifierPublic }
func (m *messageDialogOptionsConstruct) GetIsStatic() bool          { return false }
func (m *messageDialogOptionsConstruct) GetReturnType() data.Types  { return nil }

func (m *messageDialogOptionsConstruct) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "options", 0, data.NewArrayValue(nil), data.NewBaseType("array")),
	}
}
func (m *messageDialogOptionsConstruct) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "options", 0, data.NewBaseType("array")),
	}
}

func (m *messageDialogOptionsConstruct) Call(ctx data.Context) (data.GetValue, data.Control) {
	cv := getThis(ctx)
	if cv == nil {
		return nil, nil
	}
	setDefaultStringProperty(cv, "Type", "info")
	setDefaultStringProperty(cv, "Title", "")
	setDefaultStringProperty(cv, "Message", "")
	setDefaultStringProperty(cv, "DefaultButton", "")
	setDefaultStringProperty(cv, "CancelButton", "")
	cv.SetProperty("Buttons", data.NewArrayValue(nil))

	if v, ok := ctx.GetIndexValue(0); ok {
		if av, ok := v.(*data.ArrayValue); ok {
			applyArrayToClassValue(cv, av, []string{
				"Type", "Title", "Message",
				"DefaultButton", "CancelButton",
			})
			if v, ok := arrayGet(av, "Buttons"); ok {
				cv.SetProperty("Buttons", v)
			}
		}
	}
	return nil, nil
}
