package node

import (
	"bytes"
	"io"
	"strings"

	"github.com/php-any/origami/data"
)

// $_POST

type PostVariable struct {
	*Node `pp:"-"`
}

var postValue *data.ArrayValue

func NewPostVariable(from data.From) data.Variable {
	return &PostVariable{Node: NewNode(from)}
}

func (v *PostVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return superglobalArray(ctx, "_POST", func() *data.ArrayValue {
		if httpReq := getHTTPRequest(ctx); httpReq != nil {
			if httpReq.Method != "POST" {
				return data.NewArrayValueFromSlots(nil)
			}
			// Go 的 Form/PostForm 需先 ParseForm；未解析时 Form 为 nil，$_POST 会一直为空。
			ct := httpReq.Header.Get("Content-Type")
			if strings.HasPrefix(ct, "multipart/form-data") {
				post, _ := multipartInput(ctx)
				return post
			} else if strings.HasPrefix(ct, "application/x-www-form-urlencoded") {
				if httpReq.Body != nil {
					body, err := io.ReadAll(httpReq.Body)
					if err == nil {
						httpReq.Body = io.NopCloser(bytes.NewReader(body))
						return data.ParseFormFields(string(body))
					}
				}
			}
			obj := data.NewArrayValueFromSlots(nil)
			// 对齐 PHP：$_POST 只含请求体，不含 query string（用 PostForm 而非 Form）。
			for key, values := range httpReq.PostForm {
				if len(values) > 0 {
					for _, value := range values {
						data.SetFormField(obj, key, data.NewStringValue(value))
					}
				}
			}
			return obj
		}
		return data.NewArrayValueFromSlots(nil)

	}), nil
}

func (v *PostVariable) GetIndex() int       { return -1 }
func (v *PostVariable) GetName() string     { return "$_POST" }
func (v *PostVariable) GetType() data.Types { return nil }
func (v *PostVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return setSuperglobalArray(ctx, "_POST", value)
}

func (v *PostVariable) GetZVal(ctx data.Context) (*data.ZVal, data.Control) {
	if _, ctl := v.GetValue(ctx); ctl != nil {
		return nil, ctl
	}
	return ctx.GetVM().EnsureGlobalZVal("_POST"), nil
}
func (v *PostVariable) SuperglobalName() string { return "_POST" }
