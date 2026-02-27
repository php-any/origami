package node

import "github.com/php-any/origami/data"

// $_POST

type PostVariable struct {
	*Node `pp:"-"`
}

var postValue *data.ObjectValue

func NewPostVariable(from data.From) data.Variable {
	return &PostVariable{Node: NewNode(from)}
}

func (v *PostVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	if postValue == nil {
		postValue = data.NewObjectValue()
		if httpReq := getHTTPRequest(ctx); httpReq != nil {
			if httpReq.Form != nil {
				for key, values := range httpReq.Form {
					if len(values) > 0 {
						postValue.SetProperty(key, data.NewStringValue(values[0]))
					}
				}
			}
		}
	}
	return postValue, nil
}

func (v *PostVariable) GetIndex() int       { return 0 }
func (v *PostVariable) GetName() string     { return "$_POST" }
func (v *PostVariable) GetType() data.Types { return nil }
func (v *PostVariable) SetValue(ctx data.Context, value data.Value) data.Control {
	return data.NewErrorThrow(v.from, nil)
}
