package fyne

import "github.com/php-any/origami/data"

// ThemeIconNameClass 是 Fyne\Theme\IconName 枚举类
type ThemeIconNameClass struct{}

func NewThemeIconNameClass() data.ClassStmt { return &ThemeIconNameClass{} }

func (c *ThemeIconNameClass) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return data.NewClassValue(c, ctx), nil
}
func (c *ThemeIconNameClass) GetFrom() data.From                              { return nil }
func (c *ThemeIconNameClass) GetName() string                                 { return "Fyne\\Theme\\IconName" }
func (c *ThemeIconNameClass) GetExtend() *string                              { return nil }
func (c *ThemeIconNameClass) GetImplements() []string                         { return nil }
func (c *ThemeIconNameClass) GetProperty(name string) (data.Property, bool)   { return nil, false }
func (c *ThemeIconNameClass) GetPropertyList() []data.Property                { return nil }
func (c *ThemeIconNameClass) GetMethod(name string) (data.Method, bool)       { return nil, false }
func (c *ThemeIconNameClass) GetStaticMethod(name string) (data.Method, bool) { return nil, false }
func (c *ThemeIconNameClass) GetMethods() []data.Method                       { return nil }
func (c *ThemeIconNameClass) GetConstruct() data.Method                       { return nil }

var iconNames = map[string]string{
	"CANCEL":           "cancel",
	"CONFIRM":          "confirm",
	"DELETE":           "delete",
	"SEARCH":           "search",
	"SETTINGS":         "settings",
	"MORE_HORIZONTAL":  "moreHorizontal",
	"HOME":             "home",
	"HELP":             "help",
	"INFO":             "info",
	"WARNING":          "warning",
	"ERROR":            "error",
	"QUESTION":         "question",
	"DOCUMENT":         "document",
	"DOCUMENT_CREATE":  "documentCreate",
	"DOCUMENT_PRINT":   "documentPrint",
	"DOCUMENT_SAVE":    "documentSave",
	"FOLDER":           "folder",
	"FOLDER_NEW":       "folderNew",
	"FOLDER_OPEN":      "folderOpen",
	"VIEW_REFRESH":     "viewRefresh",
	"VIEW_FULLSCREEN":  "viewFullScreen",
	"VIEW_ZOOM_FIT":    "viewZoomFit",
	"VIEW_ZOOM_IN":     "viewZoomIn",
	"VIEW_ZOOM_OUT":    "viewZoomOut",
	"NAVIGATE_BACK":    "navigateBack",
	"NAVIGATE_FORWARD": "navigateForward",
	"CONTENT_COPY":     "contentCopy",
	"CONTENT_CUT":      "contentCut",
	"CONTENT_PASTE":    "contentPaste",
	"CONTENT_UNDO":     "contentUndo",
	"CONTENT_REDO":     "contentRedo",
}

func (c *ThemeIconNameClass) GetStaticProperty(name string) (data.Value, bool) {
	if v, ok := iconNames[name]; ok {
		return data.NewStringValue(v), true
	}
	return nil, false
}
