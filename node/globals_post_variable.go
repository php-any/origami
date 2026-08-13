package node

import (
	"strings"

	"github.com/php-any/origami/data"
)

// $_POST

type PostVariable struct {
	*Node `pp:"-"`
}

var postValue *data.ObjectValue

func NewPostVariable(from data.From) data.Variable {
	return &PostVariable{Node: NewNode(from)}
}

func (v *PostVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if httpReq := getHTTPRequest(ctx); httpReq != nil {
		// Go 的 Form/PostForm 需先 ParseForm；未解析时 Form 为 nil，$_POST 会一直为空。
		ct := httpReq.Header.Get("Content-Type")
		if strings.HasPrefix(ct, "multipart/form-data") {
			_ = httpReq.ParseMultipartForm(32 << 20)
		} else {
			_ = httpReq.ParseForm()
		}
		obj := data.NewObjectValue()
		// 对齐 PHP：$_POST 只含请求体，不含 query string（用 PostForm 而非 Form）。
		for key, values := range httpReq.PostForm {
			if len(values) > 0 {
				obj.SetProperty(key, data.NewStringValue(values[0]))
			}
		}
		return obj, nil
	}
	if postValue == nil {
		postValue = data.NewObjectValue()
	}
	return postValue, nil
}

func (v *PostVariable) GetIndex() int       { return 0 }
func (v *PostVariable) GetName() string     { return "$_POST" }
func (v *PostVariable) GetType() data.Types { return nil }
func (v *PostVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return data.NewErrorThrow(v.from, nil)
}
