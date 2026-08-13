package node

import "github.com/php-any/origami/data"

// $_FILES

type FilesVariable struct {
	*Node `pp:"-"`
}

var filesValue *data.ObjectValue

func NewFilesVariable(from data.From) data.Variable {
	return &FilesVariable{Node: NewNode(from)}
}

func (v *FilesVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if filesValue == nil {
		filesValue = data.NewObjectValue()
		if httpReq := getHTTPRequest(ctx); httpReq != nil {
			if httpReq.MultipartForm != nil {
				for key, fileHeaders := range httpReq.MultipartForm.File {
					fileArr := make([]data.Value, 0, len(fileHeaders))
					for _, fh := range fileHeaders {
						fileInfo := data.NewObjectValue()
						fileInfo.SetProperty("name", data.NewStringValue(fh.Filename))
						fileInfo.SetProperty("type", data.NewStringValue(fh.Header.Get("Content-Type")))
						fileInfo.SetProperty("tmp_name", data.NewStringValue(""))
						fileInfo.SetProperty("error", data.NewIntValue(0))
						fileInfo.SetProperty("size", data.NewIntValue(int(fh.Size)))
						fileArr = append(fileArr, fileInfo)
					}
					filesValue.SetProperty(key, data.NewArrayValue(fileArr))
				}
			}
		}
	}
	return filesValue, nil
}

func (v *FilesVariable) GetIndex() int       { return 0 }
func (v *FilesVariable) GetName() string     { return "$_FILES" }
func (v *FilesVariable) GetType() data.Types { return nil }
func (v *FilesVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return data.NewErrorThrow(v.from, nil)
}
